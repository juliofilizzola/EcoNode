package route

import (
	"EcoNode/internal/controller"

	"github.com/gin-gonic/gin"
)

func authRoutes(engine *gin.Engine, authController *controller.AuthController) {
	authGroup := engine.Group("/auth")
	authGroup.POST("/login", authController.Login)
}
