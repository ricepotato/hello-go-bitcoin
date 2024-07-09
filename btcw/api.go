package btcw

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type AddressEndpoint struct {
	Address            string  `json:"address"`
	TotalReceived      int     `json:"total_received"`
	TotalSent          int     `json:"total_sent"`
	Balance            int     `json:"balance"`
	UnconfirmedBalance int     `json:"unconfirmed_balance"`
	FinalBalance       int     `json:"final_balance"`
	NTx                int     `json:"n_tx"`
	UnconfirmedNTx     int     `json:"unconfirmed_n_tx"`
	FinalNTx           int     `json:"final_n_tx"`
	TxRefs             []TxRef `json:"txrefs"`
}

type TxRef struct {
	TxHash        string `json:"tx_hash"`
	BlockHeight   int    `json:"block_height"`
	TxInputN      int    `json:"tx_input_n"`
	TxOutputN     int    `json:"tx_output_n"`
	Value         int    `json:"value"`
	RefBalance    int    `json:"ref_balance"`
	Spent         bool   `json:"spent"`
	Confirmations int    `json:"confirmations"`
	Confirmed     string `json:"confirmed"`
	DoubleSpend   bool   `json:"double_spend"`
}

func GetAddressEndpoint(address string) (*AddressEndpoint, error) {
	/*
		response example
		{
			"address": "miruDdUTqQv9eXMPPwXL73b9iy4gv8KeuH",
			"total_received": 33115,
			"total_sent": 0,
			"balance": 33115,
			"unconfirmed_balance": 0,
			"final_balance": 33115,
			"n_tx": 2,
			"unconfirmed_n_tx": 0,
			"final_n_tx": 2,
			"txrefs": [
				{
					"tx_hash": "bc26416ce0facd6733b26f5322b21f834ec9206eea9525ffd44e6c1102810fad",
					"block_height": 2866858,
					"tx_input_n": -1,
					"tx_output_n": 1,
					"value": 17891,
					"ref_balance": 33115,
					"spent": false,
					"confirmations": 19,
					"confirmed": "2024-07-08T05:25:30Z",
					"double_spend": false
				},
				{
					"tx_hash": "ae2db75b9d7bfb9bc3c39e238e1479df3021b206039f67528024ab123990feb4",
					"block_height": 2822254,
					"tx_input_n": -1,
					"tx_output_n": 0,
					"value": 15224,
					"ref_balance": 15224,
					"spent": false,
					"confirmations": 44623,
					"confirmed": "2024-06-24T09:03:35Z",
					"double_spend": false
				}
			],
			"tx_url": "https://api.blockcypher.com/v1/btc/test3/txs/"
		}
	*/
	url := fmt.Sprintf("https://api.blockcypher.com/v1/btc/test3/addrs/%s?unspentOnly=true", address)
	resp, err := http.Get(url)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	if resp.StatusCode != 200 {
		bodyText := string(body)
		log.Panicf("http status error. status code: %d body: %s failed to get UTXOs", resp.StatusCode, bodyText)
	}

	var addressEndpoint AddressEndpoint
	err = json.Unmarshal(body, &addressEndpoint)
	if err != nil {
		log.Panic(err)
	}

	return &addressEndpoint, nil
}
