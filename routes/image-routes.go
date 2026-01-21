package routes

import (
	"go-upload/controllers"

	"github.com/gin-gonic/gin"
)

func ImageRoutes(router *gin.RouterGroup) {
	imageRouter := router.Group("/image")

	imageRouter.GET("/:id", controllers.HandleGetImage)
}
