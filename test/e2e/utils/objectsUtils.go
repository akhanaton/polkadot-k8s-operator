package utils

import (
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
)

func NewPolkadotSentry(namespace string) *polkadotv1alpha1.Polkadot{
	return &polkadotv1alpha1.Polkadot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-polkadot",
			Namespace: namespace,
		},
		Spec: polkadotv1alpha1.PolkadotSpec{
			ClientVersion:         "latest",
			Kind:                  "Sentry",
			Sentry:                *newSentry(),
			IsMetricsSupportActive: "false",
			IsDataPersistenceActive: "false",
		},
	}
}

func NewPolkadotValidator(namespace string) *polkadotv1alpha1.Polkadot{
	return &polkadotv1alpha1.Polkadot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-polkadot",
			Namespace: namespace,
		},
		Spec: polkadotv1alpha1.PolkadotSpec{
			ClientVersion:         "latest",
			Kind:                  "Validator",
			Validator:             *newValidator(),
		},
	}
}

func NewPolkadotSentryAndValidator(namespace string, isNetworkPolicyActive bool) *polkadotv1alpha1.Polkadot{
	return &polkadotv1alpha1.Polkadot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-polkadot",
			Namespace: namespace,
		},
		Spec: polkadotv1alpha1.PolkadotSpec{
			ClientVersion:         "latest",
			Kind:                  "SentryAndValidator",
			Sentry:                *newSentry(),
			Validator:             *newValidator(),
			IsNetworkPolicyActive: strconv.FormatBool(isNetworkPolicyActive),
		},
	}
}

func newSentry() *polkadotv1alpha1.Sentry{
	return &polkadotv1alpha1.Sentry{
		Replicas:            1,
		ClientName:          "IronoaSentry",
		NodeKey:             "0000000000000000000000000000000000000000000000000000000000000013",
		ReservedValidatorID: "QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM",
		Resources: v1.ResourceRequirements{
			Limits:   v1.ResourceList{
				"cpu":    resource.MustParse("0.2"),
				"memory": resource.MustParse("100Mi"),
			},
		},
		StorageClassName:    "local",
	}
}

func newValidator() *polkadotv1alpha1.Validator{
	return &polkadotv1alpha1.Validator{
		ClientName:          "IronoaValidator",
		NodeKey:             "0000000000000000000000000000000000000000000000000000000000000021",
		ReservedSentryID: 	 "QmQMTLWkNwGf7P5MQv7kUHCynMg7jje6h3vbvwd2ALPPhm",
		Resources: v1.ResourceRequirements{
			Limits:   v1.ResourceList{
				"cpu":    resource.MustParse("0.2"),
				"memory": resource.MustParse("100Mi"),
			},
		},
		StorageClassName:    "local",
	}
}