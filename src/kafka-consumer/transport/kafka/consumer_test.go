package kafka_test

import (
	"errors"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/golang/mock/gomock"
	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andream16/go-opentracing-example/src/kafka-consumer/transport/kafka"
	"github.com/andream16/go-opentracing-example/src/shared/todo"
	todocreatormock "github.com/andream16/go-opentracing-example/src/test/mock/kafka-consumer/todo/repository"
	opentracingmock "github.com/andream16/go-opentracing-example/src/test/mock/opentracing"
	tracingmock "github.com/andream16/go-opentracing-example/src/test/mock/tracing"
)

func TestNewConsumer(t *testing.T) {
	t.Run("it should return an error because the creator is invalid", func(t *testing.T) {
		consumer, err := kafka.NewConsumer(nil, nil)
		require.Error(t, err)
		assert.Equal(t, "repo must be not nil", err.Error())
		assert.Empty(t, consumer)
	})
	t.Run("it should return an error because the tracer is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		consumer, err := kafka.NewConsumer(todocreatormock.NewMockCreator(ctrl), nil)
		require.Error(t, err)
		assert.Equal(t, "tracer must be not nil", err.Error())
		assert.Empty(t, consumer)
	})
	t.Run("it should return a new consumer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		consumer, err := kafka.NewConsumer(todocreatormock.NewMockCreator(ctrl), tracingmock.NewMockTracer(ctrl))
		require.NoError(t, err)
		assert.NotEmpty(t, consumer)
	})
}

func TestConsumer_ReceivedMessage(t *testing.T) {
	t.Run("it should return an error because creating a todo failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		const spanName = "todo_consumer"

		var (
			mockCreator     = todocreatormock.NewMockCreator(ctrl)
			mockTracer      = tracingmock.NewMockTracer(ctrl)
			mockSpanContext = opentracingmock.NewMockSpanContext(ctrl)
			mockSpan        = opentracingmock.NewMockSpan(ctrl)
			kafkaHeaders    = []*sarama.RecordHeader{
				{
					Key:   []byte(`key1`),
					Value: []byte(`value1`),
				},
			}
			tracingHeaders = map[string]string{
				"key1": "value1",
			}
		)

		consumer, err := kafka.NewConsumer(mockCreator, mockTracer)
		require.NoError(t, err)
		assert.NotEmpty(t, consumer)

		gomock.InOrder(
			mockTracer.
				EXPECT().
				Extract(opentracing.TextMap, opentracing.TextMapCarrier(tracingHeaders)).
				Return(mockSpanContext, nil).
				Times(1),
			mockTracer.
				EXPECT().
				StartSpan(spanName, gomock.Any()).
				Return(mockSpan).
				Times(1),
			mockSpan.
				EXPECT().
				Tracer().
				Times(1),
			mockCreator.
				EXPECT().
				Create(gomock.Any(), &todo.Todo{}).
				Return(errors.New("someErr")).
				Times(1),
			mockSpan.EXPECT().Finish().Times(1),
		)

		require.Error(t, consumer.ReceivedMessage(&sarama.ConsumerMessage{
			Headers: kafkaHeaders,
		}))
	})
	t.Run("it should create a new todo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		const spanName = "todo_consumer"

		var (
			mockCreator     = todocreatormock.NewMockCreator(ctrl)
			mockTracer      = tracingmock.NewMockTracer(ctrl)
			mockSpanContext = opentracingmock.NewMockSpanContext(ctrl)
			mockSpan        = opentracingmock.NewMockSpan(ctrl)
			kafkaHeaders    = []*sarama.RecordHeader{
				{
					Key:   []byte(`key1`),
					Value: []byte(`value1`),
				},
			}
			tracingHeaders = map[string]string{
				"key1": "value1",
			}
		)

		consumer, err := kafka.NewConsumer(mockCreator, mockTracer)
		require.NoError(t, err)
		assert.NotEmpty(t, consumer)

		gomock.InOrder(
			mockTracer.
				EXPECT().
				Extract(opentracing.TextMap, opentracing.TextMapCarrier(tracingHeaders)).
				Return(mockSpanContext, nil).
				Times(1),
			mockTracer.
				EXPECT().
				StartSpan(spanName, gomock.Any()).
				Return(mockSpan).
				Times(1),
			mockSpan.
				EXPECT().
				Tracer().
				Times(1),
			mockCreator.
				EXPECT().
				Create(gomock.Any(), &todo.Todo{}).
				Return(nil).
				Times(1),
			mockSpan.EXPECT().Finish().Times(1),
		)

		require.NoError(t, consumer.ReceivedMessage(&sarama.ConsumerMessage{
			Headers: kafkaHeaders,
		}))
	})
}
