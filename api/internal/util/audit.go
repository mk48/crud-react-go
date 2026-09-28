package util

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Audit actions recorded in audit_history.action.
const (
	AuditCreate = "create"
	AuditUpdate = "update"
	AuditDelete = "delete"
)

// RecordAudit inserts a row into audit_history holding the full current row
// of tableName with id sourceID (as JSON, snake_case column keys) - call it
// right after the create/update/delete, inside the same transaction (see
// Insert/UpdateByID), so the snapshot is exactly the state just written.
// changedBy is the user who made the change.
func RecordAudit(ctx context.Context, db Execer, tableName string, sourceID uuid.UUID, action string, changedBy uuid.UUID) error {
	// clock_timestamp(), not now() (the transaction start), so several
	// changes made in one transaction still sort in the order they happened.
	// tableName is a trusted constant from each service (never user input),
	// quoted where needed (e.g. `"user"`); audit_history stores it bare.
	query := fmt.Sprintf(`
		INSERT INTO audit_history (id, table_name, source_id, action, changed_by, data, created_at)
		SELECT :id, :table_name, t.id, :action, :changed_by, to_jsonb(t), clock_timestamp()
		FROM %s t
		WHERE t.id = :source_id`, tableName)

	result, err := db.NamedExecContext(ctx, query, map[string]any{
		"id":         uuid.New(),
		"table_name": strings.Trim(tableName, `"`),
		"source_id":  sourceID,
		"action":     action,
		"changed_by": changedBy,
	})
	if err != nil {
		return fmt.Errorf("unable to record audit history. err: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("unable to read audit history rows affected. err: %w", err)
	}
	if affected != 1 {
		return fmt.Errorf("unable to record audit history: %s row %s not found", tableName, sourceID)
	}

	return nil
}

// auditActor returns params[key] - the created_by/updated_by/deleted_by
// value of an Insert/UpdateByID - as the audit_history changed_by.
func auditActor(params map[string]any, key string) (uuid.UUID, error) {
	actor, ok := params[key].(uuid.UUID)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("audited write must set %s to the acting user's id", key)
	}
	return actor, nil
}

// AuditSelectColumns lists the created_by/updated_by/deleted_by trio and
// their joined creator/updater/deleter emails, read through alias - the same
// columns every audited table's model.go needs alongside its own business
// columns. alias must be bound (directly or via a RETURNING CTE) to a row
// with those columns, and the query must bring in AuditSelectJoins(alias) so
// creator/updater/deleter resolve.
func AuditSelectColumns(alias string) string {
	return fmt.Sprintf(`
		%[1]s.created_at, %[1]s.created_by, creator.email AS created_by_email,
		%[1]s.updated_at, %[1]s.updated_by, updater.email AS updated_by_email,
		%[1]s.deleted_at, %[1]s.deleted_by, deleter.email AS deleted_by_email
	`, alias)
}

// AuditSelectJoins brings in the creator/updater/deleter emails needed by
// AuditSelectColumns(alias).
func AuditSelectJoins(alias string) string {
	return `
		JOIN "user" creator ON creator.id = ` + alias + `.created_by
		LEFT JOIN "user" updater ON updater.id = ` + alias + `.updated_by
		LEFT JOIN "user" deleter ON deleter.id = ` + alias + `.deleted_by
	`
}
