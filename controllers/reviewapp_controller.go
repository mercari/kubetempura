/*
Copyright 2021.

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

package controllers

import (
	"context"

	kubetempurav1 "github.com/mercari/kubetempura/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ReviewAppReconciler reconciles a ReviewApp object
type ReviewAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kubetempura.mercari.com,resources=reviewapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubetempura.mercari.com,resources=reviewapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubetempura.mercari.com,resources=reviewapps/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ReviewAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Do nothing
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ReviewAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubetempurav1.ReviewApp{}).
		Complete(r)
}
