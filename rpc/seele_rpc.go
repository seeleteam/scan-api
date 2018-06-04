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
	request := GetBlockByHeightRequest{
		Height: -1,
		FullTx: true,
	}
	rpcOutputBlock := make(map[string]interface{})
	if err := rpc.call("seele.GetBlockByHeight", request, &rpcOutputBlock); err != nil {
		return nil, err
	}

	timestamp := int64(rpcOutputBlock["timestamp"].(float64))
	difficulty := int64(rpcOutputBlock["difficulty"].(float64))
	height := uint64(rpcOutputBlock["height"].(float64))

	currentBlock = &CurrentBlock{
		HeadHash:  rpcOutputBlock["hash"].(string),
		Height:    height,
		Timestamp: big.NewInt(timestamp),
		Difficult: big.NewInt(difficulty),
		Creator:   rpcOutputBlock["creator"].(string),
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
	rpcOutputBlock := make(map[string]interface{})
	if err := rpc.call("seele.GetBlockByHeight", request, &rpcOutputBlock); err != nil {
		return nil, err
	}

	height := uint64(rpcOutputBlock["height"].(float64))
	hash := rpcOutputBlock["hash"].(string)
	parentHash := rpcOutputBlock["parentHash"].(string)
	nonce := uint64(rpcOutputBlock["nonce"].(float64))
	stateHash := rpcOutputBlock["stateHash"].(string)
	txHash := rpcOutputBlock["txHash"].(string)
	creator := rpcOutputBlock["creator"].(string)
	timestamp := int64(rpcOutputBlock["timestamp"].(float64))
	difficulty := int64(rpcOutputBlock["difficulty"].(float64))
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
			Txs = append(Txs, tx)
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
		prcPeerNetWork := rpcPeerInfo["network"].(map[string]interface{})
		localAddress := prcPeerNetWork["localAddress"].(string)
		remoteAddress := prcPeerNetWork["remoteAddress"].(string)

		peerInfo := PeerInfo{
			ID:            id,
			Caps:          caps,
			LocalAddress:  localAddress,
			RemoteAddress: remoteAddress,
		}

		peerInfos = append(peerInfos, peerInfo)
	}

	return peerInfos, nil
}
