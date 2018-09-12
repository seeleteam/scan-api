/**
*  @file
*  @copyright defined in scan-api/LICENSE
 */

package handlers

import (
	"testing"

	"github.com/seeleteam/scan-api/database"
	"github.com/stretchr/testify/assert"
)

func newTestDBBlock(t *testing.T) *database.DBBlock {
	return &database.DBBlock{
		Reward: 99,
		Height: 133853,
	}
}
func Test_CreateRetSimpleBlockInfo(t *testing.T) {
	header := newTestDBBlock(t)
	got := createRetSimpleBlockInfo(header)
	gots := int64(got.Height)
	assert.Equal(t, gots, header.Height)
	assert.Equal(t, got.Reward, header.Reward)

}
