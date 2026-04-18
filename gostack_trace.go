package gostack

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Trace starts a span named name; caller must end the span.
func Trace(ctx context.Context, name string) (context.Context, trace.Span) {
	return otel.Tracer("gostack").Start(ctx, name)
}
