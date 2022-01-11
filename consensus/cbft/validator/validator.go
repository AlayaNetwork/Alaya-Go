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

package validator

import (
	"encoding/hex"
	"errors"
	"fmt"
	"sync"

	"github.com/AlayaNetwork/Alaya-Go/x/xcom"

	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"

	"github.com/AlayaNetwork/Alaya-Go/core/state"

	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/utils"

	"github.com/AlayaNetwork/Alaya-Go/common"
	cvm "github.com/AlayaNetwork/Alaya-Go/common/vm"
	"github.com/AlayaNetwork/Alaya-Go/consensus"
	"github.com/AlayaNetwork/Alaya-Go/core"
	"github.com/AlayaNetwork/Alaya-Go/core/cbfttypes"
	"github.com/AlayaNetwork/Alaya-Go/core/types"
	"github.com/AlayaNetwork/Alaya-Go/core/vm"
	"github.com/AlayaNetwork/Alaya-Go/crypto"
	"github.com/AlayaNetwork/Alaya-Go/crypto/bls"
	"github.com/AlayaNetwork/Alaya-Go/event"
	"github.com/AlayaNetwork/Alaya-Go/log"
	"github.com/AlayaNetwork/Alaya-Go/params"
	"github.com/AlayaNetwork/Alaya-Go/rlp"
)

func newValidators(nodes []params.CbftNode, validBlockNumber uint64) *cbfttypes.Validators {
	vds := &cbfttypes.Validators{
		Nodes:            make(cbfttypes.ValidateNodeMap, len(nodes)),
		ValidBlockNumber: validBlockNumber,
	}

	for i, node := range nodes {
		pubkey := node.Node.Pubkey()
		if pubkey == nil {
			panic("pubkey should not nil")
		}

		blsPubKey := node.BlsPubKey

		vds.Nodes[node.Node.ID()] = &cbfttypes.ValidateNode{
			Index:     uint32(i),
			Address:   crypto.PubkeyToNodeAddress(*pubkey),
			PubKey:    pubkey,
			NodeID:    node.Node.ID(),
			BlsPubKey: &blsPubKey,
		}
	}
	return vds
}

type StaticAgency struct {
	consensus.Agency

	validators *cbfttypes.Validators
}

func NewStaticAgency(nodes []params.CbftNode) consensus.Agency {
	return &StaticAgency{
		validators: newValidators(nodes, 0),
	}
}

func (d *StaticAgency) Flush(header *types.Header) error {
	return nil
}

func (d *StaticAgency) Sign(interface{}) error {
	return nil
}

func (d *StaticAgency) VerifySign(interface{}) error {
	return nil
}

func (d *StaticAgency) VerifyHeader(header *types.Header, statedb *state.StateDB) error {
	return nil
}

func (d *StaticAgency) GetLastNumber(blockNumber uint64) uint64 {
	return 0
}

func (d *StaticAgency) GetValidators(blockHash common.Hash, blockNumber uint64) (*cbfttypes.Validators, error) {
	return d.validators, nil
}

func (d *StaticAgency) IsCandidateNode(nodeID enode.IDv0) bool {
	return false
}

func (d *StaticAgency) OnCommit(block *types.Block) error {
	return nil
}

type MockAgency struct {
	consensus.Agency
	validators *cbfttypes.Validators
	interval   uint64
}

func NewMockAgency(nodes []params.CbftNode, interval uint64) consensus.Agency {
	return &MockAgency{
		validators: newValidators(nodes, 0),
		interval:   interval,
	}
}

func (d *MockAgency) Flush(header *types.Header) error {
	return nil
}

func (d *MockAgency) Sign(interface{}) error {
	return nil
}

func (d *MockAgency) VerifySign(interface{}) error {
	return nil
}

func (d *MockAgency) VerifyHeader(header *types.Header, statedb *state.StateDB) error {
	return nil
}

func (d *MockAgency) GetLastNumber(blockNumber uint64) uint64 {
	if blockNumber%d.interval == 1 {
		return blockNumber + d.interval - 1
	}
	return 0
}

func (d *MockAgency) GetValidators(blockHash common.Hash, blockNumber uint64) (*cbfttypes.Validators, error) {
	if blockNumber > d.interval && blockNumber%d.interval == 1 {
		d.validators.ValidBlockNumber = d.validators.ValidBlockNumber + d.interval + 1
	}
	return d.validators, nil
}

