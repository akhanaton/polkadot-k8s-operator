[//]: <> (Copyright (c) 2020 Swisscom Blockchain AG)
[//]: <> (Licensed under MIT License)


# PolkaKop - Kubernetes Operator for Polkadot 

Kubernetes Operator for sentry nodes and validators

## Requirements

* Docker
* Kubernetes kubectl
* Go compiler
* Access to a container registry
* Operator SDK CLI (https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md)
* Oprional (SentryAndValidator secure depoloyment): kubenet plugin, se the CR Configurable Parameters section

## How To Run

Deploy to your favorite kubernetes cloud provided cluster (even minikube) a Custom Controller and a Polkadot Custom Resource. The Controller will create and supervise all the necessary resource needed to run a Polkadot (Rust) Client.

0. Configure your kubectl to work with your desidered Kubernetes cluster 
    (e.g. Azure: az aks get-credentials --resource-group myResourceGroup --name myAKSCluster)
1. Clone the repository locally
2. Configure both deploy/operator.yaml image and scripts/compileAndDeployOperator.sh images to point to your favourite Container Registry
3. execute scripts/init.sh

## Clean up resources

Execute scripts/wipeAll.sh

## CR Configurable Parameters

* Kind: Sentry|Validator|SentryAndValidator (string)
Allows you to decide what to deploy:
    * Sentry: deploy a Sentry node
    * Validator: deploy a validator node
    * SentryAndValidator: deploy a sentry and a validator in a secure configuration. The validator is allowed to communicate only through the sentry. This mechanism is enforced also via NetworkPolicy kubernetes native object, which requires a kubenet plugin installed in you cloud provided cluster (even in minikube) to work properly.
        * In the SentryAndValidator configuration it must be passed also an additional parameter to both the sentry and the validator:
        * reservedValidatorID: (string) Identiry of the validator, it must be set for the sentry
        * reservedSentryID: Identiry of the sentry, it must be set for the validator
        
            ![alt text](images/schema.png)

* replicas: (int)
Allows you to decide how many Sentry replicas will be created. Validator replica size is always hard coded to one and it is not possible to change it to prevent concurrent validation issues.

* clientName: (string)

* CPULimit: (string)
The format is the usual kubernetes and docker standard (e.g. "0.5")

* memoryLimit: (string)
The format is the usual kubernetes and docker standard (e.g. "500Mi")

* nodeKey: (string)
Identity of the node, private (e.g. "0000000000000000000000000000000000000000000000000000000000000013")

## About the Operator

An Operator is a method of packaging, deploying and managing a Kubernetes application. A Kubernetes application is an application that is both deployed on Kubernetes and managed using the Kubernetes APIs and kubectl tooling. To be able to make the most of Kubernetes, you need a set of cohesive APIs to extend in order to service and manage your applications that run on Kubernetes. You can think of Operators as the runtime that manages this type of application on Kubernetes.

Reference: https://coreos.com/blog/introducing-operator-framework

## About the Operator SDK and Framework

To help make it easier to build Kubernetes applications, Red Hat and the Kubernetes open source community today share the Operator Framework – an open source toolkit designed to manage Kubernetes native applications, called Operators, in a more effective, automated, and scalable way. 

The Operator SDK provides the tools to build, test and package Operators. Initially, the SDK facilitates the marriage of an application’s business logic (for example, how to scale, upgrade, or backup) with the Kubernetes API to execute those operations. 

Reference: https://coreos.com/blog/introducing-operator-framework



