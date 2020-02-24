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

const (
	imageName              = "parity/polkadot"
	serviceSentryName      = "sentry-service"
	serviceValidatorName   = "validator-service"
	P2PPort                = 30333
	P2PPortName            = "p2p"
	RPCPort                = 9933
	RPCPortName            = "http-rpc"
	WSPort                 = 9944
	WSPortName             = "websocket-rpc"
	validatorSSName        = "validator-sset"
	sentrySSName           = "sentry-sset"
	validatorNetworkPolicy = "validator-networkpolicy"
	volumeMountPath        = "/polkadot"
	storageRequest         = "10Gi"
)

func newSentryStatefulSetForCR(CRInstance *polkadotv1alpha1.Polkadot) *appsv1.StatefulSet {
	replicas := CRInstance.Spec.Sentry.Replicas
	version := CRInstance.Spec.ClientVersion
	clientName := CRInstance.Spec.Sentry.ClientName
	nodeKey := CRInstance.Spec.Sentry.NodeKey
	CPULimit := CRInstance.Spec.Sentry.CPULimit
	memoryLimit := CRInstance.Spec.Sentry.MemoryLimit
	volumeName := "polkadot-volume"
	storageClassName := "default"
	serviceName := "polkadot"

	labels := getSentrylabels()
	labelsWithVersion := getCopyLabelsWithVersion(labels, version)

	commands := []string{
		"polkadot",
		"--sentry",
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
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
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
			}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
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
	volumeName := "polkadot-volume"
	storageClassName := "default"
	serviceName := "polkadot"
	//user := int64(1000)
	//group := int64(1000)

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
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{{
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
			}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					//SecurityContext: &corev1.PodSecurityContext{
					//	RunAsUser:          &user,
					//	RunAsGroup:         &group,
					//},
					Containers: []corev1.Container{{
						//SecurityContext: corev1.SecurityContext{
						//	RunAsUser:          &user,
						//	RunAsGroup:         &group,
						//},
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

func getAppLabels() map[string]string {
	labels := map[string]string{"app": "polkadot"}
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
