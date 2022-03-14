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
	"strings"
	"testing"
	"time"
)

const ActiveVersion = 4352

const (
	insufficientSignatures        = "insufficient signatures"
	notGroupMember                = "the message sender is not a member of the group"
	includeNonGroupMembers        = "signers include non-group members"
	invalidConsensusSignature     = "verify consensus sign failed"
	invalidRGBlockQuorumCert      = "bls verifies signature fail"
	invalidRGViewChangeQuorumCert = "verify viewchange qc failed"
)

func MockRGNodes(t *testing.T, num int) []*TestCBFT {
	pk, sk, nodes := GenerateCbftNode(num)
	engines := make([]*TestCBFT, 0)

	for i := 0; i < num; i++ {
		e := MockNode(pk[i], sk[i], nodes, 10000, 10)
		e.engine.MockActiveVersion(ActiveVersion)
		assert.Nil(t, e.Start())
		engines = append(engines, e)
	}
	return engines
}

func TestRGBlockQuorumCert(t *testing.T) {
	nodes := MockRGNodes(t, 4)

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)
	parent := nodes[0].chain.Genesis()

	block := NewBlockWithSign(parent.Hash(), parent.NumberU64()+1, nodes[0])
	assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
	nodes[0].engine.OnSeal(block, result, nil, complete)
	<-complete

	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
	select {
	case b := <-result:
		assert.NotNil(t, b)
		groupVotes := make(map[uint32]*protocols.PrepareVote)
		blockIndex := uint32(0)

		for i := 0; i < len(nodes); i++ {
			msg := &protocols.PrepareVote{
				Epoch:          nodes[0].engine.state.Epoch(),
				ViewNumber:     nodes[0].engine.state.ViewNumber(),
				BlockIndex:     blockIndex,
				BlockHash:      b.Hash(),
				BlockNumber:    b.NumberU64(),
				ValidatorIndex: uint32(i),
				ParentQC:       qc,
			}
			assert.Nil(t, nodes[i].engine.signMsgByBls(msg))
			groupVotes[uint32(i)] = msg
		}
		rgqc := nodes[0].engine.generatePrepareQC(groupVotes)

		validatorIndex := 1
		groupID := uint32(0)
		msg := &protocols.RGBlockQuorumCert{
			GroupID:        groupID,
			BlockQC:        rgqc,
			ValidatorIndex: uint32(validatorIndex),
		}
		assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))
		assert.Nil(t, nodes[0].engine.OnRGBlockQuorumCert("id", msg), fmt.Sprintf("number:%d", b.NumberU64()))

		// Sleep for a while and wait for Node-0 to SendRGBlockQuorumCert
		time.Sleep(5 * time.Second)

		assert.NotNil(t, nodes[0].engine.state.FindRGBlockQuorumCerts(blockIndex, groupID, uint32(validatorIndex)))
		// The block is already QC, so node-0 need not to SendRGBlockQuorumCert
		assert.Equal(t, 1, nodes[0].engine.state.RGBlockQuorumCertsLen(blockIndex, groupID))
		assert.Equal(t, 1, len(nodes[0].engine.state.RGBlockQuorumCertsIndexes(blockIndex, groupID)))
		assert.True(t, true, nodes[0].engine.state.HadSendRGBlockQuorumCerts(blockIndex))

		assert.Equal(t, 1, nodes[0].engine.state.SelectRGQuorumCertsLen(blockIndex, groupID))
		assert.Equal(t, 1, len(nodes[0].engine.state.FindMaxRGQuorumCerts(blockIndex)))
	}
}