func (d *MockAgency) IsCandidateNode(nodeID enode.IDv0) bool {
	return false
}

func (d *MockAgency) OnCommit(block *types.Block) error {
	return nil
}

type InnerAgency struct {
	consensus.Agency

	blocksPerNode         uint64
	defaultBlocksPerRound uint64
	offset                uint64
	blockchain            *core.BlockChain
	defaultValidators     *cbfttypes.Validators
}

func NewInnerAgency(nodes []params.CbftNode, chain *core.BlockChain, blocksPerNode, offset int) consensus.Agency {
	return &InnerAgency{
		blocksPerNode:         uint64(blocksPerNode),
		defaultBlocksPerRound: uint64(len(nodes) * blocksPerNode),
		offset:                uint64(offset),
		blockchain:            chain,
		defaultValidators:     newValidators(nodes, 1),
	}
}

func (ia *InnerAgency) Flush(header *types.Header) error {
	return nil
}

func (ia *InnerAgency) Sign(interface{}) error {
	return nil
}

func (ia *InnerAgency) VerifySign(interface{}) error {
	return nil
}

func (ia *InnerAgency) VerifyHeader(header *types.Header, stateDB *state.StateDB) error {
	return nil
}

func (ia *InnerAgency) GetLastNumber(blockNumber uint64) uint64 {
	var lastBlockNumber uint64
	if blockNumber <= ia.defaultBlocksPerRound {
		lastBlockNumber = ia.defaultBlocksPerRound
	} else {
		vds, err := ia.GetValidators(common.ZeroHash, blockNumber)
		if err != nil {
			log.Error("Get validator fail", "blockNumber", blockNumber)
			return 0
		}

		if vds.ValidBlockNumber == 0 && blockNumber%ia.defaultBlocksPerRound == 0 {
			return blockNumber
		}

		// lastNumber = vds.ValidBlockNumber + ia.blocksPerNode * vds.Len() - 1
		lastBlockNumber = vds.ValidBlockNumber + ia.blocksPerNode*uint64(vds.Len()) - 1

		// May be `CurrentValidators ` had not updated, so we need to calcuate `lastBlockNumber`
		// via `blockNumber`.
		if lastBlockNumber < blockNumber {
			blocksPerRound := ia.blocksPerNode * uint64(vds.Len())
			if blockNumber%blocksPerRound == 0 {
				lastBlockNumber = blockNumber
			} else {
				baseNum := blockNumber - (blockNumber % blocksPerRound)
				lastBlockNumber = baseNum + blocksPerRound
			}
		}
	}
	//log.Debug("Get last block number", "blockNumber", blockNumber, "lastBlockNumber", lastBlockNumber)
	return lastBlockNumber
}

func (ia *InnerAgency) GetValidators(blockHash common.Hash, blockNumber uint64) (v *cbfttypes.Validators, err error) {
	defaultValidators := *ia.defaultValidators
	baseNumber := blockNumber
	if blockNumber == 0 {
		baseNumber = 1
	}
	defaultValidators.ValidBlockNumber = ((baseNumber-1)/ia.defaultBlocksPerRound)*ia.defaultBlocksPerRound + 1
	if blockNumber <= ia.defaultBlocksPerRound {
		return &defaultValidators, nil
	}

	// Otherwise, get validators from inner contract.
	vdsCftNum := blockNumber - ia.offset - 1
	block := ia.blockchain.GetBlockByNumber(vdsCftNum)
	if block == nil {
		log.Error("Get the block fail, use default validators", "number", vdsCftNum)
		return &defaultValidators, nil
	}
	state, err := ia.blockchain.StateAt(block.Root())
	if err != nil {
		log.Error("Get the state fail, use default validators", "number", block.Number(), "hash", block.Hash(), "error", err)
		return &defaultValidators, nil
	}
	b := state.GetState(cvm.ValidatorInnerContractAddr, []byte(vm.CurrentValidatorKey))
	if len(b) == 0 {
		return &defaultValidators, nil
	}
	var vds vm.Validators
	err = rlp.DecodeBytes(b, &vds)
	if err != nil {
		log.Error("RLP decode fail, use default validators", "number", block.Number(), "error", err)
		return &defaultValidators, nil
	}
	var validators cbfttypes.Validators
	validators.Nodes = make(cbfttypes.ValidateNodeMap, len(vds.ValidateNodes))

	for _, node := range vds.ValidateNodes {
		pubkey, _ := node.NodeID.Pubkey()
		blsPubKey := node.BlsPubKey
		id := enode.PubkeyToIDV4(pubkey)
		validators.Nodes[id] = &cbfttypes.ValidateNode{
			Index:     uint32(node.Index),
			Address:   node.Address,
			PubKey:    pubkey,
			NodeID:    id,
			BlsPubKey: &blsPubKey,
		}
	}
	validators.ValidBlockNumber = vds.ValidBlockNumber
	return &validators, nil
}

