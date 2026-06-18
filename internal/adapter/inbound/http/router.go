package http

import (
	stdhttp "net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/otel/trace"

	"hexagonalarchitecture/internal/adapter/inbound/http/handler"
	"hexagonalarchitecture/internal/adapter/inbound/http/response"
	"hexagonalarchitecture/internal/core/port"
)

func New(userService port.UserService, storage port.StoragePort, logger port.Logger, tracer trace.Tracer, metricsRegistry *prometheus.Registry) stdhttp.Handler {
	r := gin.New()
	r.Use(
		RecoveryMiddleware(logger),
		RequestContextMiddleware(logger, tracer),
		MetricsMiddleware(),
	)

	r.GET("/health", handler.Health)
	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(metricsRegistry, promhttp.HandlerOpts{})))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userHandler := handler.NewUserHandler(userService)
	var fileHandler *handler.FileHandler
	if storage != nil {
		fileHandler = handler.NewFileHandler(storage)
	}

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", func(c *gin.Context) { response.Bind(c, stdhttp.StatusCreated, userHandler.Create) })
			users.GET("", func(c *gin.Context) { response.Bind(c, stdhttp.StatusOK, userHandler.FindAll) })
			users.GET("/:id", func(c *gin.Context) { response.Bind(c, stdhttp.StatusOK, userHandler.FindByID) })
			users.PUT("/:id", func(c *gin.Context) { response.Bind(c, stdhttp.StatusOK, userHandler.Update) })
			users.DELETE("/:id", func(c *gin.Context) { response.Bind(c, stdhttp.StatusOK, userHandler.Delete) })
		}
		
		if fileHandler != nil {
			files := v1.Group("/files")
			{
				files.POST("/upload", fileHandler.Upload)
			}
		}
	}

	return r
}
