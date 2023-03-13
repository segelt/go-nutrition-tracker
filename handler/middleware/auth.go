package middleware

import (
	"nutritiontracker/internal"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	JWT_MANAGER *internal.JwtManager
}

func NewAuthMiddleware(jwt_secret string) *AuthMiddleware {
	return &AuthMiddleware{
		JWT_MANAGER: internal.New(jwt_secret),
	}
}

func (a *AuthMiddleware) WithAuthentication(jwt_secret string) gin.HandlerFunc {
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

		claims, err := a.JWT_MANAGER.ValidateToken(tokenString)
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
