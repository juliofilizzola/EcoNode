package route

import "github.com/gin-gonic/gin"

func scaleRoutes(engine *gin.Engine, authenticate gin.HandlerFunc) {
	engine.Group("/scale", authenticate)
}
