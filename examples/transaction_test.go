package main

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/bjartek/go-with-the-flow/gwtf"
)

/*
 Tests must be in the same folder as flow.json with contracts and transactions/scripts in subdirectories in order for the path resolver to work correctly
*/
func TestTransaction(t *testing.T) {
	g := gwtf.NewGoWithTheFlowInMemoryEmulator()

	t.Parallel()
	t.Run("Create NFT collection", func(t *testing.T) {
		g.TransactionFromFile("create_nft_collection").
			SignProposeAndPayAs("first").
			Test(t).
			AssertSuccess().
			AssertNoEvents()
	})

	t.Run("Mint tokens assert events", func(t *testing.T) {
		g.TransactionFromFile("mint_tokens").
			SignProposeAndPayAsService().
			AccountArgument("first").
			UFix64Argument("100.0").
			Test(t).
			AssertSuccess().
			AssertEventCount(3).
			AssertEmitEventName("A.0ae53cb6e3f42a79.FlowToken.TokensMinted").
			AssertEmitEventName("A.0ae53cb6e3f42a79.FlowToken.TokensDeposited").
			AssertEmitEventName("A.0ae53cb6e3f42a79.FlowToken.MinterCreated").
			AssertEmitEventName("A.0ae53cb6e3f42a79.FlowToken.TokensMinted", "A.0ae53cb6e3f42a79.FlowToken.TokensDeposited", "A.0ae53cb6e3f42a79.FlowToken.MinterCreated").
			AssertEmitEvent(gwtf.NewTestEvent("A.0ae53cb6e3f42a79.FlowToken.TokensMinted", map[string]interface{}{ "amount": "100.00000000" }))


	})

	t.Run("Inline transaction with debug log", func(t *testing.T){
		g.Transaction(`
import Debug from "../contracts/Debug.cdc"
transaction(message:String) {
  prepare(acct: AuthAccount, account2: AuthAccount) {
	Debug.log(message)
 }
}`).
			SignProposeAndPayAs("first").
			PayloadSigner("second").
			StringArgument("foobar").
			Test(t).
			AssertSuccess().
			AssertDebugLog("foobar")


	})

	t.Run("Upload test file", func(t *testing.T) {
		err := g.UploadFile("testFile.txt", "first")
		assert.NoError(t, err)
		g.Transaction(`
import Debug from "../contracts/Debug.cdc"
transaction {
  prepare(account: AuthAccount) {
    var content= account.load<String>(from: /storage/upload) ?? panic("could not load content")
	Debug.log(content)
 }
}`).
			SignProposeAndPayAs("first").
			Test(t).
			AssertSuccess().
			AssertDebugLog("VGhpcyBpcyBhIGZpbGU=")
	})


}



