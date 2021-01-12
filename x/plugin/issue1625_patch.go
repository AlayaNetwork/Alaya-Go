// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package plugin

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common/vm"

	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/PlatONnetwork/PlatON-Go/x/staking"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

func NewFixIssue1625Plugin(sdb snapshotdb.DB) *FixIssue1625Plugin {
	fix := new(FixIssue1625Plugin)
	fix.sdb = sdb
	return fix
}

type FixIssue1625Plugin struct {
	sdb snapshotdb.DB
}

// 有挪用锁仓资金的账户
// 因为restrictInfo.CachePlanAmount计算错误，造成用户可以挪用锁仓合约中的资金。
// 因此，首先需要修复restrictInfo.CachePlanAmount的问题；还需要修复利用挪用资金，做了质押、或者委托的，也要修复。（用挪用的资金做了委托的，仍然可以取得委托收益）
func (a *FixIssue1625Plugin) fix(blockHash common.Hash, head *types.Header, state xcom.StateDB, chainID *big.Int) error {
	if chainID.Cmp(params.AlayaChainConfig.ChainID) != 0 {
		return nil
	}
	issue1625, err := NewIssue1625Accounts()
	if err != nil {
		return err
	}
	for _, issue1625Account := range issue1625 {
		//获取释放给当前账户的锁仓信息
		restrictingKey, restrictInfo, err := rt.mustGetRestrictingInfoByDecode(state, issue1625Account.addr)
		if err != nil {
			return err
		}
		log.Debug("fix issue 1625 begin", "account", issue1625Account.addr, "fix amount", issue1625Account.amount, "info", restrictInfo)

		// 锁仓实际可用金额 = 锁仓当前可用金额（包含挪用总金额）- 挪用总金额
		actualRestrictingAmount := new(big.Int).Sub(restrictInfo.CachePlanAmount, issue1625Account.amount)
		if actualRestrictingAmount.Cmp(common.Big0) < 0 {
			log.Error("seems not good here", "info", restrictInfo, "amount", issue1625Account.amount, "account", issue1625Account.addr)
			return fmt.Errorf("the account restrictInfo seems not right")
		}

		//实际用于（质押+委托）的挪用金额 = 锁仓当前用于（质押+委托）的金额 - 锁仓实际可用金额
		wrongStakingAmount := new(big.Int).Sub(restrictInfo.StakingAmount, actualRestrictingAmount)
		if wrongStakingAmount.Cmp(common.Big0) > 0 {
			//用于（质押+委托）的锁仓金额，必须 <= 锁仓实际可用金额，现在是>了
			//说明有挪用金额用于质押或者委托
			//If the user uses the wrong amount,Roll back the unused part first
			//优先回滚没有使用的那部分锁仓余额
			//则首先把没有用于（质押+委托）的挪用金额，直接从锁仓当前可用金额中扣除。
			wrongNoUseAmount := new(big.Int).Sub(issue1625Account.amount, wrongStakingAmount)
			if wrongNoUseAmount.Cmp(common.Big0) > 0 {
				restrictInfo.CachePlanAmount.Sub(restrictInfo.CachePlanAmount, wrongNoUseAmount)
				rt.storeRestrictingInfo(state, restrictingKey, restrictInfo)
				log.Debug("fix issue 1625  at no use", "no use", wrongNoUseAmount)
			}

			//roll back del,回滚委托
			if err := a.rollBackDel(blockHash, head.Number, issue1625Account.addr, wrongStakingAmount, state); err != nil {
				return err
			}
			//roll back staking,回滚质押
			if wrongStakingAmount.Cmp(common.Big0) > 0 {
				if err := a.rollBackStaking(blockHash, head.Number, issue1625Account.addr, wrongStakingAmount, state); err != nil {
					return err
				}
			}
		} else {
			//当用户没有使用因为漏洞产生的钱，直接减去漏洞的钱就是正确的余额
			restrictInfo.CachePlanAmount.Sub(restrictInfo.CachePlanAmount, issue1625Account.amount)
			if restrictInfo.StakingAmount.Cmp(common.Big0) == 0 &&
				len(restrictInfo.ReleaseList) == 0 && restrictInfo.CachePlanAmount.Cmp(common.Big0) == 0 {
				state.SetState(vm.RestrictingContractAddr, restrictingKey, []byte{})
				log.Debug("fix issue 1625 finished,set info empty", "account", issue1625Account.addr, "fix amount", issue1625Account.amount)
			} else {
				rt.storeRestrictingInfo(state, restrictingKey, restrictInfo)
				log.Debug("fix issue 1625 finished", "account", issue1625Account.addr, "info", restrictInfo, "fix amount", issue1625Account.amount)
			}
		}
	}
	return nil
}

