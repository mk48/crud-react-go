package audithistory

import (
	"encoding/json"
	"time"
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
		Data      json.RawMessage `json:"data"`
		CreatedAt time.Time       `json:"createdAt"`
	}
)
