package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

// https://hou27.tistory.com/entry/Go%EB%A1%9C-%EB%A7%8C%EB%93%9C%EB%8A%94-%EB%B8%94%EB%A1%9D%EC%B2%B4%EC%9D%B8-part-5-Wallet

const BITCOIN_MAINET_VERSION = 0x00
const BITCOIN_TESTNET_VERSION = 0x6f

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}

func createNewWallet() *Wallet {
	// 함수는 개인키와 공개키를 생성한다.
	// https://pkg.go.dev/crypto/ecdsa
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Panic(err) // 현재 함수를 즉시 멈춘다. 상위 함수로 전파되어 프로그램이 종료된다.
	}

	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	return &Wallet{*privateKey, publicKey}
}

func getWalletFromPrivateKeyString(privateKeyStr string) *Wallet {
	// Convert the private key from a hex string to a byte slice
	privateKeyBytes, err := hex.DecodeString(privateKeyStr)
	if err != nil {
		log.Panic(err)
	}

	// Parse the private key bytes into an ecdsa.PrivateKey object
	privateKey, err := x509.ParseECPrivateKey(privateKeyBytes)
	if err != nil {
		log.Panic(err)
	}

	// Generate the public key from the private key
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	// Return a new Wallet instance with the private and public keys
	return &Wallet{*privateKey, publicKey}
}

func hashPublicKey(publicKey []byte) []byte {
	// Public key를 먼저 SHA-256으로 해싱하고, 이 결과를 RIPEMD-160으로 또다시 해싱
	publicSHA256 := sha256.Sum256(publicKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}

	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)
	return publicRIPEMD160

}

func (w Wallet) getAddress(bitcoinVersion byte) string {
	publicKeyHash := hashPublicKey(w.PublicKey)
	return base58.CheckEncode(publicKeyHash, bitcoinVersion)
}

func main() {
	fmt.Println("Hello go bitcoin")

	mainnetWallet := createNewWallet()
	address := mainnetWallet.getAddress(BITCOIN_MAINET_VERSION)
	fmt.Printf("Your private key : %x\n", mainnetWallet.PrivateKey.D)
	fmt.Printf("Your public key : %x\n", mainnetWallet.PublicKey)
	fmt.Printf("Your bitcoin address : %s\n", address)

	testnetWallet := createNewWallet()
	address = testnetWallet.getAddress(BITCOIN_TESTNET_VERSION)
	fmt.Printf("Your private key : %x\n", testnetWallet.PrivateKey.D)
	fmt.Printf("Your public key : %x\n", testnetWallet.PublicKey)
	fmt.Printf("Your bitcoin address : %s\n", address)
}
