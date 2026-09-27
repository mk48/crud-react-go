package migrations

import (
	"database/sql"
	"embed"

	migrate "github.com/rubenv/sql-migrate"
)

//go:embed *
var migrationsFiles embed.FS

func Run(db *sql.DB) error {

	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: migrationsFiles,
		Root:       ".",
	}

	if _, err := migrate.Exec(db, "postgres", migrations, migrate.Up); err != nil {
		return err
	}

	return nil
}
