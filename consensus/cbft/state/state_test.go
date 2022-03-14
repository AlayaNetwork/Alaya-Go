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

package state

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"encoding/json"
	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/common/math"
	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/protocols"
	ctypes "github.com/AlayaNetwork/Alaya-Go/consensus/cbft/types"
	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/utils"
	"github.com/AlayaNetwork/Alaya-Go/core/types"
)

func TestNewViewState(t *testing.T) {
	viewState := NewViewState(BaseMs, nil)
	viewState.ResetView(1, 1)

	assert.Equal(t, uint64(1), viewState.Epoch())
	assert.Equal(t, uint64(1), viewState.ViewNumber())
	assert.Equal(t, 0, viewState.ViewBlockSize())
	assert.Equal(t, uint32(0), viewState.NextViewBlockIndex())
	assert.Equal(t, uint32(math.MaxUint32), viewState.MaxQCIndex())
	assert.Equal(t, 0, viewState.ViewVoteSize())

	assert.Equal(t, uint64(1), viewState.view.Epoch())
	assert.Equal(t, uint64(1), viewState.view.ViewNumber())
	_, err := viewState.view.MarshalJSON()
	assert.Nil(t, err)

	assert.Equal(t, 0, viewState.HadSendPrepareVote().Len())
	assert.Equal(t, 0, viewState.PendingPrepareVote().Len())

	viewState.SetExecuting(uint32(1), true)
	_, finish := viewState.Executing()
	assert.True(t, finish)

	viewState.SetLastViewChangeQC(&ctypes.ViewChangeQC{})
	assert.NotNil(t, viewState.LastViewChangeQC())

	viewState.SetHighestCommitBlock(newBlock(1))
	viewState.SetHighestLockBlock(newBlock(2))
	viewState.SetHighestQCBlock(newBlock(3))
	assert.NotNil(t, viewState.HighestCommitBlock())
	assert.NotNil(t, viewState.HighestLockBlock())
	assert.NotNil(t, viewState.HighestQCBlock())
	assert.NotNil(t, viewState.HighestBlockString())

	_, err = viewState.MarshalJSON()
	assert.Nil(t, err)

	viewState.SetViewTimer(1)

	select {
	case <-viewState.ViewTimeout():
		assert.True(t, viewState.IsDeadline())
		assert.True(t, viewState.IsDeadline())
	}
}

func TestPrepareVoteQueue(t *testing.T) {
	queue := newPrepareVoteQueue()

	for i := 0; i < 10; i++ {
		b := &protocols.PrepareVote{BlockIndex: uint32(i)}
		queue.Push(b)
	}

	expect := uint32(0)
	for !queue.Empty() {
		assert.Equal(t, queue.Top().BlockIndex, expect)
		assert.True(t, queue.Had(expect))
		assert.False(t, queue.Had(expect+10))
		queue.Pop()
		expect++
	}

	assert.Equal(t, 0, queue.Len())
	assert.Len(t, queue.Peek(), 0)
}

func TestPrepareVotes(t *testing.T) {
	viewState := NewViewState(BaseMs, nil)

	var b *protocols.PrepareVote
	for i := 0; i < 10; i++ {
		b = &protocols.PrepareVote{BlockIndex: uint32(i)}
		viewState.viewVotes.addVote(uint32(i), b)
	}

	assert.Equal(t, 1, viewState.PrepareVoteLenByIndex(uint32(0)))
	assert.NotNil(t, viewState.FindPrepareVote(uint32(1), uint32(1)))
	assert.True(t, viewState.viewVotes.Votes[uint32(9)].hadVote(b))

	viewState.viewVotes.clear()
}

