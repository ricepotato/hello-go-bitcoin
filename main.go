package main

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ricepotato/hello-go-bitcoin/btcw"
)

func main() {

	fmt.Println("Hello go bitcoin")

	address1PrivKey := "18e14a7b6a307f426a94f8114701e7c8e774e7f9a47e2c2035db29a206321725"
	address2PrivKey := "d010d7a9b9a57f30e38a700fb5e2e367531f950e7a548939d4cfc8d5efc867b8"

	wallet1 := btcw.GetWalletFromPrivateKeyString(address1PrivKey)
	wallet2 := btcw.GetWalletFromPrivateKeyString(address2PrivKey)

	// miruDdUTqQv9eXMPPwXL73b9iy4gv8KeuH
	wallet1TestnetAddress := wallet1.GetLegacyAddress(btcw.BITCOIN_TESTNET_VERSION)
	// mshprpAi3YVjsiLQNNBP1z535K1YK6NLfa
	wallet2TestnetAddress := wallet2.GetLegacyAddress(btcw.BITCOIN_TESTNET_VERSION)

	fmt.Printf("Your bitcoin testnet address : %s\n", wallet1TestnetAddress)
	fmt.Printf("Your bitcoin testnet address : %s\n", wallet2TestnetAddress)

	segwitAddress := wallet1.GetSegwitAddress(&chaincfg.MainNetParams)
	fmt.Printf("Your bitcoin segwit address : %s\n", segwitAddress)
}
