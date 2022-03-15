// Copyright 2021 The Alaya Network Authors
// This file is part of the Alaya-Go library.
//
// The Alaya-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Alaya-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Alaya-Go library. If not, see <http://www.gnu.org/licenses/>.

package cbft

import (
	"fmt"
	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/protocols"
	ctypes "github.com/AlayaNetwork/Alaya-Go/consensus/cbft/types"
	"github.com/AlayaNetwork/Alaya-Go/core/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCbft_SyncPrepareVoteV2(t *testing.T) {
	nodesNum := 50
	nodes := MockRGNodes(t, nodesNum)

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)
	parent := nodes[0].chain.Genesis()

	block := NewBlockWithSign(parent.Hash(), parent.NumberU64()+1, nodes[0])
	assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
	nodes[0].engine.OnSeal(block, result, nil, complete)
	<-complete

	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
	epoch, viewNumber, blockIndex := nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), uint32(0)
	missValidatorIndex, syncValidatorIndex := uint32(0), uint32(0)

	select {
	case b := <-result:
		assert.NotNil(t, b)
		proposal, _ := nodes[0].engine.isCurrentValidator()
		groupID, anotherGroupID := uint32(0), uint32(1)
		indexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), groupID)
		anotherIndexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), anotherGroupID)
		fmt.Println(indexes)
		fmt.Println(anotherIndexes)

		for _, validatorIndex := range indexes {
			if validatorIndex != proposal.Index {
				missValidatorIndex = validatorIndex
				break
			}
		}
		for _, validatorIndex := range anotherIndexes {
			if validatorIndex != proposal.Index {
				syncValidatorIndex = validatorIndex
				break
			}
		}

		prepareBlock := &protocols.PrepareBlock{
			Epoch:         epoch,
			ViewNumber:    viewNumber,
			Block:         block,
			BlockIndex:    blockIndex,
			ProposalIndex: proposal.Index,
		}
		assert.Nil(t, nodes[0].engine.signMsgByBls(prepareBlock))
		nodes[missValidatorIndex].engine.state.AddPrepareBlock(prepareBlock)
		nodes[syncValidatorIndex].engine.state.AddPrepareBlock(prepareBlock)

		for i, validatorIndex := range indexes {
			msg := &protocols.PrepareVote{
				Epoch:          epoch,
				ViewNumber:     viewNumber,
				BlockIndex:     blockIndex,
				BlockHash:      b.Hash(),
				BlockNumber:    b.NumberU64(),
				ValidatorIndex: validatorIndex,
				ParentQC:       qc,
			}
			assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))
			if i < 17 {
				assert.Nil(t, nodes[missValidatorIndex].engine.OnPrepareVote("id", msg))
			} else {
				assert.Nil(t, nodes[syncValidatorIndex].engine.OnPrepareVote("id", msg))
			}
		}

		for i, validatorIndex := range anotherIndexes {
			msg := &protocols.PrepareVote{
				Epoch:          nodes[0].engine.state.Epoch(),
				ViewNumber:     nodes[0].engine.state.ViewNumber(),
				BlockIndex:     blockIndex,
				BlockHash:      b.Hash(),
				BlockNumber:    b.NumberU64(),
				ValidatorIndex: validatorIndex,
				ParentQC:       qc,
			}
			assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))
			if i%2 == 0 {
				assert.Nil(t, nodes[missValidatorIndex].engine.OnPrepareVote("id", msg))
			}
			assert.Nil(t, nodes[syncValidatorIndex].engine.OnPrepareVote("id", msg))
		}
	}

	// 1、missValidatorIndex节点和 syncValidatorIndex节点缺失的 votes
	getPrepareVoteV2, _ := testMissingPrepareVotev2(t, nodes[missValidatorIndex], nodes[syncValidatorIndex])

	// 2、missValidatorIndex节点 向 syncValidatorIndex节点 请求同步, syncValidatorIndex节点返回的 votes
	prepareVotesV2 := testOnGetPrepareVoteV2(t, nodes[syncValidatorIndex], getPrepareVoteV2)
	assert.Equal(t, 4, len(prepareVotesV2.Votes))
	assert.Equal(t, 0, len(prepareVotesV2.RGBlockQuorumCerts))

	// 3、missValidatorIndex节点处理同步的 votes
	assert.Equal(t, 1, len(nodes[missValidatorIndex].engine.state.FindMaxRGQuorumCerts(blockIndex)))
	assert.Equal(t, 17, nodes[missValidatorIndex].engine.state.FindMaxRGQuorumCerts(blockIndex)[0].ValidatorSet.HasLength())
	testOnPrepareVotesV2(t, nodes[missValidatorIndex], prepareVotesV2) // node2 响应 node1 同步
	assert.Equal(t, 1, len(nodes[missValidatorIndex].engine.state.FindMaxRGQuorumCerts(blockIndex)))
	assert.Equal(t, 21, nodes[missValidatorIndex].engine.state.FindMaxRGQuorumCerts(blockIndex)[0].ValidatorSet.HasLength())
}