// 委托回滚处理：
//1.委托的节点已经完全退出,并且委托的时间靠后
//2.委托的节点处于解质押状态,并且委托的时间靠后
//3.根据委托节点的分红比例从小到大排序，如果委托比例相同，根据节点id从小到大排序
// param: hash			区块hash
// param: blockNumber 	区块高度
// param: account 		账户
// param: amount 		实际用来委托的挪用的资金（<=统计出的挪用总金额）
// param: state
func (a *FixIssue1625Plugin) rollBackDel(hash common.Hash, blockNumber *big.Int, account common.Address, amount *big.Int, state xcom.StateDB) error {

	delAddrByte := account.Bytes()

	markPre := len(staking.DelegateKeyPrefix)
	markDelAddr := markPre + len(delAddrByte)

	key := make([]byte, markDelAddr)
	copy(key[:markPre], staking.DelegateKeyPrefix)
	copy(key[markPre:markDelAddr], delAddrByte)

	iter := a.sdb.Ranking(hash, key, 0)
	if err := iter.Error(); nil != err {
		return err
	}
	defer iter.Release()

	var dels issue1625AccountDelInfos

	//遍历用户的所有委托
	for iter.Valid(); iter.Next(); {
		var del staking.Delegation
		if err := rlp.DecodeBytes(iter.Value(), &del); nil != err {
			return err
		}
		_, nodeID, stakingBlock := staking.DecodeDelegateKey(iter.Key())
		canAddr, err := xutil.NodeId2Addr(nodeID)
		if nil != err {
			return err
		}
		can, err := stk.db.GetCandidateStore(hash, canAddr)
		if snapshotdb.NonDbNotFoundErr(err) {
			return err
		}

		delInfo := new(issue1625AccountDelInfo)
		delInfo.del = &del
		if can.IsNotEmpty() {
			//确保节点是委托的时的那个节点
			if can.StakingBlockNum == stakingBlock {
				delInfo.candidate = can
			}
		} else {
			delInfo.candidate = can
		}
		delInfo.nodeID = nodeID
		delInfo.stakingBlock = stakingBlock
		//来自锁仓的委托金额
		delInfo.originRestrictingAmount = new(big.Int).Add(del.RestrictingPlan, del.RestrictingPlanHes)
		//来自用户钱包的委托金额
		delInfo.originFreeAmount = new(big.Int).Add(del.Released, del.ReleasedHes)
		delInfo.canAddr = canAddr
		//如果该委托没有用锁仓，无需回滚
		if delInfo.originRestrictingAmount.Cmp(common.Big0) == 0 {
			continue
		}
		dels = append(dels, delInfo)
	}
	sort.Sort(dels)
	epoch := xutil.CalculateEpoch(blockNumber.Uint64())
	stakingdb := staking.NewStakingDBWithDB(a.sdb)
	for i := 0; i < len(dels); i++ {
		if err := dels[i].handleDelegate(hash, blockNumber, epoch, account, amount, state, stakingdb); err != nil {
			return err
		}
		if amount.Cmp(common.Big0) <= 0 {
			break
		}
	}
	return nil
}

//
//
// param: hash
// param: blockNumber
// param: account
// param: amount	实际用来质押的挪用的资金（<=统计出的挪用总金额）
// param: state
func (a *FixIssue1625Plugin) rollBackStaking(hash common.Hash, blockNumber *big.Int, account common.Address, amount *big.Int, state xcom.StateDB) error {

	iter := a.sdb.Ranking(hash, staking.CanBaseKeyPrefix, 0)
	if err := iter.Error(); nil != err {
		return err
	}
	defer iter.Release()

	stakingdb := staking.NewStakingDBWithDB(a.sdb)

	var stakings issue1625AccountStakingInfos
	for iter.Valid(); iter.Next(); {
		var canbase staking.CandidateBase
		if err := rlp.DecodeBytes(iter.Value(), &canbase); nil != err {
			return err
		}

		if canbase.StakingAddress == account {
			canAddr, err := xutil.NodeId2Addr(canbase.NodeId)
			if nil != err {
				return err
			}
			canmu, err := stakingdb.GetCanMutableStore(hash, canAddr)
			if nil != err {
				return err
			}
			candidate := staking.Candidate{
				&canbase, canmu,
			}
			//如果该质押没有用锁仓，无需回滚
			if candidate.IsNotEmpty() {
				restrictingAmount := new(big.Int).Add(candidate.RestrictingPlan, candidate.RestrictingPlanHes)
				if restrictingAmount.Cmp(common.Big0) == 0 {
					continue
				}
			}
			stakings = append(stakings, newIssue1625AccountStakingInfo(&candidate, canAddr))
		}
	}
	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	sort.Sort(stakings)

	//遍历质押信息
	for i := 0; i < len(stakings); i++ {
		if err := stakings[i].handleStaking(hash, blockNumber, epoch, amount, state, stakingdb); err != nil {
			return err
		}
		if amount.Cmp(common.Big0) <= 0 {
			break
		}
	}
	return nil
}