func TestRGBlockQuorumCert_insufficientSignatures(t *testing.T) {
	nodes := MockRGNodes(t, 4)

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)
	parent := nodes[0].chain.Genesis()

	block := NewBlockWithSign(parent.Hash(), parent.NumberU64()+1, nodes[0])
	assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
	nodes[0].engine.OnSeal(block, result, nil, complete)
	<-complete

	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
	select {
	case b := <-result:
		assert.NotNil(t, b)
		groupVotes := make(map[uint32]*protocols.PrepareVote)
		blockIndex := uint32(0)

		for i := 0; i < len(nodes)-2; i++ {
			msg := &protocols.PrepareVote{
				Epoch:          nodes[0].engine.state.Epoch(),
				ViewNumber:     nodes[0].engine.state.ViewNumber(),
				BlockIndex:     blockIndex,
				BlockHash:      b.Hash(),
				BlockNumber:    b.NumberU64(),
				ValidatorIndex: uint32(i),
				ParentQC:       qc,
			}
			assert.Nil(t, nodes[i].engine.signMsgByBls(msg))
			groupVotes[uint32(i)] = msg
		}
		rgqc := nodes[0].engine.generatePrepareQC(groupVotes)

		validatorIndex := 1
		groupID := uint32(0)
		msg := &protocols.RGBlockQuorumCert{
			GroupID:        groupID,
			BlockQC:        rgqc,
			ValidatorIndex: uint32(validatorIndex),
		}
		assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))
		err := nodes[0].engine.OnRGBlockQuorumCert("id", msg)
		assert.True(t, strings.HasPrefix(err.Error(), insufficientSignatures))
	}
}

func TestRGBlockQuorumCert_notGroupMember(t *testing.T) {
	nodes := MockRGNodes(t, 4)

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)
	parent := nodes[0].chain.Genesis()

	block := NewBlockWithSign(parent.Hash(), parent.NumberU64()+1, nodes[0])
	assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
	nodes[0].engine.OnSeal(block, result, nil, complete)
	<-complete

	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
	select {
	case b := <-result:
		assert.NotNil(t, b)
		groupVotes := make(map[uint32]*protocols.PrepareVote)
		blockIndex := uint32(0)

		for i := 0; i < len(nodes); i++ {
			msg := &protocols.PrepareVote{
				Epoch:          nodes[0].engine.state.Epoch(),
				ViewNumber:     nodes[0].engine.state.ViewNumber(),
				BlockIndex:     blockIndex,
				BlockHash:      b.Hash(),
				BlockNumber:    b.NumberU64(),
				ValidatorIndex: uint32(i),
				ParentQC:       qc,
			}
			assert.Nil(t, nodes[i].engine.signMsgByBls(msg))
			groupVotes[uint32(i)] = msg
		}
		rgqc := nodes[0].engine.generatePrepareQC(groupVotes)

		validatorIndex := 1
		groupID := uint32(0)
		msg := &protocols.RGBlockQuorumCert{
			GroupID:        groupID,
			BlockQC:        rgqc,
			ValidatorIndex: uint32(len(nodes)),
		}
		assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))
		err := nodes[0].engine.OnRGBlockQuorumCert("id", msg)
		assert.True(t, strings.HasPrefix(err.Error(), notGroupMember))
	}
}

func TestRGBlockQuorumCert_includeNonGroupMembers(t *testing.T) {
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
	select {
	case b := <-result:
		assert.NotNil(t, b)
		groupVotes := make(map[uint32]*protocols.PrepareVote)
		blockIndex := uint32(0)

		groupID := uint32(0)
		indexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), groupID)

		for i := 0; i < len(nodes); i++ {
			msg := &protocols.PrepareVote{
				Epoch:          nodes[0].engine.state.Epoch(),
				ViewNumber:     nodes[0].engine.state.ViewNumber(),
				BlockIndex:     blockIndex,
				BlockHash:      b.Hash(),
				BlockNumber:    b.NumberU64(),
				ValidatorIndex: uint32(i),
				ParentQC:       qc,
			}
			assert.Nil(t, nodes[i].engine.signMsgByBls(msg))
			groupVotes[uint32(i)] = msg
		}
		rgqc := nodes[0].engine.generatePrepareQC(groupVotes)

		validatorIndex := indexes[0]
		msg := &protocols.RGBlockQuorumCert{
			GroupID:        groupID,
			BlockQC:        rgqc,
			ValidatorIndex: validatorIndex,
		}
		assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))
		err := nodes[0].engine.OnRGBlockQuorumCert("id", msg)
		assert.True(t, strings.HasPrefix(err.Error(), includeNonGroupMembers))
	}
}

