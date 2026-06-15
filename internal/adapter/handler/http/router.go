package http

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"hexagonalarchitecture/internal/core/service"
)

func New(userService service.UserService, logger *zap.Logger, tracer trace.Tracer, metricsRegistry *prometheus.Registry) *gin.Engine {
	r := gin.New()
	r.Use(
		RecoveryMiddleware(logger),
		RequestContextMiddleware(logger, tracer),
		MetricsMiddleware(),
	)

	r.GET("/health", Health)
	r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(metricsRegistry, promhttp.HandlerOpts{})))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userHandler := NewUserHandler(userService)

	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.POST("", userHandler.Create)
			users.GET("", userHandler.FindAll)
			users.GET("/:id", userHandler.FindByID)
			users.PUT("/:id", userHandler.Update)
			users.DELETE("/:id", userHandler.Delete)
		}
	}

	return r
}
