/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package rpc

import (
	"fmt"
	"math/big"
)

// CurrentBlock is the informations about the best block
type CurrentBlock struct {
	HeadHash  string   `json:"headHash"`
	Height    uint64   `json:"height"`
	Timestamp *big.Int `json:"timestamp"`
	Difficult *big.Int `json:"difficult"`
	Creator   string   `json:"creator"`
	TxCount   int      `json:"txcount"`
}

//Transaction is the transaction data send from seele node
type Transaction struct {
	Hash         string   `json:"hash"`
	From         string   `json:"from"`
	To           string   `json:"to"`
	Amount       *big.Int `json:"amount"`
	AccountNonce uint64   `json:"accountNonce"`
	Payload      string   `json:"payload"`
	Timestamp    uint64   `json:"timestamp"`
	Fee          int64    `json:"fee"`
	Block        uint64   `json:"block"`
	Idx          uint64   `json:"idx"`
	TxType       int      `json:"txtype"`
}

//BlockInfo is the block data send from seele node
type BlockInfo struct {
	Hash            string        `json:"hash"`
	ParentHash      string        `json:"parentHash"`
	Height          uint64        `json:"height"`
	StateHash       string        `json:"stateHash"`
	Timestamp       *big.Int      `json:"timestamp"`
	Difficulty      *big.Int      `json:"difficulty"`
	TotalDifficulty *big.Int      `json:"totaldifficulty"`
	Creator         string        `json:"creator"`
	Nonce           uint64        `json:"nonce"`
	TxHash          string        `json:"txHash"`
	Txs             []Transaction `json:"txs"`
}

// GetBlockByHeightRequest request param for GetBlockByHeight api
type GetBlockByHeightRequest struct {
	Height int64 `json:"height"`
	FullTx bool  `json:"fullTx"`
}

//PeerInfo is the peer info send from seele node
type PeerInfo struct {
	ID            string   `json:"id"`            // Unique of the node
	Caps          []string `json:"caps"`          // Sum-protocols advertised by this particular peer
	LocalAddress  string   `json:"localAddress"`  // Local endpoint of the TCP data connection
	RemoteAddress string   `json:"remoteAddress"` // Remote endpoint of the TCP data connection
	ShardNumber   int      `json:"shardNumber"`
}

// Receipt is the receipt information of seele
type Receipt struct {
	Result          string   `json:"result"`
	PostState       string   `json:"result"`
	TxHash          string   `json:"txhash"`
	ContractAddress string   `json:"contractaddress"`
	Failed          bool     `json:"failed"`
	TotalFee        *big.Int `json:"totalFee"`
	UsedGas         *big.Int `json:"usedGas"`
}

// String implement fmt.Stringer interface
func (r *Receipt) String() string {
	return fmt.Sprintf("result:%s, poststate:%s, txhash:%s, contractAddress:%s, failed:%v, totalfee:%v, usedgas:%v",
		r.Result, r.PostState, r.TxHash, r.ContractAddress, r.Failed, r.TotalFee.Int64(), r.UsedGas.Int64())
}
