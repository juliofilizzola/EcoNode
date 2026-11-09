package route

import (
	"EcoNode/internal/controller"

	"github.com/gin-gonic/gin"
)

func SetupEngine(authController *controller.AuthController, authenticate gin.HandlerFunc) *gin.Engine {
	engine := gin.Default()
	health(engine)
	authRoutes(engine, authController)
	scaleRoutes(engine, authenticate)
	return engine
}
