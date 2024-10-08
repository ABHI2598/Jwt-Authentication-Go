package routes

import (
	"github.com/abhinavpandey/jwtProject/controllers"
	"github.com/abhinavpandey/jwtProject/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.AuthMiddleWare())
	incomingRoutes.GET("/users", controllers.GetUsers())
	incomingRoutes.GET("/user/:id", controllers.GetUser())
}
