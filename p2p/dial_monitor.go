package p2p

import (
	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/x/staking"
)

func ConvertToCommonNodeIdList(verifierList []*staking.Validator) []common.NodeID {
	nodeIdList := make([]common.NodeID, len(verifierList))
	for i, verifier := range verifierList {
		nodeIdList[i] = common.NodeID(verifier.NodeId)
	}
	return nodeIdList
}
