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
	"encoding/json"
	"fmt"
	"github.com/AlayaNetwork/Alaya-Go/consensus/cbft/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/crypto/bls"
)

func Test_QuorumCert(t *testing.T) {
	qc := &QuorumCert{
		Epoch:        1,
		ViewNumber:   1,
		BlockHash:    common.BytesToHash(utils.Rand32Bytes(32)),
		BlockNumber:  2,
		BlockIndex:   0,
		Signature:    BytesToSignature(utils.Rand32Bytes(64)),
		ValidatorSet: utils.NewBitArray(25),
	}
	qc.ValidatorSet.SetIndex(0, true)
	qc.ValidatorSet.SetIndex(24, true)
	_, err := qc.CannibalizeBytes()
	assert.Nil(t, err)
	assert.Equal(t, 2, qc.Len())
	assert.NotEmpty(t, qc.String())
	assert.True(t, qc.HigherQuorumCert(1, 1, 0))
	assert.True(t, qc.HigherQuorumCert(2, 1, 0))
}

func Test_ViewChangeQC(t *testing.T) {
	viewChangeQC := new(ViewChangeQC)
	hash1 := common.BytesToHash(utils.Rand32Bytes(32))
	hash2 := common.BytesToHash(utils.Rand32Bytes(32))
	hash3 := common.BytesToHash(utils.Rand32Bytes(32))
	viewChangeQC.AppendQuorumCert(makeViewChangeQuorumCert(2, 3, hash1, 9, 2, 2))
	viewChangeQC.AppendQuorumCert(makeViewChangeQuorumCert(2, 3, hash2, 9, 2, 3))
	assert.True(t, viewChangeQC.ExistViewChange(2, 3, hash2))
	assert.NotEmpty(t, viewChangeQC.String())
	assert.Equal(t, 2, viewChangeQC.Len())
	viewChangeQC.AppendQuorumCert(makeViewChangeQuorumCert(2, 4, hash3, 9, 2, 3))
	assert.NotNil(t, viewChangeQC.EqualAll(2, 3))
	last := viewChangeQC.QCs[len(viewChangeQC.QCs)-1]
	copy := last.Copy()
	assert.NotEmpty(t, copy.String())
	_, err := copy.CannibalizeBytes()
	assert.Nil(t, err)
	assert.Equal(t, hash3, copy.BlockHash)
}

func makeViewChangeQuorumCert(epoch, viewNumber uint64, blockHash common.Hash, blockNumber uint64, blockEpoch, blockViewNumber uint64) *ViewChangeQuorumCert {
	cert := &ViewChangeQuorumCert{
		Epoch:           epoch,
		ViewNumber:      viewNumber,
		BlockHash:       blockHash,
		BlockNumber:     blockNumber,
		BlockEpoch:      blockEpoch,
		BlockViewNumber: blockViewNumber,
		Signature:       BytesToSignature(utils.Rand32Bytes(64)),
		ValidatorSet:    utils.NewBitArray(25),
	}
	cert.ValidatorSet.SetIndex(0, true)
	return cert
}

func Test_ViewChangeQC_MaxBlock(t *testing.T) {
	certs := []*ViewChangeQuorumCert{
		makeViewChangeQuorumCert(2, 3, common.BytesToHash(utils.Rand32Bytes(32)), 9, 2, 1),
		makeViewChangeQuorumCert(2, 3, common.BytesToHash(utils.Rand32Bytes(32)), 9, 2, 3),
		makeViewChangeQuorumCert(2, 3, common.BytesToHash(utils.Rand32Bytes(32)), 10, 2, 1),
		makeViewChangeQuorumCert(2, 3, common.BytesToHash(utils.Rand32Bytes(32)), 10, 2, 1),
		makeViewChangeQuorumCert(2, 3, common.BytesToHash(utils.Rand32Bytes(32)), 10, 2, 2),
		makeViewChangeQuorumCert(2, 3, common.BytesToHash(utils.Rand32Bytes(32)), 10, 1, 25),
	}
	viewChangeQC := &ViewChangeQC{
		QCs: certs,
	}
	epoch, viewNumber, blockEpoch, blockViewNumber, blockHash, blockNumber := viewChangeQC.MaxBlock()
	assert.Equal(t, certs[4].Epoch, epoch)
	assert.Equal(t, certs[4].ViewNumber, viewNumber)
	assert.Equal(t, certs[4].BlockEpoch, blockEpoch)
	assert.Equal(t, certs[4].BlockViewNumber, blockViewNumber)
	assert.Equal(t, certs[4].BlockHash, blockHash)
	assert.Equal(t, certs[4].BlockNumber, blockNumber)

	viewChangeQC.QCs = nil
	epoch, viewNumber, blockEpoch, blockViewNumber, blockHash, blockNumber = viewChangeQC.MaxBlock()
	assert.Equal(t, uint64(0), epoch)
}

