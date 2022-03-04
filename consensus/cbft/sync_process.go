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
	"container/list"
	"fmt"
	"sort"
	"time"

	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/state"

	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/network"
	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/protocols"
	ctypes "github.com/AlayaNetwork/Alaya-Go/consensus/cbft/types"
	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/utils"
	"github.com/AlayaNetwork/Alaya-Go/core/types"
)

var syncPrepareVotesInterval = 8 * time.Second

// Get the block from the specified connection, get the block into the fetcher, and execute the block CBFT update state machine
func (cbft *Cbft) fetchBlock(id string, hash common.Hash, number uint64, qc *ctypes.QuorumCert) {
	if cbft.fetcher.Len() != 0 {
		cbft.log.Trace("Had fetching block")
		return
	}
	highestQC := cbft.state.HighestQCBlock()
	if highestQC.NumberU64()+3 < number {
		cbft.log.Debug(fmt.Sprintf("Local state too low, local.highestQC:%s,%d, remote.msg:%s,%d", highestQC.Hash().String(), highestQC.NumberU64(), hash.String(), number))
		return
	}

	baseBlockHash, baseBlockNumber := common.Hash{}, uint64(0)
	var parentBlock *types.Block

	if highestQC.NumberU64() < number {
		parentBlock = highestQC
		baseBlockHash, baseBlockNumber = parentBlock.Hash(), parentBlock.NumberU64()
	} else if highestQC.NumberU64() == number {
		parentBlock = cbft.state.HighestLockBlock()
		baseBlockHash, baseBlockNumber = parentBlock.Hash(), parentBlock.NumberU64()
	} else {
		cbft.log.Trace("No suitable block need to request")
		return
	}

	if qc != nil {
		if err := cbft.verifyPrepareQC(number, hash, qc); err != nil {
			cbft.log.Warn(fmt.Sprintf("Verify prepare qc failed, error:%s", err.Error()))
			return
		}
	}

	match := func(msg ctypes.Message) bool {
		_, ok := msg.(*protocols.QCBlockList)
		return ok
	}

	executor := func(msg ctypes.Message) {
		defer func() {
			cbft.log.Debug("Close fetching")
			//utils.SetFalse(&cbft.fetching)
		}()
		cbft.log.Debug("Receive QCBlockList", "msg", msg.String())
		if blockList, ok := msg.(*protocols.QCBlockList); ok {
			// Execution block
			for _, block := range blockList.Blocks {
				if block.ParentHash() != parentBlock.Hash() {
					cbft.log.Debug("Response block's is error",
						"blockHash", block.Hash(), "blockNumber", block.NumberU64(),
						"parentHash", parentBlock.Hash(), "parentNumber", parentBlock.NumberU64())
					return
				}
				if err := cbft.blockCache.Execute(block, parentBlock); err != nil {
					cbft.log.Error("Execute block failed", "hash", block.Hash(), "number", block.NumberU64(), "error", err)
					return
				}
				parentBlock = block
			}

			// Update the results to the CBFT state machine
			cbft.asyncCallCh <- func() {
				for i, block := range blockList.Blocks {
					if err := cbft.verifyPrepareQC(block.NumberU64(), block.Hash(), blockList.QC[i]); err != nil {
						cbft.log.Error("Verify block prepare qc failed", "hash", block.Hash(), "number", block.NumberU64(), "error", err)
						return
					}
				}
				if err := cbft.OnInsertQCBlock(blockList.Blocks, blockList.QC); err != nil {
					cbft.log.Error("Insert block failed", "error", err)
				}
			}
			if blockList.ForkedBlocks == nil || len(blockList.ForkedBlocks) == 0 {
				cbft.log.Trace("No forked block need to handle")
				return
			}
			// Remove local forks that already exist.
			filteredForkedBlocks := make([]*types.Block, 0)
			filteredForkedQCs := make([]*ctypes.QuorumCert, 0)
			//localForkedBlocks, _ := cbft.blockTree.FindForkedBlocksAndQCs(parentBlock.Hash(), parentBlock.NumberU64())
			localForkedBlocks, _ := cbft.blockTree.FindBlocksAndQCs(parentBlock.NumberU64())

			if localForkedBlocks != nil && len(localForkedBlocks) > 0 {
				cbft.log.Debug("LocalForkedBlocks", "number", localForkedBlocks[0].NumberU64(), "hash", localForkedBlocks[0].Hash().TerminalString())
			}

			for i, forkedBlock := range blockList.ForkedBlocks {
				for _, localForkedBlock := range localForkedBlocks {
					if forkedBlock.NumberU64() == localForkedBlock.NumberU64() && forkedBlock.Hash() != localForkedBlock.Hash() {
						filteredForkedBlocks = append(filteredForkedBlocks, forkedBlock)
						filteredForkedQCs = append(filteredForkedQCs, blockList.ForkedQC[i])
						break
					}
				}
			}
			if filteredForkedBlocks != nil && len(filteredForkedBlocks) > 0 {
				cbft.log.Debug("FilteredForkedBlocks", "number", filteredForkedBlocks[0].NumberU64(), "hash", filteredForkedBlocks[0].Hash().TerminalString())
			}

			// Execution forked block.
			//var forkedParentBlock *types.Block
			for _, forkedBlock := range filteredForkedBlocks {
				if forkedBlock.NumberU64() != parentBlock.NumberU64() {
					cbft.log.Error("Invalid forked block", "lastParentNumber", parentBlock.NumberU64(), "forkedBlockNumber", forkedBlock.NumberU64())
					break
				}
				//for _, block := range blockList.Blocks {
				//	if block.Hash() == forkedBlock.ParentHash() && block.NumberU64() == forkedBlock.NumberU64()-1 {
				//		forkedParentBlock = block
				//		break
				//	}
				//}
				//if forkedParentBlock != nil {
				//	break
				//}
			}

			// Verify forked block and execute.
			for _, forkedBlock := range filteredForkedBlocks {
				parentBlock := cbft.blockTree.FindBlockByHash(forkedBlock.ParentHash())
				if parentBlock == nil {
					cbft.log.Debug("Response forked block's is error",
						"blockHash", forkedBlock.Hash(), "blockNumber", forkedBlock.NumberU64())
					return
				}
				//if forkedParentBlock == nil || forkedBlock.ParentHash() != forkedParentBlock.Hash() {
				//	cbft.log.Debug("Response forked block's is error",
				//		"blockHash", forkedBlock.Hash(), "blockNumber", forkedBlock.NumberU64(),
				//		"parentHash", parentBlock.Hash(), "parentNumber", parentBlock.NumberU64())
				//	return
				//}

				if err := cbft.blockCache.Execute(forkedBlock, parentBlock); err != nil {
					cbft.log.Error("Execute forked block failed", "hash", forkedBlock.Hash(), "number", forkedBlock.NumberU64(), "error", err)
					return
				}
			}

			cbft.asyncCallCh <- func() {
				for i, forkedBlock := range filteredForkedBlocks {
					if err := cbft.verifyPrepareQC(forkedBlock.NumberU64(), forkedBlock.Hash(), blockList.ForkedQC[i]); err != nil {
						cbft.log.Error("Verify forked block prepare qc failed", "hash", forkedBlock.Hash(), "number", forkedBlock.NumberU64(), "error", err)
						return
					}
				}
				if err := cbft.OnInsertQCBlock(filteredForkedBlocks, filteredForkedQCs); err != nil {
					cbft.log.Error("Insert forked block failed", "error", err)
				}
			}
		}
	}

	expire := func() {
		cbft.log.Debug("Fetch timeout, close fetching", "targetId", id, "baseBlockHash", baseBlockHash, "baseBlockNumber", baseBlockNumber)
		utils.SetFalse(&cbft.fetching)
	}

	cbft.log.Debug("Start fetching", "id", id, "basBlockHash", baseBlockHash, "baseBlockNumber", baseBlockNumber)

	//utils.SetTrue(&cbft.fetching)
	cbft.fetcher.AddTask(id, match, executor, expire)
	cbft.network.Send(id, &protocols.GetQCBlockList{BlockHash: baseBlockHash, BlockNumber: baseBlockNumber})
}

