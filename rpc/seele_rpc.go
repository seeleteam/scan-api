/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package rpc

import (
	"math/big"
)

// CurrentBlock returns the current block info.
func (rpc *SeeleRPC) CurrentBlock() (currentBlock *CurrentBlock, err error) {
	var request []interface{}
	request = append(request, -1)
	request = append(request, true)

	rpcOutputBlock := make(map[string]interface{})
	if err := rpc.call("seele_getBlockByHeight", request, &rpcOutputBlock); err != nil {
		return nil, err
	}

	result := rpcOutputBlock["header"].(map[string]interface{})

	timestamp := int64(result["CreateTimestamp"].(float64))
	difficulty := int64(result["Difficulty"].(float64))
	height := uint64(result["Height"].(float64))
	currentBlock = &CurrentBlock{
		HeadHash:  rpcOutputBlock["hash"].(string),
		Height:    height,
		Timestamp: big.NewInt(timestamp),
		Difficult: big.NewInt(difficulty),
		Creator:   result["Creator"].(string),
		TxCount:   len(rpcOutputBlock["transactions"].([]interface{})),
	}

	return currentBlock, err
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
	nonce := uint64(headerMp["Nonce"].(float64))
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
			tx.Fee = int64(rpcTx["fee"].(float64))
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
			txsdbet := rpcDebtinfo["Data"].(map[string]interface{})
			de.TxHash = txsdbet["TxHash"].(string)
			de.To = txsdbet["Account"].(string)
			de.Shard = int64(txsdbet["Shard"].(float64))
			amount := int64(txsdbet["Amount"].(float64))
			de.Amount = big.NewInt(amount)
			de.Code = txsdbet["Code"].(string)
			de.Fee = int64(txsdbet["Fee"].(float64))
			Debts = append(Debts, de)
		}
	}

	block = &BlockInfo{
		Height:          height,
		Hash:            hash,
		ParentHash:      parentHash,
		Nonce:           nonce,
		StateHash:       stateHash,
		TxHash:          txHash,
		Creator:         creator,
		Timestamp:       big.NewInt(timestamp),
		Difficulty:      big.NewInt(difficulty),
		TotalDifficulty: big.NewInt(totalDifficulty),
		Txs:             Txs,
		TxDebt:          Debts,
	}
	return block, err
}

//GetPeersInfo get peers info from connected seele node
func (rpc *SeeleRPC) GetPeersInfo() (result []PeerInfo, err error) {
	var rcpPeerInfos []interface{}
	if err := rpc.call("network.GetPeersInfo", nil, &rcpPeerInfos); err != nil {
		return nil, err
	}

	var peerInfos []PeerInfo
	for i := 0; i < len(rcpPeerInfos); i++ {
		rpcPeerInfo := rcpPeerInfos[i].(map[string]interface{})

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

	return peerInfos, nil
}

//GetBalance get the balance of the account
func (rpc *SeeleRPC) GetBalance(address string) (int64, error) {
	var result interface{}
	if err := rpc.call("seele.GetBalance", &address, &result); err != nil {
		return 0, err
	}

	balance := int64(result.(float64))
	return balance, nil
}

//GetReceiptByTxHash
func (rpc *SeeleRPC) GetReceiptByTxHash(txhash string) (*Receipt, error) {

	rpcOutputReceipt := make(map[string]interface{})
	if err := rpc.call("txpool.GetReceiptByTxHash", &txhash, &rpcOutputReceipt); err != nil {
		return nil, err
	}

	result := rpcOutputReceipt["result"].(string)
	postState := rpcOutputReceipt["poststate"].(string)
	txHash := rpcOutputReceipt["txhash"].(string)
	contractAddress := rpcOutputReceipt["contract"].(string)

	receipt := Receipt{
		Result:          result,
		PostState:       postState,
		TxHash:          txHash,
		ContractAddress: contractAddress,
	}
	return &receipt, nil
}

//GetPendingTransactions
func (rpc *SeeleRPC) GetPendingTransactions() ([]Transaction, error) {
	var rpcOutputTxs []interface{}
	if err := rpc.call("txpool.GetPendingTransactions", nil, &rpcOutputTxs); err != nil {
		return nil, err
	}

	var Txs []Transaction

	for i := 0; i < len(rpcOutputTxs); i++ {
		var tx Transaction
		rpcTx := rpcOutputTxs[i].(map[string]interface{})
		tx.Hash = rpcTx["hash"].(string)
		tx.From = rpcTx["from"].(string)
		tx.To = rpcTx["to"].(string)
		amount := int64(rpcTx["amount"].(float64))
		tx.Amount = big.NewInt(amount)
		tx.AccountNonce = uint64(rpcTx["accountNonce"].(float64))
		tx.Payload = rpcTx["payload"].(string)
		tx.Timestamp = uint64(rpcTx["timestamp"].(float64))
		tx.Fee = int64(rpcTx["fee"].(float64))
		Txs = append(Txs, tx)
	}

	return Txs, nil
}
