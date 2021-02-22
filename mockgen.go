package gen

// Internal
//go:generate mockgen -package tracingmock -destination src/test/mock/tracing/tracing_mock.go -source src/shared/tracing/tracing.go Tracer
//go:generate mockgen -package transporthttpmock -destination src/test/mock/transport/http/transporthttp_mock.go -source src/shared/transport/http/doer.go Doer
//go:generate mockgen -package todoclientmock -destination src/test/mock/todoclient/todoclient_mock.go -source contracts/build/go/go_opentracing_example/grpc_server/todo/v1/todo_service_grpc.pb.go TodoServiceClient
//go:generate mockgen -package sendermock -destination src/test/mock/kafka/sender_mock.go -source src/shared/kafka/sender.go Sender
//go:generate mockgen -package todocreatormock -destination src/test/mock/kafka-consumer/todo/repository/repository_mock.go -source src/kafka-consumer/todo/repository/repository.go Creator
//go:generate mockgen -package executormock -destination src/test/mock/database/postgres/executor_mock.go -source src/shared/database/postgres/executor.go Executor

// External
//go:generate mockgen -package opentracingmock -destination src/test/mock/opentracing/opentracing_mock.go -source vendor/github.com/opentracing/opentracing-go/span.go Span,SpanContext