func TestViewBlocks(t *testing.T) {
	viewState := NewViewState(BaseMs, nil)

	var viewBlock *prepareViewBlock
	for i := 0; i < 10; i++ {
		viewBlock = &prepareViewBlock{
			pb: &protocols.PrepareBlock{BlockIndex: uint32(i), Block: newBlock(0)},
		}
		viewState.viewBlocks.addBlock(viewBlock)
	}
	assert.Equal(t, 10, viewState.viewBlocks.len())
	assert.Equal(t, uint32(9), viewState.viewBlocks.MaxIndex())
	assert.Equal(t, viewBlock.hash(), viewState.viewBlocks.Blocks[9].hash())
	assert.Equal(t, viewBlock.number(), viewState.viewBlocks.Blocks[9].number())

	assert.NotNil(t, viewState.ViewBlockByIndex(9))
	assert.NotNil(t, viewState.PrepareBlockByIndex(9))
	assert.Equal(t, 10, viewState.ViewBlockSize())
}

var (
	BaseMs = uint64(10000)
)

func TestViewVotes(t *testing.T) {
	viewState := NewViewState(BaseMs, nil)
	votes := viewState.viewVotes
	prepareVotes := []*protocols.PrepareVote{
		{BlockIndex: uint32(5)},
		{BlockIndex: uint32(6)},
		{BlockIndex: uint32(7)},
	}

	for i, p := range prepareVotes {
		viewState.AddPrepareVote(uint32(i), p)
		votes.addVote(uint32(i), p)
	}
	assert.Equal(t, 3, len(viewState.viewVotes.Votes))
	assert.Equal(t, uint32(7), viewState.MaxViewVoteIndex())
	assert.Len(t, viewState.AllPrepareVoteByIndex(5), 1)
	assert.Equal(t, viewState.PrepareVoteLenByIndex(uint32(len(prepareVotes))), 0)
	assert.Len(t, viewState.AllPrepareVoteByIndex(uint32(len(prepareVotes))), 0)

	votes.clear()
	assert.Len(t, votes.Votes, 0)
}

func TestNewViewQC(t *testing.T) {
	viewState := NewViewState(BaseMs, nil)
	viewQCs := viewState.viewQCs

	for i := uint32(0); i < 10; i++ {
		viewState.AddQC(&ctypes.QuorumCert{BlockIndex: i})
	}

	assert.Equal(t, viewState.MaxQCIndex(), uint32(9))
	for i := uint32(0); i < 10; i++ {
		assert.NotNil(t, viewQCs.index(i))
	}

	for i := uint32(10); i < 20; i++ {
		assert.Nil(t, viewQCs.index(i))
	}

	viewQCs.clear()
	assert.Equal(t, viewQCs.len(), 0)
}

func newBlock(number uint64) *types.Block {
	header := &types.Header{
		Number:     big.NewInt(int64(number)),
		ParentHash: common.Hash{},
		Time:       uint64(time.Now().UnixNano()),
		Extra:      nil,
	}
	block := types.NewBlockWithHeader(header)
	return block
}

func TestNewViewBlock(t *testing.T) {
	viewState := NewViewState(BaseMs, nil)
	for i := uint64(0); i < 10; i++ {
		viewState.AddQCBlock(newBlock(i), &ctypes.QuorumCert{BlockNumber: i, BlockIndex: uint32(i)})
	}

	for i := uint32(0); i < 10; i++ {
		block, _ := viewState.ViewBlockAndQC(i)
		assert.NotNil(t, block)
	}

	block, qc := viewState.ViewBlockAndQC(11)
	assert.Nil(t, block)
	assert.Nil(t, qc)
}

func TestNewViewChanges(t *testing.T) {
	viewState := NewViewState(BaseMs, nil)

	var v *protocols.ViewChange
	for i := uint32(0); i < 10; i++ {
		v = &protocols.ViewChange{ValidatorIndex: i}
		viewState.AddViewChange(i, v)
	}

	assert.Equal(t, 10, viewState.ViewChangeLen())
	assert.Equal(t, 10, len(viewState.AllViewChange()))
	assert.Equal(t, uint32(9), viewState.ViewChangeByIndex(9).ValidatorIndex)
}

