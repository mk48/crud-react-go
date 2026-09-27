package audithistory

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// row is the flat scan target for audit_history queries.
type row struct {
	ID        uuid.UUID       `db:"id"`
	SourceID  uuid.UUID       `db:"source_id"`
	Data      json.RawMessage `db:"data"`
	CreatedAt time.Time       `db:"created_at"`
}

func (r row) toDto() Dto {
	return Dto{
		Id:        r.ID.String(),
		SourceId:  r.SourceID.String(),
		Data:      r.Data,
		CreatedAt: r.CreatedAt,
	}
}

const selectQuery = `SELECT id, source_id, data, created_at FROM audit_history`
