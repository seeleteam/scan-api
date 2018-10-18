/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package rpc

import (
	"encoding/json"
	"math/big"
	"testing"
)

func TestGetPeerInfos(t *testing.T) {
	peersJSON := `[
					{
						"caps": [
							"lightSeele/1",
							"seele/1"
						],
						"id": "0x0ea2a45ab5a909c309439b0e004c61b7b2a3e831",
						"network": {
							"localAddress": "127.0.0.1:55239",
							"remoteAddress": "127.0.0.1:8058"
						},
						"protocols": {
							"lightSeele": "handshake",
							"seele": "handshake"
						},
						"shard": 2
					},
					{
						"caps": [
							"lightSeele/1",
							"seele/1"
						],
						"id": "0x0da2a45ab5a909c309439b0e004c61b7b2a3e831",
						"network": {
							"localAddress": "127.0.0.1:55239",
							"remoteAddress": "127.0.0.1:8058"
						},
						"protocols": {
							"lightSeele": "handshake",
							"seele": "handshake"
						},
						"shard": 1
					}
				]`

	peersMp := make([]map[string]interface{}, 0)

	peers := getPeerInfos(peersMp)
	if len(peers) != 0 {
		t.Fatal("getPeerInfos get wrong peers number")
	}

	json.Unmarshal([]byte(peersJSON), &peersMp)

	peers = getPeerInfos(peersMp)
	if len(peers) != len(peersMp) {
		t.Fatal("parseToPeerInfos failed")
	}
	for _, info := range peers {
		if info.ID == "0x0ea2a45ab5a909c309439b0e004c61b7b2a3e831" {
			for _, cap := range info.Caps {
				if cap != "lightSeele/1" && cap != "seele/1" {
					t.Fatal("getPeerInfos get wrong data")
				}
			}
			if info.ShardNumber != 2 ||
				info.LocalAddress != "127.0.0.1:55239" ||
				info.RemoteAddress != "127.0.0.1:8058" {
				t.Fatal("getPeerInfos get wrong data")
			}
		} else if info.ID == "0x0da2a45ab5a909c309439b0e004c61b7b2a3e831" {
			for _, cap := range info.Caps {
				if cap != "lightSeele/1" && cap != "seele/1" {
					t.Fatal("getPeerInfos get wrong data")
				}
			}
			if info.ShardNumber != 1 ||
				info.LocalAddress != "127.0.0.1:55239" ||
				info.RemoteAddress != "127.0.0.1:8058" {
				t.Fatal("getPeerInfos get wrong data")
			}
		}
	}

}

func TestGetBalance(t *testing.T) {
	balanceJSON := `{
						"Account": "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21",
						"Balance": 261899990000
					}`

	balanceMp := make(map[string]interface{})
	json.Unmarshal([]byte(balanceJSON), &balanceMp)

	balance, err := getBalance(balanceMp, "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21")
	if err != nil {
		t.Fatal("parseToBalance failed")
	}
	if balance != 261899990000 {
		t.Fatal("parseToBalance failed")
	}
}

func TestGetReceiptByTxHash(t *testing.T) {
	receiptJSON := `{
						"contract": "0x",
						"failed": false,
						"poststate": "0xdd0b0fc6605bbb2e76b8c22ccd466ea5eaa1a80e4860fbdf971be58ded3d782b",
						"result": "0x",
						"totalFee": 21000,
						"txhash": "0xbd2ca4f9869c714e589ad6a3b16731c8cb066de40d0e27e220cc1e014577baff",
						"usedGas": 21000
					}`

	receiptMp := make(map[string]interface{})
	json.Unmarshal([]byte(receiptJSON), &receiptMp)

	receipt := getReceiptByTxHash(receiptMp)
	if receipt.ContractAddress != "0x" ||
		receipt.UsedGas.Cmp(big.NewInt(21000)) != 0 ||
		receipt.Failed != false ||
		receipt.PostState != "0xdd0b0fc6605bbb2e76b8c22ccd466ea5eaa1a80e4860fbdf971be58ded3d782b" ||
		receipt.Result != "0x" ||
		receipt.TotalFee.Cmp(big.NewInt(21000)) != 0 ||
		receipt.TxHash != "0xbd2ca4f9869c714e589ad6a3b16731c8cb066de40d0e27e220cc1e014577baff" {
		t.Fatal("parseToReceipt failed")
	}
}

func TestGetPendingTransactions(t *testing.T) {
	pendingTxsJSON := `[
						{
							"accountNonce":3,
							"amount":10000,
							"from":"0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21",
							"gasLimit":21000,
							"gasPrice":1,
							"hash":"0xc3a8be67dbfe3fc8f9478d91fb49610368515460ad25cba3f566bdf329cdfec6",
							"payload":"",
							"timestamp":1,
							"to":"0xddada93f414f5063cbd4cb642705b1c848fa3c01"
						},
						{
							"accountNonce":4,
							"amount":20000,
							"from":"0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21",
							"gasLimit":31000,
							"gasPrice":1,
							"hash":"0xd3a8be67dbfe3fc8f9478d91fb49610368515460ad25cba3f566bdf329cdfec6",
							"payload":"382342",
							"timestamp":2,
							"to":"0xddada93f414f5063cbd4cb642705b1c848fa3c01"
						}
						]`

	pendingTxsMp := make([]map[string]interface{}, 0)
	json.Unmarshal([]byte(pendingTxsJSON), &pendingTxsMp)
	txs := getPendingTransactions(pendingTxsMp)
	for _, tx := range txs {
		if tx.Hash == "0xc3a8be67dbfe3fc8f9478d91fb49610368515460ad25cba3f566bdf329cdfec6" {
			if tx.AccountNonce != 3 ||
				tx.Amount.Cmp(big.NewInt(10000)) != 0 ||
				tx.From != "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21" ||
				tx.To != "0xddada93f414f5063cbd4cb642705b1c848fa3c01" ||
				tx.GasLimit != 21000 ||
				tx.GasPrice != 1 ||
				tx.Payload != "" ||
				tx.Timestamp != 1 {
				t.Fatal("getPendingTransactions get wrong data")
			} else if tx.Hash == "0xd3a8be67dbfe3fc8f9478d91fb49610368515460ad25cba3f566bdf329cdfec6" {
				if tx.AccountNonce != 4 ||
					tx.Amount.Cmp(big.NewInt(20000)) != 0 ||
					tx.From != "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21" ||
					tx.To != "0xddada93f414f5063cbd4cb642705b1c848fa3c01" ||
					tx.GasLimit != 31000 ||
					tx.GasPrice != 1 ||
					tx.Payload != "382342" ||
					tx.Timestamp != 2 {
					t.Fatal("getPendingTransactions get wrong data")
				}
			}
		}
	}
}
