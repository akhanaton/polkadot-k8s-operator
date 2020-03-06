// Copyright (c) 2020 Swisscom Blockchain AG
// Licensed under MIT License
package polkadot

import (
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)

func newSentryStatefulSetForCR(CRInstance *polkadotv1alpha1.Polkadot) *appsv1.StatefulSet {
	replicas := CRInstance.Spec.Sentry.Replicas
	version := CRInstance.Spec.ClientVersion
	clientName := CRInstance.Spec.Sentry.ClientName
	nodeKey := CRInstance.Spec.Sentry.NodeKey
	CPULimit := CRInstance.Spec.Sentry.CPULimit
	memoryLimit := CRInstance.Spec.Sentry.MemoryLimit
	storageClassName := CRInstance.Spec.Sentry.StorageClassName
	serviceName := "polkadot"

	labels := getSentrylabels()
	labelsWithVersion := getCopyLabelsWithVersion(labels, version)

	commands := []string{
		"polkadot",
		"--sentry",
		"--node-key", nodeKey,
		"--name", clientName,
		"-d=" + volumeMountPath,
		"--port",
		strconv.Itoa(P2PPort),
		"--rpc-port",
		strconv.Itoa(RPCPort),
		"--ws-port",
		strconv.Itoa(WSPort),
		"--unsafe-rpc-external",
		"--unsafe-ws-external",
		"--rpc-cors=all",
		"--no-telemetry",
	}
	if CRKind(CRInstance.Spec.Kind) == SentryAndValidator {
		reservedValidatorID := CRInstance.Spec.Sentry.ReservedValidatorID
		commands = append(commands, "--reserved-nodes", "/dns4/"+serviceValidatorName+"/tcp/30333/p2p/"+reservedValidatorID)
	}

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      sentrySSName,
			Namespace: CRInstance.Namespace,
			Labels:    labelsWithVersion,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			ServiceName: serviceName,
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{ *getVolumeClaimTemplate(storageClassName) },
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					SecurityContext: getPodSecurityContext(),
					InitContainers: []corev1.Container{ *getVolumePermissionInitContainer() },
					Containers: []corev1.Container{{
						Name:  serviceName,
						Image: imageName + ":" + version,
						VolumeMounts: []corev1.VolumeMount{{
							Name:      volumeName,
							MountPath: volumeMountPath,
						}},
						Command: commands,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: P2PPort,
								Name:          P2PPortName,
							},
							{
								ContainerPort: RPCPort,
								Name:          RPCPortName,
							},
							{
								ContainerPort: WSPort,
								Name:          WSPortName,
							},
						},
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/health",
									Port: intstr.IntOrString{Type: intstr.String, StrVal: RPCPortName},
								},
							},
							InitialDelaySeconds: 10,
							PeriodSeconds:       10,
						},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								"cpu":    resource.MustParse(CPULimit),
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
	storageClassName := CRInstance.Spec.Validator.StorageClassName
	serviceName := "polkadot"

	labels := getValidatorLabels()
	labelsWithVersion := getCopyLabelsWithVersion(labels, version)

	commands := []string{
		"polkadot",
		"--validator",
		"--node-key", nodeKey,
		"--name", clientName,
		"--port",
		strconv.Itoa(P2PPort),
		"--rpc-port",
		strconv.Itoa(RPCPort),
		"--ws-port",
		strconv.Itoa(WSPort),
		"--unsafe-rpc-external",
		"--unsafe-ws-external",
		"--rpc-cors=all",
		"--no-telemetry",
	}
	if CRKind(CRInstance.Spec.Kind) == SentryAndValidator {
		reservedSentryID := CRInstance.Spec.Validator.ReservedSentryID
		commands = append(commands,
			"--reserved-only",
			"--reserved-nodes", "/dns4/"+serviceSentryName+"/tcp/30333/p2p/"+reservedSentryID)
	}

	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      validatorSSName,
			Namespace: CRInstance.Namespace,
			Labels:    labelsWithVersion,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			ServiceName: serviceName,
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{ *getVolumeClaimTemplate(storageClassName) },
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					SecurityContext: getPodSecurityContext(),
					InitContainers: []corev1.Container{ *getVolumePermissionInitContainer() },
					Containers: []corev1.Container{{
						Name:  serviceName,
						Image: imageName + ":" + version,
						VolumeMounts: []corev1.VolumeMount{{
							Name:      volumeName,
							MountPath: volumeMountPath,
						}},
						Command: commands,
						Ports: []corev1.ContainerPort{
							{
								ContainerPort: P2PPort,
								Name:          P2PPortName,
							},
							{
								ContainerPort: RPCPort,
								Name:          RPCPortName,
							},
							{
								ContainerPort: WSPort,
								Name:          WSPortName,
							},
						},
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path: "/health",
									Port: intstr.IntOrString{Type: intstr.String, StrVal: RPCPortName},
								},
							},
							InitialDelaySeconds: 10,
							PeriodSeconds:       10,
						},
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								"cpu":    resource.MustParse(CPULimit),
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
			Name:      serviceSentryName,
			Namespace: CRInstance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{
					Name:       P2PPortName,
					Port:       P2PPort,
					TargetPort: intstr.FromInt(P2PPort),
					Protocol:   "TCP",
				},
				{
					Name:       RPCPortName,
					Port:       RPCPort,
					TargetPort: intstr.FromInt(RPCPort),
					Protocol:   "TCP",
				},
				{
					Name:       WSPortName,
					Port:       WSPort,
					TargetPort: intstr.FromInt(WSPort),
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
			Name:      serviceValidatorName,
			Namespace: CRInstance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type: serviceType,
			Ports: []corev1.ServicePort{
				{
					Name:       P2PPortName,
					Port:       P2PPort,
					TargetPort: intstr.FromInt(P2PPort),
					Protocol:   "TCP",
				},
				{
					Name:       RPCPortName,
					Port:       RPCPort,
					TargetPort: intstr.FromInt(RPCPort),
					Protocol:   "TCP",
				},
				{
					Name:       WSPortName,
					Port:       WSPort,
					TargetPort: intstr.FromInt(WSPort),
					Protocol:   "TCP",
				},
			},
			Selector: labels,
		},
	}
}

