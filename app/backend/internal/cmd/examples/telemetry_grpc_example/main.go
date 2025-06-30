// SPDX-FileCopyrightText: 2025 BROKE DA EAR LLC <https://brokedaear.com>
//
// SPDX-License-Identifier: Apache-2.0

//nolint:all // This is an example file.
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc/health/grpc_health_v1"

	"go.brokedaear.com/internal/common/telemetry"
	"go.brokedaear.com/internal/common/utils/loggers"
	"go.brokedaear.com/internal/core/domain"
	"go.brokedaear.com/internal/core/server"
)

func main() {
	serviceName := flag.String("serviceName", "brokedabackend", "name of this service")
	serviceVersion := flag.String("serviceVersion", "0.0.1", "semver of this service")
	serviceID := flag.String("serviceId", "local-1", "ID of this service")
	serviceEnv := flag.String("env", "development", "service runtime environment")

	flag.Parse()

	env, err := domain.EnvFromString(*serviceEnv)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	telCfg, err := telemetry.NewConfig(
		*serviceName,
		*serviceVersion,
		*serviceID,
		telemetry.NewExporterConfig(
			telemetry.WithType(telemetry.ExporterTypeGRPC),
			telemetry.WithEndpoint("127.0.0.1:4317"),
		),
	)
	if err != nil {
		fmt.Printf("failed to setup telemetry config: %v", err)
		os.Exit(1) //nolint:gocritic // Its an exit.
	}

	tel, err := telemetry.New(ctx, telCfg)
	if err != nil {
		fmt.Printf("Failed to setup telemetry: %s", err)
		os.Exit(1)
	}

	logger, err := loggers.NewZap(&loggers.ZapConfig{
		Env:          env,
		Telemetry:    tel,
		CustomZapper: nil,
	})
	if err != nil {
		fmt.Printf("failed to setup logger %v", err)
		os.Exit(1)
	}

	srv, err := server.NewGRPCServer(
		logger,
		server.WithVersion(*serviceVersion),
		server.WithTelemetry(tel),
	)
	if err != nil {
		fmt.Printf("failed to setup grpc server %v", err)
		os.Exit(1)
	}

	logger.Info("Started server",
		"name", *serviceName,
		"environment", *serviceEnv,
		"version", *serviceVersion,
		"id", *serviceID,
	)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		setupHealthMonitoring(ctx, srv, tel, logger)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Starting gRPC server", "port", 8080)
		err := srv.ListenAndServe(ctx)
		if err != nil {
			logger.Error("Server error", "error", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		generateSampleTelemetry(ctx, tel, logger)
	}()

	defer teardown()

	sig := <-c

	logger.Warn(fmt.Sprintf("Received %s, shutting down...", sig))

	logger.Warn("Waiting for jobs to complete...")

	cancel()

	wg.Wait()
}

func teardown(closers ...io.Closer) {
	for _, closer := range closers {
		closer.Close()
	}
	fmt.Println("Bye!")
}

func setupHealthMonitoring(
	ctx context.Context,
	srv server.GRPCServer,
	tel telemetry.Telemetry,
	logger server.Logger,
) {
	dbHealthGauge, err := tel.Gauge(
		telemetry.Metric{
			Name:        "database_connection_health",
			Unit:        "{status}",
			Description: "Database connection health status (1=healthy, 0=unhealthy)",
		},
	)
	if err != nil {
		logger.Error("Failed to create database health gauge", "error", err)
		return
	}

	cacheHealthGauge, err := tel.Gauge(
		telemetry.Metric{
			Name:        "cache_connection_health",
			Unit:        "{status}",
			Description: "Cache connection health status (1=healthy, 0=unhealthy)",
		},
	)
	if err != nil {
		logger.Error("Failed to create cache health gauge", "error", err)
		return
	}

	const healthMonitorInterval = 30 * time.Second

	for {
		select {
		case <-time.After(healthMonitorInterval):
			jCtx, span := tel.TraceStart(ctx, "health_monitoring_cycle")

			dbHealthy := simulateDatabaseHealthCheck(ctx, tel)
			dbStatus := grpc_health_v1.HealthCheckResponse_NOT_SERVING
			if dbHealthy {
				dbStatus = grpc_health_v1.HealthCheckResponse_SERVING
			}

			// Record database health metric
			var dbHealthValue int64
			if dbHealthy {
				dbHealthValue = 1
			}
			dbHealthGauge.Record(jCtx, dbHealthValue)

			// Simulate cache health check
			cacheHealthy := simulateCacheHealthCheck(jCtx, tel)
			cacheStatus := grpc_health_v1.HealthCheckResponse_NOT_SERVING
			if cacheHealthy {
				cacheStatus = grpc_health_v1.HealthCheckResponse_SERVING
			}

			// Record cache health metric
			cacheHealthValue := int64(0)
			if cacheHealthy {
				cacheHealthValue = 1
			}
			cacheHealthGauge.Record(jCtx, cacheHealthValue)

			// Update service health statuses
			srv.SetHealthStatus("database-service", dbStatus)
			srv.SetHealthStatus("cache-service", cacheStatus)

			// Overall health depends on all dependencies
			overallStatus := grpc_health_v1.HealthCheckResponse_SERVING
			if !dbHealthy || !cacheHealthy {
				overallStatus = grpc_health_v1.HealthCheckResponse_NOT_SERVING
			}
			srv.SetHealthStatus("", overallStatus)

			logger.Info(
				"Health check completed",
				"database_healthy", dbHealthy,
				"cache_healthy", cacheHealthy,
				"overall_status", overallStatus.String(),
			)

			span.End()
		case <-ctx.Done():
			return
		}
	}
}

func simulateDatabaseHealthCheck(
	ctx context.Context,
	tel telemetry.Telemetry,
) bool {
	_, span := tel.TraceStart(ctx, "database_health_check")
	defer span.End()

	// Simulate some work
	time.Sleep(10 * time.Millisecond)

	// 90% chance of being healthy
	healthy := time.Now().Unix()%10 != 0
	return healthy
}

func simulateCacheHealthCheck(
	ctx context.Context,
	tel telemetry.Telemetry,
) bool {
	_, span := tel.TraceStart(ctx, "cache_health_check")
	defer span.End()

	// Simulate some work
	time.Sleep(5 * time.Millisecond)

	// 95% chance of being healthy
	healthy := time.Now().Unix()%20 != 0
	return healthy
}

func generateSampleTelemetry(
	ctx context.Context,
	tel telemetry.Telemetry,
	logger server.Logger,
) {
	// Create sample metrics
	requestCounter, err := tel.UpDownCounter(
		telemetry.Metric{
			Name:        "sample_requests_total",
			Unit:        "{count}",
			Description: "Total number of sample requests",
		},
	)
	if err != nil {
		logger.Error("Failed to create request counter", "error", err)
		return
	}

	responseTimeHistogram, err := tel.Histogram(
		telemetry.Metric{
			Name:        "sample_response_time",
			Unit:        "ms",
			Description: "Sample response time in milliseconds",
		},
	)
	if err != nil {
		logger.Error("Failed to create response time histogram", "error", err)
		return
	}

	workDuration := time.Duration(50+time.Now().UnixNano()%100) * time.Millisecond

	for {
		select {
		case <-time.After(5 * time.Second):
			jCtx, span := tel.TraceStart(ctx, "sample_operation")

			// Simulate some work

			time.Sleep(workDuration)
			requestCounter.Add(jCtx, 1)
			responseTimeHistogram.Record(jCtx, workDuration.Milliseconds())

			span.End()

			logger.Debug(
				"Generated sample telemetry data",
				"duration_ms", workDuration.Milliseconds(),
			)
		case <-ctx.Done():
			return
		}
	}
}
