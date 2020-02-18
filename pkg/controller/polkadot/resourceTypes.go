package polkadot

import (
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func newSentryStatefulSetForCR(CRInstance *polkadotv1alpha1.Polkadot) *appsv1.StatefulSet {
	replicas := CRInstance.Spec.Sentry.Replicas
	version := CRInstance.Spec.ClientVersion
	clientName := CRInstance.Spec.Sentry.ClientName
	nodeKey := CRInstance.Spec.Sentry.NodeKey
	CPULimit := CRInstance.Spec.Validator.CPULimit
	memoryLimit := CRInstance.Spec.Validator.MemoryLimit
	volumeName := "polkadot-volume"
	storageClassName := "default"
	serviceName := "polkadot"

	labels := getSentrylabels()
	labelsWithVersion := getCopyLabelsWithVersion(labels,version)

	commands := []string{
		"polkadot",
		"--sentry",
		"--node-key", nodeKey,
		"--name", clientName,
		"--unsafe-rpc-external", //TODO check the unsafeness
		"--unsafe-ws-external",
		"--rpc-cors=all",
		"--no-telemetry",
	}
	if CRKind(CRInstance.Spec.Kind) == SentryAndValidator {
		reservedValidatorID := CRInstance.Spec.Sentry.ReservedValidatorID
		commands = append(commands,"--reserved-nodes", "/dns4/polkadot-service-validator/tcp/30333/p2p/" + reservedValidatorID)
	}

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "polkadot-statefulset-" + "sentry",
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
						Command: commands,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 30333,
								Name: "p2p",
							},
							{
								ContainerPort: 9933,
								Name: "http-rpc",
							},
							{
								ContainerPort: 9944,
								Name: "websocket-rpc",
							},
						},
						//ReadinessProbe: &corev1.Probe{
						//	Handler:             corev1.Handler{
						//		HTTPGet: &corev1.HTTPGetAction{
						//			Path: "/health",
						//			Port: intstr.IntOrString{Type: intstr.String, StrVal:"http-rpc"},
						//		},
						//	},
						//	InitialDelaySeconds: 10,
						//	PeriodSeconds:       10,
						//},
						LivenessProbe: &corev1.Probe{
							Handler:             corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/health",
									Port: intstr.IntOrString{Type: intstr.String, StrVal:"http-rpc"},
								},
							},
							InitialDelaySeconds: 10,
							PeriodSeconds:       10,
						},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								"cpu": resource.MustParse(CPULimit),
								"memory": resource.MustParse(memoryLimit),
							},
						},
					}},
				},
			},
		},
	}
}

func newValidatorStatefulSetForCR(CRInstance *polkadotv1alpha1.Polkadot) *appsv1.StatefulSet {
	replicas := int32(1)
	version := CRInstance.Spec.ClientVersion
	clientName := CRInstance.Spec.Validator.ClientName
	nodeKey := CRInstance.Spec.Validator.NodeKey
	CPULimit := CRInstance.Spec.Validator.CPULimit
	memoryLimit := CRInstance.Spec.Validator.MemoryLimit
	volumeName := "polkadot-volume"
	storageClassName := "default"
	serviceName := "polkadot"

	labels := getValidatorLabels()
	labelsWithVersion := getCopyLabelsWithVersion(labels,version)

	commands := []string{
		"polkadot",
		"--validator",
		"--node-key", nodeKey,
		"--name", clientName,
		//"--unsafe-rpc-external", //TODO check the unsafeness
		//"--unsafe-ws-external",
		//"--rpc-cors=all",
		"--no-telemetry",
	}
	if CRKind(CRInstance.Spec.Kind) == SentryAndValidator {
		reservedSentryID := CRInstance.Spec.Validator.ReservedSentryID
		commands = append(commands,
			"--reserved-only",
			"--reserved-nodes", "/dns4/polkadot-service-sentry/tcp/30333/p2p/" + reservedSentryID)
	}

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "polkadot-statefulset-" + "validator",
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
						Command: commands,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: 30333,
								Name: "p2p",
							},
							{
								ContainerPort: 9933,
								Name: "http-rpc",
							},
							{
								ContainerPort: 9944,
								Name: "websocket-rpc",
							},
						},
						//ReadinessProbe: &corev1.Probe{
						//	Handler:             corev1.Handler{
						//		HTTPGet: &corev1.HTTPGetAction{
						//			Path: "/health",
						//			Port: intstr.IntOrString{Type: intstr.String, StrVal:"http-rpc"},
						//		},
						//	},
						//	InitialDelaySeconds: 10,
						//	PeriodSeconds:       10,
						//},
						//LivenessProbe: &corev1.Probe{
						//	Handler:             corev1.Handler{
						//		HTTPGet: &corev1.HTTPGetAction{
						//			Path: "/health",
						//			Port: intstr.IntOrString{Type: intstr.Int, IntVal: 9933},
						//		},
						//	},
						//	InitialDelaySeconds: 10,
						//	PeriodSeconds:       10,
						//},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								"cpu": resource.MustParse(CPULimit),
								"memory": resource.MustParse(memoryLimit),
							},
						},
					}},
				},
			},
		},
	}
}

