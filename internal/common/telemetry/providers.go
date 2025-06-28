// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

package telemetry

import (
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// NewLoggerProvider creates a new logger provider with the OTLP gRPC exporter.
func NewLoggerProvider(res *resource.Resource, exporter log.Exporter) *log.LoggerProvider {
	processor := log.NewBatchProcessor(exporter)
	lp := log.NewLoggerProvider(
		log.WithProcessor(processor),
		log.WithResource(res),
	)

	return lp
}

// NewMeterProvider creates a new meter provider with the OTLP gRPC exporter.
func NewMeterProvider(
	res *resource.Resource,
	exporter metric.Exporter,
) *metric.MeterProvider {
	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(mp)

	return mp
}

// NewTracerProvider creates a new tracer provider with the OTLP gRPC exporter.
func NewTracerProvider(res *resource.Resource, exporter trace.SpanExporter) *trace.TracerProvider {
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return tp
}

// NewResource creates a new OTEL resource.
func NewResource(name, version, id string) *resource.Resource {
	hostName, _ := os.Hostname()

	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(name),
		semconv.ServiceVersion(version),
		semconv.ServiceInstanceIDKey.String(id),
		semconv.HostName(hostName),
	)
}
