package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed *.sql
var migrationsFiles embed.FS

// lockKey is an arbitrary app-wide id for the Postgres advisory lock that
// serializes migrations.
const lockKey = 7_311_424_001

// Run applies pending migrations. It holds a Postgres advisory lock while
// doing so, so replicas starting at the same time run them one at a time:
// the first applies them, the rest wait and then find nothing to do.
func Run(ctx context.Context, db *sql.DB) error {
	// Session-level lock: it must be taken and released on one connection.
	conn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("unable to get connection for migration lock: %w", err)
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `SELECT pg_advisory_lock($1)`, lockKey); err != nil {
		return fmt.Errorf("unable to take migration lock: %w", err)
	}
	defer conn.ExecContext(context.Background(), `SELECT pg_advisory_unlock($1)`, lockKey)

	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrationsFiles,
		Root:       ".",
	}

	if _, err := migrate.ExecContext(ctx, db, "postgres", migrations, migrate.Up); err != nil {
		return err
	}

	return nil
}
