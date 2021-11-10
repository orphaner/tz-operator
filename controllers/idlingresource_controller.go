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
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"tz/utils/array"
	"tz/utils/k8s"

	kidlev1 "tz/api/v1"
)

// IdlingResourceReconciler reconciles a IdlingResource object
type IdlingResourceReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kidle.kidle.dev,resources=idlingresources,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kidle.kidle.dev,resources=idlingresources/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kidle.kidle.dev,resources=idlingresources/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the IdlingResource object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *IdlingResourceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("idlingresource", req.NamespacedName)
	logger.V(1).Info("Starting reconcile loop")
	defer logger.V(1).Info("Finish reconcile loop")

	var instance kidlev1.IdlingResource
	if err := r.Get(ctx, req.NamespacedName, &instance); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var deployment v1.Deployment
	ref := instance.Spec.IdlingResourceRef
	key := types.NamespacedName{Namespace: instance.Namespace, Name: ref.Name}
	if err := r.Get(ctx, key, &deployment); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !instance.HasFinalizer(kidlev1.IdlingResourceFinalizerName) {
		logger.Info(fmt.Sprintf("AddFinalizer for %v", req.NamespacedName))
		if err := r.addFinalizer(ctx, &instance); err != nil {
			return reconcile.Result{}, fmt.Errorf("error when adding finalizer: %v", err)
		}
	}

	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if res, err := r.wakeup(ctx, instance, &deployment); err != nil {
			return res, err
		}
		if err := r.removeFinalizer(ctx, &instance); err != nil {
			return ctrl.Result{}, fmt.Errorf("error when deleting finalizer: %v", err)
		}
		return ctrl.Result{}, nil
	}

	// Create or update idle cronjob for the instance
	if instance.Spec.IdlingStrategy != nil && instance.Spec.IdlingStrategy.CronStrategy != nil {

		cjKey := types.NamespacedName{
			Namespace: instance.Namespace,
			Name:      k8s.ToDNSName("kidle", instance.Name, "idle"),
		}

		cronJob := &batchv1.CronJob{}
		if err := r.Get(ctx, cjKey, cronJob); err != nil {
			if errors.IsNotFound(err) {
				if err := r.CreateCronjob(ctx, cjKey, &instance); err != nil {
					return ctrl.Result{}, err
				}
			}
		}
		// TODO check that the cronjob spec is OK
		// Blah/foo/code
	} else {
		// TODO Delete the idle cronjob if necessary
		// Blah/foo/code
	}

	if needIdle(instance, &deployment) {
		return r.idle(ctx, &instance, &deployment)
	}

	if needWakeup(instance, &deployment) {
		return r.wakeup(ctx, instance, &deployment)
	}

	return ctrl.Result{}, nil
}

func needIdle(instance kidlev1.IdlingResource, deployment *v1.Deployment) bool {
	return instance.Spec.Idle && *deployment.Spec.Replicas > 0
}

func (r *IdlingResourceReconciler) idle(ctx context.Context, instance *kidlev1.IdlingResource, deployment *v1.Deployment) (ctrl.Result, error) {
	instance.Status.PreviousReplicas = deployment.Spec.Replicas

	deployment.Spec.Replicas = pointer.Int32(0)
	if err := r.Update(ctx, deployment); err != nil {
		return ctrl.Result{}, err
	}
	if err := r.Status().Update(ctx, instance); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func needWakeup(instance kidlev1.IdlingResource, deployment *v1.Deployment) bool {
	return !instance.Spec.Idle && *deployment.Spec.Replicas == 0
}

func (r *IdlingResourceReconciler) wakeup(ctx context.Context, instance kidlev1.IdlingResource, deployment *v1.Deployment) (ctrl.Result, error) {
	if instance.Spec.ResumeReplicas != nil {
		deployment.Spec.Replicas = instance.Spec.ResumeReplicas
	} else {
		var replicas int32 = 1
		deployment.Spec.Replicas = &replicas
	}
	if err := r.Update(ctx, deployment); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *IdlingResourceReconciler) addFinalizer(ctx context.Context, instance *kidlev1.IdlingResource) error {
	controllerutil.AddFinalizer(instance, kidlev1.IdlingResourceFinalizerName)
	err := r.Update(ctx, instance)
	if err != nil {
		return fmt.Errorf("failed to update idling resource finalizer: %v", err)
	}
	return nil
}

func (r *IdlingResourceReconciler) removeFinalizer(ctx context.Context, instance *kidlev1.IdlingResource) error {
	if array.ContainsString(instance.GetFinalizers(), kidlev1.IdlingResourceFinalizerName) {
		controllerutil.RemoveFinalizer(instance, kidlev1.IdlingResourceFinalizerName)
		err := r.Update(ctx, instance)
		if err != nil {
			return fmt.Errorf("error when removing idling resource finalizer: %v", err)
		}
	}
	return nil
}

func NewCronJob(key types.NamespacedName, instance *kidlev1.IdlingResource) *batchv1.CronJob {
	var cj = &batchv1.CronJob{
		TypeMeta: metav1.TypeMeta{
			Kind:       "CronJob",
			APIVersion: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.Name,
			Namespace: key.Namespace,
		},
		Spec: batchv1.CronJobSpec{
			Schedule: instance.Spec.IdlingStrategy.CronStrategy.Schedule,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						ObjectMeta: metav1.ObjectMeta{},
						Spec: corev1.PodSpec{
							RestartPolicy: corev1.RestartPolicyOnFailure,
							Containers: []corev1.Container{
								{
									Name:    "idle",
									Image:   "kidledev/kidlectl:main",
									Command: []string{"idle", instance.Spec.IdlingResourceRef.Name},
								},
							},
						},
					},
				},
			},
		},
	}
	return cj
}

func (r *IdlingResourceReconciler) CreateCronjob(ctx context.Context, key types.NamespacedName, instance *kidlev1.IdlingResource) error {
	cj := NewCronJob(key, instance)
	if err := controllerutil.SetControllerReference(instance, cj, r.Scheme); err != nil {
		return fmt.Errorf("unable to set controller reference for cronJob: %v", err)
	}
	if err := r.Create(ctx, cj); err != nil {
		return fmt.Errorf("unable to create cronJob: %v", err)
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IdlingResourceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kidlev1.IdlingResource{}).
		Owns(&batchv1.CronJob{}).
		Watches(
			&source.Kind{Type: &appsv1.Deployment{}},
			handler.EnqueueRequestsFromMapFunc(r.objectForIdlingResourceMapper),
		).
		Complete(r)
}

func (r *IdlingResourceReconciler) objectForIdlingResourceMapper(object client.Object) []reconcile.Request {
	if object.GetName() != "podinfo" {
		return nil
	}

	reqs := make([]reconcile.Request, 1)
	reqs[0].NamespacedName.Name = "idlingresource-sample"
	reqs[0].NamespacedName.Namespace = object.GetNamespace()
	return reqs
}