func (ia *InnerAgency) IsCandidateNode(nodeID enode.IDv0) bool {
	return true
}

func (ia *InnerAgency) OnCommit(block *types.Block) error {
	return nil
}

// ValidatorPool a pool storing validators.
type ValidatorPool struct {
	agency consensus.Agency
	lock   sync.RWMutex

	// Current node's public key
	nodeID enode.ID

	// A block number which validators switch to current.
	switchPoint uint64
	lastNumber  uint64

	// current epoch
	epoch uint64

	// grouped indicates if validators need grouped
	grouped bool
	// max validators in per group
	groupValidatorsLimit uint32
	// coordinator limit
	coordinatorLimit uint32

	prevValidators    *cbfttypes.Validators // Previous round validators
	currentValidators *cbfttypes.Validators // Current round validators
	nextValidators    *cbfttypes.Validators // Next round validators, to Post Pub event
}

// NewValidatorPool new a validator pool.
func NewValidatorPool(agency consensus.Agency, blockNumber, epoch uint64, nodeID enode.ID, needGroup bool, eventMux *event.TypeMux) *ValidatorPool {
	pool := &ValidatorPool{
		agency:               agency,
		nodeID:               nodeID,
		epoch:                epoch,
		grouped:              needGroup,
		groupValidatorsLimit: xcom.MaxGroupValidators(),
		coordinatorLimit:     xcom.CoordinatorsLimit(),
	}
	// FIXME: Check `GetValidators` return error
	if agency.GetLastNumber(blockNumber) == blockNumber {
		pool.prevValidators, _ = agency.GetValidators(common.ZeroHash, blockNumber)
		pool.currentValidators, _ = agency.GetValidators(common.ZeroHash, NextRound(blockNumber))
		pool.lastNumber = agency.GetLastNumber(NextRound(blockNumber))
		if blockNumber != 0 {
			pool.epoch += 1
		}
	} else {
		pool.currentValidators, _ = agency.GetValidators(common.ZeroHash, blockNumber)
		pool.prevValidators = pool.currentValidators
		pool.lastNumber = agency.GetLastNumber(blockNumber)
	}
	// When validator mode is `static`, the `ValidatorBlockNumber` always 0,
	// means we are using static validators. Otherwise, represent use current
	// validators validate start from `ValidatorBlockNumber` block,
	// so `ValidatorBlockNumber` - 1 is the switch point.
	if pool.currentValidators.ValidBlockNumber > 0 {
		pool.switchPoint = pool.currentValidators.ValidBlockNumber - 1
	}
	if needGroup {
		if err := pool.organize(pool.currentValidators, epoch, eventMux); err != nil {
			log.Error("ValidatorPool organized failed!", "error", err)
		}
		if pool.nextValidators == nil {
			nds, err := pool.agency.GetValidators(common.ZeroHash, NextRound(pool.currentValidators.ValidBlockNumber))
			if err != nil {
				log.Debug("Get nextValidators error", "blockNumber", blockNumber, "err", err)
				return pool
			}
			if nds != nil && !pool.currentValidators.Equal(nds) {
				pool.nextValidators = nds
				pool.organize(pool.nextValidators, epoch+1, eventMux)
			}
		}
	}
	log.Debug("Update validator", "validators", pool.currentValidators.String(), "switchpoint", pool.switchPoint, "epoch", pool.epoch, "lastNumber", pool.lastNumber)
	return pool
}