func TestViewChangeQC_ValidatorSet(t *testing.T) {
	testCases := []struct {
		ValidatorSetStr string
	}{
		{`"x_x_x_xxxx"`},
		{`"xxxxxx"`},
		{`"xx__________"`},
		{`"x_x_x_______"`},
		{`"xx__x_______"`},
		{`"x_x_x_xxxx"`},
		{`"______x_____"`},
		{`"______xxxx__"`},
		{`"______xx____"`},
		{`"______x_x_x_"`},
		{`"______xx__x_"`},
		{`"______xxx_x____"`},
	}

	bitArray := func(bitArrayStr string) *utils.BitArray {
		var ba *utils.BitArray
		json.Unmarshal([]byte(bitArrayStr), &ba)
		return ba

	}

	viewChangeQC := &ViewChangeQC{QCs: make([]*ViewChangeQuorumCert, 0)}
	for _, c := range testCases {
		qc := &ViewChangeQuorumCert{
			ValidatorSet: bitArray(c.ValidatorSetStr),
		}
		viewChangeQC.QCs = append(viewChangeQC.QCs, qc)
	}
	assert.Equal(t, 45, viewChangeQC.Len())
	assert.Equal(t, uint32(15), viewChangeQC.ValidatorSet().Size())
	assert.Equal(t, 11, viewChangeQC.HasLength())
}

func TestViewChangeQC_HasSign(t *testing.T) {
	b1 := utils.NewBitArray(1000)
	b1.SetIndex(uint32(100), true)
	b1.SetIndex(uint32(666), true)
	b1.SetIndex(uint32(888), true)
	b1.SetIndex(uint32(999), true)

	b2 := utils.NewBitArray(1000)
	b2.SetIndex(uint32(0), true)
	b2.SetIndex(uint32(100), true)
	b2.SetIndex(uint32(222), true)
	b2.SetIndex(uint32(333), true)

	b3 := utils.NewBitArray(1000)
	b3.SetIndex(uint32(666), false)
	b3.SetIndex(uint32(777), true)

	testCases := []*utils.BitArray{
		b1, b2, b3,
	}

	viewChangeQC := &ViewChangeQC{QCs: make([]*ViewChangeQuorumCert, 0)}
	for _, c := range testCases {
		qc := &ViewChangeQuorumCert{
			ValidatorSet: c,
		}
		viewChangeQC.QCs = append(viewChangeQC.QCs, qc)
	}
	assert.True(t, true, viewChangeQC.HasSign(uint32(0)))
	assert.True(t, true, viewChangeQC.HasSign(uint32(666)))
	assert.False(t, false, viewChangeQC.HasSign(uint32(555)))
}

func TestViewChangeQC_HigherSign(t *testing.T) {
	b1 := utils.NewBitArray(1000)
	b1.SetIndex(uint32(100), true)
	b1.SetIndex(uint32(666), true)
	b1.SetIndex(uint32(888), true)
	b1.SetIndex(uint32(999), true)

	b2 := utils.NewBitArray(1000)
	b2.SetIndex(uint32(0), true)
	b2.SetIndex(uint32(100), true)
	b2.SetIndex(uint32(222), true)
	b2.SetIndex(uint32(333), true)

	b3 := utils.NewBitArray(1000)
	b3.SetIndex(uint32(666), false)
	b3.SetIndex(uint32(777), true)

	testCases := []*utils.BitArray{
		b1, b2, b3,
	}

	viewChangeQC1 := &ViewChangeQC{QCs: make([]*ViewChangeQuorumCert, 0)}
	viewChangeQC2 := &ViewChangeQC{QCs: make([]*ViewChangeQuorumCert, 0)}
	for i, c := range testCases {
		qc := &ViewChangeQuorumCert{
			ValidatorSet: c,
		}
		if i <= len(testCases)-2 {
			viewChangeQC1.QCs = append(viewChangeQC1.QCs, qc)
			viewChangeQC2.QCs = append(viewChangeQC2.QCs, qc)
		}
		if i == len(testCases)-1 {
			viewChangeQC1.QCs = append(viewChangeQC1.QCs, qc)
		}
	}
	assert.Equal(t, 8, viewChangeQC1.HasLength())
	assert.Equal(t, 7, viewChangeQC2.HasLength())
	assert.True(t, true, viewChangeQC1.HigherSign(viewChangeQC2))
}

