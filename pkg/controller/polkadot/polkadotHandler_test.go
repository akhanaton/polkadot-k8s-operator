package polkadot

import (
	"github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func TestHandleCustomResource(t *testing.T) {

	type testStruct []struct {
		name    string
		request *reconcile.Request
	}

	testsOK := testStruct{
		{
			name:    "Polkadot healthy",
			request: getFakeRequest(),
		},
	}

	testsNotFound := testStruct{
		{
			name:    "Polkadot not found",
			request: getFakeRequest(),
		},
	}

	// A Polkadot object with metadata and spec.
	polkadot := getFakePolkadot()

	scheme := runtime.NewScheme()
	if err := apis.AddToScheme(scheme); err != nil {
		t.Errorf("apis.AddToScheme: %v", err)
	}

	for _, test := range testsOK {
		t.Run(test.name, func(t *testing.T) {
			// Objects to track in the fake client.
			objs := []runtime.Object{polkadot}

			// Create a fake client to mock API calls.
			client := fake.NewFakeClientWithScheme(scheme, objs...)
			reconciler := ReconcilerPolkadot{client: client, scheme: scheme}

			polkadot, err := reconciler.handleCustomResource(*test.request)
			if polkadot == nil || err != nil {
				t.Fatalf("handleNetworkPolicy: (%v)", err)
			}
		})
	}

	for _, test := range testsNotFound{
		t.Run(test.name, func(t *testing.T) {
			// Objects to track in the fake client.
			objs := []runtime.Object{}

			// Create a fake client to mock API calls.
			client := fake.NewFakeClientWithScheme(scheme, objs...)
			reconciler := ReconcilerPolkadot{client: client, scheme: scheme}

			polkadot, err := reconciler.handleCustomResource(*test.request)
			if polkadot != nil || err != nil {
				t.Fatalf("handleNetworkPolicy: (%v)", err)
			}
		})
	}
}



