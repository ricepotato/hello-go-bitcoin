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
	mainnetAddress := wallet.GetLegacyAddress(BITCOIN_MAINNET_VERSION)
	if mainnetAddress == "" && len(mainnetAddress) == 34 {
		t.Errorf("createNewWallet failed to generate a wallet with a non-empty address")
	}
	testnetAddress := wallet.GetLegacyAddress(BITCOIN_TESTNET_VERSION)
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
	address1 := wallet.GetLegacyAddress(BITCOIN_MAINNET_VERSION)
	privKey := fmt.Sprintf("%x", wallet.PrivateKey.D)
	wallet2 := GetWalletFromPrivateKeyString(privKey)
	address2 := wallet2.GetLegacyAddress(BITCOIN_MAINNET_VERSION)

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
	mainnetAddress := wallet.GetLegacyAddress(BITCOIN_MAINNET_VERSION)
	testnetAddress := wallet.GetLegacyAddress(BITCOIN_TESTNET_VERSION)

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

func TestGetBalance(t *testing.T) {
	// 지갑 주소로부터 잔고를 조회합니다.
	address := "miruDdUTqQv9eXMPPwXL73b9iy4gv8KeuH"
	balance, err := GetBalance(BITCOIN_TESTNET_VERSION, address)
	if err != nil {
		t.Errorf("GetBalance failed to get the balance")
	}
	fmt.Printf("balance : %d\n", &balance.FinalBalance)
}

func TestAddressFromWif(t *testing.T) {
	// wfi 로부터 새로운 주소를 생성 compress pubkey, uncompress pubkey 를 생성하여 비교합니다.
	// comp pubkey 와 uncomp pubkey 는 같지 않으며 compressedWif.CompressPubKey 값에 따라 맞는 pubkey 를 사용해야 함
	net := &chaincfg.MainNetParams
	//net := &chaincfg.TestNet3Params
	privateKey, _ := btcec.NewPrivateKey()
	compressedWif, _ := btcutil.NewWIF(privateKey, net, true)
	uncompressedWif, _ := btcutil.NewWIF(privateKey, net, false)

	if !compressedWif.CompressPubKey {
		t.Errorf("CompressPubKey should be true")
	}

	fmt.Printf("compressed wif : %s\n", compressedWif.String())
	fmt.Printf("uncompressed wif : %s\n", uncompressedWif.String())

	addressPubKeyComp, _ := btcutil.NewAddressPubKey(compressedWif.PrivKey.PubKey().SerializeCompressed(), net)
	addressPubKeyUncomp, _ := btcutil.NewAddressPubKey(uncompressedWif.PrivKey.PubKey().SerializeUncompressed(), net)
	encodedAddressPubKeyComp := addressPubKeyComp.EncodeAddress()
	encodedAddressPubKeyUncomp := addressPubKeyUncomp.EncodeAddress()

	fmt.Printf("address comp : %s\n", encodedAddressPubKeyComp)
	fmt.Printf("address uncomp: %s\n", encodedAddressPubKeyUncomp)

	if encodedAddressPubKeyComp == encodedAddressPubKeyUncomp {
		t.Errorf("encodedAddressPubKeyComp should not be equal to encodedAddressPubKeyUncomp")
	}
}

func TestWifStringFromLegacyAddressMainnet(t *testing.T) {
	net := &chaincfg.MainNetParams
	wif1 := "L1KJbwoL9R7oj1yjDmEnT3tWURcdXd3eEgSdJ8PWfrJb8YMWwwD4"
	compressedWif, _ := btcutil.DecodeWIF(wif1)

	addressPubKey, _ := btcutil.NewAddressPubKey(compressedWif.PrivKey.PubKey().SerializeCompressed(), net)
	encodedAddressPubKey := addressPubKey.EncodeAddress()
	fmt.Printf("address : %s\n", encodedAddressPubKey)

	legacyAddress := "1JTYmdzMnXHB3MDfsqWidkk422udQYLEJV"

	if encodedAddressPubKey != legacyAddress {
		t.Errorf("GetWalletFromPrivateKeyString failed to generate a wallet from a private key string")
	}
}

