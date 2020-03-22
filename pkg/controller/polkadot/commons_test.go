package polkadot

import (
	"context"
	"github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis"
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	v12 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"testing"
)

func TestCreateResource(t *testing.T) {

	tests := []struct {
		name        string
		newResource interface{}
	}{
		{
			name:        "Service request creation",
			newResource: getFakeService(ServiceSentryName, corev1.ServiceTypeClusterIP),
		},
		{
			name:        "NetworkPolicy request creation",
			newResource: getFakeNetworkPolicy(ValidatorNetworkPolicy, "status1"),
		},
		{
			name:        "StatefulSet request creation",
			newResource: getFakeStatefulSet(SentrySSName, 1),
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
	if err := v1.AddToScheme(scheme); err != nil {
		t.Errorf("apis.AddToScheme: %v", err)
	}
	if err := v12.AddToScheme(scheme); err != nil {
		t.Errorf("apis.AddToScheme: %v", err)
	}

	// Objects to track in the fake client.
	objs := []runtime.Object{polkadot}

	// Create a fake client to mock API calls.
	client := fake.NewFakeClientWithScheme(scheme, objs...)
	reconciler := ReconcilerPolkadot{client: client, scheme: scheme}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := reconciler.createResource(test.newResource, polkadot)
			if err != nil {
				t.Fatalf("createRsource: (%v)", err)
			}
		})
	}
}