func TestRGBlockQuorumCert_invalidConsensusSignature(t *testing.T) {
	nodes := MockRGNodes(t, 4)

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)
	parent := nodes[0].chain.Genesis()

	block := NewBlockWithSign(parent.Hash(), parent.NumberU64()+1, nodes[0])
	assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
	nodes[0].engine.OnSeal(block, result, nil, complete)
	<-complete

	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
	select {
	case b := <-result:
		assert.NotNil(t, b)
		groupVotes := make(map[uint32]*protocols.PrepareVote)
		blockIndex := uint32(0)

		for i := 0; i < len(nodes); i++ {
			msg := &protocols.PrepareVote{
				Epoch:          nodes[0].engine.state.Epoch(),
				ViewNumber:     nodes[0].engine.state.ViewNumber(),
				BlockIndex:     blockIndex,
				BlockHash:      b.Hash(),
				BlockNumber:    b.NumberU64(),
				ValidatorIndex: uint32(i),
				ParentQC:       qc,
			}
			assert.Nil(t, nodes[i].engine.signMsgByBls(msg))
			groupVotes[uint32(i)] = msg
		}
		rgqc := nodes[0].engine.generatePrepareQC(groupVotes)

		validatorIndex := 1
		groupID := uint32(0)
		msg := &protocols.RGBlockQuorumCert{
			GroupID:        groupID,
			BlockQC:        rgqc,
			ValidatorIndex: uint32(validatorIndex),
		}
		assert.Nil(t, nodes[validatorIndex+1].engine.signMsgByBls(msg))
		err := nodes[0].engine.OnRGBlockQuorumCert("id", msg)
		assert.True(t, strings.HasPrefix(err.Error(), invalidConsensusSignature))
	}
}

func TestRGBlockQuorumCert_invalidRGBlockQuorumCert(t *testing.T) {
	nodes := MockRGNodes(t, 4)

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)
	parent := nodes[0].chain.Genesis()

	block := NewBlockWithSign(parent.Hash(), parent.NumberU64()+1, nodes[0])
	assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
	nodes[0].engine.OnSeal(block, result, nil, complete)
	<-complete

	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
	select {
	case b := <-result:
		assert.NotNil(t, b)
		groupVotes := make(map[uint32]*protocols.PrepareVote)
		blockIndex := uint32(0)

		for i := 0; i < len(nodes); i++ {
			msg := &protocols.PrepareVote{
				Epoch:          nodes[0].engine.state.Epoch(),
				ViewNumber:     nodes[0].engine.state.ViewNumber(),
				BlockIndex:     blockIndex,
				BlockHash:      b.Hash(),
				BlockNumber:    b.NumberU64(),
				ValidatorIndex: uint32(i),
				ParentQC:       qc,
			}
			if i < len(nodes)-1 {
				assert.Nil(t, nodes[i+1].engine.signMsgByBls(msg))
			} else {
				assert.Nil(t, nodes[i].engine.signMsgByBls(msg))
			}
			groupVotes[uint32(i)] = msg
		}
		rgqc := nodes[0].engine.generatePrepareQC(groupVotes)

		validatorIndex := 1
		groupID := uint32(0)
		msg := &protocols.RGBlockQuorumCert{
			GroupID:        groupID,
			BlockQC:        rgqc,
			ValidatorIndex: uint32(validatorIndex),
		}
		assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))
		err := nodes[0].engine.OnRGBlockQuorumCert("id", msg)
		assert.True(t, strings.HasPrefix(err.Error(), invalidRGBlockQuorumCert))
	}
}

