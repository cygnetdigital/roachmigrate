package migrate

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/cygnetdigital/roachmigrate/pkg/store"
	"github.com/hashicorp/go-multierror"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Migrate is a mechanism to check and run migrations
type Migrate struct {
	storedb  StoreConn
	model    *store.Queries
	targetdb Execer
	key      string
	files    []*File
}

func New(storedb StoreConn, targetdb Execer, key string, files []*File) *Migrate {
	return &Migrate{
		storedb: storedb,
		model:   store.New(),
		key:     key,
		files:   sortFiles(files),
	}
}

// Execer is an interface for executing sql statements
type Execer interface {
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
}

// StoreConn is an interface for a connection to the store
type StoreConn interface {
	store.DBTX
	Begin(ctx context.Context) (pgx.Tx, error)
}

// Status returns the status of migrations
func (m *Migrate) Status(ctx context.Context) (*Status, error) {
	return m.status(ctx, m.storedb)
}

func (m *Migrate) status(ctx context.Context, tx store.DBTX) (*Status, error) {
	rows, err := m.model.ListForUpdate(ctx, tx, m.key)
	if err != nil {
		if store.IsNotInitialized(err) {
			return nil, ErrStoreNotInitialized
		}
		return nil, err
	}

	zipped := zip(m.files, rows)
	hasErrors := false
	inSync := true
	var next *Migration

	// warn for missing files where there are rows
	errs := &multierror.Error{}

	for _, z := range zipped {
		if z.File == nil && z.Row != nil {
			errs = multierror.Append(errs, fmt.Errorf("missing file for migration %s", z.Row.Filename))
		}

		if z.Row != nil {
			if z.Row.Failed {
				hasErrors = true
			}
			if !z.Row.Completed {
				inSync = false
			}
		}

		if z.File != nil && z.Row == nil {
			inSync = false
			if next == nil {
				next = z
			}
		}
	}

	if err := errs.ErrorOrNil(); err != nil {
		return nil, err
	}

	return &Status{
		Key:           m.key,
		Migrations:    zipped,
		InSync:        inSync,
		HasErrors:     hasErrors,
		NextMigration: next,
	}, nil
}

// Run the next migration
func (m *Migrate) Run(ctx context.Context, filename string) (*Migration, error) {
	tx, err := m.storedb.Begin(ctx)
	if err != nil {
		return nil, err
	}

	//nolint:errcheck
	defer tx.Rollback(ctx)

	status, err := m.status(ctx, tx)
	if err != nil {
		return nil, err
	}

	if status.InSync {
		return nil, fmt.Errorf("no migrations to run")
	}

	if status.HasErrors {
		return nil, fmt.Errorf("cannot run migrations with errors")
	}

	if status.NextMigration.Filename != filename {
		return nil, fmt.Errorf("next migration is %s", status.NextMigration.Filename)
	}

	next := status.NextMigration

	row, err := m.model.Create(ctx, tx, store.CreateParams{
		ID:       fmt.Sprintf("%s/%s", m.key, next.Filename),
		Key:      m.key,
		Filename: next.Filename,
	})
	if err != nil {
		return nil, err
	}

	next.Row = &row

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// reaching here means we have obtained a lease to run the next migration.
	if err := m.runMigration(ctx, next); err != nil {
		return next, err
	}

	return next, nil
}

func (m *Migrate) runMigration(ctx context.Context, next *Migration) error {

	up := store.UpdateParams{
		ID: next.Row.ID,
	}

	_, err := m.targetdb.Exec(ctx, string(next.File.Contents))
	if err != nil {
		up.Failed = true
		up.FailReason = sql.NullString{Valid: true, String: err.Error()}
	} else {
		up.Completed = true
	}

	row, updateErr := m.model.Update(ctx, m.storedb, up)
	if updateErr != nil {
		return &migrationUpdateError{origErr: err, updatedErr: updateErr}
	}

	next.Row = &row
	return err
}

type migrationUpdateError struct {
	origErr    error
	updatedErr error
}

func (e *migrationUpdateError) Error() string {
	if e.origErr == nil {
		return fmt.Sprintf("successfully ran migration but failed to record status: %s", e.updatedErr.Error())
	}

	return fmt.Sprintf("failed to both run and record migration:\n\trunner: %s\n\trecorder: %s", e.origErr.Error(), e.updatedErr.Error())
}
