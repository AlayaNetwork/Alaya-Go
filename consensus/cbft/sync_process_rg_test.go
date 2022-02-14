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
	"github.com/AlayaNetwork/Alaya-Go/core/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCbft_MissingPrepareVoteV2(t *testing.T) {
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
		blockIndex := uint32(0)

		groupID, anotherGroupID := uint32(0), uint32(1)
		indexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), groupID)
		anotherIndexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), anotherGroupID)
		fmt.Println(indexes)
		fmt.Println(anotherIndexes)

		for i, validatorIndex := range indexes {
			if i < 17 {
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
				nodes[0].engine.OnPrepareVote("id", msg)
			}
		}

		for i, validatorIndex := range anotherIndexes {
			if i%2 == 0 {
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
				nodes[0].engine.OnPrepareVote("id", msg)
			}
		}
		request, err := nodes[0].engine.MissingPrepareVote()
		assert.Nil(t, err)
		getPrepareVoteV2, ok := request.(*protocols.GetPrepareVoteV2)
		assert.True(t, true, ok)
		fmt.Println(getPrepareVoteV2.String())
		assert.Equal(t, 2, len(getPrepareVoteV2.UnKnownGroups.UnKnown))
		assert.Equal(t, 20, getPrepareVoteV2.UnKnownGroups.UnKnownSize())

		for _, unKnown := range getPrepareVoteV2.UnKnownGroups.UnKnown {
			groupID := unKnown.GroupID
			unKnownSet := unKnown.UnKnownSet
			indexes, _ := nodes[0].engine.validatorPool.GetValidatorIndexesByGroupID(nodes[0].engine.state.Epoch(), groupID)
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
	}
}

func TestCbft_MissingViewChangeNodesV2(t *testing.T) {

}

func TestCbft_OnGetPrepareVoteV2(t *testing.T) {

}

func TestCbft_OnGetViewChangeV2(t *testing.T) {

}