func newIssue1625AccountStakingInfo(candidate *staking.Candidate, canAddr common.NodeAddress) *issue1625AccountStakingInfo {
	w := &issue1625AccountStakingInfo{
		candidate: candidate,
		canAddr:   canAddr,
	}
	w.originRestrictingAmount = new(big.Int).Add(w.candidate.RestrictingPlan, w.candidate.RestrictingPlanHes)
	w.originFreeAmount = new(big.Int).Add(w.candidate.Released, w.candidate.ReleasedHes)
	return w
}

type issue1625AccountStakingInfo struct {
	candidate                                 *staking.Candidate
	canAddr                                   common.NodeAddress
	originRestrictingAmount, originFreeAmount *big.Int
}

//回退处于退出期的质押信息
//
//
// param: hash
// param: epoch
// param: rollBackAmount	实际用于（质押+委托）的挪用金额
// param: state
// param: stdb
func (a *issue1625AccountStakingInfo) handleExistStaking(hash common.Hash, blockNumber *big.Int, epoch uint64, rollBackAmount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) error {
	//stats
	fixStaking := new(common.FixStaking)
	fixStaking.NodeID = common.NodeID(a.candidate.NodeId)
	fixStaking.StakingBlockNumber = a.candidate.StakingBlockNum

	if a.originRestrictingAmount.Cmp(rollBackAmount) >= 0 {
		if err := rt.ReturnWrongLockFunds(a.candidate.StakingAddress, rollBackAmount, state); nil != err {
			return err
		}
		//a.candidate.RestrictingPlan 中可能包含多个锁仓计划的资金，所以，a.candidate.RestrictingPlan >= rollBackAmount
		a.candidate.RestrictingPlan = new(big.Int).Sub(a.candidate.RestrictingPlan, rollBackAmount)
		rollBackAmount.SetInt64(0)

		//stats
		fixStaking.ImproperValidRestrictingAmount = rollBackAmount

	} else {
		if err := rt.ReturnWrongLockFunds(a.candidate.StakingAddress, a.originRestrictingAmount, state); nil != err {
			return err
		}

		//stats
		// a.originRestrictingAmount == a.candidate.RestrictingPlan()
		fixStaking.ImproperValidRestrictingAmount = a.originRestrictingAmount

		a.candidate.RestrictingPlan = new(big.Int).SetInt64(0)
		a.candidate.RestrictingPlanHes = new(big.Int).SetInt64(0)
		rollBackAmount.Sub(rollBackAmount, a.originRestrictingAmount)
	}

	//stats
	//由于处于锁定期，所以只需要先调整当前质押的挪用的锁仓金额，从质押合约->锁仓合约，然后等待冻结期结束即可;
	//跟踪系统也是如此调整
	fixStaking.FurtherOperation = "NOP"
	common.CollectFixStaking(blockNumber.Uint64(), a.candidate.StakingAddress, fixStaking)

	a.candidate.StakingEpoch = uint32(epoch)

	if err := stdb.SetCanMutableStore(hash, a.canAddr, a.candidate.CandidateMutable); nil != err {
		return err
	}

	return nil
}

//检查是否达到质押门槛
func (a *issue1625AccountStakingInfo) shouldWithdrewStaking(hash common.Hash, blockNumber *big.Int, rollBackAmount *big.Int) bool {
	left := new(big.Int)
	if a.originRestrictingAmount.Cmp(rollBackAmount) >= 0 {
		left = new(big.Int).Add(a.originFreeAmount, new(big.Int).Sub(a.originRestrictingAmount, rollBackAmount))
	} else {
		left = new(big.Int).Set(a.originFreeAmount)
	}
	if ok, _ := CheckStakeThreshold(blockNumber.Uint64(), hash, left); !ok {
		return true
	}

	return false
}

func (a *issue1625AccountStakingInfo) calImproperRestrictingAmount(rollBackAmount *big.Int) *big.Int {
	//计算此次需要回退的钱
	improperRestrictingAmount := new(big.Int)
	if rollBackAmount.Cmp(a.originRestrictingAmount) >= 0 {
		improperRestrictingAmount = new(big.Int).Set(a.originRestrictingAmount)
	} else {
		improperRestrictingAmount = new(big.Int).Set(rollBackAmount)
	}
	return improperRestrictingAmount
}

