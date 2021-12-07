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
	"encoding/json"
	"fmt"
	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/utils"
	"sync/atomic"
	"time"

	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/common/math"
	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/protocols"

	ctypes "github.com/AlayaNetwork/Alaya-Go/consensus/cbft/types"
	"github.com/AlayaNetwork/Alaya-Go/core/types"
)

const (
	DefaultEpoch      = 1
	DefaultViewNumber = 0
	// maxSelectedRGLimit is the target maximum size of the reservation of the group aggregate signature.
	maxSelectedRGLimit = 5
)

type PrepareVoteQueue struct {
	Votes []*protocols.PrepareVote `json:"votes"`
}

func newPrepareVoteQueue() *PrepareVoteQueue {
	return &PrepareVoteQueue{
		Votes: make([]*protocols.PrepareVote, 0),
	}
}

func (p *PrepareVoteQueue) Top() *protocols.PrepareVote {
	return p.Votes[0]
}

func (p *PrepareVoteQueue) Pop() *protocols.PrepareVote {
	v := p.Votes[0]
	p.Votes = p.Votes[1:]
	return v
}

func (p *PrepareVoteQueue) Push(vote *protocols.PrepareVote) {
	p.Votes = append(p.Votes, vote)
}

func (p *PrepareVoteQueue) Peek() []*protocols.PrepareVote {
	return p.Votes
}

func (p *PrepareVoteQueue) Empty() bool {
	return len(p.Votes) == 0
}

func (p *PrepareVoteQueue) Len() int {
	return len(p.Votes)
}

func (p *PrepareVoteQueue) reset() {
	p.Votes = make([]*protocols.PrepareVote, 0)
}

func (p *PrepareVoteQueue) Had(index uint32) bool {
	for _, p := range p.Votes {
		if p.BlockIndex == index {
			return true
		}
	}
	return false
}

type prepareVotes struct {
	Votes map[uint32]*protocols.PrepareVote `json:"votes"`
}

func newPrepareVotes() *prepareVotes {
	return &prepareVotes{
		Votes: make(map[uint32]*protocols.PrepareVote),
	}
}

func (p *prepareVotes) hadVote(vote *protocols.PrepareVote) bool {
	for _, v := range p.Votes {
		if v.MsgHash() == vote.MsgHash() {
			return true
		}
	}
	return false
}

func (p *prepareVotes) len() int {
	return len(p.Votes)
}

func (p *prepareVotes) clear() {
	p.Votes = make(map[uint32]*protocols.PrepareVote)
}

type viewBlocks struct {
	Blocks map[uint32]viewBlock `json:"blocks"`
}

func (v *viewBlocks) MarshalJSON() ([]byte, error) {
	type viewBlocks struct {
		Hash   common.Hash `json:"hash"`
		Number uint64      `json:"number"`
		Index  uint32      `json:"blockIndex"`
	}

	vv := make(map[uint32]viewBlocks)
	for index, block := range v.Blocks {
		vv[index] = viewBlocks{
			Hash:   block.hash(),
			Number: block.number(),
			Index:  block.blockIndex(),
		}
	}

	return json.Marshal(vv)
}

func (vb *viewBlocks) UnmarshalJSON(input []byte) error {
	type viewBlocks struct {
		Hash   common.Hash `json:"hash"`
		Number uint64      `json:"number"`
		Index  uint32      `json:"blockIndex"`
	}

	var vv map[uint32]viewBlocks
	err := json.Unmarshal(input, &vv)
	if err != nil {
		return err
	}

	vb.Blocks = make(map[uint32]viewBlock)
	for k, v := range vv {
		vb.Blocks[k] = prepareViewBlock{
			pb: &protocols.PrepareBlock{
				BlockIndex: v.Index,
				Block:      types.NewSimplifiedBlock(v.Number, v.Hash),
			},
		}
	}
	return nil
}

func newViewBlocks() *viewBlocks {
	return &viewBlocks{
		Blocks: make(map[uint32]viewBlock),
	}
}

func (v *viewBlocks) index(i uint32) viewBlock {
	return v.Blocks[i]
}

func (v *viewBlocks) addBlock(block viewBlock) {
	v.Blocks[block.blockIndex()] = block
}

func (v *viewBlocks) clear() {
	v.Blocks = make(map[uint32]viewBlock)
}

func (v *viewBlocks) len() int {
	return len(v.Blocks)
}

func (v *viewBlocks) MaxIndex() uint32 {
	max := uint32(math.MaxUint32)
	for _, b := range v.Blocks {
		if max == math.MaxUint32 || b.blockIndex() > max {
			max = b.blockIndex()
		}
	}
	return max
}

type viewQCs struct {
	MaxIndex uint32                        `json:"maxIndex"`
	QCs      map[uint32]*ctypes.QuorumCert `json:"qcs"`
}

func newViewQCs() *viewQCs {
	return &viewQCs{
		MaxIndex: math.MaxUint32,
		QCs:      make(map[uint32]*ctypes.QuorumCert),
	}
}

func (v *viewQCs) index(i uint32) *ctypes.QuorumCert {
	return v.QCs[i]
}

func (v *viewQCs) addQC(qc *ctypes.QuorumCert) {
	v.QCs[qc.BlockIndex] = qc
	if v.MaxIndex == math.MaxUint32 {
		v.MaxIndex = qc.BlockIndex
	}
	if v.MaxIndex < qc.BlockIndex {
		v.MaxIndex = qc.BlockIndex
	}
}

func (v *viewQCs) maxQCIndex() uint32 {
	return v.MaxIndex
}

func (v *viewQCs) clear() {
	v.QCs = make(map[uint32]*ctypes.QuorumCert)
	v.MaxIndex = math.MaxUint32
}

func (v *viewQCs) len() int {
	return len(v.QCs)
}

type viewVotes struct {
	Votes map[uint32]*prepareVotes `json:"votes"`
}

func newViewVotes() *viewVotes {
	return &viewVotes{
		Votes: make(map[uint32]*prepareVotes),
	}
}

func (v *viewVotes) addVote(id uint32, vote *protocols.PrepareVote) {
	if ps, ok := v.Votes[vote.BlockIndex]; ok {
		ps.Votes[id] = vote
	} else {
		ps := newPrepareVotes()
		ps.Votes[id] = vote
		v.Votes[vote.BlockIndex] = ps
	}
}

func (v *viewVotes) index(i uint32) *prepareVotes {
	return v.Votes[i]
}

func (v *viewVotes) MaxIndex() uint32 {
	max := uint32(math.MaxUint32)
	for index, _ := range v.Votes {
		if max == math.MaxUint32 || index > max {
			max = index
		}
	}
	return max
}

func (v *viewVotes) clear() {
	v.Votes = make(map[uint32]*prepareVotes)
}

type viewChanges struct {
	ViewChanges map[uint32]*protocols.ViewChange `json:"viewchanges"`
}

func newViewChanges() *viewChanges {
	return &viewChanges{
		ViewChanges: make(map[uint32]*protocols.ViewChange),
	}
}

func (v *viewChanges) addViewChange(id uint32, viewChange *protocols.ViewChange) {
	v.ViewChanges[id] = viewChange
}

func (v *viewChanges) len() int {
	return len(v.ViewChanges)
}

func (v *viewChanges) clear() {
	v.ViewChanges = make(map[uint32]*protocols.ViewChange)
}

type executing struct {
	// Block index of current view
	BlockIndex uint32 `json:"blockIndex"`
	// Whether to complete
	Finish bool `json:"finish"`
}

