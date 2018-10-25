/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package database

import (
	"encoding/json"
	"math/big"
	"strconv"
	"testing"

	"github.com/seeleteam/scan-api/rpc"
	"github.com/stretchr/testify/assert"
)

func TestCreateDbBlock(t *testing.T) {
	blockInfoJSON := `{"hash":"0x000002069d9de64bad509239e2a121afbf7de183576457a1d1fb077d19fa3e8c","parentHash":"0x000001cba2c0b82402b3d2d2ad49f50ca0b21aee18c8123486377b2ec93aa0e0","height":10368,"stateHash":"0x8af14975f636ace27571cfcdcd9a1a1b4a5b15228977cf6207e82f63abf96ffd","timestamp":1539050098,"difficulty":6563003,"totaldifficulty":68985339754,"creator":"0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21","nonce":0,"txHash":"0xdb00575ff0cc0de89bd6c1799d37e5f600687963785176ca76e81bebfde6a03f","txs":[{"hash":"0x6fb17b265260caed33b4e8f58ad84b508dd8950b9bc93dae8518fc96912f76bb","debtTxHash":"","from":"0x0000000000000000000000000000000000000000","to":"0xd5a145191b7ca9cb4f3dc850e426c1e853d2a9f1","amount":150000000,"accountNonce":0,"payload":"","timestamp":1539931510,"fee":0,"block":0,"idx":0,"txtype":0,"gasLimit":0,"gasPrice":0},{"hash":"0xf526dc404145cd409601e951fec4f2222f3abf578381cdaaea9db3a791a79cbd","debtTxHash":"","from":"0xec759db47a65f6537d630517f6cd3ca39c6f93d1","to":"0xa00d22dc3624d4696eff8d1641b442f79c3379b1","amount":10000,"accountNonce":280,"payload":"","timestamp":0,"fee":21000,"block":0,"idx":0,"txtype":0,"gasLimit":21000,"gasPrice":1}],"debts":[{"hash":"0x0da1ed893e7f0ca2558c193b3b82ed20575a6978bea5b14f282309c69fee368e","txhash":"0x58752f8aeb2c69dd2c32059d3ad8b2d3d860c6d92aa2b3b30ff985e564f60fae","to":"0x0ea2a45ab5a909c309439b0e004c61b7b2a3e831","block":10368,"idx":0,"shardNumber":2,"fee":0,"payload":"","amount":10000}],"txDebts":[{"hash":"0xe1c24a636a7c27aea7c384f6eb61eb49168129105f4c081ffa8ca7e77198b3f6","txhash":"0x0b30a6edf95a16933a0a77ffd3eb15680d4e3cb79466f21c1181c013a68eae62","to":"0x0ea2a45ab5a909c309439b0e004c61b7b2a3e831","shardNumber":2,"fee":1,"payload":"","amount":10000}]}`

	var blockInfo rpc.BlockInfo
	json.Unmarshal([]byte(blockInfoJSON), &blockInfo)
	block := CreateDbBlock(&blockInfo)
	assert.Equal(t, block.Height, int64(10368))
	assert.Equal(t, block.HeadHash, "0x000002069d9de64bad509239e2a121afbf7de183576457a1d1fb077d19fa3e8c")
	assert.Equal(t, block.PreHash, "0x000001cba2c0b82402b3d2d2ad49f50ca0b21aee18c8123486377b2ec93aa0e0")
	assert.Equal(t, block.StateHash, "0x8af14975f636ace27571cfcdcd9a1a1b4a5b15228977cf6207e82f63abf96ffd")
	assert.Equal(t, block.TxHash, "0xdb00575ff0cc0de89bd6c1799d37e5f600687963785176ca76e81bebfde6a03f")
	assert.Equal(t, block.Creator, "0x4c10f2cd2159bb432094e3be7e17904c2b4aeb21")
	assert.Equal(t, block.Timestamp, int64(1539050098))
	assert.Equal(t, block.Difficulty, big.NewInt(6563003).String())
	assert.Equal(t, block.TotalDifficulty, big.NewInt(68985339754).String())

	assert.Equal(t, len(block.Txs), 2)
	for _, tx := range block.Txs {
		if tx.Hash == "0x6fb17b265260caed33b4e8f58ad84b508dd8950b9bc93dae8518fc96912f76bb" {
			assert.Equal(t, tx.From, "0x0000000000000000000000000000000000000000")
			assert.Equal(t, tx.To, "0xd5a145191b7ca9cb4f3dc850e426c1e853d2a9f1")
			assert.Equal(t, tx.Amount, int64(150000000))
			assert.Equal(t, tx.Timestamp, strconv.Itoa(1539931510))
		} else if tx.Hash == "0xf526dc404145cd409601e951fec4f2222f3abf578381cdaaea9db3a791a79cbd" {
			assert.Equal(t, tx.From, "0xec759db47a65f6537d630517f6cd3ca39c6f93d1")
			assert.Equal(t, tx.To, "0xa00d22dc3624d4696eff8d1641b442f79c3379b1")
			assert.Equal(t, tx.Amount, int64(10000))
			assert.Equal(t, tx.Timestamp, strconv.Itoa(0))
		} else {
			assert.Equal(t, tx.Hash, "")
		}
	}
	assert.Equal(t, len(block.Debts), 1)
	for _, debt := range block.Debts {
		if debt.Hash == "0x0da1ed893e7f0ca2558c193b3b82ed20575a6978bea5b14f282309c69fee368e" {
			assert.Equal(t, debt.TxHash, "0x58752f8aeb2c69dd2c32059d3ad8b2d3d860c6d92aa2b3b30ff985e564f60fae")
			assert.Equal(t, debt.ShardNumber, 2)
			assert.Equal(t, debt.Amount, int64(10000))
			assert.Equal(t, debt.Account, "0x0ea2a45ab5a909c309439b0e004c61b7b2a3e831")
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
			assert.Equal(t, txdebt.Amount, int64(10000))
			assert.Equal(t, txdebt.Account, "0x0ea2a45ab5a909c309439b0e004c61b7b2a3e831")
			assert.Equal(t, txdebt.Fee, int64(1))
			assert.Equal(t, txdebt.Payload, "")
		} else {
			assert.Equal(t, txdebt.Hash, "")
		}
	}
}
