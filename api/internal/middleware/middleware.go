package middleware

import (
	"kfamily/internal/util"

	"github.com/jmoiron/sqlx"
)

type Middleware struct {
	db      *sqlx.DB
	env     *util.AppENV
	clients *util.ClientRegistry
}

func NewMiddleware(db *sqlx.DB, env *util.AppENV, clients *util.ClientRegistry) *Middleware {
	return &Middleware{
		db:      db,
		env:     env,
		clients: clients,
	}
}