type view struct {
	epoch      uint64
	viewNumber uint64

	// The status of the block is currently being executed,
	// Finish indicates whether the execution is complete,
	// and the next block can be executed asynchronously after the execution is completed.
	executing executing

	// viewchange received by the current view
	viewChanges *viewChanges

	// QC of the previous view
	lastViewChangeQC *ctypes.ViewChangeQC

	// This view has been sent to other verifiers for voting
	hadSendPrepareVote *PrepareVoteQueue

	// Pending Votes of current view, parent block need receive N-f prepareVotes
	pendingVote *PrepareVoteQueue

	// Current view of the proposed block by the proposer
	viewBlocks *viewBlocks

	viewQCs *viewQCs

	// The current view generated by the vote
	viewVotes *viewVotes

	// All RGBlockQuorumCerts in the current view
	viewRGBlockQuorumCerts *viewRGBlockQuorumCerts

	// The RGBlockQuorumCerts reserved in the current view
	// used to aggregate into a complete signature and later synchronization logic
	selectedRGBlockQuorumCerts *selectedRGBlockQuorumCerts

	// This view has been sent to other verifiers for RGBlockQuorumCert
	hadSendRGBlockQuorumCerts map[uint32]struct{}

	// All RGViewChangeQuorumCerts in the current view
	viewRGViewChangeQuorumCerts *viewRGViewChangeQuorumCerts

	// The RGViewChangeQuorumCerts reserved in the current view
	// used to aggregate into a complete signature and later synchronization logic
	selectedRGViewChangeQuorumCerts *selectedRGViewChangeQuorumCerts

	// This view has been sent to other verifiers for RGViewChangeQuorumCerts
	hadSendRGViewChangeQuorumCerts map[uint64]struct{}
}

func newView() *view {
	return &view{
		executing:                       executing{math.MaxUint32, false},
		viewChanges:                     newViewChanges(),
		hadSendPrepareVote:              newPrepareVoteQueue(),
		pendingVote:                     newPrepareVoteQueue(),
		viewBlocks:                      newViewBlocks(),
		viewQCs:                         newViewQCs(),
		viewVotes:                       newViewVotes(),
		viewRGBlockQuorumCerts:          newViewRGBlockQuorumCerts(),
		selectedRGBlockQuorumCerts:      newSelectedRGBlockQuorumCerts(),
		hadSendRGBlockQuorumCerts:       make(map[uint32]struct{}),
		hadSendRGViewChangeQuorumCerts:  make(map[uint64]struct{}),
		viewRGViewChangeQuorumCerts:     newViewRGViewChangeQuorumCerts(),
		selectedRGViewChangeQuorumCerts: newSelectedRGViewChangeQuorumCerts(),
	}
}

func (v *view) Reset() {
	atomic.StoreUint64(&v.epoch, 0)
	atomic.StoreUint64(&v.viewNumber, 0)
	v.executing.BlockIndex = math.MaxUint32
	v.executing.Finish = false
	v.viewChanges.clear()
	v.hadSendPrepareVote.reset()
	v.pendingVote.reset()
	v.viewBlocks.clear()
	v.viewQCs.clear()
	v.viewVotes.clear()
	v.viewRGBlockQuorumCerts.clear()
	v.selectedRGBlockQuorumCerts.clear()
	v.hadSendRGBlockQuorumCerts = make(map[uint32]struct{})
	v.hadSendRGViewChangeQuorumCerts = make(map[uint64]struct{})
	v.viewRGViewChangeQuorumCerts.clear()
	v.selectedRGViewChangeQuorumCerts.clear()
}

func (v *view) ViewNumber() uint64 {
	return atomic.LoadUint64(&v.viewNumber)
}

func (v *view) Epoch() uint64 {
	return atomic.LoadUint64(&v.epoch)
}

func (v *view) MarshalJSON() ([]byte, error) {
	type view struct {
		Epoch                           uint64                           `json:"epoch"`
		ViewNumber                      uint64                           `json:"viewNumber"`
		Executing                       executing                        `json:"executing"`
		ViewChanges                     *viewChanges                     `json:"viewchange"`
		LastViewChangeQC                *ctypes.ViewChangeQC             `json:"lastViewchange"`
		HadSendPrepareVote              *PrepareVoteQueue                `json:"hadSendPrepareVote"`
		PendingVote                     *PrepareVoteQueue                `json:"pendingPrepareVote"`
		ViewBlocks                      *viewBlocks                      `json:"viewBlocks"`
		ViewQCs                         *viewQCs                         `json:"viewQcs"`
		ViewVotes                       *viewVotes                       `json:"viewVotes"`
		ViewRGBlockQuorumCerts          *viewRGBlockQuorumCerts          `json:"viewRGBlockQuorumCerts"`
		SelectedRGBlockQuorumCerts      *selectedRGBlockQuorumCerts      `json:"selectedRGBlockQuorumCerts"`
		HadSendRGBlockQuorumCerts       map[uint32]struct{}              `json:"hadSendRGBlockQuorumCerts"`
		HadSendRGViewChangeQuorumCerts  map[uint64]struct{}              `json:"hadSendRGViewChangeQuorumCerts"`
		ViewRGViewChangeQuorumCerts     *viewRGViewChangeQuorumCerts     `json:"viewRGViewChangeQuorumCerts"`
		SelectedRGViewChangeQuorumCerts *selectedRGViewChangeQuorumCerts `json:"selectedRGViewChangeQuorumCerts"`
	}
	vv := &view{
		Epoch:                           atomic.LoadUint64(&v.epoch),
		ViewNumber:                      atomic.LoadUint64(&v.viewNumber),
		Executing:                       v.executing,
		ViewChanges:                     v.viewChanges,
		LastViewChangeQC:                v.lastViewChangeQC,
		HadSendPrepareVote:              v.hadSendPrepareVote,
		PendingVote:                     v.pendingVote,
		ViewBlocks:                      v.viewBlocks,
		ViewQCs:                         v.viewQCs,
		ViewVotes:                       v.viewVotes,
		ViewRGBlockQuorumCerts:          v.viewRGBlockQuorumCerts,
		SelectedRGBlockQuorumCerts:      v.selectedRGBlockQuorumCerts,
		HadSendRGBlockQuorumCerts:       v.hadSendRGBlockQuorumCerts,
		HadSendRGViewChangeQuorumCerts:  v.hadSendRGViewChangeQuorumCerts,
		ViewRGViewChangeQuorumCerts:     v.viewRGViewChangeQuorumCerts,
		SelectedRGViewChangeQuorumCerts: v.selectedRGViewChangeQuorumCerts,
	}

	return json.Marshal(vv)
}

func (v *view) UnmarshalJSON(input []byte) error {
	type view struct {
		Epoch                           uint64                           `json:"epoch"`
		ViewNumber                      uint64                           `json:"viewNumber"`
		Executing                       executing                        `json:"executing"`
		ViewChanges                     *viewChanges                     `json:"viewchange"`
		LastViewChangeQC                *ctypes.ViewChangeQC             `json:"lastViewchange"`
		HadSendPrepareVote              *PrepareVoteQueue                `json:"hadSendPrepareVote"`
		PendingVote                     *PrepareVoteQueue                `json:"pendingPrepareVote"`
		ViewBlocks                      *viewBlocks                      `json:"viewBlocks"`
		ViewQCs                         *viewQCs                         `json:"viewQcs"`
		ViewVotes                       *viewVotes                       `json:"viewVotes"`
		ViewRGBlockQuorumCerts          *viewRGBlockQuorumCerts          `json:"viewRGBlockQuorumCerts"`
		SelectedRGBlockQuorumCerts      *selectedRGBlockQuorumCerts      `json:"selectedRGBlockQuorumCerts"`
		HadSendRGBlockQuorumCerts       map[uint32]struct{}              `json:"hadSendRGBlockQuorumCerts"`
		HadSendRGViewChangeQuorumCerts  map[uint64]struct{}              `json:"hadSendRGViewChangeQuorumCerts"`
		ViewRGViewChangeQuorumCerts     *viewRGViewChangeQuorumCerts     `json:"viewRGViewChangeQuorumCerts"`
		SelectedRGViewChangeQuorumCerts *selectedRGViewChangeQuorumCerts `json:"selectedRGViewChangeQuorumCerts"`
	}

	var vv view
	err := json.Unmarshal(input, &vv)
	if err != nil {
		return err
	}

	v.epoch = vv.Epoch
	v.viewNumber = vv.ViewNumber
	v.executing = vv.Executing
	v.viewChanges = vv.ViewChanges
	v.lastViewChangeQC = vv.LastViewChangeQC
	v.hadSendPrepareVote = vv.HadSendPrepareVote
	v.pendingVote = vv.PendingVote
	v.viewBlocks = vv.ViewBlocks
	v.viewQCs = vv.ViewQCs
	v.viewVotes = vv.ViewVotes
	v.viewRGBlockQuorumCerts = vv.ViewRGBlockQuorumCerts
	v.selectedRGBlockQuorumCerts = vv.SelectedRGBlockQuorumCerts
	v.hadSendRGBlockQuorumCerts = vv.HadSendRGBlockQuorumCerts
	v.hadSendRGViewChangeQuorumCerts = vv.HadSendRGViewChangeQuorumCerts
	v.viewRGViewChangeQuorumCerts = vv.ViewRGViewChangeQuorumCerts
	v.selectedRGViewChangeQuorumCerts = vv.SelectedRGViewChangeQuorumCerts
	return nil
}

