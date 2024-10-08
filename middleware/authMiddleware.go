package middleware

import (
	"net/http"

	"github.com/abhinavpandey/jwtProject/helpers"
	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientToken := ctx.Request.Header.Get("Authorization")

		if clientToken == "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no Authorization Header Provided"})
			ctx.Abort()
			return
		}

		claims, err := helpers.ValidateToken(clientToken)
		if err != "" {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			ctx.Abort()
			return
		}

		ctx.Set("email", claims.Email)
		ctx.Set("first_name", claims.FirstName)
		ctx.Set("last_name", claims.LastName)
		ctx.Set("user_type", claims.UserType)
		ctx.Set("user_id", claims.Uid)

		ctx.Next()

	}
}
