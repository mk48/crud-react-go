package util

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
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
	if _, err := db.NamedExecContext(ctx, query, params); err != nil {
		return fmt.Errorf("unable to insert into %s. err: %w", tableName, err)
	}

	return RecordAudit(ctx, db, id, params)
}

// UpdateByID builds "UPDATE tableName SET <col = :col, ...> WHERE id = :id"
// from params' keys (excluding "id" from the SET list), executes it, and
// records the row in audit_history (see RecordAudit). params must include
// "id" set to id.
func UpdateByID(ctx context.Context, db Execer, tableName string, id uuid.UUID, params map[string]any) error {
	cols := sortedKeys(params)

	setClauses := make([]string, 0, len(cols))
	for _, c := range cols {
		if c == "id" {
			continue
		}
		setClauses = append(setClauses, fmt.Sprintf("%s = :%s", c, c))
	}

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = :id`, tableName, strings.Join(setClauses, ", "))
	if _, err := db.NamedExecContext(ctx, query, params); err != nil {
		return fmt.Errorf("unable to update %s. err: %w", tableName, err)
	}

	return RecordAudit(ctx, db, id, params)
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
