package main

import (
	"context"
	"log"
	"time"

	_ "hexagonalarchitecture/docs"
	"hexagonalarchitecture/internal/adapter/http/router"
	"hexagonalarchitecture/internal/adapter/repository/postgres"
	"hexagonalarchitecture/internal/config"
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

	userService := service.NewUserService(userRepo)

	r := router.New(userService)

	log.Printf("server is running on %s", cfg.ServerAddress())
	if err := r.Run(cfg.ServerAddress()); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
