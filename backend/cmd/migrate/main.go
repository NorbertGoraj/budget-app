package main

import (
	"fmt"
	"os"
	"time"

	"budget-app/infrastructure"
	"budget-app/storage/migrations"

	"github.com/spf13/cobra"
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
		fmt.Println("running pending migrations...")
		return migrations.Run()
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the state of every migration file (applied / pending)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return migrations.Status()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the number of applied migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		v, err := migrations.Version()
		if err != nil {
			return err
		}
		fmt.Printf("applied migrations: %d\n", v)
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&wait, "wait", false, "wait 10s before connecting (useful in Docker Compose)")
	rootCmd.AddCommand(upCmd, statusCmd, versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