func TestViewChangeQuorumCert_AddSign(t *testing.T) {
	bls.Init(int(bls.BLS12_381))
	message := "test merge sign"
	var k int = 500
	msk := make([]bls.SecretKey, k)
	mpk := make([]bls.PublicKey, k)
	msig := make([]bls.Sign, k)
	for i := 0; i < k; i++ {
		msk[i].SetByCSPRNG()
		mpk[i] = *msk[i].GetPublicKey()
		msig[i] = *msk[i].Sign(message)
	}

	verifyViewChangeQuorumCert := func(qc *ViewChangeQuorumCert) bool {
		var pub bls.PublicKey
		for i := uint32(0); i < qc.ValidatorSet.Size(); i++ {
			if qc.ValidatorSet.GetIndex(i) {
				pub.Add(&mpk[i])
			}
		}
		var sig bls.Sign
		if err := sig.Deserialize(qc.Signature.Bytes()); err != nil {
			return false
		}

		if sig.Verify(&pub, message) {
			return true
		}
		return false
	}

	var sig bls.Sign
	vSet := utils.NewBitArray(uint32(k))
	for i := 0; i < len(msig)-2; i++ {
		sig.Add(&msig[i])
		vSet.SetIndex(uint32(i), true)
	}

	qc := &ViewChangeQuorumCert{
		ValidatorSet: vSet,
	}
	qc.Signature.SetBytes(sig.Serialize())
	//fmt.Println("qc Signature", qc.Signature.String())
	assert.Equal(t, true, verifyViewChangeQuorumCert(qc))

	// add sign and verify sign
	for i := len(msig) - 2; i < len(msig); i++ {
		var s Signature
		s.SetBytes(msig[i].Serialize())
		qc.AddSign(s, uint32(i))
		//fmt.Println("qc Signature", qc.Signature.String())
		assert.Equal(t, true, verifyViewChangeQuorumCert(qc))
	}

	// The public key does not match and cannot be verified
	var s Signature
	s.SetBytes(msig[0].Serialize())
	qc.AddSign(s, uint32(0))
	assert.Equal(t, false, verifyViewChangeQuorumCert(qc))
}

func TestQuorumCert_HasSign(t *testing.T) {
	b1 := utils.NewBitArray(1000)
	b1.SetIndex(uint32(100), true)
	b1.SetIndex(uint32(666), true)
	b1.SetIndex(uint32(888), true)
	b1.SetIndex(uint32(999), true)

	b2 := utils.NewBitArray(1000)
	b2.SetIndex(uint32(0), true)
	b2.SetIndex(uint32(100), true)
	b2.SetIndex(uint32(222), true)
	b2.SetIndex(uint32(333), true)

	qc1 := &QuorumCert{
		ValidatorSet: b1,
	}

	qc2 := &QuorumCert{
		ValidatorSet: b1.Or(b2),
	}
	assert.Equal(t, 4, qc1.Len())
	assert.Equal(t, 7, qc2.Len())
	assert.True(t, true, qc1.HigherSign(qc2))
}

func TestQuorumCert_HigherSign(t *testing.T) {
	b1 := utils.NewBitArray(1000)
	b1.SetIndex(uint32(100), true)
	b1.SetIndex(uint32(666), true)
	b1.SetIndex(uint32(888), true)
	b1.SetIndex(uint32(999), true)
	qc := &QuorumCert{
		ValidatorSet: b1,
	}
	assert.True(t, true, qc.HasSign(uint32(666)))
	assert.False(t, false, qc.HasSign(uint32(777)))
}

