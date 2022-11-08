package migrate

import (
	"fmt"

	"github.com/cygnetdigital/roachmigrate/pkg/store"
)

var (
	ErrStoreNotInitialized = fmt.Errorf("store not initialized")
)

// Status of all migrations
type Status struct {
	Key        string
	Migrations []*Migration

	// All migrations are complete
	InSync bool

	// Previous migrations have errors
	HasErrors bool

	// NextMigraiton to run (nil if none)
	NextMigration *Migration
}

// Migration represents a single migration
type Migration struct {
	// Filename from row, file or both
	Filename string

	// Database row if present
	Row *store.RoachMigration

	// File if present
	File *File
}

// File represents a single migration file
type File struct {
	Name     string
	Contents []byte
}
