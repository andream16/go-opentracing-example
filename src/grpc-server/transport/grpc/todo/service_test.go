package todo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	todov1 "github.com/andream16/go-opentracing-example/contracts/build/go/go_opentracing_example/grpc_server/todo/v1"
	"github.com/andream16/go-opentracing-example/src/grpc-server/transport/grpc/todo"
	sendermock "github.com/andream16/go-opentracing-example/src/test/mock/kafka"
	tracingmock "github.com/andream16/go-opentracing-example/src/test/mock/tracing"
)

func TestNewService(t *testing.T) {
	t.Run("it should return an error because the topic is not valid", func(t *testing.T) {
		svc, err := todo.NewService("", nil, nil)

		require.Error(t, err)
		var e todo.InvalidServiceParameterError
		require.True(t, errors.As(err, &e))
		assert.Equal(t, "invalid parameter kafka topic: must be not empty", err.Error())
		assert.Empty(t, svc)
	})
	t.Run("it should return an error because the sender is not valid", func(t *testing.T) {
		svc, err := todo.NewService("someTopic", nil, nil)

		require.Error(t, err)
		var e todo.InvalidServiceParameterError
		require.True(t, errors.As(err, &e))
		assert.Equal(t, "invalid parameter sender: must be not nil", err.Error())
		assert.Empty(t, svc)
	})
	t.Run("it should return an error because the tracer is not valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc, err := todo.NewService("someTopic", sendermock.NewMockSender(ctrl), nil)

		require.Error(t, err)
		var e todo.InvalidServiceParameterError
		require.True(t, errors.As(err, &e))
		assert.Equal(t, "invalid parameter tracer: must be not nil", err.Error())
		assert.Empty(t, svc)
	})
	t.Run("it should return a new service", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc, err := todo.NewService(
			"someTopic",
			sendermock.NewMockSender(ctrl),
			tracingmock.NewMockTracer(ctrl),
		)

		require.NoError(t, err)
		assert.NotEmpty(t, svc)
	})
}

func TestService_Create(t *testing.T) {
	t.Run("it should return an error because the request is not valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		svc, err := todo.NewService(
			"someTopic",
			sendermock.NewMockSender(ctrl),
			tracingmock.NewMockTracer(ctrl),
		)

		require.NoError(t, err)
		assert.NotNil(t, svc)

		resp, err := svc.Create(context.Background(), nil)
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
		assert.Equal(t, "received nil request for creating a todo", st.Message())
		assert.Nil(t, resp)
	})
	t.Run("it should return an error sending a message failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		const topic = "someTopic"

		var (
			req        = &todov1.CreateRequest{}
			mockSender = sendermock.NewMockSender(ctrl)
			mockTracer = tracingmock.NewMockTracer(ctrl)
		)

		svc, err := todo.NewService(
			topic,
			mockSender,
			mockTracer,
		)

		require.NoError(t, err)
		assert.NotNil(t, svc)

		gomock.InOrder(
			mockSender.EXPECT().SendMessage(&sarama.ProducerMessage{
				Topic:   topic,
				Value:   sarama.ByteEncoder([]byte{}),
				Headers: nil,
			}).Return(errors.New("someErr")).Times(1),
		)

		resp, err := svc.Create(context.Background(), req)
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.Internal, st.Code())
		assert.Equal(t, "could not produce message", st.Message())
		assert.Nil(t, resp)
	})
	t.Run("it should return publish a message", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		const topic = "someTopic"

		var (
			req        = &todov1.CreateRequest{}
			mockSender = sendermock.NewMockSender(ctrl)
			mockTracer = tracingmock.NewMockTracer(ctrl)
		)

		svc, err := todo.NewService(
			topic,
			mockSender,
			mockTracer,
		)

		require.NoError(t, err)
		assert.NotNil(t, svc)

		gomock.InOrder(
			mockSender.EXPECT().SendMessage(&sarama.ProducerMessage{
				Topic:   topic,
				Value:   sarama.ByteEncoder([]byte{}),
				Headers: nil,
			}).Return(nil).Times(1),
		)

		resp, err := svc.Create(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}
