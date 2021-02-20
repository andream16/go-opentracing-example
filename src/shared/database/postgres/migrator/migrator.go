package migrator

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
)

// Migrator is the migrator interface.
type Migrator interface {
	Migrate(ctx context.Context) error
}

// PgxMigrator wraps a pgx migrator.
type PgxMigrator struct {
	migrator *migrate.Migrator
}

// NewPgxMigrator returns a new PgxMigrator.
func NewPgxMigrator(ctx context.Context, conn *pgx.Conn, versionTable string) (PgxMigrator, error) {
	if conn == nil {
		return PgxMigrator{}, errors.New("pgx connection cannot be nil")
	}
	m, err := migrate.NewMigrator(ctx, conn, versionTable)
	if err != nil {
		return PgxMigrator{}, fmt.Errorf("could not create a new migrator: %w", err)
	}

	return PgxMigrator{migrator: m}, nil
}

// AppendMigration appends a migration to the migrator.
func (pm PgxMigrator) AppendMigration(name, upQuery, downQuery string) {
	pm.migrator.AppendMigration(name, upQuery, downQuery)
}

// Migrate runs the migration.
func (pm PgxMigrator) Migrate(ctx context.Context) error {
	return pm.migrator.Migrate(ctx)
}