// Obtain blocks that are not in the local according to the proposed block
func (cbft *Cbft) prepareBlockFetchRules(id string, pb *protocols.PrepareBlock) {
	if pb.Block.NumberU64() > cbft.state.HighestQCBlock().NumberU64() {
		for i := uint32(0); i <= pb.BlockIndex; i++ {
			b := cbft.state.ViewBlockByIndex(i)
			if b == nil {
				cbft.SyncPrepareBlock(id, cbft.state.Epoch(), cbft.state.ViewNumber(), i)
			}
		}
	}
}

// Get votes and blocks that are not available locally based on the height of the vote
func (cbft *Cbft) prepareVoteFetchRules(id string, msg ctypes.ConsensusMsg) {
	// Greater than QC+1 means the vote is behind
	if msg.BlockNum() > cbft.state.HighestQCBlock().NumberU64()+1 {
		for i := uint32(0); i <= msg.BlockIndx(); i++ {
			b, qc := cbft.state.ViewBlockAndQC(i)
			if b == nil {
				cbft.SyncPrepareBlock(id, cbft.state.Epoch(), cbft.state.ViewNumber(), i)
			} else if qc == nil {
				cbft.SyncBlockQuorumCert(id, b.NumberU64(), b.Hash(), i)
			}
		}
	}
}

// OnGetPrepareBlock handles the  message type of GetPrepareBlockMsg.
func (cbft *Cbft) OnGetPrepareBlock(id string, msg *protocols.GetPrepareBlock) error {
	if msg.Epoch == cbft.state.Epoch() && msg.ViewNumber == cbft.state.ViewNumber() {
		prepareBlock := cbft.state.PrepareBlockByIndex(msg.BlockIndex)
		if prepareBlock != nil {
			cbft.log.Debug("Send PrepareBlock", "peer", id, "prepareBlock", prepareBlock.String())
			cbft.network.Send(id, prepareBlock)
		}

		_, qc := cbft.state.ViewBlockAndQC(msg.BlockIndex)
		if qc != nil {
			cbft.log.Debug("Send BlockQuorumCert on GetPrepareBlock", "peer", id, "qc", qc.String())
			cbft.network.Send(id, &protocols.BlockQuorumCert{BlockQC: qc})
		}
	}
	return nil
}

// OnGetBlockQuorumCert handles the message type of GetBlockQuorumCertMsg.
func (cbft *Cbft) OnGetBlockQuorumCert(id string, msg *protocols.GetBlockQuorumCert) error {
	_, qc := cbft.blockTree.FindBlockAndQC(msg.BlockHash, msg.BlockNumber)
	if qc != nil {
		cbft.network.Send(id, &protocols.BlockQuorumCert{BlockQC: qc})
	}
	return nil
}

// OnBlockQuorumCert handles the message type of BlockQuorumCertMsg.
func (cbft *Cbft) OnBlockQuorumCert(id string, msg *protocols.BlockQuorumCert) error {
	cbft.log.Debug("Receive BlockQuorumCert", "peer", id, "msg", msg.String())
	if msg.BlockQC.Epoch != cbft.state.Epoch() || msg.BlockQC.ViewNumber != cbft.state.ViewNumber() {
		cbft.log.Trace("Receive BlockQuorumCert response failed", "local.epoch", cbft.state.Epoch(), "local.viewNumber", cbft.state.ViewNumber(), "msg", msg.String())
		return fmt.Errorf("msg is not match current state")
	}

	if _, qc := cbft.blockTree.FindBlockAndQC(msg.BlockQC.BlockHash, msg.BlockQC.BlockNumber); qc != nil {
		cbft.log.Trace("Block has exist", "msg", msg.String())
		return fmt.Errorf("block already exists")
	}

	// If blockQC comes the block must exist
	block := cbft.state.ViewBlockByIndex(msg.BlockQC.BlockIndex)
	if block == nil {
		cbft.log.Debug("Block not exist", "msg", msg.String())
		return fmt.Errorf("block not exist")
	}
	if err := cbft.verifyPrepareQC(block.NumberU64(), block.Hash(), msg.BlockQC); err != nil {
		cbft.log.Error("Failed to verify prepareQC", "err", err.Error())
		return &authFailedError{err}
	}

	cbft.log.Info("Receive blockQuorumCert, try insert prepareQC", "view", cbft.state.ViewString(), "blockQuorumCert", msg.BlockQC.String())
	cbft.insertPrepareQC(msg.BlockQC)
	return nil
}

// OnGetQCBlockList handles the message type of GetQCBlockListMsg.
func (cbft *Cbft) OnGetQCBlockList(id string, msg *protocols.GetQCBlockList) error {
	highestQC := cbft.state.HighestQCBlock()

	if highestQC.NumberU64() > msg.BlockNumber+3 ||
		(highestQC.Hash() == msg.BlockHash && highestQC.NumberU64() == msg.BlockNumber) {
		cbft.log.Debug(fmt.Sprintf("Receive GetQCBlockList failed, local.highestQC:%s,%d, msg:%s", highestQC.Hash().TerminalString(), highestQC.NumberU64(), msg.String()))
		return fmt.Errorf("peer state too low")
	}

	lock := cbft.state.HighestLockBlock()
	commit := cbft.state.HighestCommitBlock()

	qcs := make([]*ctypes.QuorumCert, 0)
	blocks := make([]*types.Block, 0)

	if commit.ParentHash() == msg.BlockHash {
		block, qc := cbft.blockTree.FindBlockAndQC(commit.Hash(), commit.NumberU64())
		qcs = append(qcs, qc)
		blocks = append(blocks, block)
	}

	if lock.ParentHash() == msg.BlockHash || commit.ParentHash() == msg.BlockHash {
		block, qc := cbft.blockTree.FindBlockAndQC(lock.Hash(), lock.NumberU64())
		qcs = append(qcs, qc)
		blocks = append(blocks, block)
	}
	if highestQC.ParentHash() == msg.BlockHash || lock.ParentHash() == msg.BlockHash || commit.ParentHash() == msg.BlockHash {
		block, qc := cbft.blockTree.FindBlockAndQC(highestQC.Hash(), highestQC.NumberU64())
		qcs = append(qcs, qc)
		blocks = append(blocks, block)
	}

	// If the height of the QC exists in the blocktree,
	// collecting forked blocks.
	forkedQCs := make([]*ctypes.QuorumCert, 0)
	forkedBlocks := make([]*types.Block, 0)
	if highestQC.ParentHash() == msg.BlockHash {
		bs, qcs := cbft.blockTree.FindForkedBlocksAndQCs(highestQC.Hash(), highestQC.NumberU64())
		if bs != nil && qcs != nil && len(bs) != 0 && len(qcs) != 0 {
			cbft.log.Debug("Find forked block and return", "forkedBlockLen", len(bs), "forkedQcLen", len(qcs))
			forkedBlocks = append(forkedBlocks, bs...)
			forkedQCs = append(forkedQCs, qcs...)
		}
	}

	if len(qcs) != 0 {
		cbft.network.Send(id, &protocols.QCBlockList{QC: qcs, Blocks: blocks, ForkedBlocks: forkedBlocks, ForkedQC: forkedQCs})
		cbft.log.Debug("Send QCBlockList", "len", len(qcs))
	}
	return nil
}

