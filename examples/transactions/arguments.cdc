
// This transaction creates an empty NFT Collection in the signer's account
transaction(test:String) {
  prepare(acct: &Account) {
    log(acct)
    log(test)

 }
}
