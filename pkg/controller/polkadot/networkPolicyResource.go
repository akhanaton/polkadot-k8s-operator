// Copyright (c) 2020 Swisscom Blockchain AG
// Licensed under MIT License
package polkadot

import (
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newNetworkPolicyValidator(CRInstance *polkadotv1alpha1.Polkadot) *v1.NetworkPolicy {
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
