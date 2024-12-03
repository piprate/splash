package examples

import (
	"context"

	"github.com/onflow/flowkit/v2/output"
	"github.com/piprate/splash"
)

// NewConnectorInMemoryEmulator this method is used to create an in memory emulator, deploy all contracts for the emulator and create all accounts
func NewConnectorInMemoryEmulator() (*splash.Connector, error) {
	ctx := context.Background()
	conn, err := splash.NewConnectorDefault("emulator", output.InfoLog)
	if err != nil {
		return nil, err
	}
	return conn.InitializeContracts(ctx).CreateAccountsE(ctx, "emulator-account")
}

// NewConnectorForNetwork creates a new splash client for the provided network
func NewConnectorForNetwork(network string) (*splash.Connector, error) {
	return splash.NewConnectorDefault(network, output.InfoLog)
}

// NewConnectorEmulator create a new splash client for the local emulator
func NewConnectorEmulator() (*splash.Connector, error) {
	return splash.NewConnectorDefault("emulator", output.InfoLog)
}

// NewConnectorTestNet creates a new splash client for devnet/testnet
func NewConnectorTestNet() (*splash.Connector, error) {
	return splash.NewConnectorDefault("testnet", output.InfoLog)
}

// NewConnectorMainNet creates a new splash client for mainnet
func NewConnectorMainNet() (*splash.Connector, error) {
	return splash.NewConnectorDefault("mainnet", output.InfoLog)
}
