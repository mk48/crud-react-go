package audithistory

import (
	"encoding/json"
	"time"

	"kfamily/internal/dto"
)

type (
	// Dto is one audit_history row: the full source row (snake_case column
	// keys) as it was right after a create/update/delete.
	Dto struct {
		Id        string          `json:"id"`
		TableName string          `json:"tableName"`
		SourceId  string          `json:"sourceId"`
		Action    string          `json:"action"` // create | update | delete
		ChangedBy string          `json:"changedBy"`
		Data      json.RawMessage `json:"data" swaggertype:"object"`
		CreatedAt time.Time       `json:"createdAt"`
		Operation OperationDto    `json:"operation"`
	}

	// OperationDto is the operation that caused a change (see
	// util.RunOperation). The change is a side effect of it when it targets
	// a different record (or none, e.g. an import) than the one changed.
	OperationDto struct {
		Id          string      `json:"id"`
		Kind        string      `json:"kind"`
		PerformedBy dto.IdEmail `json:"performedBy"`
		TargetTable *string     `json:"targetTable"`
		TargetId    *string     `json:"targetId"`
		// The app it came from (see operation.Dto.Client).
		Client    string    `json:"client"`
		CreatedAt time.Time `json:"createdAt"`
	}
)
