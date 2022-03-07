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
	blockMinedGauge  = metrics.NewRegisteredGauge("cbft/gauge/block/mined", nil)
	viewChangedTimer = metrics.NewRegisteredTimer("cbft/timer/view/changed", nil)

	blockProduceMeter          = metrics.NewRegisteredMeter("cbft/meter/block/produce", nil)
	blockCheckFailureMeter     = metrics.NewRegisteredMeter("cbft/meter/block/check_failure", nil)
	signatureCheckFailureMeter = metrics.NewRegisteredMeter("cbft/meter/signature/check_failure", nil)

	masterCounter    = metrics.NewRegisteredCounter("cbft/counter/view/count", nil)
	consensusCounter = metrics.NewRegisteredCounter("cbft/counter/consensus/count", nil)
	minedCounter     = metrics.NewRegisteredCounter("cbft/counter/mined/count", nil)

	viewNumberGauge     = metrics.NewRegisteredGauge("cbft/gauge/view/number", nil)
	epochNumberGauge    = metrics.NewRegisteredGauge("cbft/gauge/epoch/number", nil)
	proposerIndexGauge  = metrics.NewRegisteredGauge("cbft/gauge/proposer/index", nil)
	validatorCountGauge = metrics.NewRegisteredGauge("cbft/gauge/validator/count", nil)

	// for rand-grouped-consensus block
	upgradeCoordinatorBlockCounter = metrics.NewRegisteredCounter("cbft/counter/block/upgradeCoordinator/count", nil)

	blockGroupQCBySelfCounter  = metrics.NewRegisteredCounter("cbft/counter/block/groupqc/self/count", nil)  // Own group
	blockGroupQCByOtherCounter = metrics.NewRegisteredCounter("cbft/counter/block/groupqc/other/count", nil) // Own group

	blockWholeQCByVotesCounter   = metrics.NewRegisteredCounter("cbft/counter/block/wholeqc/votes/count", nil)
	blockWholeQCByRGQCCounter    = metrics.NewRegisteredCounter("cbft/counter/block/wholeqc/rgqc/count", nil)
	blockWholeQCByCombineCounter = metrics.NewRegisteredCounter("cbft/counter/block/wholeqc/combine/count", nil)

	blockGroupQCTimer = metrics.NewRegisteredTimer("cbft/timer/block/group/qc", nil) // Own group
	blockWholeQCTimer = metrics.NewRegisteredTimer("cbft/timer/block/whole/qc", nil)

	missRGBlockQuorumCertsGauge = metrics.NewRegisteredGauge("cbft/gauge/block/miss/rgqc", nil)
	missVotesGauge              = metrics.NewRegisteredGauge("cbft/gauge/block/miss/vote", nil)

	missVotesCounter     = metrics.NewRegisteredCounter("cbft/counter/block/miss/vote/count", nil)
	responseVotesCounter = metrics.NewRegisteredCounter("cbft/counter/block/response/vote/count", nil)

	// for rand-grouped-consensus viewChange
	upgradeCoordinatorViewCounter = metrics.NewRegisteredCounter("cbft/counter/view/upgradeCoordinator/count", nil)

	viewGroupQCBySelfCounter  = metrics.NewRegisteredCounter("cbft/counter/view/groupqc/self/count", nil)
	viewGroupQCByOtherCounter = metrics.NewRegisteredCounter("cbft/counter/view/groupqc/other/count", nil) // Own group

	viewWholeQCByVcsCounter     = metrics.NewRegisteredCounter("cbft/counter/view/wholeqc/vcs/count", nil)
	viewWholeQCByRGQCCounter    = metrics.NewRegisteredCounter("cbft/counter/view/wholeqc/rgqc/count", nil)
	viewWholeQCByCombineCounter = metrics.NewRegisteredCounter("cbft/counter/view/wholeqc/combine/count", nil)

	viewGroupQCTimer = metrics.NewRegisteredTimer("cbft/timer/view/group/qc", nil) // Own group
	viewWholeQCTimer = metrics.NewRegisteredTimer("cbft/timer/view/whole/qc", nil)

	missRGViewQuorumCertsGauge = metrics.NewRegisteredGauge("cbft/gauge/view/miss/rgqc", nil)
	missVcsGauge               = metrics.NewRegisteredGauge("cbft/gauge/view/miss/vote", nil)

	missVcsCounter     = metrics.NewRegisteredCounter("cbft/counter/view/miss/vcs/count", nil)
	responseVcsCounter = metrics.NewRegisteredCounter("cbft/counter/view/response/vcs/count", nil)
)
