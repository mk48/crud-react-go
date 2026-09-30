package crud

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"kfamily/internal/dto"
	"kfamily/internal/model"
	"kfamily/internal/util"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Service runs a Resource's queries and audited writes. Register returns
// it, so a custom endpoint can reuse these next to its own queries.
type Service[TRow, TDto, TCreate, TUpdate any] struct {
	db  *sqlx.DB
	res Resource[TRow, TDto, TCreate, TUpdate]
}

func NewService[TRow, TDto, TCreate, TUpdate any](db *sqlx.DB, res Resource[TRow, TDto, TCreate, TUpdate]) *Service[TRow, TDto, TCreate, TUpdate] {
	return &Service[TRow, TDto, TCreate, TUpdate]{db: db, res: res}
}

// DB is the connection pool, for custom queries.
func (s *Service[TRow, TDto, TCreate, TUpdate]) DB() *sqlx.DB {
	return s.db
}

// GetOne returns one row, or an error wrapping sql.ErrNoRows if there's no
// such (live, unless includeDeleted) row.
func (s *Service[TRow, TDto, TCreate, TUpdate]) GetOne(ctx context.Context, id uuid.UUID, includeDeleted bool) (TDto, error) {
	query := fmt.Sprintf(`%s WHERE %s.id = $1`, s.res.SelectQuery, s.res.Alias)
	if !includeDeleted {
		query += fmt.Sprintf(` AND %s.deleted_at IS NULL`, s.res.Alias)
	}

	var r TRow
	if err := s.db.GetContext(ctx, &r, query, id); err != nil {
		var zero TDto
		if errors.Is(err, sql.ErrNoRows) {
			return zero, err
		}
		return zero, fmt.Errorf("unable to get %s. err: %w", s.name(), err)
	}

	return s.res.ToDto(r), nil
}

// List returns a page of rows matching filter's search/sort/pagination -
// see util.List.
func (s *Service[TRow, TDto, TCreate, TUpdate]) List(ctx context.Context, filter dto.Filters) (*dto.PaginationResponse[TDto], error) {
	return util.List(ctx, s.db, s.res.Table, s.res.Alias, s.res.SelectQuery, s.res.SearchColumns, filter, s.res.ToDto)
}

// Query lists rows matching a query-builder WHERE fragment - see
// util.AdvancedQuery.
func (s *Service[TRow, TDto, TCreate, TUpdate]) Query(ctx context.Context, filter dto.Filters, whereCondition string, whereConditionParams []any) (*dto.PaginationResponse[TDto], error) {
	return util.AdvancedQuery(ctx, s.db, s.res.Table, s.res.Alias, s.res.SelectQuery, filter, whereCondition, whereConditionParams, s.res.ToDto)
}

// Metadata returns the table's columns and Postgres data types, so a
// frontend can discover valid sort/filter columns.
func (s *Service[TRow, TDto, TCreate, TUpdate]) Metadata(ctx context.Context) ([]dto.ColumnMeta, error) {
	columns, err := util.TableColumns(ctx, s.db, s.res.Table)
	if err != nil {
		return nil, fmt.Errorf("unable to load %s table metadata. err: %w", s.res.Table, err)
	}
	return columns, nil
}

// Create inserts a row (and its audit record) in one transaction, with
// Resource.CreateParams' columns, and returns the new id.
func (s *Service[TRow, TDto, TCreate, TUpdate]) Create(ctx context.Context, in *TCreate, actor *model.User) (uuid.UUID, error) {
	id := uuid.New()

	err := util.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		params, err := s.res.CreateParams(Write{Ctx: ctx, Tx: tx, ID: id, Actor: actor}, in)
		if err != nil {
			return err
		}
		params["id"] = id
		params["created_at"] = time.Now()
		params["created_by"] = actor.ID

		return util.Insert(ctx, tx, s.res.Table, id, params)
	})
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("unable to create %s. err: %w", s.name(), err)
	}

	return id, nil
}

// Update changes a live row (and records it) in one transaction, with
// Resource.UpdateParams' columns. An error wraps sql.ErrNoRows if there's
// no such live row.
func (s *Service[TRow, TDto, TCreate, TUpdate]) Update(ctx context.Context, id uuid.UUID, in *TUpdate, actor *model.User) error {
	err := util.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		params, err := s.res.UpdateParams(Write{Ctx: ctx, Tx: tx, ID: id, Actor: actor}, in)
		if err != nil {
			return err
		}
		params["id"] = id
		params["updated_at"] = time.Now()
		params["updated_by"] = actor.ID

		return util.UpdateByID(ctx, tx, s.res.Table, id, params)
	})
	if err != nil {
		return fmt.Errorf("unable to update %s. err: %w", s.name(), err)
	}

	return nil
}

// Delete soft-deletes a live row (and records it) in one transaction,
// after Resource.BeforeDelete. An error wraps sql.ErrNoRows if there's no
// such live row.
func (s *Service[TRow, TDto, TCreate, TUpdate]) Delete(ctx context.Context, id uuid.UUID, actor *model.User) error {
	err := util.WithTx(ctx, s.db, func(tx *sqlx.Tx) error {
		if s.res.BeforeDelete != nil {
			if err := s.res.BeforeDelete(Write{Ctx: ctx, Tx: tx, ID: id, Actor: actor}); err != nil {
				return err
			}
		}

		return util.UpdateByID(ctx, tx, s.res.Table, id, map[string]any{
			"id":         id,
			"deleted_at": time.Now(),
			"deleted_by": actor.ID,
		})
	})
	if err != nil {
		return fmt.Errorf("unable to delete %s. err: %w", s.name(), err)
	}

	return nil
}

// name is Label in lower case, for messages mid-sentence.
func (s *Service[TRow, TDto, TCreate, TUpdate]) name() string {
	return strings.ToLower(s.res.Label)
}
