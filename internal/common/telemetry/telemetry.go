// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

// Package telemetry exposes a custom OpenTelemetry implementation.
package telemetry

import (
	"context"
	"io"

	"go.brokedaear.com/pkg/errors"
	otelmetric "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Telemetry interface {
	Instruments
	OtelConfig
	io.Closer
}

type Instruments interface {
	Histogram(Metric) (otelmetric.Int64Histogram, error)
	UpDownCounter(Metric) (otelmetric.Int64UpDownCounter, error)
	Gauge(Metric) (otelmetric.Int64Gauge, error)
	TraceStart(context.Context, string) (context.Context, oteltrace.Span)
}

type OtelConfig interface {
	LoggerProvider() *log.LoggerProvider
	ServiceName() string
	ServiceVersion() string
	ServiceID() string
}

// otelTelemetry wraps OpenTelemetry's logger, meter, and tracer with some
// additional configuration for an exporter.
type otelTelemetry struct {
	lp     *log.LoggerProvider
	mp     *metric.MeterProvider
	tp     *trace.TracerProvider
	meter  otelmetric.Meter
	tracer oteltrace.Tracer
	Config
}

// New creates a new instance of Telemetry.
func New(ctx context.Context, config Config) (Telemetry, error) {
	err := config.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "invalid telemetry config")
	}

	rp := NewResource(
		config.ServiceName(),
		config.ServiceVersion(),
		config.ServiceID(),
	)

	le, err := newLoggerExporter(ctx, config.ExporterConfig())
	if err != nil {
		return nil, err
	}

	lp := NewLoggerProvider(rp, le)

	me, err := newMetricExporter(ctx, config.ExporterConfig())
	if err != nil {
		return nil, err
	}

	mp := NewMeterProvider(rp, me)

	meter := mp.Meter(config.ServiceName())

	te, err := newTraceExporter(ctx, config.ExporterConfig())
	if err != nil {
		return nil, err
	}

	tp := NewTracerProvider(rp, te)

	tracer := tp.Tracer(config.ServiceName())

	return &otelTelemetry{
		lp:     lp,
		mp:     mp,
		tp:     tp,
		meter:  meter,
		tracer: tracer,
		Config: config,
	}, nil
}

// Histogram creates a new int64 histogram meter.
func (t *otelTelemetry) Histogram(metric Metric) (
	otelmetric.Int64Histogram,
	error,
) { //nolint:ireturn // interface requires returning concrete type
	histogram, err := t.meter.Int64Histogram(
		metric.Name,
		otelmetric.WithDescription(metric.Description),
		otelmetric.WithUnit(metric.Unit),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create int64 histogram")
	}

	return histogram, nil
}

// UpDownCounter creates a new int64 up down counter meter.
func (t *otelTelemetry) UpDownCounter(metric Metric) (
	otelmetric.Int64UpDownCounter,
	error,
) { //nolint:ireturn // interface requires returning concrete type
	counter, err := t.meter.Int64UpDownCounter(
		metric.Name,
		otelmetric.WithDescription(metric.Description),
		otelmetric.WithUnit(metric.Unit),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create int64 up down counter")
	}

	return counter, nil
}

// Gauge creates a new int64 gauge meter.
func (t *otelTelemetry) Gauge(metric Metric) (otelmetric.Int64Gauge, error) {
	gauge, err := t.meter.Int64Gauge(
		metric.Name,
		otelmetric.WithDescription(metric.Description),
		otelmetric.WithUnit(metric.Unit),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create int64 gauge")
	}

	return gauge, nil
}

// TraceStart starts a new span with the given name. The span must be ended by calling End.
func (t *otelTelemetry) TraceStart(ctx context.Context, name string) (
	context.Context,
	oteltrace.Span,
) { //nolint:ireturn // interface requires returning concrete type
	//nolint:spancheck // span is intentionally returned for caller to manage
	return t.tracer.Start(ctx, name)
}

// LoggerProvider returns the OpenTelemetry logger provider for log integration.
func (t *otelTelemetry) LoggerProvider() *log.LoggerProvider {
	return t.lp
}

// Close shuts down all the otelTelemetry facilities.
func (t *otelTelemetry) Close() error {
	ctx := context.Background()

	err1 := t.lp.Shutdown(ctx)
	err2 := t.mp.Shutdown(ctx)
	err3 := t.tp.Shutdown(ctx)

	return errors.Join(err1, err2, err3)
}
