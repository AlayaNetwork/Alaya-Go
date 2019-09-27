package vm

import (
	"encoding/json"
	"math/big"
	"strings"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

const (
	CreateRestrictingPlanEvent = "4000"
)

type RestrictingContract struct {
	Plugin   *plugin.RestrictingPlugin
	Contract *Contract
	Evm      *EVM
}

func (rc *RestrictingContract) RequiredGas(input []byte) uint64 {
	return params.RestrictingPlanGas
}

func (rc *RestrictingContract) Run(input []byte) ([]byte, error) {
	return exec_platon_contract(input, rc.FnSigns())
}

func (rc *RestrictingContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		// Set
		4000: rc.createRestrictingPlan,

		// Get
		4100: rc.getRestrictingInfo,
		4101: rc.getRestrictingBalance,
	}
}

func (rc *RestrictingContract) CheckGasPrice(gasPrice *big.Int, fcode uint16) error {
	return nil
}

// createRestrictingPlan is a PlatON precompiled contract function, used for create a restricting plan
func (rc *RestrictingContract) createRestrictingPlan(account common.Address, plans []restricting.RestrictingPlan) ([]byte, error) {

	//sender := rc.Contract.Caller()
	from := rc.Contract.CallerAddress
	txHash := rc.Evm.StateDB.TxHash()
	blockNum := rc.Evm.BlockNumber
	blockHash := rc.Evm.BlockHash
	state := rc.Evm.StateDB

	log.Info("Call createRestrictingPlan of RestrictingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNum.Uint64(), "blockHash", blockHash.Hex(), "from", from.Hex())

	if !rc.Contract.UseGas(params.CreateRestrictingPlanGas) {
		return nil, ErrOutOfGas
	}
	if !rc.Contract.UseGas(params.ReleasePlanGas * uint64(len(plans))) {
		return nil, ErrOutOfGas
	}
	if txHash == common.ZeroHash {
		log.Warn("Call createRestrictingPlan current txHash is empty!!")
		return nil, nil
	}

	err := rc.Plugin.AddRestrictingRecord(from, account, plans, state)
	switch err.(type) {
	case nil:
		event := xcom.OkResultByte
		rc.goodLog(state, blockNum.Uint64(), txHash.Hex(), CreateRestrictingPlanEvent, string(event), "createRestrictingPlan")
		return event, nil
	case *common.BizError:
		event := xcom.NewFailResult(err)
		rc.badLog(state, blockNum.Uint64(), txHash.Hex(), CreateRestrictingPlanEvent, string(event), "createRestrictingPlan")
		return event, nil
	default:
		log.Error("AddRestrictingRecord failed to createRestrictingPlan", "txHash", txHash.Hex(), "blockNumber", blockNum.Uint64(), "error", err)
		return nil, err
	}
}

// createRestrictingPlan is a PlatON precompiled contract function, used for getting restricting info.
// first output param is a slice of byte of restricting info;
// the secend output param is the result what plugin executed GetRestrictingInfo returns.
func (rc *RestrictingContract) getRestrictingInfo(account common.Address) ([]byte, error) {
	txHash := rc.Evm.StateDB.TxHash()
	currNumber := rc.Evm.BlockNumber
	state := rc.Evm.StateDB

	log.Info("Call getRestrictingInfo of RestrictingContract", "txHash", txHash.Hex(), "blockNumber", currNumber.Uint64())

	result, err := rc.Plugin.GetRestrictingInfo(account, state)
	//var res xcom.Result
	if err != nil {
		//res.Status = false
		//res.Data = ""
		//res.ErrMsg = "get restricting info:" + err.Error()
		return xcom.NewFailResult(err), nil
	} else {
		//res.Status = true
		//res.Data = string(result)
		//res.ErrMsg = "ok"
		return xcom.NewSuccessResult(string(result)), nil
	}
	//return json.Marshal(res)
}

func (rc *RestrictingContract) getRestrictingBalance(accounts string) ([]byte, error) {

	accountList := strings.Split(accounts, ";")
	if(len(accountList) == 0){
		log.Error("getRestrictingBalance accountList empty","accountList:",len(accountList))
		return nil, nil
	}

	txHash := rc.Evm.StateDB.TxHash()
	currNumber := rc.Evm.BlockNumber
	state := rc.Evm.StateDB

	log.Info("Call getRestrictingBalance of RestrictingContract", "txHash", txHash.Hex(), "blockNumber", currNumber.Uint64())

	rs := make([]restricting.BalanceResult, len(accountList))
	for i, account := range accountList {
		address := common.HexToAddress(account)
		result, err := rc.Plugin.GetRestrictingBalance(address, state)
		if err != nil {
			rb := restricting.BalanceResult{
				Account : address,
			}
			rs[i] = rb
			log.Error("getRestrictingBalance err","account:",account,";err",err)
		} else {
			rs[i] = result
		}
	}
	resByte, _ := json.Marshal(rs)
	return xcom.NewSuccessResult(string(resByte)), nil
}

func (rc *RestrictingContract) goodLog(state xcom.StateDB, blockNumber uint64, txHash, eventType, eventData, callFn string) {
	xcom.AddLog(state, blockNumber, vm.RestrictingContractAddr, eventType, eventData)
	log.Info("Successed to "+callFn, "txHash", txHash, "blockNumber", blockNumber, "json: ", eventData)
}

func (rc *RestrictingContract) badLog(state xcom.StateDB, blockNumber uint64, txHash, eventType, eventData, callFn string) {
	xcom.AddLog(state, blockNumber, vm.RestrictingContractAddr, eventType, eventData)
	log.Debug("Failed to "+callFn, "txHash", txHash, "blockNumber", blockNumber, "json: ", eventData)
}