// OnGetPrepareVoteV2 is responsible for processing the business logic
// of the GetPrepareVote message. It will synchronously return a
// PrepareVotes message to the sender.
func (cbft *Cbft) OnGetPrepareVoteV2(id string, msg *protocols.GetPrepareVoteV2) (ctypes.Message, error) {
	cbft.log.Debug("Received message on OnGetPrepareVoteV2", "from", id, "msgHash", msg.MsgHash(), "message", msg.String())
	responseVotesCounter.Inc(1)
	if msg.Epoch == cbft.state.Epoch() && msg.ViewNumber == cbft.state.ViewNumber() {
		// If the block has already QC, that response QC instead of votes.
		// Avoid the sender spent a lot of time to verifies PrepareVote msg.
		_, qc := cbft.state.ViewBlockAndQC(msg.BlockIndex)
		if qc != nil {
			blockQuorumCert := &protocols.BlockQuorumCert{BlockQC: qc}
			cbft.network.Send(id, &protocols.BlockQuorumCert{BlockQC: qc})
			cbft.log.Debug("Send BlockQuorumCert", "peer", id, "qc", qc.String())
			return blockQuorumCert, nil
		}

		if len(msg.UnKnownGroups.UnKnown) > 0 {
			validatorLen := cbft.currentValidatorLen()
			threshold := cbft.threshold(validatorLen)
			remain := threshold - (validatorLen - msg.UnKnownGroups.UnKnownSize())

			// Defining an array for receiving PrepareVote.
			votes := make([]*protocols.PrepareVote, 0)
			RGBlockQuorumCerts := make([]*protocols.RGBlockQuorumCert, 0)

			prepareVoteMap := cbft.state.AllPrepareVoteByIndex(msg.BlockIndex)

			for _, un := range msg.UnKnownGroups.UnKnown {
				if un.UnKnownSet.Size() != uint32(validatorLen) {
					cbft.log.Error("Invalid request params,UnKnownGroups is not a specified length", "peer", id, "groupID", un.GroupID, "unKnownSet", un.UnKnownSet.String(), "validatorLen", validatorLen)
					break // Do not continue processing the request
				}
				// Limit response votes
				if remain <= 0 {
					break
				}
				groupLen := cbft.groupLen(cbft.state.Epoch(), un.GroupID)
				groupThreshold := cbft.groupThreshold(cbft.state.Epoch(), un.GroupID)
				// the other party has not reached a group consensus, and directly returns the group aggregation signature (if it exists locally)
				var rgqc *protocols.RGBlockQuorumCert
				if groupLen-un.UnKnownSet.HasLength() < groupThreshold {
					rgqc = cbft.state.FindMaxGroupRGBlockQuorumCert(msg.BlockIndex, un.GroupID)
					if rgqc != nil {
						RGBlockQuorumCerts = append(RGBlockQuorumCerts, rgqc)
						matched := rgqc.BlockQC.ValidatorSet.And(un.UnKnownSet).HasLength()
						remain -= matched
					}
					// Limit response votes
					if remain <= 0 {
						break
					}
				}
				if len(prepareVoteMap) > 0 {
					for i := uint32(0); i < un.UnKnownSet.Size(); i++ {
						if !un.UnKnownSet.GetIndex(i) || rgqc != nil && rgqc.BlockQC.HasSign(i) {
							continue
						}
						if v, ok := prepareVoteMap[i]; ok {
							votes = append(votes, v)
							remain--
						}
						// Limit response votes
						if remain <= 0 {
							break
						}
					}
				}
			}

			if len(votes) > 0 || len(RGBlockQuorumCerts) > 0 {
				prepareVotesV2 := &protocols.PrepareVotesV2{Epoch: msg.Epoch, ViewNumber: msg.ViewNumber, BlockIndex: msg.BlockIndex, Votes: votes, RGBlockQuorumCerts: RGBlockQuorumCerts}
				cbft.network.Send(id, prepareVotesV2)
				cbft.log.Debug("Send PrepareVotesV2", "peer", id, "blockIndex", msg.BlockIndex, "votes length", len(votes), "RGBlockQuorumCerts length", len(RGBlockQuorumCerts))
				return prepareVotesV2, nil
			}
		}
	}
	return nil, nil
}

// OnGetPrepareVoteV2 is responsible for processing the business logic
// of the GetPrepareVote message. It will synchronously return a
// PrepareVotes message to the sender.
func (cbft *Cbft) OnGetPrepareVote(id string, msg *protocols.GetPrepareVote) error {
	cbft.log.Debug("Received message on OnGetPrepareVote", "from", id, "msgHash", msg.MsgHash(), "message", msg.String())
	if msg.Epoch == cbft.state.Epoch() && msg.ViewNumber == cbft.state.ViewNumber() {
		// If the block has already QC, that response QC instead of votes.
		// Avoid the sender spent a lot of time to verifies PrepareVote msg.
		_, qc := cbft.state.ViewBlockAndQC(msg.BlockIndex)
		if qc != nil {
			cbft.network.Send(id, &protocols.BlockQuorumCert{BlockQC: qc})
			cbft.log.Debug("Send BlockQuorumCert", "peer", id, "qc", qc.String())
			return nil
		}

		prepareVoteMap := cbft.state.AllPrepareVoteByIndex(msg.BlockIndex)
		// Defining an array for receiving PrepareVote.
		votes := make([]*protocols.PrepareVote, 0, len(prepareVoteMap))
		if prepareVoteMap != nil {
			threshold := cbft.threshold(cbft.currentValidatorLen())
			remain := threshold - (cbft.currentValidatorLen() - int(msg.UnKnownSet.Size()))
			for k, v := range prepareVoteMap {
				if msg.UnKnownSet.GetIndex(k) {
					votes = append(votes, v)
				}

				// Limit response votes
				if remain > 0 && len(votes) >= remain {
					break
				}
			}
		}
		if len(votes) > 0 {
			cbft.network.Send(id, &protocols.PrepareVotes{Epoch: msg.Epoch, ViewNumber: msg.ViewNumber, BlockIndex: msg.BlockIndex, Votes: votes})
			cbft.log.Debug("Send PrepareVotes", "peer", id, "epoch", msg.Epoch, "viewNumber", msg.ViewNumber, "blockIndex", msg.BlockIndex)
		}
	}
	return nil
}

// OnPrepareVotes handling response from GetPrepareVote response.
func (cbft *Cbft) OnPrepareVotes(id string, msg *protocols.PrepareVotes) error {
	cbft.log.Debug("Received message on OnPrepareVotes", "from", id, "msgHash", msg.MsgHash(), "message", msg.String())
	for _, vote := range msg.Votes {
		_, qc := cbft.blockTree.FindBlockAndQC(vote.BlockHash, vote.BlockNumber)
		if qc == nil && !cbft.network.ContainsHistoryMessageHash(vote.MsgHash()) {
			if err := cbft.OnPrepareVote(id, vote); err != nil {
				if e, ok := err.(HandleError); ok && e.AuthFailed() {
					cbft.log.Error("OnPrepareVotes failed", "peer", id, "err", err)
				}
				return err
			}
		}
	}
	return nil
}

// OnPrepareVotes handling response from GetPrepareVote response.
func (cbft *Cbft) OnPrepareVotesV2(id string, msg *protocols.PrepareVotesV2) error {
	cbft.log.Debug("Received message on OnPrepareVotesV2", "from", id, "msgHash", msg.MsgHash(), "message", msg.String())
	alreadyQC := func(hash common.Hash, number uint64) bool {
		if _, qc := cbft.blockTree.FindBlockAndQC(hash, number); qc != nil {
			return true
		}
		return false
	}

	for _, rgqc := range msg.RGBlockQuorumCerts {
		if alreadyQC(rgqc.BHash(), rgqc.BlockNum()) {
			return nil
		}
		if !cbft.network.ContainsHistoryMessageHash(rgqc.MsgHash()) {
			if err := cbft.OnRGBlockQuorumCert(id, rgqc); err != nil {
				if e, ok := err.(HandleError); ok && e.AuthFailed() {
					cbft.log.Error("OnRGBlockQuorumCert failed", "peer", id, "err", err)
				}
				return err
			}
		}
	}

	for _, vote := range msg.Votes {
		if alreadyQC(vote.BlockHash, vote.BlockNumber) {
			return nil
		}
		if !cbft.network.ContainsHistoryMessageHash(vote.MsgHash()) {
			if err := cbft.OnPrepareVote(id, vote); err != nil {
				if e, ok := err.(HandleError); ok && e.AuthFailed() {
					cbft.log.Error("OnPrepareVotes failed", "peer", id, "err", err)
				}
				return err
			}
		}
	}
	return nil
}

