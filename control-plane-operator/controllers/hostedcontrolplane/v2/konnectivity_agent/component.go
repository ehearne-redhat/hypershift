package konnectivity

import (
	kasv2 "github.com/openshift/hypershift/control-plane-operator/controllers/hostedcontrolplane/v2/kas"
	component "github.com/openshift/hypershift/support/controlplane-component"

	rbacv1 "k8s.io/api/rbac/v1"
)

const (
	ComponentName = "konnectivity-agent"
)

var _ component.ComponentOptions = &konnectivityAgent{}

type konnectivityAgent struct {
}

// IsRequestServing implements controlplanecomponent.ComponentOptions.
func (r *konnectivityAgent) IsRequestServing() bool {
	return false
}

// MultiZoneSpread implements controlplanecomponent.ComponentOptions.
func (r *konnectivityAgent) MultiZoneSpread() bool {
	return true
}

// NeedsManagementKASAccess implements controlplanecomponent.ComponentOptions.
func (r *konnectivityAgent) NeedsManagementKASAccess() bool {
	return false
}

func NewComponent() component.ControlPlaneComponent {
	return component.NewDeploymentComponent(ComponentName, &konnectivityAgent{}).
		WithAdaptFunction(adaptDeployment).
		WithDependencies(kasv2.ComponentName).
		WithManifestAdapter(
			"clusterrolebinding.yaml",
			component.WithAdaptFunction(adaptClusterRoleBinding),
		).
		Build()
}

func adaptClusterRoleBinding(cpContext component.WorkloadContext, crb *rbacv1.ClusterRoleBinding) error {
	// Set the namespace for ServiceAccount subjects
	for i := range crb.Subjects {
		if crb.Subjects[i].Kind == "ServiceAccount" {
			crb.Subjects[i].Namespace = cpContext.HCP.Namespace
		}
	}
	return nil
}
