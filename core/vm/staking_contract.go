package vm

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/x/handler"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

const (
	AmountIllegalErrStr      = "This amount is too low"
	CanAlreadyExistsErrStr   = "This candidate is already exists"
	CanNotExistErrStr        = "This candidate is not exist"
	CreateCanErrStr          = "Create candidate failed"
	CanStatusInvalidErrStr   = "This candidate status was invalided"
	CanNoAllowDelegateErrStr = "This candidate is not allow to delegate"
	DelegateNotExistErrStr   = "This is delegate is not exist"
	DelegateErrStr           = "Delegate failed"
	DelegateVonTooLowStr     = "Delegate deposit too low"
	EditCanErrStr            = "Edit candidate failed"
	GetVerifierListErrStr    = "Getting verifierList is failed"
	GetValidatorListErrStr   = "Getting validatorList is failed"
	GetHistoryVerifierListErrStr = "Getting verifierList history is failed"
	GetHistoryValidatorListErrStr = "Getting validatorList history is failed"
	GetCandidateListErrStr   = "Getting candidateList is failed"
	GetDelegateRelatedErrStr = "Getting related of delegate is failed"
	IncreaseStakingErrStr    = "IncreaseStaking failed"
	ProgramVersionErrStr     = "The program version of the relates node's is too low"
	ProgramVersionSignErrStr = "The program version sign is wrong"
	QueryCanErrStr           = "Query candidate info failed"
	QueryDelErrSTr           = "Query delegate info failed"
	StakeVonTooLowStr        = "Staking deposit too low"
	StakingAddrNoSomeErrStr  = "Address must be the same as initiated staking"
	DescriptionLenErrStr     = "The Description length is wrong"
	WithdrewCanErrStr        = "Withdrew candidate failed"
	WithdrewDelegateErrStr   = "Withdrew delegate failed"
	WrongBlsPubKeyStr        = "The bls public key is wrong"
)

const (
	CreateStakingEvent     = "1000"
	EditorCandidateEvent   = "1001"
	IncreaseStakingEvent   = "1002"
	WithdrewCandidateEvent = "1003"
	DelegateEvent          = "1004"
	WithdrewDelegateEvent  = "1005"
	BLSPUBKEYLEN           = 192 //  the bls public key length must be 192 character
)

type StakingContract struct {
	Plugin   *plugin.StakingPlugin
	Contract *Contract
	Evm      *EVM
}

func (stkc *StakingContract) RequiredGas(input []byte) uint64 {
	return params.StakingGas
}

func (stkc *StakingContract) Run(input []byte) ([]byte, error) {
	return exec_platon_contract(input, stkc.FnSigns())
}

func (stkc *StakingContract) CheckGasPrice(gasPrice *big.Int, fcode uint16) error {
	return nil
}

func (stkc *StakingContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		// Set
		1000: stkc.createStaking,
		1001: stkc.editCandidate,
		1002: stkc.increaseStaking,
		1003: stkc.withdrewStaking,
		1004: stkc.delegate,
		1005: stkc.withdrewDelegate,

		// Get
		1100: stkc.getVerifierList,
		1101: stkc.getValidatorList,
		1102: stkc.getCandidateList,
		1103: stkc.getRelatedListByDelAddr,
		1104: stkc.getDelegateInfo,
		1105: stkc.getCandidateInfo,
		1106: stkc.getHistoryVerifierList,
		1107: stkc.getHistoryValidatorList,
	}
}

