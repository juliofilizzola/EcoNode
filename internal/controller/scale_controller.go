package controller

import (
	"errors"
	"log"
	"net/http"

	"EcoNode/internal/dto"
	"EcoNode/internal/repository"
	"EcoNode/internal/service"

	"github.com/gin-gonic/gin"
)

type ScaleController struct {
	scaleService service.ScaleService
}

func NewScaleController(scaleService service.ScaleService) *ScaleController {
	return &ScaleController{scaleService: scaleService}
}

func (c *ScaleController) ListNamespaces(ctx *gin.Context) {
	namespaces, err := c.scaleService.ListNamespaces(ctx.Request.Context())
	if err != nil {
		respondInternalError(ctx, "list namespaces", err)
		return
	}

	ctx.JSON(http.StatusOK, dto.NamespaceListResponse{Namespaces: namespaces})
}

func (c *ScaleController) ListDeployments(ctx *gin.Context) {
	var path dto.NamespacePath

	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "invalid path parameters",
			Details: validationDetails(err),
		})
		return
	}

	deployments, err := c.scaleService.ListDeployments(ctx.Request.Context(), path.Namespace)
	if err != nil {
		respondInternalError(ctx, "list deployments", err)
		return
	}

	response := dto.DeploymentListResponse{Deployments: make([]dto.DeploymentResponse, 0, len(deployments))}
	for _, deployment := range deployments {
		response.Deployments = append(response.Deployments, toDeploymentResponse(deployment))
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ScaleController) ScaleDeployment(ctx *gin.Context) {
	var path dto.DeploymentPath

	if err := ctx.ShouldBindUri(&path); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "invalid path parameters",
			Details: validationDetails(err),
		})
		return
	}

	var request dto.ScaleDeploymentRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Message: "invalid request body",
			Details: validationDetails(err),
		})
		return
	}

	deployment, err := c.scaleService.ScaleDeployment(ctx.Request.Context(), path.Namespace, path.Name, *request.Replicas)
	if err != nil {
		if errors.Is(err, repository.ErrDeploymentNotFound) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Message: "deployment not found"})
			return
		}

		respondInternalError(ctx, "scale deployment", err)
		return
	}

	ctx.JSON(http.StatusOK, dto.ScaledDeploymentResponse{
		Namespace: deployment.Namespace,
		Name:      deployment.Name,
		Replicas:  deployment.Replicas,
	})
}

func toDeploymentResponse(deployment repository.Deployment) dto.DeploymentResponse {
	return dto.DeploymentResponse{
		Namespace:         deployment.Namespace,
		Name:              deployment.Name,
		Replicas:          deployment.Replicas,
		ReadyReplicas:     deployment.ReadyReplicas,
		AvailableReplicas: deployment.AvailableReplicas,
	}
}

func respondInternalError(ctx *gin.Context, operation string, err error) {
	log.Printf("%s failed: %v", operation, err)
	ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Message: "internal server error"})
}
