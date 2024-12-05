package examples_test

import (
	"context"
	"testing"

	"github.com/piprate/splash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScript(t *testing.T) {
	ctx := context.Background()
	g, err := splash.NewInMemoryTestConnector(".", false)
	require.NoError(t, err)

	t.Parallel()

	t.Run("Raw account argument", func(t *testing.T) {
		value := g.ScriptFromFile("test").RawAccountArgument("0x01cf0e2f2f715450").RunReturnsInterface(ctx)
		assert.Equal(t, "0x01cf0e2f2f715450", value)
	})

	t.Run("Account argument", func(t *testing.T) {
		value := g.ScriptFromFile("test").AccountArgument("first").RunReturnsInterface(ctx)
		assert.Equal(t, "0x179b6b1cb6755e31", value)
	})

	t.Run("Script should report failure", func(t *testing.T) {
		value, err := g.Script("asdf").RunReturns(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Parsing failed")
		assert.Nil(t, value)

	})

}
