package vm_test

import (
	"bytes"
	"encoding/json"
	_ "fmt"
	"github.com/PlatONnetwork/PlatON-Go/eth"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"math/big"
	"strconv"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

// Custom func
func create_staking(blockNumber *big.Int, blockHash common.Hash, state *state.StateDB, index int, t *testing.T) *vm.StakingContract {

	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:      newEvm(blockNumber, blockHash, state),
	}

	state.Prepare(txHashArr[index], blockHash, index+1)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	typ, _ := rlp.EncodeToBytes(uint16(0))
	benefitAddress, _ := rlp.EncodeToBytes(addrArr[index])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", China")
	website, _ := rlp.EncodeToBytes("https://www." + nodeNameArr[index] + ".network")
	details, _ := rlp.EncodeToBytes(nodeNameArr[index] + " super node")
	StakeThreshold, _ := new(big.Int).SetString(balanceStr[index], 10) // equal or more than "1000000000000000000000000"
	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, programVersion)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	} else {
		t.Log(string(res))
	}

	return contract
}

func create_delegate(contract *vm.StakingContract, index int, t *testing.T) {
	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1004))
	typ, _ := rlp.EncodeToBytes(uint16(0))
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	StakeThreshold, _ := new(big.Int).SetString(balanceStr[index+6], 10)
	amount, _ := rlp.EncodeToBytes(StakeThreshold)

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, nodeId)
	params = append(params, amount)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Error("delegate encode rlp data fail", err)
		return
	} else {
		t.Log("delegate data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	} else {
		t.Log(string(res))
	}
}

func getCandidate(contract *vm.StakingContract, index int, t *testing.T) {
	params := make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1105))
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])

	params = append(params, fnType)
	params = append(params, nodeId)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("getCandidate encode rlp data fail: %v", err)
		return
	} else {
		t.Log("getCandidate data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error("getCandidate err", err)
	} else {

		var r xcom.Result
		err = json.Unmarshal(res, &r)
		if nil != err {
			t.Error("Failed to parse result", err)
		}

		if r.Status {
			t.Log("the Candidate info:", r.Data)
		} else {
			t.Log("getCandidate failed", r.ErrMsg)
		}
	}
}

func TestRLP_encode(t *testing.T) {

	var params [][]byte
	params = make([][]byte, 0)

	fnType, err := rlp.EncodeToBytes(uint16(1100))
	if nil != err {
		t.Error("fnType err", err)
		return
	} else {
		var num uint16
		rlp.DecodeBytes(fnType, &num)
		t.Log("num is ", num)
	}
	params = append(params, fnType)

	buf := new(bytes.Buffer)
	err = rlp.Encode(buf, params)
	if err != nil {
		t.Log(err)
		t.Errorf("rlp stakingContract encode rlp data fail")
	} else {
		t.Log("rlp stakingContract data rlp: ", hexutil.Encode(buf.Bytes()))
	}
}

/**
Standard test cases
*/

func TestStakingContract_createStaking(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
	}
	create_staking(blockNumber, blockHash, state, 1, t)
}

