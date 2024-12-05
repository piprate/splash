import NonFungibleToken from "../contracts/NonFungibleToken.cdc"
import ExampleNFT from "../contracts/ExampleNFT.cdc"

// This transaction creates an empty NFT Collection in the signer's account
transaction {
  prepare(acct: auth(BorrowValue, IssueStorageCapabilityController, PublishCapability, SaveValue, UnpublishCapability) &Account) {
    // store an empty NFT Collection in account storage
    let collection <- ExampleNFT.createEmptyCollection(nftType: Type<@ExampleNFT.NFT>())
    acct.storage.save(<-collection, to: /storage/NFTCollection)

    // publish a capability to the Collection in storage
    let collectionCap = acct.capabilities.storage.issue<&ExampleNFT.Collection>(/storage/NFTCollection)
    acct.capabilities.publish(collectionCap, at: /public/NFTReceiver)

    log("Created a new empty collection and published a reference")
  }
}
