package http_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	transporthttp "github.com/andream16/go-opentracing-example/src/http-server-initiator/transport/http"
	opentracingmock "github.com/andream16/go-opentracing-example/src/test/mock/opentracing"
	tracingmock "github.com/andream16/go-opentracing-example/src/test/mock/tracing"
	transporthttpmock "github.com/andream16/go-opentracing-example/src/test/mock/transport/http"
)

func TestHandler_CreateTodo(t *testing.T) {
	t.Run("it should return http.StatusBadRequest because the request body is invalid", func(t *testing.T) {
		const someHostname = "http://hello:8080"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			mockTracer = tracingmock.NewMockTracer(ctrl)
			mockSpan   = opentracingmock.NewMockSpan(ctrl)
			req        = httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{`))
			recorder   = httptest.NewRecorder()
		)

		handler, err := transporthttp.NewHandler(someHostname, &http.Client{}, mockTracer)
		require.NoError(t, err)

		gomock.InOrder(
			mockTracer.EXPECT().StartSpan("initiator_todo").Return(mockSpan).Times(1),
			mockSpan.EXPECT().Finish().Times(1),
		)

		handler.CreateTodo(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
	})
	t.Run("it should return http.StatusServiceUnavailable because the request failed", func(t *testing.T) {
		const someHostname = "http://hello:8080"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			mockDoer        = transporthttpmock.NewMockDoer(ctrl)
			mockTracer      = tracingmock.NewMockTracer(ctrl)
			mockSpan        = opentracingmock.NewMockSpan(ctrl)
			mockSpanContext = opentracingmock.NewMockSpanContext(ctrl)
			req             = httptest.NewRequest(
				http.MethodPost,
				"/",
				bytes.NewBufferString(`{"message" : "hello"}`),
			)
			recorder = httptest.NewRecorder()
			resp     = &http.Response{
				StatusCode: http.StatusTeapot,
			}
		)

		handler, err := transporthttp.NewHandler(someHostname, mockDoer, mockTracer)
		require.NoError(t, err)

		gomock.InOrder(
			mockTracer.EXPECT().StartSpan("initiator_todo").Return(mockSpan).Times(1),
			mockSpan.EXPECT().SetTag("http.url", someHostname+"/receiver/todo").Times(1),
			mockSpan.EXPECT().SetTag("http.method", http.MethodPost).Times(1),
			mockSpan.EXPECT().Context().Return(mockSpanContext).Times(1),
			mockTracer.EXPECT().
				Inject(
					mockSpanContext,
					opentracing.HTTPHeaders,
					opentracing.HTTPHeadersCarrier(req.Header),
				).Return(nil).Times(1),
			mockDoer.EXPECT().Do(gomock.Any()).Return(resp, nil),
			mockSpan.EXPECT().Finish().Times(1),
		)

		handler.CreateTodo(recorder, req)

		assert.Equal(t, http.StatusServiceUnavailable, recorder.Result().StatusCode)
	})
	t.Run("it should return http.StatusOK because a todo has been created", func(t *testing.T) {
		const someHostname = "http://hello:8080"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			mockDoer        = transporthttpmock.NewMockDoer(ctrl)
			mockTracer      = tracingmock.NewMockTracer(ctrl)
			mockSpan        = opentracingmock.NewMockSpan(ctrl)
			mockSpanContext = opentracingmock.NewMockSpanContext(ctrl)
			req             = httptest.NewRequest(
				http.MethodPost,
				"/",
				bytes.NewBufferString(`{"message" : "hello"}`),
			)
			recorder = httptest.NewRecorder()
			resp     = &http.Response{
				StatusCode: http.StatusOK,
			}
		)

		handler, err := transporthttp.NewHandler(someHostname, mockDoer, mockTracer)
		require.NoError(t, err)

		gomock.InOrder(
			mockTracer.EXPECT().StartSpan("initiator_todo").Return(mockSpan).Times(1),
			mockSpan.EXPECT().SetTag("http.url", someHostname+"/receiver/todo").Times(1),
			mockSpan.EXPECT().SetTag("http.method", http.MethodPost).Times(1),
			mockSpan.EXPECT().Context().Return(mockSpanContext).Times(1),
			mockTracer.EXPECT().
				Inject(
					mockSpanContext,
					opentracing.HTTPHeaders,
					opentracing.HTTPHeadersCarrier(req.Header),
				).Return(nil).Times(1),
			mockDoer.EXPECT().Do(gomock.Any()).Return(resp, nil),
			mockSpan.EXPECT().Finish().Times(1),
		)

		handler.CreateTodo(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
	})
}
