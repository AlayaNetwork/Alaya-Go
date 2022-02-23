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
	"time"

	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/network"

	"github.com/pkg/errors"

	"github.com/AlayaNetwork/Alaya-Go/crypto/bls"

	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/utils"

	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/executor"

	"github.com/AlayaNetwork/Alaya-Go/log"

	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/common/math"
	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/protocols"
	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/state"
	ctypes "github.com/AlayaNetwork/Alaya-Go/consensus/cbft/types"
	"github.com/AlayaNetwork/Alaya-Go/core/cbfttypes"
	"github.com/AlayaNetwork/Alaya-Go/core/types"
)

// OnPrepareBlock performs security rule verification，store in blockTree,
// Whether to start synchronization
func (cbft *Cbft) OnPrepareBlock(id string, msg *protocols.PrepareBlock) error {
	cbft.log.Debug("Receive PrepareBlock", "id", id, "msg", msg.String())
	if err := cbft.VerifyHeader(nil, msg.Block.Header(), false); err != nil {
		cbft.log.Error("Verify header fail", "number", msg.Block.Number(), "hash", msg.Block.Hash(), "err", err)
		return err
	}
	if err := cbft.safetyRules.PrepareBlockRules(msg); err != nil {
		blockCheckFailureMeter.Mark(1)

		if err.Common() {
			cbft.log.Debug("Prepare block rules fail", "number", msg.Block.Number(), "hash", msg.Block.Hash(), "err", err)
			return err
		}
		// verify consensus signature
		if err := cbft.verifyConsensusSign(msg); err != nil {
			signatureCheckFailureMeter.Mark(1)
			return err
		}

		if err.Fetch() {
			if cbft.isProposer(msg.Epoch, msg.ViewNumber, msg.ProposalIndex) {
				cbft.log.Info("Epoch or viewNumber higher than local, try to fetch block", "fetchHash", msg.Block.ParentHash(), "fetchNumber", msg.Block.NumberU64()-1)
				cbft.fetchBlock(id, msg.Block.ParentHash(), msg.Block.NumberU64()-1, nil)
			}
			return err
		}
		if err.FetchPrepare() {
			if cbft.isProposer(msg.Epoch, msg.ViewNumber, msg.ProposalIndex) {
				cbft.prepareBlockFetchRules(id, msg)
			}
			return err
		}
		if err.NewView() {
			var block *types.Block
			var qc *ctypes.QuorumCert
			if e := cbft.checkViewChangeQC(msg); e != nil {
				return e
			}
			if msg.ViewChangeQC != nil {
				_, _, _, _, hash, number := msg.ViewChangeQC.MaxBlock()
				block, qc = cbft.blockTree.FindBlockAndQC(hash, number)
			} else {
				block, qc = cbft.blockTree.FindBlockAndQC(msg.Block.ParentHash(), msg.Block.NumberU64()-1)
			}
			cbft.log.Debug("Receive new view's block, change view", "newEpoch", msg.Epoch, "newView", msg.ViewNumber)
			cbft.changeView(msg.Epoch, msg.ViewNumber, block, qc, msg.ViewChangeQC)
		}
	}

	if _, err := cbft.verifyConsensusMsg(msg); err != nil {
		cbft.log.Error("Failed to verify prepareBlock", "prepare", msg.String(), "error", err.Error())
		signatureCheckFailureMeter.Mark(1)
		return err
	}
	// The new block is notified by the PrepareBlockHash to the nodes in the network.
	cbft.state.AddPrepareBlock(msg)
	cbft.log.Debug("Receive new prepareBlock", "msgHash", msg.MsgHash(), "prepare", msg.String())
	cbft.findExecutableBlock()
	return nil
}

// OnPrepareVote perform security rule verification，store in blockTree,
// Whether to start synchronization
func (cbft *Cbft) OnPrepareVote(id string, msg *protocols.PrepareVote) error {
	cbft.log.Debug("Receive PrepareVote", "id", id, "msg", msg.String())
	if err := cbft.safetyRules.PrepareVoteRules(msg); err != nil {
		if err.Common() {
			cbft.log.Debug("Preparevote rules fail", "number", msg.BlockNumber, "hash", msg.BlockHash, "err", err)
			return err
		}

		// verify consensus signature
		if cbft.verifyConsensusSign(msg) != nil {
			signatureCheckFailureMeter.Mark(1)
			return err
		}

		if err.Fetch() {
			if msg.ParentQC != nil {
				cbft.log.Info("Epoch or viewNumber higher than local, try to fetch block", "fetchHash", msg.ParentQC.BlockHash, "fetchNumber", msg.ParentQC.BlockNumber)
				cbft.fetchBlock(id, msg.ParentQC.BlockHash, msg.ParentQC.BlockNumber, msg.ParentQC)
			}
		} else if err.FetchPrepare() {
			cbft.prepareVoteFetchRules(id, msg)
		}
		return err
	}

	var node *cbfttypes.ValidateNode
	var err error
	if node, err = cbft.verifyConsensusMsg(msg); err != nil {
		cbft.log.Error("Failed to verify prepareVote", "vote", msg.String(), "error", err.Error())
		return err
	}

	cbft.state.AddPrepareVote(node.Index, msg)
	cbft.mergeVoteToQuorumCerts(node, msg)
	cbft.log.Debug("Receive new prepareVote", "msgHash", msg.MsgHash(), "vote", msg.String(), "votes", cbft.state.PrepareVoteLenByIndex(msg.BlockIndex))

	cbft.trySendRGBlockQuorumCert()
	cbft.insertPrepareQC(msg.ParentQC)
	cbft.findQCBlock()
	return nil
}

// OnViewChange performs security rule verification, view switching.
func (cbft *Cbft) OnViewChange(id string, msg *protocols.ViewChange) error {
	cbft.log.Debug("Receive ViewChange", "id", id, "msg", msg.String())
	if err := cbft.safetyRules.ViewChangeRules(msg); err != nil {
		if err.Fetch() {
			if msg.PrepareQC != nil {
				cbft.log.Info("Epoch or viewNumber higher than local, try to fetch block", "fetchHash", msg.BlockHash, "fetchNumber", msg.BlockNumber)
				cbft.fetchBlock(id, msg.BlockHash, msg.BlockNumber, msg.PrepareQC)
			}
		}
		return err
	}

	var node *cbfttypes.ValidateNode
	var err error

	if node, err = cbft.verifyConsensusMsg(msg); err != nil {
		cbft.log.Error("Failed to verify viewChange", "viewChange", msg.String(), "error", err.Error())
		return err
	}

	cbft.state.AddViewChange(node.Index, msg)
	cbft.mergeViewChangeToViewChangeQuorumCerts(node, msg)
	cbft.log.Debug("Receive new viewChange", "msgHash", msg.MsgHash(), "viewChange", msg.String(), "total", cbft.state.ViewChangeLen())
	cbft.trySendRGViewChangeQuorumCert()

	// It is possible to achieve viewchangeQC every time you add viewchange
	cbft.tryChangeView()
	return nil
}

