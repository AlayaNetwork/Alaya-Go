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

package cbfttypes

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/AlayaNetwork/Alaya-Go/common/hexutil"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/x/xcom"

	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/protocols"

	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/utils"

	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/core/types"
	"github.com/AlayaNetwork/Alaya-Go/crypto/bls"
)

const (
	TopicConsensus = "consensus:%d"    // consensus:{epoch}
	TopicGroup     = "consensus:%d:%d" // consensus:{epoch}:{groupID}
)

func ConsensusTopicName(epoch uint64) string {
	return fmt.Sprintf(TopicConsensus, epoch)
}

func ConsensusGroupTopicName(epoch uint64, groupID uint32) string {
	return fmt.Sprintf(TopicGroup, epoch, groupID)
}

type UpdateChainStateFn func(qcState, lockState, commitState *protocols.State)

type CbftResult struct {
	Block              *types.Block
	ExtraData          []byte
	SyncState          chan error
	ChainStateUpdateCB func()
}

type AddValidatorEvent struct {
	Node *enode.Node
}

type RemoveValidatorEvent struct {
	Node *enode.Node
}

type ConsensusTopicEvent struct {
	Topic string // consensus:{epoch}
}

type NewGroupsEvent struct {
	Topic      string // consensus:{epoch}:{groupID}
	Validators *GroupValidators
}

type ExpiredConsensusTopicEvent struct {
	Topic string // consensus:{epoch}
}

type ExpiredTopicEvent struct {
	Topic string // consensus:{epoch}:{groupID}
}

//type UpdateValidatorEvent struct{}

type ValidateNode struct {
	Index     uint32             `json:"index"`
	Address   common.NodeAddress `json:"address"`
	PubKey    *ecdsa.PublicKey   `json:"-"`
	NodeID    enode.ID           `json:"nodeID"`
	BlsPubKey *bls.PublicKey     `json:"blsPubKey"`
}

type ValidateNodeMap map[enode.ID]*ValidateNode

type SortedIndexValidatorNode []*ValidateNode

func (sv SortedIndexValidatorNode) Len() int           { return len(sv) }
func (sv SortedIndexValidatorNode) Swap(i, j int)      { sv[i], sv[j] = sv[j], sv[i] }
func (sv SortedIndexValidatorNode) Less(i, j int) bool { return sv[i].Index < sv[j].Index }

type GroupValidators struct {
	// all nodes in this group
	Nodes []*ValidateNode
	// Coordinators' index  C0>C1>C2>C3...
	Units [][]uint32
	// The group ID
	groupID uint32
}

type Validators struct {
	Nodes            ValidateNodeMap `json:"validateNodes"`
	ValidBlockNumber uint64          `json:"validateBlockNumber"`

	// Sorting based on node index
	SortedNodes SortedIndexValidatorNode `json:"sortedNodes"`

	//// Sorting based on node index
	// Node grouping info
	GroupNodes []*GroupValidators `json:"groupNodes"`
}

func (vn *ValidateNode) String() string {
	b, _ := json.Marshal(vn)
	return string(b)
}

func (vn *ValidateNode) Verify(data, sign []byte) error {
	var sig bls.Sign
	err := sig.Deserialize(sign)
	if err != nil {
		return err
	}

	if !sig.Verify(vn.BlsPubKey, string(data)) {
		return errors.New(fmt.Sprintf("bls verifies signature fail, data:%s, sign:%s, pubkey:%s", hexutil.Encode(data), hexutil.Encode(sign), hexutil.Encode(vn.BlsPubKey.Serialize())))
	}
	return nil
}

func (vnm ValidateNodeMap) String() string {
	s := ""
	for k, v := range vnm {
		s = s + fmt.Sprintf("{%s:%s},", k, v)
	}
	return s
}

func (vs *Validators) String() string {
	b, _ := json.Marshal(vs)
	return string(b)
}

func (vs *Validators) NodeList() []enode.ID {
	nodeList := make([]enode.ID, 0)
	for id := range vs.Nodes {
		nodeList = append(nodeList, id)
	}
	return nodeList
}

func (vs *Validators) MembersCount(groupID uint32) (int, error) {
	if groupID >= uint32(len(vs.GroupNodes)) {
		return 0, fmt.Errorf("wrong groupid[%d]", groupID)
	}
	return len(vs.GroupNodes[groupID].Nodes), nil
}

func (vs *Validators) GetValidatorIndexes(groupid uint32) ([]uint32, error) {
	if groupid >= uint32(len(vs.GroupNodes)) {
		return nil, fmt.Errorf("MembersCount: wrong groupid[%d]", groupid)
	}
	ids := make([]uint32, 0)
	for _, node := range vs.GroupNodes[groupid].Nodes {
		ids = append(ids, node.Index)
	}
	return ids, nil
}

func (vs *Validators) NodeListByIndexes(indexes []uint32) ([]*ValidateNode, error) {
	if len(vs.SortedNodes) == 0 {
		vs.Sort()
	}
	l := make([]*ValidateNode, 0)
	for _, index := range indexes {
		if int(index) >= len(vs.SortedNodes) {
			return nil, errors.New("invalid index")
		}
		l = append(l, vs.SortedNodes[int(index)])
	}
	return l, nil
}

