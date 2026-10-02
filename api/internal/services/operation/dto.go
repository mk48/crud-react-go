package operation

import (
	"encoding/json"
	"time"

	"kfamily/internal/dto"
)

type (
	// Dto is one operation row - a user action that wrote data (see
	// util.RunOperation).
	Dto struct {
		Id          string      `json:"id"`
		Kind        string      `json:"kind"`
		PerformedBy dto.IdEmail `json:"performedBy"`
		// The record the user acted on directly; null for operations with
		// no single target (e.g. an import).
		TargetTable *string         `json:"targetTable"`
		TargetId    *string         `json:"targetId"`
		Metadata    json.RawMessage `json:"metadata" swaggertype:"object"`
		// The app it came from (web, mobile, admin, batch:<job>, system) -
		// verified from the access token.
		Client string `json:"client"`
		// What the app reported about itself (version, platform, ip,
		// userAgent) - hints, not proof.
		ClientInfo json.RawMessage `json:"clientInfo" swaggertype:"object"`
		// W3C trace id of the request/task that ran it.
		TraceId   *string   `json:"traceId"`
		CreatedAt time.Time `json:"createdAt"`
		// Number of audit_history rows (record changes) it caused.
		ChangeCount int `json:"changeCount"`
	}

	// DetailDto is an operation with the record changes it caused, in the
	// order they were made - at most maxChanges of them (compare with
	// ChangeCount).
	DetailDto struct {
		Dto
		Changes []ChangeDto `json:"changes"`
	}

	// ChangeDto is one audit_history row of an operation, with the record's
	// previous version so it can be diffed.
	ChangeDto struct {
		Id        string          `json:"id"`
		TableName string          `json:"tableName"`
		SourceId  string          `json:"sourceId"`
		Action    string          `json:"action"` // create | update | delete
		Data      json.RawMessage `json:"data" swaggertype:"object"`
		// The record as it was before this change; null for a create.
		PreviousData *json.RawMessage `json:"previousData" swaggertype:"object"`
		// Whether the record has been deleted since (by this or a later
		// operation) - its regular view page no longer shows it.
		SourceDeleted bool      `json:"sourceDeleted"`
		CreatedAt     time.Time `json:"createdAt"`
	}
)