func TestStakingContract_editCandidate(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	index := 1

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 error: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
		return
	}

	contract2 := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

	// edit
	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1001))

	benefitAddress, _ := rlp.EncodeToBytes(addrArr[0])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	externalId, _ := rlp.EncodeToBytes("I am Xu !?")
	nodeName, _ := rlp.EncodeToBytes("Xu, China")
	website, _ := rlp.EncodeToBytes("https://www.Xu.net")
	details, _ := rlp.EncodeToBytes("Xu super node")

	params = append(params, fnType)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("edit candidate encode rlp data fail: %v", err)
		return
	} else {
		t.Log("edit candidate data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())
	if nil != err {
		t.Error("Failed to Call editCandidate, err:", err)
		return
	} else {
		t.Log(string(res))
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 error: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

}

func TestStakingContract_increaseStaking(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	index := 1

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 error: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
		return
	}

	contract2 := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

	// increase

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1002))
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	typ, _ := rlp.EncodeToBytes(uint16(0))
	StakeThreshold, _ := new(big.Int).SetString(balanceStr[index-1], 10) // equal or more than "1000000000000000000000000"
	amount, _ := rlp.EncodeToBytes(StakeThreshold)

	params = append(params, fnType)
	params = append(params, nodeId)
	params = append(params, typ)
	params = append(params, amount)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Error("increaseStaking encode rlp data fail", err)
		return
	} else {
		t.Log("increaseStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())
	if nil != err {
		t.Error("Failed to Call increaseStaking,err:", err)
		return
	} else {
		t.Log(string(res))
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 error: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

}

func TestStakingContract_withdrewCandidate(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	index := 1

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
		return
	}

	contract2 := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

	// withdrewStaking

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1003))
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])

	params = append(params, fnType)
	params = append(params, nodeId)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Error("withdrewStaking encode rlp data fail", err)
		return
	} else {
		t.Log("withdrewStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())
	if nil != err {
		t.Error("Failed to Call withdrewStaking, err:", err)
		return
	} else {
		t.Log(string(res))
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 err: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

}

func TestStakingContract_delegate(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	index := 1

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
		return
	}

	contract2 := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

	// delegate
	create_delegate(contract2, index, t)

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 err: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

}

func TestStakingContract_withdrewDelegate(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	index := 1

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	// delegate
	create_delegate(contract1, index, t)

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
		return
	}

	contract2 := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

	// withdrewDelegate
	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1005))
	stakingBlockNum, _ := rlp.EncodeToBytes(blockNumber.Uint64())
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	StakeThreshold, _ := new(big.Int).SetString("4600000", 10)
	amount, _ := rlp.EncodeToBytes(StakeThreshold)

	params = append(params, fnType)
	params = append(params, stakingBlockNum)
	params = append(params, nodeId)
	params = append(params, amount)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Error("delegate encode rlp data fail", err)
		return
	} else {
		t.Log("delegate data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())
	if nil != err {
		t.Error("Failed to call delegate, err:", err)
		return
	} else {
		t.Log(string(res))
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 err: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)
}

func TestStakingContract_getVerifierList(t *testing.T) {

	state, genesis, _ := newChainState()
	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:      newEvm(blockNumber, blockHash, state),
	}
	//state.Prepare(txHashArr[idx], blockHash, idx)
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Errorf("newBlock failed, blockNumber1: %d, err:%v", blockNumber, err)
		return
	}

	params := make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1100))

	params = append(params, fnType)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("getVerifierList encode rlp data fail:%v", err)
		return
	} else {
		t.Log("getVerifierList data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error("Failed to call getVerifierList, err", err)
		return
	} else {

		var r xcom.Result
		err = json.Unmarshal(res, &r)
		if nil != err {
			t.Error("Failed tp parse result", err)
			return
		}

		if r.Status {
			t.Log("the VerifierList info:", r.Data)
		} else {
			t.Error("getVerifierList failed", r.ErrMsg)
		}
	}

}

