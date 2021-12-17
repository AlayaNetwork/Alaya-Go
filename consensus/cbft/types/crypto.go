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

package types

import (
	"fmt"
	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/common/hexutil"
	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/utils"
	"github.com/AlayaNetwork/Alaya-Go/crypto"
	"github.com/AlayaNetwork/Alaya-Go/crypto/bls"
	"github.com/AlayaNetwork/Alaya-Go/rlp"
	"reflect"
)

const (
	SignatureLength = 64
)

type Signature [SignatureLength]byte

func (sig *Signature) String() string {
	return fmt.Sprintf("%x", sig[:])
}

func (sig *Signature) SetBytes(signSlice []byte) {
	copy(sig[:], signSlice[:])
}

func (sig *Signature) Bytes() []byte {
	target := make([]byte, len(sig))
	copy(target[:], sig[:])
	return target
}

// MarshalText returns the hex representation of a.
func (sig Signature) MarshalText() ([]byte, error) {
	return hexutil.Bytes(sig[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (sig *Signature) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("BlockConfirmSign", input, sig[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (sig *Signature) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(reflect.TypeOf(Signature{}), input, sig[:])
}

func BytesToSignature(signSlice []byte) Signature {
	var sign Signature
	copy(sign[:], signSlice[:])
	return sign
}

type QuorumCert struct {
	Epoch        uint64          `json:"epoch"`
	ViewNumber   uint64          `json:"viewNumber"`
	BlockHash    common.Hash     `json:"blockHash"`
	BlockNumber  uint64          `json:"blockNumber"`
	BlockIndex   uint32          `json:"blockIndex"`
	Signature    Signature       `json:"signature"`
	ValidatorSet *utils.BitArray `json:"validatorSet"`
}

func (q *QuorumCert) DeepCopyQuorumCert() *QuorumCert {
	if q == nil {
		return nil
	}
	qc := &QuorumCert{
		Epoch:        q.Epoch,
		ViewNumber:   q.ViewNumber,
		BlockHash:    q.BlockHash,
		BlockNumber:  q.BlockNumber,
		BlockIndex:   q.BlockIndex,
		Signature:    q.Signature,
		ValidatorSet: q.ValidatorSet.Copy(),
	}
	return qc
}

func (q QuorumCert) CannibalizeBytes() ([]byte, error) {
	buf, err := rlp.EncodeToBytes([]interface{}{
		q.Epoch,
		q.ViewNumber,
		q.BlockHash,
		q.BlockNumber,
		q.BlockIndex,
	})
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(buf), nil
}

func (q *QuorumCert) Len() int {
	if q == nil || q.ValidatorSet == nil {
		return 0
	}

	length := 0
	for i := uint32(0); i < q.ValidatorSet.Size(); i++ {
		if q.ValidatorSet.GetIndex(i) {
			length++
		}
	}
	return length
}

func (q *QuorumCert) String() string {
	if q == nil {
		return ""
	}
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,Hash:%s,Number:%d,Index:%d,Signature:%s,ValidatorSet:%s}", q.Epoch, q.ViewNumber, q.BlockHash.TerminalString(), q.BlockNumber, q.BlockIndex, q.Signature.String(), q.ValidatorSet.String())
}

// Add a new signature to the aggregate signature
// Note: Call this method to ensure that the new signature does not exist in the aggregate signature, otherwise the entire aggregate signature will be wrong
func (q *QuorumCert) AddSign(sign Signature, NodeIndex uint32) bool {
	if q == nil {
		return false
	}
	var (
		addSig bls.Sign
		blsSig bls.Sign
	)
	if err := addSig.Deserialize(sign.Bytes()); err != nil {
		return false
	}
	if err := blsSig.Deserialize(q.Signature.Bytes()); err != nil {
		return false
	}
	blsSig.Add(&addSig)
	q.Signature.SetBytes(blsSig.Serialize())
	q.ValidatorSet.SetIndex(NodeIndex, true)
	return true
}

func (q *QuorumCert) HigherSign(c *QuorumCert) bool {
	if q == nil && c == nil {
		return false
	}
	if q == nil && c != nil {
		return false
	}
	if c == nil {
		return true
	}
	if !q.EqualState(c) {
		return false
	}
	return q.ValidatorSet.HasLength() > c.ValidatorSet.HasLength()
}

func (q *QuorumCert) HasSign(signIndex uint32) bool {
	if q == nil || q.ValidatorSet == nil {
		return false
	}
	return q.ValidatorSet.GetIndex(signIndex)
}

func (q *QuorumCert) EqualState(c *QuorumCert) bool {
	return q.Epoch == c.Epoch &&
		q.ViewNumber == c.ViewNumber &&
		q.BlockHash == c.BlockHash &&
		q.BlockNumber == c.BlockNumber &&
		q.BlockIndex == c.BlockIndex &&
		q.ValidatorSet.Size() == c.ValidatorSet.Size()
}

// if the two quorumCert have the same blockNumber
func (q *QuorumCert) HigherBlockView(blockEpoch, blockView uint64) bool {
	return q.Epoch > blockEpoch || (q.Epoch == blockEpoch && q.ViewNumber > blockView)
}

func (q *QuorumCert) HigherQuorumCert(blockNumber uint64, blockEpoch, blockView uint64) bool {
	if q.BlockNumber > blockNumber {
		return true
	} else if q.BlockNumber == blockNumber {
		return q.HigherBlockView(blockEpoch, blockView)
	}
	return false
}

type ViewChangeQuorumCert struct {
	Epoch           uint64          `json:"epoch"`
	ViewNumber      uint64          `json:"viewNumber"`
	BlockHash       common.Hash     `json:"blockHash"`
	BlockNumber     uint64          `json:"blockNumber"`
	BlockEpoch      uint64          `json:"blockEpoch"`
	BlockViewNumber uint64          `json:"blockViewNumber"`
	Signature       Signature       `json:"signature"`
	ValidatorSet    *utils.BitArray `json:"validatorSet"`
}

func (q *ViewChangeQuorumCert) DeepCopyViewChangeQuorumCert() *ViewChangeQuorumCert {
	if q == nil {
		return nil
	}
	qc := &ViewChangeQuorumCert{
		Epoch:           q.Epoch,
		ViewNumber:      q.ViewNumber,
		BlockHash:       q.BlockHash,
		BlockNumber:     q.BlockNumber,
		BlockEpoch:      q.BlockEpoch,
		BlockViewNumber: q.BlockViewNumber,
		Signature:       q.Signature,
		ValidatorSet:    q.ValidatorSet.Copy(),
	}
	return qc
}

func (q ViewChangeQuorumCert) CannibalizeBytes() ([]byte, error) {
	buf, err := rlp.EncodeToBytes([]interface{}{
		q.Epoch,
		q.ViewNumber,
		q.BlockHash,
		q.BlockNumber,
		q.BlockEpoch,
		q.BlockViewNumber,
	})
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(buf), nil
}

func (q *ViewChangeQuorumCert) Len() int {
	if q.ValidatorSet == nil {
		return 0
	}

	length := 0
	for i := uint32(0); i < q.ValidatorSet.Size(); i++ {
		if q.ValidatorSet.GetIndex(i) {
			length++
		}
	}
	return length
}

func (q ViewChangeQuorumCert) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,Hash:%s,Number:%d,BlockEpoch:%d,BlockViewNumber:%d:ValidatorSet:%s}", q.Epoch, q.ViewNumber, q.BlockHash.TerminalString(), q.BlockNumber, q.BlockEpoch, q.BlockViewNumber, q.ValidatorSet.String())
}

// if the two quorumCert have the same blockNumber
func (q *ViewChangeQuorumCert) HigherBlockView(blockEpoch, blockView uint64) bool {
	return q.BlockEpoch > blockEpoch || (q.BlockEpoch == blockEpoch && q.BlockViewNumber > blockView)
}

func (q *ViewChangeQuorumCert) HigherQuorumCert(c *ViewChangeQuorumCert) bool {
	if q.BlockNumber > c.BlockNumber {
		return true
	} else if q.BlockNumber == c.BlockNumber {
		return q.HigherBlockView(c.BlockEpoch, c.BlockViewNumber)
	}
	return false
}

func (q *ViewChangeQuorumCert) Copy() *ViewChangeQuorumCert {
	return &ViewChangeQuorumCert{
		Epoch:        q.Epoch,
		ViewNumber:   q.ViewNumber,
		BlockHash:    q.BlockHash,
		BlockNumber:  q.BlockNumber,
		Signature:    q.Signature,
		ValidatorSet: q.ValidatorSet.Copy(),
	}
}

// Add a new signature to the aggregate signature
// Note: Call this method to ensure that the new signature does not exist in the aggregate signature, otherwise the entire aggregate signature will be wrong
func (q *ViewChangeQuorumCert) AddSign(sign Signature, NodeIndex uint32) bool {
	if q == nil {
		return false
	}
	var (
		addSig bls.Sign
		blsSig bls.Sign
	)
	if err := addSig.Deserialize(sign.Bytes()); err != nil {
		return false
	}
	if err := blsSig.Deserialize(q.Signature.Bytes()); err != nil {
		return false
	}
	blsSig.Add(&addSig)
	q.Signature.SetBytes(blsSig.Serialize())
	q.ValidatorSet.SetIndex(NodeIndex, true)
	return true
}

func (q *ViewChangeQuorumCert) HigherSign(c *ViewChangeQuorumCert) bool {
	if q == nil && c == nil {
		return false
	}
	if q == nil && c != nil {
		return false
	}
	if c == nil {
		return true
	}
	if !q.EqualState(c) {
		return false
	}
	return q.ValidatorSet.HasLength() > c.ValidatorSet.HasLength()
}

func (q *ViewChangeQuorumCert) HasSign(signIndex uint32) bool {
	if q == nil || q.ValidatorSet == nil {
		return false
	}
	return q.ValidatorSet.GetIndex(signIndex)
}

func (q *ViewChangeQuorumCert) EqualState(c *ViewChangeQuorumCert) bool {
	return q.Epoch == c.Epoch &&
		q.ViewNumber == c.ViewNumber &&
		q.BlockHash == c.BlockHash &&
		q.BlockNumber == c.BlockNumber &&
		q.BlockEpoch == c.BlockEpoch &&
		q.BlockViewNumber == c.BlockViewNumber &&
		q.ValidatorSet.Size() == c.ValidatorSet.Size()
}

func (v ViewChangeQC) EqualAll(epoch uint64, viewNumber uint64) error {
	for _, v := range v.QCs {
		if v.ViewNumber != viewNumber || v.Epoch != epoch {
			return fmt.Errorf("not equal, local:{%d}, want{%d}", v.ViewNumber, viewNumber)
		}
	}
	return nil
}

type ViewChangeQC struct {
	QCs []*ViewChangeQuorumCert `json:"qcs"`
}

func (v ViewChangeQC) MaxBlock() (uint64, uint64, uint64, uint64, common.Hash, uint64) {
	if len(v.QCs) == 0 {
		return 0, 0, 0, 0, common.Hash{}, 0
	}

	maxQC := v.QCs[0]
	for _, qc := range v.QCs {
		if qc.HigherQuorumCert(maxQC) {
			maxQC = qc
		}
	}
	return maxQC.Epoch, maxQC.ViewNumber, maxQC.BlockEpoch, maxQC.BlockViewNumber, maxQC.BlockHash, maxQC.BlockNumber
}

func (v *ViewChangeQC) Len() int {
	if v == nil || len(v.QCs) <= 0 {
		return 0
	}
	length := 0
	for _, qc := range v.QCs {
		length += qc.Len()
	}
	return length
}

func (v *ViewChangeQC) HasLength() int {
	if v == nil || len(v.QCs) <= 0 {
		return 0
	}
	return v.ValidatorSet().HasLength()
}

func (v *ViewChangeQC) ValidatorSet() *utils.BitArray {
	if len(v.QCs) > 0 {
		vSet := v.QCs[0].ValidatorSet
		for i := 1; i < len(v.QCs); i++ {
			vSet = vSet.Or(v.QCs[i].ValidatorSet)
		}
		return vSet
	}
	return nil
}

func (v *ViewChangeQC) HasSign(signIndex uint32) bool {
	if v == nil || len(v.QCs) <= 0 {
		return false
	}
	for _, qc := range v.QCs {
		if qc.HasSign(signIndex) {
			return true
		}
	}
	return false
}

func (v *ViewChangeQC) HigherSign(c *ViewChangeQC) bool {
	return v.HasLength() > c.HasLength()
}

func (v ViewChangeQC) String() string {
	epoch, view, blockEpoch, blockViewNumber, hash, number := v.MaxBlock()
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,BlockEpoch:%d,BlockViewNumber:%d,Hash:%s,Number:%d}", epoch, view, blockEpoch, blockViewNumber, hash.TerminalString(), number)
}

