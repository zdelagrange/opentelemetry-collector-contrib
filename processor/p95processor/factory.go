package p95processor

import (
	"go.opentelemetry.io/collector/processor"
)

// NewFactory returns a new factory for the Metrics transform processor.
func NewFactory() processor.Factory {
	return processor.NewFactory(
		metadata.Type,
		createDefaultConfig,
		processor.WithMetrics(createMetricsProcessor, metadata.MetricsStability))
}
