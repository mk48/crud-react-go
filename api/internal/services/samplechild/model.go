package samplechild

import (
	"kfamily/internal/dto"
	"kfamily/internal/model"
	"kfamily/internal/util"

	"github.com/google/uuid"
)

// row is the flat scan target for sample_child_items queries. It embeds
// model.AuditRow so it joins in the creator's, updater's and deleter's email,
// plus the parent sample item's name, so the response can embed both without
// an N+1 lookup per row.
type row struct {
	ID             uuid.UUID `db:"id"`
	SampleItemID   uuid.UUID `db:"sample_item_id"`
	SampleItemName string    `db:"sample_item_name"`
	Name           string    `db:"name"`

	model.AuditRow
}

func (r row) toDto() Dto {
	return Dto{
		Id: r.ID.String(),
		SampleItem: dto.IdName{
			Id:   r.SampleItemID.String(),
			Name: r.SampleItemName,
		},
		Name:     r.Name,
		AuditDto: r.AuditRow.ToDto(),
	}
}

// selectColumns lists every column a sample_child_items query returns,
// aliased so row can scan them directly.
var selectColumns = `sc.id, sc.sample_item_id, si.name AS sample_item_name, sc.name, ` + util.AuditSelectColumns("sc")

// selectJoins brings in the parent sample item's name and the
// creator/updater/deleter emails needed by selectColumns. sample_item_id is
// NOT NULL, so the parent join is an inner join. It deliberately doesn't
// filter out soft-deleted parents - a child keeps showing its parent's name
// after the parent is deleted.
var selectJoins = `
	JOIN sample_items si ON si.id = sc.sample_item_id
` + util.AuditSelectJoins("sc")

// selectQuery is crud.Resource.SelectQuery: every column, no WHERE clause.
var selectQuery = `SELECT ` + selectColumns + ` FROM ` + tableName + ` sc ` + selectJoins
