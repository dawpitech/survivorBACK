package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func EnsureIncomingFromLocalhost(c *gin.Context) {
	if c.ClientIP() != "127.0.0.1" && c.ClientIP() != "::1" {
		if os.Getenv("GIN_MODE") == "debug" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied. Only accessible through local network"})
		} else {
			c.AbortWithStatus(http.StatusNotFound)
		}
		return
	}
	c.Next()
}