// OnGetLatestStatus hands GetLatestStatus messages.
//
// main logic:
// 1.Compare the blockNumber of the sending node with the local node,
// and if the blockNumber of local node is larger then reply LatestStatus message,
// the message contains the status information of the local node.
func (cbft *Cbft) OnGetLatestStatus(id string, msg *protocols.GetLatestStatus) error {
	cbft.log.Debug("Received message on OnGetLatestStatus", "from", id, "logicType", msg.LogicType, "msgHash", msg.MsgHash(), "message", msg.String())
	if msg.BlockNumber != 0 && msg.QuorumCert == nil || msg.LBlockNumber != 0 && msg.LQuorumCert == nil {
		cbft.log.Error("Invalid getLatestStatus,lack corresponding quorumCert", "getLatestStatus", msg.String())
		return nil
	}
	// Define a function that performs the send action.
	launcher := func(bType uint64, targetId string, blockNumber uint64, blockHash common.Hash, qc *ctypes.QuorumCert) error {
		err := cbft.network.PeerSetting(targetId, bType, blockNumber)
		if err != nil {
			cbft.log.Error("GetPeer failed", "err", err, "peerId", targetId)
			return err
		}
		// Synchronize block data with fetchBlock.
		cbft.fetchBlock(targetId, blockHash, blockNumber, qc)
		return nil
	}
	//
	if msg.LogicType == network.TypeForQCBn {
		localQCNum, localQCHash := cbft.state.HighestQCBlock().NumberU64(), cbft.state.HighestQCBlock().Hash()
		localLockNum, localLockHash := cbft.state.HighestLockBlock().NumberU64(), cbft.state.HighestLockBlock().Hash()
		if localQCNum == msg.BlockNumber && localQCHash == msg.BlockHash {
			cbft.log.Debug("Local qcBn is equal the sender's qcBn", "remoteBn", msg.BlockNumber, "localBn", localQCNum, "remoteHash", msg.BlockHash, "localHash", localQCHash)
			if forkedHash, forkedNum, forked := cbft.blockTree.IsForked(localQCHash, localQCNum); forked {
				cbft.log.Debug("Local highest QC forked", "forkedQCHash", forkedHash, "forkedQCNumber", forkedNum, "localQCHash", localQCHash, "localQCNumber", localQCNum)
				_, qc := cbft.blockTree.FindBlockAndQC(forkedHash, forkedNum)
				_, lockQC := cbft.blockTree.FindBlockAndQC(localLockHash, localLockNum)
				cbft.network.Send(id, &protocols.LatestStatus{
					BlockNumber:  forkedNum,
					BlockHash:    forkedHash,
					QuorumCert:   qc,
					LBlockNumber: localLockNum,
					LBlockHash:   localLockHash,
					LQuorumCert:  lockQC,
					LogicType:    msg.LogicType})
			}
			return nil
		}
		if localQCNum < msg.BlockNumber || (localQCNum == msg.BlockNumber && localQCHash != msg.BlockHash) {
			cbft.log.Debug("Local qcBn is less than the sender's qcBn", "remoteBn", msg.BlockNumber, "localBn", localQCNum)
			if msg.LBlockNumber == localQCNum && msg.LBlockHash != localQCHash {
				return launcher(msg.LogicType, id, msg.LBlockNumber, msg.LBlockHash, msg.LQuorumCert)
			}
			return launcher(msg.LogicType, id, msg.BlockNumber, msg.BlockHash, msg.QuorumCert)
		}
		// must carry block qc
		_, qc := cbft.blockTree.FindBlockAndQC(localQCHash, localQCNum)
		_, lockQC := cbft.blockTree.FindBlockAndQC(localLockHash, localLockNum)
		cbft.log.Debug("Local qcBn is larger than the sender's qcBn", "remoteBn", msg.BlockNumber, "localBn", localQCNum)
		cbft.network.Send(id, &protocols.LatestStatus{
			BlockNumber:  localQCNum,
			BlockHash:    localQCHash,
			QuorumCert:   qc,
			LBlockNumber: localLockNum,
			LBlockHash:   localLockHash,
			LQuorumCert:  lockQC,
			LogicType:    msg.LogicType,
		})
	}
	return nil
}

// OnLatestStatus is used to process LatestStatus messages that received from peer.
func (cbft *Cbft) OnLatestStatus(id string, msg *protocols.LatestStatus) error {
	cbft.log.Debug("Received message on OnLatestStatus", "from", id, "msgHash", msg.MsgHash(), "message", msg.String())
	if msg.BlockNumber != 0 && msg.QuorumCert == nil || msg.LBlockNumber != 0 && msg.LQuorumCert == nil {
		cbft.log.Error("Invalid LatestStatus,lack corresponding quorumCert", "latestStatus", msg.String())
		return nil
	}
	switch msg.LogicType {
	case network.TypeForQCBn:
		localQCBn, localQCHash := cbft.state.HighestQCBlock().NumberU64(), cbft.state.HighestQCBlock().Hash()
		if localQCBn < msg.BlockNumber || (localQCBn == msg.BlockNumber && localQCHash != msg.BlockHash) {
			err := cbft.network.PeerSetting(id, msg.LogicType, msg.BlockNumber)
			if err != nil {
				cbft.log.Error("PeerSetting failed", "err", err)
				return err
			}
			cbft.log.Debug("LocalQCBn is lower than sender's", "localBn", localQCBn, "remoteBn", msg.BlockNumber)
			if localQCBn == msg.LBlockNumber && localQCHash != msg.LBlockHash {
				cbft.log.Debug("OnLatestStatus ~ fetchBlock by LBlockHash and LBlockNumber")
				cbft.fetchBlock(id, msg.LBlockHash, msg.LBlockNumber, msg.LQuorumCert)
			} else {
				cbft.log.Debug("OnLatestStatus ~ fetchBlock by QCBlockHash and QCBlockNumber")
				cbft.fetchBlock(id, msg.BlockHash, msg.BlockNumber, msg.QuorumCert)
			}
		}
	}
	return nil
}

// OnPrepareBlockHash responsible for handling PrepareBlockHash message.
//
// Note: After receiving the PrepareBlockHash message, it is determined whether the
// block information exists locally. If not, send a network request to get
// the block data.
func (cbft *Cbft) OnPrepareBlockHash(id string, msg *protocols.PrepareBlockHash) error {
	cbft.log.Debug("Received message OnPrepareBlockHash", "from", id, "msgHash", msg.MsgHash(), "message", msg.String())
	if msg.Epoch == cbft.state.Epoch() && msg.ViewNumber == cbft.state.ViewNumber() {
		block := cbft.state.ViewBlockByIndex(msg.BlockIndex)
		if block == nil {
			cbft.network.RemoveMessageHash(id, msg.MsgHash())
			cbft.SyncPrepareBlock(id, msg.Epoch, msg.ViewNumber, msg.BlockIndex)
		}
	}
	return nil
}

