package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"hexagonalarchitecture/internal/adapter/http/handler"
	"hexagonalarchitecture/internal/core/service"
)

func New(userService service.UserService) *gin.Engine {
	r := gin.Default()

	r.GET("/health", handler.Health)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userHandler := handler.NewUserHandler(userService)

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
