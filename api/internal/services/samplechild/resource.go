package samplechild

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"kfamily/internal/crud"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

const tableName = "sample_child_items"

// ErrInvalidSampleItem is returned by Create/Update when sampleItemId isn't
// a valid UUID or doesn't reference an existing, non-deleted sample item -
// checked up front so the client gets a 400 instead of the foreign key
// violation surfacing as a 500.
var ErrInvalidSampleItem = errors.New("sampleItemId must reference an existing sample item")

// resource is everything specific to sample child items; the endpoints
// themselves are crud's (see crud.Register, which also explains how to add
// custom ones).
var resource = crud.Resource[row, Dto, CreateInputDto, UpdateInputDto]{
	Label:         "Sample child item",
	Table:         tableName,
	Alias:         "sc",
	SelectQuery:   selectQuery,
	SearchColumns: []string{"name"},
	ToDto:         row.toDto,

	CreateParams: func(w crud.Write, in *CreateInputDto) (map[string]any, error) {
		return params(w, in.SampleItemId, in.Name)
	},
	UpdateParams: func(w crud.Write, in *UpdateInputDto) (map[string]any, error) {
		return params(w, in.SampleItemId, in.Name)
	},
}

func params(w crud.Write, rawSampleItemId string, name string) (map[string]any, error) {
	sampleItemId, err := validSampleItemId(w, rawSampleItemId)
	if err != nil {
		return nil, err
	}
	return map[string]any{"sample_item_id": sampleItemId, "name": name}, nil
}

// validSampleItemId parses raw and checks it references a non-deleted
// sample item (the foreign key alone would accept a soft-deleted parent).
// Checked in the write's transaction.
func validSampleItemId(w crud.Write, raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.UUID{}, crud.Error(http.StatusBadRequest, ErrInvalidSampleItem)
	}

	var exists bool
	query := `SELECT EXISTS (SELECT 1 FROM sample_items WHERE id = $1 AND deleted_at IS NULL)`
	if err := w.Tx.GetContext(w.Ctx, &exists, query, id); err != nil {
		return uuid.UUID{}, fmt.Errorf("unable to check sample item. err: %w", err)
	}
	if !exists {
		return uuid.UUID{}, crud.Error(http.StatusBadRequest, ErrInvalidSampleItem)
	}

	return id, nil
}

func Init(ctx context.Context, deps crud.Deps, api *echo.Group) {
	crud.Register(ctx, deps, api, "/v1/sample-children", resource)
}
