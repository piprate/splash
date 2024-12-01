package main

import (
	"testing"

	"github.com/piprate/splash"
	"github.com/stretchr/testify/assert"
)

func TestSetupIntegration(t *testing.T) {

	t.Run("Should create inmemory emulator client", func(t *testing.T) {
		g := splash.NewConnectorInMemoryEmulator()
		assert.Equal(t, "emulator", g.Network)
	})

	t.Run("Should create local emulator client", func(t *testing.T) {
		g := splash.NewConnectorEmulator()
		assert.Equal(t, "emulator", g.Network)
	})

	t.Run("Should create testnet client", func(t *testing.T) {
		g := splash.NewConnectorTestNet()
		assert.Equal(t, "testnet", g.Network)
	})

	t.Run("Should create testnet client with for network method", func(t *testing.T) {
		g := splash.NewConnectorForNetwork("testnet")
		assert.Equal(t, "testnet", g.Network)
	})

	t.Run("Should create mainnet client", func(t *testing.T) {
		g := splash.NewConnectorMainNet()
		assert.Equal(t, "mainnet", g.Network)
		assert.True(t, g.PrependNetworkToAccountNames)
		g = g.DoNotPrependNetworkToAccountNames()
		assert.False(t, g.PrependNetworkToAccountNames)

	})
}
