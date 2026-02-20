package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"budget-app/infrastructure"
	"budget-app/storage/migrations"
)

var wait bool

var rootCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Budget App database migration tool",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if wait {
			fmt.Println("waiting 10s for the network...")
			time.Sleep(10 * time.Second)
		}
		fmt.Println("connecting to database...")
		if err := infrastructure.Connect(); err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}
		fmt.Println("connected")
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		infrastructure.Close()
	},
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Run all pending migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		oldV, newV, err := migrations.Run(infrastructure.DB, "up")
		if err != nil {
			return err
		}
		if newV != oldV {
			fmt.Printf("migrated from version %d to %d\n", oldV, newV)
		} else {
			fmt.Printf("schema is up to date at version %d\n", newV)
		}
		return nil
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll back the last applied migration",
	RunE: func(cmd *cobra.Command, args []string) error {
		oldV, newV, err := migrations.Run(infrastructure.DB, "down")
		if err != nil {
			return err
		}
		fmt.Printf("rolled back from version %d to %d\n", oldV, newV)
		return nil
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Roll back all applied migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		oldV, newV, err := migrations.Run(infrastructure.DB, "reset")
		if err != nil {
			return err
		}
		fmt.Printf("reset from version %d to %d\n", oldV, newV)
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current schema version and available migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentVersion, _, err := migrations.Run(infrastructure.DB, "version")
		if err != nil {
			return err
		}
		all := migrations.Collection.Migrations()
		fmt.Printf("current version: %d\n\n", currentVersion)
		for _, m := range all {
			if m.Version == 0 {
				continue
			}
			state := "pending"
			if m.Version <= currentVersion {
				state = "applied"
			}
			fmt.Printf("  [%-7s] version %d\n", state, m.Version)
		}
		return nil
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current schema version",
	RunE: func(cmd *cobra.Command, args []string) error {
		v, _, err := migrations.Run(infrastructure.DB, "version")
		if err != nil {
			return err
		}
		fmt.Printf("current version: %d\n", v)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&wait, "wait", false, "wait 10s before connecting (useful in Docker Compose)")
	rootCmd.AddCommand(upCmd, downCmd, resetCmd, statusCmd, versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