func (v ViewChangeQC) ExistViewChange(epoch, viewNumber uint64, blockHash common.Hash) bool {
	for _, vc := range v.QCs {
		if vc.Epoch == epoch && vc.ViewNumber == viewNumber && vc.BlockHash == blockHash {
			return true
		}
	}
	return false
}

func (v *ViewChangeQC) AppendQuorumCert(viewChangeQC *ViewChangeQuorumCert) {
	v.QCs = append(v.QCs, viewChangeQC)
}

type PrepareQCs struct {
	QCs []*QuorumCert `json:"qcs"`
}

func (p *PrepareQCs) FindPrepareQC(hash common.Hash) *QuorumCert {
	if p == nil || len(p.QCs) <= 0 {
		return nil
	}
	for _, qc := range p.QCs {
		if qc.BlockHash == hash {
			return qc
		}
	}
	return nil
}

func (p *PrepareQCs) FlattenMap() map[common.Hash]*QuorumCert {
	if p == nil || len(p.QCs) <= 0 {
		return nil
	}
	m := make(map[common.Hash]*QuorumCert)
	for _, qc := range p.QCs {
		m[qc.BlockHash] = qc
	}
	return m
}

func (p *PrepareQCs) AppendQuorumCert(qc *QuorumCert) {
	p.QCs = append(p.QCs, qc)
}

type UnKnownGroups struct {
	UnKnown []*UnKnownGroup `json:"unKnown"`
}

type UnKnownGroup struct {
	GroupID    uint32 `json:"groupID"`
	UnKnownSet *utils.BitArray
}

func (unKnowns *UnKnownGroups) UnKnownSize() int {
	if unKnowns == nil || len(unKnowns.UnKnown) <= 0 {
		return 0
	}

	var unKnownSets *utils.BitArray
	for _, un := range unKnowns.UnKnown {
		if unKnownSets == nil {
			unKnownSets = un.UnKnownSet
		} else {
			unKnownSets = unKnownSets.Or(un.UnKnownSet)
		}
	}
	return unKnownSets.HasLength()
}