func testMissingPrepareVotev2(t *testing.T, missNode, syncNode *TestCBFT) (*protocols.GetPrepareVoteV2, *protocols.GetPrepareVoteV2) {
	// check missNode
	request1, err := missNode.engine.MissingPrepareVote()
	assert.Nil(t, err)
	getPrepareVoteV2_1, ok := request1.(*protocols.GetPrepareVoteV2)
	assert.True(t, true, ok)
	fmt.Println(getPrepareVoteV2_1.String())
	assert.Equal(t, 2, len(getPrepareVoteV2_1.UnKnownGroups.UnKnown))
	assert.Equal(t, 20, getPrepareVoteV2_1.UnKnownGroups.UnKnownSize())

	for _, unKnown := range getPrepareVoteV2_1.UnKnownGroups.UnKnown {
		groupID := unKnown.GroupID
		unKnownSet := unKnown.UnKnownSet
		indexes, _ := missNode.engine.validatorPool.GetValidatorIndexesByGroupID(missNode.engine.state.Epoch(), groupID)
		for i, validatorIndex := range indexes {
			if groupID == 0 {
				if i >= 17 {
					assert.True(t, true, unKnownSet.GetIndex(validatorIndex))
				}
			} else if groupID == 1 {
				if i%2 != 0 {
					assert.True(t, true, unKnownSet.GetIndex(validatorIndex))
				}
			}
		}
	}
	// check syncNode
	request2, err := syncNode.engine.MissingPrepareVote()
	assert.Nil(t, err)
	getPrepareVoteV2_2, ok := request2.(*protocols.GetPrepareVoteV2)
	assert.True(t, true, ok)
	fmt.Println(getPrepareVoteV2_2.String())
	assert.Equal(t, 1, len(getPrepareVoteV2_2.UnKnownGroups.UnKnown))
	assert.Equal(t, 17, getPrepareVoteV2_2.UnKnownGroups.UnKnownSize())

	for _, unKnown := range getPrepareVoteV2_2.UnKnownGroups.UnKnown {
		groupID := unKnown.GroupID
		unKnownSet := unKnown.UnKnownSet
		indexes, _ := syncNode.engine.validatorPool.GetValidatorIndexesByGroupID(syncNode.engine.state.Epoch(), groupID)
		for i, validatorIndex := range indexes {
			if groupID == 0 {
				if i < 17 {
					assert.True(t, true, unKnownSet.GetIndex(validatorIndex))
				}
			}
		}
	}
	return getPrepareVoteV2_1, getPrepareVoteV2_2
}

func testOnGetPrepareVoteV2(t *testing.T, requested *TestCBFT, request *protocols.GetPrepareVoteV2) *protocols.PrepareVotesV2 {
	response, err := requested.engine.OnGetPrepareVoteV2("id", request)
	assert.Nil(t, err)
	prepareVotesV2, ok := response.(*protocols.PrepareVotesV2)
	assert.True(t, true, ok)
	return prepareVotesV2
}

func testOnPrepareVotesV2(t *testing.T, sync *TestCBFT, response *protocols.PrepareVotesV2) {
	assert.Nil(t, sync.engine.OnPrepareVotesV2("id", response))
}

