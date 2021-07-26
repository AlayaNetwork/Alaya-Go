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


package xcom

import (
	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/rlp"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDefaultEMConfig(t *testing.T) {
	t.Run("DefaultAlayaNet", func(t *testing.T) {
		if getDefaultEMConfig(DefaultAlayaNet) == nil {
			t.Error("DefaultAlayaNet can't be nil config")
		}
		if err := CheckEconomicModel(); nil != err {
			t.Error(err)
		}
	})
	t.Run("DefaultUnitTestNet", func(t *testing.T) {
		if getDefaultEMConfig(DefaultUnitTestNet) == nil {
			t.Error("DefaultUnitTestNet can't be nil config")
		}
		if err := CheckEconomicModel(); nil != err {
			t.Error(err)
		}
	})
	if getDefaultEMConfig(10) != nil {
		t.Error("the chain config not support")
	}
}

func TestEcParams0140(t *testing.T) {
	eceHash := "0xbd45f1783a2344776066ca1a88937e74dfba777c9c3eb2f6819989c66d2c0462"
	getDefaultEMConfig(DefaultAlayaNet)
	if bytes, err := EcParams0140(); nil != err {
		t.Fatal(err)
	} else {
		assert.True(t, bytes != nil)
		assert.True(t, common.RlpHash(bytes).Hex() == eceHash)
	}
}

func TestAlayaNetHash(t *testing.T) {
	alayaEc := getDefaultEMConfig(DefaultAlayaNet)
	bytes, err := rlp.EncodeToBytes(alayaEc)
	if err != nil {
		t.Error(err)
	}
	assert.True(t, common.RlpHash(bytes).Hex() == AlayaNetECHash)
}