// OnGetViewChangeV2 responds to nodes that require viewChange.
// The Epoch and viewNumber of viewChange must be consistent
// with the state of the current node.
func (cbft *Cbft) OnGetViewChangeV2(id string, msg *protocols.GetViewChangeV2) (ctypes.Message, error) {
	cbft.log.Debug("Received message on OnGetViewChangeV2", "from", id, "msgHash", msg.MsgHash(), "message", msg.String(), "local", cbft.state.ViewString())
	responseVcsCounter.Inc(1)

	localEpoch, localViewNumber := cbft.state.Epoch(), cbft.state.ViewNumber()

	isLocalView := func() bool {
		return msg.Epoch == localEpoch && msg.ViewNumber == localViewNumber
	}

	isLastView := func() bool {
		return (msg.Epoch == localEpoch && msg.ViewNumber+1 == localViewNumber) || (msg.Epoch+1 == localEpoch && localViewNumber == state.DefaultViewNumber)
	}

	isPreviousView := func() bool {
		return msg.Epoch == localEpoch && msg.ViewNumber+1 < localViewNumber
	}

	if isLocalView() && len(msg.UnKnownGroups.UnKnown) > 0 {
		validatorLen := cbft.currentValidatorLen()
		threshold := cbft.threshold(validatorLen)
		remain := threshold - (validatorLen - msg.UnKnownGroups.UnKnownSize())

		viewChanges := make([]*protocols.ViewChange, 0)
		RGViewChangeQuorumCerts := make([]*protocols.RGViewChangeQuorumCert, 0)

		viewChangeMap := cbft.state.AllViewChange()

		for _, un := range msg.UnKnownGroups.UnKnown {
			if un.UnKnownSet.Size() != uint32(validatorLen) {
				cbft.log.Error("Invalid request params,UnKnownGroups is not a specified length", "peer", id, "groupID", un.GroupID, "unKnownSet", un.UnKnownSet.String(), "validatorLen", validatorLen)
				break // Do not continue processing the request
			}
			// Limit response votes
			if remain <= 0 {
				break
			}
			groupLen := cbft.groupLen(cbft.state.Epoch(), un.GroupID)
			groupThreshold := cbft.groupThreshold(cbft.state.Epoch(), un.GroupID)
			// the other party has not reached a group consensus, and directly returns the group aggregation signature (if it exists locally)
			var rgqc *protocols.RGViewChangeQuorumCert
			if groupLen-un.UnKnownSet.HasLength() < groupThreshold {
				rgqc = cbft.state.FindMaxRGViewChangeQuorumCert(un.GroupID)
				if rgqc != nil {
					RGViewChangeQuorumCerts = append(RGViewChangeQuorumCerts, rgqc)
					matched := rgqc.ViewChangeQC.ValidatorSet().And(un.UnKnownSet).HasLength()
					remain -= matched
				}
				// Limit response votes
				if remain <= 0 {
					break
				}
			}
			if len(viewChangeMap) > 0 {
				for i := uint32(0); i < un.UnKnownSet.Size(); i++ {
					if !un.UnKnownSet.GetIndex(i) || rgqc != nil && rgqc.ViewChangeQC.HasSign(i) {
						continue
					}
					if v, ok := viewChangeMap[i]; ok {
						viewChanges = append(viewChanges, v)
						remain--
					}
					// Limit response votes
					if remain <= 0 {
						break
					}
				}
			}
		}

		if len(viewChanges) > 0 || len(RGViewChangeQuorumCerts) > 0 {
			viewChangesV2 := &protocols.ViewChangesV2{VCs: viewChanges, RGViewChangeQuorumCerts: RGViewChangeQuorumCerts}
			cbft.network.Send(id, viewChangesV2)
			cbft.log.Debug("Send ViewChangesV2", "peer", id, "viewChanges length", len(viewChanges), "RGViewChangeQuorumCerts length", len(RGViewChangeQuorumCerts))
			return viewChangesV2, nil
		}
		return nil, nil
	}
	// Return view QC in the case of less than 1.
	if isLastView() {
		lastViewChangeQC := cbft.state.LastViewChangeQC()
		if lastViewChangeQC == nil {
			cbft.log.Info("Not found lastViewChangeQC")
			return nil, nil
		}
		err := lastViewChangeQC.EqualAll(msg.Epoch, msg.ViewNumber)
		if err != nil {
			cbft.log.Error("Last view change is not equal msg.viewNumber", "err", err)
			return nil, err
		}
		cbft.network.Send(id, &protocols.ViewChangeQuorumCert{
			ViewChangeQC: lastViewChangeQC,
		})
		return nil, nil
	}
	// get previous viewChangeQC from wal db
	if isPreviousView() {
		if qc, err := cbft.bridge.GetViewChangeQC(msg.Epoch, msg.ViewNumber); err == nil && qc != nil {
			// also inform the local highest view
			highestqc, _ := cbft.bridge.GetViewChangeQC(localEpoch, localViewNumber-1)
			viewChangeQuorumCert := &protocols.ViewChangeQuorumCert{
				ViewChangeQC: qc,
			}
			if highestqc != nil {
				viewChangeQuorumCert.HighestViewChangeQC = highestqc
			}
			cbft.log.Debug("Send previous viewChange quorumCert", "viewChangeQuorumCert", viewChangeQuorumCert.String())
			cbft.network.Send(id, viewChangeQuorumCert)
			return nil, nil
		}
	}

	return nil, fmt.Errorf("request is not match local view, local:%s,msg:%s", cbft.state.ViewString(), msg.String())
}

// OnGetViewChangeV2 responds to nodes that require viewChange.
// The Epoch and viewNumber of viewChange must be consistent
// with the state of the current node.
func (cbft *Cbft) OnGetViewChange(id string, msg *protocols.GetViewChange) error {
	cbft.log.Debug("Received message on OnGetViewChange", "from", id, "msgHash", msg.MsgHash(), "message", msg.String(), "local", cbft.state.ViewString())

	localEpoch, localViewNumber := cbft.state.Epoch(), cbft.state.ViewNumber()

	isLocalView := func() bool {
		return msg.Epoch == localEpoch && msg.ViewNumber == localViewNumber
	}

	isLastView := func() bool {
		return (msg.Epoch == localEpoch && msg.ViewNumber+1 == localViewNumber) || (msg.Epoch+1 == localEpoch && localViewNumber == state.DefaultViewNumber)
	}

	isPreviousView := func() bool {
		return msg.Epoch == localEpoch && msg.ViewNumber+1 < localViewNumber
	}

	if isLocalView() {
		viewChanges := cbft.state.AllViewChange()

		vcs := &protocols.ViewChanges{}
		for k, v := range viewChanges {
			if msg.ViewChangeBits.GetIndex(k) {
				vcs.VCs = append(vcs.VCs, v)
			}
		}
		cbft.log.Debug("Send ViewChanges", "peer", id, "len", len(vcs.VCs))
		if len(vcs.VCs) != 0 {
			cbft.network.Send(id, vcs)
		}
		return nil
	}
	// Return view QC in the case of less than 1.
	if isLastView() {
		lastViewChangeQC := cbft.state.LastViewChangeQC()
		if lastViewChangeQC == nil {
			cbft.log.Info("Not found lastViewChangeQC")
			return nil
		}
		err := lastViewChangeQC.EqualAll(msg.Epoch, msg.ViewNumber)
		if err != nil {
			cbft.log.Error("Last view change is not equal msg.viewNumber", "err", err)
			return err
		}
		cbft.network.Send(id, &protocols.ViewChangeQuorumCert{
			ViewChangeQC: lastViewChangeQC,
		})
		return nil
	}
	// get previous viewChangeQC from wal db
	if isPreviousView() {
		if qc, err := cbft.bridge.GetViewChangeQC(msg.Epoch, msg.ViewNumber); err == nil && qc != nil {
			// also inform the local highest view
			highestqc, _ := cbft.bridge.GetViewChangeQC(localEpoch, localViewNumber-1)
			viewChangeQuorumCert := &protocols.ViewChangeQuorumCert{
				ViewChangeQC: qc,
			}
			if highestqc != nil {
				viewChangeQuorumCert.HighestViewChangeQC = highestqc
			}
			cbft.log.Debug("Send previous viewChange quorumCert", "viewChangeQuorumCert", viewChangeQuorumCert.String())
			cbft.network.Send(id, viewChangeQuorumCert)
			return nil
		}
	}

	return fmt.Errorf("request is not match local view, local:%s,msg:%s", cbft.state.ViewString(), msg.String())
}

