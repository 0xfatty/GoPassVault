package main

import (
	"github.com/0xfatty/GoPassVault/pkg/vault"
	"github.com/0xfatty/GoPassVault/ui"
)

func main() {
	v, err := vault.LoadOrCreateVault("vault.dat")
	if err != nil {
		panic(err)
	}

	ui.ShowUI(v)
}
