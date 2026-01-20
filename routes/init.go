package routes

import "github.com/gin-gonic/gin"

func SetupRoutes(app *gin.Engine) {
	fileRouter := app.Group("/")

	FileRoutes(fileRouter)

	router := app.Group("/api")

	AuthRoutes(router)
	UploadRoutes(router)
}