func TestRGBlockQuorumCert_richBlockQuorumCert(t *testing.T) {
	nodes := MockRGNodes(t, 4)

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)
	parent := nodes[0].chain.Genesis()

	block := NewBlockWithSign(parent.Hash(), parent.NumberU64()+1, nodes[0])
	assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
	nodes[0].engine.OnSeal(block, result, nil, complete)
	<-complete

	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
	select {
	case b := <-result:
		assert.NotNil(t, b)
		groupVotes := make(map[uint32]*protocols.PrepareVote)
		blockIndex := uint32(0)

		for i := 0; i < len(nodes); i++ {
			msg := &protocols.PrepareVote{
				Epoch:          nodes[0].engine.state.Epoch(),
				ViewNumber:     nodes[0].engine.state.ViewNumber(),
				BlockIndex:     blockIndex,
				BlockHash:      b.Hash(),
				BlockNumber:    b.NumberU64(),
				ValidatorIndex: uint32(i),
				ParentQC:       qc,
			}
			assert.Nil(t, nodes[i].engine.signMsgByBls(msg))
			if i == 0 {
				//assert.Nil(t, nodes[0].engine.OnPrepareVote("id", msg), fmt.Sprintf("number:%d", b.NumberU64()))
			} else {
				groupVotes[uint32(i)] = msg
			}
		}
		rgqc := nodes[0].engine.generatePrepareQC(groupVotes)

		validatorIndex := 1
		groupID := uint32(0)
		msg := &protocols.RGBlockQuorumCert{
			GroupID:        groupID,
			BlockQC:        rgqc,
			ValidatorIndex: uint32(validatorIndex),
		}
		assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))

		assert.Nil(t, nodes[0].engine.OnRGBlockQuorumCert("id", msg), fmt.Sprintf("number:%d", b.NumberU64()))
		assert.Equal(t, 3, nodes[0].engine.state.FindRGBlockQuorumCerts(blockIndex, groupID, uint32(validatorIndex)).BlockQC.ValidatorSet.HasLength())
		selectqc, _ := nodes[0].engine.state.FindMaxGroupRGQuorumCert(blockIndex, groupID)
		assert.Equal(t, 4, selectqc.ValidatorSet.HasLength())
	}
}

func TestRGBlockQuorumCert_mergeVoteToQuorumCerts(t *testing.T) {
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
	select {
	case b := <-result:
		assert.NotNil(t, b)
		groupVotes := make(map[uint32]*protocols.PrepareVote)
		blockIndex := uint32(0)
		groupID := uint32(0)
		var lastVote *protocols.PrepareVote

		indexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), groupID)

		for i, validatorIndex := range indexes {
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
			if i == len(indexes)-1 {
				lastVote = msg
			} else {
				groupVotes[validatorIndex] = msg
			}
		}
		rgqc := nodes[0].engine.generatePrepareQC(groupVotes)

		validatorIndex := indexes[0]

		msg := &protocols.RGBlockQuorumCert{
			GroupID:        groupID,
			BlockQC:        rgqc,
			ValidatorIndex: validatorIndex,
		}
		assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))

		assert.Nil(t, nodes[0].engine.OnRGBlockQuorumCert("id", msg), fmt.Sprintf("number:%d", b.NumberU64()))
		assert.Equal(t, 24, nodes[0].engine.state.FindMaxGroupRGBlockQuorumCert(blockIndex, groupID).BlockQC.ValidatorSet.HasLength())
		selectqc, _ := nodes[0].engine.state.FindMaxGroupRGQuorumCert(blockIndex, groupID)
		assert.Equal(t, 24, selectqc.ValidatorSet.HasLength())

		assert.Nil(t, nodes[0].engine.OnPrepareVote("id", lastVote), fmt.Sprintf("number:%d", b.NumberU64()))
		selectqc, _ = nodes[0].engine.state.FindMaxGroupRGQuorumCert(blockIndex, groupID)
		assert.Equal(t, 25, selectqc.ValidatorSet.HasLength())
	}
}

