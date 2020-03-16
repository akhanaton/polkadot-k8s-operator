package e2e

import (
	goctx "context"
	"github.com/operator-framework/operator-sdk/pkg/test/e2eutil"
	"github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis"
	polkadotv1alpha1 "github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/apis/polkadot/v1alpha1"
	"github.com/swisscom-blockchain/polkadot-k8s-operator/test/e2e/utils"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	"time"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
)

var (
	namespace string
	frameworkGlobal *framework.Framework
	ctx *framework.TestCtx
)

func TestMain(m *testing.M) {
	framework.MainEntry(m)
}

func TestPolkadot(t *testing.T) {
	polkadotList := &polkadotv1alpha1.PolkadotList{}
	err := framework.AddToFrameworkScheme(apis.AddToScheme,polkadotList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	ctx = framework.NewTestCtx(t)
	defer ctx.Cleanup()

	err = ctx.InitializeClusterResources(&framework.CleanupOptions{TestContext: ctx, Timeout: time.Second*10, RetryInterval: time.Second*10})
	if err != nil {
		t.Fatalf("failed to initialize cluster resources: %v", err)
	}

	namespace, err = ctx.GetNamespace()
	if err != nil {
		t.Fatal(err)
	}
	// get global framework
	frameworkGlobal = framework.Global
	// wait for the operator to be ready
	err = e2eutil.WaitForOperatorDeployment(t, frameworkGlobal.KubeClient, namespace, "polkadot-operator", 1, time.Second*5, time.Second*30 )
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestPolkadotSentry",testPolkadotSentry)
	t.Run("TestPolkadotValidator",testPolkadotValidator)
	t.Run("TestPolkadotSentryAndValidator",testPolkadotSentryAndValidator)
}

func testPolkadotSentry(t *testing.T) {

	polkadot := utils.NewPolkadotSentry(namespace)
	err := frameworkGlobal.Client.Create(goctx.TODO(), polkadot, &framework.CleanupOptions{TestContext: ctx, Timeout: time.Second * 5, RetryInterval: time.Second * 1})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestStatefulSetCreation", testStatefulSetCreationSentry)
	t.Run("TestServiceCreation", testServiceCreationSentry)

	err = frameworkGlobal.Client.Delete(goctx.TODO(), polkadot, &client.DeleteOptions{})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestStatefulSetDeletion", testStatefulSetDeletionSentry)
	t.Run("TestServiceDeletion", testServiceDeletionSentry)

}

func testPolkadotValidator(t *testing.T) {

	polkadot := utils.NewPolkadotValidator(namespace)
	err := frameworkGlobal.Client.Create(goctx.TODO(), polkadot, &framework.CleanupOptions{TestContext: ctx, Timeout: time.Second * 5, RetryInterval: time.Second * 1})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestStatefulSetCreation", testStatefulSetCreationValidator)
	t.Run("TestServiceCreation", testServiceCreationValidator)

	err = frameworkGlobal.Client.Delete(goctx.TODO(), polkadot, &client.DeleteOptions{})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestStatefulSetDeletion", testStatefulSetDeletionValidator)
	t.Run("TestServiceDeletion", testServiceDeletionValidator)

}

func testPolkadotSentryAndValidator(t *testing.T) {

	polkadot := utils.NewPolkadotSentryAndValidator(namespace, true)
	err := frameworkGlobal.Client.Create(goctx.TODO(), polkadot, &framework.CleanupOptions{TestContext: ctx, Timeout: time.Second * 5, RetryInterval: time.Second * 1})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestStatefulSetCreationSentry", testStatefulSetCreationSentry)
	t.Run("TestServiceCreationSentry", testServiceCreationSentry)
	t.Run("TestStatefulSetCreationValidator", testStatefulSetCreationValidator)
	t.Run("TestServiceCreationValidator", testServiceCreationValidator)
	t.Run("TestNetworkPolicyCreation", testNetworkPolicyCreation)

	err = frameworkGlobal.Client.Delete(goctx.TODO(), polkadot, &client.DeleteOptions{})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("TestStatefulSetDeletionSentry", testStatefulSetDeletionSentry)
	t.Run("TestServiceDeletionSentry", testServiceDeletionSentry)
	t.Run("TestStatefulSetDeletionValidator", testStatefulSetDeletionValidator)
	t.Run("TestServiceDeletionValidator", testServiceDeletionValidator)
	t.Run("TestNetworkPolicyDeletion", testNetworkPolicyDeletion)

}

func testStatefulSetCreationSentry(t *testing.T) {
	testStatefulSetCreation(t,"sentry-sset")
}

func testStatefulSetCreationValidator(t *testing.T) {
	testStatefulSetCreation(t,"validator-sset")
}

func testStatefulSetCreation(t *testing.T, name string) {
	err := utils.WaitForStatefulSet(t, frameworkGlobal.KubeClient, namespace, name, 1, time.Second*5, time.Second*30)
	if err != nil {
		t.Fatal(err)
	}
}

func testStatefulSetDeletionSentry(t *testing.T) {
	testStatefulSetDeletion(t,"sentry-sset")
}

func testStatefulSetDeletionValidator(t *testing.T) {
	testStatefulSetDeletion(t,"validator-sset")
}

func testStatefulSetDeletion(t *testing.T, name string) {
	err := utils.WaitForStatefulSetDelete(t, frameworkGlobal.KubeClient, namespace, name, time.Second*5, time.Second*30)
	if err != nil {
		t.Fatal(err)
	}
}

func testServiceCreationSentry(t *testing.T) {
	testServiceCreation(t, "sentry-service")
}

func testServiceCreationValidator(t *testing.T) {
	testServiceCreation(t, "validator-service")
}

func testServiceCreation(t *testing.T, name string) {
	err := utils.WaitForService(t, frameworkGlobal.KubeClient, namespace, name, time.Second*5, time.Second*30)
	if err != nil {
		t.Fatal(err)
	}
}

func testServiceDeletionSentry(t *testing.T) {
	testServiceDeletion(t, "sentry-service")
}

func testServiceDeletionValidator(t *testing.T) {
	testServiceDeletion(t, "validator-service")
}

func testServiceDeletion(t *testing.T, name string) {
	err := utils.WaitForServiceDelete(t, frameworkGlobal.KubeClient, namespace, name, time.Second*5, time.Second*30)
	if err != nil {
		t.Fatal(err)
	}
}

func testNetworkPolicyCreation(t *testing.T) {
	err := utils.WaitForNetworkPolicy(t, frameworkGlobal.KubeClient, namespace, "validator-networkpolicy", time.Second*5, time.Second*30)
	if err != nil {
		t.Fatal(err)
	}
}

func testNetworkPolicyDeletion(t *testing.T) {
	err := utils.WaitForNetworkPolicyDelete(t, frameworkGlobal.KubeClient, namespace, "validator-networkpolicy", time.Second*5, time.Second*30)
	if err != nil {
		t.Fatal(err)
	}
}