// OnRGBlockQuorumCert perform security rule verification，store in viewRGBlockQuorumCerts and selectedRGBlockQuorumCerts,
// Whether to start synchronization
func (cbft *Cbft) OnRGBlockQuorumCert(id string, msg *protocols.RGBlockQuorumCert) error {
	cbft.log.Debug("Receive RGBlockQuorumCert", "id", id, "msg", msg.String())
	if err := cbft.safetyRules.RGBlockQuorumCertRules(msg); err != nil {
		if err.Common() {
			cbft.log.Debug("RGBlockQuorumCert rules fail", "number", msg.BlockNum(), "hash", msg.BHash(), "err", err)
			return err
		}

		// verify consensus signature
		if cbft.verifyConsensusSign(msg) != nil {
			signatureCheckFailureMeter.Mark(1)
			return err
		}

		if err.Fetch() {
			if msg.ParentQC != nil {
				cbft.log.Info("Epoch or viewNumber higher than local, try to fetch block", "fetchHash", msg.ParentQC.BlockHash, "fetchNumber", msg.ParentQC.BlockNumber)
				cbft.fetchBlock(id, msg.ParentQC.BlockHash, msg.ParentQC.BlockNumber, msg.ParentQC)
			}
		} else if err.FetchPrepare() {
			cbft.prepareVoteFetchRules(id, msg)
		}
		return err
	}

	if err := cbft.AllowRGQuorumCert(msg); err != nil {
		cbft.log.Error("Failed to allow RGBlockQuorumCert", "RGBlockQuorumCert", msg.String(), "error", err.Error())
		return err
	}

	var node *cbfttypes.ValidateNode
	var err error
	if node, err = cbft.verifyConsensusMsg(msg); err != nil {
		cbft.log.Error("Failed to verify RGBlockQuorumCert", "RGBlockQuorumCert", msg.String(), "error", err.Error())
		return err
	}

	// VerifyQuorumCert,This method simply verifies the correctness of the aggregated signature itself
	// Before this, it is necessary to verify parentqc, whether the number of group signatures is sufficient, whether all signers are group members, whether the message is sent by group members.
	if err := cbft.verifyQuorumCert(msg.BlockQC); err != nil {
		cbft.log.Error("Failed to verify RGBlockQuorumCert blockQC", "blockQC", msg.BlockQC.String(), "err", err.Error())
		return &authFailedError{err}
	}

	cbft.state.AddRGBlockQuorumCert(node.Index, msg)
	blockQC, ParentQC := msg.BlockQC.DeepCopyQuorumCert(), msg.ParentQC
	cbft.richBlockQuorumCert(msg.EpochNum(), msg.BlockIndx(), msg.GroupID, blockQC)
	cbft.state.AddSelectRGQuorumCerts(msg.BlockIndx(), msg.GroupID, blockQC, ParentQC)
	cbft.log.Debug("Receive new RGBlockQuorumCert", "msgHash", msg.MsgHash(), "RGBlockQuorumCert", msg.String(), "total", cbft.state.RGBlockQuorumCertsLen(msg.BlockIndx(), msg.GroupID))

	cbft.trySendRGBlockQuorumCert()
	cbft.insertPrepareQC(ParentQC)
	cbft.findQCBlock()
	return nil
}

// Determine whether the total number of RGBlockQuorumCert signatures has reached the minimum threshold for group consensus nodes
func (cbft *Cbft) enoughSigns(epoch uint64, groupID uint32, signs int) bool {
	threshold := cbft.groupThreshold(epoch, groupID)
	return signs >= threshold
}

// Determine whether the signer of the RGBlockQuorumCert message is a member of the group
func (cbft *Cbft) isGroupMember(epoch uint64, groupID uint32, nodeIndex uint32) bool {
	// Index collection of the group members
	indexes, err := cbft.validatorPool.GetValidatorIndexesByGroupID(epoch, groupID)
	if err != nil || indexes == nil {
		return false
	}
	for _, index := range indexes {
		if index == nodeIndex {
			return true
		}
	}
	return false
}

// Determine whether the aggregate signers in the RGBlockQuorumCert message are all members of the group
func (cbft *Cbft) allGroupMember(epoch uint64, groupID uint32, validatorSet *utils.BitArray) bool {
	// Index collection of the group members
	indexes, err := cbft.validatorPool.GetValidatorIndexesByGroupID(epoch, groupID)
	if err != nil || indexes == nil {
		return false
	}
	total := cbft.validatorPool.Len(epoch)
	vSet := utils.NewBitArray(uint32(total))
	for _, index := range indexes {
		vSet.SetIndex(index, true)
	}

	return vSet.Contains(validatorSet)
}

// Verify the aggregate signer information of RGQuorumCert
func (cbft *Cbft) AllowRGQuorumCert(msg ctypes.ConsensusMsg) error {
	epoch := msg.EpochNum()
	nodeIndex := msg.NodeIndex()
	var (
		groupID      uint32
		validatorSet *utils.BitArray
		signsTotal   int
	)

	switch rg := msg.(type) {
	case *protocols.RGBlockQuorumCert:
		groupID = rg.GroupID
		signsTotal = rg.BlockQC.Len()
		validatorSet = rg.BlockQC.ValidatorSet
	case *protocols.RGViewChangeQuorumCert:
		groupID = rg.GroupID
		//signsTotal = rg.ViewChangeQC.Len()
		signsTotal = rg.ViewChangeQC.HasLength()
		validatorSet = rg.ViewChangeQC.ValidatorSet()
	}

	if !cbft.enoughSigns(epoch, groupID, signsTotal) {
		return authFailedError{
			err: fmt.Errorf("insufficient signatures"),
		}
	}
	if !cbft.isGroupMember(epoch, groupID, nodeIndex) {
		return authFailedError{
			err: fmt.Errorf("the message sender is not a member of the group"),
		}
	}
	if !cbft.allGroupMember(epoch, groupID, validatorSet) {
		return authFailedError{
			err: fmt.Errorf("signers include non-group members"),
		}
	}
	return nil
}

// OnRGViewChangeQuorumCert performs security rule verification, view switching.
func (cbft *Cbft) OnRGViewChangeQuorumCert(id string, msg *protocols.RGViewChangeQuorumCert) error {
	cbft.log.Debug("Receive RGViewChangeQuorumCert", "id", id, "msg", msg.String())
	if err := cbft.safetyRules.RGViewChangeQuorumCertRules(msg); err != nil {
		if err.Fetch() {
			viewChangeQC := msg.ViewChangeQC
			_, _, _, _, blockHash, blockNumber := viewChangeQC.MaxBlock()
			if msg.PrepareQCs != nil && msg.PrepareQCs.FindPrepareQC(blockHash) != nil {
				cbft.log.Info("Epoch or viewNumber higher than local, try to fetch block", "fetchHash", blockHash, "fetchNumber", blockNumber)
				cbft.fetchBlock(id, blockHash, blockNumber, msg.PrepareQCs.FindPrepareQC(blockHash))
			}
		}
		return err
	}

	if err := cbft.AllowRGQuorumCert(msg); err != nil {
		cbft.log.Error("Failed to allow RGViewChangeQuorumCert", "RGViewChangeQuorumCert", msg.String(), "error", err.Error())
		return err
	}

	var node *cbfttypes.ValidateNode
	var err error

	if node, err = cbft.verifyConsensusMsg(msg); err != nil {
		cbft.log.Error("Failed to verify RGViewChangeQuorumCert", "viewChange", msg.String(), "error", err.Error())
		return err
	}

	// VerifyQuorumCert,This method simply verifies the correctness of the aggregated signature itself
	// Before this, it is necessary to verify parentqc, whether the number of group signatures is sufficient, whether all signers are group members, whether the message is sent by group members.
	if err := cbft.verifyGroupViewChangeQC(msg.GroupID, msg.ViewChangeQC); err != nil {
		cbft.log.Error("Failed to verify RGViewChangeQuorumCert viewChangeQC", "err", err.Error())
		return &authFailedError{err}
	}

	cbft.state.AddRGViewChangeQuorumCert(node.Index, msg)
	viewChangeQC, prepareQCs := msg.ViewChangeQC.DeepCopyViewChangeQC(), msg.PrepareQCs.DeepCopyPrepareQCs()
	cbft.richViewChangeQuorumCert(msg.EpochNum(), msg.GroupID, viewChangeQC, prepareQCs)
	cbft.state.AddSelectRGViewChangeQuorumCerts(msg.GroupID, viewChangeQC, prepareQCs.FlattenMap())
	cbft.log.Debug("Receive new RGViewChangeQuorumCert", "msgHash", msg.MsgHash(), "RGViewChangeQuorumCert", msg.String(), "total", cbft.state.RGViewChangeQuorumCertsLen(msg.GroupID))
	cbft.trySendRGViewChangeQuorumCert()

	// It is possible to achieve viewchangeQC every time you add viewchange
	cbft.tryChangeView()
	return nil
}

