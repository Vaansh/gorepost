package api

import (
	"net/http"

	"github.com/Vaansh/gore/internal/util"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	authToken := c.GetHeader("Authorization")
	expectedToken := util.Getenv("API_AUTH_TOKEN", true)

	if authToken != expectedToken {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.Next()
}
