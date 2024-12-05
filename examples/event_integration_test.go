package examples_test

import (
	"context"
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/piprate/splash"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.Stamp})
}

func TestEvents(t *testing.T) {

	t.Run("Test that from index cannot be negative", func(t *testing.T) {
		ctx := context.Background()
		g, err := splash.NewInMemoryTestConnector(".", false)
		require.NoError(t, err)
		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			AccountArgument("zero").
			UFix64Argument("100.0").
			Test(t).
			AssertSuccess().
			AssertEventCount(4)

		_, err = g.EventFetcher().End(2).From(-10).Event("A.0ae53cb6e3f42a79.FlowToken.TokensMinted").Run(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "FromIndex is negative")
	})

	t.Run("Fetch last events", func(t *testing.T) {
		ctx := context.Background()
		g, err := splash.NewInMemoryTestConnector(".", false)
		require.NoError(t, err)
		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			AccountArgument("zero").
			UFix64Argument("100.0").
			Test(t).
			AssertSuccess().
			AssertEventCount(4)

		ev, err := g.EventFetcher().Last(1).Event("A.0ae53cb6e3f42a79.FlowToken.TokensMinted").Run(ctx)
		require.NoError(t, err)
		assert.Equal(t, 1, len(ev))
	})

	t.Run("Fetch last events and sort them", func(t *testing.T) {
		ctx := context.Background()
		g, err := splash.NewInMemoryTestConnector(".", false)
		require.NoError(t, err)
		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			AccountArgument("zero").
			UFix64Argument("100.0").
			Test(t).
			AssertSuccess().
			AssertEventCount(4)

		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			AccountArgument("zero").
			UFix64Argument("100.0").
			Test(t).
			AssertSuccess().
			AssertEventCount(4)

		ev, err := g.EventFetcher().Last(1).Event("A.0ae53cb6e3f42a79.FlowToken.TokensMinted").Run(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 2, len(ev))
		assert.True(t, ev[0].BlockHeight < ev[1].BlockHeight)
	})

	t.Run("Fetch last write progress file", func(t *testing.T) {
		ctx := context.Background()
		g, err := splash.NewInMemoryTestConnector(".", false)
		require.NoError(t, err)
		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			AccountArgument("zero").
			UFix64Argument("100.0").
			Test(t).
			AssertSuccess().
			AssertEventCount(4)

		ev, err := g.EventFetcher().Event("A.0ae53cb6e3f42a79.FlowToken.TokensMinted").TrackProgressIn("progress").Run(ctx)
		defer os.Remove("progress")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(ev))
	})

	t.Run("should fail reading invalid progress from file", func(t *testing.T) {
		ctx := context.Background()
		err := os.WriteFile("progress", []byte("invalid"), fs.ModePerm)
		assert.NoError(t, err)

		g, err := splash.NewInMemoryTestConnector(".", false)
		require.NoError(t, err)
		_, err = g.EventFetcher().Event("A.0ae53cb6e3f42a79.FlowToken.TokensMinted").TrackProgressIn("progress").Run(ctx)
		defer os.Remove("progress")
		assert.Error(t, err)
		assert.Equal(t, "could not parse progress file as block height strconv.ParseInt: parsing \"invalid\": invalid syntax", err.Error())
	})

	t.Run("Fetch last write progress file that exists", func(t *testing.T) {
		ctx := context.Background()
		err := os.WriteFile("progress", []byte("1"), fs.ModePerm)
		assert.NoError(t, err)

		g, err := splash.NewInMemoryTestConnector(".", false)
		require.NoError(t, err)

		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			AccountArgument("zero").
			UFix64Argument("100.0").
			Test(t).
			AssertSuccess().
			AssertEventCount(4)

		ev, err := g.EventFetcher().Event("A.0ae53cb6e3f42a79.FlowToken.TokensMinted").TrackProgressIn("progress").Run(ctx)
		defer os.Remove("progress")
		assert.NoError(t, err)
		assert.Equal(t, 1, len(ev))
	})

}