func newSentryServiceForCR(CRInstance *polkadotv1alpha1.Polkadot) *corev1.Service {
	labels := getSentrylabels()
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "polkadot-service-sentry",
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
					Name:       "http-rpc",
					Port:       9933,
					TargetPort: intstr.FromInt(9933),
					Protocol:   "TCP",
				},
				{
					Name:       "websocket-rpc",
					Port:       9944,
					TargetPort: intstr.FromInt(9944),
					Protocol:   "TCP",
				},
			},
			Selector: labels,
		},
	}
}

func newValidatorServiceForCR(CRInstance *polkadotv1alpha1.Polkadot) *corev1.Service {
	labels := getValidatorLabels()
	serviceType := corev1.ServiceTypeClusterIP
	if CRKind(CRInstance.Spec.Kind) == Validator {
		serviceType = corev1.ServiceTypeNodePort
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "polkadot-service-validator",
			Namespace: CRInstance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type: serviceType,
			Ports: []corev1.ServicePort{
				{
					Name:       "p2p",
					Port:       30333,
					TargetPort: intstr.FromInt(30333),
					Protocol:   "TCP",
				},
				{
					Name:       "http-rpc",
					Port:       9933,
					TargetPort: intstr.FromInt(9933),
					Protocol:   "TCP",
				},
				{
					Name:       "websocket-rpc",
					Port:       9944,
					TargetPort: intstr.FromInt(9944),
					Protocol:   "TCP",
				},
			},
			Selector: labels,
		},
	}
}

func newValidatorNetworkPolicyForCR(CRInstance *polkadotv1alpha1.Polkadot) *v1.NetworkPolicy {
	labels := getValidatorLabels()
	sentryLalbels := getSentrylabels()
	
	return &v1.NetworkPolicy{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "polkadot-networkpolicy",
			Namespace: CRInstance.Namespace,
		},
		Spec:  	v1.NetworkPolicySpec{
					PodSelector: metav1.LabelSelector{
						MatchLabels: labels,
					},
				Ingress: []v1.NetworkPolicyIngressRule{{
					From: []v1.NetworkPolicyPeer{{
						PodSelector:  &metav1.LabelSelector{
							MatchLabels: sentryLalbels,
						},
					}},
				}},
				Egress: []v1.NetworkPolicyEgressRule{{
					To: []v1.NetworkPolicyPeer{{
						PodSelector:  &metav1.LabelSelector{
							MatchLabels: sentryLalbels,
						},
					}},
				}},
		},
	}
}

func getAppLabels() map[string]string{
	labels:= map[string]string{"app":"polkadot"}
	return labels
}

func getSentrylabels() map[string]string {
	labels := getAppLabels()
	labels["role"] = "sentry"
	return labels
}

func getCopyLabelsWithVersion(labels map[string]string, version string) map[string]string {
	newLabels := getCopy(labels)
	newLabels["version"] = version
	return newLabels
}

func getCopy(originalMap map[string]string) map[string]string {
	newMap := make(map[string]string)
	for key, value := range originalMap {
		newMap[key] = value
	}
	return newMap
}

func getValidatorLabels() map[string]string {
	labels := getAppLabels()
	labels["role"] = "validator"
	return labels
}