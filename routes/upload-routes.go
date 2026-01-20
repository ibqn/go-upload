package routes

import (
	"go-upload/controllers"
	"go-upload/middleware"

	"github.com/gin-gonic/gin"
)

func UploadRoutes(router *gin.RouterGroup) {
	uploadRouter := router.Group("/upload")

	uploadRouter.POST("/", middleware.AuthRequired, controllers.CreateUpload)
	uploadRouter.GET("/", middleware.AuthRequired, controllers.ListUploads)
	uploadRouter.GET("/:id", middleware.AuthRequired, controllers.GetUpload)
	uploadRouter.DELETE("/:id", middleware.AuthRequired, controllers.DeleteUpload)
}
