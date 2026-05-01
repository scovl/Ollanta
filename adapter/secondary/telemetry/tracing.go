package telemetry

import (
	"context"
	"net"
	"net/http"
	"os"
	"net/url"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

// SetupTracing configures the global OpenTelemetry tracer provider.
// When no OTLP endpoint is configured it becomes a cheap no-op.
func SetupTracing(ctx context.Context, serviceName string, attrs ...attribute.KeyValue) (func(context.Context) error, error) {
	endpoint := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT"))
	if endpoint == "" {
		endpoint = strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	}
	if endpoint == "" {
		otel.SetTextMapPropagator(propagation.TraceContext{})
		return func(context.Context) error { return nil }, nil
	}
	if !traceEndpointReachable(endpoint) {
		otel.SetTextMapPropagator(propagation.TraceContext{})
		return func(context.Context) error { return nil }, nil
	}

	options := []otlptracehttp.Option{}
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		options = append(options, otlptracehttp.WithEndpointURL(endpoint))
	} else {
		options = append(options, otlptracehttp.WithEndpoint(endpoint))
	}
	if strings.HasPrefix(endpoint, "http://") {
		options = append(options, otlptracehttp.WithInsecure())
	}

	exporter, err := otlptracehttp.New(ctx, options...)
	if err != nil {
		return nil, err
	}

	resourceAttrs := append([]attribute.KeyValue{attribute.String("service.name", serviceName)}, attrs...)
	res, err := resource.New(ctx, resource.WithFromEnv(), resource.WithProcess(), resource.WithAttributes(resourceAttrs...))
	if err != nil {
		return nil, err
	}

	provider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(res),
	)
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return provider.Shutdown, nil
}

// WrapHTTPHandler instruments an HTTP handler with OpenTelemetry spans.
func WrapHTTPHandler(name string, next http.Handler) http.Handler {
	return otelhttp.NewHandler(next, name)
}

// SpanID returns the current span identifier, or empty string.
func SpanID(ctx context.Context) string {
	spanContext := trace.SpanContextFromContext(ctx)
	if !spanContext.IsValid() {
		return ""
	}
	return spanContext.SpanID().String()
}

// StartSpan creates a named span using the global tracer provider.
func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.Tracer("github.com/scovl/ollanta").Start(ctx, name)
}

func traceEndpointReachable(endpoint string) bool {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return false
	}
	host := parsed.Host
	if host == "" {
		host = parsed.Path
	}
	if !strings.Contains(host, ":") {
		switch parsed.Scheme {
		case "https":
			host += ":443"
		default:
			host += ":80"
		}
	}
	conn, err := net.DialTimeout("tcp", host, 500*time.Millisecond)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}