func newValidatorNetworkPolicyForCR(CRInstance *polkadotv1alpha1.Polkadot) *v1.NetworkPolicy {
	labels := getValidatorLabels()
	sentryLabels := getSentrylabels()

	return &v1.NetworkPolicy{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      validatorNetworkPolicy,
			Namespace: CRInstance.Namespace,
		},
		Spec: v1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: labels,
			},
			Ingress: []v1.NetworkPolicyIngressRule{{
				From: []v1.NetworkPolicyPeer{{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: sentryLabels,
					},
				}},
			}},
			Egress: []v1.NetworkPolicyEgressRule{{
				To: []v1.NetworkPolicyPeer{{
					PodSelector: &metav1.LabelSelector{
						MatchLabels: sentryLabels,
					},
				}},
			}},
		},
	}
}

func getVolumePermissionInitContainer() *corev1.Container {
	rootUser := int64(0)
	runAsNonRootFalse := false

	return &corev1.Container {
		Name:  "volume-mount-permissions-data",
		Image: "busybox",
		VolumeMounts: []corev1.VolumeMount{{
			Name:      volumeName,
			MountPath: volumeMountPath,
		}},
		SecurityContext: &corev1.SecurityContext{
			RunAsUser:          &rootUser,
			RunAsNonRoot: &runAsNonRootFalse,
		},
		Command: []string{"sh", "-c", "chown -R 1000:1000 " + volumeMountPath},

	}
}

func getPodSecurityContext() *corev1.PodSecurityContext {
	user := int64(1000)
	group := int64(1000)
	runAsNonRoot := true

	return &corev1.PodSecurityContext {
		RunAsUser:          &user,
		FSGroup:         &group,
		RunAsGroup:      &group,
		RunAsNonRoot: &runAsNonRoot,
	}
}

func getVolumeClaimTemplate(storageClassName string) *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: volumeName,
		},
		Spec: corev1.PersistentVolumeClaimSpec{

			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					"storage": resource.MustParse(storageRequest),
				},
			},
			StorageClassName: &storageClassName,
		},
	}
}
