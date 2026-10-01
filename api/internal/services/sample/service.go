package sample

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

const tableName = "sample_items"

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
	query := selectQuery + `WHERE s.id = $1`
	if !includeDeleted {
		query = query + ` AND s.deleted_at IS NULL`
	}
	var r row
	err := s.db.GetContext(ctx, &r, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return Dto{}, err
	} else if err != nil {
		return Dto{}, fmt.Errorf("unable to get sample item. err: %w", err)
	}

	return r.toDto(), nil
}

// Metadata returns the sample_items table's columns and Postgres data types,
// so a frontend can discover valid sort/filter column names without
// hardcoding them.
func (s *Service) Metadata(ctx context.Context) ([]dto.ColumnMeta, error) {
	columns, err := util.TableColumns(ctx, s.db, tableName)
	if err != nil {
		return nil, fmt.Errorf("unable to load sample_items table metadata. err: %w", err)
	}

	return columns, nil
}

func (s *Service) Create(ctx context.Context, input CreateInputDto, loggedInUserId uuid.UUID) (uuid.UUID, error) {
	newId := uuid.New()

	op := util.Operation{Kind: "sample.create", PerformedBy: loggedInUserId, TargetTable: tableName, TargetID: newId}
	err := util.RunOperation(ctx, s.db, op, func(ctx context.Context, tx *sqlx.Tx) error {
		return util.Insert(ctx, tx, tableName, newId, map[string]any{
			"id":          newId,
			"name":        input.Name,
			"description": input.Description,
			"created_at":  time.Now(),
			"created_by":  loggedInUserId,
		})
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("unable to create new sample item. err: %w", err)
	}

	return newId, nil
}

// List returns a page of sample items matching filter's search/sort/pagination.
// The heavy lifting is shared across every table's plain list endpoint - see
// util.List.
func (s *Service) List(ctx context.Context, filter dto.Filters) (*dto.PaginationResponse[Dto], error) {
	return util.List(ctx, s.db, tableName, "s", selectQuery, []string{"name", "description"}, filter, row.toDto)
}

// Query lists sample items matching a dynamic WHERE clause built by the
// frontend's react-querybuilder. The heavy lifting (column validation,
// alias-qualifying, parameter binding, pagination) is shared across every
// table's "/query" endpoint - see util.AdvancedQuery.
func (s *Service) Query(ctx context.Context, filter dto.Filters, whereCondition string, whereConditionParams []any) (*dto.PaginationResponse[Dto], error) {
	return util.AdvancedQuery(ctx, s.db, tableName, "s", selectQuery, filter, whereCondition, whereConditionParams, row.toDto)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateInputDto, loggedInUserId uuid.UUID) error {
	op := util.Operation{Kind: "sample.update", PerformedBy: loggedInUserId, TargetTable: tableName, TargetID: id}
	err := util.RunOperation(ctx, s.db, op, func(ctx context.Context, tx *sqlx.Tx) error {
		return util.UpdateByID(ctx, tx, tableName, id, map[string]any{
			"id":          id,
			"name":        input.Name,
			"description": input.Description,
			"updated_at":  time.Now(),
			"updated_by":  loggedInUserId,
		})
	})
	if err != nil {
		return fmt.Errorf("unable to update sample item. err: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID, loggedInUserId uuid.UUID) error {
	op := util.Operation{Kind: "sample.delete", PerformedBy: loggedInUserId, TargetTable: tableName, TargetID: id}
	err := util.RunOperation(ctx, s.db, op, func(ctx context.Context, tx *sqlx.Tx) error {
		return util.UpdateByID(ctx, tx, tableName, id, map[string]any{
			"id":         id,
			"deleted_at": time.Now(),
			"deleted_by": loggedInUserId,
		})
	})
	if err != nil {
		return fmt.Errorf("unable to delete sample item. err: %w", err)
	}

	return nil
}
