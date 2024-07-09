package btcw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type CurrentFeeResp struct {
	Result struct {
		FeeRate float64 `json:"feerate"`
		Blocks  int     `json:"blocks"`
	} `json:"result"`
	Error interface{} `json:"error"`
	ID    string      `json:"id"`
}

type SendRawTransactionResp struct {
	Result string      `json:"result"`
	Error  interface{} `json:"error"`
	ID     string      `json:"id"`
}

func SendRawTransaction(signedhex string) (string, error) {
	/*
		response example
		{
			"result": "0x0d12f3a4b5c6d7e8f9",
			"error": null,
			"id": "1"
		}
	*/
	rpcUrl := "https://bitcoin-testnet.g.allthatnode.com/archive/json_rpc/d9255f43f22848fea00de73650288453"
	maxFeeRate := 0.10
	bodyData := fmt.Sprintf("{\"jsonrpc\": \"1.0\", \"id\": \"1\", \"method\": \"sendrawtransaction\", \"params\": [\"%s\", %f]}", signedhex, maxFeeRate)
	resp, err := http.Post(rpcUrl, "plain/text", bytes.NewBuffer([]byte(bodyData)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	if resp.StatusCode != 200 {
		bodyText := string(body)
		return "", fmt.Errorf("failed to send raw transaction. status code: %d text: %s", resp.StatusCode, bodyText)
	}

	var respData SendRawTransactionResp
	err = json.Unmarshal(body, &respData)
	if err != nil {
		log.Fatal(err)
		return "0", err
	}

	return respData.Result, nil
}

// GetCurrentFee gets the current fee in bitcoin
func GetCurrentFee() (float64, error) {
	/*
		response example
		{
			"result": {
				"feerate": 0.00024190,
				"blocks": 100
			},
			"error": null,
			"id": "1"
		}
	*/

	rpcUrl := "https://bitcoin-testnet.g.allthatnode.com/archive/json_rpc/d9255f43f22848fea00de73650288453"
	bodyData := "{\"jsonrpc\": \"1.0\", \"id\": \"1\", \"method\": \"estimatesmartfee\", \"params\": [1000]}"
	resp, err := http.Post(rpcUrl, "plain/text", bytes.NewBuffer([]byte(bodyData)))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var respData CurrentFeeResp

	err = json.Unmarshal(body, &respData)
	if err != nil {
		return 0, err
	}

	return respData.Result.FeeRate, nil
}
