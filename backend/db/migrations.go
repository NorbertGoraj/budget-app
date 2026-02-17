package db

import (
	"context"
	"embed"
	"fmt"
	"log"
	"sort"
	"strings"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func RunMigrations() error {
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort files by name to ensure execution order
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		content, err := migrationFiles.ReadFile("migrations/" + entry.Name())
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", entry.Name(), err)
		}

		log.Printf("Running migration: %s", entry.Name())
		if _, err := Pool.Exec(context.Background(), string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", entry.Name(), err)
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}
