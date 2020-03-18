// Copyright (c) 2020 Swisscom Blockchain AG
// Licensed under MIT License
package polkadot

import (
	"github.com/swisscom-blockchain/polkadot-k8s-operator/config"
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)

func getCommands(nodeKey,clientName, isDataPersistenceActive string) []string{
	c := []string{
		"polkadot",
		"--node-key", nodeKey,
		"--name", clientName,
		"--port",
		strconv.Itoa(config.P2PPortEnvVar.Value),
		"--rpc-port",
		strconv.Itoa(config.RPCPortEnvVar.Value),
		"--ws-port",
		strconv.Itoa(config.WSPortEnvVar.Value),
		"--unsafe-rpc-external",
		"--unsafe-ws-external",
		"--rpc-cors=all",
		//"--no-telemetry",
	}
	if isDataPersistenceActive == "true" {
		c = append(c,"-d=" + volumeMountPath)
	}
	return c
}

type Parameters struct{
	name string
	namespace string
	labels map[string]string
	replicas int32
	storageClassName string
	version string
	commands []string
	clientContainerResources corev1.ResourceRequirements
	isDataPersistenceActive string
	isMetricsSupportActive string
}

func newStatefulSetSentry(CRInstance *polkadotv1alpha1.Polkadot) *appsv1.StatefulSet {
	replicas := CRInstance.Spec.Sentry.Replicas
	version := CRInstance.Spec.ClientVersion
	clientName := CRInstance.Spec.Sentry.ClientName
	nodeKey := CRInstance.Spec.Sentry.NodeKey
	clientContainerResources := CRInstance.Spec.Sentry.Resources
	storageClassName := CRInstance.Spec.Sentry.StorageClassName
	isDataPersistenceActive := CRInstance.Spec.IsDataPersistenceActive
	isMetricsSupportActive := CRInstance.Spec.IsMetricsSupportActive

	labels := getSentrylabels()

	commands := getCommands(nodeKey,clientName,isDataPersistenceActive)
	commands = append(commands,"--sentry")
	if CRKind(CRInstance.Spec.Kind) == SentryAndValidator {
		reservedValidatorID := CRInstance.Spec.Sentry.ReservedValidatorID
		commands = append(commands, "--reserved-nodes", "/dns4/"+ServiceValidatorName+"/tcp/30333/p2p/"+reservedValidatorID)
	}

	p := Parameters{
		name:                    SentrySSName,
		namespace:               CRInstance.Namespace,
		labels:                  labels,
		replicas:                replicas,
		storageClassName:        storageClassName,
		version:                 version,
		commands:                commands,
		clientContainerResources:clientContainerResources,
		isDataPersistenceActive: isDataPersistenceActive,
		isMetricsSupportActive:  isMetricsSupportActive,
	}

	return getStatefulSet(p)
}

func newStatefulSetValidator(CRInstance *polkadotv1alpha1.Polkadot) *appsv1.StatefulSet {
	replicas := int32(1)
	version := CRInstance.Spec.ClientVersion
	clientName := CRInstance.Spec.Validator.ClientName
	nodeKey := CRInstance.Spec.Validator.NodeKey
	clientContainerResources := CRInstance.Spec.Validator.Resources
	storageClassName := CRInstance.Spec.Validator.StorageClassName
	isDataPersistenceActive := CRInstance.Spec.IsDataPersistenceActive
	isMetricsSupportActive := CRInstance.Spec.IsMetricsSupportActive

	labels := getValidatorLabels()

	commands := getCommands(nodeKey,clientName,isDataPersistenceActive)
	commands = append(commands,"--validator")
	if CRKind(CRInstance.Spec.Kind) == SentryAndValidator {
		reservedSentryID := CRInstance.Spec.Validator.ReservedSentryID
		commands = append(commands,
			"--reserved-only",
			"--reserved-nodes", "/dns4/"+ServiceSentryName+"/tcp/30333/p2p/"+reservedSentryID)
	}

	p := Parameters{
		name:                    ValidatorSSName,
		namespace:               CRInstance.Namespace,
		labels:                  labels,
		replicas:                replicas,
		storageClassName:        storageClassName,
		version:                 version,
		commands:                commands,
		clientContainerResources:clientContainerResources,
		isDataPersistenceActive: isDataPersistenceActive,
		isMetricsSupportActive:  isMetricsSupportActive,
	}

	return getStatefulSet(p)
}

func getStatefulSet(p Parameters) *appsv1.StatefulSet{
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.name,
			Namespace: p.namespace,
			Labels:    getCopyLabelsWithVersion(p.labels, p.version),
		},
		Spec: getStatefulSetSpec(p),
	}
}

func getStatefulSetSpec(p Parameters) appsv1.StatefulSetSpec{
	sSpec := appsv1.StatefulSetSpec{
		Replicas: &p.replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: p.labels,
		},
		ServiceName: serviceName,
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: p.labels,
			},
			Spec: getPodSpec(p),
		},
	}
	if p.isDataPersistenceActive == "true"{
		sSpec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{ *getVolumeClaimTemplate(p.storageClassName) }
	}
	return sSpec
}

func getPodSpec(p Parameters) corev1.PodSpec{
	spec := corev1.PodSpec{
		SecurityContext: getPodSecurityContext(),
		Containers: []corev1.Container{
			getContainerClient(p),
		},
	}
	if p.isDataPersistenceActive == "true"{
		spec.InitContainers = []corev1.Container{ *getVolumePermissionInitContainer() }
	}
	if p.isMetricsSupportActive == "true"{
		spec.Containers = append(spec.Containers, getContainerMetrics())
	}
	return spec
}

func getContainerClient(p Parameters) corev1.Container{
	container:=corev1.Container{
			Name:           serviceName,
			Image:          config.ImageClientEnvVar.Value + ":" + p.version,
			Command:        p.commands,
			Ports:          getContainerPortsClient(),
			LivenessProbe:  getHealthProbeClient(),
			ReadinessProbe: getHealthProbeClient(),
			Resources:     p.clientContainerResources,
		}
		if p.isDataPersistenceActive == "true"{
			container.VolumeMounts=getVolumeMounts()
		}
		return container
}

func getContainerMetrics() corev1.Container{
	return corev1.Container {
		Name:          "metrics-exporter",
		Image:         config.ImageMetricsEnvVar.Value,
		Ports:         getContainerPortsMetrics(),
		LivenessProbe: getHealthProbeMetrics(),
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
		Spec: corev1.PersistentVolumeClaimSpec{ //TODO make granular from here

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

func getContainerPortsClient() []corev1.ContainerPort{
	return []corev1.ContainerPort{
		{
			ContainerPort: int32(config.P2PPortEnvVar.Value),
			Name:          P2PPortName,
		},
		{
			ContainerPort: int32(config.RPCPortEnvVar.Value),
			Name:          RPCPortName,
		},
		{
			ContainerPort: int32(config.WSPortEnvVar.Value),
			Name:          WSPortName,
		},
	}
}

func getContainerPortsMetrics() []corev1.ContainerPort{
	return []corev1.ContainerPort{
		{
			ContainerPort: int32(config.MetricsPortEnvVar.Value),
			Name:          metricsPortName,
		},
	}
}