func (vs *Validators) NodeListByBitArray(vSet *utils.BitArray) ([]*ValidateNode, error) {
	if len(vs.SortedNodes) == 0 {
		vs.Sort()
	}
	l := make([]*ValidateNode, 0)

	for index := uint32(0); index < vSet.Size(); index++ {
		if vSet.GetIndex(index) {
			if int(index) >= len(vs.SortedNodes) {
				return nil, errors.New("invalid index")
			}
			l = append(l, vs.SortedNodes[int(index)])
		}
	}
	return l, nil
}

func (vs *Validators) FindNodeByID(id enode.ID) (*ValidateNode, error) {
	node, ok := vs.Nodes[id]
	if ok {
		return node, nil
	}
	return nil, errors.New("not found the node")
}

func (vs *Validators) FindNodeByIndex(index int) (*ValidateNode, error) {
	if len(vs.SortedNodes) == 0 {
		vs.Sort()
	}
	if index >= len(vs.SortedNodes) {
		return nil, errors.New("not found the specified validator")
	} else {
		return vs.SortedNodes[index], nil
	}
}

func (vs *Validators) FindNodeByAddress(addr common.NodeAddress) (*ValidateNode, error) {
	for _, node := range vs.Nodes {
		if bytes.Equal(node.Address[:], addr[:]) {
			return node, nil
		}
	}
	return nil, errors.New("invalid address")
}

func (vs *Validators) NodeID(idx int) enode.ID {
	if len(vs.SortedNodes) == 0 {
		vs.Sort()
	}
	if idx >= vs.SortedNodes.Len() {
		return enode.ID{}
	}
	return vs.SortedNodes[idx].NodeID
}

func (vs *Validators) Index(nodeID enode.ID) (uint32, error) {
	if node, ok := vs.Nodes[nodeID]; ok {
		return node.Index, nil
	}
	return math.MaxUint32, errors.New("not found the specified validator")
}

func (vs *Validators) Len() int {
	return len(vs.Nodes)
}

func (vs *Validators) Equal(rsh *Validators) bool {
	if vs.Len() != rsh.Len() {
		return false
	}

	equal := true
	for k, v := range vs.Nodes {
		if vv, ok := rsh.Nodes[k]; !ok || vv.Index != v.Index {
			equal = false
			break
		}
	}
	return equal
}

func (vs *Validators) Sort() {
	for _, node := range vs.Nodes {
		vs.SortedNodes = append(vs.SortedNodes, node)
	}
	sort.Sort(vs.SortedNodes)
}

func (vs *Validators) GetGroupValidators(nodeID enode.ID) (*GroupValidators, error) {
	if len(vs.SortedNodes) == 0 {
		vs.Sort()
	}

	idx, err := vs.Index(nodeID)
	if err != nil {
		return nil, err
	}

	var ret *GroupValidators
	for _, gvs := range vs.GroupNodes {
		goupLen := len(gvs.Nodes)
		if idx <= gvs.Nodes[goupLen-1].Index {
			ret = gvs
			break
		}
	}
	return ret, nil
}

func (vs *Validators) UnitID(nodeID enode.ID) (uint32, error) {
	if len(vs.SortedNodes) == 0 {
		vs.Sort()
	}

	idx, err := vs.Index(nodeID)
	if err != nil {
		return idx, err
	}

	gvs, err := vs.GetGroupValidators(nodeID)
	if err != nil {
		return 0, err
	}
	for i, node := range gvs.Nodes {
		if idx == node.Index {
			return uint32(i), nil
		}
	}
	return 0, errors.New("not found the specified nodeID")
}

func (gvs *GroupValidators) GroupedUnits() {
	coordinatorLimit := xcom.CoordinatorsLimit()
	unit := make([]uint32, 0, coordinatorLimit)
	for i, n := range gvs.Nodes {
		unit = append(unit, n.Index)
		if uint32(len(unit)) >= coordinatorLimit || i == len(gvs.Nodes)-1 {
			gvs.Units = append(gvs.Units, unit)
			unit = make([]uint32, 0, coordinatorLimit)
		}
	}
}

func (gvs *GroupValidators) GetGroupID() uint32 {
	return gvs.groupID
}

// Grouped fill validators into groups
// groupValidatorsLimit is a factor to determine how many groups are grouped
// eg: [validatorCount,groupValidatorsLimit]=
// [50,25] = 25,25;[43,25] = 22,21; [101,25] = 21,20,20,20,20
func (vs *Validators) Grouped() error {
	// sort nodes by index
	if len(vs.SortedNodes) == 0 {
		vs.Sort()
	}

	validatorCount := uint32(vs.SortedNodes.Len())
	groupNum := validatorCount / xcom.MaxGroupValidators()
	mod := validatorCount % xcom.CoordinatorsLimit()
	if mod > 0 {
		groupNum = groupNum + 1
	}

	memberMinCount := validatorCount / groupNum
	remainder := validatorCount % groupNum
	vs.GroupNodes = make([]*GroupValidators, groupNum, groupNum)
	begin := uint32(0)
	end := uint32(0)
	for i := uint32(0); i < groupNum; i++ {
		begin = end
		if remainder > 0 {
			end = begin + memberMinCount + 1
			remainder = remainder - 1
		} else {
			end = begin + memberMinCount
		}
		if end > validatorCount {
			end = validatorCount
		}
		groupValidators := new(GroupValidators)
		groupValidators.Nodes = vs.SortedNodes[begin:end]
		groupValidators.groupID = i
		vs.GroupNodes[i] = groupValidators
	}

	// fill group unit
	for _, gvs := range vs.GroupNodes {
		gvs.GroupedUnits()
	}
	return nil
}