func TestCbft_SyncRGBlockQuorumCerts(t *testing.T) {
	nodesNum := 50
	nodes := MockRGNodes(t, nodesNum)

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)
	parent := nodes[0].chain.Genesis()

	block := NewBlockWithSign(parent.Hash(), parent.NumberU64()+1, nodes[0])
	assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
	nodes[0].engine.OnSeal(block, result, nil, complete)
	<-complete

	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
	epoch, viewNumber, blockIndex := nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber(), uint32(0)
	missValidatorIndex, syncValidatorIndex := uint32(0), uint32(0)

	select {
	case b := <-result:
		assert.NotNil(t, b)
		proposal, _ := nodes[0].engine.isCurrentValidator()
		groupID, anotherGroupID := uint32(0), uint32(1)
		indexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), groupID)
		anotherIndexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), anotherGroupID)
		fmt.Println(indexes)
		fmt.Println(anotherIndexes)

		for _, validatorIndex := range indexes {
			if validatorIndex != proposal.Index {
				syncValidatorIndex = validatorIndex
				break
			}
		}
		for _, validatorIndex := range anotherIndexes {
			if validatorIndex != proposal.Index {
				missValidatorIndex = validatorIndex
				break
			}
		}

		prepareBlock := &protocols.PrepareBlock{
			Epoch:         epoch,
			ViewNumber:    viewNumber,
			Block:         block,
			BlockIndex:    blockIndex,
			ProposalIndex: proposal.Index,
		}
		assert.Nil(t, nodes[0].engine.signMsgByBls(prepareBlock))
		nodes[missValidatorIndex].engine.state.AddPrepareBlock(prepareBlock)
		nodes[syncValidatorIndex].engine.state.AddPrepareBlock(prepareBlock)

		for i, validatorIndex := range indexes {
			msg := &protocols.PrepareVote{
				Epoch:          epoch,
				ViewNumber:     viewNumber,
				BlockIndex:     blockIndex,
				BlockHash:      b.Hash(),
				BlockNumber:    b.NumberU64(),
				ValidatorIndex: validatorIndex,
				ParentQC:       qc,
			}
			assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))
			if i%2 == 0 {
				assert.Nil(t, nodes[missValidatorIndex].engine.OnPrepareVote("id", msg))
			}
			assert.Nil(t, nodes[syncValidatorIndex].engine.OnPrepareVote("id", msg))
		}

		for i, validatorIndex := range anotherIndexes {
			msg := &protocols.PrepareVote{
				Epoch:          nodes[0].engine.state.Epoch(),
				ViewNumber:     nodes[0].engine.state.ViewNumber(),
				BlockIndex:     blockIndex,
				BlockHash:      b.Hash(),
				BlockNumber:    b.NumberU64(),
				ValidatorIndex: validatorIndex,
				ParentQC:       qc,
			}
			assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))
			if i >= 8 {
				assert.Nil(t, nodes[missValidatorIndex].engine.OnPrepareVote("id", msg))
			} else {
				assert.Nil(t, nodes[syncValidatorIndex].engine.OnPrepareVote("id", msg))
			}
		}
	}

	// Sleep for a while and wait for syncValidatorIndex to SendRGBlockQuorumCert
	time.Sleep(5 * time.Second)

	// 1、missValidatorIndex节点和 syncValidatorIndex节点缺失的 votes
	getPrepareVoteV2, _ := testMissingRGBlockQuorumCerts(t, nodes[missValidatorIndex], nodes[syncValidatorIndex])

	// 2、missValidatorIndex节点 向 syncValidatorIndex节点 请求同步, syncValidatorIndex节点返回的 votes
	prepareVotesV2 := testOnGetPrepareVoteV2(t, nodes[syncValidatorIndex], getPrepareVoteV2)
	assert.Equal(t, 1, len(prepareVotesV2.RGBlockQuorumCerts))

	// 3、missValidatorIndex节点处理同步的 votes
	assert.Nil(t, nodes[missValidatorIndex].engine.state.FindRGBlockQuorumCerts(0, 0, syncValidatorIndex))
	assert.Equal(t, 1, len(nodes[missValidatorIndex].engine.state.FindMaxRGQuorumCerts(blockIndex)))
	testOnPrepareVotesV2(t, nodes[missValidatorIndex], prepareVotesV2) // node2 响应 node1 同步
	xxx := nodes[missValidatorIndex].engine.state
	fmt.Println(xxx)
	assert.NotNil(t, nodes[missValidatorIndex].engine.state.FindRGBlockQuorumCerts(0, 0, syncValidatorIndex))
	assert.Equal(t, 2, len(nodes[missValidatorIndex].engine.state.FindMaxRGQuorumCerts(blockIndex)))
}

