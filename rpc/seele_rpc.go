/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package rpc

import (
	"errors"
	"fmt"
	"github.com/seeleteam/scan-api/log"
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

// GetBlockByHeight get block and transaction data from seele node
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

	// reslut data struct:
	// map[
	//   debts:[
	//     map[
	//       Data:
	//         map[
	//           Shard:2
	//           Account:0x0ea2a45ab5a909c309439b0e004c61b7b2a3e831
	//           Amount:10000
	//           Fee:0
	//           Code:
	//           TxHash:0x58752f8aeb2c69dd2c32059d3ad8b2d3d860c6d92aa2b3b30ff985e564f60fae
	//         ]
	//       Hash:0x0da1ed893e7f0ca2558c193b3b82ed20575a6978bea5b14f282309c69fee368e
	//     ]
	//   ]
	//   hash:0x000002069d9de64bad509239e2a121afbf7de183576457a1d1fb077d19fa3e8c
	//   header:
	//     map[
	//       StateHash:0x8af14975f636ace27571cfcdcd9a1a1b4a5b15228977cf6207e82f63abf96ffd
	//       ReceiptHash:0x02fa1d68e7bbf0b833f6e8719efb11b32c7f760e4ae050a4f9b58b8dd8ad1620
	//       DebtHash:0x0000000000000000000000000000000000000000000000000000000000000000
	//       CreateTimestamp:1.539050098e+09
	//       Nonce:1.782548729527727e+19
	//       ExtraData:
	//       PreviousBlockHash:0x000001cba2c0b82402b3d2d2ad49f50ca0b21aee18c8123486377b2ec93aa0e0
	//       Creator:0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21
	//       TxHash:0xdb00575ff0cc0de89bd6c1799d37e5f600687963785176ca76e81bebfde6a03f
	//       TxDebtHash:0x58d7c36b25a715f5076ccb878940920f6bb333ab142287452509f881103960d2
	//       Difficulty:6.563003e+06
	//       Height:10368
	//     ]
	//   totalDifficulty:6.8985339754e+10
	//   transactions:[
	//     map[
	//       from:0x0000000000000000000000000000000000000000
	//       hash:0x6fb17b265260caed33b4e8f58ad84b508dd8950b9bc93dae8518fc96912f76bb
	//       payload:
	//       timestamp:1.53993151e+09
	//       to:0xd5a145191b7ca9cb4f3dc850e426c1e853d2a9f1
	//       accountNonce:0
	//       amount:1.5e+08
	//       gasLimit:0
	//       gasPrice:0
	//     ]
	//     map[
	//       amount:10000
	//       hash:0xf526dc404145cd409601e951fec4f2222f3abf578381cdaaea9db3a791a79cbd
	//       payload:
	//       timestamp:0
	//       to:0xa00d22dc3624d4696eff8d1641b442f79c3379b1
	//       accountNonce:280
	//       from:0xec759db47a65f6537d630517f6cd3ca39c6f93d1
	//       gasLimit:21000
	//       gasPrice:1
	//     ]
	//   ]
	//   txDebts:[
	//     map[
	//       Hash:0xe1c24a636a7c27aea7c384f6eb61eb49168129105f4c081ffa8ca7e77198b3f6
	//       Data:
	//         map[
	//           Account:0x0ea2a45ab5a909c309439b0e004c61b7b2a3e831
	//           Amount:10000
	//           Fee:1
	//           Code:
	//           TxHash:0x0b30a6edf95a16933a0a77ffd3eb15680d4e3cb79466f21c1181c013a68eae62
	//           Shard:2
	//         ]
	//       ]
	//     ]
	//   ]
	if rpcOutputBlock == nil || rpcOutputBlock["header"] == nil{
		var err =errors.New("seele_rpc rpcOutputBlock is nil")
		log.Error("seele_rpc rpcOutputBlock is nil")
		return nil, err
	}
	return getBlockByHeight(rpcOutputBlock, fullTx), err
}

// getBlockByHeight parse block map to BlockInfo
func getBlockByHeight(rpcOutputBlock map[string]interface{}, fullTx bool) *BlockInfo {
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
			//de.ShardNumber = int(dbets["Shard"].(float64))
			amount := int64(dbets["Amount"].(float64))
			de.Amount = big.NewInt(amount)
			de.Payload = dbets["Code"].(string)
			//de.Fee = int64(dbets["Fee"].(float64))
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
			//txDebt.ShardNumber = int(txdbets["Shard"].(float64))
			amount := int64(txdbets["Amount"].(float64))
			txDebt.Amount = big.NewInt(amount)
			txDebt.Payload = txdbets["Code"].(string)
			//txDebt.Fee = int64(txdbets["Fee"].(float64))
			TxDebts = append(TxDebts, txDebt)
		}
	}

	return &BlockInfo{
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

// getPeerInfos parse peer information map to PeerInfo
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
	request = append(request, account, "", -1)
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
	if err := rpc.call("seele_getReceiptByTxHash", request, &receiptMp); err != nil {
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
		TotalFee:        totalFee,
		UsedGas:         usedGas,
	}
}

// GetPendingTransactions get pending transactions on seele node
func (rpc *SeeleRPC) GetPendingTransactions() ([]Transaction, error) {
	txsMp := make([]map[string]interface{}, 0)
	if err := rpc.call("txpool_getPendingTxs", nil, &txsMp); err != nil {
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
		tx.GasLimit = int64(rpcTx["gasLimit"].(float64))
		tx.GasPrice = int64(rpcTx["gasPrice"].(float64))
		Txs = append(Txs, tx)
	}

	return Txs
}
