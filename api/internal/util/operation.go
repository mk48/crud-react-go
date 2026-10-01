package util

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Operation describes one user action that writes data - see the operation
// table in migrations/0001-init.sql. Every Insert/UpdateByID must run
// inside one (via RunOperation), so each audit_history row can be traced
// back to the action that caused it.
type Operation struct {
	// Dotted "<resource>.<action>" key, e.g. "sample.create" - the web app
	// translates it for display (operation.kind.* in translation.json).
	Kind        string
	PerformedBy uuid.UUID
	// The record the user acted on directly; leave TargetTable empty for
	// operations with no single target (e.g. an import). Writes to any other
	// record within the operation show up in that record's history as side
	// effects of it.
	TargetTable string
	TargetID    uuid.UUID
	// Optional free-form context, e.g. an import's file name and row count.
	Metadata map[string]any
}

type operationCtxKey struct{}

// RunOperation records op and runs fn in one transaction, passing fn a
// context that carries op's id - RecordAudit stamps it on every
// audit_history row written under it. The operation, the data changes and
// their audit records all commit or roll back together.
//
// Operations don't nest: fn must use the tx it's given, not call another
// service method that would start its own operation (and transaction).
func RunOperation(ctx context.Context, db *sqlx.DB, op Operation, fn func(ctx context.Context, tx *sqlx.Tx) error) error {
	if _, ok := operationID(ctx); ok {
		return errors.New("operations can't nest - pass the existing tx down instead of starting another operation")
	}
	if op.Kind == "" {
		return errors.New("operation kind is required")
	}

	metadata := op.Metadata
	if metadata == nil {
		metadata = map[string]any{}
	}
	metadataJson, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("unable to encode operation metadata. err: %w", err)
	}

	var targetTable *string
	var targetID *uuid.UUID
	if op.TargetTable != "" {
		// Stored bare, like audit_history.table_name (see RecordAudit).
		t := strings.Trim(op.TargetTable, `"`)
		targetTable, targetID = &t, &op.TargetID
	}

	id := uuid.New()

	return WithTx(ctx, db, func(tx *sqlx.Tx) error {
		_, err := tx.NamedExecContext(ctx, `
			INSERT INTO operation (id, kind, performed_by, target_table, target_id, metadata, created_at)
			VALUES (:id, :kind, :performed_by, :target_table, :target_id, :metadata, clock_timestamp())`,
			map[string]any{
				"id":           id,
				"kind":         op.Kind,
				"performed_by": op.PerformedBy,
				"target_table": targetTable,
				"target_id":    targetID,
				"metadata":     string(metadataJson),
			})
		if err != nil {
			return fmt.Errorf("unable to record operation %s. err: %w", op.Kind, err)
		}

		return fn(context.WithValue(ctx, operationCtxKey{}, id), tx)
	})
}

// operationID returns the id of the operation ctx runs under (see
// RunOperation), if any.
func operationID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(operationCtxKey{}).(uuid.UUID)
	return id, ok
}