func (stkc *StakingContract) createStaking(typ uint16, benefitAddress common.Address, nodeId discover.NodeID,
	externalId, nodeName, website, details string, amount *big.Int, programVersion uint32,
	programVersionSign common.VersionSign, blsPubKey string) ([]byte, error) {
	txHash := stkc.Evm.StateDB.TxHash()
	txIndex := stkc.Evm.StateDB.TxIdx()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	from := stkc.Contract.CallerAddress

	state := stkc.Evm.StateDB

	log.Info("Call createStaking of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "typ", typ,
		"benefitAddress", benefitAddress.String(), "nodeId", nodeId.String(), "externalId", externalId,
		"nodeName", nodeName, "website", website, "details", details, "amount", amount,
		"programVersion", programVersion, "programVersionSign", programVersionSign.Hex(),
		"from", from.Hex(), "blsPubKey", blsPubKey)

	if !stkc.Contract.UseGas(params.CreateStakeGas) {
		return nil, ErrOutOfGas
	}

	if txHash == common.ZeroHash {
		log.Warn("Call createStaking current txHash is empty!!")
		return nil, nil
	}

	if len(blsPubKey) != BLSPUBKEYLEN {
		res := xcom.Result{false, "", WrongBlsPubKeyStr + ": " + fmt.Sprintf("got length: %d, must be: %d", len(blsPubKey), BLSPUBKEYLEN)}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, CreateStakingEvent, string(event), "createStaking")
		return event, nil
	}

	// validate programVersion sign
	if !handler.GetCryptoHandler().IsSignedByNodeID(programVersion, programVersionSign.Bytes(), nodeId) {
		res := xcom.Result{false, "", ProgramVersionSignErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, CreateStakingEvent, string(event), "createStaking")
		return event, nil
	}

	if !xutil.CheckStakeThreshold(amount) {
		res := xcom.Result{false, "", StakeVonTooLowStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, CreateStakingEvent, string(event), "createStaking")
		return event, nil
	}

	// check Description length
	desc := &staking.Description{
		NodeName:   nodeName,
		ExternalId: externalId,
		Website:    website,
		Details:    details,
	}
	if err := desc.CheckLength(); nil != err {
		res := xcom.Result{false, "", DescriptionLenErrStr + ": " + err.Error()}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, CreateStakingEvent, string(event), "createStaking")
		return event, nil
	}

	// Query current active version
	originVersion := gov.GetVersionForStaking(state)
	currVersion := xutil.CalcVersion(originVersion)
	inputVersion := xutil.CalcVersion(programVersion)

	var isDeclareVersion bool

	// Compare version
	// Just like that:
	// eg: 2.1.x == 2.1.x; 2.1.x > 2.0.x
	if inputVersion < currVersion {
		err := fmt.Errorf("input Version: %s, current valid Version: %s", xutil.ProgramVersion2Str(programVersion), xutil.ProgramVersion2Str(originVersion))
		res := xcom.Result{false, "", ProgramVersionErrStr + ": " + err.Error()}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, CreateStakingEvent, string(event), "createStaking")
		return event, nil

	} else if inputVersion > currVersion {
		isDeclareVersion = true
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to createStaking by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to createStaking by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil != canOld {
		res := xcom.Result{false, "", CanAlreadyExistsErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, CreateStakingEvent, string(event), "createStaking")
		return event, nil
	}

	// parse bls publickey
	var blsPk bls.PublicKey
	if err := blsPk.UnmarshalText([]byte(blsPubKey)); nil != err {
		res := xcom.Result{false, "", WrongBlsPubKeyStr + ": " + err.Error()}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, CreateStakingEvent, string(event), "createStaking")
		return event, nil
	}

	/**
	init candidate info
	*/
	canNew := &staking.Candidate{
		NodeId:          nodeId,
		BlsPubKey:       blsPk,
		StakingAddress:  from,
		BenefitAddress:  benefitAddress,
		StakingBlockNum: blockNumber.Uint64(),
		StakingTxIndex:  txIndex,
		Shares:          amount,

		// Prevent null pointer initialization
		Released:           common.Big0,
		ReleasedHes:        common.Big0,
		RestrictingPlan:    common.Big0,
		RestrictingPlanHes: common.Big0,

		Description: *desc,
	}

	canNew.ProgramVersion = currVersion

	err = stkc.Plugin.CreateCandidate(state, blockHash, blockNumber, amount, typ, canAddr, canNew)

	if nil != err {
		if _, ok := err.(*common.BizError); ok {
			res := xcom.Result{false, "", CreateCanErrStr + ": " + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash, CreateStakingEvent, string(event), "createStaking")
			return event, nil
		} else {
			log.Error("Failed to createStaking by CreateCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}
	}

	if isDeclareVersion {
		// Declare new Version
		err := gov.DeclareVersion(canNew.StakingAddress, canNew.NodeId,
			programVersion, programVersionSign, blockHash, blockNumber.Uint64(), stkc.Plugin, state)
		if nil != err {
			log.Error("Failed to CreateCandidate with govplugin DelareVersion failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)

			if er := stkc.Plugin.RollBackStaking(state, blockHash, blockNumber, canAddr, typ); nil != er {
				log.Error("Failed to createStaking by RollBackStaking", "txHash", txHash,
					"blockNumber", blockNumber, "err", er)
			}

			res := xcom.Result{false, "", CreateCanErrStr + ": Call DeclareVersion is failed, " + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash, CreateStakingEvent, string(event), "createStaking")
			return event, nil
		}
	}

	res := xcom.Result{true, "", "ok"}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash, CreateStakingEvent, string(event), "createStaking")
	return event, nil
}

func (stkc *StakingContract) editCandidate(benefitAddress common.Address, nodeId discover.NodeID,
	externalId, nodeName, website, details string) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	from := stkc.Contract.CallerAddress

	state := stkc.Evm.StateDB

	log.Info("Call editCandidate of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "benefitAddress", benefitAddress.String(),
		"nodeId", nodeId.String(), "externalId", externalId, "nodeName", nodeName,
		"website", website, "details", details, "from", from.Hex())

	if !stkc.Contract.UseGas(params.EditCandidatGas) {
		return nil, ErrOutOfGas
	}

	if txHash == common.ZeroHash {
		log.Warn("Call editCandidate current txHash is empty!!")
		return nil, nil
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to editCandidate by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to editCandidate by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil == canOld {
		res := xcom.Result{false, "", CanNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, EditorCandidateEvent, string(event), "editCandidate")
		return event, nil
	}

	if staking.Is_Invalid(canOld.Status) {
		res := xcom.Result{false, "", CanStatusInvalidErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, EditorCandidateEvent, string(event), "editCandidate")
		return event, nil
	}

	if from != canOld.StakingAddress {
		res := xcom.Result{false, "", StakingAddrNoSomeErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, EditorCandidateEvent, string(event), "editCandidate")
		return event, nil
	}

	canOld.BenefitAddress = benefitAddress

	// check Description length
	desc := &staking.Description{
		NodeName:   nodeName,
		ExternalId: externalId,
		Website:    website,
		Details:    details,
	}
	if err := desc.CheckLength(); nil != err {
		res := xcom.Result{false, "", DescriptionLenErrStr + ": " + err.Error()}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, EditorCandidateEvent, string(event), "editCandidate")
		return event, nil
	}

	canOld.Description = *desc

	err = stkc.Plugin.EditCandidate(blockHash, blockNumber, canOld)

	if nil != err {

		if _, ok := err.(*common.BizError); ok {
			res := xcom.Result{false, "", EditCanErrStr + ": " + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash, EditorCandidateEvent, string(event), "editCandidate")
			return event, nil
		} else {
			log.Error("Failed to editCandidate by EditCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}

	}
	res := xcom.Result{true, "", "ok"}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash, EditorCandidateEvent, string(event), "editCandidate")
	return event, nil
}

func (stkc *StakingContract) increaseStaking(nodeId discover.NodeID, typ uint16, amount *big.Int) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	from := stkc.Contract.CallerAddress

	state := stkc.Evm.StateDB

	log.Info("Call increaseStaking of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "nodeId", nodeId.String(), "typ", typ,
		"amount", amount, "from", from.Hex())

	if !stkc.Contract.UseGas(params.IncStakeGas) {
		return nil, ErrOutOfGas
	}

	if txHash == common.ZeroHash {
		log.Warn("Call increaseStaking current txHash is empty!!")
		return nil, nil
	}

	if !xutil.CheckMinimumThreshold(amount) {
		res := xcom.Result{false, "", AmountIllegalErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, IncreaseStakingEvent, string(event), "increaseStaking")
		return event, nil
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to increaseStaking by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to increaseStaking by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil == canOld {
		res := xcom.Result{false, "", CanNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, IncreaseStakingEvent, string(event), "increaseStaking")
		return event, nil
	}

	if staking.Is_Invalid(canOld.Status) {
		res := xcom.Result{false, "", CanStatusInvalidErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, IncreaseStakingEvent, string(event), "increaseStaking")
		return event, nil
	}

	if from != canOld.StakingAddress {
		res := xcom.Result{false, "", StakingAddrNoSomeErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, IncreaseStakingEvent, string(event), "increaseStaking")
		return event, nil
	}

	err = stkc.Plugin.IncreaseStaking(state, blockHash, blockNumber, amount, typ, canOld)

	if nil != err {

		if _, ok := err.(*common.BizError); ok {
			res := xcom.Result{false, "", IncreaseStakingErrStr + ": " + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash, IncreaseStakingEvent, string(event), "increaseStaking")
			return event, nil
		} else {
			log.Error("Failed to increaseStaking by EditCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}

	}
	res := xcom.Result{true, "", "ok"}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash, IncreaseStakingEvent, string(event), "increaseStaking")
	return event, nil
}

func (stkc *StakingContract) withdrewStaking(nodeId discover.NodeID) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	from := stkc.Contract.CallerAddress

	state := stkc.Evm.StateDB

	log.Info("Call withdrewStaking of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "nodeId", nodeId.String(), "from", from.Hex())

	if !stkc.Contract.UseGas(params.WithdrewStakeGas) {
		return nil, ErrOutOfGas
	}

	if txHash == common.ZeroHash {
		log.Warn("Call withdrewStaking current txHash is empty!!")
		return nil, nil
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to withdrewStaking by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to withdrewStaking by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	if nil == canOld {
		res := xcom.Result{false, "", CanNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, WithdrewCandidateEvent, string(event), "withdrewStaking")
		return event, nil
	}

	if staking.Is_Invalid(canOld.Status) {
		res := xcom.Result{false, "", CanStatusInvalidErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, WithdrewCandidateEvent, string(event), "withdrewStaking")
		return event, nil
	}

	if from != canOld.StakingAddress {
		res := xcom.Result{false, "", StakingAddrNoSomeErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, WithdrewCandidateEvent, string(event), "withdrewStaking")
		return event, nil
	}

	err = stkc.Plugin.WithdrewStaking(state, blockHash, blockNumber, canOld)
	if nil != err {

		if _, ok := err.(*common.BizError); ok {
			res := xcom.Result{false, "", WithdrewCanErrStr + ": " + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash, WithdrewCandidateEvent, string(event), "withdrewStaking")
			return event, nil
		} else {
			log.Error("Failed to withdrewStaking by WithdrewStaking", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}

	}

	res := xcom.Result{true, "", "ok"}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash, WithdrewCandidateEvent,
		string(event), "withdrewStaking")
	return event, nil
}

func (stkc *StakingContract) delegate(typ uint16, nodeId discover.NodeID, amount *big.Int) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	from := stkc.Contract.CallerAddress

	state := stkc.Evm.StateDB

	log.Info("Call delegate of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "delAddr", from.Hex(), "typ", typ,
		"nodeId", nodeId.String(), "amount", amount)

	if !stkc.Contract.UseGas(params.DelegateGas) {
		return nil, ErrOutOfGas
	}

	if txHash == common.ZeroHash {
		log.Warn("Call delegate current txHash is empty!!")
		return nil, nil
	}

	if !xutil.CheckMinimumThreshold(amount) {
		res := xcom.Result{false, "", DelegateVonTooLowStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, DelegateEvent, string(event), "delegate")
		return event, nil
	}

	// check account
	hasStake, err := stkc.Plugin.HasStake(blockHash, from)
	if nil != err {
		return nil, err
	}

	if hasStake {
		res := xcom.Result{false, "", DelegateErrStr + ": Account of Candidate(Validator)  is not allowed to be used for delegating"}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, DelegateEvent, string(event), "delegate")
		return event, nil
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to delegate by parse nodeId", "txHash", txHash, "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to delegate by GetCandidateInfo", "txHash", txHash, "blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil == canOld {
		res := xcom.Result{false, "", CanNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, DelegateEvent, string(event), "delegate")
		return event, nil
	}

	if staking.Is_Invalid(canOld.Status) {
		res := xcom.Result{false, "", CanStatusInvalidErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, DelegateEvent, string(event), "delegate")
		return event, nil
	}

	// If the candidate’s benefitaAddress is the RewardManagerPoolAddr, no delegation is allowed
	if canOld.BenefitAddress == vm.RewardManagerPoolAddr {
		res := xcom.Result{false, "", CanNoAllowDelegateErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, DelegateEvent, string(event), "delegate")
		return event, nil
	}

	// todo the delegate caller is candidate stake addr ?? How do that ?? Do not allow !!

	del, err := stkc.Plugin.GetDelegateInfo(blockHash, from, nodeId, canOld.StakingBlockNum)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to delegate by GetDelegateInfo", "txHash", txHash, "blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil == del {

		// build delegate
		del = new(staking.Delegation)

		// Prevent null pointer initialization
		del.Released = common.Big0
		del.RestrictingPlan = common.Big0
		del.ReleasedHes = common.Big0
		del.RestrictingPlanHes = common.Big0
		del.Reduction = common.Big0
	}

	err = stkc.Plugin.Delegate(state, blockHash, blockNumber, from, del, canOld, typ, amount)
	if nil != err {
		if _, ok := err.(*common.BizError); ok {
			res := xcom.Result{false, "", DelegateErrStr + ": " + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash, DelegateEvent, string(event), "delegate")
			return event, nil
		} else {
			log.Error("Failed to delegate by Delegate", "txHash", txHash, "blockNumber", blockNumber, "err", err)
			return nil, err
		}
	}

	res := xcom.Result{true, "", "ok"}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash, DelegateEvent, string(event), "delegate")
	return event, nil
}

func (stkc *StakingContract) withdrewDelegate(stakingBlockNum uint64, nodeId discover.NodeID, amount *big.Int) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	from := stkc.Contract.CallerAddress

	state := stkc.Evm.StateDB

	log.Info("Call withdrewDelegate of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "delAddr", from.Hex(), "nodeId", nodeId.String(),
		"stakingNum", stakingBlockNum, "amount", amount)

	if !stkc.Contract.UseGas(params.WithdrewDelegateGas) {
		return nil, ErrOutOfGas
	}

	if txHash == common.ZeroHash {
		log.Warn("Call withdrewDelegate current txHash is empty!!")
		return nil, nil
	}

	if !xutil.CheckMinimumThreshold(amount) {
		res := xcom.Result{false, "", AmountIllegalErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, WithdrewDelegateEvent, string(event), "withdrewDelegate")
		return event, nil
	}

	del, err := stkc.Plugin.GetDelegateInfo(blockHash, from, nodeId, stakingBlockNum)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to withdrewDelegate by GetDelegateInfo",
			"txHash", txHash.Hex(), "blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil == del {
		res := xcom.Result{false, "", DelegateNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, WithdrewDelegateEvent, string(event), "withdrewDelegate")
		return event, nil
	}

	err = stkc.Plugin.WithdrewDelegate(state, blockHash, blockNumber, amount, from, nodeId, stakingBlockNum, del)
	if nil != err {
		if _, ok := err.(*common.BizError); ok {
			res := xcom.Result{false, "", WithdrewDelegateErrStr + ": " + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash, WithdrewDelegateEvent, string(event), "withdrewDelegate")
			return event, nil
		} else {
			log.Error("Failed to withdrewDelegate by WithdrewDelegate", "txHash", txHash, "blockNumber", blockNumber, "err", err)
			return nil, err
		}
	}

	res := xcom.Result{true, "", "ok"}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash, WithdrewDelegateEvent, string(event), "withdrewDelegate")
	return event, nil
}