func TestFetchResource(t *testing.T) {

	testsFound := []struct {
		name     string
		resourceName string
		resource interface{}
	}{
		{
			name:     "Service request found",
			resourceName: ServiceSentryName,
			resource: getFakeService(ServiceSentryName, corev1.ServiceTypeClusterIP),
		},
		{
			name:     "NetworkPolicy request found",
			resourceName: ValidatorNetworkPolicy,
			resource: getFakeNetworkPolicy(ValidatorNetworkPolicy, "status1"),
		},
		{
			name:     "StatefulSet request found",
			resourceName: SentrySSName,
			resource: getFakeStatefulSet(SentrySSName, 1),
		},
	}

	testsNotFound := []struct {
		name     string
		resourceName string
		resource interface{}
	}{
		{
			name:     "Service request NOT found",
			resourceName: ServiceSentryName,
			resource: getFakeService(ServiceSentryName, corev1.ServiceTypeClusterIP),
		},
		{
			name:     "NetworkPolicy request NOT found",
			resourceName: ValidatorNetworkPolicy,
			resource: getFakeNetworkPolicy(ValidatorNetworkPolicy, "status1"),
		},
		{
			name:     "StatefulSet request NOT found",
			resourceName: SentrySSName,
			resource: getFakeStatefulSet(SentrySSName, 1),
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
	if err := v1.AddToScheme(scheme); err != nil {
		t.Errorf("apis.AddToScheme: %v", err)
	}
	if err := v12.AddToScheme(scheme); err != nil {
		t.Errorf("apis.AddToScheme: %v", err)
	}

	for _, test := range testsFound {
		t.Run(test.name, func(t *testing.T) {

			// Objects to track in the fake client.
			objs := []runtime.Object{polkadot,test.resource.(runtime.Object)}

			// Create a fake client to mock API calls.
			client := fake.NewFakeClientWithScheme(scheme, objs...)
			reconciler := ReconcilerPolkadot{client: client, scheme: scheme}

			isNotFound,err := reconciler.fetchResource(test.resource,types.NamespacedName{Name: test.resourceName})

			if err != nil || isNotFound {
				t.Fatalf("getchRsource: (%v)", err)
			}
		})
	}

	for _, test := range testsNotFound {
		t.Run(test.name, func(t *testing.T) {

			// Objects to track in the fake client.
			objs := []runtime.Object{polkadot}

			// Create a fake client to mock API calls.
			client := fake.NewFakeClientWithScheme(scheme, objs...)
			reconciler := ReconcilerPolkadot{client: client, scheme: scheme}

			isNotFound,err := reconciler.fetchResource(test.resource,types.NamespacedName{Name: test.resourceName})

			if err != nil || isNotFound != true {
				t.Fatalf("fetchResource: (%v)", err)
			}
		})
	}
}

func TestUpdateResource(t *testing.T) {

	tests := []struct {
		name             string
		resourceName 	string
 		resource         Resource
		expectedResource Resource
	}{
		{
			name:             "Service request update",
			resourceName: 	serviceName,
			resource:         Resource{getFakeService(serviceName, corev1.ServiceTypeClusterIP)},
			expectedResource: Resource{getFakeService(serviceName, corev1.ServiceTypeNodePort)},
		},
		{
			name:             "NetworkPolicy request update",
			resourceName: 	ValidatorNetworkPolicy,
			resource:         Resource{getFakeNetworkPolicy(ValidatorNetworkPolicy, "status1")},
			expectedResource: Resource{getFakeNetworkPolicy(ValidatorNetworkPolicy, "status2")},
		},
		{
			name:             "StatefulSet request update",
			resourceName: 	SentrySSName,
			resource:         Resource{getFakeStatefulSet(SentrySSName, 1)},
			expectedResource: Resource{getFakeStatefulSet(SentrySSName, 2)},
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
	if err := v1.AddToScheme(scheme); err != nil {
		t.Errorf("apis.AddToScheme: %v", err)
	}
	if err := v12.AddToScheme(scheme); err != nil {
		t.Errorf("apis.AddToScheme: %v", err)
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Objects to track in the fake client.
			objs := []runtime.Object{polkadot, test.resource.obj.(runtime.Object)}
			// Create a fake client to mock API calls.
			client := fake.NewFakeClientWithScheme(scheme, objs...)
			reconciler := ReconcilerPolkadot{client: client, scheme: scheme}

			err := reconciler.client.Get(context.TODO(), types.NamespacedName{Name: test.resourceName, Namespace: corev1.NamespaceAll}, test.resource.obj.(runtime.Object))
			if err != nil {
				t.Fatalf("updateRsource: (%v)", err)
			}

			err = reconciler.updateResource(test.expectedResource.obj)
			if err != nil {
				t.Fatalf("updateRsource: (%v)", err)
			}

			err = reconciler.client.Get(context.TODO(), types.NamespacedName{Name: test.resourceName, Namespace: corev1.NamespaceAll}, test.resource.obj.(runtime.Object))
			if err != nil {
				t.Fatalf("updateRsource: (%v)", err)
			}

			if !reflect.DeepEqual(test.resource.getSpec(),test.expectedResource.getSpec()) {
				t.Fatalf("the request doesn't match the expected result:\n (%v) \n (%v)", test.resource.getSpec(), test.expectedResource.getSpec())
			}
		})
	}
}

func getFakePolkadot() *polkadotv1alpha1.Polkadot{
	// A Polkadot object with metadata and spec.
	return &polkadotv1alpha1.Polkadot{
		ObjectMeta: metav1.ObjectMeta{
			Name:      CRName,
			Namespace: corev1.NamespaceAll,
		},
	}
}

func getFakeService(name string, serviceType corev1.ServiceType) *corev1.Service {
	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.ServiceSpec{
			Type: serviceType,
		},
	}
	return s
}

func getFakeNetworkPolicy(name, testName string) *v1.NetworkPolicy{
	return &v1.NetworkPolicy{
		TypeMeta:   metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:name,
		},
		Spec:       v1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{"test": testName},
			},
		},
	}
}

func getFakeStatefulSet(name string, replicas int32) *v12.StatefulSet{
	return &v12.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:name,
		},
		Spec:       v12.StatefulSetSpec{
			Replicas:             &replicas,
		},
	}
}

func getFakeRequest() *reconcile.Request{
	return &reconcile.Request{NamespacedName:types.NamespacedName{Name:CRName}}
}

const(
	CRName = "polkadot-cr"
)

type Resource struct {
	obj interface{}
}
func (r *Resource) getSpec() interface{}{
	switch r.obj.(type) {
	case *corev1.Service:
		return r.obj.(*corev1.Service).Spec
	case *v1.NetworkPolicy:
		return r.obj.(*v1.NetworkPolicy).Spec
	case *v12.StatefulSet:
		return r.obj.(*v12.StatefulSet).Spec
	default:
		return nil
	}
}
