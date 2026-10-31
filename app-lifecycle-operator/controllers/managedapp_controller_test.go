package controllers

import (
	platformv1alpha1 "github.com/REPLACE-ME/app-lifecycle-operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestDesiredService(t *testing.T) {
	app := &platformv1alpha1.ManagedApp{ObjectMeta: metav1.ObjectMeta{Name: "demo", Namespace: "apps"}}
	svc := desiredService(app, 8080)
	if svc.Name != "demo" || svc.Namespace != "apps" || svc.Spec.Ports[0].Port != 8080 || svc.Spec.Type != corev1.ServiceTypeClusterIP {
		t.Fatalf("unexpected service: %#v", svc)
	}
}

func TestDesiredDeployment(t *testing.T) {
	app := &platformv1alpha1.ManagedApp{ObjectMeta: metav1.ObjectMeta{Name: "demo"}, Spec: platformv1alpha1.ManagedAppSpec{Image: "nginx:1.27"}}
	dep := desiredDeployment(app, 3, 8080)
	if *dep.Spec.Replicas != 3 || dep.Spec.Template.Spec.Containers[0].Image != "nginx:1.27" {
		t.Fatalf("unexpected deployment: %#v", dep.Spec)
	}
}