//func (v *view) HadSendPrepareVote(vote *protocols.PrepareVote) bool {
//	return v.hadSendPrepareVote.hadVote(vote)
//}

//The block of current view, there two types, prepareBlock and block
type viewBlock interface {
	hash() common.Hash
	number() uint64
	blockIndex() uint32
	block() *types.Block
	//If prepareBlock is an implementation of viewBlock, return prepareBlock, otherwise nil
	prepareBlock() *protocols.PrepareBlock
}

type prepareViewBlock struct {
	pb *protocols.PrepareBlock
}

func (p prepareViewBlock) hash() common.Hash {
	return p.pb.Block.Hash()
}

func (p prepareViewBlock) number() uint64 {
	return p.pb.Block.NumberU64()
}

func (p prepareViewBlock) blockIndex() uint32 {
	return p.pb.BlockIndex
}

func (p prepareViewBlock) block() *types.Block {
	return p.pb.Block
}

func (p prepareViewBlock) prepareBlock() *protocols.PrepareBlock {
	return p.pb
}

type qcBlock struct {
	b  *types.Block
	qc *ctypes.QuorumCert
}

func (q qcBlock) hash() common.Hash {
	return q.b.Hash()
}

func (q qcBlock) number() uint64 {
	return q.b.NumberU64()
}
func (q qcBlock) blockIndex() uint32 {
	if q.qc == nil {
		return 0
	}
	return q.qc.BlockIndex
}

func (q qcBlock) block() *types.Block {
	return q.b
}

func (q qcBlock) prepareBlock() *protocols.PrepareBlock {
	return nil
}

type ViewState struct {

	//Include ViewNumber, ViewChanges, prepareVote , proposal block of current view
	*view

	highestQCBlock     atomic.Value
	highestLockBlock   atomic.Value
	highestCommitBlock atomic.Value

	//Set the timer of the view time window
	viewTimer *viewTimer

	blockTree *ctypes.BlockTree
}

func NewViewState(period uint64, blockTree *ctypes.BlockTree) *ViewState {
	return &ViewState{
		view:      newView(),
		viewTimer: newViewTimer(period),
		blockTree: blockTree,
	}
}

func (vs *ViewState) ResetView(epoch uint64, viewNumber uint64) {
	vs.view.Reset()
	atomic.StoreUint64(&vs.view.epoch, epoch)
	atomic.StoreUint64(&vs.view.viewNumber, viewNumber)
}

func (vs *ViewState) Epoch() uint64 {
	return vs.view.Epoch()
}

func (vs *ViewState) ViewNumber() uint64 {
	return vs.view.ViewNumber()
}

func (vs *ViewState) ViewString() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d}", atomic.LoadUint64(&vs.view.epoch), atomic.LoadUint64(&vs.view.viewNumber))
}

func (vs *ViewState) Deadline() time.Time {
	return vs.viewTimer.deadline
}

func (vs *ViewState) NextViewBlockIndex() uint32 {
	return vs.viewBlocks.MaxIndex() + 1
}

func (vs *ViewState) MaxViewBlockIndex() uint32 {
	max := vs.viewBlocks.MaxIndex()
	if max == math.MaxUint32 {
		return 0
	}
	return max
}

func (vs *ViewState) MaxQCIndex() uint32 {
	return vs.view.viewQCs.maxQCIndex()
}

func (vs *ViewState) ViewVoteSize() int {
	return len(vs.viewVotes.Votes)
}

func (vs *ViewState) MaxViewVoteIndex() uint32 {
	max := vs.viewVotes.MaxIndex()
	if max == math.MaxUint32 {
		return 0
	}
	return max
}

func (vs *ViewState) PrepareVoteLenByIndex(index uint32) int {
	ps := vs.viewVotes.index(index)
	if ps != nil {
		return ps.len()
	}
	return 0
}

// Find the block corresponding to the current view according to the index
func (vs *ViewState) ViewBlockByIndex(index uint32) *types.Block {
	if b := vs.view.viewBlocks.index(index); b != nil {
		return b.block()
	}
	return nil
}

func (vs *ViewState) PrepareBlockByIndex(index uint32) *protocols.PrepareBlock {
	if b := vs.view.viewBlocks.index(index); b != nil {
		return b.prepareBlock()
	}
	return nil
}

func (vs *ViewState) ViewBlockSize() int {
	return len(vs.viewBlocks.Blocks)
}

func (vs *ViewState) HadSendPrepareVote() *PrepareVoteQueue {
	return vs.view.hadSendPrepareVote
}

func (vs *ViewState) PendingPrepareVote() *PrepareVoteQueue {
	return vs.view.pendingVote
}

func (vs *ViewState) AllPrepareVoteByIndex(index uint32) map[uint32]*protocols.PrepareVote {
	ps := vs.viewVotes.index(index)
	if ps != nil {
		return ps.Votes
	}
	return nil
}

func (vs *ViewState) FindPrepareVote(blockIndex, validatorIndex uint32) *protocols.PrepareVote {
	ps := vs.viewVotes.index(blockIndex)
	if ps != nil {
		if v, ok := ps.Votes[validatorIndex]; ok {
			return v
		}
	}
	return nil
}

func (vs *ViewState) AllViewChange() map[uint32]*protocols.ViewChange {
	return vs.viewChanges.ViewChanges
}

// Returns the block index being executed, has it been completed
func (vs *ViewState) Executing() (uint32, bool) {
	return vs.view.executing.BlockIndex, vs.view.executing.Finish
}

func (vs *ViewState) SetLastViewChangeQC(qc *ctypes.ViewChangeQC) {
	vs.view.lastViewChangeQC = qc
}

func (vs *ViewState) LastViewChangeQC() *ctypes.ViewChangeQC {
	return vs.view.lastViewChangeQC
}

// Set Executing block status
func (vs *ViewState) SetExecuting(index uint32, finish bool) {
	vs.view.executing.BlockIndex, vs.view.executing.Finish = index, finish
}

func (vs *ViewState) ViewBlockAndQC(blockIndex uint32) (*types.Block, *ctypes.QuorumCert) {
	qc := vs.viewQCs.index(blockIndex)
	if b := vs.view.viewBlocks.index(blockIndex); b != nil {
		return b.block(), qc
	}
	return nil, qc
}

func (vs *ViewState) AddPrepareBlock(pb *protocols.PrepareBlock) {
	vs.view.viewBlocks.addBlock(&prepareViewBlock{pb})
}

func (vs *ViewState) AddQCBlock(block *types.Block, qc *ctypes.QuorumCert) {
	vs.view.viewBlocks.addBlock(&qcBlock{b: block, qc: qc})
}

func (vs *ViewState) AddQC(qc *ctypes.QuorumCert) {
	vs.view.viewQCs.addQC(qc)
}

