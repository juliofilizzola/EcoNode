package service

import (
	"context"

	"EcoNode/internal/integrations/kubernetes"
)

type ScaleService interface {
	ListNamespaces(ctx context.Context) ([]string, error)
	ListDeployments(ctx context.Context, namespace string) ([]kubernetes.Deployment, error)
	ScaleDeployment(ctx context.Context, namespace string, name string, replicas int32) (kubernetes.Deployment, error)
}

type scaleService struct {
	deployments kubernetes.DeploymentGateway
}

func NewScaleService(deployments kubernetes.DeploymentGateway) ScaleService {
	return &scaleService{deployments: deployments}
}

func (s *scaleService) ListNamespaces(ctx context.Context) ([]string, error) {
	return s.deployments.ListNamespaces(ctx)
}

func (s *scaleService) ListDeployments(ctx context.Context, namespace string) ([]kubernetes.Deployment, error) {
	return s.deployments.ListDeployments(ctx, namespace)
}

func (s *scaleService) ScaleDeployment(ctx context.Context, namespace string, name string, replicas int32) (kubernetes.Deployment, error) {
	return s.deployments.ScaleDeployment(ctx, namespace, name, replicas)
}
