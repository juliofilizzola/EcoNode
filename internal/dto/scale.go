package dto

type NamespacePath struct {
	Namespace string `uri:"namespace" binding:"required,max=63,hostname_rfc1123"`
}

type DeploymentPath struct {
	Namespace string `uri:"namespace" binding:"required,max=63,hostname_rfc1123"`
	Name      string `uri:"name" binding:"required,max=63,hostname_rfc1123"`
}

type ScaleDeploymentRequest struct {
	Replicas *int32 `json:"replicas" binding:"required,min=0,max=1000"`
}

type NamespaceListResponse struct {
	Namespaces []string `json:"namespaces"`
}

type DeploymentResponse struct {
	Namespace         string `json:"namespace"`
	Name              string `json:"name"`
	Replicas          int32  `json:"replicas"`
	ReadyReplicas     int32  `json:"ready_replicas"`
	AvailableReplicas int32  `json:"available_replicas"`
}

type ScaledDeploymentResponse struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
	Replicas  int32  `json:"replicas"`
}

type DeploymentListResponse struct {
	Deployments []DeploymentResponse `json:"deployments"`
}