func TestStakingContract_getHistoryVerifierList(t *testing.T) {

	state, genesis, _ := newChainState()
	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:      newEvm(blockNumber, blockHash, state),
	}
	//state.Prepare(txHashArr[idx], blockHash, idx)
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Errorf("newBlock failed, blockNumber1: %d, err:%v", blockNumber, err)
		return
	}

	// Request th opening/creation of an ephemeral database and ensure it's not persisted
	ctx := node.NewServiceContext(&node.Config{DataDir: ""}, nil, new(event.TypeMux), nil)
	config := &eth.Config{
	}
	hDB, _ := eth.CreateDB(ctx, config, "historydata")
	plugin.STAKING_DB = &plugin.StakingDB{
		HistoryDB:  hDB,
	}
	validatorQueue := make(staking.ValidatorQueue, len(nodeIdArr))
	for index, value := range nodeIdArr {
		validatorQueue[index] = &staking.Validator{
			NodeId:value,
		}
	}
	current := staking.Validator_array{Start:uint64(0), End:uint64(50), Arr:validatorQueue}
	data, _ := rlp.EncodeToBytes(current);
	numStr := strconv.FormatUint(xutil.ConsensusSize() - xcom.ElectionDistance(), 10)
	hDB.Put([]byte(plugin.ValidatorName+numStr), data)

	params := make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1106))
	param0, _ := rlp.EncodeToBytes(big.NewInt(int64(50)))

	params = append(params, fnType)
	params = append(params, param0)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("getHistoryVerifierList encode rlp data fail:%v", err)
		return
	} else {
		t.Log("getHistoryVerifierList data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error("Failed to call getHistoryVerifierList, err", err)
		return
	} else {

		var r xcom.Result
		err = json.Unmarshal(res, &r)
		if nil != err {
			t.Error("Failed tp parse result", err)
			return
		}

		if r.Status {
			t.Log("the VerifierList info:", r.Data)
		} else {
			t.Error("getHistoryVerifierList failed", r.ErrMsg)
		}
	}

}

func TestStakingContract_getHistoryValidatorList(t *testing.T) {

	state, genesis, _ := newChainState()
	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:      newEvm(blockNumber, blockHash, state),
	}
	//state.Prepare(txHashArr[idx], blockHash, idx)
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Errorf("newBlock failed, blockNumber1: %d, err:%v", blockNumber, err)
		return
	}

	// Request th opening/creation of an ephemeral database and ensure it's not persisted
	ctx := node.NewServiceContext(&node.Config{DataDir: ""}, nil, new(event.TypeMux), nil)
	config := &eth.Config{
	}
	hDB, _ := eth.CreateDB(ctx, config, "historydata")
	plugin.STAKING_DB = &plugin.StakingDB{
		HistoryDB:  hDB,
	}
	validatorQueue := make(staking.ValidatorQueue, len(nodeIdArr))
	for index, value := range nodeIdArr {
		validatorQueue[index] = &staking.Validator{
			NodeId:value,
		}
	}
	current := staking.Validator_array{Start:uint64(0), End:uint64(50), Arr:validatorQueue}
	data, _ := rlp.EncodeToBytes(current);
	numStr := strconv.FormatUint(xutil.CalcBlocksEachEpoch(), 10)
	hDB.Put([]byte(plugin.ValidatorName+numStr), data)

	params := make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1107))
	param0, _ := rlp.EncodeToBytes(big.NewInt(int64(50)))

	params = append(params, fnType)
	params = append(params, param0)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("getHistoryValidatorList encode rlp data fail:%v", err)
		return
	} else {
		t.Log("getHistoryValidatorList data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error("Failed to call getHistoryValidatorList, err", err)
		return
	} else {

		var r xcom.Result
		err = json.Unmarshal(res, &r)
		if nil != err {
			t.Error("Failed tp parse result", err)
			return
		}

		if r.Status {
			t.Log("the Validator info:", r.Data)
		} else {
			t.Error("getHistoryValidatorList failed", r.ErrMsg)
		}
	}

}

func TestStakingContract_getValidatorList(t *testing.T) {

	state, genesis, _ := newChainState()
	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:      newEvm(blockNumber, blockHash, state),
	}
	//state.Prepare(txHashArr[idx], blockHash, idx)
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Errorf("newBlock failed, blockNumber1: %d, err:%v", blockNumber, err)
		return
	}

	params := make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1101))

	params = append(params, fnType)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("getValidatorList encode rlp data fail:%v", err)
		return
	} else {
		t.Log("getValidatorList data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error("Failed to Call getValidatorList, err", err)
		return
	} else {

		var r xcom.Result
		err = json.Unmarshal(res, &r)
		if nil != err {
			t.Error("Failed to parse result", err)
			return
		}

		if r.Status {
			t.Log("the ValidatorList info:", r.Data)
		} else {
			t.Error("getValidatorList failed", r.ErrMsg)
		}
	}

}

