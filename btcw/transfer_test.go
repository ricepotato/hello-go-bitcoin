package btcw

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

func TestAddressEndpoint(t *testing.T) {
	address := "miruDdUTqQv9eXMPPwXL73b9iy4gv8KeuH"
	result, err := GetAddressEndpoint(address)
	if err != nil {
		t.Error(err)
	}

	// balance > 0
	if result.Balance <= 0 {
		t.Errorf("Balance should be greater than 0")
	}
}

func TestWifPubkey(t *testing.T) {
	chainParams := &chaincfg.TestNet3Params
	privWif := "cS5LWK2aUKgP9LmvViG3m9HkfwjaEJpGVbrFHuGZKvW2ae3W9aUe"

	decodedWif, _ := btcutil.DecodeWIF(privWif)
	addressPubKey, _ := btcutil.NewAddressPubKey(decodedWif.PrivKey.PubKey().SerializeUncompressed(), chainParams)
	pubkeyEncoded := addressPubKey.EncodeAddress()
	sourceAddress, _ := btcutil.DecodeAddress(pubkeyEncoded, chainParams)
	fromWalletPublicAddress := "mgjHgKi1g6qLFBM1gQwuMjjVBGMJdrs9pP"

	if fromWalletPublicAddress != sourceAddress.EncodeAddress() {
		t.Errorf("Address should be equal")
	}
}

func TestCreateTransferTransaction(t *testing.T) {
	fromWifString := "cQBQdEv2uqp3pTZGCMqNCoChweHrvFphZBmocyFanpgqvey7BrLP"
	wif, _ := btcutil.DecodeWIF(fromWifString)
	fromAddress := "myQCR5hm5R6NWoKn4o5MSLGiLTrKdk2AbD"
	toAddress := "mkgMpQSbsJMBffCbVMUk6KryCuFGAtnFo9"
	var amountSatoshi int64 = 100

	signedHex, err := CreateTransferTransaction(fromAddress, toAddress, wif.PrivKey.Serialize(), amountSatoshi)
	if err != nil {
		t.Error(err)
	}

	if signedHex == "" {
		t.Errorf("Signed hex should not be empty")
	}

	fmt.Printf("signedHex : %s", signedHex)
}

func HexToBytes(hexStr string) ([]byte, error) {
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func TestWifFromWallet(t *testing.T) {
	privateKey := "18e14a7b6a307f426a94f8114701e7c8e774e7f9a47e2c2035db29a206321725"
	privKeyBytes, _ := HexToBytes(privateKey)
	privKey, _ := btcec.PrivKeyFromBytes(privKeyBytes)
	testnetWif, _ := btcutil.NewWIF(privKey, &chaincfg.TestNet3Params, true)

	addressPubKey, _ := btcutil.NewAddressPubKey(testnetWif.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.TestNet3Params)
	addressStr := addressPubKey.EncodeAddress()

	fmt.Printf("addressStr : %s\n", addressStr)
}

func TestTransferCoinFailed(t *testing.T) {
	privateKey := "18e14a7b6a307f426a94f8114701e7c8e774e7f9a47e2c2035db29a206321725"
	//wallet := GetWalletFromPrivateKeyString(privateKey)
	fromAddress := "miruDdUTqQv9eXMPPwXL73b9iy4gv8KeuH"
	toAddress := "mshprpAi3YVjsiLQNNBP1z535K1YK6NLfa"
	var amountSatoshi int64 = 1000

	privKeyBytes, _ := HexToBytes(privateKey)

	txHash, err := TransferCoin(fromAddress, toAddress, privKeyBytes, amountSatoshi)
	if err != nil {
		t.Error(err)
	}

	if txHash == "" {
		t.Errorf("Tx hash should not be empty")
	}

	fmt.Printf("txHash : %s", txHash)
}

func TestTransferCoinLegacy(t *testing.T) {
	fromWifString := "cUsfNynj7UsBvjLPeb6TjnxA4SFThiZwXr7Az5TxJGDWUFGQbbZv"
	wif, _ := btcutil.DecodeWIF(fromWifString)
	fromAddress := "myQCR5hm5R6NWoKn4o5MSLGiLTrKdk2AbD"
	toAddress := "mkgMpQSbsJMBffCbVMUk6KryCuFGAtnFo9"
	var amountSatoshi int64 = 100

	txHash, err := TransferCoin(fromAddress, toAddress, wif.PrivKey.Serialize(), amountSatoshi)
	if err != nil {
		t.Error(err)
	}

	if txHash == "" {
		t.Errorf("Tx hash should not be empty")
	}

	fmt.Printf("txHash : %s", txHash)
}

func TestTransferCoinNativeSegwit(t *testing.T) {
	fromWifString := "L3F6LJgS4RJm1SJpFcQuZVFBoJ1veBowNm5Vwz8sLb4RWtPFpjPH"
	wif, _ := btcutil.DecodeWIF(fromWifString)
	fromAddress := "tb1qz40mujlemrru7t8t3yn3u5v3e9htmu5kektgme"
	toAddress := "tb1qf9k7gahvkcngazw3hwaclh6dqmc0g38ke3295q"
	var amountSatoshi int64 = 1000

	txHash, err := TransferCoin(fromAddress, toAddress, wif.PrivKey.Serialize(), amountSatoshi)
	if err != nil {
		t.Error(err)
	}

	if txHash == "" {
		t.Errorf("Tx hash should not be empty")
	}

	fmt.Printf("txHash : %s", txHash)
}
