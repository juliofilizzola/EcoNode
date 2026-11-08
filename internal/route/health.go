package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func health(engine *gin.Engine) {
	engine.GET("/health", healthStatus)
}

func healthStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"server":  "running",
		"version": "1.0.0",
	})
}
