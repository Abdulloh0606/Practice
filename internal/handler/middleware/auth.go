package middleware

import (
	"context"
	"minitrello/pkg/auth"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProjectRoleChecker interface {
	GetUserProjectRole(ctx context.Context, projectID int, userID int) (string, error)
}

func JWTauth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(401, gin.H{"error": "Invalid Authorization header format"})
			c.Abort()
			return
		}
		token, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(401, gin.H{"error:": "invalid token"})
			c.Abort()
			return
		}
		c.Set("user_id", token.UserID)
		c.Set("role", token.Role)
		c.Next()

	}
}

func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, exists := c.Get("role")
		if !exists {
			c.JSON(401, gin.H{"error": "no role in token"})
			c.Abort()
			return
		}

		userRole := r.(string)

		if userRole != role {
			c.JSON(403, gin.H{"error": "insufficient rights"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireProjectRole(service ProjectRoleChecker, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		userID := id.(int)

		projectIDStr := c.Param("project_id")
		projectID, err := strconv.Atoi(projectIDStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid project id"})
			c.Abort()
			return
		}

		role, err := service.GetUserProjectRole(c.Request.Context(), projectID, userID)
		if err != nil {
			c.JSON(403, gin.H{"error": "not a member"})
			c.Abort()
			return
		}

		if role != requiredRole {
			c.JSON(403, gin.H{"error": "insufficient rights"})
			c.Abort()
			return
		}

		c.Next()
	}
}
