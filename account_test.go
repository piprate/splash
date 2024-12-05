package splash_test

import (
	"context"
	"testing"

	. "github.com/piprate/splash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorsInAccountCreation(t *testing.T) {

	t.Run("Should give error on wrong saAccount name", func(t *testing.T) {
		g, err := NewInMemoryTestConnector("examples", false)
		require.NoError(t, err)
		_, err = g.CreateAccountsE(context.Background(), "foobar")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not find account with name foobar")
	})

	t.Run("Should give erro on wrong account name", func(t *testing.T) {
		_, err := NewInMemoryConnector([]string{"invalid_account_in_deployment.json"}, NewFileSystemLoader("fixtures"), false, NewZeroLogger())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "deployment contains nonexisting account emulator-firs")
	})

	t.Run("Should fail when creating local accounts in the wrong order", func(t *testing.T) {
		g, err := NewInMemoryConnector([]string{"wrong_account_order_emulator.json"}, NewFileSystemLoader("fixtures"), false, NewZeroLogger())
		require.NoError(t, err)
		_, err = g.CreateAccountsE(context.Background(), "emulator-first")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "could not find account with address 179b6b1cb6755e3")
	})
}