// OnViewTimeout performs timeout logic for view.
func (cbft *Cbft) OnViewTimeout() {
	cbft.log.Info("Current view timeout", "view", cbft.state.ViewString())
	node, err := cbft.isCurrentValidator()
	if err != nil {
		cbft.log.Info("ViewTimeout local node is not validator")
		return
	}

	if cbft.state.ViewChangeByIndex(node.Index) != nil {
		cbft.log.Debug("Had send viewchange, don't send again")
		return
	}

	hash, number := cbft.state.HighestQCBlock().Hash(), cbft.state.HighestQCBlock().NumberU64()
	_, qc := cbft.blockTree.FindBlockAndQC(hash, number)

	viewChange := &protocols.ViewChange{
		Epoch:          cbft.state.Epoch(),
		ViewNumber:     cbft.state.ViewNumber(),
		BlockHash:      hash,
		BlockNumber:    number,
		ValidatorIndex: node.Index,
		PrepareQC:      qc,
	}

	if err := cbft.signMsgByBls(viewChange); err != nil {
		cbft.log.Error("Sign ViewChange failed", "err", err)
		return
	}

	// write sendViewChange info to wal
	if !cbft.isLoading() {
		cbft.bridge.SendViewChange(viewChange)
	}

	cbft.state.AddViewChange(node.Index, viewChange)
	// send viewChange use pubsub
	if err := cbft.publishTopicMsg(viewChange); err != nil {
		cbft.log.Error("Publish viewChange failed", "err", err.Error(), "view", cbft.state.ViewString(), "viewChange", viewChange.String())
	}
	//cbft.network.Broadcast(viewChange)
	cbft.log.Info("Local add viewChange", "index", node.Index, "viewChange", viewChange.String(), "total", cbft.state.ViewChangeLen())

	cbft.tryChangeView()
}

// OnInsertQCBlock performs security rule verification, view switching.
func (cbft *Cbft) OnInsertQCBlock(blocks []*types.Block, qcs []*ctypes.QuorumCert) error {
	if len(blocks) != len(qcs) {
		return fmt.Errorf("block qc is inconsistent")
	}
	for i := 0; i < len(blocks); i++ {
		block, qc := blocks[i], qcs[i]

		if err := cbft.safetyRules.QCBlockRules(block, qc); err != nil {
			if err.NewView() {
				cbft.log.Info("The block to be written belongs to the next view, need change view", "view", cbft.state.ViewString(), "qcBlock", qc.String())
				cbft.changeView(qc.Epoch, qc.ViewNumber, block, qc, nil)
			} else {
				return err
			}
		}
		cbft.insertQCBlock(block, qc)
		cbft.log.Info("Insert QC block success", "qcBlock", qc.String())
	}

	cbft.findExecutableBlock()
	return nil
}

// Update blockTree, try commit new block
func (cbft *Cbft) insertQCBlock(block *types.Block, qc *ctypes.QuorumCert) {
	cbft.log.Info("Insert QC block", "qc", qc.String())
	if block.NumberU64() <= cbft.state.HighestLockBlock().NumberU64() || cbft.HasBlock(block.Hash(), block.NumberU64()) {
		cbft.log.Debug("The inserted block has exists in chain",
			"number", block.Number(), "hash", block.Hash(),
			"lockedNumber", cbft.state.HighestLockBlock().Number(),
			"lockedHash", cbft.state.HighestLockBlock().Hash())
		return
	}
	if cbft.state.Epoch() == qc.Epoch && cbft.state.ViewNumber() == qc.ViewNumber {
		if cbft.state.ViewBlockByIndex(qc.BlockIndex) == nil {
			cbft.state.AddQCBlock(block, qc)
			cbft.state.AddQC(qc)
		} else {
			cbft.state.AddQC(qc)
		}
	}

	lock, commit := cbft.blockTree.InsertQCBlock(block, qc)
	cbft.TrySetHighestQCBlock(block)
	isOwn := func() bool {
		node, err := cbft.isCurrentValidator()
		if err != nil {
			return false
		}
		proposer := cbft.currentProposer()
		// The current node is the proposer and the block is generated by itself
		if node.Index == proposer.Index && cbft.state.Epoch() == qc.Epoch && cbft.state.ViewNumber() == qc.ViewNumber {
			return true
		}
		return false
	}()
	if !isOwn {
		cbft.txPool.Reset(block)
	}
	cbft.tryCommitNewBlock(lock, commit, block)
	cbft.trySendPrepareVote()
	cbft.tryChangeView()
	if cbft.insertBlockQCHook != nil {
		// test hook
		cbft.insertBlockQCHook(block, qc)
	}
}

func (cbft *Cbft) TrySetHighestQCBlock(block *types.Block) {
	_, qc := cbft.blockTree.FindBlockAndQC(block.Hash(), block.NumberU64())
	h := cbft.state.HighestQCBlock()
	_, hqc := cbft.blockTree.FindBlockAndQC(h.Hash(), h.NumberU64())
	if hqc == nil || qc.HigherQuorumCert(hqc.BlockNumber, hqc.Epoch, hqc.ViewNumber) {
		cbft.state.SetHighestQCBlock(block)
	}
}

func (cbft *Cbft) insertPrepareQC(qc *ctypes.QuorumCert) {
	if qc != nil {
		block := cbft.state.ViewBlockByIndex(qc.BlockIndex)

		linked := func(blockNumber uint64) bool {
			if block != nil {
				parent, _ := cbft.blockTree.FindBlockAndQC(block.ParentHash(), block.NumberU64()-1)
				return parent != nil && cbft.state.HighestQCBlock().NumberU64()+1 == blockNumber
			}
			return false
		}
		hasExecuted := func() bool {
			if cbft.validatorPool.IsValidator(qc.Epoch, cbft.config.Option.Node.ID()) {
				return cbft.state.HadSendPrepareVote().Had(qc.BlockIndex) && linked(qc.BlockNumber)
			} else {
				blockIndex, finish := cbft.state.Executing()
				return blockIndex != math.MaxUint32 && (qc.BlockIndex < blockIndex || (qc.BlockIndex == blockIndex && finish)) && linked(qc.BlockNumber)
			}
			return false
		}
		if block != nil && hasExecuted() {
			cbft.log.Trace("Insert prepareQC", "qc", qc.String())
			cbft.insertQCBlock(block, qc)
		}
	}
}

