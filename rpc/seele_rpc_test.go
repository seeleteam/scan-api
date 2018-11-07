/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package rpc

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, len(peers), 0)

	json.Unmarshal([]byte(peersJSON), &peersMp)

	peers = getPeerInfos(peersMp)
	assert.Equal(t, len(peers), len(peersMp))
	for _, info := range peers {
		if info.ID == "0x0ea2a45ab5a909c309439b0e004c61b7b2a3e831" {
			assert.Contains(t, info.Caps, "lightSeele/1")
			assert.Contains(t, info.Caps, "seele/1")
			assert.Equal(t, info.ShardNumber, 2)
			assert.Equal(t, info.LocalAddress, "127.0.0.1:55239")
			assert.Equal(t, info.RemoteAddress, "127.0.0.1:8058")
		} else if info.ID == "0x0da2a45ab5a909c309439b0e004c61b7b2a3e831" {
			assert.Contains(t, info.Caps, "lightSeele/1")
			assert.Contains(t, info.Caps, "seele/1")
			assert.Equal(t, info.ShardNumber, 1)
			assert.Equal(t, info.LocalAddress, "127.0.0.1:55239")
			assert.Equal(t, info.RemoteAddress, "127.0.0.1:8058")
		} else {
			assert.Equal(t, info.ID, "")
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
	assert.Equal(t, err, nil)
	assert.Equal(t, balance, int64(261899990000))
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
	assert.Equal(t, receipt.ContractAddress, "0x")
	assert.Equal(t, receipt.UsedGas, int64(21000))
	assert.Equal(t, receipt.Failed, false)
	assert.Equal(t, receipt.PostState, "0xdd0b0fc6605bbb2e76b8c22ccd466ea5eaa1a80e4860fbdf971be58ded3d782b")
	assert.Equal(t, receipt.Result, "0x")
	assert.Equal(t, receipt.TotalFee, int64(21000))
	assert.Equal(t, receipt.TxHash, "0xbd2ca4f9869c714e589ad6a3b16731c8cb066de40d0e27e220cc1e014577baff")
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
	assert.Equal(t, len(txs), 2)
	for _, tx := range txs {
		if tx.Hash == "0xc3a8be67dbfe3fc8f9478d91fb49610368515460ad25cba3f566bdf329cdfec6" {
			assert.Equal(t, tx.AccountNonce, uint64(3))
			assert.Equal(t, tx.Amount.Cmp(big.NewInt(10000)), 0)
			assert.Equal(t, tx.From, "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21")
			assert.Equal(t, tx.To, "0xddada93f414f5063cbd4cb642705b1c848fa3c01")
			assert.Equal(t, tx.GasLimit, int64(21000))
			assert.Equal(t, tx.GasPrice, int64(1))
			assert.Equal(t, tx.Payload, "")
			assert.Equal(t, tx.Timestamp, uint64(1))
		} else if tx.Hash == "0xd3a8be67dbfe3fc8f9478d91fb49610368515460ad25cba3f566bdf329cdfec6" {
			assert.Equal(t, tx.AccountNonce, uint64(4))
			assert.Equal(t, tx.Amount.Cmp(big.NewInt(20000)), 0)
			assert.Equal(t, tx.From, "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21")
			assert.Equal(t, tx.To, "0xddada93f414f5063cbd4cb642705b1c848fa3c01")
			assert.Equal(t, tx.GasLimit, int64(31000))
			assert.Equal(t, tx.GasPrice, int64(1))
			assert.Equal(t, tx.Payload, "382342")
			assert.Equal(t, tx.Timestamp, uint64(2))
		} else {
			assert.Equal(t, tx.Hash, "")
		}
	}
}

func TestGetBlockByHeight(t *testing.T) {
	blockJSON := `{
					"debts": [
						{
							"Hash": "0x0da1ed893e7f0ca2558c193b3b82ed20575a6978bea5b14f282309c69fee368e",
							"Data": {
								"TxHash": "0x58752f8aeb2c69dd2c32059d3ad8b2d3d860c6d92aa2b3b30ff985e564f60fae",
								"Shard": 2,
								"Account": "0x0ea2a45ab5a909c309439b0e004c61b7b2a3e831",
								"Amount": 10000,
								"Fee": 0,
								"Code": ""
							}
						}
					],
        			"hash": "0x000002069d9de64bad509239e2a121afbf7de183576457a1d1fb077d19fa3e8c",
        			"header": {
						"PreviousBlockHash": "0x000001cba2c0b82402b3d2d2ad49f50ca0b21aee18c8123486377b2ec93aa0e0",
						"Creator": "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21",
						"StateHash": "0x8af14975f636ace27571cfcdcd9a1a1b4a5b15228977cf6207e82f63abf96ffd",
						"TxHash": "0xdb00575ff0cc0de89bd6c1799d37e5f600687963785176ca76e81bebfde6a03f",
						"ReceiptHash": "0x02fa1d68e7bbf0b833f6e8719efb11b32c7f760e4ae050a4f9b58b8dd8ad1620",
						"TxDebtHash": "0x58d7c36b25a715f5076ccb878940920f6bb333ab142287452509f881103960d2",
						"DebtHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
						"Difficulty": 6563003,
						"Height": 10368,
						"CreateTimestamp": 1539050098,
						"Nonce": 17825487295277268182,
						"ExtraData": ""
        			},
        			"totalDifficulty": 68985339754,
        			"transactions": [
						{
							"accountNonce": 0,
							"amount": 150000000,
							"from": "0x0000000000000000000000000000000000000000",
							"gasLimit": 0,
							"gasPrice": 0,
							"hash": "0x6fb17b265260caed33b4e8f58ad84b508dd8950b9bc93dae8518fc96912f76bb",
							"payload": "",
							"timestamp": 1539931510,
							"to": "0xd5a145191b7ca9cb4f3dc850e426c1e853d2a9f1"
						},
						{
							"accountNonce": 280,
							"amount": 10000,
							"from": "0xec759db47a65f6537d630517f6cd3ca39c6f93d1",
							"gasLimit": 21000,
							"gasPrice": 1,
							"hash": "0xf526dc404145cd409601e951fec4f2222f3abf578381cdaaea9db3a791a79cbd",
							"payload": "",
							"timestamp": 0,
							"to": "0xa00d22dc3624d4696eff8d1641b442f79c3379b1"
						}
					],
        			"txDebts": [
						{
							"Hash": "0xe1c24a636a7c27aea7c384f6eb61eb49168129105f4c081ffa8ca7e77198b3f6",
							"Data": {
								"TxHash": "0x0b30a6edf95a16933a0a77ffd3eb15680d4e3cb79466f21c1181c013a68eae62",
								"Shard": 2,
								"Account": "0x0ea2a45ab5a909c309439b0e004c61b7b2a3e831",
								"Amount": 10000,
								"Fee": 1,
								"Code": ""
							}
						}
					]
				}
`
	blockMp := make(map[string]interface{})
	json.Unmarshal([]byte(blockJSON), &blockMp)
	block := getBlockByHeight(blockMp, true)
	assert.Equal(t, block.Height, uint64(10368))
	assert.Equal(t, block.Hash, "0x000002069d9de64bad509239e2a121afbf7de183576457a1d1fb077d19fa3e8c")
	assert.Equal(t, block.ParentHash, "0x000001cba2c0b82402b3d2d2ad49f50ca0b21aee18c8123486377b2ec93aa0e0")
	assert.Equal(t, block.StateHash, "0x8af14975f636ace27571cfcdcd9a1a1b4a5b15228977cf6207e82f63abf96ffd")
	assert.Equal(t, block.TxHash, "0xdb00575ff0cc0de89bd6c1799d37e5f600687963785176ca76e81bebfde6a03f")
	assert.Equal(t, block.Creator, "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21")
	assert.Equal(t, block.Timestamp, big.NewInt(1539050098))
	assert.Equal(t, block.Difficulty, big.NewInt(6563003))
	assert.Equal(t, block.TotalDifficulty, big.NewInt(68985339754))

	assert.Equal(t, len(block.Txs), 2)
	for _, tx := range block.Txs {
		if tx.Hash == "0x6fb17b265260caed33b4e8f58ad84b508dd8950b9bc93dae8518fc96912f76bb" {
			assert.Equal(t, tx.From, "0x0000000000000000000000000000000000000000")
			assert.Equal(t, tx.To, "0xd5a145191b7ca9cb4f3dc850e426c1e853d2a9f1")
			assert.Equal(t, tx.Amount, big.NewInt(150000000))
			assert.Equal(t, tx.AccountNonce, uint64(0))
			assert.Equal(t, tx.GasLimit, int64(0))
			assert.Equal(t, tx.GasPrice, int64(0))
			assert.Equal(t, tx.Payload, "")
			assert.Equal(t, tx.Timestamp, uint64(1539931510))
		} else if tx.Hash == "0xf526dc404145cd409601e951fec4f2222f3abf578381cdaaea9db3a791a79cbd" {
			assert.Equal(t, tx.From, "0xec759db47a65f6537d630517f6cd3ca39c6f93d1")
			assert.Equal(t, tx.To, "0xa00d22dc3624d4696eff8d1641b442f79c3379b1")
			assert.Equal(t, tx.Amount, big.NewInt(10000))
			assert.Equal(t, tx.AccountNonce, uint64(280))
			assert.Equal(t, tx.GasLimit, int64(21000))
			assert.Equal(t, tx.GasPrice, int64(1))
			assert.Equal(t, tx.Payload, "")
			assert.Equal(t, tx.Timestamp, uint64(0))
		} else {
			assert.Equal(t, tx.Hash, "")
		}
	}
	assert.Equal(t, len(block.Debts), 1)
	for _, debt := range block.Debts {
		if debt.Hash == "0x0da1ed893e7f0ca2558c193b3b82ed20575a6978bea5b14f282309c69fee368e" {
			assert.Equal(t, debt.TxHash, "0x58752f8aeb2c69dd2c32059d3ad8b2d3d860c6d92aa2b3b30ff985e564f60fae")
			assert.Equal(t, debt.ShardNumber, 2)
			assert.Equal(t, debt.Amount, big.NewInt(10000))
			assert.Equal(t, debt.To, "0x0ea2a45ab5a909c309439b0e004c61b7b2a3e831")
			assert.Equal(t, debt.Fee, int64(0))
			assert.Equal(t, debt.Payload, "")
		} else {
			assert.Equal(t, debt.Hash, "")
		}
	}
	assert.Equal(t, len(block.TxDebts), 1)
	for _, txdebt := range block.TxDebts {
		if txdebt.Hash == "0xe1c24a636a7c27aea7c384f6eb61eb49168129105f4c081ffa8ca7e77198b3f6" {
			assert.Equal(t, txdebt.TxHash, "0x0b30a6edf95a16933a0a77ffd3eb15680d4e3cb79466f21c1181c013a68eae62")
			assert.Equal(t, txdebt.ShardNumber, 2)
			assert.Equal(t, txdebt.Amount, big.NewInt(10000))
			assert.Equal(t, txdebt.To, "0x0ea2a45ab5a909c309439b0e004c61b7b2a3e831")
			assert.Equal(t, txdebt.Fee, int64(1))
			assert.Equal(t, txdebt.Payload, "")
		} else {
			assert.Equal(t, txdebt.Hash, "")
		}
	}
}