func (vs *ViewState) AddPrepareVote(id uint32, vote *protocols.PrepareVote) {
	vs.view.viewVotes.addVote(id, vote)
}

func (vs *ViewState) AddViewChange(id uint32, viewChange *protocols.ViewChange) {
	vs.view.viewChanges.addViewChange(id, viewChange)
}

func (vs *ViewState) ViewChangeByIndex(index uint32) *protocols.ViewChange {
	return vs.view.viewChanges.ViewChanges[index]
}

func (vs *ViewState) ViewChangeLen() int {
	return vs.view.viewChanges.len()
}

func (vs *ViewState) HighestBlockString() string {
	qc := vs.HighestQCBlock()
	lock := vs.HighestLockBlock()
	commit := vs.HighestCommitBlock()
	return fmt.Sprintf("{HighestQC:{hash:%s,number:%d},HighestLock:{hash:%s,number:%d},HighestCommit:{hash:%s,number:%d}}",
		qc.Hash().TerminalString(), qc.NumberU64(),
		lock.Hash().TerminalString(), lock.NumberU64(),
		commit.Hash().TerminalString(), commit.NumberU64())
}

func (vs *ViewState) HighestExecutedBlock() *types.Block {
	if vs.executing.BlockIndex == math.MaxUint32 || (vs.executing.BlockIndex == 0 && !vs.executing.Finish) {
		block := vs.HighestQCBlock()
		if vs.lastViewChangeQC != nil {
			_, _, _, _, hash, _ := vs.lastViewChangeQC.MaxBlock()
			// fixme insertQCBlock should also change the state of executing
			if b := vs.blockTree.FindBlockByHash(hash); b != nil {
				block = b
			}
		}
		return block
	}

	var block *types.Block
	if vs.executing.Finish {
		block = vs.viewBlocks.index(vs.executing.BlockIndex).block()
	} else {
		block = vs.viewBlocks.index(vs.executing.BlockIndex - 1).block()
	}
	return block
}

func (vs *ViewState) FindBlock(hash common.Hash, number uint64) *types.Block {
	for _, b := range vs.viewBlocks.Blocks {
		if b.hash() == hash && b.number() == number {
			return b.block()
		}
	}
	return nil
}

func (vs *ViewState) SetHighestQCBlock(ext *types.Block) {
	vs.highestQCBlock.Store(ext)
}

func (vs *ViewState) HighestQCBlock() *types.Block {
	if v := vs.highestQCBlock.Load(); v == nil {
		panic("Get highest qc block failed")
	} else {
		return v.(*types.Block)
	}
}

func (vs *ViewState) SetHighestLockBlock(ext *types.Block) {
	vs.highestLockBlock.Store(ext)
}

func (vs *ViewState) HighestLockBlock() *types.Block {
	if v := vs.highestLockBlock.Load(); v == nil {
		panic("Get highest lock block failed")
	} else {
		return v.(*types.Block)
	}
}

func (vs *ViewState) SetHighestCommitBlock(ext *types.Block) {
	vs.highestCommitBlock.Store(ext)
}

func (vs *ViewState) HighestCommitBlock() *types.Block {
	if v := vs.highestCommitBlock.Load(); v == nil {
		panic("Get highest commit block failed")
	} else {
		return v.(*types.Block)
	}
}

func (vs *ViewState) IsDeadline() bool {
	return vs.viewTimer.isDeadline()
}

func (vs *ViewState) ViewTimeout() <-chan time.Time {
	return vs.viewTimer.timerChan()
}

func (vs *ViewState) SetViewTimer(viewInterval uint64) {
	vs.viewTimer.setupTimer(viewInterval)
}

func (vs *ViewState) String() string {
	return fmt.Sprintf("")
}

// viewRGBlockQuorumCerts
func (vs *ViewState) FindRGBlockQuorumCerts(blockIndex, groupID, validatorIndex uint32) *protocols.RGBlockQuorumCert {
	return vs.viewRGBlockQuorumCerts.FindRGBlockQuorumCerts(blockIndex, groupID, validatorIndex)
}

func (vs *ViewState) AddRGBlockQuorumCert(nodeIndex uint32, rgb *protocols.RGBlockQuorumCert) {
	vs.viewRGBlockQuorumCerts.AddRGBlockQuorumCerts(rgb.BlockIndx(), rgb)
}

func (vs *ViewState) RGBlockQuorumCertsLen(blockIndex, groupID uint32) int {
	return vs.viewRGBlockQuorumCerts.RGBlockQuorumCertsLen(blockIndex, groupID)
}

func (vs *ViewState) RGBlockQuorumCertsIndexes(blockIndex, groupID uint32) []uint32 {
	return vs.viewRGBlockQuorumCerts.RGBlockQuorumCertsIndexes(blockIndex, groupID)
}

func (vs *ViewState) FindMaxGroupRGBlockQuorumCert(blockIndex, groupID uint32) *protocols.RGBlockQuorumCert {
	return vs.viewRGBlockQuorumCerts.FindMaxGroupRGBlockQuorumCert(blockIndex, groupID)
}

// selectedRGBlockQuorumCerts
func (vs *ViewState) AddSelectRGQuorumCerts(blockIndex, groupID uint32, rgqc *ctypes.QuorumCert, parentQC *ctypes.QuorumCert) {
	vs.selectedRGBlockQuorumCerts.AddRGQuorumCerts(blockIndex, groupID, rgqc, parentQC)
}

func (vs *ViewState) SelectRGQuorumCertsLen(blockIndex, groupID uint32) int {
	return vs.selectedRGBlockQuorumCerts.RGQuorumCertsLen(blockIndex, groupID)
}

func (vs *ViewState) FindMaxRGQuorumCerts(blockIndex uint32) []*ctypes.QuorumCert {
	return vs.selectedRGBlockQuorumCerts.FindMaxRGQuorumCerts(blockIndex)
}

func (vs *ViewState) FindMaxGroupRGQuorumCert(blockIndex, groupID uint32) (*ctypes.QuorumCert, *ctypes.QuorumCert) {
	return vs.selectedRGBlockQuorumCerts.FindMaxGroupRGQuorumCert(blockIndex, groupID)
}

func (vs *ViewState) MergePrepareVotes(blockIndex, groupID uint32, votes []*protocols.PrepareVote) {
	for _, v := range votes {
		vs.selectedRGBlockQuorumCerts.MergePrepareVote(blockIndex, groupID, v)
	}
}

// viewRGViewChangeQuorumCerts
func (vs *ViewState) AddRGViewChangeQuorumCert(nodeIndex uint32, rgb *protocols.RGViewChangeQuorumCert) {
	vs.viewRGViewChangeQuorumCerts.AddRGViewChangeQuorumCerts(rgb)
}

func (vs *ViewState) FindRGViewChangeQuorumCerts(groupID uint32, validatorIndex uint32) *protocols.RGViewChangeQuorumCert {
	return vs.viewRGViewChangeQuorumCerts.FindRGViewChangeQuorumCerts(groupID, validatorIndex)
}

func (vs *ViewState) RGViewChangeQuorumCertsLen(groupID uint32) int {
	return vs.viewRGViewChangeQuorumCerts.RGViewChangeQuorumCertsLen(groupID)
}

func (vs *ViewState) RGViewChangeQuorumCertsIndexes(groupID uint32) []uint32 {
	return vs.viewRGViewChangeQuorumCerts.RGViewChangeQuorumCertsIndexes(groupID)
}

func (vs *ViewState) FindMaxRGViewChangeQuorumCert(groupID uint32) *protocols.RGViewChangeQuorumCert {
	return vs.viewRGViewChangeQuorumCerts.FindMaxRGViewChangeQuorumCert(groupID)
}

// selectedRGViewChangeQuorumCerts
func (vs *ViewState) AddSelectRGViewChangeQuorumCerts(groupID uint32, rgqc *ctypes.ViewChangeQC, prepareQCs map[common.Hash]*ctypes.QuorumCert) {
	rgqcs := make(map[common.Hash]*ctypes.ViewChangeQuorumCert)
	for _, qc := range rgqc.QCs {
		rgqcs[qc.BlockHash] = qc
	}
	vs.selectedRGViewChangeQuorumCerts.AddRGViewChangeQuorumCerts(groupID, rgqcs, prepareQCs)
}

