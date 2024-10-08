package main

import (
	"log"
	"os"

	"github.com/abhinavpandey/jwtProject/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	//set gin.Release mode for production mode you can use .env file
	gin.SetMode(gin.DebugMode)

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Unable to load .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("/api-1", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"Success": "Access Granted for API-1"})
	})

	router.Run(":" + port)

}
