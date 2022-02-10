module github.com/GoogleCloudPlatform/golang-samples/run/helloworld

go 1.13

require (
	go.opentelemetry.io/contrib/detectors/gcp v1.3.0
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.28.0
	go.opentelemetry.io/otel v1.3.0
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.3.0
	go.opentelemetry.io/otel/sdk v1.3.0
)