// OnViewChangeQuorumCert handles the message type of ViewChangeQuorumCertMsg.
func (cbft *Cbft) OnViewChangeQuorumCert(id string, msg *protocols.ViewChangeQuorumCert) error {
	cbft.log.Debug("Received message on OnViewChangeQuorumCert", "from", id, "msgHash", msg.MsgHash(), "message", msg.String())
	viewChangeQC := msg.ViewChangeQC
	epoch, viewNumber, _, _, _, _ := viewChangeQC.MaxBlock()
	if cbft.state.Epoch() == epoch && cbft.state.ViewNumber() == viewNumber {
		if err := cbft.verifyViewChangeQC(msg.ViewChangeQC); err == nil {
			cbft.log.Info("Receive viewChangeQuorumCert, try change view by viewChangeQC", "view", cbft.state.ViewString(), "viewChangeQC", viewChangeQC.String())
			cbft.tryChangeViewByViewChange(msg.ViewChangeQC)
		} else {
			cbft.log.Debug("Verify ViewChangeQC failed", "err", err)
			return &authFailedError{err}
		}
	}
	// if the other party's view is still higher than the local one, continue to synchronize the view
	cbft.trySyncViewChangeQuorumCert(id, msg)
	return nil
}

// Synchronize view one by one according to the highest view notified by the other party
func (cbft *Cbft) trySyncViewChangeQuorumCert(id string, msg *protocols.ViewChangeQuorumCert) {
	highestViewChangeQC := msg.HighestViewChangeQC
	if highestViewChangeQC == nil {
		return
	}
	epoch, viewNumber, _, _, _, _ := highestViewChangeQC.MaxBlock()
	if cbft.state.Epoch() != epoch {
		return
	}
	if cbft.state.ViewNumber() == viewNumber {
		if err := cbft.verifyViewChangeQC(highestViewChangeQC); err == nil {
			cbft.log.Debug("The highest view is equal to local, change view by highestViewChangeQC directly", "localView", cbft.state.ViewString(), "futureView", highestViewChangeQC.String())
			cbft.tryChangeViewByViewChange(highestViewChangeQC)
		}
		return
	}
	if cbft.state.ViewNumber() < viewNumber {
		// if local view lags, synchronize view one by one
		if err := cbft.verifyViewChangeQC(highestViewChangeQC); err == nil {
			cbft.log.Debug("Receive future viewChange quorumCert, sync viewChangeQC with fast mode", "localView", cbft.state.ViewString(), "futureView", highestViewChangeQC.String())
			// request viewChangeQC for the current view
			cbft.network.Send(id, &protocols.GetViewChangeV2{
				Epoch:         cbft.state.Epoch(),
				ViewNumber:    cbft.state.ViewNumber(),
				UnKnownGroups: &ctypes.UnKnownGroups{UnKnown: make([]*ctypes.UnKnownGroup, 0)},
			})
		}
	}
}

// OnViewChanges handles the message type of ViewChangesMsg.
func (cbft *Cbft) OnViewChanges(id string, msg *protocols.ViewChanges) error {
	cbft.log.Debug("Received message on OnViewChanges", "from", id, "msgHash", msg.MsgHash(), "message", msg.String())
	for _, v := range msg.VCs {
		if !cbft.network.ContainsHistoryMessageHash(v.MsgHash()) {
			if err := cbft.OnViewChange(id, v); err != nil {
				if e, ok := err.(HandleError); ok && e.AuthFailed() {
					cbft.log.Error("OnViewChanges failed", "peer", id, "err", err)
				}
				return err
			}
		}
	}
	return nil
}

// OnViewChanges handles the message type of ViewChangesMsg.
func (cbft *Cbft) OnViewChangesV2(id string, msg *protocols.ViewChangesV2) error {
	cbft.log.Debug("Received message on OnViewChangesV2", "from", id, "msgHash", msg.MsgHash(), "message", msg.String())

	for _, rgqc := range msg.RGViewChangeQuorumCerts {
		if !cbft.network.ContainsHistoryMessageHash(rgqc.MsgHash()) {
			if err := cbft.OnRGViewChangeQuorumCert(id, rgqc); err != nil {
				if e, ok := err.(HandleError); ok && e.AuthFailed() {
					cbft.log.Error("OnRGViewChangeQuorumCert failed", "peer", id, "err", err)
				}
				return err
			}
		}
	}

	for _, v := range msg.VCs {
		if !cbft.network.ContainsHistoryMessageHash(v.MsgHash()) {
			if err := cbft.OnViewChange(id, v); err != nil {
				if e, ok := err.(HandleError); ok && e.AuthFailed() {
					cbft.log.Error("OnViewChanges failed", "peer", id, "err", err)
				}
				return err
			}
		}
	}
	return nil
}

func (cbft *Cbft) KnownVoteIndexes(blockIndex uint32) []uint32 {
	groupNodes := cbft.getGroupIndexes(cbft.state.Epoch())
	allVotes := cbft.state.AllPrepareVoteByIndex(blockIndex)
	cbft.log.Trace("KnownVoteIndexes", "blockIndex", blockIndex, "groupNodes", len(groupNodes), "allVotes", len(allVotes))

	known := make([]uint32, 0)
	for groupID, indexes := range groupNodes {
		qc, _ := cbft.state.FindMaxGroupRGQuorumCert(blockIndex, groupID)
		for _, index := range indexes {
			if _, ok := allVotes[index]; ok {
				known = append(known, index)
			} else if qc.HasSign(index) {
				known = append(known, index)
			}
		}
	}
	return known
}

func (cbft *Cbft) KnownViewChangeIndexes() []uint32 {
	groupNodes := cbft.getGroupIndexes(cbft.state.Epoch())
	allViewChanges := cbft.state.AllViewChange()
	cbft.log.Trace("KnownViewChangeIndexes", "groupNodes", len(groupNodes), "allViewChanges", len(allViewChanges))

	known := make([]uint32, 0)
	for groupID, indexes := range groupNodes {
		qc, _ := cbft.state.FindMaxGroupRGViewChangeQuorumCert(groupID)
		for _, index := range indexes {
			if _, ok := allViewChanges[index]; ok {
				known = append(known, index)
			} else if qc.HasSign(index) {
				known = append(known, index)
			}
		}
	}
	return known
}

func (cbft *Cbft) MissGroupVotes(blockIndex uint32) *ctypes.UnKnownGroups {
	groupNodes := cbft.getGroupIndexes(cbft.state.Epoch())
	allVotes := cbft.state.AllPrepareVoteByIndex(blockIndex)

	validatorLen := cbft.currentValidatorLen()
	cbft.log.Trace("MissGroupVotes", "groupNodes", len(groupNodes), "allVotes", len(allVotes), "validatorLen", validatorLen)

	unKnowns := &ctypes.UnKnownGroups{UnKnown: make([]*ctypes.UnKnownGroup, 0)}

	// just for record metrics
	missGroups := 0
	missVotes := 0

	for groupID, indexes := range groupNodes {
		qc, _ := cbft.state.FindMaxGroupRGQuorumCert(blockIndex, groupID)
		groupLen := cbft.groupLen(cbft.state.Epoch(), groupID)
		// just for record metrics
		groupThreshold := cbft.groupThreshold(cbft.state.Epoch(), groupID)
		if qc.Len() < groupThreshold {
			missGroups++
		}
		if qc.Len() < groupLen {
			unKnownSet := utils.NewBitArray(uint32(validatorLen))
			for _, index := range indexes {
				if _, ok := allVotes[index]; !ok && !qc.HasSign(index) {
					if vote := cbft.csPool.GetPrepareVote(cbft.state.Epoch(), cbft.state.ViewNumber(), blockIndex, index); vote != nil {
						go cbft.ReceiveMessage(ctypes.NewInnerMsgInfo(vote.Msg, vote.PeerID))
						continue
					}
					unKnownSet.SetIndex(index, true)
				}
			}
			if unKnownSet.HasLength() > 0 {
				unKnownGroup := &ctypes.UnKnownGroup{
					GroupID:    groupID,
					UnKnownSet: unKnownSet,
				}
				unKnowns.UnKnown = append(unKnowns.UnKnown, unKnownGroup)
				// just for record metrics
				missVotes += unKnownSet.HasLength()
			}
		}
	}
	// just for record metrics
	missRGBlockQuorumCertsGauge.Update(int64(missGroups))
	missVotesGauge.Update(int64(missVotes))
	return unKnowns
}

