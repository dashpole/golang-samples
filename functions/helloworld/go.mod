module github.com/GoogleCloudPlatform/golang-samples/functions/helloworld

go 1.16

require (
	cloud.google.com/go/functions v1.0.0
	github.com/GoogleCloudPlatform/opentelemetry-operations-go v0.31.0
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.7.0
	go.opentelemetry.io/contrib/detectors/gcp v1.7.1-0.20220601202618-b7910afac6ff
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.32.1-0.20220601202618-b7910afac6ff // indirect
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/sdk v1.7.0
)