func testMissingRGBlockQuorumCerts(t *testing.T, missNode, syncNode *TestCBFT) (*protocols.GetPrepareVoteV2, *protocols.GetPrepareVoteV2) {
	// check missNode
	request1, err := missNode.engine.MissingPrepareVote()
	assert.Nil(t, err)
	getPrepareVoteV2_1, ok := request1.(*protocols.GetPrepareVoteV2)
	assert.True(t, true, ok)
	fmt.Println(getPrepareVoteV2_1.String())
	assert.Equal(t, 2, len(getPrepareVoteV2_1.UnKnownGroups.UnKnown))
	assert.Equal(t, 20, getPrepareVoteV2_1.UnKnownGroups.UnKnownSize())

	for _, unKnown := range getPrepareVoteV2_1.UnKnownGroups.UnKnown {
		groupID := unKnown.GroupID
		unKnownSet := unKnown.UnKnownSet
		indexes, _ := missNode.engine.validatorPool.GetValidatorIndexesByGroupID(missNode.engine.state.Epoch(), groupID)
		for i, validatorIndex := range indexes {
			if groupID == 0 {
				if i%2 != 0 {
					assert.True(t, true, unKnownSet.GetIndex(validatorIndex))
				}
			} else if groupID == 1 {
				if i < 8 {
					assert.True(t, true, unKnownSet.GetIndex(validatorIndex))
				}
			}
		}
	}
	// check syncNode
	request2, err := syncNode.engine.MissingPrepareVote()
	assert.Nil(t, err)
	getPrepareVoteV2_2, ok := request2.(*protocols.GetPrepareVoteV2)
	assert.True(t, true, ok)
	fmt.Println(getPrepareVoteV2_2.String())
	assert.Equal(t, 1, len(getPrepareVoteV2_2.UnKnownGroups.UnKnown))
	assert.Equal(t, 17, getPrepareVoteV2_2.UnKnownGroups.UnKnownSize())

	for _, unKnown := range getPrepareVoteV2_2.UnKnownGroups.UnKnown {
		groupID := unKnown.GroupID
		unKnownSet := unKnown.UnKnownSet
		indexes, _ := syncNode.engine.validatorPool.GetValidatorIndexesByGroupID(syncNode.engine.state.Epoch(), groupID)
		for i, validatorIndex := range indexes {
			if groupID == 1 {
				if i >= 8 {
					assert.True(t, true, unKnownSet.GetIndex(validatorIndex))
				}
			}
		}
	}
	return getPrepareVoteV2_1, getPrepareVoteV2_2
}

func TestCbft_SyncViewChangeV2(t *testing.T) {
	nodesNum := 50
	nodes := MockRGNodes(t, nodesNum)
	ReachBlock(t, nodes, 5)

	block := nodes[0].engine.state.HighestQCBlock()
	block, qc := nodes[0].engine.blockTree.FindBlockAndQC(block.Hash(), block.NumberU64())
	epoch, viewNumber := nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()

	groupID, anotherGroupID := uint32(0), uint32(1)
	indexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), groupID)
	anotherIndexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), anotherGroupID)
	fmt.Println(indexes)
	fmt.Println(anotherIndexes)
	missValidatorIndex, syncValidatorIndex := uint32(0), uint32(0)

	for i, validatorIndex := range indexes {
		if i%2 == 0 {
			missValidatorIndex = validatorIndex
			break
		}
	}
	for i, validatorIndex := range anotherIndexes {
		if i >= 17 {
			syncValidatorIndex = validatorIndex
			break
		}
	}

	for i, validatorIndex := range indexes {
		viewChange := &protocols.ViewChange{
			Epoch:          epoch,
			ViewNumber:     viewNumber,
			BlockHash:      block.Hash(),
			BlockNumber:    block.NumberU64(),
			ValidatorIndex: validatorIndex,
			PrepareQC:      qc,
		}
		assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(viewChange))
		if i%2 == 0 {
			assert.Nil(t, nodes[missValidatorIndex].engine.OnViewChange("id", viewChange))
		} else {
			assert.Nil(t, nodes[syncValidatorIndex].engine.OnViewChange("id", viewChange))
		}
	}

	for i, validatorIndex := range anotherIndexes {
		viewChange := &protocols.ViewChange{
			Epoch:          epoch,
			ViewNumber:     viewNumber,
			BlockHash:      block.Hash(),
			BlockNumber:    block.NumberU64(),
			ValidatorIndex: validatorIndex,
			PrepareQC:      qc,
		}
		assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(viewChange))
		if i < 17 {
			assert.Nil(t, nodes[missValidatorIndex].engine.OnViewChange("id", viewChange))
		} else {
			assert.Nil(t, nodes[syncValidatorIndex].engine.OnViewChange("id", viewChange))
		}
	}

	// 1、missNode 节点缺失的 viewChanges
	getViewChangeV2, _ := testMissingViewChangev2(t, nodes[missValidatorIndex], nodes[syncValidatorIndex])

	// 2、missNode 节点 向 syncNode节点 请求同步, syncNode 节点返回的 viewChanges
	viewChangesV2 := testOnGetViewChangeV2(t, nodes[syncValidatorIndex], getViewChangeV2)
	assert.Equal(t, 4, len(viewChangesV2.VCs))
	assert.Equal(t, 0, len(viewChangesV2.RGViewChangeQuorumCerts))

	// 3、missNode节点处理同步的 viewChanges
	assert.Equal(t, uint64(0), nodes[missValidatorIndex].engine.state.ViewNumber())
	testOnViewChangeV2(t, nodes[missValidatorIndex], viewChangesV2) // node2 响应 node1 同步
	assert.Equal(t, uint64(1), nodes[missValidatorIndex].engine.state.ViewNumber())
}