func TestWifStringFromNativeSegwitAddressTestnet(t *testing.T) {
	// P2WPKH 방식 bc1 으로 시작하는 주소 Testnet
	wif1 := "L3F6LJgS4RJm1SJpFcQuZVFBoJ1veBowNm5Vwz8sLb4RWtPFpjPH"
	compressedWif, _ := btcutil.DecodeWIF(wif1)
	if !compressedWif.CompressPubKey {
		t.Errorf("CompressPubKey should be true")
	}

	hashedPubKey := hashPublicKey(compressedWif.PrivKey.PubKey().SerializeCompressed())
	addressPubKey, _ := btcutil.NewAddressWitnessPubKeyHash(hashedPubKey, &chaincfg.TestNet3Params)
	encodedAddressPubKey := addressPubKey.EncodeAddress()
	expectedAddress := "tb1qz40mujlemrru7t8t3yn3u5v3e9htmu5kektgme"

	if encodedAddressPubKey != expectedAddress {
		t.Errorf("GetWalletFromPrivateKeyString failed to generate a wallet from a private key string")
	}
}

func TestWifStringFromNativeSegwitAddressMainnet(t *testing.T) {
	// P2WPKH 방식 bc1 으로 시작하는 주소 Mainnet
	wif1 := "L3F6LJgS4RJm1SJpFcQuZVFBoJ1veBowNm5Vwz8sLb4RWtPFpjPH"
	compressedWif, _ := btcutil.DecodeWIF(wif1)
	if !compressedWif.CompressPubKey {
		t.Errorf("CompressPubKey should be true")
	}

	hashedPubKey := hashPublicKey(compressedWif.PrivKey.PubKey().SerializeCompressed())
	addressPubKey, _ := btcutil.NewAddressWitnessPubKeyHash(hashedPubKey, &chaincfg.MainNetParams)
	encodedAddressPubKey := addressPubKey.EncodeAddress()
	expectedAddress := "bc1qz40mujlemrru7t8t3yn3u5v3e9htmu5knssmq2"

	if encodedAddressPubKey != expectedAddress {
		t.Errorf("GetWalletFromPrivateKeyString failed to generate a wallet from a private key string")
	}
}

func TestWif1StringFromKLegacyAddress(t *testing.T) {
	// wif 로부터 testnet legacy address 를 생성함
	net := &chaincfg.TestNet3Params
	wif1 := "cUsfNynj7UsBvjLPeb6TjnxA4SFThiZwXr7Az5TxJGDWUFGQbbZv"
	compressedWif, _ := btcutil.DecodeWIF(wif1)
	if !compressedWif.CompressPubKey {
		t.Errorf("CompressPubKey should be true")
	}

	addressPubKey, _ := btcutil.NewAddressPubKey(compressedWif.PrivKey.PubKey().SerializeCompressed(), net)
	encodedAddressPubKey := addressPubKey.EncodeAddress()
	expectedAddress := "myQCR5hm5R6NWoKn4o5MSLGiLTrKdk2AbD"

	compressedWif.PrivKey.Serialize()

	if encodedAddressPubKey != expectedAddress {
		t.Errorf("GetWalletFromPrivateKeyString failed to generate a wallet from a private key string")
	}
}

func TestWif2StringFromLegacyAddress(t *testing.T) {
	// wif 로부터 testnet legacy address 를 생성함
	net := &chaincfg.TestNet3Params
	wif1 := "cQBQdEv2uqp3pTZGCMqNCoChweHrvFphZBmocyFanpgqvey7BrLP"
	compressedWif, _ := btcutil.DecodeWIF(wif1)
	if !compressedWif.CompressPubKey {
		t.Errorf("CompressPubKey should be true")
	}

	addressPubKey, _ := btcutil.NewAddressPubKey(compressedWif.PrivKey.PubKey().SerializeCompressed(), net)
	encodedAddressPubKey := addressPubKey.EncodeAddress()
	expectedAddress := "mkgMpQSbsJMBffCbVMUk6KryCuFGAtnFo9"

	if encodedAddressPubKey != expectedAddress {
		t.Errorf("GetWalletFromPrivateKeyString failed to generate a wallet from a private key string")
	}
}
