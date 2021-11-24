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

package xutil

import (
	"bytes"
	"fmt"

	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"

	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/crypto"
	"github.com/AlayaNetwork/Alaya-Go/x/xcom"
)

func NodeId2Addr(nodeId enode.IDv0) (common.NodeAddress, error) {
	if pk, err := nodeId.Pubkey(); nil != err {
		return common.ZeroNodeAddr, err
	} else {
		return common.NodeAddress(crypto.PubkeyToAddress(*pk)), nil
	}
}

// The ProgramVersion: Major.Minor.Patch eg. 1.1.0
// Calculate the LargeVersion
// eg: 1.1.x ==> 1.1.0
func CalcVersion(programVersion uint32) uint32 {
	programVersion = programVersion >> 8
	return programVersion << 8
}

func IsWorker(extra []byte) bool {
	return len(extra) > 32 && len(extra[32:]) >= common.ExtraSeal && bytes.Equal(extra[32:97], make([]byte, common.ExtraSeal))
}

// eg. 65536 => 1.0.0
func ProgramVersion2Str(programVersion uint32) string {
	if programVersion == 0 {
		return "0.0.0"
	}
	major := programVersion << 8
	major = major >> 24

	minor := programVersion << 16
	minor = minor >> 24

	patch := programVersion << 24
	patch = patch >> 24

	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

// CalcBlocksEachEpoch return how many blocks per epoch
func CalcBlocksEachEpoch(version uint32) uint64 {
	return xcom.ConsensusSize(version) * xcom.EpochSize(version)
}

func EstimateConsensusRoundsForGov(seconds uint64, version uint32) uint64 {
	//v0.7.5, hard code 1 second for block interval for estimating.
	blockInterval := uint64(1)
	return seconds / (blockInterval * xcom.ConsensusSize(version))
}

func EstimateEndVotingBlockForParaProposal(blockNumber uint64, seconds uint64, version uint32) uint64 {
	consensusSize := xcom.ConsensusSize(version)
	epochMaxDuration := xcom.MaxEpochMinutes() //minutes
	//estimate how many consensus rounds in a epoch.
	consensusRoundsEachEpoch := epochMaxDuration * 60 / (xcom.Interval() * consensusSize)
	blocksEachEpoch := consensusRoundsEachEpoch * consensusSize

	//v0.7.5, hard code 1 second for block interval for estimating.
	blockInterval := uint64(1)
	durationEachEpoch := blocksEachEpoch * blockInterval

	epochRounds := seconds / durationEachEpoch
	return blockNumber + blocksEachEpoch - blockNumber%blocksEachEpoch + epochRounds*blocksEachEpoch
}

// calculate the Epoch number by blockNumber
func CalculateEpoch(blockNumber uint64, version uint32) uint64 {
	size := CalcBlocksEachEpoch(version)
	return calculateQuotient(blockNumber, size)
}

// calculate the Consensus number by blockNumber
func CalculateRound(blockNumber uint64, version uint32) uint64 {
	size := xcom.ConsensusSize(version)
	return calculateQuotient(blockNumber, size)
}

func calculateQuotient(blockNumber, size uint64) uint64 {
	var res uint64
	div := blockNumber / size
	mod := blockNumber % size
	switch {
	// first consensus round
	case div == 0:
		res = 1
	case div > 0 && mod == 0:
		res = div
	case div > 0 && mod > 0:
		res = div + 1
	}

	return res
}

func InNodeIDList(nodeID enode.IDv0, nodeIDList []enode.IDv0) bool {
	for _, v := range nodeIDList {
		if nodeID == v {
			return true
		}
	}
	return false
}

func InHashList(hash common.Hash, hashList []common.Hash) bool {
	for _, v := range hashList {
		if hash == v {
			return true
		}
	}
	return false
}

// end-voting-block = the end block of a consensus period - electionDistance, end-voting-block must be a Consensus Election block
func CalEndVotingBlock(blockNumber uint64, endVotingRounds uint64, version uint32) uint64 {
	electionDistance := xcom.ElectionDistance()
	consensusSize := xcom.ConsensusSize(version)
	return blockNumber + consensusSize - blockNumber%consensusSize + endVotingRounds*consensusSize - electionDistance
}

// active-block = the begin of a consensus period, so, it is possible that active-block also is the begin of a epoch.
func CalActiveBlock(endVotingBlock uint64) uint64 {
	return endVotingBlock + xcom.ElectionDistance() + 1
}

// IsBeginOfEpoch returns true if current block is the first block of a Epoch
func IsBeginOfEpoch(blockNumber uint64, version uint32) bool {
	return isEpochBeginOrEnd(blockNumber, true, version)
}

func IsEndOfEpoch(blockNumber uint64, version uint32) bool {
	return isEpochBeginOrEnd(blockNumber, false, version)
}

func isEpochBeginOrEnd(blockNumber uint64, checkBegin bool, version uint32) bool {
	size := CalcBlocksEachEpoch(version)
	mod := blockNumber % size
	if checkBegin {
		return mod == 1
	} else {
		//check end
		return mod == 0
	}
}

// IsBeginOfConsensus returns true if current block is the first block of a Consensus Cycle
func IsBeginOfConsensus(blockNumber uint64, version uint32) bool {
	size := xcom.ConsensusSize(version)
	mod := blockNumber % size
	return mod == 1
}

func IsElection(blockNumber uint64, version uint32) bool {
	tmp := blockNumber + xcom.ElectionDistance()
	mod := tmp % xcom.ConsensusSize(version)
	return mod == 0
}
