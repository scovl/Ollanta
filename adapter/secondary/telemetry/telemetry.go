package telemetry

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// Counter is a monotonically increasing counter.
type Counter struct {
	val  atomic.Int64
	help string
	name string
}

// Inc increments the counter by 1.
func (c *Counter) Inc() { c.val.Add(1) }

// Add increments the counter by n.
func (c *Counter) Add(n int64) { c.val.Add(n) }

// Gauge is a value that can go up or down.
type Gauge struct {
	val  atomic.Int64
	help string
	name string
}

// Set sets the gauge to v.
func (g *Gauge) Set(v int64) { g.val.Store(v) }

// Inc increments the gauge by 1.
func (g *Gauge) Inc() { g.val.Add(1) }

// Dec decrements the gauge by 1.
func (g *Gauge) Dec() { g.val.Add(-1) }

// Histogram tracks a distribution of observations via fixed buckets.
type Histogram struct {
	mu      sync.Mutex
	buckets []float64
	counts  []int64
	sum     float64
	count   int64
	help    string
	name    string
}

// Observe records a single observation.
func (h *Histogram) Observe(v float64) {
	h.mu.Lock()
	h.sum += v
	h.count++
	for i, bucket := range h.buckets {
		if v <= bucket {
			h.counts[i]++
		}
	}
	h.mu.Unlock()
}

// ObserveDuration records a duration as seconds.
func (h *Histogram) ObserveDuration(d time.Duration) {
	h.Observe(d.Seconds())
}

// Registry holds all application metrics.
type Registry struct {
	mu         sync.RWMutex
	counters   map[string]*Counter
	gauges     map[string]*Gauge
	histograms map[string]*Histogram
}

var defaultHistogramBuckets = []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}

const metricHelpFormat = "# HELP %s %s\n"

// NewRegistry creates an empty metrics registry.
func NewRegistry() *Registry {
	return &Registry{
		counters:   map[string]*Counter{},
		gauges:     map[string]*Gauge{},
		histograms: map[string]*Histogram{},
	}
}

// Counter registers and returns a counter by name.
func (reg *Registry) Counter(name, help string) *Counter {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	if counter, ok := reg.counters[name]; ok {
		return counter
	}
	counter := &Counter{name: name, help: help}
	reg.counters[name] = counter
	return counter
}

// Gauge registers and returns a gauge by name.
func (reg *Registry) Gauge(name, help string) *Gauge {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	if gauge, ok := reg.gauges[name]; ok {
		return gauge
	}
	gauge := &Gauge{name: name, help: help}
	reg.gauges[name] = gauge
	return gauge
}

// Histogram registers and returns a histogram by name.
func (reg *Registry) Histogram(name, help string) *Histogram {
	reg.mu.Lock()
	defer reg.mu.Unlock()
	if histogram, ok := reg.histograms[name]; ok {
		return histogram
	}
	counts := make([]int64, len(defaultHistogramBuckets))
	histogram := &Histogram{name: name, help: help, buckets: defaultHistogramBuckets, counts: counts}
	reg.histograms[name] = histogram
	return histogram
}

// Handler returns an http.HandlerFunc that serves /metrics in Prometheus text format.
func (reg *Registry) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		reg.mu.RLock()
		defer reg.mu.RUnlock()

		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		var sb strings.Builder

		for _, counter := range reg.counters {
			fmt.Fprintf(&sb, metricHelpFormat, counter.name, counter.help)
			fmt.Fprintf(&sb, "# TYPE %s counter\n", counter.name)
			fmt.Fprintf(&sb, "%s %d\n", counter.name, counter.val.Load())
		}
		for _, gauge := range reg.gauges {
			fmt.Fprintf(&sb, metricHelpFormat, gauge.name, gauge.help)
			fmt.Fprintf(&sb, "# TYPE %s gauge\n", gauge.name)
			fmt.Fprintf(&sb, "%s %d\n", gauge.name, gauge.val.Load())
		}
		for _, histogram := range reg.histograms {
			histogram.mu.Lock()
			fmt.Fprintf(&sb, metricHelpFormat, histogram.name, histogram.help)
			fmt.Fprintf(&sb, "# TYPE %s histogram\n", histogram.name)
			for i, bucket := range histogram.buckets {
				fmt.Fprintf(&sb, "%s_bucket{le=\"%g\"} %d\n", histogram.name, bucket, histogram.counts[i])
			}
			fmt.Fprintf(&sb, "%s_bucket{le=\"+Inf\"} %d\n", histogram.name, histogram.count)
			fmt.Fprintf(&sb, "%s_sum %g\n", histogram.name, histogram.sum)
			fmt.Fprintf(&sb, "%s_count %d\n", histogram.name, histogram.count)
			histogram.mu.Unlock()
		}

		_, _ = w.Write([]byte(sb.String()))
	}
}

