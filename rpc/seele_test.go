/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package rpc

import (
	"errors"
	"fmt"
	"testing"
	"math/big"
)

//import (
//	"fmt"
//	"testing"
//)
//
//func TestCurrentBlock(t *testing.T) {
//	defer ReleaseSeeleRPC()
//	rpcSeeleRPC, err := GetSeeleRPC()
//	if err != nil {
//		t.Fatalf("rpc error, %v", err)
//	}
//
//	currentBlock, err := rpcSeeleRPC.CurrentBlock()
//	if err != nil {
//		t.Fatalf("rpc error, %v", err)
//	}
//
//	fmt.Println(currentBlock)
//}
//
//func TestGetBlockByHeight(t *testing.T) {
//	defer ReleaseSeeleRPC()
//	rpcSeeleRPC, err := GetSeeleRPC()
//	if err != nil {
//		t.Fatalf("rpc error, %v", err)
//	}
//
//	currentBlock, err := rpcSeeleRPC.CurrentBlock()
//	if err != nil {
//		t.Fatalf("rpc error, %v", err)
//	}
//
//	rpcBlock, err := rpcSeeleRPC.GetBlockByHeight(currentBlock.Height-1, true)
//	if err != nil {
//		t.Fatalf("rpc error, %v", err)
//	}
//
//	fmt.Println(rpcBlock)
//}
//
//func TestGetPeersInfo(t *testing.T) {
//	defer ReleaseSeeleRPC()
//	rpcSeeleRPC, err := GetSeeleRPC()
//	if err != nil {
//		t.Fatalf("rpc error, %v", err)
//	}
//
//	rpcSeeleRPC.GetPeersInfo()
//}

const (
	SEELEADDRESS      = "172.16.0.197:8027"
	SEELEACCOUNT      = "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21"
	WRONGSEELEACCOUNT = "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb211"
	TXHASH            = "0x02c240f019adc8b267b82026aef6b677c67867624e2acc1418149e7f8083ba0e"
	WRONGTXHASH       = "0x02c240f019adc8b267b82026aef6b677c67867624e2acc1418149e7f8083ba0e1"
)

var zero = big.NewInt(0)

func newRPC(address string) (*SeeleRPC, error) {
	rpc := NewRPC(address)
	if rpc == nil {
		return nil, errors.New("newRPC failed")
	}
	if err := rpc.Connect(); err != nil {
		return nil, fmt.Errorf("rpc connect failed:%s", err.Error())
	}
	return rpc, nil
}

func TestGetPeersInfo(t *testing.T) {
	rpc, err := newRPC(SEELEADDRESS)
	if err != nil {
		t.Fatal("newRPC error:", err)
	}

	peers, err := rpc.GetPeersInfo()
	if err != nil {
		t.Fatal("GetPeersInfo failed:", err)
	}

	if len(peers) < 0 {
		t.Fatal("GetPeersInfo get invalid peers number")
	}
}

func TestGetBalance(t *testing.T) {
	rpc, err := newRPC(SEELEADDRESS)
	if err != nil {
		t.Fatal("newRPC error:", err)
	}

	balance, err := rpc.GetBalance(SEELEACCOUNT)
	if err != nil {
		t.Fatal("GetBalance failed", err)
	}

	if balance < 0 {
		t.Fatal("invalid balance")
	}

	if _, err := rpc.GetBalance(WRONGSEELEACCOUNT); err == nil {
		t.Fatal("GetReceiptByTxHash on wrong account test fail")
	}
}

func TestGetReceiptByTxHash(t *testing.T) {
	rpc, err := newRPC(SEELEADDRESS)
	if err != nil {
		t.Fatal("newRPC error:", err)
	}

	receipt, err := rpc.GetReceiptByTxHash(TXHASH)
	if err != nil {
		t.Fatal("GetReceiptByTxHash failed", err)
	}

	if receipt.Result == "" || receipt.PostState == "" || receipt.TxHash == "" || receipt.ContractAddress == "" ||
		receipt.TotalFee.Cmp(zero) < 0 || receipt.UsedGas.Cmp(zero) < 0 {
			t.Fatal("GetReceiptByTxHash get invalid data")
	}

	if _, err := rpc.GetReceiptByTxHash(WRONGTXHASH); err == nil {
		t.Fatal("GetReceiptByTxHash on wrong txhash test fail")
	}
}

func TestGetPendingTransactions(t *testing.T) {
	rpc, err := newRPC(SEELEADDRESS)
	if err != nil {
		t.Fatal("newRPC error:", err)
	}

	pendingTxs, err := rpc.GetPendingTransactions()
	if err != nil {
		t.Fatal("GetPendingTransactions failed", err)
	}

	if len(pendingTxs) < 0 {
		t.Fatal("GetPendingTransactions get invalid tx number")
	}

	for _, tx := range pendingTxs {
		if tx.Hash == "" || tx.From == "" || tx.To == "" || tx.Amount.Cmp(zero) < 0 ||
			tx.AccountNonce < 0 || tx.Fee < 0 {
			t.Fatal("GetPendingTransactions get invalid data")
		}
	}
}
