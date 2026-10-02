package audithistory

import (
	"encoding/json"
	"time"

	"kfamily/internal/dto"

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

	OperationID          uuid.UUID  `db:"operation_id"`
	OperationKind        string     `db:"operation_kind"`
	OperationPerformedBy uuid.UUID  `db:"operation_performed_by"`
	OperationPerformer   string     `db:"operation_performer_email"`
	OperationTargetTable *string    `db:"operation_target_table"`
	OperationTargetID    *uuid.UUID `db:"operation_target_id"`
	OperationClient      string     `db:"operation_client"`
	OperationCreatedAt   time.Time  `db:"operation_created_at"`
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
		Operation: OperationDto{
			Id:          r.OperationID.String(),
			Kind:        r.OperationKind,
			PerformedBy: dto.IdEmail{Id: r.OperationPerformedBy.String(), Email: r.OperationPerformer},
			TargetTable: r.OperationTargetTable,
			TargetId:    uuidString(r.OperationTargetID),
			Client:      r.OperationClient,
			CreatedAt:   r.OperationCreatedAt,
		},
	}
}

func uuidString(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}
	s := id.String()
	return &s
}

const selectQuery = `
	SELECT a.id, a.table_name, a.source_id, a.action, a.changed_by, a.data, a.created_at,
		o.id AS operation_id, o.kind AS operation_kind, o.performed_by AS operation_performed_by,
		performer.email AS operation_performer_email, o.target_table AS operation_target_table,
		o.target_id AS operation_target_id, o.client AS operation_client,
		o.created_at AS operation_created_at
	FROM audit_history a
	JOIN operation o ON o.id = a.operation_id
	JOIN "user" performer ON performer.id = o.performed_by`