func newQuorumCert(blockIndex uint32, set *utils.BitArray) *ctypes.QuorumCert {
	return &ctypes.QuorumCert{
		Epoch:        1,
		ViewNumber:   1,
		BlockHash:    common.Hash{},
		BlockNumber:  1,
		BlockIndex:   blockIndex,
		Signature:    ctypes.Signature{},
		ValidatorSet: set,
	}
}

func unmarshalBitArray(bitArrayStr string) *utils.BitArray {
	var ba *utils.BitArray
	json.Unmarshal([]byte(bitArrayStr), &ba)
	return ba
}

func marshalBitArray(arr *utils.BitArray) string {
	if b, err := json.Marshal(arr); err == nil {
		return string(b)
	}
	return ""
}

func TestViewRGBlockQuorumCerts(t *testing.T) {
	testCases := []struct {
		blockIndex      uint32
		groupID         uint32
		validatorIndex  uint32
		validatorSetStr string
	}{
		{0, 1, 11, `"x_x_x_"`},
		{0, 2, 22, `"x_x_x_"`},
		{0, 1, 12, `"x_x_x_"`},
		{1, 3, 33, `"x_x_x_"`},
		{1, 5, 55, `"x_x_x_"`},
		{2, 1, 11, `"x_x_x_"`},
		{2, 2, 22, `"x_x_x_"`},
		{2, 2, 23, `"x_xxx_"`},
		{0, 1, 12, `"x_x_x_"`}, // duplicate data
	}

	v := newViewRGBlockQuorumCerts()
	for _, c := range testCases {
		v.AddRGBlockQuorumCerts(c.blockIndex, &protocols.RGBlockQuorumCert{
			GroupID:        c.groupID,
			BlockQC:        newQuorumCert(c.blockIndex, unmarshalBitArray(c.validatorSetStr)),
			ValidatorIndex: c.validatorIndex,
		})
	}

	//fmt.Println(v.String())

	assert.Equal(t, 2, v.RGBlockQuorumCertsLen(0, 1))
	assert.Equal(t, 1, v.RGBlockQuorumCertsLen(0, 2))
	assert.Equal(t, 1, v.RGBlockQuorumCertsLen(1, 3))
	assert.Equal(t, 1, v.RGBlockQuorumCertsLen(1, 5))
	assert.Equal(t, 1, v.RGBlockQuorumCertsLen(2, 1))
	assert.Equal(t, 2, v.RGBlockQuorumCertsLen(2, 2))
	assert.Equal(t, 0, v.RGBlockQuorumCertsLen(2, 3))
	assert.Equal(t, 0, v.RGBlockQuorumCertsLen(3, 1))

	assert.Nil(t, v.FindRGBlockQuorumCerts(3, 1, 1))
	assert.Nil(t, v.FindRGBlockQuorumCerts(0, 3, 1))
	assert.Nil(t, v.FindRGBlockQuorumCerts(0, 1, 13))
	assert.NotNil(t, v.FindRGBlockQuorumCerts(1, 3, 33))
	assert.Equal(t, uint32(11), v.FindRGBlockQuorumCerts(0, 1, 11).ValidatorIndex)

	//assert.Equal(t, []uint32{11, 12}, v.RGBlockQuorumCertsIndexes(0, 1))
	assert.Equal(t, 2, len(v.RGBlockQuorumCertsIndexes(0, 1)))
	assert.Equal(t, `"x_xxx_"`, marshalBitArray(v.FindMaxGroupRGBlockQuorumCert(2, 2).BlockQC.ValidatorSet))
}

