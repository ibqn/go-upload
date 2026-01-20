package routes

import (
	"go-upload/controllers"

	"github.com/gin-gonic/gin"
)

func FileRoutes(router *gin.RouterGroup) {
	fileRouter := router.Group("/file")

	fileRouter.GET("/:id", controllers.HandleGetFile)
}
