package postgres

import "context"

// Executor describes the executor interface.
type Executor interface {
	// Exec abstracts the query execution. queryName is used for tracing and prepared statements.
	Exec(ctx context.Context, queryName, sql string, args ...interface{}) error
}
