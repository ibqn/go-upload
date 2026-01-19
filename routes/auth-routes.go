package routes

import (
	"go-upload/controllers"
	"go-upload/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup) {
	authRouter := router.Group("/auth")

	authRouter.POST("/signin", controllers.HandleSignIn)
	authRouter.POST("/signup", controllers.HandleSignUp)
	authRouter.POST("/signout", middleware.AuthRequired, controllers.HandleSignOut)
	authRouter.GET("/user", middleware.AuthRequired, controllers.GetUser)
}
