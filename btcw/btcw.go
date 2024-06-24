package btcw

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
	"math/big"

	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

const BITCOIN_MAINNET_VERSION = 0x00
const BITCOIN_TESTNET_VERSION = 0x6f

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}

func CreateNewWallet() *Wallet {
	// 개인키와 공개키를 생성한다.
	// https://pkg.go.dev/crypto/ecdsa
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Panic(err) // 현재 함수를 즉시 멈춘다. 상위 함수로 전파되어 프로그램이 종료된다.
	}

	return GetWalletFromPrivateKey(privateKey)

	//publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)

	//return &Wallet{*privateKey, publicKey}
}

func GetWalletFromPrivateKey(privateKey *ecdsa.PrivateKey) *Wallet {
	publicKey := append(privateKey.PublicKey.X.Bytes(), privateKey.PublicKey.Y.Bytes()...)
	return &Wallet{*privateKey, publicKey}
}

func GetECDSAPrivateKeyFromPrivateKeyString(privateKey string) ecdsa.PrivateKey {
	var e ecdsa.PrivateKey
	e.D, _ = new(big.Int).SetString(privateKey, 16)
	e.PublicKey.Curve = elliptic.P256()
	e.PublicKey.X, e.PublicKey.Y = e.PublicKey.Curve.ScalarBaseMult(e.D.Bytes())
	return e
}

func GetWalletFromPrivateKeyString(privateKeyStr string) *Wallet {
	privateKey := GetECDSAPrivateKeyFromPrivateKeyString(privateKeyStr)
	return GetWalletFromPrivateKey(&privateKey)
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

func (w Wallet) GetAddress(bitcoinVersion byte) string {
	publicKeyHash := hashPublicKey(w.PublicKey)
	return base58.CheckEncode(publicKeyHash, bitcoinVersion)
}
