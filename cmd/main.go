package main

import (
	"github.com/0xfatty/GoPassVault/pkg/vault"
	"github.com/0xfatty/GoPassVault/ui"
)

func main() {
	v := vault.NewVault("vault.txt")
	v.Load() // Load from disk
	ui.ShowUI(v)
}
