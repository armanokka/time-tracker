package main

import (
	"context"
	"fmt"
	"github.com/armanokka/time_tracker/config"
	"github.com/armanokka/time_tracker/internal/server"
	"github.com/armanokka/time_tracker/pkg/db/postgres"
	"github.com/armanokka/time_tracker/pkg/db/redis"
	"github.com/armanokka/time_tracker/pkg/logger"
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

/*
1. Сделать возможность обновить на пустое значение
2. Поправить структуру проекта по советам Вячеслава
3. Сделать проверку CSRF-токенов
*/

// @title           Time tracker REST API
// @version         1.0
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

	// Creating tracer provider
	tracerProvider, err := initProvider(ctx, cfg.Tracer.OtelGRPCReceiverDSN)
	if err != nil {
		panic(err)
	}
	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{}) // set global propagator to tracecontext (the default is no-op).

	defer func() {
		if err = tracerProvider.Shutdown(ctx); err != nil {
			log.Error("shutdown", zap.Error(err))
		}
	}()

	// Connecting to Postgres
	db, err := postgres.NewPsqlDB(ctx, &postgres.Config{
		User:     cfg.Postgres.User,
		DB:       cfg.Postgres.DB,
		Password: cfg.Postgres.Password,
		Driver:   cfg.Postgres.Driver,
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
	})
	if err != nil {
		panic(err)
	}

	// Connecting to Redis
	rdb, err := redis.NewRedisClient(ctx, &redis.Config{
		Addr:         cfg.Redis.Addr,
		DB:           cfg.Redis.DB,
		Password:     cfg.Redis.Password,
		MinIdleConns: cfg.Redis.MinIdleConns,
		PoolTimeout:  cfg.Redis.PoolTimeout,
		PoolSize:     cfg.Redis.PoolSize,
	})
	if err != nil {
		panic(err)
	}

	if err = server.NewServer(cfg, db, rdb, log).Run(ctx); err != nil {
		panic(err)
	}
}

// Initializes an OTLP exporter, and configures the corresponding trace and
// metric providers.
func initProvider(ctx context.Context, otelGrpcReceiverDSN string) (*tracesdk.TracerProvider, error) {
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

	return tracesdk.NewTracerProvider(
		//tracesdk.WithBatcher(exporter),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		tracesdk.WithResource(res),
		tracesdk.WithSpanProcessor(bsp),
	), nil
}
