package migrations

import (
	"context"
	"embed"
	"fmt"
	"log"
	"sort"
	"strings"

	"budget-app/infrastructure"

	"github.com/go-pg/pg/v10"
)

//go:embed *.sql
var migrationFiles embed.FS

const createTrackingTable = `
CREATE TABLE IF NOT EXISTS schema_migrations (
	version    TEXT        PRIMARY KEY,
	applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
)`

// Run checks which migrations have already been applied and runs only the
// pending ones. Each migration is recorded in schema_migrations on success.
func Run() error {
	if _, err := infrastructure.DB.Exec(createTrackingTable); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	applied, err := loadApplied()
	if err != nil {
		return err
	}

	pending, err := pendingMigrations(applied)
	if err != nil {
		return err
	}

	if len(pending) == 0 {
		log.Println("migrations: schema is up to date")
		return nil
	}

	for _, name := range pending {
		if err := apply(name); err != nil {
			return err
		}
	}

	log.Printf("migrations: %d migration(s) applied", len(pending))
	return nil
}

// Version returns the number of migrations that have been applied.
func Version() (int, error) {
	if _, err := infrastructure.DB.Exec(createTrackingTable); err != nil {
		return 0, fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	var count int
	_, err := infrastructure.DB.QueryOne(pg.Scan(&count), "SELECT COUNT(*) FROM schema_migrations")
	if err != nil {
		return 0, fmt.Errorf("failed to query migration version: %w", err)
	}
	return count, nil
}

// Status prints the state of every migration file: applied or pending.
func Status() error {
	if _, err := infrastructure.DB.Exec(createTrackingTable); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	applied, err := loadApplied()
	if err != nil {
		return err
	}

	entries, err := sortedEntries()
	if err != nil {
		return err
	}

	log.Println("migrations status:")
	for _, name := range entries {
		if ts, ok := applied[name]; ok {
			log.Printf("  [applied]  %s (at %s)", name, ts)
		} else {
			log.Printf("  [pending]  %s", name)
		}
	}
	return nil
}

// loadApplied returns a map of version → applied_at timestamp string.
func loadApplied() (map[string]string, error) {
	var rows []struct {
		Version   string
		AppliedAt string `pg:"applied_at"`
	}
	_, err := infrastructure.DB.Query(&rows,
		"SELECT version, applied_at::text FROM schema_migrations ORDER BY applied_at")
	if err != nil {
		return nil, fmt.Errorf("failed to query schema_migrations: %w", err)
	}

	applied := make(map[string]string, len(rows))
	for _, r := range rows {
		applied[r.Version] = r.AppliedAt
	}
	return applied, nil
}

// pendingMigrations returns sorted file names that are not yet in applied.
func pendingMigrations(applied map[string]string) ([]string, error) {
	all, err := sortedEntries()
	if err != nil {
		return nil, err
	}

	var pending []string
	for _, name := range all {
		if _, ok := applied[name]; !ok {
			pending = append(pending, name)
		}
	}
	return pending, nil
}

// sortedEntries returns all .sql file names from the embedded FS, sorted.
func sortedEntries() ([]string, error) {
	entries, err := migrationFiles.ReadDir(".")
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	var names []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			names = append(names, e.Name())
		}
	}
	return names, nil
}

// apply executes a single migration file and records it in schema_migrations.
func apply(name string) error {
	content, err := migrationFiles.ReadFile(name)
	if err != nil {
		return fmt.Errorf("failed to read migration %s: %w", name, err)
	}

	err = infrastructure.DB.RunInTransaction(context.Background(), func(tx *pg.Tx) error {
		if _, err := tx.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", name, err)
		}
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES (?)", name); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", name, err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("migrations: applied %s", name)
	return nil
}
