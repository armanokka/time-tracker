package server

import (
	"context"
	"errors"
	"github.com/armanokka/test_task_Effective_mobile/config"
	_ "github.com/armanokka/test_task_Effective_mobile/docs"
	"github.com/armanokka/test_task_Effective_mobile/pkg/logger"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const (
	ctxTimeout = 5
)

// Server struct
type Server struct {
	router *gin.Engine
	cfg    *config.Config
	db     *sqlx.DB
	rdb    *redis.Client
	logger logger.Logger
}

// NewServer New Server constructor
func NewServer(cfg *config.Config, db *sqlx.DB, redisClient *redis.Client, logger logger.Logger) *Server {
	return &Server{router: gin.Default(), cfg: cfg, db: db, rdb: redisClient, logger: logger}
}

func (s Server) Run(ctx context.Context) error {
	s.router.Use(gin.Recovery(), requestid.New())
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	server := &http.Server{
		Addr:    "0.0.0.0:" + strconv.Itoa(s.cfg.Server.Port),
		Handler: s.router,
	}

	go func() {
		s.logger.Infof("Server is listening on PORT: %d", s.cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatalf("Error starting Server: ", err)
		}
	}()

	go func() {
		s.logger.Infof("Starting Debug Server on PORT: %d", s.cfg.Server.PprofPort)
		if err := http.ListenAndServe(":"+strconv.Itoa(s.cfg.Server.PprofPort), http.DefaultServeMux); err != nil {
			s.logger.Errorf("Error PPROF ListenAndServe: %s", err)
		}
	}()

	s.MapHandlers(s.router.Group("/api"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return server.Shutdown(ctx)
}
