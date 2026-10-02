package user

import (
	"kfamily/internal/model"
	"kfamily/internal/util"

	"github.com/google/uuid"
)

// row is the flat scan target for user queries. It embeds model.AuditRow so
// it joins in the creator's, updater's and deleter's email and the response
// can embed a full IdEmail without an N+1 lookup per row.
type row struct {
	ID        uuid.UUID `db:"id"`
	Sub       string    `db:"sub"`
	Email     string    `db:"email"`
	Name      *string   `db:"name"`
	IsAdmin   bool      `db:"is_admin"`
	IsService bool      `db:"is_service"`

	model.AuditRow
}

func (r row) toDto() Dto {
	return Dto{
		Id:        r.ID.String(),
		Sub:       r.Sub,
		Email:     r.Email,
		Name:      r.Name,
		IsAdmin:   r.IsAdmin,
		IsService: r.IsService,
		AuditDto:  r.AuditRow.ToDto(),
	}
}

// selectColumns lists every column a user query returns, aliased so row can
// scan them directly. It reads through whatever alias "u" is bound to, so
// the same columns work whether u is the user table itself or a RETURNING
// CTE from an update (see service.go).
var selectColumns = `u.id, u.sub, u.email, u.name, u.is_admin, u.is_service, ` + util.AuditSelectColumns("u")

// selectJoins brings in the creator/updater/deleter emails needed by
// selectColumns. This is a self-join back onto the user table (aliased under
// creator/updater/deleter) since users audit each other.
var selectJoins = util.AuditSelectJoins("u")

// selectQuery reads directly from the user table, for plain reads (GetOne,
// List).
var selectQuery = `SELECT ` + selectColumns + ` FROM ` + tableName + ` u ` + selectJoins