// Reset reset validator pool.
func (vp *ValidatorPool) Reset(blockNumber uint64, epoch uint64, eventMux *event.TypeMux) {
	if vp.agency.GetLastNumber(blockNumber) == blockNumber {
		vp.prevValidators, _ = vp.agency.GetValidators(common.ZeroHash, blockNumber)
		vp.currentValidators, _ = vp.agency.GetValidators(common.ZeroHash, NextRound(blockNumber))
		vp.lastNumber = vp.agency.GetLastNumber(NextRound(blockNumber))
		vp.epoch = epoch + 1
	} else {
		vp.currentValidators, _ = vp.agency.GetValidators(common.ZeroHash, blockNumber)
		vp.prevValidators = vp.currentValidators
		vp.lastNumber = vp.agency.GetLastNumber(blockNumber)
		vp.epoch = epoch
	}
	if vp.currentValidators.ValidBlockNumber > 0 {
		vp.switchPoint = vp.currentValidators.ValidBlockNumber - 1
	}
	if vp.grouped {
		vp.organize(vp.currentValidators, epoch, eventMux)
	}
	log.Debug("Update validator", "validators", vp.currentValidators.String(), "switchpoint", vp.switchPoint, "epoch", vp.epoch, "lastNumber", vp.lastNumber)
}

// ShouldSwitch check if should switch validators at the moment.
func (vp *ValidatorPool) ShouldSwitch(blockNumber uint64) bool {
	if blockNumber == 0 {
		return false
	}
	if blockNumber == vp.switchPoint {
		return true
	}
	return blockNumber == vp.lastNumber
}

// EqualSwitchPoint returns boolean which representment the switch point
// equal the inputs number.
func (vp *ValidatorPool) EqualSwitchPoint(number uint64) bool {
	return vp.switchPoint > 0 && vp.switchPoint == number
}

func (vp *ValidatorPool) EnableVerifyEpoch(epoch uint64) error {
	if epoch+1 == vp.epoch || epoch == vp.epoch {
		return nil
	}
	return fmt.Errorf("unable verify epoch:%d,%d, request:%d", vp.epoch-1, vp.epoch, epoch)
}

func (vp *ValidatorPool) MockSwitchPoint(number uint64) {
	vp.switchPoint = 0
	vp.lastNumber = number
}

// Update switch validators.
func (vp *ValidatorPool) Update(blockHash common.Hash, blockNumber uint64, epoch uint64, version uint32, eventMux *event.TypeMux) error {
	vp.lock.Lock()
	defer vp.lock.Unlock()

	needGroup := version >= params.FORKVERSION_0_17_0
	//分组提案生效后第一个共识round到ElectionPoint时初始化分组信息
	if !vp.grouped && needGroup {
		vp.grouped = true
		vp.groupValidatorsLimit = xcom.MaxGroupValidators()
		vp.coordinatorLimit = xcom.CoordinatorsLimit()
	}
	// 生效后第一个共识周期的Election block已经是新值（2130）所以第一次触发update是cbft.tryChangeView->shouldSwitch
	if blockNumber <= vp.switchPoint {
		log.Trace("Already update validator before", "blockNumber", blockNumber, "switchPoint", vp.switchPoint)
		return errors.New("already updated before")
	}

	var err error
	var nds *cbfttypes.Validators
	if vp.nextValidators == nil {
		nds, err = vp.agency.GetValidators(blockHash, NextRound(blockNumber))
		if err != nil {
			log.Error("Get validator error", "blockNumber", blockNumber, "err", err)
			return err
		}
		if needGroup {
			//生效后第一个共识周期的switchpoint是旧值，此时不能切换
			//判断依据是新validators和current完全相同且nextValidators为空
			if vp.currentValidators.Equal(nds) {
				vp.currentValidators = nds
				vp.switchPoint = nds.ValidBlockNumber - 1
				vp.lastNumber = vp.agency.GetLastNumber(blockNumber)
				//不切换，所以epoch不增
				vp.grouped = false
				log.Debug("update currentValidators success!", "lastNumber", vp.lastNumber, "grouped", vp.grouped, "switchPoint", vp.switchPoint)
				return nil
			}
		}

		// 节点中间重启过， nextValidators没有赋值
		vp.nextValidators = nds
		if vp.grouped {
			vp.organize(vp.nextValidators, epoch, eventMux)
		}
	}

	vp.prevValidators = vp.currentValidators
	vp.currentValidators = vp.nextValidators
	vp.switchPoint = vp.currentValidators.ValidBlockNumber - 1
	vp.lastNumber = vp.agency.GetLastNumber(NextRound(blockNumber))
	currEpoch := vp.epoch
	vp.epoch = epoch
	vp.nextValidators = nil

	//切换共识轮时需要将上一轮分组的topic取消订阅
	vp.dissolve(currEpoch, eventMux)
	//旧版本（非分组共识）需要发events断开落选的共识节点
	if !vp.grouped {
		vp.dealWithOldVersionEvents(epoch, eventMux)
	}
	log.Info("Update validators", "validators.len", vp.currentValidators.Len(), "switchpoint", vp.switchPoint, "epoch", vp.epoch, "lastNumber", vp.lastNumber)
	return nil
}