func TestSelectedRGBlockQuorumCerts(t *testing.T) {
	testCases := []struct {
		blockIndex      uint32
		groupID         uint32
		ValidatorSetStr string
	}{
		{0, 1, `"x"`},
		{0, 1, `"xxxx__"`},
		{0, 1, `"xx"`},
		{0, 1, `"x_x_x_"`},
		{0, 1, `"xx__x_"`},
		{0, 1, `"xxx_x_"`},

		{1, 1, `"x"`},
		{1, 1, `"xxxx"`},
		{1, 1, `"xx"`},
		{1, 1, `"x_x_x_"`},
		{1, 1, `"xx__x_"`},
		{1, 1, `"xxx_x_"`},
		{1, 1, `"xxxxx_"`}, // contains all

		{0, 2, `"x______"`},
		{0, 2, `"_x_____"`},
		{0, 2, `"__x____"`},
		{0, 2, `"___x___"`},
		{0, 2, `"____x__"`},
		{0, 2, `"_____xx"`}, // exceed the limit,but more sign,accept

		{1, 2, `"x"`},
		{1, 2, `"xxxxx_"`}, // contains all
		{1, 2, `"xx"`},
		{1, 2, `"x_x_x_"`},
		{1, 2, `"xx__x_"`},
		{1, 2, `"xxx_x_"`},
		{1, 2, `"xxxx"`},

		{2, 2, `"x_x_x_"`},
		{2, 2, `"x_xx_x"`},
		{2, 2, `"_x_xx_"`},
		{2, 2, `"xx__x_"`},
		{2, 2, `"_x__xx"`},
	}

	s := newSelectedRGBlockQuorumCerts()
	for _, c := range testCases {
		s.AddRGQuorumCerts(c.blockIndex, c.groupID, &ctypes.QuorumCert{
			BlockIndex:   c.blockIndex,
			ValidatorSet: unmarshalBitArray(c.ValidatorSetStr),
		}, newQuorumCert(c.blockIndex, unmarshalBitArray(c.ValidatorSetStr)))
	}

	//fmt.Println(s.String())

	assert.Equal(t, 2, len(s.FindRGQuorumCerts(0, 1)))
	assert.Equal(t, 1, len(s.FindRGQuorumCerts(1, 1)))
	assert.Equal(t, 6, s.RGQuorumCertsLen(0, 2))
	assert.Equal(t, 1, s.RGQuorumCertsLen(1, 2))
	assert.Equal(t, 5, s.RGQuorumCertsLen(2, 2))
	assert.Equal(t, 0, s.RGQuorumCertsLen(0, 3))
	assert.Equal(t, 0, s.RGQuorumCertsLen(3, 0))

	max, parentQC := s.FindMaxGroupRGQuorumCert(0, 1)
	assert.Equal(t, uint32(0), max.BlockIndex)
	assert.Equal(t, `"xxxx__"`, marshalBitArray(max.ValidatorSet))
	assert.Equal(t, uint32(0), parentQC.BlockIndex)

	maxs := s.FindMaxRGQuorumCerts(1)
	assert.Equal(t, uint32(1), maxs[1].BlockIndex)
	assert.Equal(t, `"xxxxx_"`, marshalBitArray(maxs[1].ValidatorSet))

	// test merge vote
	s.MergePrepareVote(0, 2, &protocols.PrepareVote{
		BlockIndex:     0,
		ValidatorIndex: 4,
	})
	assert.Equal(t, 5, s.RGQuorumCertsLen(0, 2)) // after merge, Removes the contained element, changing length from 6 to 5
	max, parentQC = s.FindMaxGroupRGQuorumCert(0, 2)
	assert.Equal(t, `"____xxx"`, marshalBitArray(max.ValidatorSet))
	// test merge vote
	s.MergePrepareVote(2, 2, &protocols.PrepareVote{
		BlockIndex:     2,
		ValidatorIndex: 2,
	})
	assert.Equal(t, 4, s.RGQuorumCertsLen(2, 2))
	// test merge vote
	s.MergePrepareVote(2, 2, &protocols.PrepareVote{
		BlockIndex:     2,
		ValidatorIndex: 0,
	})
	assert.Equal(t, 3, s.RGQuorumCertsLen(2, 2))
	// test merge vote
	s.MergePrepareVote(2, 2, &protocols.PrepareVote{
		BlockIndex:     2,
		ValidatorIndex: 3,
	})
	assert.Equal(t, 1, s.RGQuorumCertsLen(2, 2))
	max, _ = s.FindMaxGroupRGQuorumCert(2, 2)
	assert.Equal(t, `"xxxxxx"`, marshalBitArray(max.ValidatorSet))
	//a := s.FindRGQuorumCerts(2, 2)
	//for _, v := range a {
	//	fmt.Println(marshalBitArray(v.ValidatorSet))
	//}
}