func (vs *ViewState) MergeViewChanges(groupID uint32, vcs []*protocols.ViewChange, validatorLen int) {
	for _, vc := range vcs {
		vs.selectedRGViewChangeQuorumCerts.MergeViewChange(groupID, vc, uint32(validatorLen))
	}
}

func (vs *ViewState) SelectRGViewChangeQuorumCertsLen(groupID uint32) int {
	return vs.selectedRGViewChangeQuorumCerts.RGViewChangeQuorumCertsLen(groupID)
}

func (vs *ViewState) FindMaxGroupRGViewChangeQuorumCert(groupID uint32) (*ctypes.ViewChangeQC, *ctypes.PrepareQCs) {
	return vs.selectedRGViewChangeQuorumCerts.FindMaxGroupRGViewChangeQuorumCert(groupID)
}

func (vs *ViewState) FindMaxRGViewChangeQuorumCerts() []*ctypes.ViewChangeQC {
	return vs.selectedRGViewChangeQuorumCerts.FindMaxRGViewChangeQuorumCert()
}

// hadSendRGBlockQuorumCerts
func (vs *ViewState) AddSendRGBlockQuorumCerts(blockIndex uint32) {
	vs.hadSendRGBlockQuorumCerts[blockIndex] = struct{}{}
}

func (vs *ViewState) HadSendRGBlockQuorumCerts(blockIndex uint32) bool {
	if _, ok := vs.hadSendRGBlockQuorumCerts[blockIndex]; ok {
		return true
	}
	return false
}

// hadSendRGViewChangeQuorumCerts
func (vs *ViewState) AddSendRGViewChangeQuorumCerts(viewNumber uint64) {
	vs.hadSendRGViewChangeQuorumCerts[viewNumber] = struct{}{}
}

func (vs *ViewState) HadSendRGViewChangeQuorumCerts(viewNumber uint64) bool {
	if _, ok := vs.hadSendRGViewChangeQuorumCerts[viewNumber]; ok {
		return true
	}
	return false
}

func (vs *ViewState) MarshalJSON() ([]byte, error) {
	type hashNumber struct {
		Hash   common.Hash `json:"hash"`
		Number uint64      `json:"number"`
	}
	type state struct {
		View               *view      `json:"view"`
		HighestQCBlock     hashNumber `json:"highestQCBlock"`
		HighestLockBlock   hashNumber `json:"highestLockBlock"`
		HighestCommitBlock hashNumber `json:"highestCommitBlock"`
	}

	s := &state{
		View:               vs.view,
		HighestQCBlock:     hashNumber{Hash: vs.HighestQCBlock().Hash(), Number: vs.HighestQCBlock().NumberU64()},
		HighestLockBlock:   hashNumber{Hash: vs.HighestLockBlock().Hash(), Number: vs.HighestLockBlock().NumberU64()},
		HighestCommitBlock: hashNumber{Hash: vs.HighestCommitBlock().Hash(), Number: vs.HighestCommitBlock().NumberU64()},
	}
	return json.Marshal(s)
}

func (vs *ViewState) UnmarshalJSON(input []byte) error {
	type hashNumber struct {
		Hash   common.Hash `json:"hash"`
		Number uint64      `json:"number"`
	}
	type state struct {
		View               *view      `json:"view"`
		HighestQCBlock     hashNumber `json:"highestQCBlock"`
		HighestLockBlock   hashNumber `json:"highestLockBlock"`
		HighestCommitBlock hashNumber `json:"highestCommitBlock"`
	}

	var s state
	err := json.Unmarshal(input, &s)
	if err != nil {
		return err
	}

	vs.view = s.View
	vs.SetHighestQCBlock(types.NewSimplifiedBlock(s.HighestQCBlock.Number, s.HighestQCBlock.Hash))
	vs.SetHighestLockBlock(types.NewSimplifiedBlock(s.HighestLockBlock.Number, s.HighestLockBlock.Hash))
	vs.SetHighestCommitBlock(types.NewSimplifiedBlock(s.HighestCommitBlock.Number, s.HighestCommitBlock.Hash))
	return nil
}

// viewRGBlockQuorumCerts
type viewRGBlockQuorumCerts struct {
	BlockRGBlockQuorumCerts map[uint32]*groupRGBlockQuorumCerts `json:"blockRGBlockQuorumCerts"` // The map key is blockIndex
}

type groupRGBlockQuorumCerts struct {
	GroupRGBlockQuorumCerts map[uint32]*validatorRGBlockQuorumCerts `json:"groupRGBlockQuorumCerts"` // The map key is groupID
}

type validatorRGBlockQuorumCerts struct {
	ValidatorRGBlockQuorumCerts map[uint32]*protocols.RGBlockQuorumCert `json:"validatorRGBlockQuorumCerts"` // The map key is ValidatorIndex
}

func newViewRGBlockQuorumCerts() *viewRGBlockQuorumCerts {
	return &viewRGBlockQuorumCerts{
		BlockRGBlockQuorumCerts: make(map[uint32]*groupRGBlockQuorumCerts),
	}
}

func newGroupRGBlockQuorumCerts() *groupRGBlockQuorumCerts {
	return &groupRGBlockQuorumCerts{
		GroupRGBlockQuorumCerts: make(map[uint32]*validatorRGBlockQuorumCerts),
	}
}

func newValidatorRGBlockQuorumCerts() *validatorRGBlockQuorumCerts {
	return &validatorRGBlockQuorumCerts{
		ValidatorRGBlockQuorumCerts: make(map[uint32]*protocols.RGBlockQuorumCert),
	}
}

func (vrg *validatorRGBlockQuorumCerts) addRGBlockQuorumCerts(validatorIndex uint32, rg *protocols.RGBlockQuorumCert) bool {
	if _, ok := vrg.ValidatorRGBlockQuorumCerts[validatorIndex]; !ok {
		vrg.ValidatorRGBlockQuorumCerts[validatorIndex] = rg
		return true
	}
	return false
}

func (grg *groupRGBlockQuorumCerts) addRGBlockQuorumCerts(groupID uint32, rg *protocols.RGBlockQuorumCert) {
	if ps, ok := grg.GroupRGBlockQuorumCerts[groupID]; ok {
		ps.addRGBlockQuorumCerts(rg.ValidatorIndex, rg)
	} else {
		vrg := newValidatorRGBlockQuorumCerts()
		vrg.addRGBlockQuorumCerts(rg.ValidatorIndex, rg)
		grg.GroupRGBlockQuorumCerts[groupID] = vrg
	}
}

func (brg *viewRGBlockQuorumCerts) AddRGBlockQuorumCerts(blockIndex uint32, rg *protocols.RGBlockQuorumCert) {
	groupID := rg.GroupID
	if ps, ok := brg.BlockRGBlockQuorumCerts[blockIndex]; ok {
		ps.addRGBlockQuorumCerts(groupID, rg)
	} else {
		grg := newGroupRGBlockQuorumCerts()
		grg.addRGBlockQuorumCerts(groupID, rg)
		brg.BlockRGBlockQuorumCerts[blockIndex] = grg
	}
}

func (brg *viewRGBlockQuorumCerts) clear() {
	brg.BlockRGBlockQuorumCerts = make(map[uint32]*groupRGBlockQuorumCerts)
}

func (brg *viewRGBlockQuorumCerts) String() string {
	if s, err := json.Marshal(brg); err == nil {
		return string(s)
	}
	return ""
}

func (vrg *validatorRGBlockQuorumCerts) findRGBlockQuorumCerts(validatorIndex uint32) *protocols.RGBlockQuorumCert {
	if ps, ok := vrg.ValidatorRGBlockQuorumCerts[validatorIndex]; ok {
		return ps
	}
	return nil
}

