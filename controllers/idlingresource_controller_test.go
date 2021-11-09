package controllers

import (
	"context"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/retry"
	"k8s.io/utils/pointer"
	"time"
	kidlev1 "tz/api/v1"
)

var _ = Describe("idling/wakeup Deployments", func() {
	const (
		timeout = time.Second * 10
		//duration = time.Second * 10
		interval = time.Millisecond * 250
	)
	var (
		ctx   = context.Background()
		irKey = types.NamespacedName{Name: "ir-idler-deploy", Namespace: "default"}
	)

	Context("Deployment suite", func() {
		var deployKey = types.NamespacedName{Name: "nginx", Namespace: "default"}
		var deploy = appsv1.Deployment{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Deployment",
				APIVersion: "apps/appsv1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      deployKey.Name,
				Namespace: deployKey.Namespace,
				Labels: map[string]string{
					"app": "nginx",
				},
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: pointer.Int32(1),
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "nginx",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "nginx",
						Namespace: "default",
						Labels: map[string]string{
							"app": "nginx",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:  "nginx",
								Image: "nginx",
							},
						},
					},
				},
			},
		}

		It("Should idle the Deployment", func() {

			By("Creating the IdlingResource object")
			ref := kidlev1.CrossVersionObjectReference{
				Kind:       "Deployment",
				Name:       deployKey.Name,
				APIVersion: "apps/appsv1",
			}
			Expect(k8sClient.Create(ctx, newIdlingResource(irKey, &ref))).Should(Succeed())

			By("Validation of the IdlingResource creation")
			createdIR := &kidlev1.IdlingResource{}
			Eventually(func() error {
				return k8sClient.Get(ctx, irKey, createdIR)
			}, timeout, interval).Should(Succeed())

			By("Creating the Deployment object")
			Expect(k8sClient.Create(ctx, &deploy)).Should(Succeed())

			d := &appsv1.Deployment{}
			Eventually(func() error {
				return k8sClient.Get(ctx, deployKey, d)
			}, timeout, interval).Should(Succeed())

			By("Idling the deployment")
			Expect(setIdleFlag(ctx, irKey, true)).Should(Succeed())

			By("Checking that Replicas == 0")
			Eventually(func() (*int32, error) {
				d := &appsv1.Deployment{}
				err := k8sClient.Get(ctx, deployKey, d)
				if err != nil {
					return nil, err
				}
				return d.Spec.Replicas, nil
			}, timeout, interval).Should(Equal(pointer.Int32(0)))
		})
	})
})

func newIdlingResource(key types.NamespacedName, ref *kidlev1.CrossVersionObjectReference) *kidlev1.IdlingResource {
	return &kidlev1.IdlingResource{
		TypeMeta: metav1.TypeMeta{
			Kind:       "IdlingResource",
			APIVersion: "kidle.kidle.dev/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.Name,
			Namespace: key.Namespace,
		},
		Spec: kidlev1.IdlingResourceSpec{
			IdlingResourceRef: *ref,
			Idle:              false,
		},
	}
}
func setIdleFlag(ctx context.Context, irKey types.NamespacedName, idle bool) error {
	ir := &kidlev1.IdlingResource{}
	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		if err := k8sClient.Get(ctx, irKey, ir); err != nil {
			return err
		}
		ir.Spec.Idle = idle
		return k8sClient.Update(ctx, ir)
	})
}