func (a *issue1625AccountStakingInfo) fixCandidateInfo(improperRestrictingAmount *big.Int, fixStaking *common.FixStaking) {
	//修正质押信息
	if a.candidate.RestrictingPlanHes.Cmp(improperRestrictingAmount) >= 0 {
		a.candidate.RestrictingPlanHes.Sub(a.candidate.RestrictingPlanHes, improperRestrictingAmount)

		//stats, 记录犹豫期的锁仓金额，应该退回的数量
		fixStaking.ImproperHesitatingRestrictingAmount = improperRestrictingAmount
	} else {
		hes := new(big.Int).Set(a.candidate.RestrictingPlanHes)
		a.candidate.RestrictingPlanHes = new(big.Int)
		a.candidate.RestrictingPlan = new(big.Int).Sub(a.candidate.RestrictingPlan, new(big.Int).Sub(improperRestrictingAmount, hes))

		//stats, 记录犹豫期的，以及有效锁仓金额，应该退回的数量
		fixStaking.ImproperHesitatingRestrictingAmount = hes
		fixStaking.ImproperValidRestrictingAmount = new(big.Int).Sub(improperRestrictingAmount, hes)
	}
	a.candidate.SubShares(improperRestrictingAmount)
}

//减持质押
func (a *issue1625AccountStakingInfo) decreaseStaking(hash common.Hash, blockNumber *big.Int, epoch uint64, rollBackAmount *big.Int, state xcom.StateDB) error {
	if err := stk.db.DelCanPowerStore(hash, a.candidate); nil != err {
		return err
	}

	//stats
	//要先调整当前质押的挪用的锁仓金额，从质押合约->锁仓合约，然后走减持质押流程，这是个新的流程。
	//跟踪系统也是如此调整
	fixStaking := new(common.FixStaking)
	fixStaking.NodeID = common.NodeID(a.candidate.NodeId)
	fixStaking.StakingBlockNumber = a.candidate.StakingBlockNum
	fixStaking.FurtherOperation = "REDUCE"

	lazyCalcStakeAmount(epoch, a.candidate.CandidateMutable)

	//计算此次需要回退的钱
	improperRestrictingAmount := a.calImproperRestrictingAmount(rollBackAmount)

	//回退因为漏洞产生的金额
	if err := rt.ReturnWrongLockFunds(a.candidate.StakingAddress, improperRestrictingAmount, state); nil != err {
		return err
	}

	//修正质押信息
	//a.fixCandidateInfo(improperRestrictingAmount)
	a.fixCandidateInfo(improperRestrictingAmount, fixStaking)

	//stats
	common.CollectFixStaking(blockNumber.Uint64(), a.candidate.StakingAddress, fixStaking)

	a.candidate.StakingEpoch = uint32(epoch)

	if err := stk.db.SetCanPowerStore(hash, a.canAddr, a.candidate); nil != err {
		return err
	}
	if err := stk.db.SetCanMutableStore(hash, a.canAddr, a.candidate.CandidateMutable); nil != err {
		return err
	}

	rollBackAmount.Sub(rollBackAmount, improperRestrictingAmount)
	return nil
}