// pre-init validator nodes for the next round.version >= 0.17.0 only
func (vp *ValidatorPool) InitComingValidators(blockHash common.Hash, blockNumber uint64, eventMux *event.TypeMux) error {
	vp.lock.Lock()
	defer vp.lock.Unlock()

	//分组提案生效后第一个共识round到ElectionPoint时初始化分组信息
	if !vp.grouped {
		vp.grouped = true
		vp.groupValidatorsLimit = xcom.MaxGroupValidators()
		vp.coordinatorLimit = xcom.CoordinatorsLimit()
	}

	// 提前更新nextValidators，为了p2p早一步订阅分组事件以便建链接
	nds, err := vp.agency.GetValidators(blockHash, blockNumber+xcom.ElectionDistance()+1)
	if err != nil {
		log.Error("Get validators error", "blockNumber", blockNumber, "err", err)
		return err
	}
	vp.nextValidators = nds
	vp.organize(vp.nextValidators, vp.epoch+1, eventMux)
	log.Debug("Update nextValidators OK", "blockNumber", blockNumber, "epoch", vp.epoch+1)
	return nil
}

// dealWithOldVersionEvents process version <= 0.16.0 logics
func (vp *ValidatorPool) dealWithOldVersionEvents(epoch uint64, eventMux *event.TypeMux) {
	isValidatorBefore := vp.isValidator(epoch-1, vp.nodeID)
	isValidatorAfter := vp.isValidator(epoch, vp.nodeID)

	if isValidatorBefore {
		// If we are still a consensus node, that adding
		// new validators as consensus peer, and removing
		// validators. Added as consensus peersis because
		// we need to keep connect with other validators
		// in the consensus stages. Also we are not needed
		// to keep connect with old validators.
		if isValidatorAfter {
			for _, n := range vp.currentValidators.Nodes {
				if node, _ := vp.prevValidators.FindNodeByID(n.NodeID); node == nil {
					eventMux.Post(cbfttypes.AddValidatorEvent{Node: enode.NewV4(n.PubKey, nil, 0, 0)})
					log.Trace("Post AddValidatorEvent", "node", n.String())
				}
			}

			for _, n := range vp.prevValidators.Nodes {
				if node, _ := vp.currentValidators.FindNodeByID(n.NodeID); node == nil {
					eventMux.Post(cbfttypes.RemoveValidatorEvent{Node: enode.NewV4(n.PubKey, nil, 0, 0)})
					log.Trace("Post RemoveValidatorEvent", "node", n.String())
				}
			}
		} else {
			for _, node := range vp.prevValidators.Nodes {
				eventMux.Post(cbfttypes.RemoveValidatorEvent{Node: enode.NewV4(node.PubKey, nil, 0, 0)})
				log.Trace("Post RemoveValidatorEvent", "nodeID", node.String())
			}
		}
	} else {
		// We are become a consensus node, that adding all
		// validators as consensus peer except us. Added as
		// consensus peers is because we need to keep connecting
		// with other validators in the consensus stages.
		if isValidatorAfter {
			for _, node := range vp.currentValidators.Nodes {
				eventMux.Post(cbfttypes.AddValidatorEvent{Node: enode.NewV4(node.PubKey, nil, 0, 0)})
				log.Trace("Post AddValidatorEvent", "nodeID", node.String())
			}
		}
	}
}

// GetValidatorByNodeID get the validator by node id.
func (vp *ValidatorPool) GetValidatorByNodeID(epoch uint64, nodeID enode.ID) (*cbfttypes.ValidateNode, error) {
	vp.lock.RLock()
	defer vp.lock.RUnlock()
	return vp.getValidatorByNodeID(epoch, nodeID)
}

func (vp *ValidatorPool) getValidatorByNodeID(epoch uint64, nodeID enode.ID) (*cbfttypes.ValidateNode, error) {
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators.FindNodeByID(nodeID)
	}
	return vp.currentValidators.FindNodeByID(nodeID)
}