func (cbft *Cbft) MissGroupViewChanges() *ctypes.UnKnownGroups {
	groupNodes := cbft.getGroupIndexes(cbft.state.Epoch())
	allViewChanges := cbft.state.AllViewChange()

	validatorLen := cbft.currentValidatorLen()
	cbft.log.Trace("MissGroupViewChanges", "groupNodes", len(groupNodes), "allViewChanges", len(allViewChanges), "validatorLen", validatorLen)

	unKnowns := &ctypes.UnKnownGroups{UnKnown: make([]*ctypes.UnKnownGroup, 0)}

	// just for record metrics
	missGroups := 0
	missViewChanges := 0

	for groupID, indexes := range groupNodes {
		qc, _ := cbft.state.FindMaxGroupRGViewChangeQuorumCert(groupID)
		groupLen := cbft.groupLen(cbft.state.Epoch(), groupID)
		// just for record metrics
		groupThreshold := cbft.groupThreshold(cbft.state.Epoch(), groupID)
		if qc.Len() < groupThreshold {
			missGroups++
		}
		if qc.Len() < groupLen {
			unKnownSet := utils.NewBitArray(uint32(validatorLen))
			for _, index := range indexes {
				if _, ok := allViewChanges[index]; !ok && !qc.HasSign(index) {
					unKnownSet.SetIndex(index, true)
				}
			}
			if unKnownSet.HasLength() > 0 {
				unKnownGroup := &ctypes.UnKnownGroup{
					GroupID:    groupID,
					UnKnownSet: unKnownSet,
				}
				unKnowns.UnKnown = append(unKnowns.UnKnown, unKnownGroup)
				// just for record metrics
				missViewChanges += unKnownSet.HasLength()
			}
		}
	}
	// just for record metrics
	missRGViewQuorumCertsGauge.Update(int64(missGroups))
	missVcsGauge.Update(int64(missViewChanges))
	return unKnowns
}

// MissingPrepareVote returns missing vote.
func (cbft *Cbft) MissingPrepareVote() (v ctypes.Message, err error) {
	if !cbft.NeedGroup() {
		return cbft.MissingPrepareVoteV1()
	}
	result := make(chan struct{})

	cbft.asyncCallCh <- func() {
		defer func() { result <- struct{}{} }()

		begin := cbft.state.MaxQCIndex() + 1
		end := cbft.state.NextViewBlockIndex()
		threshold := cbft.threshold(cbft.currentValidatorLen())
		cbft.log.Trace("Synchronize votes by grouped channel", "threshold", threshold)

		block := cbft.state.HighestQCBlock()
		blockTime := common.MillisToTime(int64(block.Time()))

		for index := begin; index < end; index++ {
			if time.Since(blockTime) < syncPrepareVotesInterval {
				err = fmt.Errorf("not need sync prepare vote")
				break
			}

			size := len(cbft.KnownVoteIndexes(index))
			// We need sync prepare votes when a long time not arrived QC.
			if size < threshold { // need sync prepare votes
				unKnownGroups := cbft.MissGroupVotes(index)
				cbft.log.Trace("Synchronize votes by grouped channel,missGroupVotes", "blockIndex", index, "threshold", threshold, "size", size, "unKnownGroups", unKnownGroups.String())
				if len(unKnownGroups.UnKnown) > 0 {
					v, err = &protocols.GetPrepareVoteV2{
						Epoch:         cbft.state.Epoch(),
						ViewNumber:    cbft.state.ViewNumber(),
						BlockIndex:    index,
						UnKnownGroups: unKnownGroups,
					}, nil
					cbft.log.Debug("Synchronize votes by grouped channel,missingPrepareVote", "blockIndex", index, "known", size, "threshold", threshold, "request", v.String())
					missVotesCounter.Inc(1)
					break
				}
			}
		}
		if v == nil {
			err = fmt.Errorf("not need sync prepare vote")
		}
	}
	<-result
	return
}

// MissingViewChangeNodes returns the node ID of the missing viewChanges.
// Use the channel to complete serial execution to prevent concurrency.
func (cbft *Cbft) MissingViewChangeNodes() (v ctypes.Message, err error) {
	if !cbft.NeedGroup() {
		return cbft.MissingViewChangeNodesV1()
	}

	result := make(chan struct{})

	cbft.asyncCallCh <- func() {
		defer func() { result <- struct{}{} }()

		if !cbft.state.IsDeadline() {
			v, err = nil, fmt.Errorf("no need sync viewchange")
			return
		}

		threshold := cbft.threshold(cbft.currentValidatorLen())
		size := len(cbft.KnownViewChangeIndexes())
		cbft.log.Trace("Synchronize viewChanges by grouped channel", "threshold", threshold, "size", size)

		if size < threshold {
			unKnownGroups := cbft.MissGroupViewChanges()
			cbft.log.Trace("Synchronize viewChanges by grouped channel,missGroupViewChanges", "threshold", threshold, "size", size, "unKnownGroups", unKnownGroups.String())
			if len(unKnownGroups.UnKnown) > 0 {
				v, err = &protocols.GetViewChangeV2{
					Epoch:         cbft.state.Epoch(),
					ViewNumber:    cbft.state.ViewNumber(),
					UnKnownGroups: unKnownGroups,
				}, nil
				cbft.log.Debug("Synchronize viewChanges by grouped channel,missingViewChangeNodes", "known", size, "threshold", threshold, "request", v.String())
				missVcsCounter.Inc(1)
			}
		}
		if v == nil {
			err = fmt.Errorf("not need sync prepare vote")
		}
	}
	<-result
	return
}

// MissingPrepareVote returns missing vote.
func (cbft *Cbft) MissingPrepareVoteV1() (v *protocols.GetPrepareVote, err error) {
	result := make(chan struct{})

	cbft.asyncCallCh <- func() {
		defer func() { result <- struct{}{} }()

		begin := cbft.state.MaxQCIndex() + 1
		end := cbft.state.NextViewBlockIndex()
		len := cbft.currentValidatorLen()
		cbft.log.Debug("MissingPrepareVote", "epoch", cbft.state.Epoch(), "viewNumber", cbft.state.ViewNumber(), "beginIndex", begin, "endIndex", end, "validatorLen", len)

		block := cbft.state.HighestQCBlock()
		blockTime := common.MillisToTime(int64(block.Time()))

		for index := begin; index < end; index++ {
			size := cbft.state.PrepareVoteLenByIndex(index)
			cbft.log.Debug("The length of prepare vote", "index", index, "size", size)

			// We need sync prepare votes when a long time not arrived QC.
			if size < cbft.threshold(len) && time.Since(blockTime) >= syncPrepareVotesInterval { // need sync prepare votes
				knownVotes := cbft.state.AllPrepareVoteByIndex(index)
				unKnownSet := utils.NewBitArray(uint32(len))
				for i := uint32(0); i < unKnownSet.Size(); i++ {
					if vote := cbft.csPool.GetPrepareVote(cbft.state.Epoch(), cbft.state.ViewNumber(), index, i); vote != nil {
						go cbft.ReceiveMessage(vote)
						continue
					}
					if _, ok := knownVotes[i]; !ok {
						unKnownSet.SetIndex(i, true)
					}
				}

				v, err = &protocols.GetPrepareVote{
					Epoch:      cbft.state.Epoch(),
					ViewNumber: cbft.state.ViewNumber(),
					BlockIndex: index,
					UnKnownSet: unKnownSet,
				}, nil
				break
			}
		}
		if v == nil {
			err = fmt.Errorf("not need sync prepare vote")
		}
	}
	<-result
	return
}

