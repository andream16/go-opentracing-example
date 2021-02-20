package pgxwrapper

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/opentracing/opentracing-go"
)

// PgxWrapper is a wrapper to jackc/pgx/v4.
type PgxWrapper struct {
	tracer opentracing.Tracer
	pool   *pgxpool.Pool
}

// New returns a new PgxWrapper given a postgresql dsn.
// The wrapper has built in tracing.
// The connection will be retried until completion.
func New(
	ctx context.Context,
	dsn string,
	waitFor time.Duration,
	tracer opentracing.Tracer,
) (PgxWrapper, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return PgxWrapper{}, fmt.Errorf("could not create new connection configuration: %w", err)
	}

	pool, err := newPgxPool(ctx, cfg, waitFor)
	if err != nil {
		return PgxWrapper{}, fmt.Errorf("could not create new connection pool: %w", err)
	}

	return PgxWrapper{
		pool:   pool,
		tracer: tracer,
	}, nil
}

// Exec is pgx's concrete implementation for executing a query with tracing.
func (p PgxWrapper) Exec(ctx context.Context, queryName, sql string, args ...interface{}) error {
	span, _ := opentracing.StartSpanFromContext(ctx, queryName)
	span.Finish()

	// TODO implement me!

	return nil
}

func newPgxPool(ctx context.Context, config *pgxpool.Config, waitFor time.Duration) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		time.Sleep(waitFor)
		return newPgxPool(ctx, config, waitFor)
	}

	return pool, nil
}