func TestStakingContract_getCandidateList(t *testing.T) {

	state, genesis, _ := newChainState()

	//state.Prepare(txHashArr[idx], blockHash, idx)
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Errorf("newBlock failed, blockNumber1: %d, err:%v", blockNumber, err)
		return
	}

	for i := 0; i < 2; i++ {
		create_staking(blockNumber, blockHash, state, i, t)
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}
	//sndb.Compaction()

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
		return
	}

	for i := 2; i < 4; i++ {
		create_staking(blockNumber2, blockHash2, state, i, t)
	}

	// getCandidate List
	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}
	params := make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1102))

	params = append(params, fnType)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("getCandidateList encode rlp data fail:%v", err)
		return
	} else {
		t.Log("getCandidateList data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error("Failed to Call getCandidateList, err", err)
		return
	} else {

		var r xcom.Result
		err = json.Unmarshal(res, &r)
		if nil != err {
			t.Error("Failed to parse result", err)
			return
		}

		if r.Status {
			t.Log("the CandidateList info:", r.Data)
		} else {
			t.Error("CandidateList failed", r.ErrMsg)
		}
	}

}

func TestStakingContract_getRelatedListByDelAddr(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	for i := 0; i < 4; i++ {
		create_staking(blockNumber, blockHash, state, i, t)
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
		return
	}

	contract2 := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}

	// delegate
	for i := 0; i < 3; i++ {
		create_delegate(contract2, i, t)
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 err: %v", err)
		return
	}

	// get RelatedListByDelAddr
	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1103))
	delAddr, _ := rlp.EncodeToBytes(sender)

	params = append(params, fnType)
	params = append(params, delAddr)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Error("getRelatedListByDelAddr encode rlp data fail", err)
		return
	} else {
		t.Log("getRelatedListByDelAddr data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())
	if nil != err {
		t.Error("Failed to call getRelatedListByDelAddr, err:", err)
		return
	} else {

		var r *xcom.Result
		if err := json.Unmarshal(res, &r); nil != err {
			t.Error("Failed to parse json", err)
		} else {
			t.Log("the Related list is:", r.Data)
		}
	}
}

func TestStakingContract_getDelegateInfo(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	index := 1

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	// delegate
	create_delegate(contract1, index, t)

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Failed to commit block 1, err: %v", err)
		return
	}
	//sndb.Compaction()

	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
		return
	}

	contract2 := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

	// get DelegateInfo
	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1104))
	stakingBlockNum, _ := rlp.EncodeToBytes(blockNumber.Uint64())
	delAddr, _ := rlp.EncodeToBytes(sender)
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])

	params = append(params, fnType)
	params = append(params, stakingBlockNum)
	params = append(params, delAddr)
	params = append(params, nodeId)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Error("getDelegateInfo encode rlp data fail", err)
		return
	} else {
		t.Log("getDelegateInfo data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())
	if nil != err {
		t.Error("Failed to call getDelegateInfo, err:", err)
	} else {
		var r xcom.Result

		err = json.Unmarshal(res, &r)
		if nil != err {
			t.Errorf("parse json failed, err: %v", err)
		} else {
			t.Log(r.Data)
		}
	}
}

func TestStakingContract_getCandidateInfo(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("Failed to newBlock", err)
		return
	}
	contract := create_staking(blockNumber, blockHash, state, 1, t)
	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}
	//sndb.Compaction()

	// get candidate Info
	getCandidate(contract, 1, t)
}

/**
Expand test cases
*/

func TestStakingContract_batchCreateStaking(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("Failed to newBlock", err)
		return
	}

	for i := 0; i < 4; i++ {
		create_staking(blockNumber, blockHash, state, i, t)
	}

}

func TestStakingContract_cleanSnapshotDB(t *testing.T) {
	sndb := snapshotdb.Instance()
	sndb.Clear()
}
