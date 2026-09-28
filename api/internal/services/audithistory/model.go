package audithistory

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// row is the flat scan target for audit_history queries.
type row struct {
	ID        uuid.UUID       `db:"id"`
	TableName string          `db:"table_name"`
	SourceID  uuid.UUID       `db:"source_id"`
	Action    string          `db:"action"`
	ChangedBy uuid.UUID       `db:"changed_by"`
	Data      json.RawMessage `db:"data"`
	CreatedAt time.Time       `db:"created_at"`
}

func (r row) toDto() Dto {
	return Dto{
		Id:        r.ID.String(),
		TableName: r.TableName,
		SourceId:  r.SourceID.String(),
		Action:    r.Action,
		ChangedBy: r.ChangedBy.String(),
		Data:      r.Data,
		CreatedAt: r.CreatedAt,
	}
}

const selectQuery = `SELECT id, table_name, source_id, action, changed_by, data, created_at FROM audit_history`
