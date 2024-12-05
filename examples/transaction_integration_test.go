package examples_test

import (
	"bytes"
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/piprate/splash"
)

/*
Tests must be in the same folder as flow.json with contracts and transactions/scripts in subdirectories in order for the path resolver to work correctly
*/
func TestTransaction(t *testing.T) {
	g, err := splash.NewInMemoryTestConnector(".", false)
	require.NoError(t, err)

	ctx := context.Background()
	err = g.CreateAccounts(ctx, "emulator-account").InitializeContractsE(ctx)
	require.NoError(t, err)

	t.Parallel()

	t.Run("Fail on missing signer", func(t *testing.T) {
		g.TransactionFromFile("create_nft_collection").
			ProposeAs("first").
			PayAs("first").
			Test(t).                                                                                      //This method will return a TransactionResult that we can assert upon
			AssertFailure("provided authorizers length mismatch, required authorizers 1, but provided 0") //we assert that there is a failure
	})

	t.Run("Fail on wrong transaction name", func(t *testing.T) {
		g.TransactionFromFile("create_nf_collection").
			SignProposeAndPayAs("first").
			Test(t).                                                                                           //This method will return a TransactionResult that we can assert upon
			AssertFailure("could not read transaction file from path=./transactions/create_nf_collection.cdc") //we assert that there is a failure
	})

	t.Run("Create NFT collection", func(t *testing.T) {
		g.TransactionFromFile("create_nft_collection").
			SignProposeAndPayAs("first").
			Test(t).         //This method will return a TransactionResult that we can assert upon
			AssertSuccess(). //Assert that there are no errors and that the transactions succeeds
			AssertEventCount(2)
	})

	t.Run("Mint tokens assert events", func(t *testing.T) {
		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			AccountArgument("first").
			UFix64Argument("100.0").
			Test(t).
			AssertSuccess().
			AssertEventCount(4).                                                                                                                                                                           //assert the number of events returned
			AssertPartialEvent(splash.NewTestEvent("A.0ae53cb6e3f42a79.FlowToken.TokensDeposited", map[string]interface{}{"amount": "100.00000000"})).                                                     //assert a given event, can also take multiple events if you like
			AssertEmitEventName("A.0ae53cb6e3f42a79.FlowToken.TokensMinted").                                                                                                                              //assert the name of a single event
			AssertEmitEventName("A.0ae53cb6e3f42a79.FlowToken.TokensMinted", "A.0ae53cb6e3f42a79.FlowToken.TokensDeposited", "A.0ae53cb6e3f42a79.FlowToken.MinterCreated").                                //or assert more then one eventname in a go
			AssertEmitEvent(splash.NewTestEvent("A.0ae53cb6e3f42a79.FlowToken.TokensMinted", map[string]interface{}{"amount": "100.00000000"})).                                                           //assert a given event, can also take multiple events if you like
			AssertEmitEventJSON("{\n  \"name\": \"A.0ae53cb6e3f42a79.FlowToken.MinterCreated\",\n  \"time\": \"1970-01-01T00:00:00Z\",\n  \"fields\": {\n    \"allowedAmount\": \"100.00000000\"\n  }\n}") //assert a given event using json, can also take multiple events if you like

	})

	t.Run("Inline transaction with debug log", func(t *testing.T) {
		g.Transaction(`
import Debug from "../contracts/Debug.cdc"
transaction(message:String) {
  prepare(acct: auth(BorrowValue) &Account, account2: auth(BorrowValue) &Account) {
	Debug.log(message)
 }
}`).
			SignProposeAndPayAs("first").
			PayloadSigner("second").
			StringArgument("foobar").
			Test(t).
			AssertSuccess().
			AssertDebugLog("foobar") //assert that we have debug logged something. The assertion is contains so you do not need to write the entire debug log output if you do not like

	})

	t.Run("Raw account argument", func(t *testing.T) {
		g.Transaction(`
import Debug from "../contracts/Debug.cdc"
transaction(user:Address) {
  prepare(acct: auth(BorrowValue) &Account) {
	Debug.log(user.toString())
 }
}`).
			SignProposeAndPayAsService().
			RawAccountArgument("0x01cf0e2f2f715450").
			Test(t).
			AssertSuccess().
			AssertDebugLog("0x01cf0e2f2f715450")
	})

	t.Run("transaction that should fail", func(t *testing.T) {
		g.Transaction(`
import Debug from "../contracts/Debug.cdc"
transaction(user:Address) {
  prepare(acct: auth(BorrowValue) &Account) {
	Debug.log(user.toStrig())
 }
}`).
			SignProposeAndPayAsService().
			RawAccountArgument("0x01cf0e2f2f715450").
			Test(t).
			AssertFailure("has no member `toStrig`") //assert failure with an error message. uses contains so you do not need to write entire message
	})

	t.Run("Assert print events", func(t *testing.T) {
		ctx := context.Background()
		var str bytes.Buffer
		log.SetOutput(&str)
		defer log.SetOutput(os.Stdout)

		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			AccountArgument("first").
			UFix64Argument("100.0").RunPrintEventsFull(ctx)

		assert.Contains(t, str.String(), "A.0ae53cb6e3f42a79.FlowToken.MinterCreated")

	})

	t.Run("Assert print events", func(t *testing.T) {
		var str bytes.Buffer
		log.SetOutput(&str)
		defer log.SetOutput(os.Stdout)

		ctx := context.Background()
		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			AccountArgument("first").
			UFix64Argument("100.0").
			RunPrintEvents(ctx, map[string][]string{"A.0ae53cb6e3f42a79.FlowToken.TokensDeposited": {"to"}})

		assert.NotContains(t, str.String(), "0x01cf0e2f2f715450")
	})
}
