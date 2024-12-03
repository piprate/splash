package examples_test

import (
	"testing"

	"github.com/piprate/splash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupIntegration(t *testing.T) {

	t.Run("Should create inmemory emulator client", func(t *testing.T) {
		g, err := splash.NewConnectorInMemoryEmulator()
		require.NoError(t, err)
		assert.Equal(t, "emulator", g.Network)
	})

	t.Run("Should create local emulator client", func(t *testing.T) {
		g, err := splash.NewConnectorInMemoryEmulator()
		require.NoError(t, err)
		assert.Equal(t, "emulator", g.Network)
	})

	t.Run("Should create testnet client", func(t *testing.T) {
		g, err := splash.NewConnectorTestNet()
		require.NoError(t, err)
		assert.Equal(t, "testnet", g.Network)
	})

	t.Run("Should create testnet client with for network method", func(t *testing.T) {
		g, err := splash.NewConnectorForNetwork("testnet")
		require.NoError(t, err)
		assert.Equal(t, "testnet", g.Network)
	})

	t.Run("Should create mainnet client", func(t *testing.T) {
		g, err := splash.NewConnectorMainNet()
		require.NoError(t, err)
		assert.Equal(t, "mainnet", g.Network)
		assert.True(t, g.PrependNetworkToAccountNames)
		g = g.DoNotPrependNetworkToAccountNames()
		assert.False(t, g.PrependNetworkToAccountNames)

	})
}
