package splash

import (
	"context"
	"fmt"
	"log"
	"sort"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flowkit/v2"
	"github.com/onflow/flowkit/v2/accounts"
)

func (c *Connector) CreateAccounts(ctx context.Context, saAccountName string) *Connector {
	conn, err := c.CreateAccountsE(ctx, saAccountName)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

// CreateAccountsE ensures that all accounts present in the deployment block for the given network is present
func (c *Connector) CreateAccountsE(ctx context.Context, saAccountName string) (*Connector, error) {
	p := c.State
	signerAccount, err := p.Accounts().ByName(saAccountName)
	if err != nil {
		return nil, err
	}

	accountList := *p.AccountsForNetwork(c.Services.Network())
	accountNames := accountList.Names()
	sort.Strings(accountNames)

	c.Logger.Info(fmt.Sprintf("%v\n", accountNames))

	for _, accountName := range accountNames {
		c.Logger.Debug(fmt.Sprintf("Ensuring account with name '%s' is present", accountName))

		// this error can never happen here, there is a test for it.
		account, _ := p.Accounts().ByName(accountName)

		if _, err := c.Services.GetAccount(ctx, account.Address); err == nil {
			c.Logger.Debug("Account is present")
			continue
		}

		a, _, err := c.Services.CreateAccount(
			ctx,
			signerAccount,
			[]accounts.PublicKey{{
				Public:   account.Key.ToConfig().PrivateKey.PublicKey(),
				Weight:   flow.AccountKeyWeightThreshold,
				SigAlgo:  account.Key.SigAlgo(),
				HashAlgo: account.Key.HashAlgo(),
			}})
		if err != nil {
			return nil, err
		}
		c.Logger.Info("Account created " + a.Address.String())
		if a.Address.String() != account.Address.String() {
			// this condition happens when we create accounts defined in flow.json
			// after some other accounts were created manually.
			// In this case, account addresses may not match the expected values
			c.Logger.Error("Account address mismatch. Expected " + account.Address.String() + ", got " + a.Address.String())
		}
	}
	return c, nil
}

// InitializeContracts installs all contracts in the deployment block for the configured network
func (c *Connector) InitializeContracts(ctx context.Context) *Connector {
	if err := c.InitializeContractsE(ctx); err != nil {
		log.Fatal(err)
	}

	return c
}

// InitializeContractsE installs all contracts in the deployment block for the configured network
// and returns an error if it fails.
func (c *Connector) InitializeContractsE(ctx context.Context) error {
	c.Logger.Info("Deploying contracts")
	if _, err := c.Services.DeployProject(ctx, flowkit.UpdateExistingContract(true)); err != nil {
		return err
	}

	return nil
}