func TestViewRGViewChangeQuorumCerts(t *testing.T) {
	testCases := []struct {
		groupID         uint32
		validatorIndex  uint32
		validatorSetStr string
	}{
		{1, 11, `"x_x_x_"`},
		{2, 22, `"x_x_x_"`},
		{1, 12, `"x_x_x_"`},
		{3, 33, `"x_x_x_"`},
		{5, 55, `"x_x_x_"`},
		{1, 11, `"x_x_x_"`}, // duplicate data
		{2, 22, `"x_x_x_"`}, // duplicate data
		{2, 23, `"x_xxx_"`},
		{1, 12, `"x_x_x_"`}, // duplicate data
	}

	v := newViewRGViewChangeQuorumCerts()
	for _, c := range testCases {
		v.AddRGViewChangeQuorumCerts(&protocols.RGViewChangeQuorumCert{
			GroupID: c.groupID,
			ViewChangeQC: &ctypes.ViewChangeQC{
				QCs: []*ctypes.ViewChangeQuorumCert{
					{ValidatorSet: unmarshalBitArray(c.validatorSetStr)},
				},
			},
			ValidatorIndex: c.validatorIndex,
		})
	}

	//fmt.Println(v.String())

	assert.Equal(t, 2, v.RGViewChangeQuorumCertsLen(1))
	assert.Equal(t, 2, v.RGViewChangeQuorumCertsLen(2))
	assert.Equal(t, 1, v.RGViewChangeQuorumCertsLen(3))
	assert.Equal(t, 0, v.RGViewChangeQuorumCertsLen(4))
	assert.Equal(t, 1, v.RGViewChangeQuorumCertsLen(5))

	rg := v.FindRGViewChangeQuorumCerts(1, 12)
	assert.NotNil(t, rg)
	assert.Equal(t, uint32(12), rg.ValidatorIndex)

	assert.Equal(t, 2, len(v.RGViewChangeQuorumCertsIndexes(2)))
	assert.Equal(t, 4, v.FindMaxRGViewChangeQuorumCert(2).ViewChangeQC.HasLength())
}

