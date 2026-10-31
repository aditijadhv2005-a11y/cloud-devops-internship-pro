package controllers

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	platformv1alpha1 "github.com/REPLACE-ME/app-lifecycle-operator/api/v1alpha1"
)

// ManagedAppReconciler reconciles ManagedApp resources.
// +kubebuilder:rbac:groups=platform.example.com,resources=managedapps,verbs=get;list;watch
// +kubebuilder:rbac:groups=platform.example.com,resources=managedapps/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=platform.example.com,resources=managedapps/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

type ManagedAppReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *ManagedAppReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var app platformv1alpha1.ManagedApp
	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Owner references provide garbage-collection cleanup; no finalizer is needed.
	if !app.DeletionTimestamp.IsZero() {
		return ctrl.Result{}, nil
	}

	replicas := int32(2)
	if app.Spec.Replicas != nil {
		replicas = *app.Spec.Replicas
	}
	port := app.Spec.Port
	if port == 0 {
		port = 8080
	}
	if app.Spec.Image == "" {
		meta.SetStatusCondition(&app.Status.Conditions, metav1.Condition{Type: "Degraded", Status: metav1.ConditionTrue, Reason: "InvalidSpec", Message: "spec.image must not be empty", ObservedGeneration: app.Generation})
		if err := r.Status().Update(ctx, &app); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	desiredDeployment := desiredDeployment(&app, replicas, port)
	desiredDeploymentSpec := desiredDeployment.Spec
	desiredDeploymentLabels := map[string]string{"app.kubernetes.io/managed-by": "app-lifecycle-operator"}
	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, desiredDeployment, func() error {
		if err := controllerutil.SetControllerReference(&app, desiredDeployment, r.Scheme); err != nil {
			return err
		}
		desiredDeployment.Labels = desiredDeploymentLabels
		desiredDeployment.Spec = desiredDeploymentSpec
		return nil
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("reconcile deployment: %w", err)
	}

	desiredSvc := desiredService(&app, port)
	desiredServiceSpec := desiredSvc.Spec
	desiredServiceLabels := map[string]string{"app.kubernetes.io/managed-by": "app-lifecycle-operator"}
	_, err = controllerutil.CreateOrUpdate(ctx, r.Client, desiredSvc, func() error {
		if err := controllerutil.SetControllerReference(&app, desiredSvc, r.Scheme); err != nil {
			return err
		}
		desiredSvc.Labels = desiredServiceLabels
		desiredSvc.Spec = desiredServiceSpec
		return nil
	})
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("reconcile service: %w", err)
	}

	var current appsv1.Deployment
	if err := r.Get(ctx, types.NamespacedName{Namespace: app.Namespace, Name: app.Name}, &current); err != nil {
		return ctrl.Result{}, err
	}
	app.Status.ObservedGeneration = app.Generation
	app.Status.ReadyReplicas = current.Status.ReadyReplicas
	conditionStatus, reason, message := metav1.ConditionFalse, "Progressing", "Waiting for desired replicas to become ready"
	if current.Status.ReadyReplicas >= replicas {
		conditionStatus, reason, message = metav1.ConditionTrue, "Available", "Desired replicas are ready"
	}
	meta.SetStatusCondition(&app.Status.Conditions, metav1.Condition{Type: "Ready", Status: conditionStatus, Reason: reason, Message: message, ObservedGeneration: app.Generation})
	meta.SetStatusCondition(&app.Status.Conditions, metav1.Condition{Type: "Degraded", Status: metav1.ConditionFalse, Reason: "AsExpected", Message: "No reconciliation error reported", ObservedGeneration: app.Generation})
	if err := r.Status().Update(ctx, &app); err != nil {
		return ctrl.Result{}, err
	}
	logger.Info("reconciled ManagedApp", "name", app.Name, "namespace", app.Namespace, "readyReplicas", current.Status.ReadyReplicas)
	return ctrl.Result{}, nil
}

func desiredDeployment(app *platformv1alpha1.ManagedApp, replicas, port int32) *appsv1.Deployment {
	labels := map[string]string{"app.kubernetes.io/name": app.Name, "app.kubernetes.io/managed-by": "app-lifecycle-operator"}
	return &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: app.Name, Namespace: app.Namespace}, Spec: appsv1.DeploymentSpec{
		Replicas: &replicas, Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"app.kubernetes.io/name": app.Name}},
		Template: corev1.PodTemplateSpec{ObjectMeta: metav1.ObjectMeta{Labels: labels}, Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: "app", Image: app.Spec.Image, Ports: []corev1.ContainerPort{{ContainerPort: port}}, Resources: app.Spec.Resources, ReadinessProbe: &corev1.Probe{ProbeHandler: corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{Path: "/healthz", Port: intstr.FromInt32(port)}}, InitialDelaySeconds: 3, PeriodSeconds: 5}, LivenessProbe: &corev1.Probe{ProbeHandler: corev1.ProbeHandler{HTTPGet: &corev1.HTTPGetAction{Path: "/healthz", Port: intstr.FromInt32(port)}}, InitialDelaySeconds: 10, PeriodSeconds: 10}}}}},
	}}
}

func desiredService(app *platformv1alpha1.ManagedApp, port int32) *corev1.Service {
	return &corev1.Service{ObjectMeta: metav1.ObjectMeta{Name: app.Name, Namespace: app.Namespace, Labels: map[string]string{"app.kubernetes.io/managed-by": "app-lifecycle-operator"}}, Spec: corev1.ServiceSpec{Selector: map[string]string{"app.kubernetes.io/name": app.Name}, Ports: []corev1.ServicePort{{Name: "http", Port: port, TargetPort: intstr.FromInt32(port)}}, Type: corev1.ServiceTypeClusterIP}}
}