//撤销质押
func (a *issue1625AccountStakingInfo) withdrewStaking(hash common.Hash, epoch uint64, blockNumber *big.Int, rollBackAmount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) error {
	if err := stdb.DelCanPowerStore(hash, a.candidate); nil != err {
		return err
	}

	//stats
	//要先调整当前质押的挪用的锁仓金额，从质押合约->锁仓合约，然后走解除质押流程，并锁定一个结算周期。
	//跟踪系统也是如此调整
	fixStaking := new(common.FixStaking)
	fixStaking.NodeID = common.NodeID(a.candidate.NodeId)
	fixStaking.StakingBlockNumber = a.candidate.StakingBlockNum
	fixStaking.FurtherOperation = "WITHDRAW"

	//计算此次需要回退的钱
	improperRestrictingAmount := a.calImproperRestrictingAmount(rollBackAmount)

	//回退因为漏洞产生的金额
	if err := rt.ReturnWrongLockFunds(a.candidate.StakingAddress, improperRestrictingAmount, state); nil != err {
		return err
	}

	//修正质押信息
	//a.fixCandidateInfo(improperRestrictingAmount)
	//stats
	a.fixCandidateInfo(improperRestrictingAmount, fixStaking)

	//stats
	common.CollectFixStaking(blockNumber.Uint64(), a.candidate.StakingAddress, fixStaking)

	//开始解质押
	//回退犹豫期的自由金
	if a.candidate.ReleasedHes.Cmp(common.Big0) > 0 {
		rt.transferAmount(state, vm.StakingContractAddr, a.candidate.StakingAddress, a.candidate.ReleasedHes)
		a.candidate.ReleasedHes = new(big.Int).SetInt64(0)
	}

	//回退犹豫期的锁仓
	if a.candidate.RestrictingPlanHes.Cmp(common.Big0) > 0 {
		err := rt.ReturnLockFunds(a.candidate.StakingAddress, a.candidate.RestrictingPlanHes, state)
		if nil != err {
			return err
		}
		a.candidate.RestrictingPlanHes = new(big.Int).SetInt64(0)
	}

	a.candidate.CleanShares()
	a.candidate.Status |= staking.Invalided | staking.Withdrew

	a.candidate.StakingEpoch = uint32(epoch)

	if a.candidate.Released.Cmp(common.Big0) > 0 || a.candidate.RestrictingPlan.Cmp(common.Big0) > 0 {
		//如果质押处于生效期，需要锁定
		if err := stk.addErrorAccountUnStakeItem(blockNumber.Uint64(), hash, a.candidate.NodeId, a.canAddr, a.candidate.StakingBlockNum); nil != err {
			return err
		}
		// sub the account staking Reference Count
		if err := stdb.SubAccountStakeRc(hash, a.candidate.StakingAddress); nil != err {
			return err
		}
		if err := stdb.SetCanMutableStore(hash, a.canAddr, a.candidate.CandidateMutable); nil != err {
			return err
		}
	} else {
		//如果质押还处于犹豫期，不用锁定
		if err := stdb.DelCandidateStore(hash, a.canAddr); nil != err {
			return err
		}
	}

	rollBackAmount.Sub(rollBackAmount, improperRestrictingAmount)
	return nil
}

//
//
// param: hash
// param: blockNumber
// param: epoch
// param: rollBackAmount	实际用于（质押+委托）的挪用金额
// param: state
// param: stdb
func (a *issue1625AccountStakingInfo) handleStaking(hash common.Hash, blockNumber *big.Int, epoch uint64, rollBackAmount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) error {
	lazyCalcStakeAmount(epoch, a.candidate.CandidateMutable)
	log.Debug("fix issue 1625 for staking begin", "account", a.candidate.StakingAddress, "nodeID", a.candidate.NodeId.TerminalString(), "return", rollBackAmount, "restrictingPlan",
		a.candidate.RestrictingPlan, "restrictingPlanRes", a.candidate.RestrictingPlanHes, "released", a.candidate.Released, "releasedHes", a.candidate.ReleasedHes, "share", a.candidate.Shares)
	if a.candidate.Status.IsWithdrew() {
		//已经解质押,节点处于退出锁定期
		if err := a.handleExistStaking(hash, blockNumber, epoch, rollBackAmount, state, stdb); err != nil {
			return err
		}
		log.Debug("fix issue 1625 for staking end", "account", a.candidate.StakingAddress, "nodeID", a.candidate.NodeId.TerminalString(), "status", a.candidate.Status, "return",
			rollBackAmount, "restrictingPlan", a.candidate.RestrictingPlan, "restrictingPlanRes", a.candidate.RestrictingPlanHes, "released", a.candidate.Released, "releasedHes", a.candidate.ReleasedHes, "share", a.candidate.Shares)
	} else {
		//如果质押中,根据回退后的剩余质押金额是否达到质押门槛来判断时候需要撤销或者减持质押
		shouldWithdrewStaking := a.shouldWithdrewStaking(hash, blockNumber, rollBackAmount)
		if shouldWithdrewStaking {
			//撤销质押
			if err := a.withdrewStaking(hash, epoch, blockNumber, rollBackAmount, state, stdb); err != nil {
				return err
			}
		} else {
			//减持质押
			if err := a.decreaseStaking(hash, blockNumber, epoch, rollBackAmount, state); err != nil {
				return err
			}
		}
		log.Debug("fix issue 1625 for staking end", "account", a.candidate.StakingAddress, "nodeID", a.candidate.NodeId.TerminalString(), "status", a.candidate.Status, "withdrewStaking",
			shouldWithdrewStaking, "return", rollBackAmount, "restrictingPlan", a.candidate.RestrictingPlan, "restrictingPlanRes", a.candidate.RestrictingPlanHes, "released", a.candidate.Released,
			"releasedHes", a.candidate.ReleasedHes, "share", a.candidate.Shares)
	}
	return nil
}

type issue1625AccountStakingInfos []*issue1625AccountStakingInfo

