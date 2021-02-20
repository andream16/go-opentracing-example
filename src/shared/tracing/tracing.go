package tracing

import (
	"io"
	"log"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

// Tracer defines the open tracing interface.
type Tracer interface {
	StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span
	Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error
	Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error)
	Close() error
}

// JaegerTracer wraps an opentracing tracer with a jaeger reporter.
type JaegerTracer struct {
	closer       io.Closer
	globalTracer opentracing.Tracer
}

// NewJaegerTracer returns a traces with defaults.
func NewJaegerTracer(serviceName, reporterHostName, reporterPort string) Tracer {
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: reporterHostName + ":" + reporterPort,
		},
	}

	tracer, closer, err := cfg.NewTracer(
		config.Logger(jaegerlog.StdLogger),
		config.Metrics(metrics.NullFactory),
	)
	if err != nil {
		log.Fatalf("could not initialise tracer: %v", err)
	}

	opentracing.SetGlobalTracer(tracer)

	return JaegerTracer{
		closer:       closer,
		globalTracer: tracer,
	}
}

// StartSpan wraps StartSpan.
func (jt JaegerTracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	return jt.globalTracer.StartSpan(operationName, opts...)
}

// StartSpan wraps Inject.
func (jt JaegerTracer) Inject(sm opentracing.SpanContext, format interface{}, carrier interface{}) error {
	return jt.globalTracer.Inject(sm, format, carrier)
}

// Extract wraps Extract.
func (jt JaegerTracer) Extract(format interface{}, carrier interface{}) (opentracing.SpanContext, error) {
	return jt.globalTracer.Extract(format, carrier)
}

// Close closes the tracer.
func (jt JaegerTracer) Close() error {
	return jt.closer.Close()
}
