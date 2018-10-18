/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package rpc

import (
	"fmt"
	"math/big"
)

// CurrentBlockHeight gets the current blockchain height
func (rpc *SeeleRPC) CurrentBlockHeight() (uint64, error) {
	var height uint64
	if err := rpc.call("seele_getBlockHeight", nil, &height); err != nil {
		return 0, err
	}

	return height, nil
}

//GetBlockByHeight get block and transaction data from seele node
func (rpc *SeeleRPC) GetBlockByHeight(h uint64, fullTx bool) (block *BlockInfo, err error) {
	request := GetBlockByHeightRequest{
		Height: int64(h),
		FullTx: fullTx,
	}
	var req []interface{}
	req = append(req, request.Height)
	req = append(req, request.FullTx)
	rpcOutputBlock := make(map[string]interface{})
	if err := rpc.call("seele_getBlockByHeight", req, &rpcOutputBlock); err != nil {
		return nil, err
	}

	headerMp := rpcOutputBlock["header"].(map[string]interface{})
	height := uint64(headerMp["Height"].(float64))
	hash := rpcOutputBlock["hash"].(string)
	parentHash := headerMp["PreviousBlockHash"].(string)
	stateHash := headerMp["StateHash"].(string)
	txHash := headerMp["TxHash"].(string)
	creator := headerMp["Creator"].(string)
	timestamp := int64(headerMp["CreateTimestamp"].(float64))
	difficulty := int64(headerMp["Difficulty"].(float64))
	totalDifficulty := int64(rpcOutputBlock["totalDifficulty"].(float64))

	var Txs []Transaction
	if fullTx {
		var rpcTxs []interface{}
		rpcTxs = rpcOutputBlock["transactions"].([]interface{})
		for i := 0; i < len(rpcTxs); i++ {
			var tx Transaction
			rpcTx := rpcTxs[i].(map[string]interface{})
			tx.Hash = rpcTx["hash"].(string)
			tx.From = rpcTx["from"].(string)
			tx.To = rpcTx["to"].(string)
			amount := int64(rpcTx["amount"].(float64))
			tx.Amount = big.NewInt(amount)
			tx.AccountNonce = uint64(rpcTx["accountNonce"].(float64))
			tx.Payload = rpcTx["payload"].(string)
			tx.Timestamp = uint64(rpcTx["timestamp"].(float64))
			tx.GasLimit = int64(rpcTx["gasLimit"].(float64))
			tx.GasPrice = int64(rpcTx["gasPrice"].(float64))
			Txs = append(Txs, tx)
		}
	}

	var Debts []Debt
	if fullTx {
		var rpcDebt []interface{}
		rpcDebt = rpcOutputBlock["debts"].([]interface{})
		for i := 0; i < len(rpcDebt); i++ {
			var de Debt
			rpcDebtinfo := rpcDebt[i].(map[string]interface{})
			dbets := rpcDebtinfo["Data"].(map[string]interface{})
			de.TxHash = dbets["TxHash"].(string)
			de.Hash = rpcDebtinfo["Hash"].(string)
			de.Block = height
			de.To = dbets["Account"].(string)
			de.ShardNumber = int(dbets["Shard"].(float64))
			amount := int64(dbets["Amount"].(float64))
			de.Amount = big.NewInt(amount)
			de.Payload = dbets["Code"].(string)
			de.Fee = int64(dbets["Fee"].(float64))
			Debts = append(Debts, de)
		}
	}

	var TxDebts []TxDebt
	if fullTx {
		var rpcTxDebt []interface{}
		rpcTxDebt = rpcOutputBlock["txDebts"].([]interface{})
		for i := 0; i < len(rpcTxDebt); i++ {
			var txDebt TxDebt
			rpcDebtinfo := rpcTxDebt[i].(map[string]interface{})
			txdbets := rpcDebtinfo["Data"].(map[string]interface{})
			txDebt.Hash = rpcDebtinfo["Hash"].(string)
			txDebt.TxHash = txdbets["TxHash"].(string)
			txDebt.To = txdbets["Account"].(string)
			txDebt.ShardNumber = int(txdbets["Shard"].(float64))
			amount := int64(txdbets["Amount"].(float64))
			txDebt.Amount = big.NewInt(amount)
			txDebt.Payload = txdbets["Code"].(string)
			txDebt.Fee = int64(txdbets["Fee"].(float64))
			TxDebts = append(TxDebts, txDebt)
		}
	}

	block = &BlockInfo{
		Height:          height,
		Hash:            hash,
		ParentHash:      parentHash,
		StateHash:       stateHash,
		TxHash:          txHash,
		Creator:         creator,
		Timestamp:       big.NewInt(timestamp),
		Difficulty:      big.NewInt(difficulty),
		TotalDifficulty: big.NewInt(totalDifficulty),
		Txs:             Txs,
		Debts:           Debts,
		TxDebts:         TxDebts,
	}
	return block, err
}

// GetPeersInfo get peers info from connected seele node
func (rpc *SeeleRPC) GetPeersInfo() (result []PeerInfo, err error) {
	rpcPeerInfos := make([]map[string]interface{}, 0)
	if err := rpc.call("network_getPeersInfo", nil, &rpcPeerInfos); err != nil {
		return nil, err
	}

	// result data struct:
	// []map[
	//   id:0x0ea2a45ab5a909c309439b0e004c61b7b2a3e831
	//   caps:[lightSeele_1/1 lightSeele_2/1 seele/1]
	//   network:
	//     map[
	//       localAddress:127.0.0.1:8057
	//       remoteAddress:127.0.0.1:54337
	//     ]
	//   protocols:
	//     map[
	//       lightSeele_2:handshake
	//       seele:
	//         map[
	//           version:1
	//           difficulty:7.926036971e+09
	//           head:0000017b5835582b259848c6b0e21d35d90408205c1a41e0aeebe6a67797b8a8
	//         ]
	//       lightSeele_1:handshake
	//     ]
	//   shard:2
	// ]
	return getPeerInfos(rpcPeerInfos), nil
}

