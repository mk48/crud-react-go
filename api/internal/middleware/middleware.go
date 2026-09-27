package middleware

import (
	"kfamily/internal/util"

	"github.com/jmoiron/sqlx"
)

type Middleware struct {
	db  *sqlx.DB
	env *util.AppENV
}

func NewMiddleware(db *sqlx.DB, env *util.AppENV) *Middleware {
	return &Middleware{
		db:  db,
		env: env,
	}
}