// GetValidatorByAddr get the validator by address.
func (vp *ValidatorPool) GetValidatorByAddr(epoch uint64, addr common.NodeAddress) (*cbfttypes.ValidateNode, error) {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	return vp.getValidatorByAddr(epoch, addr)
}

func (vp *ValidatorPool) getValidatorByAddr(epoch uint64, addr common.NodeAddress) (*cbfttypes.ValidateNode, error) {
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators.FindNodeByAddress(addr)
	}
	return vp.currentValidators.FindNodeByAddress(addr)
}

// GetValidatorByIndex get the validator by index.
func (vp *ValidatorPool) GetValidatorByIndex(epoch uint64, index uint32) (*cbfttypes.ValidateNode, error) {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	return vp.getValidatorByIndex(epoch, index)
}

func (vp *ValidatorPool) getValidatorByIndex(epoch uint64, index uint32) (*cbfttypes.ValidateNode, error) {
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators.FindNodeByIndex(index)
	}
	return vp.currentValidators.FindNodeByIndex(index)
}

// GetNodeIDByIndex get the node id by index.
func (vp *ValidatorPool) GetNodeIDByIndex(epoch uint64, index uint32) enode.ID {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	return vp.getNodeIDByIndex(epoch, index)
}

func (vp *ValidatorPool) getNodeIDByIndex(epoch uint64, index uint32) enode.ID {
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators.NodeID(index)
	}
	return vp.currentValidators.NodeID(index)
}

// GetIndexByNodeID get the index by node id.
func (vp *ValidatorPool) GetIndexByNodeID(epoch uint64, nodeID enode.ID) (uint32, error) {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	return vp.getIndexByNodeID(epoch, nodeID)
}

func (vp *ValidatorPool) getIndexByNodeID(epoch uint64, nodeID enode.ID) (uint32, error) {
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators.Index(nodeID)
	}
	return vp.currentValidators.Index(nodeID)
}

// ValidatorList get the validator list.
func (vp *ValidatorPool) ValidatorList(epoch uint64) []enode.ID {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	return vp.validatorList(epoch)
}

func (vp *ValidatorPool) validatorList(epoch uint64) []enode.ID {
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators.NodeIdList()
	}
	return vp.currentValidators.NodeIdList()
}

func (vp *ValidatorPool) Validators(epoch uint64) *cbfttypes.Validators {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators
	}
	return vp.currentValidators
}

// VerifyHeader verify block's header.
func (vp *ValidatorPool) VerifyHeader(header *types.Header) error {
	_, err := crypto.Ecrecover(header.SealHash().Bytes(), header.Signature())
	if err != nil {
		return err
	}
	return vp.agency.VerifyHeader(header, nil)
}

// IsValidator check if the node is validator.
func (vp *ValidatorPool) IsValidator(epoch uint64, nodeID enode.ID) bool {
	return vp.isValidator(epoch, nodeID)
}

func (vp *ValidatorPool) isValidator(epoch uint64, nodeID enode.ID) bool {
	_, err := vp.getValidatorByNodeID(epoch, nodeID)
	return err == nil
}

// IsCandidateNode check if the node is candidate node.
func (vp *ValidatorPool) IsCandidateNode(nodeID enode.IDv0) bool {
	return vp.agency.IsCandidateNode(nodeID)
}

// Len return number of validators.
func (vp *ValidatorPool) Len(epoch uint64) int {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators.Len()
	}
	return vp.currentValidators.Len()
}

// Verify verifies signature using the specified validator's bls public key.
func (vp *ValidatorPool) Verify(epoch uint64, validatorIndex uint32, msg, signature []byte) error {
	validator, err := vp.GetValidatorByIndex(epoch, validatorIndex)
	if err != nil {
		return err
	}

	return validator.Verify(msg, signature)
}

// VerifyAggSig verifies aggregation signature using the specified validators' public keys.
func (vp *ValidatorPool) VerifyAggSig(epoch uint64, validatorIndexes []uint32, msg, signature []byte) bool {
	vp.lock.RLock()
	validators := vp.currentValidators
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		validators = vp.prevValidators
	}

	nodeList, err := validators.NodeListByIndexes(validatorIndexes)
	if err != nil {
		vp.lock.RUnlock()
		return false
	}
	vp.lock.RUnlock()

	var pub bls.PublicKey
	for _, node := range nodeList {
		pub.Add(node.BlsPubKey)
	}

	var sig bls.Sign
	err = sig.Deserialize(signature)
	if err != nil {
		return false
	}
	return sig.Verify(&pub, string(msg))
}

