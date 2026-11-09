package route

import (
	"EcoNode/internal/controller"

	"github.com/gin-gonic/gin"
)

func scaleRoutes(engine *gin.Engine, authenticate gin.HandlerFunc, scaleController *controller.ScaleController) {
	scaleGroup := engine.Group("/scale", authenticate)
	scaleGroup.GET("/namespaces", scaleController.ListNamespaces)
	scaleGroup.GET("/namespaces/:namespace/deployments", scaleController.ListDeployments)
	scaleGroup.PATCH("/namespaces/:namespace/deployments/:name/scale", scaleController.ScaleDeployment)
}
