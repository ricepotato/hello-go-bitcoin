package btcw

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"math/big"
	"net/http"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/chaincfg"
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

func (w Wallet) GetLegacyAddress(bitcoinVersion byte) string {
	publicKeyHash := hashPublicKey(w.PublicKey)
	return base58.CheckEncode(publicKeyHash, bitcoinVersion)
}

func (w Wallet) GetLegacyAddress2(chaincfgParams *chaincfg.Params) string {
	publicKeyHash := hashPublicKey(w.PublicKey)
	rtn, _ := btcutil.NewAddressPubKeyHash(publicKeyHash, chaincfgParams)
	return rtn.EncodeAddress()
}

func (w Wallet) GetSegwitAddress(chaincfgParams *chaincfg.Params) string {
	publicKeyHash := hashPublicKey(w.PublicKey)
	rtn, _ := btcutil.NewAddressWitnessPubKeyHash(publicKeyHash, chaincfgParams)
	return rtn.EncodeAddress()
}

func (w Wallet) PrivateKeyToBytes() []byte {
	return w.PrivateKey.D.Bytes()
}

func GetLegacyAddressFromPubKeyString(serializedPubKey string, net *chaincfg.Params) string {
	// P2PKH 주소를 public key 로부터 생성
	serializedPubKeyBytes, err := hex.DecodeString(serializedPubKey)
	if err != nil {
		log.Panic(err)
	}
	return GetLegacyAddressFromPubKeyBytes(serializedPubKeyBytes, net)
}

func GetLegacyAddressFromPubKeyBytes(serializedPubKey []byte, net *chaincfg.Params) string {
	// P2PKH 주소를 public key 로부터 생성
	addressPubKey, err := btcutil.NewAddressPubKey(serializedPubKey, net)
	if err != nil {
		log.Panic(err)
	}
	return addressPubKey.EncodeAddress()
}

func GetSegwitAddressFromPubKeyBytes(serializedPubKey []byte, net *chaincfg.Params) string {
	// 네이티브 세그윗 주소(Bech32 주소)
	// segwit mainnet p2wpkh address 를 pubkey 로부터 생성
	// bc1 으로 시작하고 소문자로만 구성됨
	pkHash := hashPublicKey(serializedPubKey)
	addressPubKey, err := btcutil.NewAddressWitnessPubKeyHash(pkHash, net)
	if err != nil {
		log.Panic(err)
	}
	return addressPubKey.EncodeAddress()
}

func GetSegwitAddressFromPubKeyString(serializedPubKey string, net *chaincfg.Params) string {
	// 네이티브 세그윗 주소(Bech32 주소)
	serializedPubKeyBytes, err := hex.DecodeString(serializedPubKey)
	if err != nil {
		log.Panic(err)
	}
	return GetSegwitAddressFromPubKeyBytes(serializedPubKeyBytes, net)
}

type GetBalanceResponse struct {
	Address            string `json:"address"`
	TotalReceived      int    `json:"total_received"`
	TotalSent          int    `json:"total_sent"`
	Balance            int    `json:"balance"`
	UnconfirmedBalance int    `json:"unconfirmed_balance"`
	FinalBalance       int    `json:"final_balance"`
	NTx                int    `json:"n_tx"`
	UnconfirmedNTx     int    `json:"unconfirmed_n_tx"`
	FinalNTx           int    `json:"final_n_tx"`
}

func GetBalance(network int, address string) (*GetBalanceResponse, error) {
	/*
	   GetBalanceResponse example
	   {
	   "address": "1DEP8i3QJCsomS4BSMY2RpU1upv62aGvhD",
	   "total_received": 4433416,
	   "total_sent": 0,
	   "balance": 4433416,
	   "unconfirmed_balance": 0,
	   "final_balance": 4433416,
	   "n_tx": 7,
	   "unconfirmed_n_tx": 0,
	   "final_n_tx": 7
	   }
	*/
	// https://www.blockcypher.com/dev/?go#address-balance-endpoint
	baseUrl := "https://api.blockcypher.com/v1/btc/"
	if network == BITCOIN_MAINNET_VERSION {
		baseUrl += "main/"
	} else if network == BITCOIN_TESTNET_VERSION {
		baseUrl += "test3/"
	}
	baseUrl += "addrs/" + address + "/balance"

	res, err := http.Get(baseUrl)
	if err != nil {
		log.Panic(err)
	}

	if res.StatusCode != 200 {
		return nil, errors.New("failed to get balance")
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Panic(err)
	}

	var result GetBalanceResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Panic(err)
	}

	return &result, nil

}