func (d issue1625AccountStakingInfos) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d issue1625AccountStakingInfos) Len() int      { return len(d) }
func (d issue1625AccountStakingInfos) Less(i, j int) bool {
	//1.节点处于解质押状态,并且质押的时间靠后的排在前面
	//2.按照节点质押的时间远近进行排序，并且质押的时间靠后的排在前面
	if d[i].candidate.IsWithdrew() {
		if d[j].candidate.IsWithdrew() {
			return d.LessByStakingBlockNum(i, j)
		} else {
			return false
		}
	} else {
		if d[j].candidate.IsWithdrew() {
			return false
		} else {
			return d.LessByStakingBlockNum(i, j)
		}
	}
}

func (d issue1625AccountStakingInfos) LessByStakingBlockNum(i, j int) bool {
	if d[i].candidate.StakingBlockNum > d[j].candidate.StakingBlockNum {
		return true
	} else {
		return false
	}
}

type issue1625AccountDelInfo struct {
	del       *staking.Delegation //委托信息
	candidate *staking.Candidate  //节点信息
	//use for get staking
	stakingBlock uint64             //节点质押快高
	canAddr      common.NodeAddress //节点地址
	nodeID       discover.NodeID

	originRestrictingAmount, originFreeAmount *big.Int
}

//
// 检查当前委托，再退回挪用金额后，剩余委托金额是否满足最低要求。如果不满足最低阈值，则需要撤回此委托。
// param: hash				区块hash
// param: blockNumber		区块高度
// param: rollBackAmount	挪用的锁仓总金额
// return:
func (a *issue1625AccountDelInfo) shouldWithdrewDel(hash common.Hash, blockNumber *big.Int, rollBackAmount *big.Int) bool {
	leftTotalDelgateAmount := new(big.Int)
	if rollBackAmount.Cmp(a.originRestrictingAmount) >= 0 {
		leftTotalDelgateAmount.Set(a.originFreeAmount)
	} else {
		leftTotalDelgateAmount.Add(a.originFreeAmount, new(big.Int).Sub(a.originRestrictingAmount, rollBackAmount))
	}
	if ok, _ := CheckOperatingThreshold(blockNumber.Uint64(), hash, leftTotalDelgateAmount); ok {
		return false
	}
	return true
}

