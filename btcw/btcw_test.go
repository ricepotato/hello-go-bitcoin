package btcw

import (
	"fmt"
	"testing"
)

func TestCreateNewWallet(t *testing.T) {
	wallet := CreateNewWallet()
	if wallet.PrivateKey.D == nil || wallet.PublicKey == nil {
		t.Errorf("createNewWallet failed to generate a wallet with non-nil keys")
	}
}

func TestGetPrivateKeyFromString(t *testing.T) {
	privKeyStr := "18e14a7b6a307f426a94f8114701e7c8e774e7f9a47e2c2035db29a206321725"
	privKey := GetECDSAPrivateKeyFromPrivateKeyString(privKeyStr)
	privKey2 := fmt.Sprintf("%x", privKey.D)
	if privKey2 != privKeyStr {
		t.Errorf("GetECDSAPrivateKeyFromPrivateKeyString failed to generate a private key from a string")
	}
}

func TestCreateWalletGetWalletFromPrivStr(t *testing.T) {
	wallet := CreateNewWallet()
	address1 := wallet.GetAddress(BITCOIN_MAINNET_VERSION)
	privKey := fmt.Sprintf("%x", wallet.PrivateKey.D)
	wallet2 := GetWalletFromPrivateKeyString(privKey)
	address2 := wallet2.GetAddress(BITCOIN_MAINNET_VERSION)

	if address1 != address2 {
		t.Errorf("GetWalletFromPrivateKeyString failed to generate a wallet from a private key string")
	}
}

func TestCreateWalletFromPrivKeyStr(t *testing.T) {
	privKeyStr := "d010d7a9b9a57f30e38a700fb5e2e367531f950e7a548939d4cfc8d5efc867b8"
	mainnetAddressExpected := "1DBsZm5jEX4V6brneoD1C4riDKQqPpnGmN"
	testnetAddressExpected := "mshprpAi3YVjsiLQNNBP1z535K1YK6NLfa"
	wallet := GetWalletFromPrivateKeyString(privKeyStr)
	mainnetAddress := wallet.GetAddress(BITCOIN_MAINNET_VERSION)
	testnetAddress := wallet.GetAddress(BITCOIN_TESTNET_VERSION)

	if mainnetAddress != mainnetAddressExpected {
		t.Errorf("GetWalletFromPrivateKeyString failed to generate a wallet from a private key string")
	}
	if testnetAddress != testnetAddressExpected {
		t.Errorf("GetWalletFromPrivateKeyString failed to generate a wallet from a private key string")
	}

}