func TestRGViewChangeQuorumCert(t *testing.T) {
	nodesNum := 50
	nodes := MockRGNodes(t, nodesNum)
	ReachBlock(t, nodes, 5)

	block := nodes[0].engine.state.HighestQCBlock()
	block, qc := nodes[0].engine.blockTree.FindBlockAndQC(block.Hash(), block.NumberU64())
	epoch, view := nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()

	groupID := uint32(0)
	groupViewChanges := make(map[uint32]*protocols.ViewChange)
	indexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), groupID)

	for _, validatorIndex := range indexes {
		viewChange := &protocols.ViewChange{
			Epoch:          epoch,
			ViewNumber:     view,
			BlockHash:      block.Hash(),
			BlockNumber:    block.NumberU64(),
			ValidatorIndex: validatorIndex,
			PrepareQC:      qc,
		}
		assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(viewChange))
		groupViewChanges[validatorIndex] = viewChange
	}

	rgqc := nodes[0].engine.generateViewChangeQC(groupViewChanges)

	primaryIndex := indexes[0]
	validatorIndex := indexes[1]
	msg := &protocols.RGViewChangeQuorumCert{
		GroupID:        groupID,
		ViewChangeQC:   rgqc,
		ValidatorIndex: uint32(validatorIndex),
		PrepareQCs: &ctypes.PrepareQCs{
			QCs: []*ctypes.QuorumCert{qc},
		},
	}
	assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))
	assert.Nil(t, nodes[primaryIndex].engine.OnRGViewChangeQuorumCert("id", msg))

	// Sleep for a while and wait for Node-primaryIndex to SendRGViewChangeQuorumCert
	time.Sleep(5 * time.Second)

	assert.NotNil(t, nodes[primaryIndex].engine.state.FindRGViewChangeQuorumCerts(groupID, uint32(validatorIndex)))
	assert.Equal(t, 2, nodes[primaryIndex].engine.state.RGViewChangeQuorumCertsLen(groupID))
	assert.Equal(t, 2, len(nodes[primaryIndex].engine.state.RGViewChangeQuorumCertsIndexes(groupID)))
	assert.True(t, true, nodes[primaryIndex].engine.state.HadSendRGViewChangeQuorumCerts(view))

	assert.Equal(t, 1, nodes[primaryIndex].engine.state.SelectRGViewChangeQuorumCertsLen(groupID))
	assert.Equal(t, 1, len(nodes[primaryIndex].engine.state.FindMaxRGViewChangeQuorumCerts()))
}

func TestRGViewChangeQuorumCert_insufficientSignatures(t *testing.T) {
	nodesNum := 50
	nodes := MockRGNodes(t, nodesNum)
	ReachBlock(t, nodes, 5)

	block := nodes[0].engine.state.HighestQCBlock()
	block, qc := nodes[0].engine.blockTree.FindBlockAndQC(block.Hash(), block.NumberU64())
	epoch, view := nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()

	groupID := uint32(0)
	groupViewChanges := make(map[uint32]*protocols.ViewChange)
	indexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), groupID)

	for i := 0; i < len(indexes)/2; i++ {
		viewChange := &protocols.ViewChange{
			Epoch:          epoch,
			ViewNumber:     view,
			BlockHash:      block.Hash(),
			BlockNumber:    block.NumberU64(),
			ValidatorIndex: indexes[i],
			PrepareQC:      qc,
		}
		assert.Nil(t, nodes[indexes[i]].engine.signMsgByBls(viewChange))
		groupViewChanges[indexes[i]] = viewChange
	}

	rgqc := nodes[0].engine.generateViewChangeQC(groupViewChanges)

	validatorIndex := indexes[1]
	msg := &protocols.RGViewChangeQuorumCert{
		GroupID:        groupID,
		ViewChangeQC:   rgqc,
		ValidatorIndex: validatorIndex,
		PrepareQCs: &ctypes.PrepareQCs{
			QCs: []*ctypes.QuorumCert{qc},
		},
	}
	assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))

	err := nodes[0].engine.OnRGViewChangeQuorumCert("id", msg)
	assert.True(t, strings.HasPrefix(err.Error(), insufficientSignatures))
}

