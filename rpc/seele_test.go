/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package rpc

import (
	"fmt"
	"testing"
)

func TestCurrentBlock(t *testing.T) {
	defer ReleaseSeeleRPC()
	rpcSeeleRPC, err := GetSeeleRPC()
	if err != nil {
		t.Fatalf("rpc error, %v", err)
	}

	currentBlock, err := rpcSeeleRPC.CurrentBlock()
	if err != nil {
		t.Fatalf("rpc error, %v", err)
	}

	fmt.Println(currentBlock)
}

func TestGetBlockByHeight(t *testing.T) {
	defer ReleaseSeeleRPC()
	rpcSeeleRPC, err := GetSeeleRPC()
	if err != nil {
		t.Fatalf("rpc error, %v", err)
	}

	currentBlock, err := rpcSeeleRPC.CurrentBlock()
	if err != nil {
		t.Fatalf("rpc error, %v", err)
	}

	rpcBlock, err := rpcSeeleRPC.GetBlockByHeight(currentBlock.Height-1, true)
	if err != nil {
		t.Fatalf("rpc error, %v", err)
	}

	fmt.Println(rpcBlock)
}

func TestGetPeersInfo(t *testing.T) {
	defer ReleaseSeeleRPC()
	rpcSeeleRPC, err := GetSeeleRPC()
	if err != nil {
		t.Fatalf("rpc error, %v", err)
	}

	rpcSeeleRPC.GetPeersInfo()
}
