package splash_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/onflow/flowkit/v2/config"
	. "github.com/piprate/splash"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp})
}

func TestNewInMemoryConnector_NoTxFees(t *testing.T) {
	client, err := NewInMemoryConnector([]string{config.DefaultPath}, NewFileSystemLoader("examples"), false, NewZeroLogger())
	require.NoError(t, err)

	ctx := context.Background()

	_, err = client.CreateAccountsE(ctx, "emulator-account")
	require.NoError(t, err)

	err = client.InitializeContractsE(ctx)
	require.NoError(t, err)
}

func TestNewInMemoryConnector_WithTxFees(t *testing.T) {
	client, err := NewInMemoryConnector([]string{config.DefaultPath}, NewFileSystemLoader("examples"), true, NewZeroLogger())
	require.NoError(t, err)

	ctx := context.Background()

	_, err = client.DoNotPrependNetworkToAccountNames().CreateAccountsE(ctx, "emulator-account")
	require.NoError(t, err)
}
