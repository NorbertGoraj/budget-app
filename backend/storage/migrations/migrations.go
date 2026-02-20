package migrations

import (
	"embed"
	"fmt"

	"github.com/go-pg/migrations/v8"
	"github.com/go-pg/pg/v10"
)

//go:embed *.sql
var sqlFiles embed.FS

// Collection is the ordered set of migrations exposed to the CLI.
var Collection = migrations.NewCollection()

func init() {
	register(
		"1_create_accounts",
		"2_create_transactions",
		"3_create_budgets",
		"4_create_planned_purchases",
		"5_create_investments",
	)
}

// register adds up/down migrations for each name in order.
func register(names ...string) {
	for _, name := range names {
		up := mustSQL(name + ".up.sql")
		down := mustSQL(name + ".down.sql")

		Collection.MustRegisterTx(
			func(db migrations.DB) error { _, err := db.Exec(up); return err },
			func(db migrations.DB) error { _, err := db.Exec(down); return err },
		)
	}
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
