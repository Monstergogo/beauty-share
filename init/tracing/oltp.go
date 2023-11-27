package tracing

import (
	"context"
	"fmt"
	"github.com/Monstergogo/beauty-share/util"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var (
	Shutdowns []func(context.Context) error
	Tracer    trace.Tracer
	Meter     otelmetric.Meter
)

// InitProvider 初始化tracing、meter provider
func InitProvider() {
	res, err := newResource(util.TracingServiceName, util.ServiceVersion)
	if err != nil {
		panic(err)
	}
	initTracerProvider(res)
	initMeterProvider(res)
}

// initTracerProvider Initializes an OTLP exporter, and configures the corresponding trace providers.
func initTracerProvider(res *resource.Resource) {
	ctx := context.Background()

	// replace `localhost` with the endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, util.OTLPCollectorGrpcAddr,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection to collector: %w", err))
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		panic(fmt.Errorf("failed to create trace exporter: %w", err))
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	Shutdowns = append(Shutdowns, tracerProvider.Shutdown)
	Tracer = otel.Tracer("share")
}

func initMeterProvider(res *resource.Resource) {
	exporter, err := prometheus.New()
	if err != nil {
		panic(err)
	}

	meterProvider := metric.NewMeterProvider(metric.WithResource(res),
		metric.WithReader(exporter))
	otel.SetMeterProvider(meterProvider)
	Meter = otel.Meter("share")
	Shutdowns = append(Shutdowns, meterProvider.Shutdown)
}

func newResource(serviceName, serviceVersion string) (*resource.Resource, error) {
	return resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
}
