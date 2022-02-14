/*
Copyright 2022.

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
	appdemov1 "github.com/fumanne/appdemo-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// NginxAppReconciler reconciles a NginxApp object
type NginxAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

const NginxAppFinalizer = "appdemo.dailygn.com/finalizer"

//+kubebuilder:rbac:groups=appdemo.dailygn.com,resources=nginxapps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=appdemo.dailygn.com,resources=nginxapps/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=appdemo.dailygn.com,resources=nginxapps/finalizers,verbs=update
//+kubebuilder:rbac:groups=*,resources=*,verbs=*
//+kubebuilder:subresource:NginxAppStatus

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NginxApp object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *NginxAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	appLog := log.FromContext(ctx)
	appLog.Info("Starting Nginx App Logging..........")

	// TODO(user): your logic here
	instance := &appdemov1.NginxApp{}
	err := r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{Requeue: true}, nil
		}
		appLog.Error(err, "11111111--Get NginxApp Error", "NginxApp Info", instance)
		return ctrl.Result{}, err
	}

	appLog.Info("2222222222--Get instance Info", "NginxAPP info", instance)

	//
	//isNginxAppMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	//if isNginxAppMarkedToBeDeleted {
	//	if controllerutil.ContainsFinalizer(instance, NginxAppFinalizer) {
	//		controllerutil.RemoveFinalizer(instance, NginxAppFinalizer)
	//		if err := r.Client.Update(ctx, instance); err != nil {
	//			return ctrl.Result{}, err
	//		}
	//	}
	//	return ctrl.Result{}, nil
	//}
	//
	//if !controllerutil.ContainsFinalizer(instance, NginxAppFinalizer) {
	//	controllerutil.AddFinalizer(instance, NginxAppFinalizer)
	//	if err := r.Client.Update(ctx, instance); err != nil {
	//		return ctrl.Result{}, err
	//	}
	//}

	deploy := &appsv1.Deployment{}
	if err := r.Client.Get(ctx, req.NamespacedName, deploy); err != nil && errors.IsNotFound(err) {
		deploy := NewDeploy(instance)
		if err := r.Client.Create(ctx, deploy); err != nil {
			appLog.Error(err, "333333333333--Create Deployment Error", "Deployment", deploy)
			return ctrl.Result{}, err
		}
		if err := controllerutil.SetControllerReference(instance, deploy, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
	}
	appLog.Info("AAAAAAAAAAAA--Show Deployment ", "Deployment", deploy)

	svc := &corev1.Service{}
	if err := r.Client.Get(ctx, req.NamespacedName, svc); err != nil && errors.IsNotFound(err) {
		svc := NewService(instance)
		if err := r.Client.Create(ctx, svc); err != nil {
			appLog.Error(err, "BBBBBBBBBBB--Create Service Error", "Service", svc)
			return ctrl.Result{}, err
		}
		if err := controllerutil.SetControllerReference(instance, svc, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
	}

	appLog.Info("CCCCCCCCCCCC--Show Service Info", "Service", svc)

	cm := &corev1.ConfigMap{}
	if err := r.Client.Get(ctx, req.NamespacedName, cm); err != nil && errors.IsNotFound(err) {
		cm := NewCM(instance)
		if err := r.Client.Create(ctx, cm); err != nil {
			appLog.Error(err, "DDDDDDDDD--Create CM Error", "CM", cm)
			return ctrl.Result{}, err
		}
		if err := controllerutil.SetControllerReference(instance, cm, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}
	}
	appLog.Info("EEEEEEEEEEE--Show CM Info", "CM", cm)

	// Update resource according NginxApp Spec
	expectedDeploy := NewDeploy(instance)
	if !reflect.DeepEqual(expectedDeploy.Spec, deploy.Spec) {
		deploy.Spec = expectedDeploy.Spec
		if err := r.Client.Update(ctx, deploy); err != nil {
			appLog.Error(err, "UUUUUUUUUU-- Update Deploy Error", "DeploySpec", deploy.Spec)
			return ctrl.Result{}, err
		}
		appLog.Info("Update Deploy Spec", "DeploySpec", deploy.Spec)
	}

	expectedSvc := NewService(instance)
	if !reflect.DeepEqual(expectedSvc.Spec, svc.Spec) {
		svc.Spec = expectedSvc.Spec
		if err := r.Client.Update(ctx, svc); err != nil {
			appLog.Error(err, "UUUUUUUUUU-- Update Svc Error", "Svc Spec", svc.Spec)
			return ctrl.Result{}, err
		}
		appLog.Info("Update Svc Spec", "Svc Spec", svc.Spec)
	}

	expectedCM := NewCM(instance)
	if !reflect.DeepEqual(expectedCM.Data, cm.Data) {
		cm.Data = expectedCM.Data
		if err := r.Client.Update(ctx, cm); err != nil {
			appLog.Error(err, "UUUUUUUUUU-- Update Cm Data Error", "Cm Data", cm.Data)
			return ctrl.Result{}, err
		}
		appLog.Info("Update CM Data", "CM Data", cm.Data)
	}

	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(instance, NginxAppFinalizer) {
			controllerutil.AddFinalizer(instance, NginxAppFinalizer)
			if err := r.Client.Update(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(instance, NginxAppFinalizer) {
			if err := r.deleteExternalResource(ctx, deploy, svc, cm); err != nil {
				appLog.Error(err, "Delete Resource Error")
				return ctrl.Result{}, err
			}
			controllerutil.RemoveFinalizer(instance, NginxAppFinalizer)
			if err := r.Client.Update(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{Requeue: true}, nil
}

func (r *NginxAppReconciler) deleteExternalResource(ctx context.Context, d *appsv1.Deployment, s *corev1.Service, c *corev1.ConfigMap) error {
	if err := r.Client.Delete(ctx, d); err != nil {
		return err
	}

	if err := r.Client.Delete(ctx, s); err != nil {
		return err
	}

	if err := r.Client.Delete(ctx, c); err != nil {
		return err
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NginxAppReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appdemov1.NginxApp{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
