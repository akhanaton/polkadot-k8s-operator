// Copyright (c) 2020 Swisscom Blockchain AG
// Licensed under MIT License
package polkadot

const (
	ServiceSentryName    = "sentry-service"
	ServiceValidatorName = "validator-service"
	metricsPortName        = "http-metrics"
	P2PPortName            = "p2p"
	RPCPortName            = "http-rpc"
	WSPortName             = "websocket-rpc"
	ValidatorSSName        = "validator-sset"
	SentrySSName           = "sentry-sset"
	ValidatorNetworkPolicy = "validator-networkpolicy"
	volumeMountPath        = "/data"
	serviceName            = "polkadot"
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

func getValidatorLabels() map[string]string {
	labels := getAppLabels()
	labels["role"] = "validator"
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