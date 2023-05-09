/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	vault "github.com/hashicorp/vault/api"
	"k8s.io/apimachinery/pkg/runtime"
	vaultv1 "qts.cloud/vault-controller/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PolicyReconciler reconciles a Policy object
type PolicyReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	VaultClient *vault.Client
}

//+kubebuilder:rbac:groups=vault.qts.cloud,resources=policies,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=vault.qts.cloud,resources=policies/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=vault.qts.cloud,resources=policies/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Policy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *PolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	policy := &vaultv1.Policy{}
	err := r.Get(ctx, req.NamespacedName, policy)
	if err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Convert policy rules to HCL format
	policyHCL := ""
	for _, rule := range policy.Spec.Rules {
		policyHCL += "path \"" + rule.Path + "\" {\n"
		policyHCL += "  capabilities = " + vaultCapabilitiesToString(rule.Capabilities) + "\n"
		policyHCL += "}\n"
	}

	// Write the policy to Vault
	err = r.VaultClient.Sys().PutPolicy(policy.ObjectMeta.Name, policyHCL)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Update the policy status
	policy.Status.Applied = true
	err = r.Status().Update(ctx, policy)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&vaultv1.Policy{}).
		Complete(r)
}

func vaultCapabilitiesToString(capabilities []string) string {
	if len(capabilities) == 0 {
		return "[]"
	}

	result := "["
	for i, capability := range capabilities {
		result += "\"" + capability + "\""
		if i < len(capabilities)-1 {
			result += ", "
		}
	}
	result += "]"

	return result
}