//
// 处理委托用户在一个节点上的委托
// param: hash				区块hash
// param: blockNumber		区块高度
// param: epoch				当前epoch
// param: delAddr			委托用户地址
// param: rollBackAmount	用于（质押+委托）的挪用总金额
// param: state
// param: stdb
// return:
func (a *issue1625AccountDelInfo) handleDelegate(hash common.Hash, blockNumber *big.Int, epoch uint64, delAddr common.Address, rollBackAmount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) error {
	//当前委托中，挪用的锁仓金额
	improperRestrictingAmount := new(big.Int)

	if rollBackAmount.Cmp(a.originRestrictingAmount) >= 0 {
		//如果挪用总金额 >= 委托中来自锁仓的金额，则认为来自锁仓的金额都是挪用金额
		improperRestrictingAmount = new(big.Int).Set(a.originRestrictingAmount)
	} else {
		//如果挪用总金额 < 委托中来自锁仓的金额，则认为来自锁仓的金额，只有部分是来自挪用金额
		improperRestrictingAmount = new(big.Int).Set(rollBackAmount)
	}
	log.Debug("fix issue 1625 for delegate begin", "account", delAddr, "candidate", a.nodeID.String(), "currentReturn", improperRestrictingAmount, "leftReturn", rollBackAmount, "restrictingPlan", a.del.RestrictingPlan, "restrictingPlanRes", a.del.RestrictingPlanHes,
		"release", a.del.Released, "releaseHes", a.del.ReleasedHes, "CumulativeIncome", a.del.CumulativeIncome)
	if a.candidate.IsNotEmpty() {
		log.Debug("fix issue 1625 for delegate ,can begin info", "account", delAddr, "candidate", a.nodeID.String(), "share", a.candidate.Shares, "candidate.del", a.candidate.DelegateTotal, "candidate.delhes", a.candidate.DelegateTotalHes, "canValid", a.candidate.IsValid())
	}

	//stats
	//fix委托，构造调整记录
	fixDelegation := new(common.FixDelegation)
	fixDelegation.NodeID = common.NodeID(a.nodeID)
	fixDelegation.StakingBlockNumber = a.stakingBlock

	//先计算委托收益
	delegateRewardPerList, err := RewardMgrInstance().GetDelegateRewardPerList(hash, a.nodeID, a.stakingBlock, uint64(a.del.DelegateEpoch), xutil.CalculateEpoch(blockNumber.Uint64())-1)
	if snapshotdb.NonDbNotFoundErr(err) {
		return err
	}

	rewardsReceive := calcDelegateIncome(epoch, a.del, delegateRewardPerList)
	if err := UpdateDelegateRewardPer(hash, a.nodeID, a.stakingBlock, rewardsReceive, rm.db); err != nil {
		return err
	}
	if a.candidate.IsNotEmpty() {
		lazyCalcNodeTotalDelegateAmount(epoch, a.candidate.CandidateMutable)
	}

	a.del.DelegateEpoch = uint32(epoch)

	if a.candidate.IsNotEmpty() && a.candidate.IsValid() {
		if err := stdb.DelCanPowerStore(hash, a.candidate); nil != err {
			return err
		}
	}
	//检查是否需要撤消当前委托
	withdrewDel := a.shouldWithdrewDel(hash, blockNumber, rollBackAmount)
	if withdrewDel {
		//回滚错误金额
		if err := a.fixImproperRestrictingAmountByDel(delAddr, improperRestrictingAmount, state, fixDelegation); err != nil {
			return err
		}

		//开始撤销委托
		//更新candidate中由于撤销委托导致的记录变动
		if a.candidate.IsNotEmpty() {
			hes := new(big.Int).Add(a.del.ReleasedHes, a.del.RestrictingPlanHes)
			lock := new(big.Int).Add(a.del.Released, a.del.RestrictingPlan)
			a.candidate.DelegateTotalHes.Sub(a.candidate.DelegateTotalHes, hes)
			a.candidate.DelegateTotal.Sub(a.candidate.DelegateTotal, lock)
			if a.candidate.Shares.Cmp(new(big.Int).Add(hes, lock)) >= 0 {
				a.candidate.Shares.Sub(a.candidate.Shares, new(big.Int).Add(hes, lock))
			}
		}

		//回退锁仓
		if err := rt.ReturnLockFunds(delAddr, new(big.Int).Add(a.del.RestrictingPlan, a.del.RestrictingPlanHes), state); err != nil {
			return err
		}

		//回退自由
		rt.transferAmount(state, vm.StakingContractAddr, delAddr, new(big.Int).Add(a.del.ReleasedHes, a.del.Released))

		//领取收益
		if err := rm.ReturnDelegateReward(delAddr, a.del.CumulativeIncome, state); err != nil {
			return common.InternalError
		}

		//stats
		//用户从某节点全部撤消委托后，记录用户的委托收益
		//撤消委托，构造调整记录，记录用户的委托收益
		fixDelegation.Withdraw = true
		fixDelegation.RewardAmount = a.del.CumulativeIncome
		//删除委托
		if err := stdb.DelDelegateStore(hash, delAddr, a.nodeID, a.stakingBlock); nil != err {
			return err
		}

		log.Debug("fix issue 1625 for delegate,withdrew del", "account", delAddr, "candidate", a.nodeID.String(), "income", a.del.CumulativeIncome)
	} else {
		//不需要解除委托
		if err := a.fixImproperRestrictingAmountByDel(delAddr, improperRestrictingAmount, state, fixDelegation); err != nil {
			return err
		}
		if err := stdb.SetDelegateStore(hash, delAddr, a.nodeID, a.stakingBlock, a.del); nil != err {
			return err
		}
		//stats
		//撤消委托，构造调整记录
		fixDelegation.Withdraw = false

		log.Debug("fix issue 1625 for delegate,decrease del", "account", delAddr, "candidate", a.nodeID.String(), "restrictingPlan", a.del.RestrictingPlan, "restrictingPlanRes", a.del.RestrictingPlanHes,
			"release", a.del.Released, "releaseHes", a.del.ReleasedHes, "income", a.del.CumulativeIncome)
	}

	//stats
	common.CollectFixDelegation(blockNumber.Uint64(), common.Address(a.canAddr), fixDelegation)

	if a.candidate.IsNotEmpty() {
		if a.candidate.IsValid() {
			if err := stdb.SetCanPowerStore(hash, a.canAddr, a.candidate); nil != err {
				return err
			}
		}
		if err := stdb.SetCanMutableStore(hash, a.canAddr, a.candidate.CandidateMutable); nil != err {
			return err
		}
	}

	rollBackAmount.Sub(rollBackAmount, improperRestrictingAmount)

	if !a.candidate.IsEmpty() {
		log.Debug("fix issue 1625 for delegate,can last info", "account", delAddr, "candidate", a.nodeID.String(), "share", a.candidate.Shares, "candidate.del", a.candidate.DelegateTotal, "candidate.delhes", a.candidate.DelegateTotalHes)
	}
	log.Debug("fix issue 1625 for delegate end", "account", delAddr, "candidate", a.nodeID.String(), "leftReturn", rollBackAmount, "withdrewDel", withdrewDel)
	return nil
}

