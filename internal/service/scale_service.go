package service

import (
	"context"

	"EcoNode/internal/repository"
)

type ScaleService interface {
	ListNamespaces(ctx context.Context) ([]string, error)
	ListDeployments(ctx context.Context, namespace string) ([]repository.Deployment, error)
	ScaleDeployment(ctx context.Context, namespace string, name string, replicas int32) (repository.Deployment, error)
}

type scaleService struct {
	deployments repository.DeploymentRepository
}

func NewScaleService(deployments repository.DeploymentRepository) ScaleService {
	return &scaleService{deployments: deployments}
}

func (s *scaleService) ListNamespaces(ctx context.Context) ([]string, error) {
	return s.deployments.ListNamespaces(ctx)
}

func (s *scaleService) ListDeployments(ctx context.Context, namespace string) ([]repository.Deployment, error) {
	return s.deployments.ListDeployments(ctx, namespace)
}

func (s *scaleService) ScaleDeployment(ctx context.Context, namespace string, name string, replicas int32) (repository.Deployment, error) {
	return s.deployments.ScaleDeployment(ctx, namespace, name, replicas)
}
