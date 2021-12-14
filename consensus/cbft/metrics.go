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

package cbft

import (
	"github.com/AlayaNetwork/Alaya-Go/metrics"
)

var (
	blockMinedGauage       = metrics.NewRegisteredGauge("cbft/gauage/block/mined", nil)
	viewChangedTimer       = metrics.NewRegisteredTimer("cbft/timer/view/changed", nil)
	blockQCCollectedGauage = metrics.NewRegisteredGauge("cbft/gauage/block/qc_collected", nil)

	blockProduceMeter          = metrics.NewRegisteredMeter("cbft/meter/block/produce", nil)
	blockCheckFailureMeter     = metrics.NewRegisteredMeter("cbft/meter/block/check_failure", nil)
	signatureCheckFailureMeter = metrics.NewRegisteredMeter("cbft/meter/signature/check_failure", nil)
	blockConfirmedMeter        = metrics.NewRegisteredMeter("cbft/meter/block/confirmed", nil)

	masterCounter    = metrics.NewRegisteredCounter("cbft/counter/view/count", nil)
	consensusCounter = metrics.NewRegisteredCounter("cbft/counter/consensus/count", nil)
	minedCounter     = metrics.NewRegisteredCounter("cbft/counter/mined/count", nil)

	viewNumberGauage          = metrics.NewRegisteredGauge("cbft/gauage/view/number", nil)
	epochNumberGauage         = metrics.NewRegisteredGauge("cbft/gauage/epoch/number", nil)
	proposerIndexGauage       = metrics.NewRegisteredGauge("cbft/gauage/proposer/index", nil)
	validatorCountGauage      = metrics.NewRegisteredGauge("cbft/gauage/validator/count", nil)
	blockNumberGauage         = metrics.NewRegisteredGauge("cbft/gauage/block/number", nil)
	highestQCNumberGauage     = metrics.NewRegisteredGauge("cbft/gauage/block/qc/number", nil)
	highestLockedNumberGauage = metrics.NewRegisteredGauge("cbft/gauage/block/locked/number", nil)
	highestCommitNumberGauage = metrics.NewRegisteredGauge("cbft/gauage/block/commit/number", nil)

	// for rand-grouped-consensus block
	upgradeCoordinatorBlockCounter = metrics.NewRegisteredCounter("cbft/counter/block/upgradeCoordinator/count", nil)

	blockGroupQCBySelfCounter  = metrics.NewRegisteredCounter("cbft/counter/block/groupqc/self/count", nil)  // Own group
	blockGroupQCByOtherCounter = metrics.NewRegisteredCounter("cbft/counter/block/groupqc/other/count", nil) // Own group

	blockWholeQCByVotesCounter   = metrics.NewRegisteredCounter("cbft/counter/block/wholeqc/votes/count", nil)
	blockWholeQCByRGQCCounter    = metrics.NewRegisteredCounter("cbft/counter/block/wholeqc/rgqc/count", nil)
	blockWholeQCByCombineCounter = metrics.NewRegisteredCounter("cbft/counter/block/wholeqc/combine/count", nil)

	blockGroupQCTimer = metrics.NewRegisteredTimer("cbft/timer/block/group/qc", nil) // Own group
	blockWholeQCTimer = metrics.NewRegisteredTimer("cbft/timer/block/whole/qc", nil)

	missRGBlockQuorumCertsGauage = metrics.NewRegisteredGauge("cbft/gauage/block/miss/rgqc", nil)
	missVotesGauage              = metrics.NewRegisteredGauge("cbft/gauage/block/miss/vote", nil)

	// for rand-grouped-consensus viewChange
	upgradeCoordinatorViewCounter = metrics.NewRegisteredCounter("cbft/counter/view/upgradeCoordinator/count", nil)

	viewGroupQCBySelfCounter  = metrics.NewRegisteredCounter("cbft/counter/view/groupqc/self/count", nil)
	viewGroupQCByOtherCounter = metrics.NewRegisteredCounter("cbft/counter/view/groupqc/other/count", nil) // Own group

	viewWholeQCByVcsCounter     = metrics.NewRegisteredCounter("cbft/counter/view/wholeqc/vcs/count", nil)
	viewWholeQCByRGQCCounter    = metrics.NewRegisteredCounter("cbft/counter/view/wholeqc/rgqc/count", nil)
	viewWholeQCByCombineCounter = metrics.NewRegisteredCounter("cbft/counter/view/wholeqc/combine/count", nil)

	viewGroupQCTimer = metrics.NewRegisteredTimer("cbft/timer/view/group/qc", nil) // Own group
	viewWholeQCTimer = metrics.NewRegisteredTimer("cbft/timer/view/whole/qc", nil)

	missRGViewQuorumCertsGauage = metrics.NewRegisteredGauge("cbft/gauage/view/miss/rgqc", nil)
	missVcsGauage               = metrics.NewRegisteredGauge("cbft/gauage/view/miss/vote", nil)
)