// Asynchronous execution block callback function
func (cbft *Cbft) onAsyncExecuteStatus(s *executor.BlockExecuteStatus) {
	cbft.log.Debug("Async Execute Block", "hash", s.Hash, "number", s.Number)
	if s.Err != nil {
		cbft.log.Error("Execute block failed", "err", s.Err, "hash", s.Hash, "number", s.Number)
		return
	}
	index, finish := cbft.state.Executing()
	if !finish {
		block := cbft.state.ViewBlockByIndex(index)
		if block != nil {
			if block.Hash() == s.Hash {
				cbft.state.SetExecuting(index, true)
				if cbft.executeFinishHook != nil {
					cbft.executeFinishHook(index)
				}
				_, err := cbft.isCurrentValidator()
				if err != nil {
					cbft.log.Debug("Current node is not validator,no need to sign block", "err", err, "hash", s.Hash, "number", s.Number)
					return
				}
				if err := cbft.signBlock(block.Hash(), block.NumberU64(), index); err != nil {
					cbft.log.Error("Sign block failed", "err", err, "hash", s.Hash, "number", s.Number)
					return
				}

				cbft.log.Debug("Sign block", "hash", s.Hash, "number", s.Number)
				if msg := cbft.csPool.GetPrepareQC(cbft.state.Epoch(), cbft.state.ViewNumber(), index); msg != nil {
					go cbft.ReceiveSyncMsg(ctypes.NewInnerMsgInfo(msg.Msg, msg.PeerID))
				}
			}
		}
	}
	cbft.findQCBlock()
	cbft.findExecutableBlock()
}

// Sign the block that has been executed
// Every time try to trigger a send PrepareVote
func (cbft *Cbft) signBlock(hash common.Hash, number uint64, index uint32) error {
	// parentQC added when sending
	// Determine if the current consensus node is
	node, err := cbft.isCurrentValidator()
	if err != nil {
		return err
	}
	prepareVote := &protocols.PrepareVote{
		Epoch:          cbft.state.Epoch(),
		ViewNumber:     cbft.state.ViewNumber(),
		BlockHash:      hash,
		BlockNumber:    number,
		BlockIndex:     index,
		ValidatorIndex: node.Index,
	}

	if err := cbft.signMsgByBls(prepareVote); err != nil {
		return err
	}
	cbft.state.PendingPrepareVote().Push(prepareVote)
	// Record the number of participating consensus
	consensusCounter.Inc(1)

	cbft.trySendPrepareVote()
	return nil
}

// Send a signature,
// obtain a signature from the pending queue,
// determine whether the parent block has reached QC,
// and send a signature if it is reached, otherwise exit the sending logic.
func (cbft *Cbft) trySendPrepareVote() {
	// Check timeout
	if cbft.state.IsDeadline() {
		cbft.log.Debug("Current view had timeout, Refuse to send prepareVotes")
		return
	}

	pending := cbft.state.PendingPrepareVote()
	hadSend := cbft.state.HadSendPrepareVote()

	for !pending.Empty() {
		p := pending.Top()
		if err := cbft.voteRules.AllowVote(p); err != nil {
			cbft.log.Debug("Not allow send vote", "err", err, "msg", p.String())
			break
		}

		block := cbft.state.ViewBlockByIndex(p.BlockIndex)
		// The executed block has a signature.
		// Only when the view is switched, the block is cleared but the vote is also cleared.
		// If there is no block, the consensus process is abnormal and should not run.
		if block == nil {
			cbft.log.Crit("Try send PrepareVote failed", "err", "vote corresponding block not found", "view", cbft.state.ViewString(), "vote", p.String())
		}
		if b, qc := cbft.blockTree.FindBlockAndQC(block.ParentHash(), block.NumberU64()-1); b != nil || block.NumberU64() == 0 {
			p.ParentQC = qc
			hadSend.Push(p)
			//Determine if the current consensus node is
			node, _ := cbft.isCurrentValidator()
			cbft.log.Info("Add local prepareVote", "vote", p.String())
			cbft.state.AddPrepareVote(node.Index, p)
			pending.Pop()

			// write sendPrepareVote info to wal
			if !cbft.isLoading() {
				cbft.bridge.SendPrepareVote(block, p)
			}

			// send prepareVote use pubsub
			if err := cbft.publishTopicMsg(p); err != nil {
				cbft.log.Error("Publish PrepareVote failed", "err", err.Error(), "view", cbft.state.ViewString(), "vote", p.String())
			}
			//cbft.network.Broadcast(p)
		} else {
			break
		}
	}
}

func (cbft *Cbft) publishTopicMsg(msg ctypes.ConsensusMsg) error {
	if cbft.NeedGroup() {
		groupID, _, err := cbft.getGroupByValidatorID(cbft.state.Epoch(), cbft.Node().ID())
		if err != nil {
			return fmt.Errorf("the group info of the current node is not queried, cannot publish the topic message")
		}
		network.MeteredWriteRGMsg(protocols.MessageType(msg), msg)
		topic := cbfttypes.ConsensusGroupTopicName(cbft.state.Epoch(), groupID)
		return cbft.pubSub.Publish(topic, protocols.MessageType(msg), msg)
	} else {
		cbft.network.Broadcast(msg)
		return nil
	}
}

func (cbft *Cbft) trySendRGBlockQuorumCert() {
	if !cbft.NeedGroup() {
		return
	}
	// Check timeout
	if cbft.state.IsDeadline() {
		cbft.log.Debug("Current view had timeout, refuse to send RGBlockQuorumCert")
		return
	}

	//node, err := cbft.isCurrentValidator()
	//if err != nil || node == nil {
	//	cbft.log.Debug("Current node is not validator, no need to send RGBlockQuorumCert")
	//	return
	//}

	groupID, _, err := cbft.getGroupByValidatorID(cbft.state.Epoch(), cbft.Node().ID())
	if err != nil {
		cbft.log.Debug("Current node is not validator, no need to send RGBlockQuorumCert")
		return
	}

	enoughVotes := func(blockIndex, groupID uint32) bool {
		threshold := cbft.groupThreshold(cbft.state.Epoch(), groupID)
		groupVotes := cbft.groupPrepareVotes(cbft.state.Epoch(), blockIndex, groupID)
		if len(groupVotes) >= threshold {
			// generatePrepareQC by group votes
			rgqc := cbft.generatePrepareQC(groupVotes)
			// get parentQC
			var parentQC *ctypes.QuorumCert
			for _, v := range groupVotes {
				parentQC = v.ParentQC
				break
			}
			// Add SelectRGQuorumCerts
			cbft.state.AddSelectRGQuorumCerts(blockIndex, groupID, rgqc, parentQC)
			blockGroupQCBySelfCounter.Inc(1)
			return true
		}
		return false
	}

	alreadyRGBlockQuorumCerts := func(blockIndex, groupID uint32) bool {
		len := cbft.state.SelectRGQuorumCertsLen(blockIndex, groupID)
		if len > 0 {
			blockGroupQCByOtherCounter.Inc(1)
			return true
		}
		return false
	}

	for index := uint32(0); index <= cbft.state.MaxViewBlockIndex(); index++ {
		if cbft.state.HadSendRGBlockQuorumCerts(index) {
			cbft.log.Trace("RGBlockQuorumCert has been sent, no need to send again", "blockIndex", index, "groupID", groupID)
			continue
		}

		if alreadyRGBlockQuorumCerts(index, groupID) || enoughVotes(index, groupID) {
			cbft.RGBroadcastManager.AsyncSendRGQuorumCert(&awaitingRGBlockQC{
				groupID:    groupID,
				blockIndex: index,
				epoch:      cbft.state.Epoch(),
				viewNumber: cbft.state.ViewNumber(),
			})
			cbft.state.AddSendRGBlockQuorumCerts(index)
			cbft.log.Debug("Send RGBlockQuorumCert asynchronously", "blockIndex", index, "groupID", groupID)
			// record metrics
			block := cbft.state.ViewBlockByIndex(index)
			blockGroupQCTimer.UpdateSince(time.Unix(int64(block.Time()), 0))
		}
	}
}