// Metrics holds named application metrics.
type Metrics struct {
	HTTPRequestsTotal   *Counter
	HTTPRequestDuration *Histogram
	ScansTotal          *Counter
	IngestDuration      *Histogram
	IngestQueueDepth    *Gauge
	IndexQueueDepth     *Gauge
	IndexJobsProcessed  *Counter
	IndexJobRetries     *Counter
	WebhookQueueDepth   *Gauge
	WebhookDeliveries   *Counter
}

// NewMetrics registers all application metrics in reg.
func NewMetrics(reg *Registry) *Metrics {
	return &Metrics{
		HTTPRequestsTotal:   reg.Counter("ollanta_http_requests_total", "Total number of HTTP requests handled"),
		HTTPRequestDuration: reg.Histogram("ollanta_http_request_duration_seconds", "Duration of HTTP requests in seconds"),
		ScansTotal:          reg.Counter("ollanta_scans_total", "Total number of scans ingested"),
		IngestDuration:      reg.Histogram("ollanta_ingest_duration_seconds", "Duration of scan ingest pipeline in seconds"),
		IngestQueueDepth:    reg.Gauge("ollanta_ingest_queue_depth", "Current depth of the ingest queue"),
		IndexQueueDepth:     reg.Gauge("ollanta_index_queue_depth", "Current depth of the durable index queue"),
		IndexJobsProcessed:  reg.Counter("ollanta_index_jobs_total", "Total number of index jobs processed successfully"),
		IndexJobRetries:     reg.Counter("ollanta_index_job_retries_total", "Total number of index job retries scheduled"),
		WebhookQueueDepth:   reg.Gauge("ollanta_webhook_queue_depth", "Current depth of the durable webhook queue"),
		WebhookDeliveries:   reg.Counter("ollanta_webhook_deliveries_total", "Total webhook deliveries attempted"),
	}
}

// ObserveHTTPRequest records a handled HTTP request.
func (m *Metrics) ObserveHTTPRequest(d time.Duration) {
	if m == nil {
		return
	}
	m.HTTPRequestsTotal.Inc()
	m.HTTPRequestDuration.ObserveDuration(d)
}

// SetupLogger creates a structured logger for runtime services.
func SetupLogger(level string, attrs ...any) *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: parseLogLevel(level)})
	return slog.New(handler).With(attrs...)
}

func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// WithTraceAttrs appends trace correlation fields when they exist in ctx.
func WithTraceAttrs(ctx context.Context, attrs ...any) []any {
	traceID := TraceID(ctx)
	spanID := SpanID(ctx)
	if traceID == "" && spanID == "" {
		return attrs
	}
	out := make([]any, 0, len(attrs)+4)
	out = append(out, attrs...)
	if traceID != "" {
		out = append(out, "trace_id", traceID)
	}
	if spanID != "" {
		out = append(out, "span_id", spanID)
	}
	return out
}

type contextKey string

const traceIDKey contextKey = "trace_id"

// TraceID returns the trace ID stored in ctx, or empty string.
func TraceID(ctx context.Context) string {
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.IsValid() {
		return spanContext.TraceID().String()
	}
	value, _ := ctx.Value(traceIDKey).(string)
	return value
}

// TraceIDMiddleware injects X-Trace-Id into every request.
func TraceIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Trace-Id")
		if currentTraceID := TraceID(r.Context()); currentTraceID != "" {
			traceID = currentTraceID
		}
		if traceID == "" {
			traceID = newUUID()
		}
		ctx := context.WithValue(r.Context(), traceIDKey, traceID)
		w.Header().Set("X-Trace-Id", traceID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func newUUID() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		timeNow().UnixNano()&0xFFFFFFFF,
		timeNow().UnixNano()>>16&0xFFFF,
		(timeNow().UnixNano()>>32&0x0FFF)|0x4000,
		(timeNow().UnixNano()>>48&0x3FFF)|0x8000,
		timeNow().UnixNano()&0xFFFFFFFFFFFF,
	)
}

var timeNow = time.Now