func (stkc *StakingContract) getVerifierList() ([]byte, error) {

	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	arr, err := stkc.Plugin.GetVerifierList(blockHash, blockNumber.Uint64(), plugin.QueryStartIrr)

	if nil != err && err != snapshotdb.ErrNotFound {
		res := xcom.Result{false, "", GetVerifierListErrStr + ": " + err.Error()}
		data, _ := json.Marshal(res)
		log.Error("Failed to getVerifierList: Query VerifierList is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return data, nil
	}

	if (nil != err && err == snapshotdb.ErrNotFound) || nil == arr {
		res := xcom.Result{false, "", "VerifierList info is not found"}
		data, _ := json.Marshal(res)
		log.Error("Failed to getVerifierList: VerifierList info is not found",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex())
		return data, nil
	}

	arrByte, err := json.Marshal(arr)
	if nil != err {
		res := xcom.Result{false, "", GetVerifierListErrStr + ": " + err.Error()}
		data, _ := json.Marshal(res)
		log.Error("Failed to getVerifierList: VerifierList Marshal json is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return data, nil
	}
	res := xcom.Result{true, string(arrByte), "ok"}
	data, _ := json.Marshal(res)

	log.Info("getVerifierList", "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "verArr", string(arrByte))
	return data, nil
}

func (stkc *StakingContract) getHistoryVerifierList(blockNumber *big.Int) ([]byte, error) {

	blockHash := stkc.Evm.BlockHash

	arr, err := stkc.Plugin.GetHistoryVerifierList(blockHash, blockNumber.Uint64(), plugin.QueryStartIrr)
	if nil != err {
		res := xcom.Result{false, "", GetHistoryVerifierListErrStr + ": " + err.Error()}
		data, _ := json.Marshal(res)
		return data, nil
	}
	arrByte, _ := json.Marshal(arr)
	res := xcom.Result{true, string(arrByte), "ok"}
	data, _ := json.Marshal(res)

	// todo test
	log.Debug("getHistoryVerifierList", "valArr", string(arrByte))
	return data, nil
}

func (stkc *StakingContract) getValidatorList() ([]byte, error) {

	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	arr, err := stkc.Plugin.GetValidatorList(blockHash, blockNumber.Uint64(), plugin.CurrentRound, plugin.QueryStartIrr)
	if nil != err && err != snapshotdb.ErrNotFound {
		res := xcom.Result{false, "", GetValidatorListErrStr + ": " + err.Error()}
		data, _ := json.Marshal(res)
		log.Error("Failed to getValidatorList: Query ValidatorList is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return data, nil
	}

	if (nil != err && err == snapshotdb.ErrNotFound) || nil == arr {
		res := xcom.Result{false, "", "ValidatorList info is not found"}
		data, _ := json.Marshal(res)
		log.Error("Failed to getValidatorList: ValidatorList info is not found",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex())
		return data, nil
	}

	arrByte, err := json.Marshal(arr)
	if nil != err {
		res := xcom.Result{false, "", GetValidatorListErrStr + ": " + err.Error()}
		data, _ := json.Marshal(res)
		log.Error("Failed to getValidatorList: ValidatorList Marshal json is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return data, nil
	}

	res := xcom.Result{true, string(arrByte), "ok"}
	data, _ := json.Marshal(res)

	log.Info("getValidatorList", "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "valArr", string(arrByte))
	return data, nil
}

func (stkc *StakingContract) getHistoryValidatorList(blockNumber *big.Int) ([]byte, error) {
	blockHash := stkc.Evm.BlockHash

	arr, err := stkc.Plugin.GetHistoryValidatorList(blockHash, blockNumber.Uint64(), plugin.CurrentRound, plugin.QueryStartIrr)
	if nil != err {
		res := xcom.Result{false, "", GetHistoryValidatorListErrStr + ": " + err.Error()}
		data, _ := json.Marshal(res)
		return data, nil
	}

	if nil == arr {
		res := xcom.Result{false, "", "ValidatorList info is not found"}
		data, _ := json.Marshal(res)
		return data, nil
	}

	arrByte, _ := json.Marshal(arr)
	res := xcom.Result{true, string(arrByte), "ok"}
	data, _ := json.Marshal(res)

	// todo test
	log.Debug("getHistoryValidatorList", "valArr", string(arrByte))
	return data, nil
}

func (stkc *StakingContract) getCandidateList() ([]byte, error) {

	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	arr, err := stkc.Plugin.GetCandidateList(blockHash, blockNumber.Uint64())
	if nil != err && err != snapshotdb.ErrNotFound {
		res := xcom.Result{false, "", GetCandidateListErrStr + ": " + err.Error()}
		data, _ := json.Marshal(res)
		log.Error("Failed to getCandidateList: Query CandidateList is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return data, nil
	}

	if (nil != err && err == snapshotdb.ErrNotFound) || nil == arr {
		res := xcom.Result{false, "", "CandidateList info is not found"}
		data, _ := json.Marshal(res)
		log.Error("Failed to getCandidateList: CandidateList info is not found",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex())
		return data, nil
	}

	arrByte, err := json.Marshal(arr)
	if nil != err {
		res := xcom.Result{false, "", GetCandidateListErrStr + ": " + err.Error()}
		data, _ := json.Marshal(res)
		log.Error("Failed to getCandidateList: CandidateList Marshal json is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return data, nil
	}
	res := xcom.Result{true, string(arrByte), "ok"}
	data, _ := json.Marshal(res)

	log.Info("getCandidateList", "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "canArr", string(arrByte))
	return data, nil
}

// todo Maybe will implement
func (stkc *StakingContract) getRelatedListByDelAddr(addr common.Address) ([]byte, error) {

	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	arr, err := stkc.Plugin.GetRelatedListByDelAddr(blockHash, addr)
	if nil != err && err != snapshotdb.ErrNotFound {
		res := xcom.Result{false, "", GetDelegateRelatedErrStr + ": " + err.Error()}
		data, _ := json.Marshal(res)
		log.Error("Failed to getRelatedListByDelAddr: Query RelatedList is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", addr.Hex(), "err", err)
		return data, nil
	}

	if (nil != err && err == snapshotdb.ErrNotFound) || nil == arr {
		res := xcom.Result{false, "", "RelatedList info is not found"}
		data, _ := json.Marshal(res)
		log.Error("Failed to getRelatedListByDelAddr: RelatedList info is not found",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", addr.Hex())
		return data, nil
	}

	jsonByte, err := json.Marshal(arr)
	if nil != err {
		res := xcom.Result{false, "", GetDelegateRelatedErrStr + ": " + err.Error()}
		data, _ := json.Marshal(res)
		log.Error("Failed to getRelatedListByDelAddr: RelatedList Marshal json is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", addr.Hex(), "err", err)
		return data, nil
	}
	res := xcom.Result{true, string(jsonByte), "ok"}
	data, _ := json.Marshal(res)

	log.Info("getRelatedListByDelAddr", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"delAddr", addr.Hex(), "relateArr", string(jsonByte))
	return data, nil
}

func (stkc *StakingContract) getDelegateInfo(stakingBlockNum uint64, delAddr common.Address,
	nodeId discover.NodeID) ([]byte, error) {

	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	del, err := stkc.Plugin.GetDelegateExCompactInfo(blockHash, blockNumber.Uint64(), delAddr, nodeId, stakingBlockNum)
	if nil != err && err != snapshotdb.ErrNotFound {
		res := xcom.Result{false, "", QueryDelErrSTr + ": " + err.Error()}
		data, _ := json.Marshal(res)
		log.Error("Failed to getDelegateInfo: Query Delegate info is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"delAddr", delAddr.Hex(), "nodeId", nodeId.String(), "stakingBlockNumber", stakingBlockNum, "err", err)
		return data, nil
	}

	if (nil != err && err == snapshotdb.ErrNotFound) || nil == del {
		res := xcom.Result{false, "", "Delegate info is not found"}
		data, _ := json.Marshal(res)
		log.Error("Failed to getDelegateInfo: Delegate info is not found",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"delAddr", delAddr.Hex(), "nodeId", nodeId.String(), "stakingBlockNumber", stakingBlockNum)
		return data, nil
	}

	jsonByte, err := json.Marshal(del)
	if nil != err {
		res := xcom.Result{false, "", QueryDelErrSTr + ": " + err.Error()}
		data, _ := json.Marshal(res)
		log.Error("Failed to getDelegateInfo: Delegate Marshal json is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"delAddr", delAddr.Hex(), "nodeId", nodeId.String(), "stakingBlockNumber", stakingBlockNum, "err", err)
		return data, nil
	}
	res := xcom.Result{true, string(jsonByte), "ok"}
	data, _ := json.Marshal(res)

	log.Info("getDelegateInfo", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"delAddr", delAddr.Hex(), "nodeId", nodeId.String(), "stakingBlockNumber", stakingBlockNum, "delinfo", string(jsonByte))
	return data, nil
}

func (stkc *StakingContract) getCandidateInfo(nodeId discover.NodeID) ([]byte, error) {

	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		res := xcom.Result{false, "", QueryCanErrStr + ": " + err.Error()}
		data, _ := json.Marshal(res)
		log.Error("Failed to getCandidateInfo: Parse NodeId to Address is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return data, nil
	}
	can, err := stkc.Plugin.GetCandidateCompactInfo(blockHash, blockNumber.Uint64(), canAddr)
	if nil != err && err != snapshotdb.ErrNotFound {
		res := xcom.Result{false, "", QueryCanErrStr + ": " + err.Error()}
		data, _ := json.Marshal(res)
		log.Error("Failed to getCandidateInfo: Query Candidate info is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return data, nil
	}

	if (nil != err && err == snapshotdb.ErrNotFound) || nil == can {
		res := xcom.Result{false, "", "Candidate info is not found"}
		data, _ := json.Marshal(res)
		log.Error("Failed to getCandidateInfo: Candidate info is not found",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String())
		return data, nil
	}

	jsonByte, err := json.Marshal(can)
	if nil != err {
		res := xcom.Result{false, "", QueryDelErrSTr + ": " + err.Error()}
		data, _ := json.Marshal(res)
		log.Error("Failed to getCandidateInfo: Candidate Marshal json is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return data, nil
	}
	res := xcom.Result{true, string(jsonByte), "ok"}
	data, _ := json.Marshal(res)

	log.Info("getCandidateInfo", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"nodeId", nodeId.String(), "caninfo", string(jsonByte))
	return data, nil
}

func (stkc *StakingContract) goodLog(state xcom.StateDB, blockNumber uint64, txHash common.Hash, eventType, eventData, callFn string) {
	xcom.AddLog(state, blockNumber, vm.StakingContractAddr, eventType, eventData)
	log.Info("Call "+callFn+" of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber, "json: ", eventData)
}

func (stkc *StakingContract) badLog(state xcom.StateDB, blockNumber uint64, txHash common.Hash, eventType, eventData, callFn string) {
	xcom.AddLog(state, blockNumber, vm.StakingContractAddr, eventType, eventData)
	log.Warn("Failed to "+callFn+" of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber, "json: ", eventData)
}
