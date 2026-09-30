package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"kfamily/internal/crud"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v5"
)

// tableName is double-quoted because "user" is a reserved word in Postgres.
const tableName = `"user"`

var (
	// ErrLastAdmin: the change would leave no live admin.
	ErrLastAdmin = errors.New("can't remove the last remaining admin")
	// ErrSelfDemote/ErrSelfDelete: an admin could otherwise lock themselves
	// out (deleted users are refused by AuthMiddleware).
	ErrSelfDemote = errors.New("You can't remove your own admin access")
	ErrSelfDelete = errors.New("You can't delete your own account")
)

// resource is everything specific to users; the generic endpoints are
// crud's. There is no Create - users are provisioned automatically by
// AuthMiddleware the first time they sign in. GET /me is a custom endpoint
// (see Init and handler.go).
var resource = crud.Resource[row, Dto, crud.NoInput, UpdateInputDto]{
	Label:         "User",
	Table:         tableName,
	Alias:         "u",
	SelectQuery:   selectQuery,
	SearchColumns: []string{"name", "email"},
	ToDto:         row.toDto,
	Omit:          crud.Create,

	// Only name and is_admin change - sub and email are set once by
	// AuthMiddleware. is_admin is left as is when IsAdmin is omitted.
	UpdateParams: func(w crud.Write, in *UpdateInputDto) (map[string]any, error) {
		params := map[string]any{"name": in.Name}
		if in.IsAdmin == nil {
			return params, nil
		}
		params["is_admin"] = *in.IsAdmin

		if !*in.IsAdmin {
			if w.ID == w.Actor.ID {
				return nil, crud.Error(http.StatusBadRequest, ErrSelfDemote)
			}
			if err := guardLastAdmin(w.Ctx, w.Tx, w.ID); err != nil {
				return nil, err
			}
		}
		return params, nil
	},

	BeforeDelete: func(w crud.Write) error {
		if w.ID == w.Actor.ID {
			return crud.Error(http.StatusBadRequest, ErrSelfDelete)
		}
		return guardLastAdmin(w.Ctx, w.Tx, w.ID)
	},
}

// guardLastAdmin refuses (409) if id is the only live admin. It locks the
// live admin rows (FOR UPDATE) until tx ends, so two admins demoting or
// deleting each other at the same moment can't both pass the check and
// leave no admin at all.
func guardLastAdmin(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) error {
	adminIds := []uuid.UUID{}
	query := `SELECT id FROM "user" WHERE is_admin AND deleted_at IS NULL FOR UPDATE`
	if err := tx.SelectContext(ctx, &adminIds, query); err != nil {
		return fmt.Errorf("unable to load admins. err: %w", err)
	}

	if len(adminIds) == 1 && adminIds[0] == id {
		return crud.Error(http.StatusConflict, ErrLastAdmin)
	}

	return nil
}

func Init(ctx context.Context, deps crud.Deps, api *echo.Group) {
	g, svc := crud.Register(ctx, deps, api, "/v1/users", resource)

	// Custom endpoint next to the generic ones - see crud.Register.
	h := &Handler{svc: svc}
	g.GET("/me", h.Me)
}
