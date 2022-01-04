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

package network

import (
	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/protocols"
	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/types"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
)

// Cbft defines the network layer to use the relevant interface
// to the consensus layer.
type Cbft interface {

	// Returns the ID value of the current node.
	Node() *enode.Node

	// Return a list of all consensus nodes.
	ConsensusNodes() ([]enode.ID, error)

	// Return configuration information of CBFT consensus.
	Config() *types.Config

	// Entrance: The messages related to the consensus are entered from here.
	// The message sent from the peer node is sent to the CBFT message queue and
	// there is a loop that will distribute the incoming message.
	ReceiveMessage(msg *types.MsgInfo) error

	// ReceiveSyncMsg is used to receive messages that are synchronized from other nodes.
	ReceiveSyncMsg(msg *types.MsgInfo) error

	// Return the highest QC block number of the current node.
	HighestQCBlockBn() (uint64, common.Hash)

	// Return the highest locked block number of the current node.
	HighestLockBlockBn() (uint64, common.Hash)

	// Return the highest commit block number of the current node.
	HighestCommitBlockBn() (uint64, common.Hash)

	// Returns the missing vote.
	MissingPrepareVote() (types.Message, error)

	// Returns the node ID of the missing vote.
	MissingViewChangeNodes() (types.Message, error)

	// Returns latest status.
	LatestStatus() *protocols.GetLatestStatus

	// OnPong records net delay time.
	OnPong(nodeID string, netLatency int64) error

	// BlockExists determines if a block exists.
	BlockExists(blockNumber uint64, blockHash common.Hash) error

	// NeedGroup indicates whether grouped consensus will be used
	NeedGroup() bool

	// TODO just for log
	GetGroupByValidatorID(nodeID enode.ID) (uint32, uint32, error)
}