func (cbft *Cbft) trySendRGViewChangeQuorumCert() {
	if !cbft.NeedGroup() {
		return
	}

	groupID, _, err := cbft.getGroupByValidatorID(cbft.state.Epoch(), cbft.Node().ID())
	if err != nil {
		cbft.log.Debug("Current node is not validator, no need to send RGViewChangeQuorumCert")
		return
	}

	enoughViewChanges := func(groupID uint32) bool {
		threshold := cbft.groupThreshold(cbft.state.Epoch(), groupID)
		groupViewChanges := cbft.groupViewChanges(cbft.state.Epoch(), groupID)
		if len(groupViewChanges) >= threshold {
			// generatePrepareQC by group votes
			rgqc := cbft.generateViewChangeQC(groupViewChanges)
			// get parentQC
			prepareQCs := make(map[common.Hash]*ctypes.QuorumCert)
			for _, v := range groupViewChanges {
				if v.PrepareQC != nil {
					prepareQCs[v.BlockHash] = v.PrepareQC
				}
			}
			// Add SelectRGViewChangeQuorumCerts
			cbft.state.AddSelectRGViewChangeQuorumCerts(groupID, rgqc, prepareQCs)
			viewGroupQCBySelfCounter.Inc(1)
			return true
		}
		return false
	}

	alreadyRGViewChangeQuorumCerts := func(groupID uint32) bool {
		len := cbft.state.SelectRGViewChangeQuorumCertsLen(groupID)
		if len > 0 {
			viewGroupQCByOtherCounter.Inc(1)
			return true
		}
		return false
	}

	if cbft.state.HadSendRGViewChangeQuorumCerts(cbft.state.ViewNumber()) {
		cbft.log.Trace("RGViewChangeQuorumCert has been sent, no need to send again", "groupID", groupID)
		return
	}

	if alreadyRGViewChangeQuorumCerts(groupID) || enoughViewChanges(groupID) {
		cbft.RGBroadcastManager.AsyncSendRGQuorumCert(&awaitingRGViewQC{
			groupID:    groupID,
			epoch:      cbft.state.Epoch(),
			viewNumber: cbft.state.ViewNumber(),
		})
		cbft.state.AddSendRGViewChangeQuorumCerts(cbft.state.ViewNumber())
		cbft.log.Debug("Send RGViewChangeQuorumCert asynchronously", "groupID", groupID)
		viewGroupQCTimer.UpdateSince(cbft.state.Deadline())
	}
}

// Every time there is a new block or a new executed block result will enter this judgment, find the next executable block
func (cbft *Cbft) findExecutableBlock() {
	qcIndex := cbft.state.MaxQCIndex()
	blockIndex, finish := cbft.state.Executing()

	// If we are not execute block yet and the QC index
	// is greater 0, then starting execute block from qc index.
	if blockIndex == math.MaxUint32 && qcIndex != math.MaxUint32 {
		blockIndex, finish = qcIndex, true
	}

	if blockIndex == math.MaxUint32 {
		block := cbft.state.ViewBlockByIndex(blockIndex + 1)
		if block != nil {
			parent, _ := cbft.blockTree.FindBlockAndQC(block.ParentHash(), block.NumberU64()-1)
			if parent == nil {
				cbft.log.Error(fmt.Sprintf("Find executable block's parent failed :[%d,%d,%s]", blockIndex, block.NumberU64(), block.Hash().String()))
				return
			}

			cbft.log.Debug("Find Executable Block", "hash", block.Hash(), "number", block.NumberU64())
			if err := cbft.asyncExecutor.Execute(block, parent); err != nil {
				cbft.log.Error("Async Execute block failed", "error", err)
			}
			cbft.state.SetExecuting(0, false)
		}
	}

	if finish {
		block := cbft.state.ViewBlockByIndex(blockIndex + 1)
		if block != nil {
			parent := cbft.state.ViewBlockByIndex(blockIndex)
			if parent == nil {
				cbft.log.Error(fmt.Sprintf("Find executable block's parent failed :[%d,%d,%s]", blockIndex, block.NumberU64(), block.Hash()))
				return
			}

			cbft.log.Debug("Find Executable Block", "hash", block.Hash(), "number", block.NumberU64())
			if err := cbft.asyncExecutor.Execute(block, parent); err != nil {
				cbft.log.Error("Async Execute block failed", "error", err)
			}
			cbft.state.SetExecuting(blockIndex+1, false)
		}
	}
}

func (cbft *Cbft) mergeVoteToQuorumCerts(node *cbfttypes.ValidateNode, vote *protocols.PrepareVote) {
	if !cbft.NeedGroup() {
		return
	}
	// Query which group the PrepareVote belongs to according to the information of the node sending the PrepareVote message
	groupID, _, err := cbft.getGroupByValidatorID(vote.EpochNum(), node.NodeID)
	if err != nil {
		cbft.log.Error("Failed to find the group info of the node", "epoch", vote.EpochNum(), "nodeID", node.NodeID.TerminalString(), "error", err)
		return
	}
	cbft.state.MergePrepareVotes(vote.BlockIndex, groupID, []*protocols.PrepareVote{vote})
}

func (cbft *Cbft) richBlockQuorumCert(epoch uint64, blockIndex, groupID uint32, blockQC *ctypes.QuorumCert) {
	mergeVotes := cbft.groupPrepareVotes(epoch, blockIndex, groupID)
	if len(mergeVotes) > 0 {
		for _, v := range mergeVotes {
			if !blockQC.HasSign(v.NodeIndex()) {
				blockQC.AddSign(v.Signature, v.NodeIndex())
			}
		}
	}
}

func (cbft *Cbft) mergeViewChangeToViewChangeQuorumCerts(node *cbfttypes.ValidateNode, vc *protocols.ViewChange) {
	if !cbft.NeedGroup() {
		return
	}
	// Query which group the ViewChange belongs to according to the information of the node sending the ViewChange message
	groupID, _, err := cbft.getGroupByValidatorID(vc.EpochNum(), node.NodeID)
	if err != nil {
		cbft.log.Error("Failed to find the group info of the node", "epoch", vc.EpochNum(), "nodeID", node.NodeID.TerminalString(), "error", err)
		return
	}
	validatorLen := cbft.currentValidatorLen()
	cbft.state.MergeViewChanges(groupID, []*protocols.ViewChange{vc}, validatorLen)
}

func (cbft *Cbft) richViewChangeQuorumCert(epoch uint64, groupID uint32, viewChangeQC *ctypes.ViewChangeQC, prepareQCs *ctypes.PrepareQCs) {
	mergeVcs := cbft.groupViewChanges(epoch, groupID)
	if len(mergeVcs) > 0 {
		for _, vc := range mergeVcs {
			cbft.MergeViewChange(viewChangeQC, vc)
			if !viewChangeQC.ExistViewChange(vc.Epoch, vc.ViewNumber, vc.BlockHash) {
				if vc.PrepareQC != nil {
					prepareQCs.AppendQuorumCert(vc.PrepareQC)
				}
			}
		}
	}
}

