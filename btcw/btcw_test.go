package btcw

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

func TestCreateNewWallet(t *testing.T) {
	// 새로운 지갑을 생성합니다. private key 와 public key 가 생성됐는지 확인합니다.
	wallet := CreateNewWallet()
	if wallet.PrivateKey.D == nil || wallet.PublicKey == nil {
		t.Errorf("createNewWallet failed to generate a wallet with non-nil keys")
	}

	// mainnet, testnet 주소가 생성됐는지 확인합니다.
	mainnetAddress := wallet.GetAddress(BITCOIN_MAINNET_VERSION)
	if mainnetAddress == "" && len(mainnetAddress) == 34 {
		t.Errorf("createNewWallet failed to generate a wallet with a non-empty address")
	}
	testnetAddress := wallet.GetAddress(BITCOIN_TESTNET_VERSION)
	if testnetAddress == "" && len(mainnetAddress) == 34 {
		t.Errorf("createNewWallet failed to generate a wallet with a non-empty address")
	}
}

func TestGetPrivateKeyFromString(t *testing.T) {
	// private key string 으로부터 private key 객체를 생성합니다.
	// 생성된 private key 객체의 private key string 값과 비교합니다.
	privKeyStr := "18e14a7b6a307f426a94f8114701e7c8e774e7f9a47e2c2035db29a206321725"
	privKey := GetECDSAPrivateKeyFromPrivateKeyString(privKeyStr)
	privKey2 := fmt.Sprintf("%x", privKey.D)
	if privKey2 != privKeyStr {
		t.Errorf("GetECDSAPrivateKeyFromPrivateKeyString failed to generate a private key from a string")
	}
}

func TestCreateWalletGetWalletFromPrivStr(t *testing.T) {
	// 지갑을 생성한 다음 private key 같은 지갑 객체를 다시 생성 같은지 확인
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
	// private key string 으로부터 지갑을 생성한 다음 mainnet, testnet 주소가 같은지 확인합니다.
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

func TestCreateNewAddress(t *testing.T) {
	newPrivKey, err := btcec.NewPrivateKey()
	if err != nil {
		t.Errorf("NewPrivateKey failed to generate a new private key")
	}

	newWif, err := btcutil.NewWIF(newPrivKey, &chaincfg.MainNetParams, true)
	if err != nil {
		t.Errorf("NewWIF failed to generate a new WIF")
	}

	fmt.Printf("wif : %s\n", newWif.String())

	serializedPubKey := newPrivKey.PubKey().SerializeUncompressed()

	addressPubKey, err := btcutil.NewAddressPubKey(serializedPubKey, &chaincfg.MainNetParams)
	if err != nil {
		t.Errorf("NewAddressPubKey failed to generate a new address")
	}

	fmt.Printf("address : %s\n", addressPubKey.EncodeAddress())
}
