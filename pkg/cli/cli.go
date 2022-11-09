package cli

import (
	"context"
	"fmt"

	"github.com/cygnetdigital/roachmigrate/pkg/migrate"
	"github.com/cygnetdigital/roachmigrate/pkg/store"
)

// CLI wraps the core migration logic to provide CLI based functionality
type CLI struct {
	m *migrate.Migrate
}

// New CLI from a migrate instance
func New(m *migrate.Migrate) *CLI {
	return &CLI{m}
}

// Status prints the current status of all migrations
func (cli *CLI) Status(ctx context.Context) error {
	status, err := cli.m.Status(ctx)
	if err != nil {
		return err
	}

	printMigrations(status.Migrations)

	if status.InSync {
		fmt.Println("✅ in sync")
	}

	if status.HasErrors {
		fmt.Println("❌ migrations have errors")
	}

	return nil
}

// Run the given migration
func (cli *CLI) Run(ctx context.Context, m string) error {
	if m == "" {
		return fmt.Errorf("migration name is required")
	}

	migration, err := cli.m.Run(ctx, m)
	if err != nil {
		return err
	}

	printMigrations([]*migrate.Migration{migration})

	return nil
}

// Validate the given migration
func (cli *CLI) Validate(ctx context.Context) error {
	if len(cli.m.Files) == 0 {
		return fmt.Errorf("no migrations found")
	}

	for _, f := range cli.m.Files {
		if err := migrate.Validate(f); err != nil {
			// red cross emoji
			fmt.Printf("❌ %s - %s\n", f.Name, err)
		} else {
			fmt.Printf("✅ %s\n", f.Name)
		}
	}

	return nil
}

// Init the migrations table
func (cli *CLI) Init(ctx context.Context) error {
	if err := store.Init(ctx, cli.m.StoreConn); err != nil {
		return fmt.Errorf("failed to init migrate store: %w", err)
	}

	fmt.Println("done!")
	return nil
}
