package splash

import (
	"context"
	"fmt"

	"log"

	"github.com/onflow/flow-go-sdk/access"
	"github.com/onflow/flowkit/v2"
	"github.com/onflow/flowkit/v2/accounts"
	"github.com/onflow/flowkit/v2/config"
	"github.com/onflow/flowkit/v2/gateway"
	"github.com/onflow/flowkit/v2/output"
	"github.com/spf13/afero"
)

// Connector Entire configuration to work with Splash
type Connector struct {
	State                        *flowkit.State
	Client                       access.Client
	Services                     flowkit.Services
	Network                      string
	Logger                       output.Logger
	PrependNetworkToAccountNames bool
}

// NewConnectorInMemoryEmulator this method is used to create an in memory emulator, deploy all contracts for the emulator and create all accounts
func NewConnectorInMemoryEmulator() *Connector {
	ctx := context.Background()
	return NewConnector(config.DefaultPaths(), "emulator", true, output.InfoLog).InitializeContracts(ctx).CreateAccounts(ctx, "emulator-account")
}

// NewTestingEmulator create new emulator that ignore all log messages
func NewTestingEmulator() *Connector {
	ctx := context.Background()
	return NewConnector(config.DefaultPaths(), "emulator", true, output.NoneLog).InitializeContracts(ctx).CreateAccounts(ctx, "emulator-account")
}

// NewConnectorForNetwork creates a new splash client for the provided network
func NewConnectorForNetwork(network string) *Connector {
	return NewConnector(config.DefaultPaths(), network, false, output.InfoLog)
}

// NewConnectorEmulator create a new client
func NewConnectorEmulator() *Connector {
	return NewConnector(config.DefaultPaths(), "emulator", false, output.InfoLog)
}

// NewConnectorTestNet creates a new splash client for devnet/testnet
func NewConnectorTestNet() *Connector {
	return NewConnector(config.DefaultPaths(), "testnet", false, output.InfoLog)
}

// NewConnectorMainNet creates a new gwft client for mainnet
func NewConnectorMainNet() *Connector {
	return NewConnector(config.DefaultPaths(), "mainnet", false, output.InfoLog)
}

// NewConnector with custom file panic on error
func NewConnector(filenames []string, network string, inMemory bool, loglevel int) *Connector {
	conn, err := NewConnectorWithError(filenames, network, inMemory, loglevel)
	if err != nil {
		log.Fatalf("error %+v", err)
	}
	return conn
}

// DoNotPrependNetworkToAccountNames disable the default behavior of prefixing account names with network-
func (c *Connector) DoNotPrependNetworkToAccountNames() *Connector {
	c.PrependNetworkToAccountNames = false
	return c
}

// Account fetch an account from flow.json, prefixing the name with network- as default (can be turned off)
func (c *Connector) Account(key string) *accounts.Account {
	if c.PrependNetworkToAccountNames {
		key = fmt.Sprintf("%s-%s", c.Services.Network().Name, key)
	}

	account, err := c.State.Accounts().ByName(key)
	if err != nil {
		log.Fatal(err)
	}

	return account
}

// NewConnectorWithError creates a new local go with the flow client
func NewConnectorWithError(paths []string, network string, inMemory bool, logLevel int) (*Connector, error) {

	loader := &afero.Afero{Fs: afero.NewOsFs()}
	state, err := flowkit.Load(paths, loader)
	if err != nil {
		return nil, err
	}

	logger := output.NewStdoutLogger(logLevel)
	var service flowkit.Services
	if inMemory {
		// YAY, we can run it inline in memory!
		acc, _ := state.EmulatorServiceAccount()
		pk, _ := acc.Key.PrivateKey()
		gw := gateway.NewEmulatorGateway(&gateway.EmulatorKey{
			PublicKey: (*pk).PublicKey(),
			SigAlgo:   acc.Key.SigAlgo(),
			HashAlgo:  acc.Key.HashAlgo(),
		})
		service = flowkit.NewFlowkit(state, config.EmulatorNetwork, gw, logger)
	} else {
		networkDef, err := state.Networks().ByName(network)
		if err != nil {
			return nil, err
		}
		gw, err := gateway.NewGrpcGateway(*networkDef)
		if err != nil {
			return nil, err
		}
		service = flowkit.NewFlowkit(state, *networkDef, gw, logger)
	}
	return &Connector{
		State:                        state,
		Services:                     service,
		Network:                      network,
		Logger:                       logger,
		PrependNetworkToAccountNames: true,
	}, nil
}