//修正委托以及验证人的锁仓信息
//
// 退回前委托的挪用的锁仓金额（挪用的锁仓资金，只是体现在委托信息的RestrictingPlanHes/RestrictingPlan字段中。）
// 用户的委托信息修改后，接受委托的节点信息，也需要修改接受的委托金额。
// param: delAddr					委托用户地址
// param: improperRestrictingAmount	当前委托的挪用的锁仓金额
// param: state
func (a *issue1625AccountDelInfo) fixImproperRestrictingAmountByDel(delAddr common.Address, improperRestrictingAmount *big.Int, state xcom.StateDB, fixDelegation *common.FixDelegation) error {
	//退回当前委托挪用的锁仓金额
	if err := rt.ReturnWrongLockFunds(delAddr, improperRestrictingAmount, state); nil != err {
		return err
	}
	if a.del.RestrictingPlanHes.Cmp(improperRestrictingAmount) >= 0 {
		//委托的犹豫期委托金额（来自锁仓的） > 挪用的锁仓金额，则从 委托的犹豫期委托金额（来自锁仓的）直接扣除挪用的锁仓金额即可。
		a.del.RestrictingPlanHes.Sub(a.del.RestrictingPlanHes, improperRestrictingAmount)
		if a.candidate.IsNotEmpty() {
			//修改委托节点上的总委托金额（犹豫期的）
			a.candidate.DelegateTotalHes.Sub(a.candidate.DelegateTotalHes, improperRestrictingAmount)
		}
		//stats
		fixDelegation.ImproperHesitatingRestrictingAmount = improperRestrictingAmount
		fixDelegation.ImproperValidRestrictingAmount = common.Big0
	} else {
		// 委托的犹豫期委托金额（来自锁仓的） <= 挪用的锁仓金额， 则扣除掉委托的犹豫期委托金额（来自锁仓的），并把余下的挪用金额，从委托的有效委托金额（来自锁仓的）里扣除；
		hes := new(big.Int).Set(a.del.RestrictingPlanHes)
		a.del.RestrictingPlanHes = new(big.Int)
		a.del.RestrictingPlan = new(big.Int).Sub(a.del.RestrictingPlan, new(big.Int).Sub(improperRestrictingAmount, hes))
		if a.candidate.IsNotEmpty() {
			// 修改委托节点上的总委托金额（犹豫期的，以及有效的）
			a.candidate.DelegateTotalHes.Sub(a.candidate.DelegateTotalHes, hes)
			a.candidate.DelegateTotal.Sub(a.candidate.DelegateTotal, new(big.Int).Sub(improperRestrictingAmount, hes))
		}

		//stats
		fixDelegation.ImproperHesitatingRestrictingAmount = hes
		fixDelegation.ImproperValidRestrictingAmount = new(big.Int).Sub(improperRestrictingAmount, hes)
	}
	if a.candidate.IsNotEmpty() {
		if a.candidate.Shares.Cmp(improperRestrictingAmount) >= 0 {
			//修改委托节点上的shares（用来算权重的）
			a.candidate.SubShares(improperRestrictingAmount)
		}
	}
	return nil
}

type issue1625AccountDelInfos []*issue1625AccountDelInfo

func (d issue1625AccountDelInfos) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d issue1625AccountDelInfos) Len() int      { return len(d) }
func (d issue1625AccountDelInfos) Less(i, j int) bool {
	//排序顺序
	//1.委托的节点已经完全退出,并且委托的时间靠后
	//2.委托的节点处于解质押状态,并且委托的时间靠后
	//3.根据委托节点的分红比例从小到大排序，如果委托比例相同，根据节点id从小到大排序
	if d[i].candidate.IsEmpty() {
		if d[j].candidate.IsEmpty() {
			return d.LessDelByEpoch(i, j)
		} else {
			return true
		}
	} else {
		if d[j].candidate.IsEmpty() {
			return false
		} else {
			if d[i].candidate.IsWithdrew() {
				if d[j].candidate.IsWithdrew() {
					return d.LessDelByEpoch(i, j)
				} else {
					return true
				}
			} else {
				if d[j].candidate.IsWithdrew() {
					return false
				} else {
					return d.LessDelByRewardPer(i, j)
				}
			}
		}
	}

}

func (d issue1625AccountDelInfos) LessDelByEpoch(i, j int) bool {
	if d[i].del.DelegateEpoch > d[j].del.DelegateEpoch {
		return true
	} else {
		return false
	}
}

func (d issue1625AccountDelInfos) LessDelByRewardPer(i, j int) bool {
	if d[i].candidate.RewardPer < d[j].candidate.RewardPer {
		return true
	} else if d[i].candidate.RewardPer == d[j].candidate.RewardPer {
		if bytes.Compare(d[i].candidate.NodeId.Bytes(), d[j].candidate.NodeId.Bytes()) < 0 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
