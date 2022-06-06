module helloworld

go 1.15

require (
	github.com/GoogleCloudPlatform/opentelemetry-operations-go v0.32.1-0.20220606142555-e818bf8a84c9
	github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace v1.8.0
	go.opentelemetry.io/contrib/detectors/gcp v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.32.1-0.20220601202618-b7910afac6ff
	go.opentelemetry.io/otel v1.7.0
	go.opentelemetry.io/otel/sdk v1.7.0
)

// Use PR: https://github.com/open-telemetry/opentelemetry-go-contrib/pull/2341
replace go.opentelemetry.io/contrib/detectors/gcp => github.com/dashpole/opentelemetry-go-contrib/detectors/gcp v0.12.1-0.20220607122304-0658c36d4ce3