func (cbft *Cbft) MergeViewChange(qcs *ctypes.ViewChangeQC, vc *protocols.ViewChange) {
	validatorLen := cbft.currentValidatorLen()

	if !qcs.ExistViewChange(vc.Epoch, vc.ViewNumber, vc.BlockHash) {
		qc := &ctypes.ViewChangeQuorumCert{
			Epoch:        vc.Epoch,
			ViewNumber:   vc.ViewNumber,
			BlockHash:    vc.BlockHash,
			BlockNumber:  vc.BlockNumber,
			ValidatorSet: utils.NewBitArray(uint32(validatorLen)),
		}
		if vc.PrepareQC != nil {
			qc.BlockEpoch = vc.PrepareQC.Epoch
			qc.BlockViewNumber = vc.PrepareQC.ViewNumber
			//rgb.PrepareQCs.AppendQuorumCert(vc.PrepareQC)
		}
		qc.ValidatorSet.SetIndex(vc.ValidatorIndex, true)
		qc.Signature.SetBytes(vc.Signature.Bytes())

		qcs.AppendQuorumCert(qc)
	} else {
		for _, qc := range qcs.QCs {
			if qc.BlockHash == vc.BlockHash && !qc.HasSign(vc.NodeIndex()) {
				qc.AddSign(vc.Signature, vc.NodeIndex())
				break
			}
		}
	}
}

// Return all votes of the specified group under the current view
func (cbft *Cbft) groupPrepareVotes(epoch uint64, blockIndex, groupID uint32) map[uint32]*protocols.PrepareVote {
	indexes, err := cbft.validatorPool.GetValidatorIndexesByGroupID(epoch, groupID)
	if err != nil || indexes == nil {
		return nil
	}
	// Find votes corresponding to indexes
	votes := cbft.state.AllPrepareVoteByIndex(blockIndex)
	if len(votes) > 0 {
		groupVotes := make(map[uint32]*protocols.PrepareVote)
		for _, index := range indexes {
			if vote, ok := votes[index]; ok {
				groupVotes[index] = vote
			}
		}
		return groupVotes
	}
	return nil
}

// Return all viewChanges of the specified group under the current view
func (cbft *Cbft) groupViewChanges(epoch uint64, groupID uint32) map[uint32]*protocols.ViewChange {
	indexes, err := cbft.validatorPool.GetValidatorIndexesByGroupID(epoch, groupID)
	if err != nil || indexes == nil {
		return nil
	}
	// Find viewChanges corresponding to indexes
	vcs := cbft.state.AllViewChange()
	if len(vcs) > 0 {
		groupVcs := make(map[uint32]*protocols.ViewChange)
		for _, index := range indexes {
			if vc, ok := vcs[index]; ok {
				groupVcs[index] = vc
			}
		}
		return groupVcs
	}
	return nil
}

// Each time a new vote is triggered, a new QC Block will be triggered, and a new one can be found by the commit block.
func (cbft *Cbft) findQCBlock() {
	index := cbft.state.MaxQCIndex()
	next := index + 1
	threshold := cbft.threshold(cbft.currentValidatorLen())

	var qc *ctypes.QuorumCert

	enoughVotes := func() bool {
		size := cbft.state.PrepareVoteLenByIndex(next)
		if size >= threshold {
			qc = cbft.generatePrepareQC(cbft.state.AllPrepareVoteByIndex(next))
			blockWholeQCByVotesCounter.Inc(1)
			cbft.log.Debug("Enough prepareVote have been received, generate prepareQC", "qc", qc.String())
		}
		return qc.Len() >= threshold
	}

	enoughRGQuorumCerts := func() bool {
		if !cbft.NeedGroup() {
			return false
		}
		rgqcs := cbft.state.FindMaxRGQuorumCerts(next)
		size := 0
		if len(rgqcs) > 0 {
			for _, qc := range rgqcs {
				size += qc.Len()
			}
		}
		if size >= threshold {
			qc = cbft.combinePrepareQC(rgqcs)
			blockWholeQCByRGQCCounter.Inc(1)
			cbft.log.Debug("Enough RGBlockQuorumCerts have been received, combine prepareQC", "qc", qc.String())
		}
		return qc.Len() >= threshold
	}

	enoughCombine := func() bool {
		if !cbft.NeedGroup() {
			return false
		}
		knownIndexes := cbft.KnownVoteIndexes(next)
		if len(knownIndexes) >= threshold {
			rgqcs := cbft.state.FindMaxRGQuorumCerts(next)
			qc = cbft.combinePrepareQC(rgqcs)
			cbft.log.Trace("enoughCombine rgqcs", "qc", qc.String())

			allVotes := cbft.state.AllPrepareVoteByIndex(next)
			for index, v := range allVotes {
				if qc.Len() >= threshold {
					break // The merge can be stopped when the threshold is reached
				}
				if !qc.HasSign(index) {
					cbft.log.Trace("enoughCombine add vote", "index", v.NodeIndex(), "vote", v.String())
					qc.AddSign(v.Signature, v.NodeIndex())
				}
			}
			blockWholeQCByCombineCounter.Inc(1)
			cbft.log.Debug("Enough RGBlockQuorumCerts and prepareVote have been received, combine prepareQC", "qc", qc.String())
		}
		return qc.Len() >= threshold
	}

	linked := func(blockIndex uint32) bool {
		block := cbft.state.ViewBlockByIndex(blockIndex)
		if block != nil {
			parent, _ := cbft.blockTree.FindBlockAndQC(block.ParentHash(), block.NumberU64()-1)
			return parent != nil && cbft.state.HighestQCBlock().NumberU64()+1 == block.NumberU64()
		}
		return false
	}

	hasExecuted := func(blockIndex uint32) bool {
		if cbft.validatorPool.IsValidator(cbft.state.Epoch(), cbft.config.Option.Node.ID()) {
			return cbft.state.HadSendPrepareVote().Had(blockIndex)
		} else {
			executingIndex, finish := cbft.state.Executing()
			return blockIndex != math.MaxUint32 && (blockIndex < executingIndex || (executingIndex == blockIndex && finish)) && linked(blockIndex)
		}
		return false
	}

	alreadyQC := func() bool {
		return hasExecuted(next) && (enoughRGQuorumCerts() || enoughVotes() || enoughCombine())
	}

	if alreadyQC() {
		block := cbft.state.ViewBlockByIndex(next)
		if qc != nil {
			cbft.log.Info("New qc block have been created", "qc", qc.String())
			cbft.insertQCBlock(block, qc)
			cbft.network.Broadcast(&protocols.BlockQuorumCert{BlockQC: qc})
			// metrics
			blockQCCollectedGauage.Update(int64(block.Time()))
			blockWholeQCTimer.UpdateSince(time.Unix(int64(block.Time()), 0))
			cbft.trySendPrepareVote()
		}
	}

	cbft.tryChangeView()
}

