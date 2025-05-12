package traces

import "go.opentelemetry.io/otel"

var (
	Tracer = otel.Tracer("darkcore")
)
