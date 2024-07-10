package btcw

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"sort"
	"time"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

// UTXO ...
type UTXO struct {
	Hash      string
	TxIndex   int
	Amount    *big.Int
	Spendable bool
	PKScript  []byte
}

func CreateTransferTransaction(fromAddress string, toAddress string, privKey []byte, amountSatoshi int64) (string, error) {
	unspentTXOs, err := GetUTXO(fromAddress)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	// if fromAddress UTXO is empty, return
	if len(unspentTXOs) == 0 {
		err := errors.New("fromAddress UTXO is empty")
		log.Fatal(err)
		return "", err
	}

	chainParams := &chaincfg.TestNet3Params
	amountToSend := big.NewInt(amountSatoshi) // amount to send in satoshis (0.01 btc)
	feeRate, err := GetCurrentFeeRate()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	unspentTXOs, UTXOsAmount, err := marshalUTXOs(unspentTXOs, amountToSend, feeRate)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	tx := wire.NewMsgTx(wire.TxVersion)
	var sourceUTXOs []*UTXO

	for idx := range unspentTXOs {
		hashStr := unspentTXOs[idx].Hash

		sourceUTXOHash, err := chainhash.NewHashFromStr(hashStr)
		if err != nil {
			log.Fatal(err)
		}

		sourceUTXOIndex := uint32(unspentTXOs[idx].TxIndex)
		sourceUTXO := wire.NewOutPoint(sourceUTXOHash, sourceUTXOIndex)
		sourceUTXOs = append(sourceUTXOs, unspentTXOs[idx])
		sourceTxIn := wire.NewTxIn(sourceUTXO, nil, nil)

		tx.AddTxIn(sourceTxIn)
	}

	// calculate fees
	txByteSize := big.NewInt(int64(len(tx.TxIn)*180 + len(tx.TxOut)*34 + 10 + len(tx.TxIn)))
	totalFee := new(big.Int).Mul(feeRate, txByteSize)
	log.Printf("total fee: %s", totalFee)

	// calculate the change
	change := new(big.Int).Set(UTXOsAmount)
	change = new(big.Int).Sub(change, amountToSend)
	change = new(big.Int).Sub(change, totalFee)
	if change.Cmp(big.NewInt(0)) == -1 {
		log.Fatal(err)
	}

	// create the tx outs
	destAddress, err := btcutil.DecodeAddress(toAddress, chainParams)
	if err != nil {
		log.Fatal(err)
	}

	destScript, err := txscript.PayToAddrScript(destAddress)
	if err != nil {
		log.Fatal(err)
	}

	// tx out to send btc to user
	destOutput := wire.NewTxOut(amountToSend.Int64(), destScript)
	tx.AddTxOut(destOutput)

	// our change address
	changeSendToAddress, err := btcutil.DecodeAddress(fromAddress, chainParams)
	if err != nil {
		log.Fatal(err)
	}

	changeSendToScript, err := txscript.PayToAddrScript(changeSendToAddress)
	if err != nil {
		log.Fatal(err)
	}

	// tx out to send change back to us
	changeOutput := wire.NewTxOut(change.Int64(), changeSendToScript)
	tx.AddTxOut(changeOutput)

	sourceAddress, err := btcutil.DecodeAddress(fromAddress, chainParams)
	if err != nil {
		log.Fatal(err)
	}

	sourcePkScript, err := txscript.PayToAddrScript(sourceAddress)
	if err != nil {
		log.Fatal(err)
	}

	pKey, _ := btcec.PrivKeyFromBytes(privKey)

	for i := range sourceUTXOs {
		outputFetcher := txscript.NewCannedPrevOutputFetcher(sourcePkScript, sourceUTXOs[i].Amount.Int64())
		txSigHashes := txscript.NewTxSigHashes(tx, outputFetcher)
		txWitness, err := txscript.WitnessSignature(tx, txSigHashes, i, sourceUTXOs[i].Amount.Int64(), sourcePkScript, txscript.SigHashAll, pKey, true)
		if err != nil {
			log.Fatalf("could not generate pubSig; err: %v", err)
		}
		tx.TxIn[i].Witness = txWitness
	}

	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	tx.Serialize(buf)

	fmt.Printf("Redeem Tx: %v\n", hex.EncodeToString(buf.Bytes()))

	t := hex.EncodeToString(buf.Bytes())
	return t, nil
}

func TransferCoin(fromAddress string, toAddress string, privKey []byte, amountSatoshi int64) (string, error) {
	log.Printf("%s->%s, CreateTransferTransaction amountSatoshi: %d", fromAddress, toAddress, amountSatoshi)
	signedHex, err := CreateTransferTransaction(fromAddress, toAddress, privKey, amountSatoshi)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	log.Printf("%s->%s SendRawTransaction", fromAddress, toAddress)
	txHash, err := SendRawTransaction(signedHex)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	log.Printf("%s->%s txHash: %s", fromAddress, toAddress, txHash)
	return txHash, nil
}