// Try commit a new block
func (cbft *Cbft) tryCommitNewBlock(lock *types.Block, commit *types.Block, qc *types.Block) {
	if lock == nil || commit == nil {
		cbft.log.Warn("Try commit failed", "hadLock", lock != nil, "hadCommit", commit != nil)
		return
	}
	if commit.NumberU64()+2 != qc.NumberU64() {
		cbft.log.Warn("Try commit failed,the requirement of commit block is not achieved", "commit", commit.NumberU64(), "lock", lock.NumberU64(), "qc", qc.NumberU64())
		return
	}
	//highestqc := cbft.state.HighestQCBlock()
	highestqc := qc
	commitBlock, commitQC := cbft.blockTree.FindBlockAndQC(commit.Hash(), commit.NumberU64())

	_, oldCommit := cbft.state.HighestLockBlock(), cbft.state.HighestCommitBlock()
	// Incremental commit block
	if oldCommit.NumberU64()+1 == commit.NumberU64() {
		cbft.commitBlock(commit, commitQC, lock, highestqc)
		cbft.state.SetHighestLockBlock(lock)
		cbft.state.SetHighestCommitBlock(commit)
		cbft.blockTree.PruneBlock(commit.Hash(), commit.NumberU64(), nil)
		cbft.blockTree.NewRoot(commit)
		// metrics
		blockNumberGauage.Update(int64(commit.NumberU64()))
		highestQCNumberGauage.Update(int64(highestqc.NumberU64()))
		highestLockedNumberGauage.Update(int64(lock.NumberU64()))
		highestCommitNumberGauage.Update(int64(commit.NumberU64()))
		blockConfirmedMeter.Mark(1)
	} else if oldCommit.NumberU64() == commit.NumberU64() && oldCommit.NumberU64() > 0 {
		cbft.log.Info("Fork block", "number", highestqc.NumberU64(), "hash", highestqc.Hash())
		lockBlock, lockQC := cbft.blockTree.FindBlockAndQC(lock.Hash(), lock.NumberU64())
		qcBlock, qcQC := cbft.blockTree.FindBlockAndQC(highestqc.Hash(), highestqc.NumberU64())

		qcState := &protocols.State{Block: qcBlock, QuorumCert: qcQC}
		lockState := &protocols.State{Block: lockBlock, QuorumCert: lockQC}
		commitState := &protocols.State{Block: commitBlock, QuorumCert: commitQC}
		cbft.bridge.UpdateChainState(qcState, lockState, commitState)
	}
}

// Asynchronous processing of errors generated by the submission block
func (cbft *Cbft) OnCommitError(err error) {
	// FIXME Do you want to panic and stop the consensus?
	cbft.log.Error("Commit block error", "err", err)
}

// According to the current view QC situation, try to switch view
func (cbft *Cbft) tryChangeView() {
	// Had receive all qcs of current view
	block, qc := cbft.blockTree.FindBlockAndQC(cbft.state.HighestQCBlock().Hash(), cbft.state.HighestQCBlock().NumberU64())

	increasing := func() uint64 {
		return cbft.state.ViewNumber() + 1
	}

	shouldSwitch := func() bool {
		return cbft.validatorPool.ShouldSwitch(block.NumberU64())
	}()

	enough := func() bool {
		return cbft.state.MaxQCIndex()+1 == cbft.config.Sys.Amount ||
			(qc != nil && qc.Epoch == cbft.state.Epoch() && shouldSwitch)
	}()

	activeVersion, err := cbft.blockCache.GetActiveVersion(block.Header())
	if err != nil {
		log.Error("GetActiveVersion failed", "err", err)
	}

	if shouldSwitch {
		if err := cbft.validatorPool.Update(block.Hash(), block.NumberU64(), cbft.state.Epoch()+1, activeVersion, cbft.eventMux); err == nil {
			cbft.log.Info("Update validator success", "number", block.NumberU64())
		} else {
			cbft.log.Trace("Update validator failed!", "number", block.NumberU64(), "err", err)
		}
	}

	if enough {
		// If current has produce enough blocks, then change view immediately.
		// Otherwise, waiting for view's timeout.
		if cbft.validatorPool.EqualSwitchPoint(block.NumberU64()) {
			cbft.log.Info("BlockNumber is equal to switchPoint, change epoch", "blockNumber", block.NumberU64(), "view", cbft.state.ViewString())
			cbft.changeView(cbft.state.Epoch()+1, state.DefaultViewNumber, block, qc, nil)
		} else {
			cbft.log.Info("Produce enough blocks, change view", "view", cbft.state.ViewString())
			cbft.changeView(cbft.state.Epoch(), increasing(), block, qc, nil)
		}
		return
	}

	threshold := cbft.threshold(cbft.currentValidatorLen())
	var viewChangeQC *ctypes.ViewChangeQC

	enoughViewChanges := func() bool {
		size := cbft.state.ViewChangeLen()
		if size >= threshold {
			viewChangeQC = cbft.generateViewChangeQC(cbft.state.AllViewChange())
			viewWholeQCByVcsCounter.Inc(1)
			cbft.log.Info("Receive enough viewChange, generate viewChangeQC", "viewChangeQC", viewChangeQC.String())
		}
		return viewChangeQC.HasLength() >= threshold
	}

	enoughRGQuorumCerts := func() bool {
		if !cbft.NeedGroup() {
			return false
		}
		viewChangeQCs := cbft.state.FindMaxRGViewChangeQuorumCerts()
		size := 0
		if len(viewChangeQCs) > 0 {
			for _, qcs := range viewChangeQCs {
				size += qcs.HasLength()
			}
		}
		if size >= threshold {
			viewChangeQC = cbft.combineViewChangeQC(viewChangeQCs)
			viewWholeQCByRGQCCounter.Inc(1)
			cbft.log.Debug("Enough RGViewChangeQuorumCerts have been received, combine ViewChangeQC", "viewChangeQC", viewChangeQC.String())
		}
		return viewChangeQC.HasLength() >= threshold
	}

	enoughCombine := func() bool {
		if !cbft.NeedGroup() {
			return false
		}
		knownIndexes := cbft.KnownViewChangeIndexes()
		if len(knownIndexes) >= threshold {
			viewChangeQCs := cbft.state.FindMaxRGViewChangeQuorumCerts()
			viewChangeQC = cbft.combineViewChangeQC(viewChangeQCs)
			cbft.log.Trace("enoughCombine viewChangeQCs", "viewChangeQC", viewChangeQC.String())

			allViewChanges := cbft.state.AllViewChange()
			for index, v := range allViewChanges {
				if viewChangeQC.HasLength() >= threshold {
					break // The merge can be stopped when the threshold is reached
				}
				if !viewChangeQC.HasSign(index) {
					cbft.log.Trace("enoughCombine add viewChange", "index", v.NodeIndex(), "viewChange", v.String())
					cbft.MergeViewChange(viewChangeQC, v)
				}
			}
			viewWholeQCByCombineCounter.Inc(1)
			cbft.log.Debug("Enough RGViewChangeQuorumCerts and viewChange have been received, combine ViewChangeQC", "viewChangeQC", viewChangeQC.String())
		}
		return viewChangeQC.HasLength() >= threshold
	}

	alreadyQC := func() bool {
		return enoughRGQuorumCerts() || enoughViewChanges() || enoughCombine()
	}

	if alreadyQC() {
		viewWholeQCTimer.UpdateSince(cbft.state.Deadline())
		cbft.tryChangeViewByViewChange(viewChangeQC)
	}
}

