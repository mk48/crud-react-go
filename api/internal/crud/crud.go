// Package crud implements the endpoints every audited, soft-deletable table
// shares - get one, list, advanced query, column metadata, create, update
// and delete - once, generically. An entity package describes only what is
// specific to it in a Resource and calls Register (see register.go, which
// also shows how to add custom endpoints next to the generic ones).
package crud

import (
	"context"
	"errors"

	"kfamily/internal/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Resource describes one table for the generic endpoints.
//
// TRow is the sqlx scan target of SelectQuery, TDto the API response shape,
// TCreate/TUpdate the request bodies of POST and PUT. If *TCreate/*TUpdate
// has a `Validate() error` method it runs right after binding (trimming
// and checking the input; an error is a 400).
type Resource[TRow, TDto, TCreate, TUpdate any] struct {
	// Label names one record in messages, e.g. "Sample item" gives
	// "Sample item not found".
	Label string
	// Table is the table name, double-quoted if it's a reserved word
	// (`"user"`). Alias is the alias SelectQuery reads it through.
	Table string
	Alias string
	// SelectQuery selects every column TRow scans, FROM Table Alias plus any
	// joins, without a WHERE clause.
	SelectQuery string
	// SearchColumns are matched (ILIKE) by the list endpoint's searchText.
	SearchColumns []string
	// ToDto converts a scanned row to its response shape.
	ToDto func(TRow) TDto

	// Omit lists generic routes this resource doesn't have, e.g. Create for
	// users (they're provisioned on sign-in).
	Omit Route

	// CreateParams returns the entity's own column values for an INSERT -
	// id, created_at and created_by are added for you. It runs inside the
	// insert's transaction, so it can also check references (see
	// samplechild) or anything else; return Error(status, err) to answer
	// with a specific status. Required unless Create is omitted.
	CreateParams func(w Write, in *TCreate) (map[string]any, error)
	// UpdateParams does the same for an UPDATE (id, updated_at and
	// updated_by are added). Required unless Update is omitted.
	UpdateParams func(w Write, in *TUpdate) (map[string]any, error)
	// BeforeDelete optionally vetoes a soft delete, inside its transaction.
	BeforeDelete func(w Write) error
}

// Write is what the write hooks get: the transaction the change will run
// in, the target row's id (the new id, for a create) and who's acting.
type Write struct {
	Ctx   context.Context
	Tx    *sqlx.Tx
	ID    uuid.UUID
	Actor *model.User
}

// NoInput is the TCreate/TUpdate of a resource that omits that route.
type NoInput struct{}

// Route identifies the generic endpoints, as flags for Resource.Omit.
type Route uint

const (
	Get    Route = 1 << iota // GET    /:id
	List                     // GET    /
	Query                    // GET    /query
	Meta                     // GET    /meta
	Create                   // POST   /
	Update                   // PUT    /:id
	Delete                   // DELETE /:id
)

// HTTPError makes a hook's error reach the client with a specific status
// and message instead of a 500 - create one with Error.
type HTTPError struct {
	Status int
	Err    error
}

func (e *HTTPError) Error() string { return e.Err.Error() }
func (e *HTTPError) Unwrap() error { return e.Err }

// Error wraps err so the handler answers with status and err's message,
// e.g. crud.Error(http.StatusConflict, ErrLastAdmin). errors.Is still sees
// err through it.
func Error(status int, err error) error {
	return &HTTPError{Status: status, Err: err}
}

func asHTTPError(err error) (*HTTPError, bool) {
	var he *HTTPError
	ok := errors.As(err, &he)
	return he, ok
}
