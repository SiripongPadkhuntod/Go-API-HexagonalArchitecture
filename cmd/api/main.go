package main

import (
	"context"
	"log"
	"time"

	_ "hexagonalarchitecture/docs"
	httpadapter "hexagonalarchitecture/internal/adapter/handler/http"
	"hexagonalarchitecture/internal/adapter/outboundapi/httpclient"
	"hexagonalarchitecture/internal/adapter/outboundapi/noop"
	"hexagonalarchitecture/internal/adapter/repository/postgres"
	"hexagonalarchitecture/internal/config"
	"hexagonalarchitecture/internal/core/port"
	"hexagonalarchitecture/internal/core/service"
)

// @title Hexagonal Architecture CRUD API
// @version 1.0
// @description A CRUD REST API built with Gin, PostgreSQL, Docker, and Hexagonal Architecture.
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userRepo, err := postgres.NewUserRepository(ctx, cfg.Database.URL())
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer userRepo.Close()

	outboundClient := newOutboundAPIClient(cfg)
	userService := service.NewUserService(userRepo, outboundClient)

	r := httpadapter.New(userService)

	log.Printf("server is running on %s", cfg.ServerAddress())
	if err := r.Run(cfg.ServerAddress()); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func newOutboundAPIClient(cfg config.Config) port.OutboundAPIClient {
	if cfg.OutboundAPI.BaseURL == "" {
		return noop.NewClient()
	}

	client, err := httpclient.New(httpclient.Config{
		BaseURL: cfg.OutboundAPI.BaseURL,
	})
	if err != nil {
		log.Fatalf("failed to create outbound API client: %v", err)
	}

	return client
}
