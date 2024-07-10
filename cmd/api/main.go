package main

import (
	"context"
	"fmt"
	"github.com/armanokka/test_task_Effective_mobile/config"
	"github.com/armanokka/test_task_Effective_mobile/internal/server"
	"github.com/armanokka/test_task_Effective_mobile/pkg/db/postgres"
	"github.com/armanokka/test_task_Effective_mobile/pkg/db/redis"
	"github.com/armanokka/test_task_Effective_mobile/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title           Test task for Effective Mobile
// @version         1.0
// @description     REST API for Effective Mobile
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    https://t.me/armanokka
// @contact.email  armangokka@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:80
// @BasePath  /api/

// @tag.name 		auth
// @tag.description Auth section

// @tag.name 		projects
// @tag.description Projects section

// @tag.name 		tasks
// @tag.description Tasks section

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	// Creating context
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGQUIT)
	defer signal.Stop(c)
	go func() {
		<-c
		cancel()
	}()

	// Creating logger
	log := logger.NewApiLogger(cfg)
	log.InitLogger()

	// Creating tracer
	shutdownTracerProvider, err := initProvider(ctx, "otel-collector:4317")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err = shutdownTracerProvider(ctx); err != nil {
			log.Error("shutdown", zap.Error(err))
		}
	}()
	otel.SetTracerProvider(otel.GetTracerProvider()) // setting global tracer provider

	// Connecting to Postgres
	db, err := postgres.NewPsqlDB(ctx, &cfg.Postgres)
	if err != nil {
		panic(err)
	}

	// Connecting to Redis
	rdb, err := redis.NewRedisClient(ctx, cfg)
	if err != nil {
		panic(err)
	}

	if err = server.NewServer(cfg, db, rdb, log).Run(ctx); err != nil {
		panic(err)
	}
}

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider(ctx context.Context, otelGrpcReceiverDSN string) (func(context.Context) error, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String("api"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	conn, err := grpc.NewClient(otelGrpcReceiverDSN,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := tracesdk.NewBatchSpanProcessor(traceExporter)

	tracerProvider := tracesdk.NewTracerProvider(
		//tracesdk.WithBatcher(exporter),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithResource(res),
		tracesdk.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider.Shutdown, nil
}
