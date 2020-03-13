// Copyright (c) 2020 Swisscom Blockchain AG
// Licensed under MIT License
package polkadot

const (
	imageName              = "parity/polkadot"
	imageNameMetrics 	   = "ironoa/polkadot-metrics:v0.0.1" //define your favourite
	serviceSentryName      = "sentry-service"
	serviceValidatorName   = "validator-service"
	metricsPort			   = 8000
	metricsPortName		   = "http-metrics"
	P2PPort                = 30333
	P2PPortName            = "p2p"
	RPCPort                = 9933
	RPCPortName            = "http-rpc"
	WSPort                 = 9944
	WSPortName             = "websocket-rpc"
	validatorSSName        = "validator-sset"
	sentrySSName           = "sentry-sset"
	validatorNetworkPolicy = "validator-networkpolicy"
	volumeMountPath        = "/data"
	volumeName 			   = "polkadot-volume"
	storageRequest         = "10Gi"
	serviceName 		   = "polkadot"
)


func getAppLabels() map[string]string {
	labels := map[string]string{"app": "polkadot"}
	return labels
}

func getSentrylabels() map[string]string {
	labels := getAppLabels()
	labels["role"] = "sentry"
	return labels
}

func getCopyLabelsWithVersion(labels map[string]string, version string) map[string]string {
	newLabels := getCopy(labels)
	newLabels["version"] = version
	return newLabels
}

func getCopy(originalMap map[string]string) map[string]string {
	newMap := make(map[string]string)
	for key, value := range originalMap {
		newMap[key] = value
	}
	return newMap
}

func getValidatorLabels() map[string]string {
	labels := getAppLabels()
	labels["role"] = "validator"
	return labels
}