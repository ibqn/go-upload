package main

import (
	"fmt"
	"go-upload/routes"
	"go-upload/utils"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	PORT := os.Getenv("PORT")

	utils.SetupDb()

	app := gin.Default()

	routes.SetupRoutes(app)

	app.Run(fmt.Sprintf(":%s", PORT))
}
