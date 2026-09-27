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
		query = query + ` AND u.deleted_by IS NULL`
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
// once by AuthMiddleware at first sign-in.
func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateInputDto, loggedInUserId uuid.UUID) error {
	err := util.UpdateByID(ctx, s.db, tableName, id, map[string]any{
		"id":         id,
		"name":       input.Name,
		"is_admin":   input.IsAdmin,
		"updated_at": time.Now(),
		"updated_by": loggedInUserId,
	})
	if err != nil {
		return fmt.Errorf("unable to update user. err: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID, loggedInUserId uuid.UUID) error {
	err := util.UpdateByID(ctx, s.db, tableName, id, map[string]any{
		"id":         id,
		"deleted_at": time.Now(),
		"deleted_by": loggedInUserId,
	})
	if err != nil {
		return fmt.Errorf("unable to delete user. err: %w", err)
	}

	return nil
}
