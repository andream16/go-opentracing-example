package http_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	transporthttp "github.com/andream16/go-opentracing-example/src/http-server-receiver/transport/http"
	"github.com/andream16/go-opentracing-example/src/shared/tracing"
	todoclientmock "github.com/andream16/go-opentracing-example/src/test/mock/todoclient"
)

func TestNewHandler(t *testing.T) {
	t.Run("it should return an error because the todo client is invalid", func(t *testing.T) {
		handler, err := transporthttp.NewHandler(nil, nil)

		require.Error(t, err)
		var e transporthttp.InvalidHandlerParameterError
		assert.True(t, errors.As(err, &e))
		assert.Equal(t, "invalid parameter todoSvcClient: cannot be nil", err.Error())
		assert.Empty(t, handler)
	})
	t.Run("it should return an error because the tracer is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		handler, err := transporthttp.NewHandler(todoclientmock.NewMockTodoServiceClient(ctrl), nil)

		require.Error(t, err)
		var e transporthttp.InvalidHandlerParameterError
		assert.True(t, errors.As(err, &e))
		assert.Equal(t, "invalid parameter tracer: cannot be nil", err.Error())
		assert.Empty(t, handler)
	})
	t.Run("it should return a new handler", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		handler, err := transporthttp.NewHandler(todoclientmock.NewMockTodoServiceClient(ctrl), tracing.JaegerTracer{})

		require.NoError(t, err)
		assert.NotEmpty(t, handler)
	})
}

func TestHandler_Router(t *testing.T) {
	t.Run("it should return the internal mux router", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		handler, err := transporthttp.NewHandler(todoclientmock.NewMockTodoServiceClient(ctrl), tracing.JaegerTracer{})

		require.NoError(t, err)
		assert.NotEmpty(t, handler, handler.Router())
	})
}
