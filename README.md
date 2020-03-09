Copyright (c) 2020 Swisscom Blockchain AG  
Licensed under MIT License

# PolkaKop - Kubernetes Operator for Polkadot 

Kubernetes Operator for Polkadot Sentry and Validators nodes.

Client - Rust implementation of the Polkadot Host: https://github.com/paritytech/polkadot

The Polkadot Project: https://wiki.polkadot.network/en/

## Polkadot Custom Resource 

The deployable CR (Custom Resource) is called "Polkadot"

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

* Kubernetes Cluster Network Plugin: network plugin, see the Secure Communications section (SentryAndValidator secure deployment)

* Minikube, a tool that runs a single-node Kubernetes cluster in your local environment   
Mac: https://kubernetes.io/docs/tasks/tools/install-minikube/#install-minikube  
Linux: https://kubernetes.io/docs/tasks/tools/install-minikube/#install-minikube-using-a-package

## How To Run

Deploy to your favorite kubernetes cloud provided cluster (even minikube) a Custom Controller and a Polkadot Custom Resource. The Controller will create and supervise all the necessary resources needed to run a Polkadot Client configuration.

0. Configure your kubectl to work with your desired Kubernetes cluster 
    (e.g. Azure: az aks get-credentials --resource-group myResourceGroup --name myAKSCluster)
1. Clone the repository locally
2. In both deploy/operator.yaml and scripts/utils/compileAndDeployOperator.sh configure the images to point to your favourite Container Registry
3. execute scripts/init.sh

## Clean up resources

Execute scripts/wipeAll.sh

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
  isNetworkPolicyActive: "true"
  sentry:
    replicas: 1
    clientName: "IronoaSentry"
    nodeKey: "0000000000000000000000000000000000000000000000000000000000000013" # Local node id: QmQMTLWkNwGf7P5MQv7kUHCynMg7jje6h3vbvwd2ALPPhm
    reservedValidatorID: "QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM"
    CPULimit: "0.5"
    memoryLimit: "512Mi"
    storageClassName: "default" #["default","managed-premium"]
  validator:
    clientName: "IronoaValidator"
    nodeKey: "0000000000000000000000000000000000000000000000000000000000000021" # Local node id: QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM
    reservedSentryID: "QmQMTLWkNwGf7P5MQv7kUHCynMg7jje6h3vbvwd2ALPPhm"
    CPULimit: "0.5"
    memoryLimit: "512Mi"
    storageClassName: "default" #["default","managed-premium"]
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
            - name: WATCH_NAMESPACE
              valueFrom:
                fieldRef:
                  fieldPath: metadata.namespace
            - name: POD_NAME
              valueFrom:
                fieldRef:
                  fieldPath: metadata.name
            - name: OPERATOR_NAME
              value: "polkadot-operator"
```

Change scripts/utils/compileAndDeployOperator.sh accordingly to the previous configured image value.
```sh
operator-sdk build ironoa/customresource-operator:v0.0.8 # define your favourite
docker push ironoa/customresource-operator:v0.0.8 #define your favourite
kubectl create -f deploy/operator.yaml
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

## Polkadot CR Configurable Parameters

* clientVersion: (string)  
Image version of the clients. See the Updating of Node Versions section.

* isNetworkPolicyActive: (string)  
If set to "true", the operator will handle the creation and the deployment of a Network Policy object that will ensure the secureness of the Validator (it only affects the Kind "SentryAndValidator"). 
With the parameter active, the Validator is allowed to communicate only with the Sentry layer. Being this mechanism enforced via NetworkPolicy (kubernetes native object), it requires a network plugin installed in you cloud provided cluster (even in minikube) to work properly.  
See the Secure Communications section.

* replicas: (int)  
Allows to decide how many Sentry replicas will be created. See the Node Cluster Scaling Support section.

* clientName: (string)

* CPULimit: (string)  
The format is the usual kubernetes and docker standard (e.g. "0.5")

* memoryLimit: (string)  
The format is the usual kubernetes and docker standard (e.g. "500Mi")

* nodeKey: (string)  
Identity of the node, private (e.g. "0000000000000000000000000000000000000000000000000000000000000013")

* storageClassName: (string)  
Desired volume type. For instance, Azure provides two built in storage classes:
    * "default": HHD backed
    * "managed-premium": SSD backed, high performance
    * See Data Persistence Support section for more information.     

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

## Node Cluster Scaling Support

This is the ability of the operator to respond to scale operations defined in the deployed configuration, for example to extend the amount of sentry nodes from 3 to 4. The correct functioning can be tested by executing such an operation and checking the number of deployed instances before and afterwards.  
In any case, Validator replica size is always hard coded to one and it is not possible to change it to prevent concurrent validation issues.
            
## Secure Communications (Kind:SentryAndValidator)

The configuration is based on the "polkadot-secure-validator" guidelines: https://github.com/w3f/polkadot-secure-validator

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

The current solution is using a Kubernetes Persistent Volume Claim with Azure Disk.
Reference: https://docs.microsoft.com/en-us/azure/aks/azure-disks-dynamic-pv

### How To Tutorial with Minikube

If you want to test it locally, you first have to manually provide a class of few persistent volumes (at least two, one for each client you deploy) to minikube. Minikube will extract from this pool an available volume thanks to the Persistent Volume Claim mechanism.   
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

The fields "hostPath -> path" describe where the folder will be mounted inside minikube. We have to create that folders and grant some permission first:

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
  isNetworkPolicyActive: "true"
  sentry:
    replicas: 1
    clientName: "IronoaSentry"
    nodeKey: "0000000000000000000000000000000000000000000000000000000000000013" # Local node id: QmQMTLWkNwGf7P5MQv7kUHCynMg7jje6h3vbvwd2ALPPhm
    reservedValidatorID: "QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM"
    CPULimit: "0.5"
    memoryLimit: "512Mi"
    storageClassName: "local"
  validator:
    clientName: "IronoaValidator"
    nodeKey: "0000000000000000000000000000000000000000000000000000000000000021" # Local node id: QmQtR1cdEaJM11qBWQBd34FoSgFichCjhtsBfrUFsVAjZM
    reservedSentryID: "QmQMTLWkNwGf7P5MQv7kUHCynMg7jje6h3vbvwd2ALPPhm"
    CPULimit: "0.5"
    memoryLimit: "512Mi"
    storageClassName: "local"
```

You can now deploy the operator as usual, also with the init.sh script.

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