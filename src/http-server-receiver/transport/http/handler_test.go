package http_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	todov1 "github.com/andream16/go-opentracing-example/contracts/build/go/go_opentracing_example/grpc_server/todo/v1"
	transporthttp "github.com/andream16/go-opentracing-example/src/http-server-receiver/transport/http"
	opentracingmock "github.com/andream16/go-opentracing-example/src/test/mock/opentracing"
	todoclientmock "github.com/andream16/go-opentracing-example/src/test/mock/todoclient"
	tracingmock "github.com/andream16/go-opentracing-example/src/test/mock/tracing"
)

func TestHandler_CreateTodo(t *testing.T) {
	t.Run("it should return http.StatusBadRequest because the payload is malformed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			mockTodoClient  = todoclientmock.NewMockTodoServiceClient(ctrl)
			mockTracer      = tracingmock.NewMockTracer(ctrl)
			mockSpan        = opentracingmock.NewMockSpan(ctrl)
			mockSpanContext = opentracingmock.NewMockSpanContext(ctrl)
			req             = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{`))
			recorder        = httptest.NewRecorder()
		)

		handler, err := transporthttp.NewHandler(mockTodoClient, mockTracer)
		require.NoError(t, err)

		gomock.InOrder(
			mockTracer.EXPECT().Extract(opentracing.HTTPHeaders, gomock.Any()).Return(mockSpanContext, nil).Times(1),
			mockTracer.EXPECT().StartSpan("receiver_todo", ext.RPCServerOption(mockSpanContext)).Return(mockSpan).Times(1),
			mockSpan.EXPECT().Finish().Times(1),
		)

		handler.CreateTodo(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
	})
	t.Run("it should return http.StatusServiceUnavailable because the request failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			mockTodoClient  = todoclientmock.NewMockTodoServiceClient(ctrl)
			mockTracer      = tracingmock.NewMockTracer(ctrl)
			mockSpan        = opentracingmock.NewMockSpan(ctrl)
			mockSpanContext = opentracingmock.NewMockSpanContext(ctrl)
			req             = httptest.NewRequest(
				http.MethodPost,
				"/",
				bytes.NewBufferString(`{"message" : "hey there"}`),
			)
			recorder = httptest.NewRecorder()
		)

		handler, err := transporthttp.NewHandler(mockTodoClient, mockTracer)
		require.NoError(t, err)

		gomock.InOrder(
			mockTracer.EXPECT().Extract(opentracing.HTTPHeaders, gomock.Any()).Return(mockSpanContext, nil).Times(1),
			mockTracer.EXPECT().StartSpan("receiver_todo", ext.RPCServerOption(mockSpanContext)).Return(mockSpan).Times(1),
			mockSpan.EXPECT().Tracer().Times(1),
			mockTodoClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("someErr")).Times(1),
			mockSpan.EXPECT().Finish().Times(1),
		)

		handler.CreateTodo(recorder, req)

		assert.Equal(t, http.StatusServiceUnavailable, recorder.Result().StatusCode)
	})
	t.Run("it should return http.StatusOK because the request was successful", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			mockTodoClient  = todoclientmock.NewMockTodoServiceClient(ctrl)
			mockTracer      = tracingmock.NewMockTracer(ctrl)
			mockSpan        = opentracingmock.NewMockSpan(ctrl)
			mockSpanContext = opentracingmock.NewMockSpanContext(ctrl)
			req             = httptest.NewRequest(
				http.MethodPost,
				"/",
				bytes.NewBufferString(`{"message" : "hey there"}`),
			)
			recorder = httptest.NewRecorder()
		)

		handler, err := transporthttp.NewHandler(mockTodoClient, mockTracer)
		require.NoError(t, err)

		gomock.InOrder(
			mockTracer.EXPECT().Extract(opentracing.HTTPHeaders, gomock.Any()).Return(mockSpanContext, nil).Times(1),
			mockTracer.EXPECT().StartSpan("receiver_todo", ext.RPCServerOption(mockSpanContext)).Return(mockSpan).Times(1),
			mockSpan.EXPECT().Tracer().Times(1),
			mockTodoClient.EXPECT().Create(gomock.Any(), gomock.Any()).Return(&todov1.CreateResponse{}, nil).Times(1),
			mockSpan.EXPECT().Finish().Times(1),
		)

		handler.CreateTodo(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
	})
}
