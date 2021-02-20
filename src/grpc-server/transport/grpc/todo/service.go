package todo

import (
	"context"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	todov1 "github.com/andream16/go-opentracing-example/contracts/build/go/go_opentracing_example/grpc_server/todo/v1"
	"github.com/andream16/go-opentracing-example/src/shared/kafka"
	"github.com/andream16/go-opentracing-example/src/shared/tracing"
)

// Service implements the grpc service.
type Service struct {
	kafkaTopic    string
	kafkaProducer kafka.SyncProducer
	tracer        tracing.Tracer
}

// NewService returns a new Service.
func NewService(kafkaTopic string, kafkaProducer kafka.SyncProducer, tracer tracing.Tracer) Service {
	return Service{
		kafkaTopic:    kafkaTopic,
		kafkaProducer: kafkaProducer,
		tracer:        tracer,
	}
}

// Creates a new todo.
func (svc Service) Create(ctx context.Context, req *todov1.CreateRequest) (*todov1.CreateResponse, error) {
	if req == nil {
		log.Println("received nil request for creating a todo")
		return nil, status.Error(codes.InvalidArgument, "received nil request for creating a todo")
	}

	var saramaHeaders []sarama.RecordHeader

	headers := make(map[string]string)
	if span := opentracing.SpanFromContext(ctx); span != nil {
		_ = svc.tracer.Inject(
			span.Context(),
			opentracing.TextMap,
			opentracing.TextMapCarrier(headers),
		)
	}

	for k, v := range headers {
		saramaHeaders = append(saramaHeaders, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}

	b, err := proto.Marshal(req)
	if err != nil {
		log.Println(fmt.Sprintf("could not marshal request: %v", err))
		return nil, status.Error(codes.Internal, "could not marshal request")
	}

	if err := svc.kafkaProducer.SendMessage(&sarama.ProducerMessage{
		Topic:   svc.kafkaTopic,
		Value:   sarama.ByteEncoder(b),
		Headers: saramaHeaders,
	}); err != nil {
		log.Println(fmt.Sprintf("could not produce message: %v", err))
		return nil, status.Error(codes.Internal, "could not produce message")
	}

	return &todov1.CreateResponse{}, nil
}
