package infrastructure

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtService *JWTService
}

func NewAuthMiddleware(jwtService *JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (a *AuthMiddleware) AuthMiddlewareRole(reqRole []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		token, err := a.jwtService.ValidateToken(authParts[1])
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT"})
			c.Abort()
			return
		}

		claims, err := a.jwtService.ExtractClaims(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT claims"})
			c.Abort()
			return
		}

		role := claims["role"].(string)
		roleFound := false
		for _, r := range reqRole {
			if r == role {
				roleFound = true
				break
			}
		}

		if !roleFound {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		c.Set("userId", claims["userId"])
		c.Set("role", role)

		c.Next()
	}
}

func (a *AuthMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		authParts := strings.Split(authHeader, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header"})
			c.Abort()
			return
		}

		token, err := a.jwtService.ValidateToken(authParts[1])
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT"})
			c.Abort()
			return
		}

		claims, err := a.jwtService.ExtractClaims(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid JWT claims"})
			c.Abort()
			return
		}

		role := claims["role"].(string)

		c.Set("userId", claims["userId"])
		c.Set("role", role)

		c.Next()
	}
}
