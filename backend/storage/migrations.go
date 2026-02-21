package storage

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
)

//go:embed migrations/*.sql
var sqlFiles embed.FS

// Collection is the ordered set of migrations exposed to the CLI.
var Collection = migrations.NewCollection()

func init() {
	if err := Collection.DiscoverSQLMigrationsFromFilesystem(http.FS(sqlFiles), "migrations"); err != nil {
		panic(fmt.Sprintf("migrations: %v", err))
	}
}

// Run delegates to the collection, returning the old and new schema versions.
func Run(db *pg.DB, args ...string) (oldVersion, newVersion int64, err error) {
	return Collection.Run(db, args...)
}
