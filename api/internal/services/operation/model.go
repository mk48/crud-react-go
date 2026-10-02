package operation

import (
	"encoding/json"
	"time"

	"kfamily/internal/dto"

	"github.com/google/uuid"
)

// row is the flat scan target for operation queries, with the performer's
// email and the number of changes joined in.
type row struct {
	ID               uuid.UUID       `db:"id"`
	Kind             string          `db:"kind"`
	PerformedBy      uuid.UUID       `db:"performed_by"`
	PerformedByEmail string          `db:"performed_by_email"`
	TargetTable      *string         `db:"target_table"`
	TargetID         *uuid.UUID      `db:"target_id"`
	Metadata         json.RawMessage `db:"metadata"`
	Client           string          `db:"client"`
	ClientInfo       json.RawMessage `db:"client_info"`
	TraceID          *string         `db:"trace_id"`
	CreatedAt        time.Time       `db:"created_at"`
	ChangeCount      int             `db:"change_count"`
}

func (r row) toDto() Dto {
	d := Dto{
		Id:          r.ID.String(),
		Kind:        r.Kind,
		PerformedBy: dto.IdEmail{Id: r.PerformedBy.String(), Email: r.PerformedByEmail},
		TargetTable: r.TargetTable,
		Metadata:    r.Metadata,
		Client:      r.Client,
		ClientInfo:  r.ClientInfo,
		TraceId:     r.TraceID,
		CreatedAt:   r.CreatedAt,
		ChangeCount: r.ChangeCount,
	}
	if r.TargetID != nil {
		id := r.TargetID.String()
		d.TargetId = &id
	}
	return d
}

// selectQuery reads from the operation table aliased "o", as util.List and
// util.AdvancedQuery require.
const selectQuery = `
	SELECT o.id, o.kind, o.performed_by, performer.email AS performed_by_email,
		o.target_table, o.target_id, o.metadata, o.client, o.client_info, o.trace_id, o.created_at,
		(SELECT COUNT(*) FROM audit_history a WHERE a.operation_id = o.id) AS change_count
	FROM operation o
	JOIN "user" performer ON performer.id = o.performed_by `

// changeRow is the flat scan target for an operation's audit_history rows.
type changeRow struct {
	ID            uuid.UUID        `db:"id"`
	TableName     string           `db:"table_name"`
	SourceID      uuid.UUID        `db:"source_id"`
	Action        string           `db:"action"`
	Data          json.RawMessage  `db:"data"`
	PreviousData  *json.RawMessage `db:"previous_data"`
	SourceDeleted bool             `db:"source_deleted"`
	CreatedAt     time.Time        `db:"created_at"`
}

func (r changeRow) toDto() ChangeDto {
	return ChangeDto{
		Id:            r.ID.String(),
		TableName:     r.TableName,
		SourceId:      r.SourceID.String(),
		Action:        r.Action,
		Data:          r.Data,
		PreviousData:  r.PreviousData,
		SourceDeleted: r.SourceDeleted,
		CreatedAt:     r.CreatedAt,
	}
}

// changesQuery lists operation $1's changes in the order they were made,
// each with the same record's preceding audit_history snapshot (if any) to
// diff against.
const changesQuery = `
	SELECT a.id, a.table_name, a.source_id, a.action, a.data, a.created_at,
		prev.data AS previous_data,
		EXISTS (
			SELECT 1 FROM audit_history d WHERE d.source_id = a.source_id AND d.action = 'delete'
		) AS source_deleted
	FROM audit_history a
	LEFT JOIN LATERAL (
		SELECT p.data FROM audit_history p
		WHERE p.source_id = a.source_id AND (p.created_at, p.id) < (a.created_at, a.id)
		ORDER BY p.created_at DESC, p.id DESC
		LIMIT 1
	) prev ON TRUE
	WHERE a.operation_id = $1
	ORDER BY a.created_at, a.id
	LIMIT $2`
