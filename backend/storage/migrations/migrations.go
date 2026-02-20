package migrations

import (
	"embed"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
)

//go:embed *.sql
var sqlFiles embed.FS

// Collection is the ordered set of migrations exposed to the CLI.
var Collection = migrations.NewCollection()

func init() {
	if err := autoRegister(); err != nil {
		panic(fmt.Sprintf("migrations: %v", err))
	}
}

// autoRegister discovers all *.up.sql files from the embedded FS, sorts them
// by their numeric prefix, verifies each has a matching *.down.sql, then
// registers the pair. Adding a new migration requires only adding the two SQL
// files — no changes to Go code are needed.
func autoRegister() error {
	entries, err := sqlFiles.ReadDir(".")
	if err != nil {
		return fmt.Errorf("failed to read embedded SQL files: %w", err)
	}

	var upFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}

	sort.Slice(upFiles, func(i, j int) bool {
		return versionPrefix(upFiles[i]) < versionPrefix(upFiles[j])
	})

	for _, upFile := range upFiles {
		downFile := strings.TrimSuffix(upFile, ".up.sql") + ".down.sql"
		if _, err := sqlFiles.Open(downFile); err != nil {
			return fmt.Errorf("missing down migration for %s (expected %s)", upFile, downFile)
		}

		up := mustSQL(upFile)
		down := mustSQL(downFile)

		Collection.MustRegisterTx(
			func(db migrations.DB) error { _, err := db.Exec(up); return err },
			func(db migrations.DB) error { _, err := db.Exec(down); return err },
		)
	}

	return nil
}

// versionPrefix parses the leading integer from a filename like "3_create_budgets.up.sql".
func versionPrefix(filename string) int64 {
	part, _, _ := strings.Cut(filename, "_")
	v, _ := strconv.ParseInt(part, 10, 64)
	return v
}

// mustSQL reads an embedded SQL file and panics if it is missing.
func mustSQL(name string) string {
	b, err := sqlFiles.ReadFile(name)
	if err != nil {
		panic(fmt.Sprintf("migrations: missing SQL file %s: %v", name, err))
	}
	return string(b)
}

// Run delegates to the collection, returning the old and new schema versions.
func Run(db *pg.DB, args ...string) (oldVersion, newVersion int64, err error) {
	return Collection.Run(db, args...)
}
