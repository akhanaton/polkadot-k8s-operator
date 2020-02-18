package controller

import (
	"github.com/swisscom-blockchain/polkadot-k8s-operator/pkg/controller/polkadot"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, polkadot.Add)
}