func (grg *groupRGBlockQuorumCerts) findRGBlockQuorumCerts(groupID uint32, validatorIndex uint32) *protocols.RGBlockQuorumCert {
	if ps, ok := grg.GroupRGBlockQuorumCerts[groupID]; ok {
		return ps.findRGBlockQuorumCerts(validatorIndex)
	}
	return nil
}

func (brg *viewRGBlockQuorumCerts) FindRGBlockQuorumCerts(blockIndex, groupID, validatorIndex uint32) *protocols.RGBlockQuorumCert {
	if ps, ok := brg.BlockRGBlockQuorumCerts[blockIndex]; ok {
		return ps.findRGBlockQuorumCerts(groupID, validatorIndex)
	}
	return nil
}

func (brg *viewRGBlockQuorumCerts) FindMaxGroupRGBlockQuorumCert(blockIndex, groupID uint32) *protocols.RGBlockQuorumCert {
	if ps, ok := brg.BlockRGBlockQuorumCerts[blockIndex]; ok {
		if gs, ok := ps.GroupRGBlockQuorumCerts[groupID]; ok {
			var max *protocols.RGBlockQuorumCert
			for _, rg := range gs.ValidatorRGBlockQuorumCerts {
				if max == nil {
					max = rg
				} else if rg.BlockQC.HigherSign(max.BlockQC) {
					max = rg
				}
			}
			return max
		}
	}
	return nil
}

func (brg *viewRGBlockQuorumCerts) RGBlockQuorumCertsLen(blockIndex, groupID uint32) int {
	if ps, ok := brg.BlockRGBlockQuorumCerts[blockIndex]; ok {
		if gs, ok := ps.GroupRGBlockQuorumCerts[groupID]; ok {
			return len(gs.ValidatorRGBlockQuorumCerts)
		}
	}
	return 0
}

func (brg *viewRGBlockQuorumCerts) RGBlockQuorumCertsIndexes(blockIndex, groupID uint32) []uint32 {
	if ps, ok := brg.BlockRGBlockQuorumCerts[blockIndex]; ok {
		if gs, ok := ps.GroupRGBlockQuorumCerts[groupID]; ok {
			indexes := make([]uint32, 0, len(gs.ValidatorRGBlockQuorumCerts))
			for i, _ := range gs.ValidatorRGBlockQuorumCerts {
				indexes = append(indexes, i)
			}
			return indexes
		}
	}
	return nil
}

// selectedRGBlockQuorumCerts
type selectedRGBlockQuorumCerts struct {
	BlockRGBlockQuorumCerts map[uint32]*QuorumCerts `json:"blockRGBlockQuorumCerts"` // The map key is blockIndex
}

type QuorumCerts struct {
	GroupQuorumCerts map[uint32][]*ctypes.QuorumCert `json:"groupQuorumCerts"` // The map key is groupID
	ParentQC         *ctypes.QuorumCert              `json:"parentQC"`
}

func newSelectedRGBlockQuorumCerts() *selectedRGBlockQuorumCerts {
	return &selectedRGBlockQuorumCerts{
		BlockRGBlockQuorumCerts: make(map[uint32]*QuorumCerts),
	}
}

func newQuorumCerts() *QuorumCerts {
	return &QuorumCerts{
		GroupQuorumCerts: make(map[uint32][]*ctypes.QuorumCert),
	}
}

func (grg *QuorumCerts) addRGQuorumCerts(groupID uint32, rgqc *ctypes.QuorumCert) {
	if ps, ok := grg.GroupQuorumCerts[groupID]; ok {
		if len(ps) > 0 {
			for i := len(ps) - 1; i >= 0; i-- {
				if ps[i].ValidatorSet.Contains(rgqc.ValidatorSet) {
					return
				}
				if rgqc.ValidatorSet.Contains(ps[i].ValidatorSet) {
					ps = append(ps[:i], ps[i+1:]...)
					//grg.GroupQuorumCerts[groupID] = append(grg.GroupQuorumCerts[groupID][:i], grg.GroupQuorumCerts[groupID][i+1:]...)
					//return
				}
			}
		}
		if len(ps) < maxSelectedRGLimit || rgqc.HigherSign(findMaxQuorumCert(ps)) {
			ps = append(ps, rgqc)
			grg.GroupQuorumCerts[groupID] = ps
			//grg.GroupQuorumCerts[groupID] = append(grg.GroupQuorumCerts[groupID], rgqc)
		}
	} else {
		qcs := make([]*ctypes.QuorumCert, 0, maxSelectedRGLimit)
		qcs = append(qcs, rgqc)
		grg.GroupQuorumCerts[groupID] = qcs
	}
}

func (srg *selectedRGBlockQuorumCerts) AddRGQuorumCerts(blockIndex, groupID uint32, rgqc *ctypes.QuorumCert, parentQC *ctypes.QuorumCert) {
	if ps, ok := srg.BlockRGBlockQuorumCerts[blockIndex]; ok {
		ps.addRGQuorumCerts(groupID, rgqc)
		if ps.ParentQC == nil && parentQC != nil {
			ps.ParentQC = parentQC
		}
	} else {
		grg := newQuorumCerts()
		grg.addRGQuorumCerts(groupID, rgqc)
		if parentQC != nil {
			grg.ParentQC = parentQC
		}
		srg.BlockRGBlockQuorumCerts[blockIndex] = grg
	}
}

func (grg *QuorumCerts) findRGQuorumCerts(groupID uint32) []*ctypes.QuorumCert {
	if gs, ok := grg.GroupQuorumCerts[groupID]; ok {
		return gs
	}
	return nil
}

func (srg *selectedRGBlockQuorumCerts) FindRGQuorumCerts(blockIndex, groupID uint32) []*ctypes.QuorumCert {
	if ps, ok := srg.BlockRGBlockQuorumCerts[blockIndex]; ok {
		return ps.findRGQuorumCerts(groupID)
	}
	return nil
}

func (srg *selectedRGBlockQuorumCerts) RGQuorumCertsLen(blockIndex, groupID uint32) int {
	if ps, ok := srg.BlockRGBlockQuorumCerts[blockIndex]; ok {
		gs := ps.findRGQuorumCerts(groupID)
		return len(gs)
	}
	return 0
}

func findMaxQuorumCert(qcs []*ctypes.QuorumCert) *ctypes.QuorumCert {
	if len(qcs) > 0 {
		m := qcs[0]
		for i := 1; i < len(qcs); i++ {
			if qcs[i].HigherSign(m) {
				m = qcs[i]
			}
		}
		return m
	}
	return nil
}

// Returns the QuorumCert with the most signatures in each group
func (srg *selectedRGBlockQuorumCerts) FindMaxRGQuorumCerts(blockIndex uint32) []*ctypes.QuorumCert {
	if ps, ok := srg.BlockRGBlockQuorumCerts[blockIndex]; ok {
		var groupMaxs []*ctypes.QuorumCert // The QuorumCert with the largest number of signatures per group
		if len(ps.GroupQuorumCerts) > 0 {
			groupMaxs = make([]*ctypes.QuorumCert, 0, len(ps.GroupQuorumCerts))
			for _, qcs := range ps.GroupQuorumCerts {
				max := findMaxQuorumCert(qcs)
				if max != nil {
					groupMaxs = append(groupMaxs, max)
				}
			}
		}
		return groupMaxs
	}
	return nil
}

// Returns the QuorumCert with the most signatures in specified group
func (srg *selectedRGBlockQuorumCerts) FindMaxGroupRGQuorumCert(blockIndex, groupID uint32) (*ctypes.QuorumCert, *ctypes.QuorumCert) {
	gs := srg.FindRGQuorumCerts(blockIndex, groupID)
	max := findMaxQuorumCert(gs)
	if max != nil {
		parentQC := srg.BlockRGBlockQuorumCerts[blockIndex].ParentQC
		return max, parentQC
	}
	return nil, nil
}

