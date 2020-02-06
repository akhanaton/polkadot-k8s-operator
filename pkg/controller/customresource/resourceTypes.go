package customresource

import (
	cachev1alpha1 "github.com/ironoa/kubernetes-customresource-operator/pkg/apis/cache/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func newSentryStatefulSetForCR(CRInstance *cachev1alpha1.CustomResource) *appsv1.StatefulSet {
	labels := labelsForSentry()
	replicas := CRInstance.Spec.Size
	version := CRInstance.Spec.Version
	labelsWithVersion := labelsForSentryWithVersion(version)
	volumeName := "polkadot-volume"
	storageClassName := "default"
	serviceName := "polkadot"

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CRInstance.Name + "-set-" + "sentry",
			Namespace: CRInstance.Namespace,
			Labels:    labelsWithVersion,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			ServiceName: serviceName,
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
				ObjectMeta: metav1.ObjectMeta{
					Name: volumeName,
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					Resources:        corev1.ResourceRequirements{
						Requests: map[corev1.ResourceName]resource.Quantity{
							corev1.ResourceStorage: *resource.NewQuantity(5*1000*1000*1000, resource.DecimalSI), //5GB
						},
					},
					StorageClassName: &storageClassName,
				},
			}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  serviceName,
						Image: "parity/polkadot:" + version,
						VolumeMounts: []corev1.VolumeMount{{
							Name: volumeName,
							MountPath: "/data",
						}},
						Command: []string{
							"polkadot",
							"--sentry",
							"--node-key", "0000000000000000000000000000000000000000000000000000000000000013", // Local node id: QmQMTLWkNwGf7P5MQv7kUHCynMg7jje6h3vbvwd2ALPPhm
							"--reserved-nodes", "/dns4/polkadot-service-validator/tcp/30333/p2p/QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM",
							"--name", "IronoaSentry",
							//"--unsafe-rpc-external",
							//"--unsafe-ws-external",
							//"--rpc-cors=all",
							"--no-telemetry",
						},
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 30333,
							},
							{
								ContainerPort: 9933,
							},
							{
								ContainerPort: 9944,
							},
						},
					}},
				},
			},
		},
	}
}

func newValidatorStatefulSetForCR(CRInstance *cachev1alpha1.CustomResource) *appsv1.StatefulSet {
	labels := labelsForValidator()
	replicas := int32(1)
	version := CRInstance.Spec.Version
	labelsWithVersion := labelsForValidatorWithVersion(version)
	volumeName := "polkadot-volume"
	storageClassName := "default"
	serviceName := "polkadot"

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CRInstance.Name + "-set-" + "validator",
			Namespace: CRInstance.Namespace,
			Labels:    labelsWithVersion,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			ServiceName: serviceName,
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
				ObjectMeta: metav1.ObjectMeta{
					Name: volumeName,
				},
				Spec: corev1.PersistentVolumeClaimSpec{
					AccessModes:      []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
					Resources:        corev1.ResourceRequirements{
						Requests: map[corev1.ResourceName]resource.Quantity{
							corev1.ResourceStorage: *resource.NewQuantity(5*1000*1000*1000, resource.DecimalSI), //5GB
						},
					},
					StorageClassName: &storageClassName,
				},
			}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  serviceName,
						Image: "parity/polkadot:" + version,
						VolumeMounts: []corev1.VolumeMount{{
							Name: volumeName,
							MountPath: "/data",
						}},
						Command: []string{
							"polkadot",
							"--validator",
							"--node-key", "0000000000000000000000000000000000000000000000000000000000000021", // Local node id: QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM
							"--reserved-only",
							"--reserved-nodes", "/dns4/polkadot-service-sentry/tcp/30333/p2p/QmQMTLWkNwGf7P5MQv7kUHCynMg7jje6h3vbvwd2ALPPhm",
							"--name", "IronoaValidator",
							//"--unsafe-rpc-external",
							//"--unsafe-ws-external",
							//"--rpc-cors=all",
							"--no-telemetry",
						},
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 30333,
							},
							{
								ContainerPort: 9933,
							},
							{
								ContainerPort: 9944,
							},
						},
					}},
				},
			},
		},
	}
}

