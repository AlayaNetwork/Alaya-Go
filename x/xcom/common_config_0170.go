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

import "github.com/AlayaNetwork/Alaya-Go/rlp"

// New parameters added in version 0.17.0 need to be saved on the chain.
// Calculate the rlp of the new parameter and return it to the upper storage.
func EcParams0170() ([]byte, error) {
	params := struct {
		MaxGroupValidators uint32 `json:"caxGroupValidators"` // max validators count in 1 group
		CoordinatorsLimit  uint32 `json:"coordinatorLimit"`   // max Coordinators count in 1 group
		MaxConsensusVals   uint64 `json:"maxConsensusVals"`   // The consensus validators count
	}{
		MaxGroupValidators: ece.Extend0170.Common.MaxGroupValidators,
		CoordinatorsLimit:  ece.Extend0170.Common.CoordinatorsLimit,
		MaxConsensusVals:   ece.Extend0170.Common.MaxConsensusVals,
	}
	bytes, err := rlp.EncodeToBytes(params)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

type EconomicModel0170Extend struct {
	Common   EconomicModel0170CommonConfig   `json:"common"`
	Staking  EconomicModel0170StakingConfig  `json:"staking"`
	Slashing EconomicModel0170SlashingConfig `json:"slashing"`
}

type EconomicModel0170CommonConfig struct {
	MaxGroupValidators uint32 `json:"maxGroupValidators"` // max validators count in 1 group
	CoordinatorsLimit  uint32 `json:"coordinatorLimit"`   // max Coordinators count in 1 group
	MaxConsensusVals   uint64 `json:"maxConsensusVals"`   // The consensus validators count
}

type EconomicModel0170StakingConfig struct {
	MaxValidators uint64 `json:"maxValidators"` // The epoch (billing cycle) validators count
}

type EconomicModel0170SlashingConfig struct {
	ZeroProduceCumulativeTime uint16 `json:"zeroProduceCumulativeTime"` // Count the number of zero-production blocks in this time range and check it. If it reaches a certain number of times, it can be punished (unit is consensus round)
}

// 主网升级0170时需要更新此可治理参数
func MaxValidators0170() uint64 {
	return ece.Extend0170.Staking.MaxValidators
}

// 主网升级0170时需要更新此可治理参数
func ZeroProduceCumulativeTime0170() uint16 {
	return ece.Extend0170.Slashing.ZeroProduceCumulativeTime
}
