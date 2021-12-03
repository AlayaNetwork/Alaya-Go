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
	MaxGroupValidators uint32 `json:"caxGroupValidators"` // max validators count in 1 group
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