func marshalUTXOs(utxos []*UTXO, amount, feeRate *big.Int) ([]*UTXO, *big.Int, error) {
	// same strategy as bitcoin core
	// from: https://blog.lopp.net/the-challenges-of-optimizing-unspent-output-selection/
	// 1. sort the UTXOs from smallest to largest amounts
	sort.Slice(utxos, func(i, j int) bool {
		return utxos[i].Amount.Cmp(utxos[j].Amount) == -1
	})

	// 2. search for exact match
	for idx := range utxos {
		exactTxSize := calculateTotalTxBytes(1, 2)
		totalFee := new(big.Int).Mul(feeRate, big.NewInt(int64(exactTxSize)))
		totalTxAmount := new(big.Int).Add(totalFee, amount)

		switch utxos[idx].Amount.Cmp(totalTxAmount) {
		case 0:
			var resp []*UTXO
			resp = append(resp, utxos[idx])
			// TODO: store these in the DB to be sure they aren't being claimed??
			return resp, sumUTXOs(resp), nil

		case 1:
			break
		}
	}

	// 3. calculate the sum of all UTXOs smaller than amount
	sumSmall := big.NewInt(0)
	var sumSmallUTXOs []*UTXO
	for idx := range utxos {
		switch utxos[idx].Amount.Cmp(amount) {
		case -1:
			_ = sumSmall.Add(sumSmall, utxos[idx].Amount)
			sumSmallUTXOs = append(sumSmallUTXOs, utxos[idx])

		default:
			break
		}
	}

	exactTxSize := calculateTotalTxBytes(len(sumSmallUTXOs), 2)
	totalFee := new(big.Int).Mul(feeRate, big.NewInt(int64(exactTxSize)))
	totalTxAmount := new(big.Int).Add(totalFee, amount)

	log.Printf("exactTxSize: %d", totalTxAmount)
	log.Printf("totalFee: %d", totalFee)
	log.Printf("totalTxAmount: %d", totalTxAmount)

	switch sumSmall.Cmp(totalTxAmount) {
	case 0:
		return sumSmallUTXOs, sumUTXOs(sumSmallUTXOs), nil

	case -1:
		for idx := range utxos {
			exactTxSize := calculateTotalTxBytes(1, 2)
			totalFee := new(big.Int).Mul(feeRate, big.NewInt(int64(exactTxSize)))
			totalTxAmount := new(big.Int).Add(totalFee, amount)
			if utxos[idx].Amount.Cmp(totalTxAmount) == 1 {
				var resp []*UTXO
				resp = append(resp, utxos[idx])
				return resp, sumUTXOs(resp), nil
			}
		}

		// should reach here if not enought UXOs
		log.Fatal("not enough UTXOs to meet target amount")

	case 1:
		return roundRobinSelectUTXOs(sumSmallUTXOs, amount, feeRate)

	default:
		log.Fatal("unknown comparison")
	}

	return nil, nil, nil
}

func roundRobinSelectUTXOs(utxos []*UTXO, amount, feeRate *big.Int) ([]*UTXO, *big.Int, error) {
	var possibilities [][]*UTXO
	lenInput := len(utxos)
	log.Printf("round robin select; lenInput: %v", lenInput)
	if lenInput == 0 {
		log.Fatal("expected utxos size to be greater than 0")
	}

	for i := 0; i < 1000; i++ {
		selectedIdxs := make(map[int]bool)
		var sum *big.Int
		var possibility []*UTXO
		for {
			for {
				rand.Seed(time.Now().Unix())
				tmp := 0
				if lenInput > 1 {
					tmp = rand.Intn(lenInput - 1)
				}

				if !selectedIdxs[tmp] {
					selectedIdxs[tmp] = true
					_ = sum.Add(sum, utxos[tmp].Amount)
					possibility = append(possibility, utxos[tmp])

					break
				}
			}

			exactTxSize := calculateTotalTxBytes(len(possibility), 2)
			totalFee := new(big.Int).Mul(feeRate, big.NewInt(int64(exactTxSize)))
			totalTxAmount := new(big.Int).Add(totalFee, amount)

			if sum.Cmp(totalTxAmount) == 0 {
				return possibility, sum, nil
			}

			if sum.Cmp(totalTxAmount) == 1 {
				possibilities = append(possibilities, possibility)
				break
			}
		}
	}

	if len(possibilities) < 1 {
		return nil, nil, errors.New("no possible utxo combos")
	}

	smallestLen := len(possibilities[0])
	smallestIdx := 0

	for idx := 1; idx < len(possibilities); idx++ {
		l := len(possibilities[idx])
		if l < smallestLen {
			smallestLen = l
			smallestIdx = idx
		}
	}

	return possibilities[smallestIdx], sumUTXOs(possibilities[smallestIdx]), nil
}

func sumUTXOs(utxos []*UTXO) *big.Int {
	sum := big.NewInt(0)
	for idx := range utxos {
		sum = sum.Add(sum, utxos[idx].Amount)
	}

	return sum
}

// https://bitcoin.stackexchange.com/questions/1195/how-to-calculate-transaction-size-before-sending-legacy-non-segwit-p2pkh-p2sh
func calculateTotalTxBytes(txInLength, txOutLength int) int {
	return txInLength*180 + txOutLength*34 + 10 + txInLength
}

// GetCurrentFeeRate gets the current fee in satoshis per kb
func GetCurrentFeeRate() (*big.Int, error) {
	fee, err := GetCurrentFee()
	if err != nil {
		return nil, err
	}

	// convert to satoshis to bytes
	// feeRate := big.NewInt(int64(msg.Result * 1.0E8))
	// convert to satoshis to kb
	feeRate := big.NewInt(int64(fee * 1.0e5))

	fmt.Printf("fee rate: %s\n", feeRate)

	return feeRate, nil
}

func GetUTXO(address string) ([]*UTXO, error) {
	addressEndpoint, err := GetAddressEndpoint(address)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	// filter spent false from TxRefs
	var utxos []*UTXO
	for _, txRef := range addressEndpoint.TxRefs {
		if !txRef.Spent {
			utxos = append(utxos, &UTXO{
				Hash:      txRef.TxHash,
				TxIndex:   txRef.TxOutputN,
				Amount:    big.NewInt(int64(txRef.Value)),
				Spendable: true,
				PKScript:  nil,
			})
			log.Printf("UTXO address: %s, hash: %s, amount: %d", address, txRef.TxHash, txRef.Value)
		}

	}
	return utxos, nil
}
