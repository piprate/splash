package main

import (
	"context"
	"testing"

	"github.com/piprate/splash"
	"github.com/stretchr/testify/assert"
)

func TestScript(t *testing.T) {
	ctx := context.Background()
	g := splash.NewTestingEmulator()
	//t.Parallel()

	t.Run("Raw account argument", func(t *testing.T) {
		value := g.ScriptFromFile("test").RawAccountArgument("0x1cf0e2f2f715450").RunReturnsInterface(ctx)
		assert.Equal(t, "0x1cf0e2f2f715450", value)
	})

	t.Run("Raw account argument", func(t *testing.T) {
		value := g.ScriptFromFile("test").AccountArgument("first").RunReturnsInterface(ctx)
		assert.Equal(t, "0x1cf0e2f2f715450", value)
	})

	t.Run("Script should report failure", func(t *testing.T) {
		value, err := g.Script("asdf").RunReturns(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Parsing failed")
		assert.Nil(t, value)

	})

}
