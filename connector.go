package splash

import (
	"fmt"
	"log"

	"github.com/onflow/flow-emulator/emulator"
	"github.com/onflow/flow-go-sdk/access"
	grpcAccess "github.com/onflow/flow-go-sdk/access/grpc"
	"github.com/onflow/flowkit/v2"
	"github.com/onflow/flowkit/v2/accounts"
	"github.com/onflow/flowkit/v2/config"
	"github.com/onflow/flowkit/v2/gateway"
	"github.com/onflow/flowkit/v2/output"
	"github.com/spf13/afero"
	"google.golang.org/grpc"
)

// Connector Entire configuration to work with Splash
type Connector struct {
	State                        *flowkit.State
	GRPCClient                   access.Client
	Services                     flowkit.Services
	Network                      string
	Logger                       output.Logger
	PrependNetworkToAccountNames bool
}

// maxGRPCMessageSize 60mb
const maxGRPCMessageSize = 1024 * 1024 * 60

// NewNetworkConnector creates a new local go with the flow client
func NewNetworkConnector(paths []string, baseLoader flowkit.ReaderWriter, network string, logger output.Logger) (*Connector, error) {

	state, err := flowkit.Load(paths, baseLoader)
	if err != nil {
		return nil, err
	}

	networkDef, err := state.Networks().ByName(network)
	if err != nil {
		return nil, err
	}
	gw, err := gateway.NewGrpcGateway(*networkDef)
	if err != nil {
		return nil, err
	}
	service := flowkit.NewFlowkit(state, *networkDef, gw, logger)

	grpcClient, err := grpcAccess.NewClient(
		networkDef.Host,
		grpcAccess.WithGRPCDialOptions(
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxGRPCMessageSize)),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to host %s", networkDef.Host)
	}

	return &Connector{
		State:                        state,
		Services:                     service,
		GRPCClient:                   grpcClient,
		Logger:                       logger,
		PrependNetworkToAccountNames: true,
	}, nil
}

func NewInMemoryConnector(paths []string, baseLoader flowkit.ReaderWriter, enableTxFees bool, logger output.Logger) (*Connector, error) {

	state, err := flowkit.Load(paths, baseLoader)
	if err != nil {
		return nil, err
	}

	acc, _ := state.EmulatorServiceAccount()
	pk, _ := acc.Key.PrivateKey()

	key := &gateway.EmulatorKey{
		PublicKey: (*pk).PublicKey(),
		SigAlgo:   acc.Key.SigAlgo(),
		HashAlgo:  acc.Key.HashAlgo(),
	}
	var gw *gateway.EmulatorGateway
	if enableTxFees {
		gw = gateway.NewEmulatorGatewayWithOpts(key, gateway.WithEmulatorOptions(emulator.WithTransactionFeesEnabled(true)))
	} else {
		gw = gateway.NewEmulatorGateway(key)
	}
	service := flowkit.NewFlowkit(state, config.EmulatorNetwork, gw, logger)

	return &Connector{
		State:                        state,
		Services:                     service,
		Logger:                       logger,
		PrependNetworkToAccountNames: true,
	}, nil
}

func NewConnectorDefault(network string, logLevel int) (*Connector, error) {
	loader := &afero.Afero{Fs: afero.NewOsFs()}
	stdoutLogger := output.NewStdoutLogger(logLevel)
	return NewNetworkConnector(config.DefaultPaths(), loader, network, stdoutLogger)
}

func NewInMemoryTestConnector(baseDir string, enableTxFees bool) (*Connector, error) {
	return NewInMemoryConnector([]string{config.DefaultPath}, NewFileSystemLoader(baseDir), enableTxFees, NewZeroLogger())
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
