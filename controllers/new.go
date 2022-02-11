package controllers

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	appdemov1 "github.com/fumanne/appdemo-operator/api/v1"
)

func lbls(k, v string) map[string]string {
	m := make(map[string]string)
	m[k] = v
	return m
}

func NewDeploy(nginxDeploy *appdemov1.NginxApp) *appsv1.Deployment {

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nginxDeploy.Name,
			Namespace: nginxDeploy.Namespace,
			Labels:    lbls("release", nginxDeploy.Name),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(nginxDeploy, schema.GroupVersionKind{
					Group:   "appdemo",
					Version: "v1",
					Kind:    "NginxApp",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: nginxDeploy.Spec.Num,
			Selector: &metav1.LabelSelector{
				MatchLabels: lbls("release", nginxDeploy.Name),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      nginxDeploy.Name,
					Namespace: nginxDeploy.Namespace,
					Labels:    lbls("release", nginxDeploy.Name),
				},
				Spec: corev1.PodSpec{
					Containers: NewContainers(nginxDeploy),
				},
			},
		},
		Status: appsv1.DeploymentStatus{},
	}

}

func NewContainers(nginxContainer *appdemov1.NginxApp) []corev1.Container {
	return []corev1.Container{
		{
			Name:  nginxContainer.Name,
			Image: nginxContainer.Spec.Image,
			Env:   nginxContainer.Spec.Env,
		},
	}

}

func NewService(nginxService *appdemov1.NginxApp) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nginxService.Name,
			Namespace: nginxService.Namespace,
			Labels:    lbls("release", nginxService.Name),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(nginxService, schema.GroupVersionKind{
					Group:   "appdemo",
					Version: "v1",
					Kind:    "NginxApp",
				}),
			},
		},
		Spec: corev1.ServiceSpec{
			Ports:    nginxService.Spec.Ports,
			Selector: lbls("release", nginxService.Name),
			Type:     corev1.ServiceTypeClusterIP,
		},
		Status: corev1.ServiceStatus{},
	}

}

func NewCM(nginxCM *appdemov1.NginxApp) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      nginxCM.Name,
			Namespace: nginxCM.Namespace,
			Labels:    lbls("release", nginxCM.Name),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(nginxCM, schema.GroupVersionKind{
					Group:   "appdemo",
					Version: "v1",
					Kind:    "NginxApp",
				}),
			},
		},
		Data: nginxCM.Spec.Data,
	}

}