func (srg *selectedRGBlockQuorumCerts) MergePrepareVote(blockIndex, groupID uint32, vote *protocols.PrepareVote) {
	rgqcs := srg.FindRGQuorumCerts(blockIndex, groupID)
	if len(rgqcs) <= 0 {
		return
	}

	for _, qc := range rgqcs {
		if !qc.HasSign(vote.NodeIndex()) {
			qc.AddSign(vote.Signature, vote.NodeIndex())
		}
	}
	// merge again
	deleteIndexes := make(map[int]struct{})
	for i := 0; i < len(rgqcs); i++ {
		if _, ok := deleteIndexes[i]; ok {
			continue
		}
		for j := 0; j < len(rgqcs); j++ {
			if _, ok := deleteIndexes[j]; ok || j == i {
				continue
			}
			if rgqcs[i].ValidatorSet.Contains(rgqcs[j].ValidatorSet) {
				deleteIndexes[j] = struct{}{}
			} else if rgqcs[j].ValidatorSet.Contains(rgqcs[i].ValidatorSet) {
				deleteIndexes[i] = struct{}{}
				break
			}
		}
	}
	if len(deleteIndexes) > 0 {
		merged := make([]*ctypes.QuorumCert, 0)
		for i := 0; i < len(rgqcs); i++ {
			if _, ok := deleteIndexes[i]; !ok {
				merged = append(merged, rgqcs[i])
			}
		}
		srg.BlockRGBlockQuorumCerts[blockIndex].GroupQuorumCerts[groupID] = merged
	}
}

func (srg *selectedRGBlockQuorumCerts) clear() {
	srg.BlockRGBlockQuorumCerts = make(map[uint32]*QuorumCerts)
}

func (srg *selectedRGBlockQuorumCerts) String() string {
	if s, err := json.Marshal(srg); err == nil {
		return string(s)
	}
	return ""
}

// viewRGViewChangeQuorumCerts
type viewRGViewChangeQuorumCerts struct {
	GroupRGViewChangeQuorumCerts map[uint32]*validatorRGViewChangeQuorumCerts `json:"groupRGViewChangeQuorumCerts"` // The map key is groupID
}

type validatorRGViewChangeQuorumCerts struct {
	ValidatorRGViewChangeQuorumCerts map[uint32]*protocols.RGViewChangeQuorumCert `json:"validatorRGViewChangeQuorumCerts"` // The map key is ValidatorIndex
}

func newViewRGViewChangeQuorumCerts() *viewRGViewChangeQuorumCerts {
	return &viewRGViewChangeQuorumCerts{
		GroupRGViewChangeQuorumCerts: make(map[uint32]*validatorRGViewChangeQuorumCerts),
	}
}

func newValidatorRGViewChangeQuorumCerts() *validatorRGViewChangeQuorumCerts {
	return &validatorRGViewChangeQuorumCerts{
		ValidatorRGViewChangeQuorumCerts: make(map[uint32]*protocols.RGViewChangeQuorumCert),
	}
}

func (vrg *validatorRGViewChangeQuorumCerts) addRGViewChangeQuorumCerts(validatorIndex uint32, rg *protocols.RGViewChangeQuorumCert) bool {
	if _, ok := vrg.ValidatorRGViewChangeQuorumCerts[validatorIndex]; !ok {
		vrg.ValidatorRGViewChangeQuorumCerts[validatorIndex] = rg
		return true
	}
	return false
}

func (brg *viewRGViewChangeQuorumCerts) AddRGViewChangeQuorumCerts(rg *protocols.RGViewChangeQuorumCert) {
	groupID := rg.GroupID
	validatorIndex := rg.ValidatorIndex
	if ps, ok := brg.GroupRGViewChangeQuorumCerts[groupID]; ok {
		ps.addRGViewChangeQuorumCerts(validatorIndex, rg)
	} else {
		vrg := newValidatorRGViewChangeQuorumCerts()
		vrg.addRGViewChangeQuorumCerts(validatorIndex, rg)
		brg.GroupRGViewChangeQuorumCerts[groupID] = vrg
	}
}

func (vrg *validatorRGViewChangeQuorumCerts) findRGViewChangeQuorumCerts(validatorIndex uint32) *protocols.RGViewChangeQuorumCert {
	if ps, ok := vrg.ValidatorRGViewChangeQuorumCerts[validatorIndex]; ok {
		return ps
	}
	return nil
}

func (brg *viewRGViewChangeQuorumCerts) FindRGViewChangeQuorumCerts(groupID uint32, validatorIndex uint32) *protocols.RGViewChangeQuorumCert {
	if ps, ok := brg.GroupRGViewChangeQuorumCerts[groupID]; ok {
		return ps.findRGViewChangeQuorumCerts(validatorIndex)
	}
	return nil
}

func (brg *viewRGViewChangeQuorumCerts) RGViewChangeQuorumCertsLen(groupID uint32) int {
	if ps, ok := brg.GroupRGViewChangeQuorumCerts[groupID]; ok {
		return len(ps.ValidatorRGViewChangeQuorumCerts)
	}
	return 0
}

func (brg *viewRGViewChangeQuorumCerts) RGViewChangeQuorumCertsIndexes(groupID uint32) []uint32 {
	if ps, ok := brg.GroupRGViewChangeQuorumCerts[groupID]; ok {
		indexes := make([]uint32, 0, len(ps.ValidatorRGViewChangeQuorumCerts))
		for i, _ := range ps.ValidatorRGViewChangeQuorumCerts {
			indexes = append(indexes, i)
		}
		return indexes
	}
	return nil
}

func (brg *viewRGViewChangeQuorumCerts) FindMaxRGViewChangeQuorumCert(groupID uint32) *protocols.RGViewChangeQuorumCert {
	if ps, ok := brg.GroupRGViewChangeQuorumCerts[groupID]; ok {
		var max *protocols.RGViewChangeQuorumCert
		for _, rg := range ps.ValidatorRGViewChangeQuorumCerts {
			if max == nil {
				max = rg
			} else if rg.ViewChangeQC.HigherSign(max.ViewChangeQC) {
				max = rg
			}
		}
		return max
	}
	return nil
}

func (brg *viewRGViewChangeQuorumCerts) clear() {
	brg.GroupRGViewChangeQuorumCerts = make(map[uint32]*validatorRGViewChangeQuorumCerts)
}

func (brg *viewRGViewChangeQuorumCerts) String() string {
	if s, err := json.Marshal(brg); err == nil {
		return string(s)
	}
	return ""
}

// selectedRGViewChangeQuorumCerts
type selectedRGViewChangeQuorumCerts struct {
	GroupRGViewChangeQuorumCerts map[uint32]*ViewChangeQuorumCerts  `json:"groupRGViewChangeQuorumCerts"` // The map key is groupID
	PrepareQCs                   map[common.Hash]*ctypes.QuorumCert `json:"prepareQCs"`
}

type ViewChangeQuorumCerts struct {
	QuorumCerts map[common.Hash][]*ctypes.ViewChangeQuorumCert `json:"quorumCerts"` // The map key is blockHash
}

func newSelectedRGViewChangeQuorumCerts() *selectedRGViewChangeQuorumCerts {
	return &selectedRGViewChangeQuorumCerts{
		GroupRGViewChangeQuorumCerts: make(map[uint32]*ViewChangeQuorumCerts),
		PrepareQCs:                   make(map[common.Hash]*ctypes.QuorumCert),
	}
}

func newViewChangeQuorumCerts() *ViewChangeQuorumCerts {
	return &ViewChangeQuorumCerts{
		QuorumCerts: make(map[common.Hash][]*ctypes.ViewChangeQuorumCert),
	}
}

func findMaxViewChangeQuorumCert(qcs []*ctypes.ViewChangeQuorumCert) *ctypes.ViewChangeQuorumCert {
	if len(qcs) > 0 {
		m := qcs[0]
		for i := 1; i < len(qcs); i++ {
			if qcs[i].HigherSign(m) {
				m = qcs[i]
			}
		}
		return m
	}
	return nil
}

