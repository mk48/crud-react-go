package util

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Insert builds "INSERT INTO tableName (<cols>) VALUES (<:cols>)" from
// params' keys, executes it, and records the row in audit_history (see
// RecordAudit). params must include "id" set to id. Deriving the column
// list from params instead of writing it out separately means there's no
// place left for the two to drift apart as columns are added.
//
// tx and ctx must be the ones RunOperation hands out, so the change is
// recorded as part of that operation.
func Insert(ctx context.Context, tx *sqlx.Tx, tableName string, id uuid.UUID, params map[string]any) error {
	cols := sortedKeys(params)

	placeholders := make([]string, len(cols))
	for i, c := range cols {
		placeholders[i] = ":" + c
	}

	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, tableName, strings.Join(cols, ", "), strings.Join(placeholders, ", "))

	actor, err := auditActor(params, "created_by")
	if err != nil {
		return err
	}

	if _, err := tx.NamedExecContext(ctx, query, params); err != nil {
		return fmt.Errorf("unable to insert into %s. err: %w", tableName, err)
	}

	return RecordAudit(ctx, tx, tableName, id, AuditCreate, actor)
}

// UpdateByID builds "UPDATE tableName SET <col = :col, ...> WHERE id = :id
// AND deleted_at IS NULL" from params' keys (excluding "id" from the SET
// list), executes it, and records the row in audit_history (see
// RecordAudit). params must include "id" set to id. As with Insert, tx and
// ctx must come from RunOperation.
//
// Soft-deleted rows are read-only - they can't be edited or deleted again
// (which would overwrite deleted_at/deleted_by). If no live row matches id,
// nothing is changed or audited and an error wrapping sql.ErrNoRows is
// returned, so handlers can answer 404 the same way GetOne does.
func UpdateByID(ctx context.Context, tx *sqlx.Tx, tableName string, id uuid.UUID, params map[string]any) error {
	cols := sortedKeys(params)

	setClauses := make([]string, 0, len(cols))
	for _, c := range cols {
		if c == "id" {
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = :%s", c, c))
	}

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = :id AND deleted_at IS NULL`, tableName, strings.Join(setClauses, ", "))

	// A soft delete is an update that sets deleted_at/deleted_by.
	action, actorKey := AuditUpdate, "updated_by"
	if _, ok := params["deleted_by"]; ok {
		action, actorKey = AuditDelete, "deleted_by"
	}
	actor, err := auditActor(params, actorKey)
	if err != nil {
		return err
	}

	result, err := tx.NamedExecContext(ctx, query, params)
	if err != nil {
		return fmt.Errorf("unable to update %s. err: %w", tableName, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to read rows affected for %s. err: %w", tableName, err)
	}
	if affected == 0 {
		return fmt.Errorf("no live %s row with id %s. err: %w", tableName, id, sql.ErrNoRows)
	}

	return RecordAudit(ctx, tx, tableName, id, action, actor)
}

// WithTx runs fn in a new transaction on db - committed if fn returns nil,
// rolled back otherwise. Writes to audited tables use RunOperation instead,
// which runs on top of this.
func WithTx(ctx context.Context, db *sqlx.DB, fn func(tx *sqlx.Tx) error) error {
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to begin transaction. err: %w", err)
	}
	// No-op once Commit has succeeded.
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction. err: %w", err)
	}

	return nil
}

// sortedKeys returns params' keys sorted, so the generated query text (and
// column order in it) is deterministic across calls - map iteration order
// isn't, which would otherwise make query logs/plans churn for no reason.
func sortedKeys(params map[string]any) []string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
