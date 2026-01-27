package router

import (
	"github.com/gin-gonic/gin"
	"go-upload/internal/handler"
	"go-upload/internal/middleware"
)

// SetupRoutes configures all application routes
func SetupRoutes(
	app *gin.Engine,
	authHandler *handler.AuthHandler,
	uploadHandler *handler.UploadHandler,
	fileHandler *handler.FileHandler,
	imageHandler *handler.ImageHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	// Public routes (no authentication required)
	fileRouter := app.Group("/")
	{
		fileRouter.GET("/file/:id", fileHandler.GetFile)
		fileRouter.GET("/image/:id", imageHandler.GetImage)
	}

	// API routes
	apiRouter := app.Group("/api")
	{
		// Auth routes
		authRouter := apiRouter.Group("/auth")
		{
			authRouter.POST("/signin", authHandler.SignIn)
			authRouter.POST("/signup", authHandler.SignUp)
			authRouter.POST("/signout", authMiddleware.RequireAuth(), authHandler.SignOut)
			authRouter.GET("/user", authMiddleware.RequireAuth(), authHandler.GetUser)
		}

		// Upload routes (all require authentication)
		uploadRouter := apiRouter.Group("/upload")
		uploadRouter.Use(authMiddleware.RequireAuth())
		{
			uploadRouter.POST("/", uploadHandler.CreateUpload)
			uploadRouter.GET("/", uploadHandler.ListUploads)
			uploadRouter.GET("/:id", uploadHandler.GetUpload)
			uploadRouter.DELETE("/:id", uploadHandler.DeleteUpload)
		}
	}
}