func TestRGViewChangeQuorumCert_notGroupMember(t *testing.T) {
	nodesNum := 50
	nodes := MockRGNodes(t, nodesNum)
	ReachBlock(t, nodes, 5)

	block := nodes[0].engine.state.HighestQCBlock()
	block, qc := nodes[0].engine.blockTree.FindBlockAndQC(block.Hash(), block.NumberU64())
	epoch, view := nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()

	groupID, anotherGroupID := uint32(0), uint32(1)
	groupViewChanges := make(map[uint32]*protocols.ViewChange)
	indexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), groupID)
	anotherIndexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), anotherGroupID)

	for _, validatorIndex := range indexes {
		viewChange := &protocols.ViewChange{
			Epoch:          epoch,
			ViewNumber:     view,
			BlockHash:      block.Hash(),
			BlockNumber:    block.NumberU64(),
			ValidatorIndex: validatorIndex,
			PrepareQC:      qc,
		}
		assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(viewChange))
		groupViewChanges[validatorIndex] = viewChange
	}

	rgqc := nodes[0].engine.generateViewChangeQC(groupViewChanges)

	validatorIndex := anotherIndexes[0]
	msg := &protocols.RGViewChangeQuorumCert{
		GroupID:        groupID,
		ViewChangeQC:   rgqc,
		ValidatorIndex: validatorIndex,
		PrepareQCs: &ctypes.PrepareQCs{
			QCs: []*ctypes.QuorumCert{qc},
		},
	}
	assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))

	err := nodes[0].engine.OnRGViewChangeQuorumCert("id", msg)
	assert.True(t, strings.HasPrefix(err.Error(), notGroupMember))
}

func TestRGViewChangeQuorumCert_includeNonGroupMembers(t *testing.T) {
	nodesNum := 50
	nodes := MockRGNodes(t, nodesNum)
	ReachBlock(t, nodes, 5)

	block := nodes[0].engine.state.HighestQCBlock()
	block, qc := nodes[0].engine.blockTree.FindBlockAndQC(block.Hash(), block.NumberU64())
	epoch, view := nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()

	groupID := uint32(0)
	groupViewChanges := make(map[uint32]*protocols.ViewChange)
	indexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), groupID)

	for i := 0; i < len(nodes); i++ {
		viewChange := &protocols.ViewChange{
			Epoch:          epoch,
			ViewNumber:     view,
			BlockHash:      block.Hash(),
			BlockNumber:    block.NumberU64(),
			ValidatorIndex: uint32(i),
			PrepareQC:      qc,
		}
		assert.Nil(t, nodes[i].engine.signMsgByBls(viewChange))
		groupViewChanges[uint32(i)] = viewChange
	}

	rgqc := nodes[0].engine.generateViewChangeQC(groupViewChanges)

	validatorIndex := indexes[0]
	msg := &protocols.RGViewChangeQuorumCert{
		GroupID:        groupID,
		ViewChangeQC:   rgqc,
		ValidatorIndex: validatorIndex,
		PrepareQCs: &ctypes.PrepareQCs{
			QCs: []*ctypes.QuorumCert{qc},
		},
	}
	assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))

	err := nodes[0].engine.OnRGViewChangeQuorumCert("id", msg)
	assert.True(t, strings.HasPrefix(err.Error(), includeNonGroupMembers))
}

