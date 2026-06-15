package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"

	_ "hexagonalarchitecture/docs"
	databasepostgres "hexagonalarchitecture/internal/adapter/database/postgres"
	httpadapter "hexagonalarchitecture/internal/adapter/handler/http"
	"hexagonalarchitecture/internal/adapter/outboundapi/httpclient"
	"hexagonalarchitecture/internal/adapter/outboundapi/noop"
	"hexagonalarchitecture/internal/adapter/repository/postgres"
	"hexagonalarchitecture/internal/config"
	"hexagonalarchitecture/internal/core/port"
	"hexagonalarchitecture/internal/core/service"
	observabilitylogger "hexagonalarchitecture/internal/observability/logger"
	"hexagonalarchitecture/internal/observability/tracer"
)

// @title Hexagonal Architecture CRUD API
// @version 1.0
// @description A CRUD REST API built with Gin, PostgreSQL, Docker, and Hexagonal Architecture.
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.Load()

	logger, err := observabilitylogger.New()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	defer logger.Sync()

	tracerProvider, err := tracer.NewProvider()
	if err != nil {
		logger.Fatal("failed to initialize tracer", zap.Error(err))
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tracer.Shutdown(ctx, tracerProvider); err != nil {
			logger.Error("failed to shutdown tracer", zap.Error(err))
		}
	}()

	metricsRegistry := prometheus.NewRegistry()
	httpadapter.RegisterMetrics(metricsRegistry)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := databasepostgres.NewPool(ctx, cfg.Database.URL())
	if err != nil {
		logger.Fatal("failed to connect database", zap.Error(err))
	}
	defer dbPool.Close()

	userRepo, err := postgres.NewUserRepository(ctx, dbPool)
	if err != nil {
		logger.Fatal("failed to initialize user repository", zap.Error(err))
	}

	outboundClient := newOutboundAPIClient(cfg)
	userService := service.NewUserService(userRepo, outboundClient)

	r := httpadapter.New(userService, logger, otel.Tracer("hexagonalarchitecture-api"), metricsRegistry)
	server := &http.Server{
		Addr:    cfg.ServerAddress(),
		Handler: r,
	}

	logger.Info("server is running", zap.String("address", cfg.ServerAddress()))
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("failed to run server", zap.Error(err))
		}
	}()

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-shutdownCtx.Done()
	stop()

	logger.Info("shutdown signal received")

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("failed to shutdown server gracefully", zap.Error(err))
	}

	logger.Info("server stopped gracefully")
}

func newOutboundAPIClient(cfg config.Config) port.OutboundAPIClient {
	if cfg.OutboundAPI.BaseURL == "" {
		return noop.NewClient()
	}

	client, err := httpclient.New(httpclient.Config{
		BaseURL: cfg.OutboundAPI.BaseURL,
	})
	if err != nil {
		zap.L().Fatal("failed to create outbound API client", zap.Error(err))
	}

	return client
}
