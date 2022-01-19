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

package plugin

import (
	"context"
	"fmt"
	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/core/snapshotdb"
	"github.com/AlayaNetwork/Alaya-Go/core/state"
	"github.com/AlayaNetwork/Alaya-Go/core/types"
	"github.com/AlayaNetwork/Alaya-Go/rpc"
	"github.com/AlayaNetwork/Alaya-Go/x/gov"
	"github.com/AlayaNetwork/Alaya-Go/x/staking"
	"github.com/AlayaNetwork/Alaya-Go/x/xutil"
)

type BackendAPI interface {
	StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error)
}

// Provides an API interface to obtain data related to the economic model
type PublicPPOSAPI struct {
	snapshotDB snapshotdb.DB
	bkApi      BackendAPI
}

func NewPublicPPOSAPI(api BackendAPI) *PublicPPOSAPI {
	return &PublicPPOSAPI{snapshotDB: snapshotdb.Instance(), bkApi: api}
}

// Get node list of zero-out blocks
func (p *PublicPPOSAPI) GetWaitSlashingNodeList() string {
	list, err := slash.getWaitSlashingNodeList(0, common.ZeroHash)
	if nil != err || len(list) == 0 {
		return ""
	}
	return fmt.Sprintf("%+v", list)
}

// Get the list of consensus nodes for the current consensus cycle.
func (p *PublicPPOSAPI) GetConsensusNodeList(ctx context.Context) (staking.ValidatorExQueue, error) {
	state, header, err := p.bkApi.StateAndHeaderByNumber(ctx, rpc.LatestBlockNumber)
	if err != nil {
		return nil, err
	}
	blockHash := common.ZeroHash
	if !xutil.IsWorker(header.Extra) {
		blockHash = header.CacheHash()
	}
	return StakingInstance().GetValidatorList(blockHash, header.Number.Uint64(), CurrentRound, QueryStartNotIrr, gov.GetCurrentActiveVersion(state))
}

// Get the nodes in the current settlement cycle.
func (p *PublicPPOSAPI) GetValidatorList(ctx context.Context) (staking.ValidatorExQueue, error) {
	_, header, err := p.bkApi.StateAndHeaderByNumber(ctx, rpc.LatestBlockNumber)
	if err != nil {
		return nil, err
	}
	blockHash := common.ZeroHash
	if !xutil.IsWorker(header.Extra) {
		blockHash = header.CacheHash()
	}
	return StakingInstance().GetVerifierList(blockHash, header.Number.Uint64(), QueryStartNotIrr)
}

// Get all staking nodes.
func (p *PublicPPOSAPI) GetCandidateList(ctx context.Context) (staking.CandidateHexQueue, error) {
	state, header, err := p.bkApi.StateAndHeaderByNumber(ctx, rpc.LatestBlockNumber)
	if err != nil {
		return nil, err
	}
	blockHash := common.ZeroHash
	if !xutil.IsWorker(header.Extra) {
		blockHash = header.CacheHash()
	}
	return StakingInstance().GetCandidateList(blockHash, header.Number.Uint64(), state)
}
