package sample

import (
	"context"

	"kfamily/internal/crud"

	"github.com/labstack/echo/v5"
)

const tableName = "sample_items"

// resource is everything specific to sample items; the endpoints
// themselves are crud's (see crud.Register, which also explains how to add
// custom ones).
var resource = crud.Resource[row, Dto, CreateInputDto, UpdateInputDto]{
	Label:         "Sample item",
	Table:         tableName,
	Alias:         "s",
	SelectQuery:   selectQuery,
	SearchColumns: []string{"name", "description"},
	ToDto:         row.toDto,

	CreateParams: func(_ crud.Write, in *CreateInputDto) (map[string]any, error) {
		return map[string]any{"name": in.Name, "description": in.Description}, nil
	},
	UpdateParams: func(_ crud.Write, in *UpdateInputDto) (map[string]any, error) {
		return map[string]any{"name": in.Name, "description": in.Description}, nil
	},
}

func Init(ctx context.Context, deps crud.Deps, api *echo.Group) {
	crud.Register(ctx, deps, api, "/v1/samples", resource)
}
