package audithistory

import (
	"encoding/json"
	"time"
)

type (
	// Dto is one audit_history row: a snapshot of a source record's data at
	// the time of a create/update/delete.
	Dto struct {
		Id        string          `json:"id"`
		SourceId  string          `json:"sourceId"`
		Data      json.RawMessage `json:"data"`
		CreatedAt time.Time       `json:"createdAt"`
	}
)