func (vp *ValidatorPool) VerifyAggSigByBA(epoch uint64, vSet *utils.BitArray, msg, signature []byte) error {
	vp.lock.RLock()
	validators := vp.currentValidators
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		validators = vp.prevValidators
	}

	nodeList, err := validators.NodeListByBitArray(vSet)
	if err != nil || len(nodeList) == 0 {
		vp.lock.RUnlock()
		return fmt.Errorf("not found validators: %v", err)
	}
	vp.lock.RUnlock()

	var pub bls.PublicKey
	pub = *nodeList[0].BlsPubKey
	//pub.Deserialize(nodeList[0].BlsPubKey.Serialize())
	for i := 1; i < len(nodeList); i++ {
		pub.Add(nodeList[i].BlsPubKey)
	}

	var sig bls.Sign
	err = sig.Deserialize(signature)
	if err != nil {
		return err
	}
	if !sig.Verify(&pub, string(msg)) {
		log.Error("Verify signature fail", "epoch", epoch, "vSet", vSet.String(), "msg", hex.EncodeToString(msg), "signature", hex.EncodeToString(signature), "nodeList", nodeList, "validators", validators.String())
		return errors.New("bls verifies signature fail")
	}
	return nil
}

func (vp *ValidatorPool) epochToBlockNumber(epoch uint64) uint64 {
	if epoch > vp.epoch {
		panic(fmt.Sprintf("get unknown epoch, current:%d, request:%d", vp.epoch, epoch))
	}
	if epoch+1 == vp.epoch {
		return vp.switchPoint
	}
	return vp.switchPoint + 1
}

func (vp *ValidatorPool) Flush(header *types.Header) error {
	return vp.agency.Flush(header)
}

func (vp *ValidatorPool) Commit(block *types.Block) error {
	return vp.agency.OnCommit(block)
}

// NeedGroup return if currentValidators need grouped
func (vp *ValidatorPool) NeedGroup() bool {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	return vp.grouped
}

// GetGroupID return GroupID according epoch & NodeID
func (vp *ValidatorPool) GetGroupID(epoch uint64, nodeID enode.ID) (uint32, error) {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	var validators *cbfttypes.Validators
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		validators = vp.prevValidators
	} else {
		validators = vp.currentValidators
	}
	gvs, err := validators.GetGroupValidators(nodeID)
	if err != nil || gvs == nil {
		return 0, err
	}
	return gvs.GetGroupID(), nil
}

// GetUnitID return index according epoch & NodeID
func (vp *ValidatorPool) GetUnitID(epoch uint64, nodeID enode.ID) (uint32, error) {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators.UnitID(nodeID)
	}
	return vp.currentValidators.UnitID(nodeID)
}

func NextRound(blockNumber uint64) uint64 {
	return blockNumber + 1
}

// Len return number of validators by groupID.
// 返回指定epoch和分组下的节点总数
func (vp *ValidatorPool) LenByGroupID(epoch uint64, groupID uint32) int {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	grouplen := 0
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		grouplen, _ = vp.prevValidators.MembersCount(groupID)
	}
	grouplen, _ = vp.currentValidators.MembersCount(groupID)
	return grouplen
}

//// 查询指定epoch下对应nodeIndex的节点信息，没有对应信息返回nil
//func (vp *ValidatorPool) GetValidatorByGroupIdAndIndex(epoch uint64, nodeIndex uint32) (*cbfttypes.ValidateNode,error) {
//	vp.lock.RLock()
//	defer vp.lock.RUnlock()
//
//	if epoch+1 == vp.epoch {
//		return vp.prevValidators.FindNodeByIndex(int(nodeIndex))
//	}
//	return vp.currentValidators.FindNodeByIndex(int(nodeIndex))
//}

// 返回指定epoch和分组下所有共识节点的index集合，e.g. [25,26,27,28,29,30...49]
func (vp *ValidatorPool) GetValidatorIndexesByGroupID(epoch uint64, groupID uint32) ([]uint32, error) {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		return vp.prevValidators.GetValidatorIndexes(groupID)
	}
	return vp.currentValidators.GetValidatorIndexes(groupID)
}

