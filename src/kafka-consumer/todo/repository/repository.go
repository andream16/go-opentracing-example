package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/andream16/go-opentracing-example/src/shared/database/postgres"
)

const (
	todosTableName              = "todos"
	todosTableMessageColumnName = "message"
)

// Creator describes the creator interface.
type Creator interface {
	// Create creates a new todo.
	Create(ctx context.Context, todo *Todo) error
}

// TodoCreator is the todos repository.
type TodoCreator struct {
	executor postgres.Executor
}

// New returns a new TodoCreator.
func New(executor postgres.Executor) (TodoCreator, error) {
	if executor == nil {
		return TodoCreator{}, errors.New("executor cannot be nil")
	}
	return TodoCreator{
		executor: executor,
	}, nil
}

// Create inserts a new todo in the todosTableName table.
func (tc TodoCreator) Create(ctx context.Context, todo *Todo) error {
	const createTodosQueryName = "create_todos"

	if err := tc.executor.Exec(
		ctx,
		createTodosQueryName,
		`INSERT INTO $1($2) VALUES($3)`,
		todosTableName,
		todosTableMessageColumnName,
		todo.Message,
	); err != nil {
		return fmt.Errorf("could not insert todo: %w", err)
	}

	return nil
}