func newDeploymentForCR(CRInstance *cachev1alpha1.CustomResource) *appsv1.Deployment {
	labels := labelsForApp(CRInstance)
	replicas := CRInstance.Spec.Size
	version := CRInstance.Spec.Version
	labelsWithVersion := labelsForAppWithVersion(CRInstance, version)
	volumeName := "polkadot-data"

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CRInstance.Name + "-deployment",
			Namespace: CRInstance.Namespace,
			Labels:    labelsWithVersion,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{{
						Name:         volumeName,
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: getPVCName(CRInstance),
							},
						},
					}},
					Containers: []corev1.Container{{
						Name:  "polkadot",
						Image: "chevdor/polkadot:" + version,
						VolumeMounts: []corev1.VolumeMount{{
							Name: volumeName,
							MountPath: "/data",
						}},
						Command: []string{
							"polkadot", "--name", "Ironoa", "--rpc-external", "--rpc-cors=all",
						},
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 30333,
							},
							{
								ContainerPort: 9933,
							},
							{
								ContainerPort: 9944,
							},
						},
					}},
				},
			},
		},
	}
}

func newSentryServiceForCR(CRInstance *cachev1alpha1.CustomResource) *corev1.Service {
	labels := labelsForSentry()
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CRInstance.Name + "-service-sentry",
			Namespace: CRInstance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Name:       "p2p",
					Port:       30333,
					TargetPort: intstr.FromInt(30333),
					Protocol:   "TCP",
				},
				{
					Name:       "rpc",
					Port:       9933,
					TargetPort: intstr.FromInt(9933),
					Protocol:   "TCP",
				},
				{
					Name:       "to-be-defined-c",
					Port:       9944,
					TargetPort: intstr.FromInt(9944),
					Protocol:   "TCP",
				},
			},
			Selector: labels,
		},
	}
}

func newValidatorServiceForCR(CRInstance *cachev1alpha1.CustomResource) *corev1.Service {
	labels := labelsForValidator()
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CRInstance.Name + "-service-validator",
			Namespace: CRInstance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Name:       "p2p",
					Port:       30333,
					TargetPort: intstr.FromInt(30333),
					Protocol:   "TCP",
				},
				{
					Name:       "rpc",
					Port:       9933,
					TargetPort: intstr.FromInt(9933),
					Protocol:   "TCP",
				},
				{
					Name:       "to-be-defined-c",
					Port:       9944,
					TargetPort: intstr.FromInt(9944),
					Protocol:   "TCP",
				},
			},
			Selector: labels,
		},
	}
}

func newServiceForCR(CRInstance *cachev1alpha1.CustomResource) *corev1.Service {
	labels := labelsForApp(CRInstance)
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CRInstance.Name + "-service-sentry",
			Namespace: CRInstance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Name:       "p2p",
					Port:       30333,
					TargetPort: intstr.FromInt(30333),
					Protocol:   "TCP",
				},
				{
					Name:       "rpc",
					Port:       9933,
					TargetPort: intstr.FromInt(9933),
					Protocol:   "TCP",
				},
				{
					Name:       "to-be-defined-c",
					Port:       9944,
					TargetPort: intstr.FromInt(9944),
					Protocol:   "TCP",
				},
			},
			Selector: labels,
		},
	}
}

func newPVCForCR(CRInstance *cachev1alpha1.CustomResource) *corev1.PersistentVolumeClaim {
	labels := labelsForApp(CRInstance)
	storageClassName := "polkadot"
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getPVCName(CRInstance),
			Namespace: CRInstance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &storageClassName,
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("2Gi"),
				},
			},
		},
	}
}

func labelsForSentry() map[string]string {
	labels := map[string]string{"app":"sentry"}
	labels["app"] = "sentry"
	return labels
}

func labelsForSentryWithVersion(version string) map[string]string {
	labels := labelsForSentry()
	labels["version"] = version
	return labels
}

func labelsForValidator() map[string]string {
	labels := map[string]string{"app":"validator"}
	labels["app"] = "sentry"
	return labels
}

func labelsForValidatorWithVersion(version string) map[string]string {
	labels := labelsForSentry()
	labels["version"] = version
	return labels
}

// labelsForApp creates a simple set of labels for App.
func labelsForApp(cr *cachev1alpha1.CustomResource) map[string]string {
	return map[string]string{"app": cr.Name, "app_cr": cr.Name}
}

func labelsForAppWithVersion(cr *cachev1alpha1.CustomResource, version string) map[string]string {
	labels := labelsForApp(cr)
	labels["version"] = version
	return labels
}

func matchingLabels(cr *cachev1alpha1.CustomResource) map[string]string {
	return map[string]string{
		"app":    cr.Name,
		"server": cr.Name,
	}
}

func serverLabels(cr *cachev1alpha1.CustomResource) map[string]string {
	labels := map[string]string{
		"version": cr.Spec.Version,
	}
	for k, v := range matchingLabels(cr) {
		labels[k] = v
	}
	return labels
}

func getPVCName(CRInstance *cachev1alpha1.CustomResource) string {
	return CRInstance.Name + "-pvc"
}