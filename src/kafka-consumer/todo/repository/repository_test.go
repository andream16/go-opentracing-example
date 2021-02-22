package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andream16/go-opentracing-example/src/kafka-consumer/todo/repository"
	"github.com/andream16/go-opentracing-example/src/shared/todo"
	executormock "github.com/andream16/go-opentracing-example/src/test/mock/database/postgres"
)

func TestNew(t *testing.T) {
	t.Run("it should return an error because the executor is not valid", func(t *testing.T) {
		creator, err := repository.New(nil)
		require.Error(t, err)
		assert.Empty(t, creator)
	})
	t.Run("it should return a new creator", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		creator, err := repository.New(executormock.NewMockExecutor(ctrl))
		require.NoError(t, err)
		assert.NotEmpty(t, creator)
	})
}

func TestTodoCreator_Create(t *testing.T) {
	t.Run("it should return an error because the execution of the query failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		const (
			queryName   = "create_todos"
			todoMessage = "hello"
		)

		var (
			ctx          = context.Background()
			executorMock = executormock.NewMockExecutor(ctrl)
		)

		creator, err := repository.New(executorMock)
		require.NoError(t, err)
		assert.NotEmpty(t, creator)

		executorMock.EXPECT().Exec(
			ctx,
			queryName,
			`INSERT INTO todos(message) VALUES($1::text)`,
			todoMessage,
		).Return(errors.New("someErr")).Times(1)

		require.Error(t, creator.Create(ctx, &todo.Todo{Message: todoMessage}))
	})
	t.Run("it should create a todo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		const (
			queryName   = "create_todos"
			todoMessage = "hello"
		)

		var (
			ctx          = context.Background()
			executorMock = executormock.NewMockExecutor(ctrl)
		)

		creator, err := repository.New(executorMock)
		require.NoError(t, err)
		assert.NotEmpty(t, creator)

		executorMock.EXPECT().Exec(
			ctx,
			queryName,
			`INSERT INTO todos(message) VALUES($1::text)`,
			todoMessage,
		).Return(nil).Times(1)

		require.NoError(t, creator.Create(ctx, &todo.Todo{Message: todoMessage}))
	})
}