func TestRGViewChangeQuorumCert_invalidConsensusSignature(t *testing.T) {
	nodes := MockRGNodes(t, 4)
	ReachBlock(t, nodes, 5)

	block := nodes[0].engine.state.HighestQCBlock()
	block, qc := nodes[0].engine.blockTree.FindBlockAndQC(block.Hash(), block.NumberU64())
	epoch, view := nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()

	groupID := uint32(0)
	groupViewChanges := make(map[uint32]*protocols.ViewChange)

	for i := 0; i < len(nodes); i++ {
		viewChange := &protocols.ViewChange{
			Epoch:          epoch,
			ViewNumber:     view,
			BlockHash:      block.Hash(),
			BlockNumber:    block.NumberU64(),
			ValidatorIndex: uint32(i),
			PrepareQC:      qc,
		}
		assert.Nil(t, nodes[i].engine.signMsgByBls(viewChange))
		groupViewChanges[uint32(i)] = viewChange
	}

	rgqc := nodes[0].engine.generateViewChangeQC(groupViewChanges)

	validatorIndex := 1
	msg := &protocols.RGViewChangeQuorumCert{
		GroupID:        groupID,
		ViewChangeQC:   rgqc,
		ValidatorIndex: uint32(validatorIndex),
		PrepareQCs: &ctypes.PrepareQCs{
			QCs: []*ctypes.QuorumCert{qc},
		},
	}
	assert.Nil(t, nodes[validatorIndex+1].engine.signMsgByBls(msg))

	err := nodes[0].engine.OnRGViewChangeQuorumCert("id", msg)
	assert.True(t, strings.HasPrefix(err.Error(), invalidConsensusSignature))
}

func TestRGViewChangeQuorumCert_invalidRGViewChangeQuorumCert(t *testing.T) {
	nodes := MockRGNodes(t, 4)
	ReachBlock(t, nodes, 5)

	block := nodes[0].engine.state.HighestQCBlock()
	block, qc := nodes[0].engine.blockTree.FindBlockAndQC(block.Hash(), block.NumberU64())
	epoch, view := nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()

	groupID := uint32(0)
	groupViewChanges := make(map[uint32]*protocols.ViewChange)

	for i := 0; i < len(nodes); i++ {
		viewChange := &protocols.ViewChange{
			Epoch:          epoch,
			ViewNumber:     view,
			BlockHash:      block.Hash(),
			BlockNumber:    block.NumberU64(),
			ValidatorIndex: uint32(i),
			PrepareQC:      qc,
		}
		if i < len(nodes)-1 {
			assert.Nil(t, nodes[i+1].engine.signMsgByBls(viewChange))
		} else {
			assert.Nil(t, nodes[i].engine.signMsgByBls(viewChange))
		}
		groupViewChanges[uint32(i)] = viewChange
	}

	rgqc := nodes[0].engine.generateViewChangeQC(groupViewChanges)

	validatorIndex := 1
	msg := &protocols.RGViewChangeQuorumCert{
		GroupID:        groupID,
		ViewChangeQC:   rgqc,
		ValidatorIndex: uint32(validatorIndex),
		PrepareQCs: &ctypes.PrepareQCs{
			QCs: []*ctypes.QuorumCert{qc},
		},
	}
	assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))

	err := nodes[0].engine.OnRGViewChangeQuorumCert("id", msg)
	assert.True(t, strings.HasPrefix(err.Error(), invalidRGViewChangeQuorumCert))
}

