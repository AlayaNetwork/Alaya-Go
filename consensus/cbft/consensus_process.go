package cbft

import (
	"fmt"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/pkg/errors"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/executor"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

// OnPrepareBlock performs security rule verification，store in blockTree,
// Whether to start synchronization
func (cbft *Cbft) OnPrepareBlock(id string, msg *protocols.PrepareBlock) error {
	cbft.log.Debug("Receive PrepareBlock", "id", id, "msg", msg.String())
	if err := cbft.safetyRules.PrepareBlockRules(msg); err != nil {
		blockCheckFailureMeter.Mark(1)

		if err.Common() {
			cbft.csPool.AddPrepareBlock(msg.BlockIndex, &ctypes.MsgInfo{PeerID: id, Msg: msg})
			cbft.log.Debug("Prepare block rules fail", "number", msg.Block.Number(), "hash", msg.Block.Hash(), "err", err)
			return err
		}
		// verify consensus signature
		if err := cbft.verifyConsensusSign(msg); err != nil {
			signatureCheckFailureMeter.Mark(1)
			return err
		}

		cbft.csPool.AddPrepareBlock(msg.BlockIndex, &ctypes.MsgInfo{PeerID: id, Msg: msg})

		if err.Fetch() {
			if cbft.isProposer(msg.Epoch, msg.ViewNumber, msg.ProposalIndex) {
				cbft.fetchBlock(id, msg.Block.ParentHash(), msg.Block.NumberU64()-1, nil)
			}
			return err
		}
		if err.FetchPrepare() {
			cbft.prepareBlockFetchRules(id, msg)
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

	var node *cbfttypes.ValidateNode
	var err error
	if node, err = cbft.verifyConsensusMsg(msg); err != nil {
		signatureCheckFailureMeter.Mark(1)
		return err
	}

	if err := cbft.evPool.AddPrepareBlock(msg, node); err != nil {
		if _, ok := err.(*evidence.DuplicatePrepareBlockEvidence); ok {
			cbft.log.Warn("Receive DuplicatePrepareBlockEvidence msg", "err", err.Error())
			return err
		}
	}

	// The new block is notified by the PrepareBlockHash to the nodes in the network.
	cbft.state.AddPrepareBlock(msg)
	cbft.findExecutableBlock()
	return nil
}

// OnPrepareVote perform security rule verification，store in blockTree,
// Whether to start synchronization
func (cbft *Cbft) OnPrepareVote(id string, msg *protocols.PrepareVote) error {
	if err := cbft.safetyRules.PrepareVoteRules(msg); err != nil {

		if err.Common() {
			cbft.csPool.AddPrepareVote(msg.BlockIndex, msg.ValidatorIndex, &ctypes.MsgInfo{PeerID: id, Msg: msg})
			cbft.log.Debug("Preparevote rules fail", "number", msg.BlockHash, "hash", msg.BlockHash, "err", err)
			return err
		}

		// verify consensus signature
		if cbft.verifyConsensusSign(msg) != nil {
			signatureCheckFailureMeter.Mark(1)
			return err
		}

		cbft.csPool.AddPrepareVote(msg.BlockIndex, msg.ValidatorIndex, &ctypes.MsgInfo{PeerID: id, Msg: msg})

		if err.Fetch() {
			if msg.ParentQC != nil {
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
		return err
	}

	if err := cbft.evPool.AddPrepareVote(msg, node); err != nil {
		if _, ok := err.(*evidence.DuplicatePrepareVoteEvidence); ok {
			cbft.log.Warn("Receive DuplicatePrepareVoteEvidence msg", "err", err.Error())
			return err
		}
	}

	cbft.insertPrepareQC(msg.ParentQC)

	cbft.state.AddPrepareVote(uint32(node.Index), msg)
	cbft.log.Debug("Add prepare vote", "msgHash", msg.MsgHash(), "validatorIndex", msg.ValidatorIndex, "blockIndex", msg.BlockIndex, "number", msg.BlockNumber, "hash", msg.BlockHash, "votes", cbft.state.PrepareVoteLenByIndex(msg.BlockIndex))

	cbft.findQCBlock()
	return nil
}

// OnViewChange performs security rule verification, view switching.
func (cbft *Cbft) OnViewChange(id string, msg *protocols.ViewChange) error {
	cbft.log.Debug("Receive ViewChange", "msg", msg.String())
	if err := cbft.safetyRules.ViewChangeRules(msg); err != nil {
		if err.Fetch() {
			if msg.PrepareQC != nil {
				cbft.fetchBlock(id, msg.BlockHash, msg.BlockNumber, msg.PrepareQC)
			}
		}
		return err
	}

	var node *cbfttypes.ValidateNode
	var err error

	if node, err = cbft.verifyConsensusMsg(msg); err != nil {
		return err
	}

	if err := cbft.evPool.AddViewChange(msg, node); err != nil {
		if _, ok := err.(*evidence.DuplicateViewChangeEvidence); ok {
			cbft.log.Warn("Receive DuplicateViewChangeEvidence msg", "err", err.Error())
			return err
		}
	}

	cbft.state.AddViewChange(uint32(node.Index), msg)
	cbft.log.Debug("Receive new viewchange", "index", node.Index, "total", cbft.state.ViewChangeLen())
	// It is possible to achieve viewchangeQC every time you add viewchange
	cbft.tryChangeView()
	return nil
}

// OnViewTimeout performs timeout logic for view.
func (cbft *Cbft) OnViewTimeout() {
	cbft.log.Info("Current view timeout", "view", cbft.state.ViewString())
	node, err := cbft.isCurrentValidator()
	if err != nil {
		cbft.log.Error("ViewTimeout local node is not validator")
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
		ValidatorIndex: uint32(node.Index),
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

	cbft.state.AddViewChange(uint32(node.Index), viewChange)
	cbft.log.Debug("Local add viewchange", "index", node.Index, "total", cbft.state.ViewChangeLen())

	cbft.network.Broadcast(viewChange)
	cbft.tryChangeView()
}

// OnInsertQCBlock performs security rule verification, view switching.
func (cbft *Cbft) OnInsertQCBlock(blocks []*types.Block, qcs []*ctypes.QuorumCert) error {
	if len(blocks) != len(qcs) {
		return fmt.Errorf("block qc is inconsistent")
	}
	//todo insert tree, update view
	for i := 0; i < len(blocks); i++ {
		block, qc := blocks[i], qcs[i]
		//todo verify qc

		if err := cbft.safetyRules.QCBlockRules(block, qc); err != nil {
			if err.NewView() {
				cbft.changeView(qc.Epoch, qc.ViewNumber, block, qc, nil)
			} else {
				return err
			}
		}
		cbft.insertQCBlock(block, qc)
		cbft.log.Debug("Insert QC block success", "hash", qc.BlockHash, "number", qc.BlockNumber)

	}

	return nil
}

// Update blockTree, try commit new block
func (cbft *Cbft) insertQCBlock(block *types.Block, qc *ctypes.QuorumCert) {
	cbft.log.Debug("Insert QC block", "qc", qc.String())
	if block.NumberU64() <= cbft.state.HighestLockBlock().NumberU64() || cbft.HasBlock(block.Hash(), block.NumberU64()) {
		cbft.log.Debug("The inserted block has exists in chain",
			"number", block.Number(), "hash", block.Hash(),
			"lockedNumber", cbft.state.HighestLockBlock().Number(),
			"lockedHash", cbft.state.HighestLockBlock().Hash())
		return
	}
	if cbft.state.Epoch() == qc.Epoch && cbft.state.ViewNumber() == qc.ViewNumber {
		cbft.state.AddQC(qc)
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
	cbft.tryChangeView()
	if cbft.insertBlockQCHook != nil {
		// test hook
		cbft.insertBlockQCHook(block, qc)
	}
	cbft.trySendPrepareVote()
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
			if cbft.validatorPool.IsValidator(qc.Epoch, cbft.config.Option.NodeID) {
				return cbft.state.HadSendPrepareVote().Had(qc.BlockIndex) && linked(qc.BlockNumber)
			} else if cbft.validatorPool.IsCandidateNode(cbft.config.Option.NodeID) {
				blockIndex, finish := cbft.state.Executing()
				return blockIndex != math.MaxUint32 && (qc.BlockIndex < blockIndex || (qc.BlockIndex == blockIndex && finish)) && linked(qc.BlockNumber)
			}
			return false
		}

		if block != nil && hasExecuted() {
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
				if err := cbft.signBlock(block.Hash(), block.NumberU64(), index); err != nil {
					cbft.log.Error("Sign block failed", "err", err, "hash", s.Hash, "number", s.Number)
					return
				}

				cbft.log.Debug("Sign block", "hash", s.Hash, "number", s.Number)
				if msg := cbft.csPool.GetPrepareQC(index); msg != nil {
					go cbft.ReceiveMessage(msg)
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
	// todo sign vote
	// parentQC added when sending
	// Determine if the current consensus node is
	node, err := cbft.validatorPool.GetValidatorByNodeID(cbft.state.Epoch(), cbft.config.Option.NodeID)
	if err != nil {
		return err
	}
	prepareVote := &protocols.PrepareVote{
		Epoch:          cbft.state.Epoch(),
		ViewNumber:     cbft.state.ViewNumber(),
		BlockHash:      hash,
		BlockNumber:    number,
		BlockIndex:     index,
		ValidatorIndex: uint32(node.Index),
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
			cbft.log.Crit("Try send PrepareVote failed", "err", "vote corresponding block not found", "view", cbft.state.ViewString(), p.String())
		}
		if b, qc := cbft.blockTree.FindBlockAndQC(block.ParentHash(), block.NumberU64()-1); b != nil || block.NumberU64() == 0 {
			p.ParentQC = qc
			hadSend.Push(p)
			//Determine if the current consensus node is
			node, _ := cbft.validatorPool.GetValidatorByNodeID(cbft.state.Epoch(), cbft.config.Option.NodeID)
			cbft.log.Debug("Add local prepareVote", "vote", p.String())
			cbft.state.AddPrepareVote(uint32(node.Index), p)
			pending.Pop()

			// write sendPrepareVote info to wal
			if !cbft.isLoading() {
				cbft.bridge.SendPrepareVote(block, p)
			}

			cbft.network.Broadcast(p)
		} else {
			break
		}
	}
}

// Every time there is a new block or a new executed block result will enter this judgment, find the next executable block
func (cbft *Cbft) findExecutableBlock() {
	blockIndex, finish := cbft.state.Executing()
	if blockIndex == math.MaxUint32 {
		block := cbft.state.ViewBlockByIndex(blockIndex + 1)
		if block != nil {
			parent, _ := cbft.blockTree.FindBlockAndQC(block.ParentHash(), block.NumberU64()-1)
			if parent == nil {
				cbft.log.Error(fmt.Sprintf("Find executable block's parent failed :[%d,%d,%s]", blockIndex, block.NumberU64(), block.Hash()))
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

			if err := cbft.asyncExecutor.Execute(block, parent); err != nil {
				cbft.log.Error("Async Execute block failed", "error", err)
			}
			cbft.state.SetExecuting(blockIndex+1, false)
		}
	}
}

// Each time a new vote is triggered, a new QC Block will be triggered, and a new one can be found by the commit block.
func (cbft *Cbft) findQCBlock() {
	index := cbft.state.MaxQCIndex()
	next := index + 1
	size := cbft.state.PrepareVoteLenByIndex(next)

	prepareQC := func() bool {
		return size >= cbft.threshold(cbft.currentValidatorLen()) && cbft.state.HadSendPrepareVote().Had(next)
	}

	if prepareQC() {
		block := cbft.state.ViewBlockByIndex(next)
		qc := cbft.generatePrepareQC(cbft.state.AllPrepareVoteByIndex(next))
		if qc != nil {
			cbft.insertQCBlock(block, qc)
			cbft.network.Broadcast(&protocols.BlockQuorumCert{BlockQC: qc})
			// metrics
			blockQCCollectedTimer.UpdateSince(time.Unix(block.Time().Int64(), 0))
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
		cbft.log.Debug("fork block", "number", highestqc.NumberU64(), "hash", highestqc.Hash())
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
			//(qc != nil && qc.Epoch == cbft.state.Epoch() && qc.BlockIndex+1 == cbft.config.Sys.Amount && cbft.validatorPool.ShouldSwitch(block.NumberU64()))
			(qc != nil && qc.Epoch == cbft.state.Epoch() && shouldSwitch)
	}()

	if shouldSwitch {
		if err := cbft.validatorPool.Update(block.NumberU64(), cbft.state.Epoch()+1, cbft.eventMux); err == nil {
			cbft.log.Debug("Update validator success", "number", block.NumberU64())
			//if !enough {
			//	cbft.OnViewTimeout()
			//}
		}
	}

	if enough {
		cbft.log.Debug("Produce enough blocks, change view", "view", cbft.state.ViewString())
		// If current has produce enough blocks, then change view immediately.
		// Otherwise, waiting for view's timeout.
		if cbft.validatorPool.EqualSwitchPoint(block.NumberU64()) {
			cbft.changeView(cbft.state.Epoch()+1, state.DefaultViewNumber, block, qc, nil)
		} else {
			cbft.changeView(cbft.state.Epoch(), increasing(), block, qc, nil)
		}
		return
	}

	viewChangeQC := func() bool {
		if cbft.state.ViewChangeLen() >= cbft.threshold(cbft.currentValidatorLen()) {
			return true
		}
		cbft.log.Debug("Try change view failed, had receive viewchange", "len", cbft.state.ViewChangeLen(), "view", cbft.state.ViewString())
		return false
	}

	if viewChangeQC() {
		viewChangeQC := cbft.generateViewChangeQC(cbft.state.AllViewChange())
		cbft.tryChangeViewByViewChange(viewChangeQC)
	}
}

func (cbft *Cbft) richViewChangeQC(viewChangeQC *ctypes.ViewChangeQC) {
	node, err := cbft.isCurrentValidator()
	if err != nil {
		cbft.log.Error("Local node is not validator")
		return
	}
	hadSend := cbft.state.ViewChangeByIndex(uint32(node.Index))
	if hadSend != nil && !viewChangeQC.ExistViewChange(hadSend.Epoch, hadSend.ViewNumber, hadSend.BlockHash) {
		cert, err := cbft.generateViewChangeQuorumCert(hadSend)
		if err != nil {
			cbft.log.Error("Generate viewChangeQuorumCert error", "err", err)
			return
		}
		cbft.log.Debug("Already send viewChange,append viewChangeQuorumCert to ViewChangeQC", "cert", cert.String())
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
			cbft.log.Debug("Not send viewChange,append viewChangeQuorumCert to ViewChangeQC", "cert", cert.String())
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
	total := uint32(cbft.validatorPool.Len(cbft.state.Epoch()))
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
		ValidatorIndex: uint32(node.Index),
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
		return viewNumber - minuend
	}
	// syncingCache is belong to last view request, clear all sync cache
	cbft.syncingCache.Purge()
	cbft.csPool.Purge()

	cbft.state.ResetView(epoch, viewNumber)
	cbft.state.SetViewTimer(interval())
	cbft.state.SetLastViewChangeQC(viewChangeQC)

	// metrics.
	viewNumberGauage.Update(int64(viewNumber))
	epochNumberGauage.Update(int64(epoch))
	viewChangedTimer.UpdateSince(time.Unix(block.Time().Int64(), 0))

	// write confirmed viewChange info to wal
	if !cbft.isLoading() {
		cbft.bridge.ConfirmViewChange(epoch, viewNumber, block, qc, viewChangeQC)
	}
	cbft.clearInvalidBlocks(block)
	cbft.evPool.Clear(epoch, viewNumber)
	// view change maybe lags behind the other nodes,active sync prepare block
	cbft.SyncPrepareBlock("", epoch, viewNumber, 0)
	cbft.log = log.New("epoch", cbft.state.Epoch(), "view", cbft.state.ViewNumber())
	cbft.log.Debug(fmt.Sprintf("Current view deadline:%v", cbft.state.Deadline()))
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
	cbft.blockCacheWriter.ClearCache(cbft.state.HighestCommitBlock())

	//todo proposer is myself
	cbft.txPool.ForkedReset(newHead, rollback)
}
