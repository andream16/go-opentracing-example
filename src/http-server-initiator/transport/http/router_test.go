package http_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	transporthttp "github.com/andream16/go-opentracing-example/src/http-server-initiator/transport/http"
	"github.com/andream16/go-opentracing-example/src/shared/tracing"
)

func TestNewHandler(t *testing.T) {
	t.Run("it should return an error because the receiver host is invalid", func(t *testing.T) {
		handler, err := transporthttp.NewHandler("", nil, nil)

		require.Error(t, err)
		var e transporthttp.InvalidHandlerParameterError
		assert.True(t, errors.As(err, &e))
		assert.Equal(t, "invalid parameter receiverHostname: cannot be empty", err.Error())
		assert.Empty(t, handler)
	})
	t.Run("it should return an error because the http client is invalid", func(t *testing.T) {
		handler, err := transporthttp.NewHandler("someHostName", nil, nil)

		require.Error(t, err)
		var e transporthttp.InvalidHandlerParameterError
		assert.True(t, errors.As(err, &e))
		assert.Equal(t, "invalid parameter httpClient: cannot be nil", err.Error())
		assert.Empty(t, handler)
	})
	t.Run("it should return an error because the tracer is invalid", func(t *testing.T) {
		handler, err := transporthttp.NewHandler("someHostName", &http.Client{}, nil)

		require.Error(t, err)
		var e transporthttp.InvalidHandlerParameterError
		assert.True(t, errors.As(err, &e))
		assert.Equal(t, "invalid parameter tracer: cannot be nil", err.Error())
		assert.Empty(t, handler)
	})
	t.Run("it should return a new handler", func(t *testing.T) {
		handler, err := transporthttp.NewHandler(
			"someHostName",
			&http.Client{},
			tracing.JaegerTracer{},
		)

		require.NoError(t, err)
		assert.NotEmpty(t, handler)
	})
}

func TestHandler_Router(t *testing.T) {
	t.Run("it should return the internal mux router", func(t *testing.T) {
		handler, err := transporthttp.NewHandler(
			"someHostName",
			&http.Client{},
			tracing.JaegerTracer{},
		)

		require.NoError(t, err)
		assert.NotEmpty(t, handler, handler.Router())
	})
}
