package kubernetes

import (
	"context"
	"errors"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeclient "k8s.io/client-go/kubernetes"
)

var ErrDeploymentNotFound = errors.New("deployment not found")

type Deployment struct {
	Namespace         string
	Name              string
	Replicas          int32
	ReadyReplicas     int32
	AvailableReplicas int32
}

type DeploymentGateway interface {
	ListNamespaces(ctx context.Context) ([]string, error)
	ListDeployments(ctx context.Context, namespace string) ([]Deployment, error)
	ScaleDeployment(ctx context.Context, namespace string, name string, replicas int32) (Deployment, error)
}

type clientDeploymentGateway struct {
	clientset kubeclient.Interface
}

func NewDeploymentGateway(clientset kubeclient.Interface) DeploymentGateway {
	return &clientDeploymentGateway{clientset: clientset}
}

func (g *clientDeploymentGateway) ListNamespaces(ctx context.Context) ([]string, error) {
	namespaceList, err := g.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list namespaces: %w", err)
	}

	names := make([]string, 0, len(namespaceList.Items))
	for _, namespace := range namespaceList.Items {
		names = append(names, namespace.Name)
	}

	return names, nil
}

func (g *clientDeploymentGateway) ListDeployments(ctx context.Context, namespace string) ([]Deployment, error) {
	deploymentList, err := g.clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("list deployments in namespace %s: %w", namespace, err)
	}

	deployments := make([]Deployment, 0, len(deploymentList.Items))
	for _, item := range deploymentList.Items {
		deployments = append(deployments, toDeployment(item))
	}

	return deployments, nil
}

func (g *clientDeploymentGateway) ScaleDeployment(ctx context.Context, namespace string, name string, replicas int32) (Deployment, error) {
	deployments := g.clientset.AppsV1().Deployments(namespace)

	scale, err := deployments.GetScale(ctx, name, metav1.GetOptions{})
	if err != nil {
		return Deployment{}, mapDeploymentError(err, namespace, name, "read scale of deployment")
	}

	scale.Spec.Replicas = replicas

	updatedScale, err := deployments.UpdateScale(ctx, name, scale, metav1.UpdateOptions{})
	if err != nil {
		return Deployment{}, mapDeploymentError(err, namespace, name, "update scale of deployment")
	}

	return Deployment{
		Namespace: namespace,
		Name:      name,
		Replicas:  updatedScale.Spec.Replicas,
	}, nil
}

func toDeployment(item appsv1.Deployment) Deployment {
	var replicas int32
	if item.Spec.Replicas != nil {
		replicas = *item.Spec.Replicas
	}

	return Deployment{
		Namespace:         item.Namespace,
		Name:              item.Name,
		Replicas:          replicas,
		ReadyReplicas:     item.Status.ReadyReplicas,
		AvailableReplicas: item.Status.AvailableReplicas,
	}
}

func mapDeploymentError(err error, namespace string, name string, action string) error {
	if apierrors.IsNotFound(err) {
		return fmt.Errorf("%s %s/%s: %w", action, namespace, name, ErrDeploymentNotFound)
	}

	return fmt.Errorf("%s %s/%s: %w", action, namespace, name, err)
}
