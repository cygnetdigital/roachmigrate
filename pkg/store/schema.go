package store

import (
	"context"
	"embed"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
)

//go:embed schema.sql
var schema embed.FS

type SchemaExecer interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
}

// Init schema
func Init(ctx context.Context, conn SchemaExecer) error {
	bts, err := schema.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("failed to read schema: %w", err)
	}

	if _, err := conn.Exec(ctx, string(bts)); err != nil {
		return fmt.Errorf("error executing schema creation: %w", err)
	}

	return nil
}

// migrationsTable is the name of the table used to store migrations
const migrationsTable = "roach_migrations"

// IsNotInitialized returns true if the error is due to the migrations table not existing
func IsNotInitialized(err error) bool {
	return err != nil && strings.Contains(err.Error(), fmt.Sprintf("relation %q does not exist", migrationsTable))
}
