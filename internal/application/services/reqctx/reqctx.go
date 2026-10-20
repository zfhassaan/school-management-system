// Package reqctx carries request-scoped metadata (correlation id, source).
package reqctx

import (
	"context"

	"school-management/internal/domain/audit"
	"school-management/internal/domain/shared"
)

type key int

const (
	correlationKey key = iota
	sourceKey
)

// WithCorrelationID attaches a correlation id used in logs and audit entries.
func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, correlationKey, id)
}

// CorrelationID returns the request correlation id, generating none if absent.
func CorrelationID(ctx context.Context) string {
	if v, ok := ctx.Value(correlationKey).(string); ok {
		return v
	}
	return ""
}

// WithSource records the transport that initiated the request.
func WithSource(ctx context.Context, s audit.Source) context.Context {
	return context.WithValue(ctx, sourceKey, s)
}

// Source returns the request source, defaulting to system.
func Source(ctx context.Context) audit.Source {
	if v, ok := ctx.Value(sourceKey).(audit.Source); ok {
		return v
	}
	return audit.SourceSystem
}

// New returns a context with a fresh correlation id and the given source.
func New(ctx context.Context, ids shared.IDGenerator, s audit.Source) context.Context {
	return WithSource(WithCorrelationID(ctx, ids.NewID().String()), s)
}

