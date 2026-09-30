package sample

import (
	"kfamily/internal/model"
	"kfamily/internal/util"

	"github.com/google/uuid"
)

// row is the flat scan target for sample_item queries. It embeds
// model.AuditRow so it joins in the creator's, updater's and deleter's email
// and the response can embed a full IdEmail without an N+1 lookup per row.
type row struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description *string   `db:"description"`

	model.AuditRow
}

func (r row) toDto() Dto {
	return Dto{
		Id:          r.ID.String(),
		Name:        r.Name,
		Description: r.Description,
		AuditDto:    r.AuditRow.ToDto(),
	}
}

// selectColumns lists every column a sample_item query returns, aliased so
// row can scan them directly. It reads through whatever alias "s" is bound
// to, so the same columns work whether s is the sample_items table itself or
// a RETURNING CTE from an insert/update (see service.go).
var selectColumns = `s.id, s.name, s.description, ` + util.AuditSelectColumns("s")

// selectJoins brings in the creator/updater/deleter emails needed by
// selectColumns.
var selectJoins = util.AuditSelectJoins("s")

// selectQuery reads directly from the sample_items table, for plain reads
// (GetOne, List). Mutations build their own query around selectColumns and
// selectJoins so they can read from a RETURNING CTE instead.
var selectQuery = `SELECT ` + selectColumns + ` FROM ` + tableName + ` s ` + selectJoins
