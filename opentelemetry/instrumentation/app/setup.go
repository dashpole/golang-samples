// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/propagators/autoprop"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/oauth"
)

// setupOpenTelemetry sets up the OpenTelemetry SDK and exporters for metrics and
// traces. If it does not return an error, call shutdown for proper cleanup.
// [START opentelemetry_instrumentation_setup_opentelemetry]
func setupOpenTelemetry(ctx context.Context) (shutdown func(context.Context) error, err error) {
	res, err := resource.New(ctx,
		resource.WithDetectors(gcp.NewDetector()),
		resource.WithFromEnv(),
	)
	if errors.Is(err, resource.ErrPartialResource) || errors.Is(err, resource.ErrSchemaURLConflict) {
		log.Println(err)
	} else if err != nil {
		return nil, err
	}
	var shutdownFuncs []func(context.Context) error

	// shutdown combines shutdown functions from multiple OpenTelemetry
	// components into a single function.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// Configure Context Propagation to use the default W3C traceparent format
	otel.SetTextMapPropagator(autoprop.NewTextMapPropagator())

	// Get Google Application Default credentials
	creds, err := oauth.NewApplicationDefault(ctx)
	if err != nil {
		err = errors.Join(err, shutdown(ctx))
		return
	}
	// Provide credentials to the OTLP trace grpc exporter.
	texporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithDialOption(grpc.WithPerRPCCredentials(creds)))
	if err != nil {
		panic(err)
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(texporter),
		trace.WithResource(res),
	)
	shutdownFuncs = append(shutdownFuncs, tp.Shutdown)
	otel.SetTracerProvider(tp)

	// Configure Metric Export to send metrics as OTLP
	mreader, err := autoexport.NewMetricReader(ctx)
	if err != nil {
		err = errors.Join(err, shutdown(ctx))
		return
	}
	mp := metric.NewMeterProvider(
		metric.WithReader(mreader),
		metric.WithResource(res),
	)
	shutdownFuncs = append(shutdownFuncs, mp.Shutdown)
	otel.SetMeterProvider(mp)

	return shutdown, nil
}

// [END opentelemetry_instrumentation_setup_opentelemetry]

// setupLogging configures logs to write JSON logs to stdout, and add span
// context attributes.
// [START opentelemetry_instrumentation_setup_logging]
func setupLogging() {
	// Use json as our base logging format.
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{ReplaceAttr: replacer})
	// Add span context attributes when Context is passed to logging calls.
	instrumentedHandler := handlerWithSpanContext(jsonHandler)
	// Set this handler as the global slog handler.
	slog.SetDefault(slog.New(instrumentedHandler))
}

// [END opentelemetry_instrumentation_setup_logging]
