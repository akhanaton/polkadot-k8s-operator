// Copyright (c) 2020 Swisscom Blockchain AG
// Licensed under MIT License
package polkadot

import (
	"github.com/swisscom-blockchain/polkadot-k8s-operator/config"
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)

func getCommands(nodeKey,clientName string, isDataPersistenceEnabled, isMetricsSupportEnabled bool) []string{
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
	if isDataPersistenceEnabled == true {
		c = append(c,"-d=" + volumeMountPath)
	}
	if isMetricsSupportEnabled == true {
		c = append(c, "--prometheus-external", "--prometheus-port", strconv.Itoa(config.MetricsPortEnvVar.Value))
	}
	return c
}

type Parameters struct{
	name                     string
	namespace                string
	labels                   map[string]string
	replicas                 int32
	version                  string
	commands                 []string
	clientContainerResources corev1.ResourceRequirements
	dataPersistence          polkadotv1alpha1.DataPersistenceSupport
	isMetricsSupportEnabled  bool
}

func newStatefulSetSentry(CRInstance *polkadotv1alpha1.Polkadot) *appsv1.StatefulSet {
	replicas := CRInstance.Spec.Sentry.Replicas
	version := CRInstance.Spec.ClientVersion
	clientName := CRInstance.Spec.Sentry.ClientName
	nodeKey := CRInstance.Spec.Sentry.NodeKey
	clientContainerResources := CRInstance.Spec.Sentry.Resources
	dataPersistence := CRInstance.Spec.Sentry.DataPersistenceSupport
	isMetricsSupportEnabled := CRInstance.Spec.MetricsSupport.Enabled

	labels := getSentrylabels()

	commands := getCommands(nodeKey,clientName,dataPersistence.Enabled,isMetricsSupportEnabled)
	commands = append(commands,"--sentry")
	if CRKind(CRInstance.Spec.Kind) == SentryAndValidator {
		reservedValidatorID := CRInstance.Spec.Sentry.ReservedValidatorID
		commands = append(commands, "--reserved-nodes", "/dns4/"+ServiceValidatorName+"/tcp/30333/p2p/"+reservedValidatorID)
	}

	p := Parameters{
		name:                     SentrySSName,
		namespace:                CRInstance.Namespace,
		labels:                   labels,
		replicas:                 replicas,
		version:                  version,
		commands:                 commands,
		clientContainerResources: clientContainerResources,
		dataPersistence:          dataPersistence,
		isMetricsSupportEnabled:  isMetricsSupportEnabled,
	}

	return getStatefulSet(p)
}

func newStatefulSetValidator(CRInstance *polkadotv1alpha1.Polkadot) *appsv1.StatefulSet {
	replicas := int32(1)
	version := CRInstance.Spec.ClientVersion
	clientName := CRInstance.Spec.Validator.ClientName
	nodeKey := CRInstance.Spec.Validator.NodeKey
	clientContainerResources := CRInstance.Spec.Validator.Resources
	dataPersistence := CRInstance.Spec.Validator.DataPersistenceSupport
	isMetricsSupportEnabled := CRInstance.Spec.MetricsSupport.Enabled

	labels := getValidatorLabels()

	commands := getCommands(nodeKey,clientName,dataPersistence.Enabled,isMetricsSupportEnabled)
	commands = append(commands,"--validator")
	if CRKind(CRInstance.Spec.Kind) == SentryAndValidator {
		reservedSentryID := CRInstance.Spec.Validator.ReservedSentryID
		commands = append(commands,
			"--reserved-only",
			"--reserved-nodes", "/dns4/"+ServiceSentryName+"/tcp/30333/p2p/"+reservedSentryID)
	}

	p := Parameters{
		name:                     ValidatorSSName,
		namespace:                CRInstance.Namespace,
		labels:                   labels,
		replicas:                 replicas,
		version:                  version,
		commands:                 commands,
		clientContainerResources: clientContainerResources,
		dataPersistence:          dataPersistence,
		isMetricsSupportEnabled:  isMetricsSupportEnabled,
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
	if p.dataPersistence.Enabled == true{
		sSpec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{ p.dataPersistence.PersistentVolumeClaim }
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
	if p.dataPersistence.Enabled == true{
		spec.InitContainers = []corev1.Container{ *getVolumePermissionInitContainer(p.dataPersistence.PersistentVolumeClaim.ObjectMeta.Name) }
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
		if p.dataPersistence.Enabled == true{
			container.VolumeMounts=getVolumeMounts(p.dataPersistence.PersistentVolumeClaim.ObjectMeta.Name)
		}
		return container
}

func getVolumePermissionInitContainer(volumeMountName string) *corev1.Container {
	rootUser := int64(0)
	runAsNonRootFalse := false

	return &corev1.Container {
		Name:  "volume-mount-permissions-data",
		Image: "busybox",
		VolumeMounts: getVolumeMounts(volumeMountName),
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

func getVolumeMounts(volumeName string) []corev1.VolumeMount{
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