func (cbft *Cbft) richViewChangeQC(viewChangeQC *ctypes.ViewChangeQC) {
	node, err := cbft.isCurrentValidator()
	if err != nil {
		cbft.log.Info("Local node is not validator")
		return
	}
	hadSend := cbft.state.ViewChangeByIndex(node.Index)
	if hadSend != nil && !viewChangeQC.ExistViewChange(hadSend.Epoch, hadSend.ViewNumber, hadSend.BlockHash) {
		cert, err := cbft.generateViewChangeQuorumCert(hadSend)
		if err != nil {
			cbft.log.Error("Generate viewChangeQuorumCert error", "err", err)
			return
		}
		cbft.log.Info("Already send viewChange, append viewChangeQuorumCert to ViewChangeQC", "cert", cert.String())
		viewChangeQC.AppendQuorumCert(cert)
	}

	_, qc := cbft.blockTree.FindBlockAndQC(cbft.state.HighestQCBlock().Hash(), cbft.state.HighestQCBlock().NumberU64())
	_, _, blockEpoch, blockView, _, number := viewChangeQC.MaxBlock()

	if qc.HigherQuorumCert(number, blockEpoch, blockView) {
		if hadSend == nil {
			v, err := cbft.generateViewChange(qc)
			if err != nil {
				cbft.log.Error("Generate viewChange error", "err", err)
				return
			}
			cert, err := cbft.generateViewChangeQuorumCert(v)
			if err != nil {
				cbft.log.Error("Generate viewChangeQuorumCert error", "err", err)
				return
			}
			cbft.log.Info("Not send viewChange, append viewChangeQuorumCert to ViewChangeQC", "cert", cert.String())
			viewChangeQC.AppendQuorumCert(cert)
		}
	}
}

func (cbft *Cbft) tryChangeViewByViewChange(viewChangeQC *ctypes.ViewChangeQC) {
	increasing := func() uint64 {
		return cbft.state.ViewNumber() + 1
	}

	_, _, blockEpoch, _, hash, number := viewChangeQC.MaxBlock()
	block, qc := cbft.blockTree.FindBlockAndQC(cbft.state.HighestQCBlock().Hash(), cbft.state.HighestQCBlock().NumberU64())
	if block.NumberU64() != 0 {
		cbft.richViewChangeQC(viewChangeQC)
		_, _, blockEpoch, _, hash, number = viewChangeQC.MaxBlock()
		block, qc := cbft.blockTree.FindBlockAndQC(hash, number)
		if block == nil || qc == nil {
			// fixme get qc block
			cbft.log.Warn("Local node is behind other validators", "blockState", cbft.state.HighestBlockString(), "viewChangeQC", viewChangeQC.String())
			return
		}
	}

	if cbft.validatorPool.EqualSwitchPoint(number) && blockEpoch == cbft.state.Epoch() {
		// Validator already switch, new epoch
		cbft.log.Info("BlockNumber is equal to switchPoint, change epoch", "blockNumber", number, "view", cbft.state.ViewString())
		cbft.changeView(cbft.state.Epoch()+1, state.DefaultViewNumber, block, qc, viewChangeQC)
	} else {
		cbft.changeView(cbft.state.Epoch(), increasing(), block, qc, viewChangeQC)
	}
}

func (cbft *Cbft) generateViewChangeQuorumCert(v *protocols.ViewChange) (*ctypes.ViewChangeQuorumCert, error) {
	node, err := cbft.isCurrentValidator()
	if err != nil {
		return nil, errors.Wrap(err, "local node is not validator")
	}
	total := uint32(cbft.currentValidatorLen())
	var aggSig bls.Sign
	if err := aggSig.Deserialize(v.Sign()); err != nil {
		return nil, err
	}

	blockEpoch, blockView := uint64(0), uint64(0)
	if v.PrepareQC != nil {
		blockEpoch, blockView = v.PrepareQC.Epoch, v.PrepareQC.ViewNumber
	}
	cert := &ctypes.ViewChangeQuorumCert{
		Epoch:           v.Epoch,
		ViewNumber:      v.ViewNumber,
		BlockHash:       v.BlockHash,
		BlockNumber:     v.BlockNumber,
		BlockEpoch:      blockEpoch,
		BlockViewNumber: blockView,
		ValidatorSet:    utils.NewBitArray(total),
	}
	cert.Signature.SetBytes(aggSig.Serialize())
	cert.ValidatorSet.SetIndex(node.Index, true)
	return cert, nil
}

func (cbft *Cbft) generateViewChange(qc *ctypes.QuorumCert) (*protocols.ViewChange, error) {
	node, err := cbft.isCurrentValidator()
	if err != nil {
		return nil, errors.Wrap(err, "local node is not validator")
	}
	v := &protocols.ViewChange{
		Epoch:          cbft.state.Epoch(),
		ViewNumber:     cbft.state.ViewNumber(),
		BlockHash:      qc.BlockHash,
		BlockNumber:    qc.BlockNumber,
		ValidatorIndex: node.Index,
		PrepareQC:      qc,
	}
	if err := cbft.signMsgByBls(v); err != nil {
		return nil, errors.Wrap(err, "Sign ViewChange failed")
	}

	return v, nil
}

// change view
func (cbft *Cbft) changeView(epoch, viewNumber uint64, block *types.Block, qc *ctypes.QuorumCert, viewChangeQC *ctypes.ViewChangeQC) {
	interval := func() uint64 {
		if block.NumberU64() == 0 {
			return viewNumber - state.DefaultViewNumber + 1
		}
		if qc.ViewNumber+1 == viewNumber {
			return uint64((cbft.config.Sys.Amount-qc.BlockIndex)/3) + 1
		}
		minuend := qc.ViewNumber
		if qc.Epoch != epoch {
			minuend = state.DefaultViewNumber
		}
		return viewNumber - minuend + 1
	}
	// last epoch and last viewNumber
	// when cbft is started or fast synchronization ends, the preEpoch, preViewNumber defaults to 0, 0
	// but cbft is now in the loading state and lastViewChangeQC is nil, does not save the lastViewChangeQC
	preEpoch, preViewNumber := cbft.state.Epoch(), cbft.state.ViewNumber()
	// syncingCache is belong to last view request, clear all sync cache
	cbft.syncingCache.Purge()
	cbft.csPool.Purge(epoch, viewNumber)

	cbft.state.ResetView(epoch, viewNumber)
	cbft.state.SetViewTimer(interval())
	cbft.state.SetLastViewChangeQC(viewChangeQC)

	// record metrics
	viewNumberGauage.Update(int64(viewNumber))
	epochNumberGauage.Update(int64(epoch))
	viewChangedTimer.UpdateSince(time.Unix(int64(block.Time()), 0))

	// write confirmed viewChange info to wal
	if !cbft.isLoading() {
		cbft.bridge.ConfirmViewChange(epoch, viewNumber, block, qc, viewChangeQC, preEpoch, preViewNumber)
	}
	cbft.RGBroadcastManager.Reset()
	cbft.clearInvalidBlocks(block)
	cbft.evPool.Clear(epoch, viewNumber)
	// view change maybe lags behind the other nodes,active sync prepare block
	cbft.SyncPrepareBlock("", epoch, viewNumber, 0)
	cbft.log = log.New("epoch", cbft.state.Epoch(), "view", cbft.state.ViewNumber())
	cbft.log.Info("Success to change view, current view deadline", "deadline", cbft.state.Deadline())
}

// Clean up invalid blocks in the previous view
func (cbft *Cbft) clearInvalidBlocks(newBlock *types.Block) {
	var rollback []*types.Block
	newHead := newBlock.Header()
	for _, p := range cbft.state.HadSendPrepareVote().Peek() {
		if p.BlockNumber > newBlock.NumberU64() {
			block := cbft.state.ViewBlockByIndex(p.BlockIndex)
			rollback = append(rollback, block)
		}
	}
	for _, p := range cbft.state.PendingPrepareVote().Peek() {
		if p.BlockNumber > newBlock.NumberU64() {
			block := cbft.state.ViewBlockByIndex(p.BlockIndex)
			rollback = append(rollback, block)
		}
	}
	cbft.blockCache.ClearCache(cbft.state.HighestCommitBlock())

	//todo proposer is myself
	cbft.txPool.ForkedReset(newHead, rollback)
}
