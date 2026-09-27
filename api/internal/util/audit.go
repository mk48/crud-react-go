package util

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

// RecordAudit inserts a row into audit_history capturing data as the new
// state of sourceID after a create/update/delete. Call it right after the
// mutation succeeds, typically passing the same params map used for the
// insert/update/delete. db may be *sqlx.DB or *sqlx.Tx (see Execer) so it
// can be folded into an existing transaction.
func RecordAudit(ctx context.Context, db Execer, sourceID uuid.UUID, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("unable to marshal audit data. err: %w", err)
	}

	_, err = db.NamedExecContext(ctx, `INSERT INTO audit_history (id, source_id, data) VALUES (:id, :source_id, :data)`,
		map[string]any{
			"id":        uuid.New(),
			"source_id": sourceID,
			"data":      payload,
		})
	if err != nil {
		return fmt.Errorf("unable to record audit history. err: %w", err)
	}

	return nil
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