// 返回指定epoch和分组下所有共识节点的协调节点编组信息，e.g. [[25,26],[27,28],[29,30]...[49]]
// 严格按编组顺序返回
func (vp *ValidatorPool) GetCoordinatorIndexesByGroupID(epoch uint64, groupID uint32) ([][]uint32, error) {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	var validators *cbfttypes.Validators
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		validators = vp.prevValidators
	} else {
		validators = vp.currentValidators
	}

	if groupID >= uint32(len(validators.GroupNodes)) {
		return nil, fmt.Errorf("GetCoordinatorIndexesByGroupID: wrong groupid[%d]", groupID)
	}
	return validators.GroupNodes[groupID].Units, nil
}

// 返回指定epoch下，nodeID所在的groupID和unitID（两者都是0开始计数），没有对应信息需返回error
func (vp *ValidatorPool) GetGroupByValidatorID(epoch uint64, nodeID enode.ID) (uint32, uint32, error) {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	var validators *cbfttypes.Validators
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		validators = vp.prevValidators
	} else {
		validators = vp.currentValidators
	}
	gvs, err := validators.GetGroupValidators(nodeID)
	if nil != err || gvs == nil {
		return 0, 0, err
	}
	unitID, err := validators.UnitID(nodeID)
	return gvs.GetGroupID(), unitID, err
}

// 返回指定epoch下节点的分组信息，key=groupID，value=分组节点index集合
func (vp *ValidatorPool) GetGroupIndexes(epoch uint64) map[uint32][]uint32 {
	vp.lock.RLock()
	defer vp.lock.RUnlock()

	validators := vp.currentValidators
	if vp.epochToBlockNumber(epoch) <= vp.switchPoint {
		validators = vp.prevValidators
	}
	groupIdxs := make(map[uint32][]uint32, len(validators.GroupNodes))
	var err error
	for i, _ := range validators.GroupNodes {
		gid := uint32(i)
		groupIdxs[gid], err = validators.GetValidatorIndexes(gid)
		if nil != err {
			log.Error("GetValidatorIndexes failed!", "err", err)
		}
	}
	return groupIdxs
}

// organize validators into groups
func (vp *ValidatorPool) organize(validators *cbfttypes.Validators, epoch uint64, eventMux *event.TypeMux) error {
	if validators == nil {
		return errors.New("validators is nil")
	}
	err := validators.Grouped()
	if err != nil {
		return err
	}

	gvs, err := validators.GetGroupValidators(vp.nodeID)
	if nil != err || gvs == nil {
		// 当前节点不是共识节点
		return err
	}
	log.Debug("ValidatorPool organized OK!", "epoch", epoch, "validators", validators.String())

	consensusNodes := validators.NodeList()
	groupNodes := gvs.NodeList()
	consensusTopic := cbfttypes.ConsensusTopicName(epoch)
	groupTopic := cbfttypes.ConsensusGroupTopicName(epoch, gvs.GetGroupID())

	eventMux.Post(cbfttypes.NewTopicEvent{Topic: consensusTopic, Nodes: consensusNodes})
	eventMux.Post(cbfttypes.NewTopicEvent{Topic: groupTopic, Nodes: groupNodes})
	eventMux.Post(cbfttypes.GroupTopicEvent{Topic: groupTopic})
	return nil
}

// dissolve prevValidators group
func (vp *ValidatorPool) dissolve(epoch uint64, eventMux *event.TypeMux) {
	if !vp.grouped || vp.prevValidators == nil {
		return
	}
	gvs, err := vp.prevValidators.GetGroupValidators(vp.nodeID)
	if nil != err || gvs == nil {
		// nil != err 说明当前节点上一轮不是共识节点，gvs == nil意味着上一轮共识没分组
		return
	}

	consensusTopic := cbfttypes.ConsensusTopicName(epoch)
	groupTopic := cbfttypes.ConsensusGroupTopicName(epoch, gvs.GetGroupID())

	eventMux.Post(cbfttypes.ExpiredTopicEvent{Topic: consensusTopic})  // for p2p
	eventMux.Post(cbfttypes.ExpiredTopicEvent{Topic: groupTopic})      // for p2p
	eventMux.Post(cbfttypes.ExpiredGroupTopicEvent{Topic: groupTopic}) // for cbft
}