func (grg *ViewChangeQuorumCerts) addRGViewChangeQuorumCert(hash common.Hash, rgqc *ctypes.ViewChangeQuorumCert) {
	if ps, ok := grg.QuorumCerts[hash]; ok {
		if len(ps) > 0 {
			for i := len(ps) - 1; i >= 0; i-- {
				if ps[i].ValidatorSet.Contains(rgqc.ValidatorSet) {
					return
				}
				if rgqc.ValidatorSet.Contains(ps[i].ValidatorSet) {
					ps = append(ps[:i], ps[i+1:]...)
					//grg.QuorumCerts[hash] = append(grg.QuorumCerts[hash][:i], grg.QuorumCerts[hash][i+1:]...)
					//return
				}
			}
		}
		if len(ps) < maxSelectedRGLimit || rgqc.HigherSign(findMaxViewChangeQuorumCert(ps)) {
			ps = append(ps, rgqc)
			grg.QuorumCerts[hash] = ps
			//grg.QuorumCerts[hash] = append(grg.QuorumCerts[hash], rgqc)
		}
	} else {
		qcs := make([]*ctypes.ViewChangeQuorumCert, 0, maxSelectedRGLimit)
		qcs = append(qcs, rgqc)
		grg.QuorumCerts[hash] = qcs
	}
}

func (grg *ViewChangeQuorumCerts) addRGViewChangeQuorumCerts(rgqcs map[common.Hash]*ctypes.ViewChangeQuorumCert) {
	for hash, qc := range rgqcs {
		grg.addRGViewChangeQuorumCert(hash, qc)
	}
}

func (srg *selectedRGViewChangeQuorumCerts) AddRGViewChangeQuorumCerts(groupID uint32, rgqcs map[common.Hash]*ctypes.ViewChangeQuorumCert, prepareQCs map[common.Hash]*ctypes.QuorumCert) {
	if ps, ok := srg.GroupRGViewChangeQuorumCerts[groupID]; ok {
		ps.addRGViewChangeQuorumCerts(rgqcs)
	} else {
		grg := newViewChangeQuorumCerts()
		grg.addRGViewChangeQuorumCerts(rgqcs)
		srg.GroupRGViewChangeQuorumCerts[groupID] = grg
	}
	if len(prepareQCs) > 0 {
		for hash, qc := range prepareQCs {
			if srg.PrepareQCs[hash] == nil {
				srg.PrepareQCs[hash] = qc
			}
		}
	}
}

func (srg *selectedRGViewChangeQuorumCerts) findRGQuorumCerts(groupID uint32) map[common.Hash][]*ctypes.ViewChangeQuorumCert {
	if ps, ok := srg.GroupRGViewChangeQuorumCerts[groupID]; ok {
		if ps != nil && len(ps.QuorumCerts) > 0 {
			return ps.QuorumCerts
		}
	}
	return nil
}

func (srg *selectedRGViewChangeQuorumCerts) MergeViewChange(groupID uint32, vc *protocols.ViewChange, validatorLen uint32) {
	rgqcs := srg.findRGQuorumCerts(groupID)
	if len(rgqcs) <= 0 {
		// If there is no aggregate signature under the group at this time, then the single ViewChange is not merged
		return
	}

	if qcs, ok := rgqcs[vc.BHash()]; ok {
		for _, qc := range qcs {
			if !qc.HasSign(vc.NodeIndex()) {
				qc.AddSign(vc.Signature, vc.NodeIndex())
			}
		}
		// merge again
		deleteIndexes := make(map[int]struct{})
		for i := 0; i < len(qcs); i++ {
			if _, ok := deleteIndexes[i]; ok {
				continue
			}
			for j := 0; j < len(qcs); j++ {
				if _, ok := deleteIndexes[j]; ok || j == i {
					continue
				}
				if qcs[i].ValidatorSet.Contains(qcs[j].ValidatorSet) {
					deleteIndexes[j] = struct{}{}
				} else if qcs[j].ValidatorSet.Contains(qcs[i].ValidatorSet) {
					deleteIndexes[i] = struct{}{}
					break
				}
			}
		}
		if len(deleteIndexes) > 0 {
			merged := make([]*ctypes.ViewChangeQuorumCert, 0)
			for i := 0; i < len(qcs); i++ {
				if _, ok := deleteIndexes[i]; !ok {
					merged = append(merged, qcs[i])
				}
			}
			srg.GroupRGViewChangeQuorumCerts[groupID].QuorumCerts[vc.BHash()] = merged
		}
	} else {
		qcs := make([]*ctypes.ViewChangeQuorumCert, 0, maxSelectedRGLimit)
		qc := &ctypes.ViewChangeQuorumCert{
			Epoch:        vc.Epoch,
			ViewNumber:   vc.ViewNumber,
			BlockHash:    vc.BlockHash,
			BlockNumber:  vc.BlockNumber,
			ValidatorSet: utils.NewBitArray(validatorLen),
		}
		if vc.PrepareQC != nil {
			qc.BlockEpoch = vc.PrepareQC.Epoch
			qc.BlockViewNumber = vc.PrepareQC.ViewNumber
			srg.PrepareQCs[vc.BHash()] = vc.PrepareQC
		}
		qc.Signature.SetBytes(vc.Signature.Bytes())
		qc.ValidatorSet.SetIndex(vc.ValidatorIndex, true)
		qcs = append(qcs, qc)

		srg.GroupRGViewChangeQuorumCerts[groupID].QuorumCerts[vc.BHash()] = qcs
	}
}

func (srg *selectedRGViewChangeQuorumCerts) RGViewChangeQuorumCertsLen(groupID uint32) int {
	rgqcs := srg.findRGQuorumCerts(groupID)
	return len(rgqcs)
}

// Returns the QuorumCert with the most signatures in specified group
func (srg *selectedRGViewChangeQuorumCerts) FindMaxGroupRGViewChangeQuorumCert(groupID uint32) (*ctypes.ViewChangeQC, *ctypes.PrepareQCs) {
	rgqcs := srg.findRGQuorumCerts(groupID)
	if len(rgqcs) <= 0 {
		return nil, nil
	}

	viewChangeQC := &ctypes.ViewChangeQC{QCs: make([]*ctypes.ViewChangeQuorumCert, 0)}
	prepareQCs := &ctypes.PrepareQCs{QCs: make([]*ctypes.QuorumCert, 0)}
	for hash, qcs := range rgqcs {
		max := findMaxViewChangeQuorumCert(qcs)
		viewChangeQC.QCs = append(viewChangeQC.QCs, max)
		if srg.PrepareQCs != nil && srg.PrepareQCs[hash] != nil {
			prepareQCs.QCs = append(prepareQCs.QCs, srg.PrepareQCs[hash])
		}
	}
	return viewChangeQC, prepareQCs
}

func (srg *selectedRGViewChangeQuorumCerts) FindMaxRGViewChangeQuorumCert() []*ctypes.ViewChangeQC {
	viewChangeQCs := make([]*ctypes.ViewChangeQC, 0, len(srg.GroupRGViewChangeQuorumCerts))
	for groupID, _ := range srg.GroupRGViewChangeQuorumCerts {
		viewChangeQC, _ := srg.FindMaxGroupRGViewChangeQuorumCert(groupID)
		if viewChangeQC != nil {
			viewChangeQCs = append(viewChangeQCs, viewChangeQC)
		}
	}
	return viewChangeQCs
}

func (srg *selectedRGViewChangeQuorumCerts) clear() {
	srg.GroupRGViewChangeQuorumCerts = make(map[uint32]*ViewChangeQuorumCerts)
	srg.PrepareQCs = make(map[common.Hash]*ctypes.QuorumCert)
}

func (srg *selectedRGViewChangeQuorumCerts) String() string {
	if s, err := json.Marshal(srg); err == nil {
		return string(s)
	}
	return ""
}
