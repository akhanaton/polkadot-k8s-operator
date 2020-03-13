// Copyright (c) 2020 Swisscom Blockchain AG
// Licensed under MIT License
package polkadot

import (
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)

func getCommands(nodeKey,clientName string) []string{
	return []string{
		"polkadot",
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
		//"--no-telemetry",
	}
}

func newStatefulSetSentry(CRInstance *polkadotv1alpha1.Polkadot) *appsv1.StatefulSet {
	replicas := CRInstance.Spec.Sentry.Replicas
	version := CRInstance.Spec.ClientVersion
	clientName := CRInstance.Spec.Sentry.ClientName
	nodeKey := CRInstance.Spec.Sentry.NodeKey
	CPULimit := CRInstance.Spec.Sentry.CPULimit
	memoryLimit := CRInstance.Spec.Sentry.MemoryLimit
	storageClassName := CRInstance.Spec.Sentry.StorageClassName

	labels := getSentrylabels()

	commands := getCommands(nodeKey,clientName)
	commands = append(commands,"--sentry")
	if CRKind(CRInstance.Spec.Kind) == SentryAndValidator {
		reservedValidatorID := CRInstance.Spec.Sentry.ReservedValidatorID
		commands = append(commands, "--reserved-nodes", "/dns4/"+serviceValidatorName+"/tcp/30333/p2p/"+reservedValidatorID)
	}

	return getStatefulSet(sentrySSName,CRInstance.Namespace,labels,replicas,storageClassName,version,commands,CPULimit,memoryLimit)
}

func newStatefulSetValidator(CRInstance *polkadotv1alpha1.Polkadot) *appsv1.StatefulSet {
	replicas := int32(1)
	version := CRInstance.Spec.ClientVersion
	clientName := CRInstance.Spec.Validator.ClientName
	nodeKey := CRInstance.Spec.Validator.NodeKey
	CPULimit := CRInstance.Spec.Validator.CPULimit
	memoryLimit := CRInstance.Spec.Validator.MemoryLimit
	storageClassName := CRInstance.Spec.Validator.StorageClassName

	labels := getValidatorLabels()

	commands := getCommands(nodeKey,clientName)
	commands = append(commands,"--validator")
	if CRKind(CRInstance.Spec.Kind) == SentryAndValidator {
		reservedSentryID := CRInstance.Spec.Validator.ReservedSentryID
		commands = append(commands,
			"--reserved-only",
			"--reserved-nodes", "/dns4/"+serviceSentryName+"/tcp/30333/p2p/"+reservedSentryID)
	}

	return getStatefulSet(validatorSSName,CRInstance.Namespace,labels,replicas,storageClassName,version,commands,CPULimit,memoryLimit)
}

func getStatefulSet(name string, namespace string, labels map[string]string, replicas int32, storageClassName string, version string, commands []string, CPULimit,memoryLimit string) *appsv1.StatefulSet{
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    getCopyLabelsWithVersion(labels, version),
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
					Containers: []corev1.Container{
						{
						Name:           serviceName,
						Image:          imageName + ":" + version,
						VolumeMounts:   getVolumeMounts(),
						Command:        commands,
						Ports:          getContainerPortsClient(),
						LivenessProbe:  getHealthProbeClient(),
						ReadinessProbe: getHealthProbeClient(),
						Resources:      getResourceLimits(CPULimit,memoryLimit),
					},
					{
						Name:          "metrics-exporter",
						Image:         imageNameMetrics,
						Ports:         getContainerPortsMetrics(),
						LivenessProbe: getHealthProbeMetrics(),
					},
					},
				},
			},
		},
	}
}

func getVolumePermissionInitContainer() *corev1.Container {
	rootUser := int64(0)
	runAsNonRootFalse := false

	return &corev1.Container {
		Name:  "volume-mount-permissions-data",
		Image: "busybox",
		VolumeMounts: getVolumeMounts(),
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

func getVolumeMounts() []corev1.VolumeMount{
	return []corev1.VolumeMount{{
		Name:      volumeName,
		MountPath: volumeMountPath,
	}}
}

func getHealthProbeClient() *corev1.Probe{
	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/health",
				Port: intstr.IntOrString{Type: intstr.String, StrVal: RPCPortName},
			},
		},
		InitialDelaySeconds: 10,
		PeriodSeconds:       10,
	}
}

func getHealthProbeMetrics() *corev1.Probe{
	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path: "/metrics",
				Port: intstr.IntOrString{Type: intstr.String, StrVal: metricsPortName},
			},
		},
		InitialDelaySeconds: 10,
		PeriodSeconds:       10,
	}
}

func getResourceLimits(cpu,memory string) corev1.ResourceRequirements{
	return corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			"cpu":    resource.MustParse(cpu),
			"memory": resource.MustParse(memory),
		},
	}
}

func getContainerPortsClient() []corev1.ContainerPort{
	return []corev1.ContainerPort{
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
	}
}

func getContainerPortsMetrics() []corev1.ContainerPort{
	return []corev1.ContainerPort{
		{
			ContainerPort: metricsPort,
			Name:          metricsPortName,
		},
	}
}
