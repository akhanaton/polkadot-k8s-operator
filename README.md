Copyright (c) 2020 Swisscom Blockchain AG  
Licensed under MIT License

# PolkaKop - Kubernetes Operator for Polkadot 

Kubernetes Operator for Polkadot Sentry and Validators nodes.

Client - Rust implementation of the Polkadot Host: https://github.com/paritytech/polkadot

The Polkadot Project: https://wiki.polkadot.network/en/

## Table Of Contents

* [Polkadot Custom Resource](#polkadot-custom-resource)  
* [Requirements](#requirements)  
    * [Optionals](#optionals)  
* [How To Run](#how-to-run)  
* [Clean up resources](#clean-up-resources)  
* [How To Run the Unit Tests](#how-to-run-the-unit-tests)  
* [How To Tutorial with Minikube](#how-to-tutorial-with-minikube)  
    * [Clone the repository](#clone-the-repository)  
    * [Parameters tuning](#parameters-tuning)  
    * [Deployment phase](#deployment-phase)  
* [Operator Configurable Environment Variables](#operator-configurable-environment-variables)     
* [Polkadot CR Configurable Parameters](#polkadot-cr-configurable-parameters)  
* [Updating of Node Versions](#updating-of-node-versions)  
* [Node Cluster Scaling Support](#node-cluster-scaling-support)  
    * [Please Note](#please-note)  
* [Secure Communications (Kind:SentryAndValidator)](#secure-communications-kindsentryandvalidator)  
* [Network Policies](#network-policies)  
    * [Default configuration](#default-configuration)  
    * [Prerequisites](#prerequisites)  
    * [Azure Example](#azure-example)  
* [Data Persistence Support](#data-persistence-support)  
    * [How To Tutorial with Minikube](#how-to-tutorial-with-minikube-1)  
* [Metrics Support](#metrics-support)  
    * [Default configuration](#default-configuration-1)  
    * [How to access to the metrics: Example in Minikube](#how-to-access-to-the-metrics-example-in-minikube)  
* [E2E Testing](#e2e-testing)  
    * [Build and run test](#build-and-run-test)  
* [About Kubernetes](#about-kubernetes)  
* [About the Operator](#about-the-operator)  
* [About the Operator SDK and Framework](#about-the-operator-sdk-and-framework)  
* [Project Directory Structure](#project-directory-structure)  

## Polkadot Custom Resource 

The deployable CR (Custom Resource) is called "Polkadot".  

```yaml
# Copyright (c) 2020 Swisscom Blockchain AG
# Licensed under MIT License
apiVersion: polkadot.swisscomblockchain.com/v1alpha1
kind: Polkadot
metadata:
  name: polkadot-cr
spec:
  clientVersion: latest
  kind: "SentryAndValidator"
  secureCommunicationSupport:
    enabled: true
  metricsSupport:
    enabled: true
  sentry:
    replicas: 1
    clientName: "IronoaSentry"
    nodeKey: "0000000000000000000000000000000000000000000000000000000000000013" # Local node id: QmQMTLWkNwGf7P5MQv7kUHCynMg7jje6h3vbvwd2ALPPhm
    reservedValidatorID: "QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM"
    resources:
      limits:
        memory: "512Mi"
        cpu: "0.5"
    dataPersistenceSupport:
      enabled: false
  validator:
    clientName: "IronoaValidator"
    nodeKey: "0000000000000000000000000000000000000000000000000000000000000021" # Local node id: QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM
    reservedSentryID: "QmQMTLWkNwGf7P5MQv7kUHCynMg7jje6h3vbvwd2ALPPhm"
    resources:
      limits:
        memory: "512Mi"
        cpu: "0.5"
    dataPersistenceSupport:
      enabled: true
      persistentVolumeClaim:
        metadata:
          name: "polkadot-volume"
        spec:
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 10Gi
          storageClassName: managed-premium #[default,managed-premium,local]
```

## Requirements

* Docker  
Mac: https://docs.docker.com/docker-for-mac/  
Linux: https://docs.docker.com/install/linux/docker-ce/ubuntu/

* The Kubernetes command-line tool, kubectl  
Mac: https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-with-homebrew-on-macos  
Linux: https://kubernetes.io/docs/tasks/tools/install-kubectl/#install-using-native-package-management

* Go compiler  
Mac, Linux, Windows, official website: https://golang.org/doc/install  
Mac, Homebrew: https://ahmadawais.com/install-go-lang-on-macos-with-homebrew/

* Access to a Container Registry  
Docker Hub: https://hub.docker.com/signup

* Operator SDK CLI tool  
Mac, Homebrew: https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md#install-from-homebrew-macos  
Linux: https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md#install-from-github-release

### Optionals

* Kubernetes Cluster Network Plugin: network plugin, see the [Secure Communications section (SentryAndValidator secure deployment)](#secure-communications-kindsentryandvalidator)

* Minikube, a tool that runs a single-node Kubernetes cluster in your local environment   
Mac: https://kubernetes.io/docs/tasks/tools/install-minikube/#install-minikube  
Linux: https://kubernetes.io/docs/tasks/tools/install-minikube/#install-minikube-using-a-package

## How To Run

Deploy to your favorite kubernetes cloud provided cluster (even minikube) a Custom Controller and a Polkadot Custom Resource. The Controller will create and supervise all the necessary resources needed to run a Polkadot Client configuration.

0. Configure your kubectl to work with your desired Kubernetes cluster 
    (e.g. Azure: az aks get-credentials --resource-group myResourceGroup --name myAKSCluster)
1. Clone the repository locally
2. Configure your application, the single entry points are:
    * deploy/operator.yaml
    * scripts/config/config.sh  
    In both these files configure the images to point to your favourite Container Registry
    * deploy/crds/polkadot.swisscomblockchain.com_v1alpha1_polkadot_cr.yaml
3. execute scripts/init.sh

## Clean up resources

Execute scripts/wipeAll.sh

## How To Run the Unit Tests

Move to the controller folder and then run the tests as normal GO unit tests:

```sh
# move to the controller folder
$ cd pkg/controller/polkadot

# run the tests (-v = verbose)
$ go test -v
=== RUN   TestCreateResource
=== RUN   TestCreateResource/Service_request_creation
=== RUN   TestCreateResource/NetworkPolicy_request_creation
=== RUN   TestCreateResource/StatefulSet_request_creation
--- PASS: TestCreateResource (0.00s)
    --- PASS: TestCreateResource/Service_request_creation (0.00s)
    --- PASS: TestCreateResource/NetworkPolicy_request_creation (0.00s)
    --- PASS: TestCreateResource/StatefulSet_request_creation (0.00s)
=== RUN   TestFetchResource
=== RUN   TestFetchResource/Service_request_found
=== RUN   TestFetchResource/NetworkPolicy_request_found
=== RUN   TestFetchResource/StatefulSet_request_found
=== RUN   TestFetchResource/Service_request_NOT_found
=== RUN   TestFetchResource/NetworkPolicy_request_NOT_found
=== RUN   TestFetchResource/StatefulSet_request_NOT_found
--- PASS: TestFetchResource (0.01s)
    --- PASS: TestFetchResource/Service_request_found (0.00s)
    --- PASS: TestFetchResource/NetworkPolicy_request_found (0.00s)
    --- PASS: TestFetchResource/StatefulSet_request_found (0.01s)
    --- PASS: TestFetchResource/Service_request_NOT_found (0.00s)
    --- PASS: TestFetchResource/NetworkPolicy_request_NOT_found (0.00s)
    --- PASS: TestFetchResource/StatefulSet_request_NOT_found (0.00s)
=== RUN   TestUpdateResource
=== RUN   TestUpdateResource/Service_request_update
=== RUN   TestUpdateResource/NetworkPolicy_request_update
=== RUN   TestUpdateResource/StatefulSet_request_update
--- PASS: TestUpdateResource (0.00s)
    --- PASS: TestUpdateResource/Service_request_update (0.00s)
    --- PASS: TestUpdateResource/NetworkPolicy_request_update (0.00s)
    --- PASS: TestUpdateResource/StatefulSet_request_update (0.00s)
=== RUN   TestHandleNetworkPolicyGeneric
=== RUN   TestHandleNetworkPolicyGeneric/NetworkPolicy_healthy
=== RUN   TestHandleNetworkPolicyGeneric/NetworkPolicy_not_found
--- PASS: TestHandleNetworkPolicyGeneric (0.00s)
    --- PASS: TestHandleNetworkPolicyGeneric/NetworkPolicy_healthy (0.00s)
    --- PASS: TestHandleNetworkPolicyGeneric/NetworkPolicy_not_found (0.00s)
=== RUN   TestHandleCustomResource
=== RUN   TestHandleCustomResource/Polkadot_healthy
=== RUN   TestHandleCustomResource/Polkadot_not_found
--- PASS: TestHandleCustomResource (0.00s)
    --- PASS: TestHandleCustomResource/Polkadot_healthy (0.00s)
    --- PASS: TestHandleCustomResource/Polkadot_not_found (0.00s)
=== RUN   TestServiceGeneric
=== RUN   TestServiceGeneric/Service_healthy
=== RUN   TestServiceGeneric/Service_not_found
--- PASS: TestServiceGeneric (0.00s)
    --- PASS: TestServiceGeneric/Service_healthy (0.00s)
    --- PASS: TestServiceGeneric/Service_not_found (0.00s)
=== RUN   TestHandleStatefulSetGeneric
=== RUN   TestHandleStatefulSetGeneric/StatefulSet_healthy
=== RUN   TestHandleStatefulSetGeneric/StatefulSet_not_found
--- PASS: TestHandleStatefulSetGeneric (0.00s)
    --- PASS: TestHandleStatefulSetGeneric/StatefulSet_healthy (0.00s)
    --- PASS: TestHandleStatefulSetGeneric/StatefulSet_not_found (0.00s)
PASS
ok      github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/controller/polkadot    0.305s
```

## How To Tutorial with Minikube

### Clone the repository

```sh
# Clone the repository
$ git clone https://github.com/swisscom-blockchain/polkadot-k8s-operator.git
$ cd polkadot-k8s-operator
```

### Parameters tuning

Example of a deployable deploy/crds/polkadot.swisscomblockchain.com_v1alpha1_polkadot_cr.yaml in a "SentryAndValidator" configuration.  
Note that if you deploy the operator locally, it is important to limit the the CPU and the memory usage (due to minikube limitations)
```yaml
# Copyright (c) 2020 Swisscom Blockchain AG
# Licensed under MIT License
apiVersion: polkadot.swisscomblockchain.com/v1alpha1
kind: Polkadot
metadata:
  name: polkadot-cr
spec:
  clientVersion: latest
  kind: "SentryAndValidator"
  secureCommunicationSupport:
    enabled: false
  metricsSupport:
    enabled: false
  sentry:
    replicas: 1
    clientName: "IronoaSentry"
    nodeKey: "0000000000000000000000000000000000000000000000000000000000000013" # Local node id: QmQMTLWkNwGf7P5MQv7kUHCynMg7jje6h3vbvwd2ALPPhm
    reservedValidatorID: "QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM"
    resources:
      limits:
        memory: "512Mi"
        cpu: "0.5"
    dataPersistenceSupport:
      enabled: false
  validator:
    clientName: "IronoaValidator"
    nodeKey: "0000000000000000000000000000000000000000000000000000000000000021" # Local node id: QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM
    reservedSentryID: "QmQMTLWkNwGf7P5MQv7kUHCynMg7jje6h3vbvwd2ALPPhm"
    resources:
      limits:
        memory: "512Mi"
        cpu: "0.5"
    dataPersistenceSupport:
      enabled: false
```

Example of a deployable deploy/operator.yaml, configured to work with my docker hub account (please change the image parameter).
```yaml
# Copyright (c) 2020 Swisscom Blockchain AG
# Licensed under MIT License
apiVersion: apps/v1
kind: Deployment
metadata:
  name: polkadot-operator
spec:
  replicas: 1
  selector:
    matchLabels:
      name: polkadot-operator
  template:
    metadata:
      labels:
        name: polkadot-operator
    spec:
      serviceAccountName: polkadot-operator
      containers:
        - name: polkadot-operator
          image: ironoa/customresource-operator:v0.0.8 #define your favourite
          command:
          - polkadot-operator
          imagePullPolicy: Always
          env:
            - name: OPERATOR_NAME
              value: "polkadot-operator"
            - name: CONTROLLER_NAME
              value: "polkadot-controller"
            - name: IMAGE_CLIENT
              value: "parity/polkadot"
            - name: METRICS_PORT
              value: "9615"
            - name: P2P_PORT
              value: "30333"
            - name: RPC_PORT
              value: "9933"
            - name: WS_PORT
              value: "9944"
```

Change scripts/config/config.sh accordingly to the previous configured image value.
```sh
IMAGE_OPERATOR=ironoa/customresource-operator:v0.0.8 #define your favourite
# The above parameter have to match with the ones in the deployed resource defined in the deploy/operator.yaml file
```

### Deployment phase

```sh
# start Minikube
$ minikube start

# verify Minikube is running and kubectl is linked to the node
$ kubectl get nodes
NAME       STATUS   ROLES    AGE   VERSION
minikube   Ready    master   6d    v1.17.3

# compile the go project, create the docker image, push the image to the container registry, deploy the k8s resources to the cluster
$ ./scripts/init.sh
serviceaccount/polkadot-operator created
role.rbac.authorization.k8s.io/polkadot-operator created
rolebinding.rbac.authorization.k8s.io/polkadot-operator created
customresourcedefinition.apiextensions.k8s.io/polkadots.polkadot.swisscomblockchain.com created
INFO[0017] Building OCI image ironoa/customresource-operator:v0.0.8
Sending build context to Docker daemon  57.73MB
Step 1/7 : FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
 ---> d17cc1f9d041
Step 2/7 : ENV OPERATOR=/usr/local/bin/polkadot-operator     USER_UID=1001     USER_NAME=polkadot-operator
 ---> Using cache
 ---> 7e017ac07d9a
Step 3/7 : COPY build/_output/bin/polkadot-operator ${OPERATOR}
 ---> 5a7722bb05e4
Step 4/7 : COPY build/bin /usr/local/bin
 ---> 0ae2300af3ce
Step 5/7 : RUN  /usr/local/bin/user_setup
 ---> Running in 7d470167a9de
+ echo 'polkadot-operator:x:1001:0:polkadot-operator user:/root:/sbin/nologin'
+ mkdir -p /root
+ chown 1001:0 /root
+ chmod ug+rwx /root
+ rm /usr/local/bin/user_setup
Removing intermediate container 7d470167a9de
 ---> 85e8d9678a99
Step 6/7 : ENTRYPOINT ["/usr/local/bin/entrypoint"]
 ---> Running in d0dc06133f0f
Removing intermediate container d0dc06133f0f
 ---> d8b7a29276e6
Step 7/7 : USER ${USER_UID}
 ---> Running in 02e4db32c816
Removing intermediate container 02e4db32c816
 ---> 681c495435da
Successfully built 681c495435da
Successfully tagged ironoa/customresource-operator:v0.0.8
INFO[0022] Operator build complete.
The push refers to repository [docker.io/ironoa/customresource-operator]
fe82d8c4fc3f: Pushed
a69c4df39fea: Pushed
1117800c0a97: Pushed
27cd2023d60a: Layer already exists
4b52dfd1f9d9: Layer already exists
v0.0.8: digest: sha256:2bda05786ab1586a844e21bcb5f23b09c588ec7517c64f7652c46c8d97f4f74e size: 1363
deployment.apps/polkadot-operator created
polkadot.polkadot.swisscomblockchain.com/polkadot-cr created

# verify the success of the deployment
$ kubectl get pod
NAME                                READY   STATUS              RESTARTS   AGE
polkadot-operator-78b5fc54f-njv9h   0/1     ContainerCreating   0          2s

# the controller is deployed, after few seconds it will be ready and it will take care of deploying the Polkadot resources automatically
NAME                                READY   STATUS    RESTARTS   AGE
polkadot-operator-78b5fc54f-njv9h   1/1     Running   0          11s
sentry-sset-0                       1/1     Running   0          7s
validator-sset-0                    1/1     Running   0          7s

# verify the single pod behaviour
$ kubectl logs sentry-sset-0
2020-03-03 16:43:09 It isn't safe to expose RPC publicly without a proxy server that filters available set of RPC methods.
2020-03-03 16:43:09 It isn't safe to expose RPC publicly without a proxy server that filters available set of RPC methods.
2020-03-03 16:43:09 Parity Polkadot
2020-03-03 16:43:09   version 0.7.20-3738158-x86_64-linux-gnu
2020-03-03 16:43:09   by Parity Team <admin@parity.io>, 2017-2020
2020-03-03 16:43:09 Chain specification: Kusama CC3
2020-03-03 16:43:09 Node name: IronoaSentry
2020-03-03 16:43:09 Roles: SENTRY
2020-03-03 16:43:09 Native runtime: kusama-1045:2(parity-kusama-0)
2020-03-03 16:43:09 ----------------------------
2020-03-03 16:43:09 This chain is not in any way
2020-03-03 16:43:09       endorsed by the
2020-03-03 16:43:09      KUSAMA FOUNDATION
2020-03-03 16:43:09 ----------------------------
2020-03-03 16:43:09 Initializing Genesis block/state (state: 0xb000…ef6b, header-hash: 0xb0a8…dafe)
2020-03-03 16:43:09 Loading GRANDPA authority set from genesis on what appears to be first startup.

# note how the validator is connected always only with one peer, the sentry
$ kubectl logs validator-sset-0
2020-03-03 16:43:10 It isn't safe to expose RPC publicly without a proxy server that filters available set of RPC methods.
2020-03-03 16:43:10 It isn't safe to expose RPC publicly without a proxy server that filters available set of RPC methods.
2020-03-03 16:43:10 Parity Polkadot
2020-03-03 16:43:10   version 0.7.20-3738158-x86_64-linux-gnu
2020-03-03 16:43:10   by Parity Team <admin@parity.io>, 2017-2020
2020-03-03 16:43:10 Chain specification: Kusama CC3
2020-03-03 16:43:10 Node name: IronoaValidator
2020-03-03 16:43:10 Roles: AUTHORITY
2020-03-03 16:43:10 Native runtime: kusama-1045:2(parity-kusama-0)
2020-03-03 16:43:10 ----------------------------
2020-03-03 16:43:10 This chain is not in any way
2020-03-03 16:43:10       endorsed by the
2020-03-03 16:43:10      KUSAMA FOUNDATION
2020-03-03 16:43:10 ----------------------------
2020-03-03 16:43:10 Initializing Genesis block/state (state: 0xb000…ef6b, header-hash: 0xb0a8…dafe)
2020-03-03 16:43:10 Loading GRANDPA authority set from genesis on what appears to be first startup.
2020-03-03 16:43:11 Loaded block-time = BabeConfiguration { slot_duration: 6000, epoch_length: 600, c: (1, 4), genesis_authorities: [(Public(ca239392960473fe1bc65f94ee27d890a49c1b200c006ff5dcc525330ecc1677 (5Gdk6etL...)), 1), (Public(b46f01874ce7abbb5220e8fd89bede0adad14c73039d91e28e881823433e723f (5G9HTB1d...)), 1), (Public(d684d9176d6eb69887540c9a89fa6097adea82fc4b0ff26d1062b488f352e179 (5GuyZvzU...)), 1), (Public(68195a71bdde49117a616424bdc60a1733e96acb1da5aeab5d268cf2a572e941 (5ERCNLU4...)), 1), (Public(1a0575ef4ae24bdfd31f4cb5bd61239ae67c12d4e64ae51ac756044aa6ad8200 (5Cepixt1...)), 1), (Public(18168f2aad0081a25728961ee00627cfe35e39833c805016632bf7c14da58009 (5CcHi1bG...)), 1)], randomness: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0], secondary_slots: true } milliseconds from genesis on first-launch
2020-03-03 16:43:11 Creating empty BABE epoch changes on what appears to be first startup.
2020-03-03 16:43:11 Highest known block at #0
2020-03-03 16:43:11 Local node identity is: QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM
2020-03-03 16:43:11 Starting BABE Authorship worker
2020-03-03 16:43:11 Grafana data source server started at 127.0.0.1:9955
2020-03-03 16:43:11 Discovered new external address for our node: /ip4/172.17.0.6/tcp/30333/p2p/QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM
2020-03-03 16:43:11 New epoch 0 launching at block 0xcd9b…5ad3 (block slot 262493679 >= start slot 262493679).
2020-03-03 16:43:11 Next epoch starts at slot 262494279
2020-03-03 16:43:11 Discovered new external address for our node: /ip4/10.0.1.134/tcp/30333/p2p/QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM
2020-03-03 16:43:11 Reserved peer QmQMTLWkNwGf7P5MQv7kUHCynMg7jje6h3vbvwd2ALPPhm disconnected
2020-03-03 16:43:12 Discovered new external address for our node: /dns4/validator-service/tcp/30333/p2p/QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM
2020-03-03 16:43:13 Discovered new external address for our node: /ip4/178.197.224.81/tcp/30333/p2p/QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM
2020-03-03 16:43:16 Idle (0 peers), best: #18 (0x48b2…00e2), finalized #0 (0xb0a8…dafe), ⬇ 6.4kiB/s ⬆ 3.3kiB/s
2020-03-03 16:43:21 Idle (1 peers), best: #58 (0x29a8…dfc2), finalized #0 (0xb0a8…dafe), ⬇ 7.4kiB/s ⬆ 1.2kiB/s
2020-03-03 16:43:26 Idle (1 peers), best: #134 (0x55c2…7604), finalized #0 (0xb0a8…dafe), ⬇ 2.6kiB/s ⬆ 1.0kiB/s

# clean up the resources
$ ./scripts/wipeAll.sh
deployment.apps "polkadot-operator" deleted
polkadot.polkadot.swisscomblockchain.com "polkadot-cr" deleted
customresourcedefinition.apiextensions.k8s.io "polkadots.polkadot.swisscomblockchain.com" deleted
rolebinding.rbac.authorization.k8s.io "polkadot-operator" deleted
role.rbac.authorization.k8s.io "polkadot-operator" deleted
serviceaccount "polkadot-operator" deleted

# stop Minikube
$ minikube stop
```

## Operator Configurable Environment Variables

* IMAGE_CLIENT: (string)  
Client Image on the Container Registry.

* METRICS_PORT: (string)  
Port of the service where it is possible to scrape the metrics from. See the [Metrics Support section](#metrics-support).

* P2P_PORT: (string)  
P2P port of both the service and the client.

* RPC_PORT: (string)  
RPC port of both the service and the client.

* WS_PORT: (string)  
Web Socket port of both the service and the client.

## Polkadot CR Configurable Parameters

* clientVersion: (string)  
Image version of the clients. See the [Updating of Node Versions section](#updating-of-node-versions).

* secureCommunicationSupport: (struct)
    * enabled: (bool)    
If set to "true", the operator will handle the creation and the deployment of a Network Policy object that will ensure the secureness of the Validator (it only affects the Kind "SentryAndValidator"). 
With the parameter active, the Validator is allowed to communicate only with the Sentry layer. Being this mechanism enforced via NetworkPolicy (kubernetes native object), it requires a network plugin installed in you cloud provided cluster (even in minikube) to work properly.  
See the [Secure Communications section](#secure-communications-kindsentryandvalidator).

* metricsSupport: (struct)
    * enabled: (bool)  
See the [Metrics support section](#metrics-support).    

* replicas: (int)  
Allows to decide how many Sentry replicas will be created. See the [Node Cluster Scaling Support section](#node-cluster-scaling-support).

* clientName: (string)

* resources: (ResourceRequirements)  
You can limit or require a specif amount of CPU or memory from your cluster, for example.  
See the official godoc: https://godoc.org/k8s.io/api/core/v1#ResourceRequirements

* nodeKey: (string)  
Identity of the node, private (e.g. "0000000000000000000000000000000000000000000000000000000000000013")

* dataPersistenceSupport: (struct)
    * enabled: (bool)
    * persistentVolumeClaim: (PersistentVolumeClaim)  
    See the official godoc: https://godoc.org/k8s.io/api/core/v1#PersistentVolumeClaim  
        * storageClassName: (string)  
        This field is particularly interesting. For instance, Azure provides two built in storage classes:  
            * "default": HHD backed
            * "managed-premium": SSD backed, high performance  
        See [Data Persistence Support section](#data-persistence-support) for more information.     

* kind: Sentry | Validator | SentryAndValidator (string)  
Desired deployable configuration:
    * Sentry: deploy a Sentry only configuration
    * Validator: deploy a Validator only configuration
    * SentryAndValidator: deploy a Sentry and Validator configuration (please take a look at the Secure Communications section). In the SentryAndValidator configuration it must be passed an additional parameter to both the sentry and the validator:
        * reservedValidatorID: (string) Identity of the Validator, it must be set on the Sentry
        * reservedSentryID: (string) Identity of the Sentry, it must be set on the Validator
        
            ![alt text](images/schema.png)

## Updating of Node Versions

It is possible to change the Client Nodes Version at runtime (kubectl apply): the operator will automatically handle the clients version update of all the running pods.  

```yaml
# You can control the feature customizing this parameter on the yourCRfile.yaml
spec:
  clientVersion: latest
```

```sh
# You can apply the change with this command, the operator will take care of the rest for you
$ kubectl apply -f yourCRfile.yaml
```

## Node Cluster Scaling Support

This is the ability of the operator to respond to scale operations defined in the deployed configuration, for example to extend the amount of sentry nodes from 3 to 4. The correct functioning can be tested by executing such an operation and checking the number of deployed instances before and afterwards.  
In any case, Validator replica size is always hard coded to one and it is not possible to change it to prevent concurrent validation issues.

```yaml
# You can control the feature customizing this parameter on the yourCRfile.yaml
sentry:
  replicas: 2
```

```sh
# You can apply the change with this command, the operator will take care of the rest for you
$ kubectl apply -f yourCRfile.yaml
```

### Please Note

* If you are running your Polkadot Custom Resources with the Data Persistence Feature enabled, each one of your replicas will need a dedicated Persistent Volume available. Please make sure that there are enough of those, one for each instance.
* you can even use the utils script to apply your updated Custom Resource: scripts/utils/updateCR.yaml

            
## Secure Communications (Kind:SentryAndValidator)

The configuration is based on the "polkadot-secure-validator" guidelines: https://github.com/w3f/polkadot-secure-validator  

The implementation uses the NetworkPolicy kubernetes object.

### Network Policies

By default, pods are non-isolated; they accept traffic from any source. Pods become isolated by having a NetworkPolicy that selects them. A network policy is a specification of how groups of pods are allowed to communicate with each other and other network endpoints.
Reference: https://kubernetes.io/docs/concepts/services-networking/network-policies/

### Prerequisites

Network policies are implemented by the network plugin. To use network policies, you must be using a networking solution which supports NetworkPolicy. Creating a NetworkPolicy resource without a controller that implements it will have no effect.

### Azure Example

A tested working solution is using "Calico Network Policies" as network plugin on an Azure Kubernetes Service. 
Reference: https://docs.microsoft.com/en-us/azure/aks/use-network-policies

You can test the effectiveness of the network policy creating a new "default deny" one for the validator: it will not be able to communicate with the sentry (and even whit the external world) anymore. 

## Data Persistence Support

Deployments on Kubernetes are by their nature ephemeral. Thus it is important to  provide Kubernetes with support for data persistence – such as a virtual SSD in the cloud – so that new instances of the application can resume the state of the previous instance. It can be tested by killing a Stateful Set instance and then checking whether the state (block number synchronization) is resumed by the new instance.  

A tested working solution is using a Kubernetes Persistent Volume Claim with Azure Built in storage classes (Azure Disk).
Reference: https://docs.microsoft.com/en-us/azure/aks/azure-disks-dynamic-pv

You can even decide to deploy your own StorageClass according to the specs provided by your favourite Cluster Provider, making this solution cluster agnostic.  
Reference: https://kubernetes.io/docs/concepts/storage/storage-classes/

### How To Tutorial with Minikube

If you want to test it locally, you first have to manually provide a few persistent volumes (at least two, one for each client you deploy) to minikube. Minikube will extract from this named pool (storageClassName) an available volume thanks to the Persistent Volume Claim mechanism.   
Please note, in this example we decided to name the storageClassName as "local".

```yaml
# Copyright (c) 2020 Swisscom Blockchain AG
# Licensed under MIT License
# persistentVolume.yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-volume0
  labels:
    type: local
spec:
  storageClassName: local
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: "/mnt/vda1/data0"
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-volume1
  labels:
    type: local
spec:
  storageClassName: local
  capacity:
    storage: 10Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: "/mnt/vda1/data1"
```

The fields "hostPath -> path" describe where the folder will be mounted inside minikube. We have first to create that folders and grant some permission:

```sh
# ssh into Minikube
$ minikube ssh
                         _             _
            _         _ ( )           ( )
  ___ ___  (_)  ___  (_)| |/')  _   _ | |_      __
/' _ ` _ `\| |/' _ `\| || , <  ( ) ( )| '_`\  /'__`\
| ( ) ( ) || || ( ) || || |\`\ | (_) || |_) )(  ___/
(_) (_) (_)(_)(_) (_)(_)(_) (_)`\___/'(_,__/'`\____)

$ cd /mnt/vda1/
$ sudo mkdir data0
$ sudo mkdir data1
$ sudo chmod 0777 data0
$ sudo chmod 0777 data1
# Done it! now exit from the ssh 
```

Now we can create the two Persistent Volume we defined previously:

```sh
$ kubectl apply -f persistentVolume.yaml
persistentvolume/pv-volume0 created
persistentvolume/pv-volume1 created
```

Finally, change the storageClassName parameter in the CustomResource definition, for example:

```yaml
# Copyright (c) 2020 Swisscom Blockchain AG
# Licensed under MIT License
apiVersion: polkadot.swisscomblockchain.com/v1alpha1
kind: Polkadot
metadata:
  name: polkadot-cr
spec:
  clientVersion: latest
  kind: "SentryAndValidator"
  secureCommunicationSupport:
    enabled: false
  metricsSupport:
    enabled: false
  sentry:
    replicas: 1
    clientName: "IronoaSentry"
    nodeKey: "0000000000000000000000000000000000000000000000000000000000000013" # Local node id: QmQMTLWkNwGf7P5MQv7kUHCynMg7jje6h3vbvwd2ALPPhm
    reservedValidatorID: "QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM"
    resources:
      limits:
        memory: "512Mi"
        cpu: "0.5"
    dataPersistenceSupport:
      enabled: true
      persistentVolumeClaim:
        metadata:
          name: "polkadot-volume"
        spec:
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 10Gi
          storageClassName: local 
  validator:
    clientName: "IronoaValidator"
    nodeKey: "0000000000000000000000000000000000000000000000000000000000000021" # Local node id: QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM
    reservedSentryID: "QmQMTLWkNwGf7P5MQv7kUHCynMg7jje6h3vbvwd2ALPPhm"
    resources:
      limits:
        memory: "512Mi"
        cpu: "0.5"
    dataPersistenceSupport:
      enabled: true
      persistentVolumeClaim:
        metadata:
          name: "polkadot-volume"
        spec:
          accessModes:
            - ReadWriteOnce
          resources:
            requests:
              storage: 10Gi
          storageClassName: local 
```

You can now deploy the operator as usual, also with the init.sh script.

## Metrics Support

Substrate exposes an endpoint which serves metrics in the Prometheus exposition format available on port 9615. You can change the port with --prometheus-port <PORT> and enable it to be accessed over an interface other than local host with --prometheus-external.  
Official documentation at: https://substrate.dev/docs/en/next/tutorials/visualizing-node-metrics/

### Default configuration

* Metrics functionality is not active by default, you have to explicitly activate it by setting the parameter metricsSupport->enabled to "true"
* Each Pod Service provide access to the metrics:
    * at port 9615
    * at /metrics endpoint
    * "client-service-ip:9615/metrics"
    
You can change the metrics port via the parameter METRICS_PORT in the deploy/operator.yaml

### How to access to the metrics: Example in Minikube

```sh
# accordingly to the previuous tutorials, deploy the operator via the init script
$ ./init.sh
serviceaccount/polkadot-operator created
role.rbac.authorization.k8s.io/polkadot-operator created
rolebinding.rbac.authorization.k8s.io/polkadot-operator created
customresourcedefinition.apiextensions.k8s.io/polkadots.polkadot.swisscomblockchain.com create
INFO[0004] Building OCI image ironoa/customresource-operator:v0.0.8
Sending build context to Docker daemon  58.03MB
Step 1/7 : FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
 ---> d17cc1f9d041
Step 2/7 : ENV OPERATOR=/usr/local/bin/polkadot-operator     USER_UID=1001     USER_NAME=polkadot-operator
 ---> Using cache
 ---> 7e017ac07d9a
Step 3/7 : COPY build/_output/bin/polkadot-operator ${OPERATOR}
 ---> 27824852294b
Step 4/7 : COPY build/bin /usr/local/bin
 ---> 9c402e3c4e3a
Step 5/7 : RUN  /usr/local/bin/user_setup
 ---> Running in 83b229aa3959
+ echo 'polkadot-operator:x:1001:0:polkadot-operator user:/root:/sbin/nologin'
+ mkdir -p /root
+ chown 1001:0 /root
+ chmod ug+rwx /root
+ rm /usr/local/bin/user_setup
Removing intermediate container 83b229aa3959
 ---> 867af0f89f07
Step 6/7 : ENTRYPOINT ["/usr/local/bin/entrypoint"]
 ---> Running in a31e1d3c6fe9
Removing intermediate container a31e1d3c6fe9
 ---> c9584177aefe
Step 7/7 : USER ${USER_UID}
 ---> Running in 4c65e80f42cf
Removing intermediate container 4c65e80f42cf
 ---> 990d538921ab
Successfully built 990d538921ab
Successfully tagged ironoa/customresource-operator:v0.0.8
INFO[0010] Operator build complete.
The push refers to repository [docker.io/ironoa/customresource-operator]
265a1cc36a93: Pushed
926eb1253b1e: Pushed
352257d20965: Pushed
27cd2023d60a: Layer already exists
4b52dfd1f9d9: Layer already exists
v0.0.8: digest: sha256:3c21d9ac7cf5d0cf1e0ecf88d0a69d3929eb5754912733fe560090f8fc582354 size: 1363
deployment.apps/polkadot-operator created
polkadot.polkadot.swisscomblockchain.com/polkadot-cr created

# check the status of the deployment
$ kubectl get pods
polkadot-operator-78b5fc54f-v9d6d   1/1     Running   0          33s
sentry-sset-0                       1/1     Running   0          29s
validator-sset-0                    1/1     Running   0          29s

# retrieve the services and check the IP addresses of the polkadot clients
$ kubectl get services
kubernetes                  ClusterIP   10.96.0.1        <none>        443/TCP                                                        77m
polkadot-operator-metrics   ClusterIP   10.100.143.145   <none>        8383/TCP,8686/TCP                                              87s
sentry-service              NodePort    10.96.76.21      <none>        30333:31945/TCP,9933:30506/TCP,9944:30586/TCP,9615:31836/TCP   87s
validator-service           ClusterIP   10.101.249.247   <none>        30333/TCP,9933/TCP,9944/TCP,9615/TCP                           87s

# access inside the minikube cluster
$ minikube ssh
                         _             _
            _         _ ( )           ( )
  ___ ___  (_)  ___  (_)| |/')  _   _ | |_      __
/' _ ` _ `\| |/' _ `\| || , <  ( ) ( )| '_`\  /'__`\
| ( ) ( ) || || ( ) || || |\`\ | (_) || |_) )(  ___/
(_) (_) (_)(_)(_) (_)(_)(_) (_)`\___/'(_,__/'`\____)

# test the metrics enpoints
$ curl 10.96.76.21:9615/metrics

$ curl 10.101.249.247:9615/metrics
# HELP polkadot_block_height_number Height of the chain
# TYPE polkadot_block_height_number gauge
polkadot_block_height_number{status="best"} 1196
polkadot_block_height_number{status="finalized"} 1024
polkadot_block_height_number{status="sync_target"} 1667996
# HELP polkadot_build_info A metric with a constant '1' value labeled by name, version, and commit.
# TYPE polkadot_build_info gauge
polkadot_build_info{commit="7f59f2c",name="parity-polkadot",version="0.7.28"} 1
# HELP polkadot_cpu_usage_percentage Node CPU usage
# TYPE polkadot_cpu_usage_percentage gauge
polkadot_cpu_usage_percentage 51.579490661621094
# HELP polkadot_database_cache_bytes RocksDB cache size in bytes
# TYPE polkadot_database_cache_bytes gauge
polkadot_database_cache_bytes 13318400
# HELP polkadot_finality_grandpa_communication_gossip_validator_messages Number of messages validated by the finality grandpa gossip validator.
# TYPE polkadot_finality_grandpa_communication_gossip_validator_messages counter
polkadot_finality_grandpa_communication_gossip_validator_messages{action="discard",message="neighbor"} 519
# HELP polkadot_finality_grandpa_round Highest completed GRANDPA round.
# TYPE polkadot_finality_grandpa_round gauge
polkadot_finality_grandpa_round 0
# HELP polkadot_import_queue_blocks_submitted Number of blocks submitted to the import queue.
# TYPE polkadot_import_queue_blocks_submitted counter
polkadot_import_queue_blocks_submitted 40
# HELP polkadot_import_queue_finality_proofs_submitted Number of finality proofs submitted to the import queue.
# TYPE polkadot_import_queue_finality_proofs_submitted counter
polkadot_import_queue_finality_proofs_submitted 0
# HELP polkadot_import_queue_justifications_submitted Number of justifications submitted to the import queue.
# TYPE polkadot_import_queue_justifications_submitted counter
polkadot_import_queue_justifications_submitted 0
# HELP polkadot_memory_usage_bytes Node memory usage
# TYPE polkadot_memory_usage_bytes gauge
polkadot_memory_usage_bytes 105020
# HELP polkadot_network_per_sec_bytes Networking bytes per second
# TYPE polkadot_network_per_sec_bytes gauge
polkadot_network_per_sec_bytes{direction="download"} 4252
polkadot_network_per_sec_bytes{direction="upload"} 2885
# HELP polkadot_node_roles The roles the node is running as
# TYPE polkadot_node_roles gauge
polkadot_node_roles 4
# HELP polkadot_ready_transactions_number Number of transactions in the ready queue
# TYPE polkadot_ready_transactions_number gauge
polkadot_ready_transactions_number 0
# HELP polkadot_state_cache_bytes State cache size in bytes
# TYPE polkadot_state_cache_bytes gauge
polkadot_state_cache_bytes 1253522
# HELP polkadot_state_db_cache_bytes State DB cache in bytes
# TYPE polkadot_state_db_cache_bytes gauge
polkadot_state_db_cache_bytes{subtype="non_canonical"} 168
polkadot_state_db_cache_bytes{subtype="pinned"} 44
# HELP polkadot_sub_libp2p_connections Number of libp2p connections
# TYPE polkadot_sub_libp2p_connections gauge
polkadot_sub_libp2p_connections 100
# HELP polkadot_sub_libp2p_is_major_syncing Whether the node is performing a major sync or not.
# TYPE polkadot_sub_libp2p_is_major_syncing gauge
polkadot_sub_libp2p_is_major_syncing 1
# HELP polkadot_sub_libp2p_kbuckets_num_nodes Number of nodes in the Kademlia k-buckets
# TYPE polkadot_sub_libp2p_kbuckets_num_nodes gauge
polkadot_sub_libp2p_kbuckets_num_nodes 76
# HELP polkadot_sub_libp2p_network_per_sec_bytes Average bandwidth usage per second
# TYPE polkadot_sub_libp2p_network_per_sec_bytes gauge
polkadot_sub_libp2p_network_per_sec_bytes{direction="in"} 4406
polkadot_sub_libp2p_network_per_sec_bytes{direction="out"} 2864
# HELP polkadot_sub_libp2p_notifications_total Number of notification received from all nodes
# TYPE polkadot_sub_libp2p_notifications_total counter
polkadot_sub_libp2p_notifications_total{direction="in",protocol="FRNK"} 519
polkadot_sub_libp2p_notifications_total{direction="in",protocol="dot2"} 20
polkadot_sub_libp2p_notifications_total{direction="out",protocol="FRNK"} 22
polkadot_sub_libp2p_notifications_total{direction="out",protocol="dot2"} 22
# HELP polkadot_sub_libp2p_num_event_stream_channels Number of internal active channels that broadcast network events
# TYPE polkadot_sub_libp2p_num_event_stream_channels gauge
polkadot_sub_libp2p_num_event_stream_channels 3
# HELP polkadot_sub_libp2p_opened_notification_streams Number of open notification substreams
# TYPE polkadot_sub_libp2p_opened_notification_streams gauge
polkadot_sub_libp2p_opened_notification_streams{protocol="FRNK"} 22
polkadot_sub_libp2p_opened_notification_streams{protocol="dot1"} 22
polkadot_sub_libp2p_opened_notification_streams{protocol="dot2"} 22
# HELP polkadot_sub_libp2p_peers_count Number of network gossip peers
# TYPE polkadot_sub_libp2p_peers_count gauge
polkadot_sub_libp2p_peers_count 22
# HELP polkadot_sub_libp2p_peerset_num_discovered Number of nodes stored in the peerset manager
# TYPE polkadot_sub_libp2p_peerset_num_discovered gauge
polkadot_sub_libp2p_peerset_num_discovered 262
# HELP polkadot_sub_libp2p_peerset_num_requested Number of nodes that the peerset manager wants us to be connected to
# TYPE polkadot_sub_libp2p_peerset_num_requested gauge
polkadot_sub_libp2p_peerset_num_requested 26
# HELP polkadot_sub_libp2p_random_kademalia_queries_total Number of random Kademlia queries started
# TYPE polkadot_sub_libp2p_random_kademalia_queries_total counter
polkadot_sub_libp2p_random_kademalia_queries_total 5
# HELP polkadot_sync_extra_finality_proofs Number of extra finality proof requests
# TYPE polkadot_sync_extra_finality_proofs gauge
polkadot_sync_extra_finality_proofs{status="active"} 0
polkadot_sync_extra_finality_proofs{status="failed"} 0
polkadot_sync_extra_finality_proofs{status="importing"} 0
polkadot_sync_extra_finality_proofs{status="pending"} 0
# HELP polkadot_sync_extra_justifications Number of extra justifications requests
# TYPE polkadot_sync_extra_justifications gauge
polkadot_sync_extra_justifications{status="active"} 0
polkadot_sync_extra_justifications{status="failed"} 0
polkadot_sync_extra_justifications{status="importing"} 0
polkadot_sync_extra_justifications{status="pending"} 0
# HELP polkadot_sync_fork_targets Number of fork sync targets
# TYPE polkadot_sync_fork_targets gauge
polkadot_sync_fork_targets 0
# HELP polkadot_sync_handshaking_peers Number of newly connected peers
# TYPE polkadot_sync_handshaking_peers gauge
polkadot_sync_handshaking_peers 0
# HELP polkadot_sync_obsolete_requests Number of obsolete requests
# TYPE polkadot_sync_obsolete_requests gauge
polkadot_sync_obsolete_requests 0
# HELP polkadot_sync_peers Number of peers we sync with
# TYPE polkadot_sync_peers gauge
polkadot_sync_peers 22
# HELP polkadot_sync_queued_blocks Number of blocks in import queue
# TYPE polkadot_sync_queued_blocks gauge
polkadot_sync_queued_blocks 2304
```

## E2E Testing

End-to-end (e2e) testing is automated testing written as Go test.   
All the e2e tests are located in test/e2e.  
The Operator SDK includes a testing framework to make writing tests simpler and quicker by removing boilerplate code and providing common test utilities.  
Reference: https://github.com/operator-framework/operator-sdk/blob/master/doc/test-framework/writing-e2e-tests.md

### Build and run test

Requirements:

* a running k8s cluster (e.g. Minikube)

You can run the tests in your local environment with the following command:
```sh
# from the root of the project...
$ operator-sdk test local ./test/e2e --go-test-flags "-v"
INFO[0000] Testing operator locally.                    
=== RUN   TestPolkadot
    TestPolkadot: client.go:62: resource type ServiceAccount with namespace/name (osdk-e2e-6112de68-81ad-43f4-85d9-a38e3834381c/polkadot-operator) created
    TestPolkadot: client.go:62: resource type Role with namespace/name (osdk-e2e-6112de68-81ad-43f4-85d9-a38e3834381c/polkadot-operator) created
    TestPolkadot: client.go:62: resource type RoleBinding with namespace/name (osdk-e2e-6112de68-81ad-43f4-85d9-a38e3834381c/polkadot-operator) created
    TestPolkadot: client.go:62: resource type Deployment with namespace/name (osdk-e2e-6112de68-81ad-43f4-85d9-a38e3834381c/polkadot-operator) created
    TestPolkadot: wait_util.go:70: Deployment available (1/1)
=== RUN   TestPolkadot/TestPolkadotSentry
    TestPolkadot: client.go:62: resource type  with namespace/name (osdk-e2e-6112de68-81ad-43f4-85d9-a38e3834381c/example-polkadot) created
=== RUN   TestPolkadot/TestPolkadotSentry/TestStatefulSetCreation
    TestPolkadot/TestPolkadotSentry/TestStatefulSetCreation: waitUtils.go:26: Waiting for full availability of sentry-sset stateful set (0/1)
    TestPolkadot/TestPolkadotSentry/TestStatefulSetCreation: waitUtils.go:26: Waiting for full availability of sentry-sset stateful set (0/1)
    TestPolkadot/TestPolkadotSentry/TestStatefulSetCreation: waitUtils.go:32: Stateful Set available (1/1)
=== RUN   TestPolkadot/TestPolkadotSentry/TestServiceCreation
    TestPolkadot/TestPolkadotSentry/TestServiceCreation: waitUtils.go:74: Service available
=== RUN   TestPolkadot/TestPolkadotSentry/TestStatefulSetDeletion
    TestPolkadot/TestPolkadotSentry/TestStatefulSetDeletion: waitUtils.go:55: Stateful Set deleted
=== RUN   TestPolkadot/TestPolkadotSentry/TestServiceDeletion
    TestPolkadot/TestPolkadotSentry/TestServiceDeletion: waitUtils.go:93: Service deleted
=== RUN   TestPolkadot/TestPolkadotValidator
    TestPolkadot: client.go:62: resource type  with namespace/name (osdk-e2e-6112de68-81ad-43f4-85d9-a38e3834381c/example-polkadot) created
=== RUN   TestPolkadot/TestPolkadotValidator/TestStatefulSetCreation
    TestPolkadot/TestPolkadotValidator/TestStatefulSetCreation: waitUtils.go:26: Waiting for full availability of validator-sset stateful set (0/1)
    TestPolkadot/TestPolkadotValidator/TestStatefulSetCreation: waitUtils.go:26: Waiting for full availability of validator-sset stateful set (0/1)
    TestPolkadot/TestPolkadotValidator/TestStatefulSetCreation: waitUtils.go:32: Stateful Set available (1/1)
=== RUN   TestPolkadot/TestPolkadotValidator/TestServiceCreation
    TestPolkadot/TestPolkadotValidator/TestServiceCreation: waitUtils.go:74: Service available
=== RUN   TestPolkadot/TestPolkadotValidator/TestStatefulSetDeletion
    TestPolkadot/TestPolkadotValidator/TestStatefulSetDeletion: waitUtils.go:55: Stateful Set deleted
=== RUN   TestPolkadot/TestPolkadotValidator/TestServiceDeletion
    TestPolkadot/TestPolkadotValidator/TestServiceDeletion: waitUtils.go:93: Service deleted
=== RUN   TestPolkadot/TestPolkadotSentryAndValidator
    TestPolkadot: client.go:62: resource type  with namespace/name (osdk-e2e-6112de68-81ad-43f4-85d9-a38e3834381c/example-polkadot) created
=== RUN   TestPolkadot/TestPolkadotSentryAndValidator/TestStatefulSetCreationSentry
    TestPolkadot/TestPolkadotSentryAndValidator/TestStatefulSetCreationSentry: waitUtils.go:26: Waiting for full availability of sentry-sset stateful set (0/1)
    TestPolkadot/TestPolkadotSentryAndValidator/TestStatefulSetCreationSentry: waitUtils.go:26: Waiting for full availability of sentry-sset stateful set (0/1)
    TestPolkadot/TestPolkadotSentryAndValidator/TestStatefulSetCreationSentry: waitUtils.go:26: Waiting for full availability of sentry-sset stateful set (0/1)
    TestPolkadot/TestPolkadotSentryAndValidator/TestStatefulSetCreationSentry: waitUtils.go:26: Waiting for full availability of sentry-sset stateful set (0/1)
    TestPolkadot/TestPolkadotSentryAndValidator/TestStatefulSetCreationSentry: waitUtils.go:32: Stateful Set available (1/1)
=== RUN   TestPolkadot/TestPolkadotSentryAndValidator/TestServiceCreationSentry
    TestPolkadot/TestPolkadotSentryAndValidator/TestServiceCreationSentry: waitUtils.go:74: Service available
=== RUN   TestPolkadot/TestPolkadotSentryAndValidator/TestStatefulSetCreationValidator
    TestPolkadot/TestPolkadotSentryAndValidator/TestStatefulSetCreationValidator: waitUtils.go:32: Stateful Set available (1/1)
=== RUN   TestPolkadot/TestPolkadotSentryAndValidator/TestServiceCreationValidator
    TestPolkadot/TestPolkadotSentryAndValidator/TestServiceCreationValidator: waitUtils.go:74: Service available
=== RUN   TestPolkadot/TestPolkadotSentryAndValidator/TestNetworkPolicyCreation
    TestPolkadot/TestPolkadotSentryAndValidator/TestNetworkPolicyCreation: waitUtils.go:112: Network policy available
=== RUN   TestPolkadot/TestPolkadotSentryAndValidator/TestStatefulSetDeletionSentry
    TestPolkadot/TestPolkadotSentryAndValidator/TestStatefulSetDeletionSentry: waitUtils.go:55: Stateful Set deleted
=== RUN   TestPolkadot/TestPolkadotSentryAndValidator/TestServiceDeletionSentry
    TestPolkadot/TestPolkadotSentryAndValidator/TestServiceDeletionSentry: waitUtils.go:93: Service deleted
=== RUN   TestPolkadot/TestPolkadotSentryAndValidator/TestStatefulSetDeletionValidator
    TestPolkadot/TestPolkadotSentryAndValidator/TestStatefulSetDeletionValidator: waitUtils.go:55: Stateful Set deleted
=== RUN   TestPolkadot/TestPolkadotSentryAndValidator/TestServiceDeletionValidator
    TestPolkadot/TestPolkadotSentryAndValidator/TestServiceDeletionValidator: waitUtils.go:93: Service deleted
=== RUN   TestPolkadot/TestPolkadotSentryAndValidator/TestNetworkPolicyDeletion
    TestPolkadot/TestPolkadotSentryAndValidator/TestNetworkPolicyDeletion: waitUtils.go:131: Network policy deleted
    TestPolkadot: client.go:82: resource type  with namespace/name (osdk-e2e-6112de68-81ad-43f4-85d9-a38e3834381c/example-polkadot) successfully deleted
    TestPolkadot: client.go:82: resource type  with namespace/name (osdk-e2e-6112de68-81ad-43f4-85d9-a38e3834381c/example-polkadot) successfully deleted
    TestPolkadot: client.go:82: resource type  with namespace/name (osdk-e2e-6112de68-81ad-43f4-85d9-a38e3834381c/example-polkadot) successfully deleted
    TestPolkadot: client.go:82: resource type Deployment with namespace/name (osdk-e2e-6112de68-81ad-43f4-85d9-a38e3834381c/polkadot-operator) successfully deleted
    TestPolkadot: client.go:82: resource type RoleBinding with namespace/name (osdk-e2e-6112de68-81ad-43f4-85d9-a38e3834381c/polkadot-operator) successfully deleted
    TestPolkadot: client.go:82: resource type Role with namespace/name (osdk-e2e-6112de68-81ad-43f4-85d9-a38e3834381c/polkadot-operator) successfully deleted
    TestPolkadot: client.go:82: resource type ServiceAccount with namespace/name (osdk-e2e-6112de68-81ad-43f4-85d9-a38e3834381c/polkadot-operator) successfully deleted
--- PASS: TestPolkadot (136.37s)
    --- PASS: TestPolkadot/TestPolkadotSentry (30.05s)
        --- PASS: TestPolkadot/TestPolkadotSentry/TestStatefulSetCreation (15.01s)
        --- PASS: TestPolkadot/TestPolkadotSentry/TestServiceCreation (5.01s)
        --- PASS: TestPolkadot/TestPolkadotSentry/TestStatefulSetDeletion (5.01s)
        --- PASS: TestPolkadot/TestPolkadotSentry/TestServiceDeletion (5.01s)
    --- PASS: TestPolkadot/TestPolkadotValidator (30.05s)
        --- PASS: TestPolkadot/TestPolkadotValidator/TestStatefulSetCreation (15.01s)
        --- PASS: TestPolkadot/TestPolkadotValidator/TestServiceCreation (5.01s)
        --- PASS: TestPolkadot/TestPolkadotValidator/TestStatefulSetDeletion (5.01s)
        --- PASS: TestPolkadot/TestPolkadotValidator/TestServiceDeletion (5.01s)
    --- PASS: TestPolkadot/TestPolkadotSentryAndValidator (70.09s)
        --- PASS: TestPolkadot/TestPolkadotSentryAndValidator/TestStatefulSetCreationSentry (25.01s)
        --- PASS: TestPolkadot/TestPolkadotSentryAndValidator/TestServiceCreationSentry (5.01s)
        --- PASS: TestPolkadot/TestPolkadotSentryAndValidator/TestStatefulSetCreationValidator (5.01s)
        --- PASS: TestPolkadot/TestPolkadotSentryAndValidator/TestServiceCreationValidator (5.01s)
        --- PASS: TestPolkadot/TestPolkadotSentryAndValidator/TestNetworkPolicyCreation (5.01s)
        --- PASS: TestPolkadot/TestPolkadotSentryAndValidator/TestStatefulSetDeletionSentry (5.00s)
        --- PASS: TestPolkadot/TestPolkadotSentryAndValidator/TestServiceDeletionSentry (5.01s)
        --- PASS: TestPolkadot/TestPolkadotSentryAndValidator/TestStatefulSetDeletionValidator (5.01s)
        --- PASS: TestPolkadot/TestPolkadotSentryAndValidator/TestServiceDeletionValidator (5.01s)
        --- PASS: TestPolkadot/TestPolkadotSentryAndValidator/TestNetworkPolicyDeletion (5.01s)
PASS
ok      github.com/swisscom-blockchain/polkadot-k8s-operator/test/e2e   136.908s
?       github.com/swisscom-blockchain/polkadot-k8s-operator/test/e2e/utils     [no test files]
```

## About Kubernetes

Kubernetes (K8s) is an open-source system for automating deployment, scaling, and management of containerized applications.  

Reference: https://kubernetes.io/

## About the Operator

An Operator is a method of packaging, deploying and managing a Kubernetes application. A Kubernetes application is an application that is both deployed on Kubernetes and managed using the Kubernetes APIs and kubectl tooling. You can think of Operators as the runtime that manages this type of application on Kubernetes.  

Reference: https://coreos.com/blog/introducing-operator-framework

## About the Operator SDK and Framework

To help make it easier to build Kubernetes applications, Red Hat and the Kubernetes open source community today share the Operator Framework – an open source toolkit designed to manage Kubernetes native applications, called Operators, in a more effective, automated, and scalable way. 

The Operator SDK provides the tools to build, test and package Operators. Initially, the SDK facilitates the marriage of an application’s business logic (for example, how to scale, upgrade, or backup) with the Kubernetes API to execute those operations. 

Reference: https://coreos.com/blog/introducing-operator-framework

## Project Directory Structure

The directory structure is based on the one generated by the operator-sdk CLI.  
See the table at: https://github.com/operator-framework/operator-sdk/blob/master/doc/project_layout.md

The most interesting part, the controller business logic, is located in pkg/controller/polkadot/polkadotController.go  

The deployable resources, the CR and the CRD, are located under deploy/crds/  
Other deployable resources such as the controller operator and the service account are located under deploy/  

The scripts to easily compile and deploy the operator are located under scripts/

The end-to-end tests are located under tests/e2e/