func TestSelectedRGViewChangeQuorumCerts(t *testing.T) {
	testCases := []struct {
		groupID         uint32
		blockNumber     int64
		ValidatorSetStr string
	}{
		{0, 1, `"x___________"`},
		{0, 1, `"xxxx________"`},
		{0, 1, `"xx__________"`},
		{0, 1, `"x_x_x_______"`},
		{0, 1, `"xx__x_______"`},
		{0, 1, `"xxx_x_______"`},
		{0, 2, `"______x_____"`},
		{0, 2, `"______xxxx__"`},
		{0, 2, `"______xx____"`},
		{0, 2, `"______x_x_x_"`},
		{0, 2, `"______xx__x_"`},
		{0, 2, `"______xxx_x_"`},

		{2, 1, `"x___________"`},
		{2, 1, `"xxxx________"`},
		{2, 1, `"xx__________"`},
		{2, 1, `"x_x_x_______"`},
		{2, 1, `"xx__x_______"`},
		{2, 1, `"xxx_x_______"`},
		{2, 3, `"______x_____"`},
		{2, 3, `"______xxxx__"`},
		{2, 3, `"______xx____"`},
		{2, 3, `"______x_x_x_"`},
		{2, 3, `"______xx__x_"`},
		{2, 3, `"______xxx_x_"`},
	}

	s := newSelectedRGViewChangeQuorumCerts()
	for _, c := range testCases {
		hash := common.BigToHash(big.NewInt(c.blockNumber))
		rgqcs := map[common.Hash]*ctypes.ViewChangeQuorumCert{
			hash: {
				BlockNumber:  uint64(c.blockNumber),
				BlockHash:    hash,
				Signature:    ctypes.Signature{},
				ValidatorSet: unmarshalBitArray(c.ValidatorSetStr),
			},
		}
		prepareQCs := map[common.Hash]*ctypes.QuorumCert{
			hash: {
				BlockNumber: uint64(c.blockNumber),
				BlockHash:   hash,
			},
		}
		s.AddRGViewChangeQuorumCerts(c.groupID, rgqcs, prepareQCs)
	}

	//fmt.Println(s.String())

	assert.Equal(t, 2, len(s.findRGQuorumCerts(0)))
	assert.Equal(t, 0, len(s.findRGQuorumCerts(1)))
	assert.Equal(t, 0, s.RGViewChangeQuorumCertsLen(1))
	assert.Equal(t, 2, s.RGViewChangeQuorumCertsLen(2))
	viewChangeQC, prepareQCs := s.FindMaxGroupRGViewChangeQuorumCert(0)
	assert.Equal(t, 2, len(viewChangeQC.QCs))
	for _, qc := range viewChangeQC.QCs {
		if qc.BlockNumber == uint64(1) {
			assert.Equal(t, `"xxxx________"`, marshalBitArray(qc.ValidatorSet))
		} else if qc.BlockNumber == uint64(2) {
			assert.Equal(t, `"______xxxx__"`, marshalBitArray(qc.ValidatorSet))
		}
	}

	assert.Equal(t, 2, len(prepareQCs.QCs))
	//assert.Equal(t, uint64(1), prepareQCs.QCs[0].BlockNumber)
	//assert.Equal(t, uint64(2), prepareQCs.QCs[1].BlockNumber)

	maxs := s.FindMaxRGViewChangeQuorumCert()
	assert.Equal(t, 2, len(maxs))
	//if maxs[1].QCs[0].BlockNumber == uint64(1) {
	//	assert.Equal(t, `"xxxx________"`, marshalBitArray(maxs[1].QCs[0].ValidatorSet))
	//} else if maxs[1].QCs[0].BlockNumber == uint64(3) {
	//	assert.Equal(t, `"______xxx_x_"`, marshalBitArray(maxs[1].QCs[0].ValidatorSet))
	//}

	// merge viewchange
	s.MergeViewChange(0, &protocols.ViewChange{
		BlockNumber:    3,
		BlockHash:      common.BigToHash(big.NewInt(int64(3))),
		ValidatorIndex: 6,
	}, 12)
	assert.Equal(t, 3, s.RGViewChangeQuorumCertsLen(0))
	viewChangeQC, prepareQCs = s.FindMaxGroupRGViewChangeQuorumCert(0)
	assert.Equal(t, 3, len(viewChangeQC.QCs))
	for _, qc := range viewChangeQC.QCs {
		if qc.BlockNumber == uint64(1) {
			assert.Equal(t, `"xxxx________"`, marshalBitArray(qc.ValidatorSet))
		} else if qc.BlockNumber == uint64(2) {
			assert.Equal(t, `"______xxxx__"`, marshalBitArray(qc.ValidatorSet))
		} else if qc.BlockNumber == uint64(3) {
			assert.Equal(t, `"______x_____"`, marshalBitArray(qc.ValidatorSet))
		}
	}
	assert.Equal(t, 3, len(prepareQCs.QCs))

	// merge viewchange
	s.MergeViewChange(0, &protocols.ViewChange{
		BlockNumber:    1,
		BlockHash:      common.BigToHash(big.NewInt(int64(1))),
		ValidatorIndex: 3,
	}, 12)
	//fmt.Println(s.String())
	m := s.findRGQuorumCerts(0) // after merge, Removes the contained element, changing length from 2 to 1
	v, ok := m[common.BigToHash(big.NewInt(int64(1)))]
	assert.True(t, true, ok)
	assert.Equal(t, 1, len(v))
	assert.Equal(t, `"xxxxx_______"`, marshalBitArray(v[0].ValidatorSet))
}
