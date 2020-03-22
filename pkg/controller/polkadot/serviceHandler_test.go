package polkadot

import (
	"github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestServiceGeneric(t *testing.T) {

	type testStruct []struct {
		name        string
		newResource *corev1.Service
	}
	testsOK := testStruct{
		{
			name:        "Service healthy",
			newResource: getFakeService(ServiceSentryName, corev1.ServiceTypeClusterIP),
		},
	}

	testsNotFound := testStruct{
		{
			name:        "Service not found",
			newResource: getFakeService(ServiceSentryName, corev1.ServiceTypeClusterIP),
		},
	}

	// A Polkadot object with metadata and spec.
	polkadot := getFakePolkadot()

	scheme := runtime.NewScheme()
	if err := apis.AddToScheme(scheme); err != nil {
		t.Errorf("apis.AddToScheme: %v", err)
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Errorf("apis.AddToScheme: %v", err)
	}

	for _, test := range testsOK {
		t.Run(test.name, func(t *testing.T) {
			// Objects to track in the fake client.
			objs := []runtime.Object{polkadot,test.newResource}

			// Create a fake client to mock API calls.
			client := fake.NewFakeClientWithScheme(scheme, objs...)
			reconciler := ReconcilerPolkadot{client: client, scheme: scheme}

			isRequeueForced, err := reconciler.handleServiceGeneric(polkadot,test.newResource)
			if isRequeueForced || err != nil {
				t.Fatalf("handleNetworkPolicy: (%v)", isRequeueForced)
			}
		})
	}

	for _, test := range testsNotFound{
		t.Run(test.name, func(t *testing.T) {
			// Objects to track in the fake client.
			objs := []runtime.Object{polkadot}

			// Create a fake client to mock API calls.
			client := fake.NewFakeClientWithScheme(scheme, objs...)
			reconciler := ReconcilerPolkadot{client: client, scheme: scheme}

			isRequeueForced, err := reconciler.handleServiceGeneric(polkadot,test.newResource)
			if !isRequeueForced || err != nil {
				t.Fatalf("handleNetworkPolicy: (%v)", isRequeueForced)
			}
		})
	}
}