func testMissingViewChangev2(t *testing.T, missNode, syncNode *TestCBFT) (*protocols.GetViewChangeV2, *protocols.GetViewChangeV2) {
	var getViewChangeV2_1, getViewChangeV2_2 *protocols.GetViewChangeV2
	var request ctypes.Message
	var err error
	// check missNode
	select {
	case <-missNode.engine.state.ViewTimeout():
		request, err = missNode.engine.MissingViewChangeNodes()
	case <-time.After(10000 * time.Millisecond):
		request, err = missNode.engine.MissingViewChangeNodes()
	}
	assert.Nil(t, err)
	getViewChangeV2_1, _ = request.(*protocols.GetViewChangeV2)
	fmt.Println(getViewChangeV2_1.String())
	assert.Equal(t, 2, len(getViewChangeV2_1.UnKnownGroups.UnKnown))
	assert.Equal(t, 20, getViewChangeV2_1.UnKnownGroups.UnKnownSize())

	for _, unKnown := range getViewChangeV2_1.UnKnownGroups.UnKnown {
		groupID := unKnown.GroupID
		unKnownSet := unKnown.UnKnownSet
		indexes, _ := missNode.engine.validatorPool.GetValidatorIndexesByGroupID(missNode.engine.state.Epoch(), groupID)
		for i, validatorIndex := range indexes {
			if groupID == 0 {
				if i%2 != 0 {
					assert.True(t, true, unKnownSet.GetIndex(validatorIndex))
				}
			} else if groupID == 1 {
				if i >= 17 {
					assert.True(t, true, unKnownSet.GetIndex(validatorIndex))
				}
			}
		}
	}

	// check syncNode
	select {
	case <-syncNode.engine.state.ViewTimeout():
		request, err = syncNode.engine.MissingViewChangeNodes()
	case <-time.After(10000 * time.Millisecond):
		request, err = syncNode.engine.MissingViewChangeNodes()
	}
	assert.Nil(t, err)
	getViewChangeV2_2, _ = request.(*protocols.GetViewChangeV2)
	fmt.Println(getViewChangeV2_2.String())
	assert.Equal(t, 2, len(getViewChangeV2_2.UnKnownGroups.UnKnown))
	assert.Equal(t, 30, getViewChangeV2_2.UnKnownGroups.UnKnownSize())

	for _, unKnown := range getViewChangeV2_2.UnKnownGroups.UnKnown {
		groupID := unKnown.GroupID
		unKnownSet := unKnown.UnKnownSet
		indexes, _ := syncNode.engine.validatorPool.GetValidatorIndexesByGroupID(syncNode.engine.state.Epoch(), groupID)
		for i, validatorIndex := range indexes {
			if groupID == 0 {
				if i%2 == 0 {
					assert.True(t, true, unKnownSet.GetIndex(validatorIndex))
				}
				if i < 17 {
					assert.True(t, true, unKnownSet.GetIndex(validatorIndex))
				}
			}
		}
	}

	return getViewChangeV2_1, getViewChangeV2_2
}

func testOnGetViewChangeV2(t *testing.T, requested *TestCBFT, request *protocols.GetViewChangeV2) *protocols.ViewChangesV2 {
	response, err := requested.engine.OnGetViewChangeV2("id", request)
	assert.Nil(t, err)
	viewChangesV2, ok := response.(*protocols.ViewChangesV2)
	assert.True(t, true, ok)
	return viewChangesV2
}

func testOnViewChangeV2(t *testing.T, sync *TestCBFT, response *protocols.ViewChangesV2) {
	assert.Nil(t, sync.engine.OnViewChangesV2("id", response))
}
