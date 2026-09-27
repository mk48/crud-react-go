package util

import (
	"context"
	"log/slog"
	"maps"
	"strings"

	"kfamily/internal/dto"

	"github.com/jmoiron/sqlx"
)

// TableColumns returns the column names and Postgres data types for tableName,
// in table definition order. Lets list endpoints expose a table's shape so a
// frontend can discover which columns it's allowed to sort/filter on.
//
// tableName may be double-quoted (e.g. `"user"`, needed elsewhere to embed it
// raw in SQL when it's a reserved word) - the quotes are stripped here since
// information_schema stores the bare name.
func TableColumns(ctx context.Context, db *sqlx.DB, tableName string) ([]dto.ColumnMeta, error) {
	query := `
		SELECT column_name, data_type
		FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = $1
		ORDER BY ordinal_position`

	columns := []dto.ColumnMeta{}
	if err := db.SelectContext(ctx, &columns, query, strings.Trim(tableName, `"`)); err != nil {
		return nil, err
	}

	return columns, nil
}

// ToCamelCase converts a snake_case Postgres column name (e.g. "created_at")
// to the camelCase form frontend DTOs use (e.g. "createdAt"), so a table's
// physical column names can be matched against the sort keys a JSON API
// client actually sends.
func ToCamelCase(snakeCase string) string {
	parts := strings.Split(snakeCase, "_")
	for i := 1; i < len(parts); i++ {
		if parts[i] == "" {
			continue
		}
		parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
	}
	return strings.Join(parts, "")
}

// BuildSortableColumns fetches tableName's columns (via TableColumns) and
// maps each one's camelCase form to itself, for use as
// middleware.PaginationFilter's sortableColumns - the frontend sends sort
// keys matching its DTO field names, which need mapping to the real
// snake_case Postgres column to order by. extra adds (or overrides) further
// sort keys on top, e.g. a composite UI-only column like "actionAt" that has
// no single matching DB column of its own.
//
// Logs and returns just extra if the table's metadata can't be loaded,
// rather than failing service startup.
func BuildSortableColumns(ctx context.Context, log *slog.Logger, db *sqlx.DB, tableName string, extra map[string]string) map[string]string {
	sortable := map[string]string{}

	columns, err := TableColumns(ctx, db, tableName)
	if err != nil {
		log.Error("Error in getting metadata of table "+tableName, "err", err)
	} else {
		for _, c := range columns {
			sortable[ToCamelCase(c.Name)] = c.Name
		}
	}

	maps.Copy(sortable, extra)

	return sortable
}
