package config

import (
	"fmt"
	"os"
	"strconv"
)

type EnvVar struct {
	name string
	Value string
}

type EnvVarInt struct {
	name string
	Value int
}

// these vars are set by the main function at the startup
var (
	ControllerNameEnvVar = EnvVar{"CONTROLLER_NAME", ""}
	ImageClientEnvVar    = EnvVar{"IMAGE_CLIENT", ""}
	MetricsPortEnvVar   = EnvVarInt{"METRICS_PORT",-1}
	P2PPortEnvVar   = EnvVarInt{"P2P_PORT",-1}
	RPCPortEnvVar   = EnvVarInt{"RPC_PORT",-1}
	WSPortEnvVar   = EnvVarInt{"WS_PORT",-1}
)

// this function is called by the main at the startup
func LoadAllEnvVar() error {
	var err error
	if err = loadEnvVar(&ControllerNameEnvVar); err != nil {return err }
	if err = loadEnvVar(&ImageClientEnvVar); err != nil {return err }
	if err = loadEnvVarInt(&MetricsPortEnvVar); err != nil {return err }
	if err = loadEnvVarInt(&P2PPortEnvVar); err != nil {return err }
	if err = loadEnvVarInt(&RPCPortEnvVar); err != nil {return err }
	if err = loadEnvVarInt(&WSPortEnvVar); err != nil {return err }
	return err
}

func loadEnvVar(envVar *EnvVar) error{
	var err error
	envVar.Value, err = getEnvVar(envVar.name)
	return err
}

func loadEnvVarInt(envVar *EnvVarInt) error{
	var err error
	stringValue, err := getEnvVar(envVar.name)
	if err != nil {
		return err
	}
	intValue, err := strconv.Atoi(stringValue)
	if err != nil {
		return err
	}
	envVar.Value = intValue
	return err
}

func getEnvVar(name string) (string,error){
	n, isFound := os.LookupEnv(name)
	if !isFound {
		return "", fmt.Errorf("%s must be set", name)
	}
	return n, nil
}