// getPeerInfos parse peer informations map to PeerInfo
func getPeerInfos(infos []map[string]interface{}) []PeerInfo {
	var peerInfos []PeerInfo
	for _, rpcPeerInfo := range infos {
		id := rpcPeerInfo["id"].(string)
		rpcCaps := rpcPeerInfo["caps"].([]interface{})
		var caps []string
		for j := 0; j < len(rpcCaps); j++ {
			capString := rpcCaps[j].(string)
			caps = append(caps, capString)
		}
		rpcPeerNetWork := rpcPeerInfo["network"].(map[string]interface{})
		localAddress := rpcPeerNetWork["localAddress"].(string)
		remoteAddress := rpcPeerNetWork["remoteAddress"].(string)
		shardNumber := int(rpcPeerInfo["shard"].(float64))

		peerInfo := PeerInfo{
			ID:            id,
			Caps:          caps,
			LocalAddress:  localAddress,
			RemoteAddress: remoteAddress,
			ShardNumber:   shardNumber,
		}

		peerInfos = append(peerInfos, peerInfo)
	}

	return peerInfos
}

// GetBalance get the balance of the account
func (rpc *SeeleRPC) GetBalance(account string) (int64, error) {
	balanceMp := make(map[string]interface{})
	var request []interface{}
	request = append(request, account)
	if err := rpc.call("seele_getBalance", request, &balanceMp); err != nil {
		return 0, err
	}

	// result data struct:
	// map[
	//   Balance:1.9975499e+12
	//   Account:0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21
	// ]

	return getBalance(balanceMp, account)
}

// getBalance parse balance informations map to int64
func getBalance(balanceMp map[string]interface{}, account string) (int64, error) {
	retAccount := balanceMp["Account"].(string)
	if retAccount != account {
		return 0, fmt.Errorf("expected balance '%s', actually '%s'", account, retAccount)
	}
	return int64(balanceMp["Balance"].(float64)), nil
}

// GetReceiptByTxHash get the receipt by tx hash
func (rpc *SeeleRPC) GetReceiptByTxHash(txhash string) (*Receipt, error) {
	receiptMp := make(map[string]interface{})
	var request []interface{}
	abiJSON := ""
	request = append(request, txhash, abiJSON)
	if err := rpc.call("txpool_getReceiptByTxHash", request, &receiptMp); err != nil {
		return nil, err
	}

	// result data struct:
	// map[
	//   poststate:0x95645120bcdc5f07dc3b8f30f0f3d4069d3374cf0167575f8be474d6c3ad7038
	//   result:0x
	//   totalFee:1
	//   txhash:0x02c240f019adc8b267b82026aef6b677c67867624e2acc1418149e7f8083ba0e
	//   usedGas:0
	//   contract:0x
	//   failed:false
	// ]

	return getReceiptByTxHash(receiptMp), nil
}

// getReceiptByTxHash parse receipt map to Receipt
func getReceiptByTxHash(receiptMp map[string]interface{}) *Receipt {
	result := receiptMp["result"].(string)
	postState := receiptMp["poststate"].(string)
	txHash := receiptMp["txhash"].(string)
	contractAddress := receiptMp["contract"].(string)
	failed := receiptMp["failed"].(bool)
	totalFee := int64(receiptMp["totalFee"].(float64))
	usedGas := int64(receiptMp["usedGas"].(float64))

	return &Receipt{
		Result:          result,
		PostState:       postState,
		TxHash:          txHash,
		ContractAddress: contractAddress,
		Failed:          failed,
		TotalFee:        big.NewInt(totalFee),
		UsedGas:         big.NewInt(usedGas),
	}
}

// GetPendingTransactions get pending transactions on seele node
func (rpc *SeeleRPC) GetPendingTransactions() ([]Transaction, error) {
	txsMp := make([]map[string]interface{}, 0)
	if err := rpc.call("debug_getPendingTransactions", nil, &txsMp); err != nil {
		return nil, err
	}

	// result data struct:
	// []map[
	//   from:0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21
	//   to:0xddada93f414f5063cbd4cb642705b1c848fa3c01
	//   hash:0x1c54966c31a174a404f0feb2ade442e87c09a2802bef41b7efd6ac088d2a0edb
	//   payload:
	//   timestamp:0
	//	 accountNonce:1
	//   amount:10000
	//   gasLimit:21000
	//   gasPrice:1
	// ]

	return getPendingTransactions(txsMp), nil
}

// getPendingTransactions parse txs map to Transaction
func getPendingTransactions(txsMp []map[string]interface{}) []Transaction {
	var Txs []Transaction
	for _, rpcTx := range txsMp {
		var tx Transaction
		tx.Hash = rpcTx["hash"].(string)
		tx.From = rpcTx["from"].(string)
		tx.To = rpcTx["to"].(string)
		amount := int64(rpcTx["amount"].(float64))
		tx.Amount = big.NewInt(amount)
		tx.AccountNonce = uint64(rpcTx["accountNonce"].(float64))
		tx.Payload = rpcTx["payload"].(string)
		tx.Timestamp = uint64(rpcTx["timestamp"].(float64))
		tx.GasLimit = int64(rpcTx["gasLimit"].(float64))
		tx.GasPrice = int64(rpcTx["gasPrice"].(float64))
		Txs = append(Txs, tx)
	}

	return Txs
}