func TestRGViewChangeQuorumCert_richViewChangeQuorumCert(t *testing.T) {
	nodesNum := 50
	nodes := MockRGNodes(t, nodesNum)
	ReachBlock(t, nodes, 5)

	block := nodes[0].engine.state.HighestQCBlock()
	block, qc := nodes[0].engine.blockTree.FindBlockAndQC(block.Hash(), block.NumberU64())
	epoch, view := nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()

	groupID := uint32(0)
	groupViewChanges := make(map[uint32]*protocols.ViewChange)
	indexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), groupID)
	primaryIndex := indexes[0]
	validatorIndex := indexes[1]

	for i := 0; i < len(indexes); i++ {
		viewChange := &protocols.ViewChange{
			Epoch:          epoch,
			ViewNumber:     view,
			BlockHash:      block.Hash(),
			BlockNumber:    block.NumberU64(),
			ValidatorIndex: indexes[i],
			PrepareQC:      qc,
		}
		assert.Nil(t, nodes[indexes[i]].engine.signMsgByBls(viewChange))
		if i < 8 {
			assert.Nil(t, nodes[primaryIndex].engine.OnViewChange("id", viewChange))
		} else {
			groupViewChanges[indexes[i]] = viewChange
		}
	}

	rgqc := nodes[0].engine.generateViewChangeQC(groupViewChanges)

	msg := &protocols.RGViewChangeQuorumCert{
		GroupID:        groupID,
		ViewChangeQC:   rgqc,
		ValidatorIndex: validatorIndex,
		PrepareQCs: &ctypes.PrepareQCs{
			QCs: []*ctypes.QuorumCert{qc},
		},
	}
	assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))

	assert.Nil(t, nodes[primaryIndex].engine.OnRGViewChangeQuorumCert("id", msg))
	assert.Equal(t, 17, nodes[primaryIndex].engine.state.FindRGViewChangeQuorumCerts(groupID, validatorIndex).ViewChangeQC.HasLength())
	selectqc, _ := nodes[primaryIndex].engine.state.FindMaxGroupRGViewChangeQuorumCert(groupID)
	assert.Equal(t, 25, selectqc.HasLength())
}

func TestRGViewChangeQuorumCert_mergeViewChangeToViewChangeQuorumCerts(t *testing.T) {
	nodesNum := 50
	nodes := MockRGNodes(t, nodesNum)
	ReachBlock(t, nodes, 5)

	block := nodes[0].engine.state.HighestQCBlock()
	block, qc := nodes[0].engine.blockTree.FindBlockAndQC(block.Hash(), block.NumberU64())
	epoch, view := nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()

	groupID := uint32(0)
	groupViewChanges := make(map[uint32]*protocols.ViewChange)
	lastViewChanges := make(map[uint32]*protocols.ViewChange)
	indexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), groupID)
	primaryIndex := indexes[0]
	validatorIndex := indexes[1]

	for i := 0; i < len(indexes); i++ {
		viewChange := &protocols.ViewChange{
			Epoch:          epoch,
			ViewNumber:     view,
			BlockHash:      block.Hash(),
			BlockNumber:    block.NumberU64(),
			ValidatorIndex: indexes[i],
			PrepareQC:      qc,
		}
		assert.Nil(t, nodes[indexes[i]].engine.signMsgByBls(viewChange))
		if i < 8 {
			lastViewChanges[indexes[i]] = viewChange
		} else {
			groupViewChanges[indexes[i]] = viewChange
		}
	}

	rgqc := nodes[0].engine.generateViewChangeQC(groupViewChanges)

	msg := &protocols.RGViewChangeQuorumCert{
		GroupID:        groupID,
		ViewChangeQC:   rgqc,
		ValidatorIndex: validatorIndex,
		PrepareQCs: &ctypes.PrepareQCs{
			QCs: []*ctypes.QuorumCert{qc},
		},
	}
	assert.Nil(t, nodes[validatorIndex].engine.signMsgByBls(msg))

	assert.Nil(t, nodes[primaryIndex].engine.OnRGViewChangeQuorumCert("id", msg))
	assert.Equal(t, 17, nodes[primaryIndex].engine.state.FindRGViewChangeQuorumCerts(groupID, validatorIndex).ViewChangeQC.HasLength())
	selectqc, _ := nodes[primaryIndex].engine.state.FindMaxGroupRGViewChangeQuorumCert(groupID)
	assert.Equal(t, 17, selectqc.HasLength())

	for _, viewChange := range lastViewChanges {
		assert.Nil(t, nodes[primaryIndex].engine.OnViewChange("id", viewChange))
	}
	assert.Equal(t, 17, nodes[primaryIndex].engine.state.FindRGViewChangeQuorumCerts(groupID, validatorIndex).ViewChangeQC.HasLength())
	selectqc, _ = nodes[primaryIndex].engine.state.FindMaxGroupRGViewChangeQuorumCert(groupID)
	assert.Equal(t, 25, selectqc.HasLength())
}
