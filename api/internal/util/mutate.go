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

// Execer is satisfied by both *sqlx.DB and *sqlx.Tx, so Insert/UpdateByID
// (and the RecordAudit call they make) can run either directly or inside an
// existing transaction - see photosubmission.Service.Approve, which needs
// several inserts/updates across different tables to commit atomically.
type Execer interface {
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
}

// Insert builds "INSERT INTO tableName (<cols>) VALUES (<:cols>)" from
// params' keys, executes it, and records the row in audit_history (see
// RecordAudit). params must include "id" set to id. Deriving the column
// list from params instead of writing it out separately means there's no
// place left for the two to drift apart as columns are added.
func Insert(ctx context.Context, db Execer, tableName string, id uuid.UUID, params map[string]any) error {
	cols := sortedKeys(params)

	placeholders := make([]string, len(cols))
	for i, c := range cols {
		placeholders[i] = ":" + c
	}

	query := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, tableName, strings.Join(cols, ", "), strings.Join(placeholders, ", "))

	return inTx(ctx, db, func(tx Execer) error {
		if _, err := tx.NamedExecContext(ctx, query, params); err != nil {
			return fmt.Errorf("unable to insert into %s. err: %w", tableName, err)
		}

		return RecordAudit(ctx, tx, id, params)
	})
}

// UpdateByID builds "UPDATE tableName SET <col = :col, ...> WHERE id = :id
// AND deleted_at IS NULL" from params' keys (excluding "id" from the SET
// list), executes it, and records the row in audit_history (see
// RecordAudit). params must include "id" set to id.
//
// Soft-deleted rows are read-only - they can't be edited or deleted again
// (which would overwrite deleted_at/deleted_by). If no live row matches id,
// nothing is changed or audited and an error wrapping sql.ErrNoRows is
// returned, so handlers can answer 404 the same way GetOne does.
func UpdateByID(ctx context.Context, db Execer, tableName string, id uuid.UUID, params map[string]any) error {
	cols := sortedKeys(params)

	setClauses := make([]string, 0, len(cols))
	for _, c := range cols {
		if c == "id" {
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = :%s", c, c))
	}

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = :id AND deleted_at IS NULL`, tableName, strings.Join(setClauses, ", "))

	return inTx(ctx, db, func(tx Execer) error {
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

		return RecordAudit(ctx, tx, id, params)
	})
}

// inTx runs fn so that the row change and its audit_history record commit
// or roll back together. Given a *sqlx.DB it opens (and commits/rolls back)
// its own transaction; given anything else - i.e. a caller's *sqlx.Tx - fn
// just runs on it, and the caller stays in charge of commit/rollback.
func inTx(ctx context.Context, db Execer, fn func(tx Execer) error) error {
	sqlDB, ok := db.(*sqlx.DB)
	if !ok {
		return fn(db)
	}

	return WithTx(ctx, sqlDB, func(tx *sqlx.Tx) error { return fn(tx) })
}

// WithTx runs fn in a new transaction on db - committed if fn returns nil,
// rolled back otherwise. Pass tx on to Insert/UpdateByID to fold them (and
// their audit records) into the same transaction.
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
