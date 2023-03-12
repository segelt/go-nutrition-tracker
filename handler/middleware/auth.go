package middleware

import (
	"nutritiontracker/internal"

	"github.com/gin-gonic/gin"
)

func WithAuthentication() gin.HandlerFunc {
	return func(context *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := context.GetHeader("Authorization")
		// tokenString := authHeader[len(BEARER_SCHEMA):]
		if authHeader == "" {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}
		tokenString := authHeader[len(BEARER_SCHEMA):]
		if tokenString == "" {
			context.JSON(401, gin.H{"error": "request does not contain an access token"})
			context.Abort()
			return
		}

		claims, err := internal.ValidateToken(tokenString)
		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		context.Set("userid", claims.Id)
		context.Set("username", claims.UserName)
		context.Next()
	}
}
