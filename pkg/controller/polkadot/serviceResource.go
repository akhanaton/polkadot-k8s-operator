// Copyright (c) 2020 Swisscom Blockchain AG
// Licensed under MIT License
package polkadot

import (
	"github.com/swisscom-blockchain/polkadot-k8s-operator/config"
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func newServiceSentry(CRInstance *polkadotv1alpha1.Polkadot) *corev1.Service {
	labels := getSentrylabels()
	return getService(ServiceSentryName,CRInstance,labels,corev1.ServiceTypeNodePort)
}

func newServiceValidator(CRInstance *polkadotv1alpha1.Polkadot) *corev1.Service {
	labels := getValidatorLabels()
	serviceType := corev1.ServiceTypeClusterIP
	if CRKind(CRInstance.Spec.Kind) == Validator {
		serviceType = corev1.ServiceTypeNodePort
	}
	return getService(ServiceValidatorName,CRInstance,labels,serviceType)
}

func getService(name string, CRInstance *polkadotv1alpha1.Polkadot, labels  map[string]string, serviceType corev1.ServiceType) *corev1.Service{
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: CRInstance.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type:     serviceType,
			Ports:    getServicePorts(CRInstance),
			Selector: labels,
		},
	}
}

func getServicePorts(CRInstance *polkadotv1alpha1.Polkadot) []corev1.ServicePort{
	isMetricsSupportEnabled := CRInstance.Spec.MetricsSupport.Enabled

	service := []corev1.ServicePort{
		{
			Name:       P2PPortName,
			Port:       int32(config.P2PPortEnvVar.Value),
			TargetPort: intstr.FromInt(config.P2PPortEnvVar.Value),
			Protocol:   "TCP",
		},
		{
			Name:       RPCPortName,
			Port:       int32(config.RPCPortEnvVar.Value),
			TargetPort: intstr.FromInt(config.RPCPortEnvVar.Value),
			Protocol:   "TCP",
		},
		{
			Name:       WSPortName,
			Port:       int32(config.WSPortEnvVar.Value),
			TargetPort: intstr.FromInt(config.WSPortEnvVar.Value),
			Protocol:   "TCP",
		},
	}

	if isMetricsSupportEnabled == true{
		service = append(service,*getMetricsPort())
	}

	return service
}

func getMetricsPort() *corev1.ServicePort{
	return &corev1.ServicePort{
		Name:       metricsPortName,
		Port:       int32(config.MetricsPortEnvVar.Value),
		TargetPort: intstr.FromInt(config.MetricsPortEnvVar.Value),
		Protocol:   "TCP",
	}
}