// MissingViewChangeNodes returns the node ID of the missing vote.
// Use the channel to complete serial execution to prevent concurrency.
func (cbft *Cbft) MissingViewChangeNodesV1() (v *protocols.GetViewChange, err error) {
	result := make(chan struct{})

	cbft.asyncCallCh <- func() {
		defer func() { result <- struct{}{} }()
		allViewChange := cbft.state.AllViewChange()

		length := cbft.currentValidatorLen()
		vbits := utils.NewBitArray(uint32(length))

		// enough qc or did not reach deadline
		if len(allViewChange) >= cbft.threshold(length) || !cbft.state.IsDeadline() {
			v, err = nil, fmt.Errorf("no need sync viewchange")
			return
		}
		for i := uint32(0); i < vbits.Size(); i++ {
			if _, ok := allViewChange[i]; !ok {
				vbits.SetIndex(i, true)
			}
		}

		v, err = &protocols.GetViewChange{
			Epoch:          cbft.state.Epoch(),
			ViewNumber:     cbft.state.ViewNumber(),
			ViewChangeBits: vbits,
		}, nil
	}
	<-result
	return
}

// LatestStatus returns latest status.
func (cbft *Cbft) LatestStatus() (v *protocols.GetLatestStatus) {
	result := make(chan struct{})

	cbft.asyncCallCh <- func() {
		defer func() { result <- struct{}{} }()

		qcBn, qcHash := cbft.HighestQCBlockBn()
		_, qc := cbft.blockTree.FindBlockAndQC(qcHash, qcBn)

		lockBn, lockHash := cbft.HighestLockBlockBn()
		_, lockQC := cbft.blockTree.FindBlockAndQC(lockHash, lockBn)

		v = &protocols.GetLatestStatus{
			BlockNumber:  qcBn,
			BlockHash:    qcHash,
			QuorumCert:   qc,
			LBlockNumber: lockBn,
			LBlockHash:   lockHash,
			LQuorumCert:  lockQC,
		}
	}
	<-result
	return
}

// OnPong is used to receive the average delay time.
func (cbft *Cbft) OnPong(nodeID string, netLatency int64) error {
	cbft.log.Trace("OnPong", "nodeID", nodeID, "netLatency", netLatency)
	cbft.netLatencyLock.Lock()
	defer cbft.netLatencyLock.Unlock()
	latencyList, exist := cbft.netLatencyMap[nodeID]
	if !exist {
		cbft.netLatencyMap[nodeID] = list.New()
		cbft.netLatencyMap[nodeID].PushBack(netLatency)
	} else {
		if latencyList.Len() > 5 {
			e := latencyList.Front()
			cbft.netLatencyMap[nodeID].Remove(e)
		}
		cbft.netLatencyMap[nodeID].PushBack(netLatency)
	}
	return nil
}

// BlockExists is used to query whether the specified block exists in this node.
func (cbft *Cbft) BlockExists(blockNumber uint64, blockHash common.Hash) error {
	result := make(chan error, 1)
	cbft.asyncCallCh <- func() {
		if (blockHash == common.Hash{}) {
			result <- fmt.Errorf("invalid blockHash")
			return
		}
		block := cbft.blockTree.FindBlockByHash(blockHash)
		if block = cbft.blockChain.GetBlock(blockHash, blockNumber); block == nil {
			result <- fmt.Errorf("not found block by hash:%s, number:%d", blockHash.TerminalString(), blockNumber)
			return
		}
		if block.Hash() != blockHash || blockNumber != block.NumberU64() {
			result <- fmt.Errorf("not match from block, hash:%s, number:%d, queriedHash:%s, queriedNumber:%d",
				blockHash.TerminalString(), blockNumber,
				block.Hash().TerminalString(), block.NumberU64())
			return
		}
		result <- nil
	}
	return <-result
}

// AvgLatency returns the average delay time of the specified node.
//
// The average is the average delay between the current
// node and all consensus nodes.
// Return value unit: milliseconds.
func (cbft *Cbft) AvgLatency() time.Duration {
	cbft.netLatencyLock.Lock()
	defer cbft.netLatencyLock.Unlock()
	// The intersection of peerSets and consensusNodes.
	target, err := cbft.network.AliveConsensusNodeIDs()
	if err != nil {
		return time.Duration(0)
	}
	var (
		avgSum     int64
		result     int64
		validCount int64
	)
	// Take 2/3 nodes from the target.
	var pair utils.KeyValuePairList
	for _, v := range target {
		if latencyList, exist := cbft.netLatencyMap[v]; exist {
			avg := calAverage(latencyList)
			pair.Push(utils.KeyValuePair{Key: v, Value: avg})
		}
	}
	sort.Sort(pair)
	if pair.Len() == 0 {
		return time.Duration(0)
	}
	validCount = int64(pair.Len() * 2 / 3)
	if validCount == 0 {
		validCount = 1
	}
	for _, v := range pair[:validCount] {
		avgSum += v.Value
	}

	result = avgSum / validCount
	cbft.log.Debug("Get avg latency", "avg", result)
	return time.Duration(result) * time.Millisecond
}

// DefaultAvgLatency returns the avg latency of default.
func (cbft *Cbft) DefaultAvgLatency() time.Duration {
	return time.Duration(protocols.DefaultAvgLatency) * time.Millisecond
}

func calAverage(latencyList *list.List) int64 {
	var (
		sum    int64
		counts int64
	)
	for e := latencyList.Front(); e != nil; e = e.Next() {
		if latency, ok := e.Value.(int64); ok {
			counts++
			sum += latency
		}
	}
	if counts > 0 {
		return sum / counts
	}
	return 0
}

func (cbft *Cbft) SyncPrepareBlock(id string, epoch uint64, viewNumber uint64, blockIndex uint32) {
	if msg := cbft.csPool.GetPrepareBlock(epoch, viewNumber, blockIndex); msg != nil {
		go cbft.ReceiveMessage(ctypes.NewInnerMsgInfo(msg.Msg, msg.PeerID))
	}
	if cbft.syncingCache.AddOrReplace(blockIndex) {
		msg := &protocols.GetPrepareBlock{Epoch: epoch, ViewNumber: viewNumber, BlockIndex: blockIndex}
		if id == "" {
			cbft.network.PartBroadcast(msg)
			cbft.log.Debug("Send GetPrepareBlock by part broadcast", "msg", msg.String())
		} else {
			cbft.network.Send(id, msg)
			cbft.log.Debug("Send GetPrepareBlock", "peer", id, "msg", msg.String())
		}
	}
}

func (cbft *Cbft) SyncBlockQuorumCert(id string, blockNumber uint64, blockHash common.Hash, blockIndex uint32) {
	if msg := cbft.csPool.GetPrepareQC(cbft.state.Epoch(), cbft.state.ViewNumber(), blockIndex); msg != nil {
		go cbft.ReceiveSyncMsg(ctypes.NewInnerMsgInfo(msg.Msg, msg.PeerID))
	}
	if cbft.syncingCache.AddOrReplace(blockHash) {
		msg := &protocols.GetBlockQuorumCert{BlockHash: blockHash, BlockNumber: blockNumber}
		cbft.network.Send(id, msg)
		cbft.log.Debug("Send GetBlockQuorumCert", "peer", id, "msg", msg.String())
	}
}
