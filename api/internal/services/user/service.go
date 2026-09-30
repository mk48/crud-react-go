package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"kfamily/internal/dto"
	"kfamily/internal/util"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// tableName is double-quoted because "user" is a reserved word in Postgres.
const tableName = `"user"`

// ErrLastAdmin is returned by Update/Delete when the change would leave no
// live admin.
var ErrLastAdmin = errors.New("can't remove the last remaining admin")

type (
	Service struct {
		db *sqlx.DB
	}
)

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

//------------------------------------------------------------------------------

func (s *Service) GetOne(ctx context.Context, id uuid.UUID, includeDeleted bool) (Dto, error) {
	query := selectQuery + `WHERE u.id = $1`
	if !includeDeleted {
		query = query + ` AND u.deleted_at IS NULL`
	}
	var r row
	err := s.db.GetContext(ctx, &r, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return Dto{}, err
	} else if err != nil {
		return Dto{}, fmt.Errorf("unable to get user. err: %w", err)
	}

	return r.toDto(), nil
}

// Metadata returns the user table's columns and Postgres data types, so a
// frontend can discover valid sort/filter column names without hardcoding them.
func (s *Service) Metadata(ctx context.Context) ([]dto.ColumnMeta, error) {
	columns, err := util.TableColumns(ctx, s.db, tableName)
	if err != nil {
		return nil, fmt.Errorf("unable to load user table metadata. err: %w", err)
	}

	return columns, nil
}

// List returns a page of users matching filter's search/sort/pagination.
// The heavy lifting is shared across every table's plain list endpoint - see
// util.List. There is no Create - users are provisioned automatically by
// AuthMiddleware the first time they sign in.
func (s *Service) List(ctx context.Context, filter dto.Filters) (*dto.PaginationResponse[Dto], error) {
	return util.List(ctx, s.db, tableName, "u", selectQuery, []string{"name", "email"}, filter, row.toDto)
}

// Query lists users matching a dynamic WHERE clause built by the frontend's
// react-querybuilder (see AdvancedQueryBuilder.tsx). The heavy lifting
// (column validation, alias-qualifying, parameter binding, pagination) is
// shared across every table's "/query" endpoint - see util.AdvancedQuery.
func (s *Service) Query(ctx context.Context, filter dto.Filters, whereCondition string, whereConditionParams []any) (*dto.PaginationResponse[Dto], error) {
	return util.AdvancedQuery(ctx, s.db, tableName, "u", selectQuery, filter, whereCondition, whereConditionParams, row.toDto)
}

// Update only touches name and is_admin - sub and email are immutable, set
// once by AuthMiddleware at first sign-in. is_admin is left unchanged when
// input.IsAdmin is nil.
func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateInputDto, loggedInUserId uuid.UUID) error {
	params := map[string]any{
		"id":         id,
		"name":       input.Name,
		"updated_at": time.Now(),
		"updated_by": loggedInUserId,
	}
	if input.IsAdmin != nil {
		params["is_admin"] = *input.IsAdmin
	}
	demoting := input.IsAdmin != nil && !*input.IsAdmin

	err := util.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		if demoting {
			if err := guardLastAdmin(ctx, tx, id); err != nil {
				return err
			}
		}
		return util.UpdateByID(ctx, tx, tableName, id, params)
	})
	if err != nil {
		return fmt.Errorf("unable to update user. err: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID, loggedInUserId uuid.UUID) error {
	err := util.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		if err := guardLastAdmin(ctx, tx, id); err != nil {
			return err
		}
		return util.UpdateByID(ctx, tx, tableName, id, map[string]any{
			"id":         id,
			"deleted_at": time.Now(),
			"deleted_by": loggedInUserId,
		})
	})
	if err != nil {
		return fmt.Errorf("unable to delete user. err: %w", err)
	}

	return nil
}

// guardLastAdmin returns ErrLastAdmin if id is the only live admin. It locks
// the live admin rows (FOR UPDATE) until tx ends, so two admins demoting or
// deleting each other at the same moment can't both pass the check and
// leave no admin at all.
func guardLastAdmin(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) error {
	adminIds := []uuid.UUID{}
	query := `SELECT id FROM "user" WHERE is_admin AND deleted_at IS NULL FOR UPDATE`
	if err := tx.SelectContext(ctx, &adminIds, query); err != nil {
		return fmt.Errorf("unable to load admins. err: %w", err)
	}

	if len(adminIds) == 1 && adminIds[0] == id {
		return ErrLastAdmin
	}

	return nil
}
