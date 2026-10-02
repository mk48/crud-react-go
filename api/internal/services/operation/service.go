package operation

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"kfamily/internal/dto"
	"kfamily/internal/util"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const tableName = "operation"

// maxChanges caps how many changes GetOne returns - an import can touch
// thousands of rows, more than a page can usefully show.
const maxChanges = 1000

type Service struct {
	db *sqlx.DB
}

func NewService(db *sqlx.DB) *Service {
	return &Service{db: db}
}

//------------------------------------------------------------------------------

// GetOne returns operation id with (up to maxChanges of) the record changes
// it caused.
func (s *Service) GetOne(ctx context.Context, id uuid.UUID) (DetailDto, error) {
	var r row
	err := s.db.GetContext(ctx, &r, selectQuery+`WHERE o.id = $1`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return DetailDto{}, err
	} else if err != nil {
		return DetailDto{}, fmt.Errorf("unable to get operation. err: %w", err)
	}

	changeRows := []changeRow{}
	if err := s.db.SelectContext(ctx, &changeRows, changesQuery, id, maxChanges); err != nil {
		return DetailDto{}, fmt.Errorf("unable to load operation changes. err: %w", err)
	}

	changes := make([]ChangeDto, len(changeRows))
	for i, c := range changeRows {
		changes[i] = c.toDto()
	}

	return DetailDto{Dto: r.toDto(), Changes: changes}, nil
}

// Metadata returns the operation table's columns and Postgres data types,
// for the advanced query builder.
func (s *Service) Metadata(ctx context.Context) ([]dto.ColumnMeta, error) {
	columns, err := util.TableColumns(ctx, s.db, tableName)
	if err != nil {
		return nil, fmt.Errorf("unable to load operation table metadata. err: %w", err)
	}

	return columns, nil
}

// List returns a page of operations matching filter's search (over kind,
// target table and client), sort and pagination. See util.List.
func (s *Service) List(ctx context.Context, filter dto.Filters) (*dto.PaginationResponse[Dto], error) {
	// Operations are never deleted - there's no deleted_at for util.List to
	// filter on.
	filter.IncludeDeletedRecords = true
	return util.List(ctx, s.db, tableName, "o", selectQuery, []string{"kind", "target_table", "client"}, filter, row.toDto)
}

// Query lists operations matching a dynamic WHERE clause built by the
// frontend's react-querybuilder. See util.AdvancedQuery.
func (s *Service) Query(ctx context.Context, filter dto.Filters, whereCondition string, whereConditionParams []any) (*dto.PaginationResponse[Dto], error) {
	return util.AdvancedQuery(ctx, s.db, tableName, "o", selectQuery, filter, whereCondition, whereConditionParams, row.toDto)
}