func TestQuorumCert_AddSign(t *testing.T) {
	bls.Init(int(bls.BLS12_381))
	message := "test merge sign"
	var k int = 500
	msk := make([]bls.SecretKey, k)
	mpk := make([]bls.PublicKey, k)
	msig := make([]bls.Sign, k)
	for i := 0; i < k; i++ {
		msk[i].SetByCSPRNG()
		mpk[i] = *msk[i].GetPublicKey()
		msig[i] = *msk[i].Sign(message)
	}

	verifyQuorumCert := func(qc *QuorumCert) bool {
		var pub bls.PublicKey
		for i := uint32(0); i < qc.ValidatorSet.Size(); i++ {
			if qc.ValidatorSet.GetIndex(i) {
				pub.Add(&mpk[i])
			}
		}
		var sig bls.Sign
		if err := sig.Deserialize(qc.Signature.Bytes()); err != nil {
			return false
		}

		if sig.Verify(&pub, message) {
			return true
		}
		return false
	}

	var sig bls.Sign
	vSet := utils.NewBitArray(uint32(k))
	for i := 0; i < len(msig)-2; i++ {
		sig.Add(&msig[i])
		vSet.SetIndex(uint32(i), true)
	}

	qc := &QuorumCert{
		ValidatorSet: vSet,
	}
	qc.Signature.SetBytes(sig.Serialize())
	//fmt.Println("qc Signature", qc.Signature.String())
	assert.Equal(t, true, verifyQuorumCert(qc))

	// add sign and verify sign
	for i := len(msig) - 2; i < len(msig); i++ {
		var s Signature
		s.SetBytes(msig[i].Serialize())
		qc.AddSign(s, uint32(i))
		//fmt.Println("qc Signature", qc.Signature.String())
		assert.Equal(t, true, verifyQuorumCert(qc))
	}

	// The public key does not match and cannot be verified
	var s Signature
	s.SetBytes(msig[0].Serialize())
	qc.AddSign(s, uint32(0))
	assert.Equal(t, false, verifyQuorumCert(qc))
}

func TestAddSign(t *testing.T) {
	bls.Init(int(bls.BLS12_381))
	message := "test merge sign"
	var k int = 100000
	msk := make([]bls.SecretKey, k)
	mpk := make([]bls.PublicKey, k)
	msig := make([]bls.Sign, k)
	msignature := make([]Signature, k)
	for i := 0; i < 1; i++ {
		msk[i].SetByCSPRNG()
		mpk[i] = *msk[i].GetPublicKey()
		msig[i] = *msk[i].Sign(message)
		msignature[i].SetBytes(msig[i].Serialize())
	}

	qc := &QuorumCert{
		Signature:    Signature{},
		ValidatorSet: utils.NewBitArray(uint32(300)),
	}
	qc.Signature.SetBytes(msig[0].Serialize())
	qc.ValidatorSet.SetIndex(0, true)

	start := common.Millis(time.Now())
	fmt.Println("test", "start", start)
	for i := 1; i < k; i++ {
		qc.AddSign(qc.Signature, uint32(i%300))
	}
	end := common.Millis(time.Now())
	fmt.Println("test", "end", end, "v", qc.ValidatorSet.HasLength())
}

func TestUnKnownGroups_UnKnownSize(t *testing.T) {
	testCases := []struct {
		groupID         uint32
		ValidatorSetStr string
	}{
		{1, `"xx__________"`},
		{2, `"______xxxxxx"`},
		{3, `"xx__x___xx__"`},
	}

	bitArray := func(bitArrayStr string) *utils.BitArray {
		var ba *utils.BitArray
		json.Unmarshal([]byte(bitArrayStr), &ba)
		return ba

	}

	unKnownGroups := &UnKnownGroups{UnKnown: make([]*UnKnownGroup, 0)}
	for _, c := range testCases {
		unKnownGroups.UnKnown = append(unKnownGroups.UnKnown, &UnKnownGroup{
			GroupID:    c.groupID,
			UnKnownSet: bitArray(c.ValidatorSetStr),
		})
	}
	assert.Equal(t, 9, unKnownGroups.UnKnownSize())
}
