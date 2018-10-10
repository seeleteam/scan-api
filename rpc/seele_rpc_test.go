package rpc

import (
	"errors"
	"fmt"
	"testing"
)

const (
	SEELEADDRESS      = "172.16.0.197:8027"
	SEELEACCOUNT      = "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21"
	WRONGSEELEACCOUNT = "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb211"
	Height            = 10386
	fullTx            = true
	TXHASH            = "0x02c240f019adc8b267b82026aef6b677c67867624e2acc1418149e7f8083ba0e"
	WRONGTXHASH       = "0x02c240f019adc8b267b82026aef6b677c67867624e2acc1418149e7f8083ba0e1"
)

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
func TestCurrentBlock(t *testing.T) {
	rpc, err := newRPC(SEELEADDRESS)
	if err != nil {
		t.Fatal("newRPC error:", err)
	}
	currentBlock, err := rpc.CurrentBlock()
	if err != nil {
		t.Fatal("GetPeersInfo failed:", err)
	}
	t.Log(currentBlock)
}
func TestGetBlockByHeight(t *testing.T) {
	rpc, err := newRPC(SEELEADDRESS)
	if err != nil {
		t.Fatal("newRPC error:", err)
	}
	block, err := rpc.GetBlockByHeight(Height, fullTx)
	if err != nil {
		t.Fatal("GetBalance failed", err)
	}
	t.Log(block)
}
