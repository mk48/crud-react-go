package samplechild

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

const tableName = "sample_child_items"

// ErrInvalidSampleItem is returned by Create/Update when sampleItemId isn't
// a valid UUID or doesn't reference an existing, non-deleted sample item -
// checked up front so the handler can answer 400 instead of surfacing the
// foreign key violation as a 500.
var ErrInvalidSampleItem = errors.New("sampleItemId must reference an existing sample item")

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
	query := selectQuery + `WHERE sc.id = $1`
	if !includeDeleted {
		query = query + ` AND sc.deleted_at IS NULL`
	}
	var r row
	err := s.db.GetContext(ctx, &r, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return Dto{}, err
	} else if err != nil {
		return Dto{}, fmt.Errorf("unable to get sample child item. err: %w", err)
	}

	return r.toDto(), nil
}

// Metadata returns the sample_child_items table's columns and Postgres data
// types, so a frontend can discover valid sort/filter column names without
// hardcoding them.
func (s *Service) Metadata(ctx context.Context) ([]dto.ColumnMeta, error) {
	columns, err := util.TableColumns(ctx, s.db, tableName)
	if err != nil {
		return nil, fmt.Errorf("unable to load sample_child_items table metadata. err: %w", err)
	}

	return columns, nil
}

func (s *Service) Create(ctx context.Context, input CreateInputDto, loggedInUserId uuid.UUID) (uuid.UUID, error) {
	sampleItemId, err := s.validSampleItemId(ctx, input.SampleItemId)
	if err != nil {
		return uuid.UUID{}, err
	}

	newId := uuid.New()

	op := util.Operation{Kind: "sample_child.create", PerformedBy: loggedInUserId, TargetTable: tableName, TargetID: newId}
	err = util.RunOperation(ctx, s.db, op, func(ctx context.Context, tx *sqlx.Tx) error {
		return util.Insert(ctx, tx, tableName, newId, map[string]any{
			"id":             newId,
			"sample_item_id": sampleItemId,
			"name":           input.Name,
			"created_at":     time.Now(),
			"created_by":     loggedInUserId,
		})
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("unable to create new sample child item. err: %w", err)
	}

	return newId, nil
}

// List returns a page of sample child items matching filter's
// search/sort/pagination. See util.List.
func (s *Service) List(ctx context.Context, filter dto.Filters) (*dto.PaginationResponse[Dto], error) {
	return util.List(ctx, s.db, tableName, "sc", selectQuery, []string{"name"}, filter, row.toDto)
}

// Query lists sample child items matching a dynamic WHERE clause built by the
// frontend's react-querybuilder. See util.AdvancedQuery.
func (s *Service) Query(ctx context.Context, filter dto.Filters, whereCondition string, whereConditionParams []any) (*dto.PaginationResponse[Dto], error) {
	return util.AdvancedQuery(ctx, s.db, tableName, "sc", selectQuery, filter, whereCondition, whereConditionParams, row.toDto)
}

func (s *Service) Update(ctx context.Context, id uuid.UUID, input UpdateInputDto, loggedInUserId uuid.UUID) error {
	sampleItemId, err := s.validSampleItemId(ctx, input.SampleItemId)
	if err != nil {
		return err
	}

	op := util.Operation{Kind: "sample_child.update", PerformedBy: loggedInUserId, TargetTable: tableName, TargetID: id}
	err = util.RunOperation(ctx, s.db, op, func(ctx context.Context, tx *sqlx.Tx) error {
		return util.UpdateByID(ctx, tx, tableName, id, map[string]any{
			"id":             id,
			"sample_item_id": sampleItemId,
			"name":           input.Name,
			"updated_at":     time.Now(),
			"updated_by":     loggedInUserId,
		})
	})
	if err != nil {
		return fmt.Errorf("unable to update sample child item. err: %w", err)
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID, loggedInUserId uuid.UUID) error {
	op := util.Operation{Kind: "sample_child.delete", PerformedBy: loggedInUserId, TargetTable: tableName, TargetID: id}
	err := util.RunOperation(ctx, s.db, op, func(ctx context.Context, tx *sqlx.Tx) error {
		return util.UpdateByID(ctx, tx, tableName, id, map[string]any{
			"id":         id,
			"deleted_at": time.Now(),
			"deleted_by": loggedInUserId,
		})
	})
	if err != nil {
		return fmt.Errorf("unable to delete sample child item. err: %w", err)
	}

	return nil
}

// validSampleItemId parses raw and checks it references a non-deleted
// sample item, returning ErrInvalidSampleItem otherwise. The foreign key
// alone would still accept a soft-deleted parent.
func (s *Service) validSampleItemId(ctx context.Context, raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.UUID{}, ErrInvalidSampleItem
	}

	var exists bool
	err = s.db.GetContext(ctx, &exists, `SELECT EXISTS (SELECT 1 FROM sample_items WHERE id = $1 AND deleted_at IS NULL)`, id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("unable to check sample item. err: %w", err)
	}
	if !exists {
		return uuid.UUID{}, ErrInvalidSampleItem
	}

	return id, nil
}
