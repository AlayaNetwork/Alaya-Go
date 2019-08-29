package cbft

import (
	"fmt"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/stretchr/testify/assert"
)

func TestViewChange(t *testing.T) {
	pk, sk, cbftnodes := GenerateCbftNode(4)
	nodes := make([]*TestCBFT, 0)
	for i := 0; i < 4; i++ {
		node := MockNode(pk[i], sk[i], cbftnodes, 10000, 10)
		assert.Nil(t, node.Start())

		nodes = append(nodes, node)
	}

	// TestTryViewChange
	testTryViewChange(t, nodes)

	// TestTryChangeViewByViewChange
	testTryChangeViewByViewChange(t, nodes)
}

func testTryViewChange(t *testing.T, nodes []*TestCBFT) {

	result := make(chan *types.Block, 1)

	parent := nodes[0].chain.Genesis()
	for i := 0; i < 4; i++ {
		block := NewBlock(parent.Hash(), parent.NumberU64()+1)
		assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		nodes[0].engine.OnSeal(block, result, nil)

		_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
		select {
		case b := <-result:
			assert.NotNil(t, b)
			assert.Equal(t, uint32(i-1), nodes[0].engine.state.MaxQCIndex())
			for j := 1; j < 3; j++ {
				msg := &protocols.PrepareVote{
					Epoch:          nodes[0].engine.state.Epoch(),
					ViewNumber:     nodes[0].engine.state.ViewNumber(),
					BlockIndex:     uint32(i),
					BlockHash:      b.Hash(),
					BlockNumber:    b.NumberU64(),
					ValidatorIndex: uint32(j),
					ParentQC:       qc,
				}
				assert.Nil(t, nodes[j].engine.signMsgByBls(msg))
				assert.Nil(t, nodes[0].engine.OnPrepareVote("id", msg), fmt.Sprintf("number:%d", b.NumberU64()))
			}
			parent = b
		}
	}
	time.Sleep(10 * time.Second)

	block := nodes[0].engine.state.HighestQCBlock()
	block, qc := nodes[0].engine.blockTree.FindBlockAndQC(block.Hash(), block.NumberU64())

	for i := 0; i < 4; i++ {
		epoch, view := nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()
		viewchange := &protocols.ViewChange{
			Epoch:          epoch,
			ViewNumber:     view,
			BlockHash:      block.Hash(),
			BlockNumber:    block.NumberU64(),
			ValidatorIndex: uint32(i),
			PrepareQC:      qc,
		}
		assert.Nil(t, nodes[i].engine.signMsgByBls(viewchange))
		assert.Nil(t, nodes[0].engine.OnViewChanges("id", &protocols.ViewChanges{
			VCs: []*protocols.ViewChange{
				viewchange,
			},
		}))
	}
	assert.NotNil(t, nodes[0].engine.state.LastViewChangeQC())

	assert.Equal(t, uint64(1), nodes[0].engine.state.ViewNumber())

}

func testTryChangeViewByViewChange(t *testing.T, nodes []*TestCBFT) {
	// note:
	// node-0 has been successfully switched to view-1, HighestQC blockNumber = 4

	// build a duplicate block-4
	number, hash := nodes[0].engine.HighestQCBlockBn()
	block, _ := nodes[0].engine.blockTree.FindBlockAndQC(hash, number)
	dulBlock := NewBlock(block.ParentHash(), block.NumberU64())
	_, preQC := nodes[0].engine.blockTree.FindBlockAndQC(block.ParentHash(), block.NumberU64()-1)
	// Vote and generate prepareQC for dulBlock
	votes := make(map[uint32]*protocols.PrepareVote)
	for j := 1; j < 3; j++ {
		vote := &protocols.PrepareVote{
			Epoch:          nodes[0].engine.state.Epoch(),
			ViewNumber:     nodes[0].engine.state.ViewNumber(),
			BlockIndex:     uint32(0),
			BlockHash:      dulBlock.Hash(),
			BlockNumber:    dulBlock.NumberU64(),
			ValidatorIndex: uint32(j),
			ParentQC:       preQC,
		}
		assert.Nil(t, nodes[j].engine.signMsgByBls(vote))
		votes[uint32(j)] = vote
	}
	dulQC := nodes[0].engine.generatePrepareQC(votes)
	// build new viewChange
	viewChanges := make(map[uint32]*protocols.ViewChange)
	for j := 1; j < 3; j++ {
		viewchange := &protocols.ViewChange{
			Epoch:          nodes[0].engine.state.Epoch(),
			ViewNumber:     nodes[0].engine.state.ViewNumber(),
			BlockHash:      dulBlock.Hash(),
			BlockNumber:    dulBlock.NumberU64(),
			ValidatorIndex: uint32(j),
			PrepareQC:      dulQC,
		}
		assert.Nil(t, nodes[j].engine.signMsgByBls(viewchange))
		viewChanges[uint32(j)] = viewchange
	}
	viewChangeQC := nodes[0].engine.generateViewChangeQC(viewChanges)

	// local highestqc is behind other validators and not exist viewChangeQC.maxBlock, sync qc block
	nodes[0].engine.tryChangeViewByViewChange(viewChangeQC)
	assert.Equal(t, uint64(1), nodes[0].engine.state.ViewNumber())
	assert.Equal(t, hash, nodes[0].engine.state.HighestQCBlock().Hash())

	// local highestqc is behind other validators but exist viewChangeQC.maxBlock, change the view
	nodes[0].engine.blockTree.InsertQCBlock(dulBlock, dulQC)
	nodes[0].engine.tryChangeViewByViewChange(viewChangeQC)
	assert.Equal(t, uint64(2), nodes[0].engine.state.ViewNumber())
	assert.Equal(t, dulQC.BlockHash, nodes[0].engine.state.HighestQCBlock().Hash())

	// based on the view-2 build a duplicate block-4
	dulBlock = NewBlock(block.ParentHash(), block.NumberU64())
	// Vote and generate prepareQC for dulBlock
	votes = make(map[uint32]*protocols.PrepareVote)
	for j := 1; j < 3; j++ {
		vote := &protocols.PrepareVote{
			Epoch:          nodes[0].engine.state.Epoch(),
			ViewNumber:     nodes[0].engine.state.ViewNumber(),
			BlockIndex:     uint32(0),
			BlockHash:      dulBlock.Hash(),
			BlockNumber:    dulBlock.NumberU64(),
			ValidatorIndex: uint32(j),
			ParentQC:       preQC,
		}
		assert.Nil(t, nodes[j].engine.signMsgByBls(vote))
		votes[uint32(j)] = vote
	}
	dulQC = nodes[0].engine.generatePrepareQC(votes)
	nodes[0].engine.blockTree.InsertQCBlock(dulBlock, dulQC)
	nodes[0].engine.state.SetHighestQCBlock(dulBlock)
	// local highestqc is ahead other validators, generate new viewChange quorumCert and change the view
	nodes[0].engine.tryChangeViewByViewChange(viewChangeQC)
	assert.Equal(t, uint64(3), nodes[0].engine.state.ViewNumber())
	assert.Equal(t, dulQC.BlockHash, nodes[0].engine.state.HighestQCBlock().Hash())
	_, _, _, blockView, _, _ := viewChangeQC.MaxBlock()
	assert.Equal(t, uint64(2), blockView)
}
