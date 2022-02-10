// Copyright 2020 Google LLC
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

// [START cloudrun_helloworld_service]
// [START run_helloworld_service]

// Sample run-helloworld is a minimal Cloud Run service.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.opentelemetry.io/contrib/detectors/gcp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
	log.Print("starting server...")
	initTracer()
	http.Handle("/", otelhttp.NewHandler(&baseHandler{}, "helloworld"))

	// Determine port for HTTP service.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
	}

	// Start HTTP server.
	log.Printf("listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// initTracer creates and registers trace provider instance.
func initTracer() {
	var err error
	// Keep things simple by using a standard out exporter
	exp, err := stdouttrace.New()
	if err != nil {
		log.Fatal(err)
		return
	}
	bsp := sdktrace.NewBatchSpanProcessor(exp)
	resource, err := resource.New(context.Background(),
		// Detect resource information from cloud run detector
		resource.WithDetectors(gcp.NewCloudRun()),
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tp)
}

type baseHandler struct{}

func (*baseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(w, "Hello %s!\n", name)
}

// [END run_helloworld_service]
// [END cloudrun_helloworld_service]
