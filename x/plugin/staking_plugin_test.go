package plugin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/eth"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"math/big"
	mrand "math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common/vm"

	"github.com/PlatONnetwork/PlatON-Go/x/handler"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

/**
test tool
*/
func Test_CleanSnapshotDB(t *testing.T) {
	sndb := snapshotdb.Instance()
	sndb.Clear()
}

func watching(eventMux *event.TypeMux, t *testing.T) {
	events := eventMux.Subscribe(cbfttypes.AddValidatorEvent{})
	defer events.Unsubscribe()

	for {
		select {
		case ev := <-events.Chan():
			if ev == nil {
				t.Error("ev is nil, may be Server closing")
				continue
			}

			switch ev.Data.(type) {
			case cbfttypes.AddValidatorEvent:
				addEv, ok := ev.Data.(cbfttypes.AddValidatorEvent)
				if !ok {
					t.Error("Received add validator event type error")
					continue
				}

				str, _ := json.Marshal(addEv)
				t.Log("P2P Received the add validator is:", string(str))
			default:
				t.Error("Received unexcepted event")
			}

		}
	}
}

func build_vrf_Nonce() ([]byte, [][]byte) {
	preNonces := make([][]byte, 0)
	curentNonce := crypto.Keccak256([]byte(string("nonce")))
	for i := 0; i < int(xcom.EpochValidatorNum()); i++ {
		preNonces = append(preNonces, crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(i))))[:])
		time.Sleep(time.Microsecond * 10)
	}
	return curentNonce, preNonces
}

func create_staking(state xcom.StateDB, blockNumber *big.Int, blockHash common.Hash, index int, typ uint16, t *testing.T) error {

	balance, _ := new(big.Int).SetString(balanceStr[index], 10)
	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()
	canTmp := &staking.Candidate{
		NodeId:          nodeIdArr[index],
		BlsPubKey:       *blsKey.GetPublicKey(),
		StakingAddress:  sender,
		BenefitAddress:  addrArr[index],
		StakingBlockNum: blockNumber.Uint64(),
		StakingTxIndex:  uint32(index),
		Shares:          balance,
		ProgramVersion:  xutil.CalcVersion(initProgramVersion),

		// Prevent null pointer initialization
		Released:           common.Big0,
		ReleasedHes:        common.Big0,
		RestrictingPlan:    common.Big0,
		RestrictingPlanHes: common.Big0,

		Description: staking.Description{
			NodeName:   nodeNameArr[index],
			ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+1)] + "balabalala" + chaList[index],
			Website:    "www." + nodeNameArr[index] + ".org",
			Details:    "This is " + nodeNameArr[index] + " Super Node",
		},
	}

	canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

	return StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, typ, canAddr, canTmp)
}

func getCandidate(blockHash common.Hash, index int) (*staking.Candidate, error) {
	addr, _ := xutil.NodeId2Addr(nodeIdArr[index])

	if can, err := StakingInstance().GetCandidateInfo(blockHash, addr); nil != err {
		return nil, err
	} else {

		return can, nil
	}
}

func delegate(state xcom.StateDB, blockHash common.Hash, blockNumber *big.Int,
	can *staking.Candidate, typ uint16, index int, t *testing.T) (*staking.Delegation, error) {

	delAddr := addrArr[index+1]

	// build delegate
	del := new(staking.Delegation)

	// Prevent null pointer initialization
	del.Released = common.Big0
	del.RestrictingPlan = common.Big0
	del.ReleasedHes = common.Big0
	del.RestrictingPlanHes = common.Big0
	del.Reduction = common.Big0

	//amount := common.Big257  // FAIL
	amount, _ := new(big.Int).SetString(balanceStr[index+1], 10) // PASS

	return del, StakingInstance().Delegate(state, blockHash, blockNumber, delAddr, del, can, 0, amount)
}

func getDelegate(blockHash common.Hash, stakingNum uint64, index int, t *testing.T) *staking.Delegation {

	del, err := StakingInstance().GetDelegateInfo(blockHash, addrArr[index+1], nodeIdArr[index], stakingNum)
	if nil != err {
		t.Log("Failed to GetDelegateInfo:", err)
	} else {
		delByte, _ := json.Marshal(del)
		t.Log("Get Delegate Info is:", string(delByte))
	}
	return del
}

/**
Standard test cases
*/

func TestStakingPlugin_BeginBlock(t *testing.T) {
	// nothings in that
}

func TestStakingPlugin_EndBlock(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// New VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	// build vrf proof
	// build ancestor nonces
	_, nonces := build_vrf_Nonce()
	enValue, err := rlp.EncodeToBytes(nonces)
	if nil != err {
		t.Error("Failed to rlp vrf nonces", "err", err)
		return
	}

	// new block
	privateKey, err := crypto.GenerateKey()
	if nil != err {
		t.Errorf("Failed to generate random Address private key: %v", err)
		return
	}
	nodeId := discover.PubkeyID(&privateKey.PublicKey)
	currentHash := crypto.Keccak256Hash([]byte(nodeId.String()))
	currentNumber := big.NewInt(1)

	// build genesis veriferList and validatorList
	validatorQueue := make(staking.ValidatorQueue, xcom.EpochValidatorNum())

	for j := 0; j < 1000; j++ {
		var index int = j % 25

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(j),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(j) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(j) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		// Store Candidate power
		powerKey := staking.TallyPowerKey(canTmp.Shares, canTmp.StakingBlockNum, canTmp.StakingTxIndex, canTmp.ProgramVersion)
		if err := sndb.PutBaseDB(powerKey, canAddr.Bytes()); nil != err {
			t.Errorf("Failed to Store Candidate Power: PutBaseDB failed. error:%s", err.Error())
			return
		}

		// Store Candidate info
		canKey := staking.CandidateKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canTmp); nil != err {
			t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canKey, val); nil != err {
				t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		if j < int(xcom.EpochValidatorNum()) {
			v := &staking.Validator{
				NodeAddress: canAddr,
				NodeId:      canTmp.NodeId,
				BlsPubKey:   canTmp.BlsPubKey,
				StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
					fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
				ValidatorTerm: 0,
			}
			validatorQueue[j] = v
		}

	}

	/**
	*******
	build genesis epoch validators
	*******
	*/
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}

	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
		return
	}

	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	/**
	*******
	build genesis curr round validators
	*******
	*/
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
		return
	}

	xcom.PrintObject("Test round", validatorQueue[:xcom.ConsValidatorNum()])
	roundArr, err := rlp.EncodeToBytes(validatorQueue[:xcom.ConsValidatorNum()])
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Current Round validator
	if err := sndb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), roundArr); nil != err {
		t.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	// Store vrf nonces
	if err := sndb.PutBaseDB(handler.NonceStorageKey, enValue); nil != err {
		t.Errorf("Failed to Store Current Vrf nonces : PutBaseDB failed. error:%s", err.Error())
		return
	}

	// SetCurrent to snapshotDB
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50
	preNum1 := new(big.Int).Sub(currentNumber, big.NewInt(1))
	if err := sndb.SetCurrent(currentHash, *preNum1, *preNum1); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error()))
	}

	/**
	EndBlock to Election()
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50

	nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
	header := &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	fmt.Println("currentHash Election", currentHash.Hex())

	err = StakingInstance().EndBlock(currentHash, header, state)
	if nil != err {
		t.Errorf("Failed to Election, blockNumber: %d, err: %v", currentNumber, err)
		return
	}

	if err := sndb.Commit(currentHash); nil != err {
		t.Errorf("Failed to Commit, blockNumber: %d, blockHHash: %s, err: %v", currentNumber, currentHash.Hex(), err)
		return
	}

	if err := sndb.Compaction(); nil != err {
		t.Errorf("Failed to Compaction, blockNumber: %d, blockHHash: %s, err: %v", currentNumber, currentHash.Hex(), err)
		return
	}

	// new block
	privateKey2, err := crypto.GenerateKey()
	if nil != err {
		t.Errorf("Failed to generate random Address private key: %v", err)
		return
	}
	nodeId2 := discover.PubkeyID(&privateKey2.PublicKey)
	currentHash = crypto.Keccak256Hash([]byte(nodeId2.String()))

	/**
	Elect Epoch validator list  == ElectionNextList()
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() * xutil.EpochSize())) // 600

	preNum := new(big.Int).Sub(currentNumber, big.NewInt(1)) // 599

	if err := sndb.SetCurrent(currentHash, *preNum, *preNum); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error()))
	}

	nonce = crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
	header = &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	fmt.Println("currentHash ElectionNextList", currentHash.Hex())
	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	err = StakingInstance().EndBlock(currentHash, header, state)
	if nil != err {
		t.Errorf("Failed to Election, blockNumber: %d, err: %v", currentNumber, err)
		return
	}
}

func TestStakingPlugin_Confirmed(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// New VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	// build vrf proof
	// build ancestor nonces
	_, nonces := build_vrf_Nonce()
	enValue, err := rlp.EncodeToBytes(nonces)
	if nil != err {
		t.Error("Failed to rlp vrf nonces", "err", err)
		return
	}

	// new block
	privateKey, err := crypto.GenerateKey()
	if nil != err {
		t.Errorf("Failed to generate random Address private key: %v", err)
		return
	}
	nodeId := discover.PubkeyID(&privateKey.PublicKey)
	currentHash := crypto.Keccak256Hash([]byte(nodeId.String()))
	currentNumber := big.NewInt(1)

	// build genesis veriferList and validatorList
	validatorQueue := make(staking.ValidatorQueue, xcom.EpochValidatorNum())

	for j := 0; j < 1000; j++ {
		var index int = j % 25

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(j),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(j) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(j) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		// Store Candidate power
		powerKey := staking.TallyPowerKey(canTmp.Shares, canTmp.StakingBlockNum, canTmp.StakingTxIndex, canTmp.ProgramVersion)
		if err := sndb.PutBaseDB(powerKey, canAddr.Bytes()); nil != err {
			t.Errorf("Failed to Store Candidate Power: PutBaseDB failed. error:%s", err.Error())
			return
		}

		// Store Candidate info
		canKey := staking.CandidateKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canTmp); nil != err {
			t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canKey, val); nil != err {
				t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		if j < int(xcom.EpochValidatorNum()) {
			v := &staking.Validator{
				NodeAddress: canAddr,
				NodeId:      canTmp.NodeId,
				BlsPubKey:   canTmp.BlsPubKey,
				StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
					fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
				ValidatorTerm: 0,
			}
			validatorQueue[j] = v
		}

	}
	/**
	*******
	build genesis epoch validators
	*******
	*/
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}

	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
		return
	}

	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	/**
	*******
	build genesis curr round validators
	*******
	*/
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
		return
	}

	xcom.PrintObject("Test round", validatorQueue[:xcom.ConsValidatorNum()])
	roundArr, err := rlp.EncodeToBytes(validatorQueue[:xcom.ConsValidatorNum()])
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Current Round validator
	if err := sndb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), roundArr); nil != err {
		t.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	// Store vrf nonces
	if err := sndb.PutBaseDB(handler.NonceStorageKey, enValue); nil != err {
		t.Errorf("Failed to Store Current Vrf nonces : PutBaseDB failed. error:%s", err.Error())
		return
	}

	// SetCurrent to snapshotDB
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50
	preNum1 := new(big.Int).Sub(currentNumber, big.NewInt(1))
	if err := sndb.SetCurrent(currentHash, *preNum1, *preNum1); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error()))
	}

	/**
	EndBlock to Election()
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50

	nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
	header := &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	fmt.Println("currentHash Election", currentHash.Hex())

	err = StakingInstance().EndBlock(currentHash, header, state)
	if nil != err {
		t.Errorf("Failed to Election, blockNumber: %d, err: %v", currentNumber, err)
		return
	}

	/**
	Start Confirmed
	*/

	eventMux := &event.TypeMux{}
	StakingInstance().SetEventMux(eventMux)
	go watching(eventMux, t)

	blockElection := types.NewBlock(header, nil, nil)

	err = StakingInstance().Confirmed(blockElection)
	if nil != err {
		t.Errorf("Failed to Confirmed, blockNumber: %d, err: %v", blockElection.Number().Uint64(), err)
		return
	}

}

func TestStakingPlugin_CreateCandidate(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	/**
	Start Create Staking
	*/
	if err := create_staking(state, blockNumber, blockHash, 1, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}

}

func TestStakingPlugin_GetCandidateInfo(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	/**
	Start Get Candidate Info
	*/
	can, err := getCandidate(blockHash, index)

	assert.True(t, nil == err)

	assert.True(t, nil != can)

	t.Log("Get Candidate Info is:", can)

	//canByte, _ := json.Marshal(can)
	//t.Log("Get Candidate Info is:", string(canByte))

}

func TestStakingPlugin_GetCandidateInfoByIrr(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	/**
	Start GetCandidateInfoByIrr

	Get Candidate Info
	*/
	addr, _ := xutil.NodeId2Addr(nodeIdArr[index])

	can, err := StakingInstance().GetCandidateInfoByIrr(addr)
	assert.True(t, nil == err)

	assert.True(t, nil != can)
	canByte, _ := json.Marshal(can)
	t.Log("Get Candidate Info is:", string(canByte))

}

func TestStakingPlugin_GetCandidateList(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	count := 0
	for i := 0; i < 4; i++ {
		if err := create_staking(state, blockNumber, blockHash, i, 0, t); nil != err {
			t.Error("Failed to Create num: "+fmt.Sprint(i)+" Staking", err)
			return
		}
		count++
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	/**
	Start GetCandidateList
	*/

	queue, err := StakingInstance().GetCandidateList(blockHash, blockNumber.Uint64())
	assert.True(t, nil == err)

	assert.Equal(t, count, len(queue))
	queueByte, _ := json.Marshal(queue)
	t.Log("Get CandidateList Info is:", string(queueByte))
}

func TestStakingPlugin_EditorCandidate(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}

	var c *staking.Candidate
	// Get Candidate Info
	if can, err := getCandidate(blockHash, index); nil != err {
		t.Errorf("Failed to Get candidate info, err: %v", err)
		return
	} else {
		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		c = can
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock2 err", err)
		return
	}

	/**
	Start Edit Candidate
	*/
	c.NodeName = nodeNameArr[index+1]
	c.ExternalId = "What is this ?"
	c.Website = "www.baidu.com"
	c.Details = "This is buidu website ?"
	if err := StakingInstance().EditCandidate(blockHash2, blockNumber2, c); nil != err {
		t.Error("Failed to EditCandidate", err)
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 err: %v", err)
		return
	}

	// get Candidate info after edit
	if can, err := getCandidate(blockHash2, index); nil != err {
		t.Errorf("Failed to Get candidate info, err: %v", err)
		return
	} else {
		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		return
	}
}

func TestStakingPlugin_IncreaseStaking(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}

	var c *staking.Candidate
	// Get Candidate Info
	if can, err := getCandidate(blockHash, index); nil != err {
		t.Errorf("Failed to Get candidate info, err: %v", err)
		return
	} else {
		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		c = can
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock2 err", err)
		return
	}

	/**
	Start IncreaseStaking
	*/
	if err := StakingInstance().IncreaseStaking(state, blockHash2, blockNumber2, common.Big256, uint16(0), c); nil != err {
		t.Error("Failed to IncreaseStaking", err)
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 err: %v", err)
		return
	}

	// get Candidate info
	addr, _ := xutil.NodeId2Addr(nodeIdArr[index])
	if can, err := StakingInstance().GetCandidateInfoByIrr(addr); nil != err {
		//if can, err := StakingInstance().GetCandidateInfo(blockHash2, addr); nil != err {
		t.Error("Failed to Get Candidate info After Increase", err)
	} else {

		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info After Increase is:", string(canByte))
		c = can
	}

}

func TestStakingPlugin_WithdrewCandidate(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	var c *staking.Candidate
	// Get Candidate Info
	if can, err := getCandidate(blockHash, index); nil != err {
		t.Errorf("Failed to Get candidate info, err: %v", err)
		return
	} else {
		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		c = can
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock2 err", err)
		return
	}

	/**
	Start WithdrewStaking
	*/
	if err := StakingInstance().WithdrewStaking(state, blockHash2, blockNumber2, c); nil != err {
		t.Error("Failed to WithdrewStaking", err)
		return
	}

	t.Log("Finish WithdrewStaking ~~")
	// get Candidate info
	if _, err := getCandidate(blockHash2, index); nil != err && err == snapshotdb.ErrNotFound {
		t.Logf("expect candidate info is no found, err: %v", err)
		return
	} else {
		t.Error("It is not expect~")
		return
	}

}

func TestStakingPlugin_HandleUnCandidateItem(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()
	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	// Add UNStakingItems
	//stakingDB := staking.NewStakingDB()

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())
	canAddr, _ := xutil.NodeId2Addr(nodeIdArr[index])

	if err := StakingInstance().addUnStakeItem(state, blockNumber.Uint64(), blockHash, epoch, nodeIdArr[index], canAddr, blockNumber.Uint64()); nil != err {
		t.Error("Failed to AddUnStakeItemStore:", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}

	// Get Candidate Info
	if can, err := getCandidate(blockHash, index); nil != err {
		t.Errorf("Failed to Get candidate info, err: %v", err)
		return
	} else {
		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock2 err", err)
		return
	}

	/**
	Start HandleUnCandidateItem
	*/
	err = StakingInstance().HandleUnCandidateItem(state, blockNumber2.Uint64(), blockHash2, epoch+xcom.UnStakeFreezeRatio())
	if nil != err {
		t.Error("Failed to HandleUnCandidateItem:", err)
		return
	}

	t.Log("Finish HandleUnCandidateItem ~~")

	// get Candidate info
	if _, err := getCandidate(blockHash2, index); nil != err && err == snapshotdb.ErrNotFound {
		t.Logf("expect candidate info is no found, err: %v", err)
		return
	} else {
		t.Error("It is not expect~")
		return
	}

}

func TestStakingPlugin_Delegate(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()

	defer func() {
		sndb.Clear()
	}()
	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	var c *staking.Candidate
	// Get Candidate Info
	if can, err := getCandidate(blockHash, index); nil != err {
		t.Errorf("Failed to Get candidate info, err: %v", err)
		return
	} else {
		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		c = can
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	/**
	Start Delegate
	*/
	_, err = delegate(state, blockHash2, blockNumber2, c, 0, index, t)
	if nil != err {
		t.Error("Failed to Delegate:", err)
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)

	}
	t.Log("Finish Delegate ~~")
	if can, err := getCandidate(blockHash2, index); nil != err {
		t.Errorf("Failed to Get candidate info, err: %v", err)
		return
	} else {
		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		return
	}

}

func TestStakingPlugin_WithdrewDelegate(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()
	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	var c *staking.Candidate
	// Get Candidate Info
	if can, err := getCandidate(blockHash, index); nil != err {
		t.Errorf("Failed to Get candidate info, err: %v", err)
		return
	} else {
		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		c = can
	}

	// Delegate
	del, err := delegate(state, blockHash, blockNumber, c, 0, index, t)
	if nil != err {
		t.Error("Failed to Delegate:", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	t.Log("Finish delegate ~~")
	if can, err := getCandidate(blockHash, index); nil != err {
		t.Errorf("Failed to Get candidate info, err: %v", err)
		return
	} else {
		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	/**
	Start Withdrew Delegate
	*/
	err = StakingInstance().WithdrewDelegate(state, blockHash2, blockNumber2, common.Big257, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del)
	if nil != err {
		t.Error("Failed to WithdrewDelegate:", err)
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
	}
	t.Log("Finish WithdrewDelegate ~~")
	if can, err := getCandidate(blockHash2, index); nil != err {
		t.Errorf("Failed to Get candidate info, err: %v", err)
		return
	} else {
		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		return
	}
}

func TestStakingPlugin_GetDelegateInfo(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()
	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	t.Log("Finish delegate ~~")

	var c *staking.Candidate

	if can, err := getCandidate(blockHash, index); nil != err {
		t.Errorf("Failed to Get candidate info, err: %v", err)
		return
	} else {
		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		c = can
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	// Delegate
	_, err = delegate(state, blockHash2, blockNumber2, c, 0, index, t)
	if nil != err {
		t.Error("Failed to Delegate:", err)
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
		return
	}

	t.Log("Finished Delegate ~~")
	// get Delegate info
	getDelegate(blockHash2, blockNumber.Uint64(), index, t)
}

func TestStakingPlugin_GetDelegateInfoByIrr(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	var c *staking.Candidate
	if can, err := getCandidate(blockHash, index); nil != err {
		t.Errorf("Failed to Get candidate info, err: %v", err)
		return
	} else {
		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		c = can
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	t.Log("Start delegate ~~")
	// Delegate
	_, err = delegate(state, blockHash2, blockNumber2, c, 0, index, t)
	if nil != err {
		t.Error("Failed to Delegate:", err)
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
		return
	}

	t.Log("Finished Delegate ~~")
	/**
	Start get Delegate info
	*/
	del, err := StakingInstance().GetDelegateInfoByIrr(addrArr[index+1], nodeIdArr[index], blockNumber.Uint64())
	if nil != err {
		t.Error("Failed to GetDelegateInfoByIrr:", err)
		return
	}

	delByte, _ := json.Marshal(del)
	t.Log("Get Delegate is:", string(delByte))

}

func TestStakingPlugin_GetRelatedListByDelAddr(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()
	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	// staking 0, 1, 2, 3
	for i := 0; i < 4; i++ {
		if err := create_staking(state, blockNumber, blockHash, i, 0, t); nil != err {
			t.Error("Failed to Create Staking", err)
			return
		}
	}

	t.Log("First delegate ~~")
	for i := 0; i < 2; i++ {
		// 0, 1
		var c *staking.Candidate

		if can, err := getCandidate(blockHash, i); nil != err {
			t.Errorf("Failed to Get candidate info, err: %v", err)
			return
		} else {
			canByte, _ := json.Marshal(can)
			t.Log("Get Candidate Info is:", string(canByte))
			c = can
		}
		// Delegate  0, 1
		_, err := delegate(state, blockHash, blockNumber, c, 0, i, t)
		if nil != err {
			t.Errorf("Failed to Delegate: Num: %d, error: %v", i, err)
			return
		}

		//t.Log("First: Del => Can:", addrArr[i+1].Hex(), c.NodeId.String(), c.StakingBlockNum)

	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	t.Log("Second delegate ~~")
	for i := 1; i < 3; i++ {
		// 0, 1
		var c *staking.Candidate
		if can, err := getCandidate(blockHash2, i-1); nil != err {
			t.Errorf("Failed to Get candidate info, err: %v", err)
			return
		} else {
			canByte, _ := json.Marshal(can)
			t.Log("Get Candidate Info is:", string(canByte))
			c = can
		}

		// Delegate
		_, err := delegate(state, blockHash2, blockNumber2, c, 0, i, t)
		if nil != err {
			t.Errorf("Failed to Delegate: Num: %d, error: %v", i, err)
			return
		}

		//t.Log("Second: Del => Can:", addrArr[i+1].Hex(), c.NodeId.String(), c.StakingBlockNum)

	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
		return
	}

	t.Log("Finished Delegate ~~")
	/**
	Start get RelatedList
	*/
	rel, err := StakingInstance().GetRelatedListByDelAddr(blockHash2, addrArr[1+1])
	if nil != err {
		t.Error("Failed to GetRelatedListByDelAddr:", err)
		return
	}

	relByte, _ := json.Marshal(rel)
	t.Log("Get RelatedList is:", string(relByte))
}

func TestStakingPlugin_HandleUnDelegateItem(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	var c *staking.Candidate

	if can, err := getCandidate(blockHash, index); nil != err {
		t.Errorf("Failed to Get candidate info, err: %v", err)
		return
	} else {
		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		c = can
	}

	// Delegate
	_, err = delegate(state, blockHash, blockNumber, c, 0, index, t)
	if nil != err {
		t.Error("Failed to Delegate:", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	t.Log("Finished Delegate ~~")
	// get Delegate info
	del := getDelegate(blockHash2, blockNumber.Uint64(), index, t)

	// Add UnDelegateItem
	stakingDB := staking.NewStakingDB()

	epoch := xutil.CalculateEpoch(blockNumber2.Uint64())

	amount := new(big.Int).Add(c.Released, c.RestrictingPlan)

	delAddr := addrArr[index+1]

	err = StakingInstance().addUnDelegateItem(blockNumber2.Uint64(), blockHash2, delAddr, c.NodeId, epoch, c.StakingBlockNum, amount)
	if nil != err {
		t.Error("Failed to AddUnDelegateItemStore:", err)
		return
	}
	del.Reduction = new(big.Int).Add(del.Reduction, amount)
	// update del
	if err := stakingDB.SetDelegateStore(blockHash2, delAddr, c.NodeId, c.StakingBlockNum, del); nil != err {
		t.Error("Failed to Update Delegate When AddUnDelegateItemStore:", err)
		return
	}

	if err := stakingDB.DelCanPowerStore(blockHash2, c); nil != err {
		t.Error("Failed to DelCanPowerStore:", err)
		return
	}

	// change candidate shares
	c.Shares = new(big.Int).Sub(c.Shares, amount)

	canAddr, _ := xutil.NodeId2Addr(c.NodeId)

	if err := stakingDB.SetCandidateStore(blockHash2, canAddr, c); nil != err {
		t.Error("Failed to SetCandidateStore:", err)
		return
	}

	if err := stakingDB.SetCanPowerStore(blockHash2, canAddr, c); nil != err {
		t.Error("Failed to SetCanPowerStore:", err)
		return
	}

	/**
	Start HandleUnDelegateItem
	*/
	err = StakingInstance().HandleUnDelegateItem(state, blockNumber2.Uint64(), blockHash2, epoch+xcom.ActiveUnDelFreezeRatio())
	if nil != err {
		t.Error("Failed to HandleUnDelegateItem:", err)
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
		return
	}

	t.Log("Finished HandleUnDelegateItem ~~")

	// get Candiddate
	if can, err := getCandidate(blockHash2, index); nil != err {
		t.Errorf("Failed to Get candidate info, err: %v", err)
		return
	} else {
		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		c = can
	}

	// get Delegate
	getDelegate(blockHash2, c.StakingBlockNum, index, t)

}

func TestStakingPlugin_ElectNextVerifierList(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),

			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList
	start := uint64(1)
	end := xutil.EpochSize() * xutil.ConsensusSize()

	new_verifierArr := &staking.ValidatorArray{
		Start: start,
		End:   end,
	}

	queue := make(staking.ValidatorQueue, 0)

	iter := sndb.Ranking(blockHash, staking.CanPowerKeyPrefix, 0)
	if err := iter.Error(); nil != err {
		t.Errorf("Failed to build genesis VerifierList, the iter is  err: %v", err)
		return
	}

	defer iter.Release()

	// for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {

	count := 0
	for iter.Valid(); iter.Next(); {
		if uint64(count) == xcom.EpochValidatorNum() {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProgramVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        can.NodeId,
			BlsPubKey:     can.BlsPubKey,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
		count++
	}

	new_verifierArr.Arr = queue

	err = setVerifierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	/*
		Start ElectNextVerifierList
	*/
	targetNum := xutil.EpochSize() * xutil.ConsensusSize()
	fmt.Println("targetNum:", targetNum)

	targetNumInt := big.NewInt(int64(targetNum))

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	err = StakingInstance().ElectNextVerifierList(blockHash2, targetNumInt.Uint64(), state)
	if nil != err {
		t.Errorf("Failed to ElectNextVerifierList, err: %v", err)
	}

}

func TestStakingPlugin_Election(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// Must new VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList

	start := uint64(1)
	end := xutil.EpochSize() * xutil.ConsensusSize()

	new_verifierArr := &staking.ValidatorArray{
		Start: start,
		End:   end,
	}

	queue := make(staking.ValidatorQueue, 0)

	iter := sndb.Ranking(blockHash, staking.CanPowerKeyPrefix, 0)
	if err := iter.Error(); nil != err {
		t.Errorf("Failed to build genesis VerifierList, the iter is  err: %v", err)
		return
	}

	defer iter.Release()

	count := 0
	for iter.Valid(); iter.Next(); {
		if uint64(count) == xcom.EpochValidatorNum() {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProgramVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        can.NodeId,
			BlsPubKey:     can.BlsPubKey,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
		count++
	}

	new_verifierArr.Arr = queue

	err = setVerifierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}

	// build gensis current validatorList
	new_validatorArr := &staking.ValidatorArray{
		Start: start,
		End:   xutil.ConsensusSize(),
	}

	new_validatorArr.Arr = queue[:int(xcom.ConsValidatorNum())]

	err = setRoundValList(blockHash, new_validatorArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis current round validatorList, err: %v", err)
		return
	}

	// build ancestor nonces
	currNonce, nonces := build_vrf_Nonce()
	if enValue, err := rlp.EncodeToBytes(nonces); nil != err {
		t.Error("Storage previous nonce failed", "err", err)
		return
	} else {
		sndb.Put(blockHash, handler.NonceStorageKey, enValue)
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	/*
		Start Election
	*/
	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	header := &types.Header{
		ParentHash: blockHash,
		Number:     big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())),
		Nonce:      types.EncodeNonce(currNonce),
	}

	err = StakingInstance().Election(blockHash2, header, state)
	if nil != err {
		t.Errorf("Failed to Election, err: %v", err)
	}

}

func TestStakingPlugin_SlashCandidates(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	// Will be Slashing candidate
	slashQueue := make(staking.CandidateQueue, 2)

	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
		if i < len(slashQueue) {
			slashQueue[i] = canTmp
		}
	}

	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList

	start := uint64(1)
	end := xutil.EpochSize() * xutil.ConsensusSize()

	new_verifierArr := &staking.ValidatorArray{
		Start: start,
		End:   end,
	}

	queue := make(staking.ValidatorQueue, 0)

	iter := sndb.Ranking(blockHash, staking.CanPowerKeyPrefix, 0)
	if err := iter.Error(); nil != err {
		t.Errorf("Failed to build genesis VerifierList, the iter is  err: %v", err)
		return
	}

	defer iter.Release()

	// for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {

	count := 0
	for iter.Valid(); iter.Next(); {
		if uint64(count) == xcom.EpochValidatorNum() {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProgramVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        can.NodeId,
			BlsPubKey:     can.BlsPubKey,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
		count++
	}

	new_verifierArr.Arr = queue

	err = setVerifierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	/**
	Start SlashCandidates
	*/
	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock err", err)
		return
	}

	slash1 := slashQueue[0]
	slash2 := slashQueue[1]

	slashItem1 := &staking.SlashNodeItem{
		NodeId:      slash1.NodeId,
		Amount:      slash1.Released,
		SlashType:   staking.LowRatio,
		BenefitAddr: vm.RewardManagerPoolAddr,
	}

	sla := new(big.Int).Div(slash2.Released, big.NewInt(10))
	caller := common.HexToAddress("0xe4a22694827bFa617bF039c937403190477934bF")

	slashItem2 := &staking.SlashNodeItem{
		NodeId:      slash2.NodeId,
		Amount:      sla,
		SlashType:   staking.DuplicateSign,
		BenefitAddr: caller,
	}
	slashItemQueue := make(staking.SlashQueue, 0)
	slashItemQueue = append(slashItemQueue, slashItem1)
	slashItemQueue = append(slashItemQueue, slashItem2)

	err = StakingInstance().SlashCandidates(state, blockHash2, blockNumber2.Uint64(), slashItemQueue...)
	if nil != err {
		t.Errorf("Failed to SlashCandidates Second can (DuplicateSign), err: %v", err)
		return
	}

}

func TestStakingPlugin_DeclarePromoteNotify(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	handler.NewVrfHandler(genesis.Hash().Bytes())

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	queue := make(staking.CandidateQueue, 0)
	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}

		if i < 20 {
			queue = append(queue, canTmp)
		}
	}

	// Commit Block 1
	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	/**
	Start DeclarePromoteNotify
	*/
	for i, can := range queue {
		err = StakingInstance().DeclarePromoteNotify(blockHash2, blockNumber2.Uint64(), can.NodeId, promoteVersion)
		if nil != err {
			t.Errorf("Failed to DeclarePromoteNotify, index: %d, err: %v", i, err)
			return
		}
	}

}

func TestStakingPlugin_ProposalPassedNotify(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}

	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	handler.NewVrfHandler(genesis.Hash().Bytes())

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	validatorQueue := make(staking.ValidatorQueue, 0)

	nodeIdArr := make([]discover.NodeID, 0)
	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}

		if i < 20 {
			nodeIdArr = append(nodeIdArr, canTmp.NodeId)
		}

		v := &staking.Validator{
			NodeAddress: canAddr,
			NodeId:      canTmp.NodeId,
			BlsPubKey:   canTmp.BlsPubKey,
			StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
				fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
			ValidatorTerm: 0,
		}

		validatorQueue = append(validatorQueue, v)
	}

	epoch_Arr := &staking.ValidatorArray{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
		Arr:   validatorQueue,
	}

	curr_Arr := &staking.ValidatorArray{
		Start: 1,
		End:   xutil.ConsensusSize(),
		Arr:   validatorQueue,
	}

	t.Log("Store Curr Epoch VerifierList", "len", len(epoch_Arr.Arr))
	setVerifierList(blockHash, epoch_Arr)

	t.Log("Store CuRR Round Validator", "len", len(epoch_Arr.Arr))
	setRoundValList(blockHash, curr_Arr)

	// Commit Block 1
	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	/**
	Start ProposalPassedNotify
	*/
	err = StakingInstance().ProposalPassedNotify(blockHash2, blockNumber2.Uint64(), nodeIdArr, promoteVersion)
	if nil != err {
		t.Errorf("Failed to ProposalPassedNotify, err: %v", err)
	}
}

func TestStakingPlugin_GetCandidateONEpoch(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// New VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	// build vrf proof
	// build ancestor nonces
	_, nonces := build_vrf_Nonce()
	enValue, err := rlp.EncodeToBytes(nonces)
	if nil != err {
		t.Error("Failed to rlp vrf nonces", "err", err)
		return
	}

	// new block
	privateKey, err := crypto.GenerateKey()
	if nil != err {
		t.Errorf("Failed to generate random Address private key: %v", err)
		return
	}
	nodeId := discover.PubkeyID(&privateKey.PublicKey)
	currentHash := crypto.Keccak256Hash([]byte(nodeId.String()))
	currentNumber := big.NewInt(1)

	// build genesis veriferList and validatorList
	validatorQueue := make(staking.ValidatorQueue, xcom.EpochValidatorNum())

	for j := 0; j < 1000; j++ {
		var index int = j % 25

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(j),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(j) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(j) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		// Store Candidate power
		powerKey := staking.TallyPowerKey(canTmp.Shares, canTmp.StakingBlockNum, canTmp.StakingTxIndex, canTmp.ProgramVersion)
		if err := sndb.PutBaseDB(powerKey, canAddr.Bytes()); nil != err {
			t.Errorf("Failed to Store Candidate Power: PutBaseDB failed. error:%s", err.Error())
			return
		}

		// Store Candidate info
		canKey := staking.CandidateKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canTmp); nil != err {
			t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canKey, val); nil != err {
				t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		if j < int(xcom.EpochValidatorNum()) {
			v := &staking.Validator{
				NodeAddress: canAddr,
				NodeId:      canTmp.NodeId,
				BlsPubKey:   canTmp.BlsPubKey,
				StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
					fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
				ValidatorTerm: 0,
			}
			validatorQueue[j] = v
		}

	}

	/**
	*******
	build genesis epoch validators
	*******
	*/
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}

	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
		return
	}

	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	/**
	*******
	build genesis curr round validators
	*******
	*/
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
		return
	}

	xcom.PrintObject("Test round", validatorQueue[:xcom.ConsValidatorNum()])
	roundArr, err := rlp.EncodeToBytes(validatorQueue[:xcom.ConsValidatorNum()])
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Current Round validator
	if err := sndb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), roundArr); nil != err {
		t.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	// Store vrf nonces
	if err := sndb.PutBaseDB(handler.NonceStorageKey, enValue); nil != err {
		t.Errorf("Failed to Store Current Vrf nonces : PutBaseDB failed. error:%s", err.Error())
		return
	}

	// SetCurrent to snapshotDB
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50
	preNum1 := new(big.Int).Sub(currentNumber, big.NewInt(1))
	if err := sndb.SetCurrent(currentHash, *preNum1, *preNum1); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error()))
	}

	/**
	Start GetCandidateONEpoch
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50

	nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
	header := &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	fmt.Println("Start GetCandidateONEpoch CurrBlockHash", currentHash.Hex())

	/**
	Start GetCandidateONEpoch
	*/
	canQueue, err := StakingInstance().GetCandidateONEpoch(currentHash, currentNumber.Uint64(), QueryStartNotIrr)
	if nil != err {
		t.Errorf("Failed to GetCandidateONEpoch by QueryStartNotIrr, err: %v", err)
		return
	}

	canArr, _ := json.Marshal(canQueue)
	t.Log("GetCandidateONEpoch by QueryStartNotIrr:", string(canArr))

	canQueue, err = StakingInstance().GetCandidateONEpoch(currentHash, currentNumber.Uint64(), QueryStartIrr)
	if nil != err {
		t.Errorf("Failed to GetCandidateONEpoch by QueryStartIrr, err: %v", err)
	}
	canArr, _ = json.Marshal(canQueue)
	t.Log("GetCandidateONEpoch by QueryStartIrr:", string(canArr))
}

func TestStakingPlugin_GetCandidateONRound(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// New VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	// build vrf proof
	// build ancestor nonces
	_, nonces := build_vrf_Nonce()
	enValue, err := rlp.EncodeToBytes(nonces)
	if nil != err {
		t.Error("Failed to rlp vrf nonces", "err", err)
		return
	}

	// new block
	privateKey, err := crypto.GenerateKey()
	if nil != err {
		t.Errorf("Failed to generate random Address private key: %v", err)
		return
	}
	nodeId := discover.PubkeyID(&privateKey.PublicKey)
	currentHash := crypto.Keccak256Hash([]byte(nodeId.String()))
	currentNumber := big.NewInt(1)

	// build genesis veriferList and validatorList
	validatorQueue := make(staking.ValidatorQueue, xcom.EpochValidatorNum())

	for j := 0; j < 1000; j++ {
		var index int = j % 25

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(j),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(j) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(j) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		// Store Candidate power
		powerKey := staking.TallyPowerKey(canTmp.Shares, canTmp.StakingBlockNum, canTmp.StakingTxIndex, canTmp.ProgramVersion)
		if err := sndb.PutBaseDB(powerKey, canAddr.Bytes()); nil != err {
			t.Errorf("Failed to Store Candidate Power: PutBaseDB failed. error:%s", err.Error())
			return
		}

		// Store Candidate info
		canKey := staking.CandidateKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canTmp); nil != err {
			t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canKey, val); nil != err {
				t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		if j < int(xcom.EpochValidatorNum()) {
			v := &staking.Validator{
				NodeAddress: canAddr,
				NodeId:      canTmp.NodeId,
				BlsPubKey:   canTmp.BlsPubKey,
				StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
					fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
				ValidatorTerm: 0,
			}
			validatorQueue[j] = v
		}

	}

	/**
	*******
	build genesis epoch validators
	*******
	*/
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}

	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
		return
	}

	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	/**
	*******
	build genesis curr round validators
	*******
	*/
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
		return
	}

	xcom.PrintObject("Test round", validatorQueue[:xcom.ConsValidatorNum()])
	roundArr, err := rlp.EncodeToBytes(validatorQueue[:xcom.ConsValidatorNum()])
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Current Round validator
	if err := sndb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), roundArr); nil != err {
		t.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	// Store vrf nonces
	if err := sndb.PutBaseDB(handler.NonceStorageKey, enValue); nil != err {
		t.Errorf("Failed to Store Current Vrf nonces : PutBaseDB failed. error:%s", err.Error())
		return
	}

	// SetCurrent to snapshotDB
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50
	preNum1 := new(big.Int).Sub(currentNumber, big.NewInt(1))
	if err := sndb.SetCurrent(currentHash, *preNum1, *preNum1); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error()))
	}

	/**
	Start GetCandidateONRound
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50

	nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
	header := &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	fmt.Println("Start GetCandidateONRound CurrBlockHash", currentHash.Hex())

	/**
	Start GetCandidateONRound
	*/
	canQueue, err := StakingInstance().GetCandidateONRound(currentHash, currentNumber.Uint64(), CurrentRound, QueryStartNotIrr)
	if nil != err {
		t.Errorf("Failed to GetCandidateONRound by QueryStartNotIrr, err: %v", err)
		return
	}

	canArr, _ := json.Marshal(canQueue)
	t.Log("GetCandidateONRound by QueryStartNotIrr:", string(canArr))

	canQueue, err = StakingInstance().GetCandidateONRound(currentHash, currentNumber.Uint64(), CurrentRound, QueryStartIrr)
	if nil != err {
		t.Errorf("Failed to GetCandidateONRound by QueryStartIrr, err: %v", err)
		return
	}
	canArr, _ = json.Marshal(canQueue)
	t.Log("GetCandidateONRound by QueryStartIrr:", string(canArr))

}

func TestStakingPlugin_GetValidatorList(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// New VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	// build vrf proof
	// build ancestor nonces
	_, nonces := build_vrf_Nonce()
	enValue, err := rlp.EncodeToBytes(nonces)
	if nil != err {
		t.Error("Failed to rlp vrf nonces", "err", err)
		return
	}

	// new block
	privateKey, err := crypto.GenerateKey()
	if nil != err {
		t.Errorf("Failed to generate random Address private key: %v", err)
		return
	}
	nodeId := discover.PubkeyID(&privateKey.PublicKey)
	currentHash := crypto.Keccak256Hash([]byte(nodeId.String()))
	currentNumber := big.NewInt(1)

	// build genesis veriferList and validatorList
	validatorQueue := make(staking.ValidatorQueue, xcom.EpochValidatorNum())

	for j := 0; j < 1000; j++ {
		var index int = j % 25

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(j),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(j) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(j) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		// Store Candidate power
		powerKey := staking.TallyPowerKey(canTmp.Shares, canTmp.StakingBlockNum, canTmp.StakingTxIndex, canTmp.ProgramVersion)
		if err := sndb.PutBaseDB(powerKey, canAddr.Bytes()); nil != err {
			t.Errorf("Failed to Store Candidate Power: PutBaseDB failed. error:%s", err.Error())
			return
		}

		// Store Candidate info
		canKey := staking.CandidateKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canTmp); nil != err {
			t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canKey, val); nil != err {
				t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		if j < int(xcom.EpochValidatorNum()) {
			v := &staking.Validator{
				NodeAddress: canAddr,
				NodeId:      canTmp.NodeId,
				BlsPubKey:   canTmp.BlsPubKey,
				StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
					fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
				ValidatorTerm: 0,
			}
			validatorQueue[j] = v
		}

	}

	/**
	*******
	build genesis epoch validators
	*******
	*/
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}

	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
		return
	}

	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	/**
	*******
	build genesis curr round validators
	*******
	*/
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
		return
	}

	xcom.PrintObject("Test round", validatorQueue[:xcom.ConsValidatorNum()])
	roundArr, err := rlp.EncodeToBytes(validatorQueue[:xcom.ConsValidatorNum()])
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Current Round validator
	if err := sndb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), roundArr); nil != err {
		t.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	// Store vrf nonces
	if err := sndb.PutBaseDB(handler.NonceStorageKey, enValue); nil != err {
		t.Errorf("Failed to Store Current Vrf nonces : PutBaseDB failed. error:%s", err.Error())
		return
	}

	// SetCurrent to snapshotDB
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50
	preNum1 := new(big.Int).Sub(currentNumber, big.NewInt(1))
	if err := sndb.SetCurrent(currentHash, *preNum1, *preNum1); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error()))
	}

	/**
	Start GetValidatorList
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50

	nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
	header := &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	fmt.Println("Start GetValidatorList CurrBlockHash", currentHash.Hex())

	/**
	Start  GetValidatorList
	*/
	validatorExQueue, err := StakingInstance().GetValidatorList(currentHash, currentNumber.Uint64(), CurrentRound, QueryStartNotIrr)
	if nil != err {
		t.Errorf("Failed to GetValidatorList by QueryStartNotIrr, err: %v", err)
		return
	}

	validatorExArr, _ := json.Marshal(validatorExQueue)
	t.Log("GetValidatorList by QueryStartNotIrr:", string(validatorExArr))

	validatorExQueue, err = StakingInstance().GetValidatorList(currentHash, currentNumber.Uint64(), CurrentRound, QueryStartIrr)
	if nil != err {
		t.Errorf("Failed to GetValidatorList by QueryStartIrr, err: %v", err)
		return
	}
	validatorExArr, _ = json.Marshal(validatorExQueue)
	t.Log("GetValidatorList by QueryStartIrr:", string(validatorExArr))

	t.Log("validatorList[0]", validatorExQueue[0])

}

func TestStakingPlugin_GetHistoryValidatorList(t *testing.T) {

	//defer plugin.ClearStakingPlugin()
	//defer plugin.ClearGovPlugin()

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	xcom.NewVrfHandler(genesis.Hash().Bytes())

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	stakingDB := staking.NewStakingDB()

	//
	headerMap := make(map[int]*types.Header, 0)
	switchNum := int(xutil.CalcBlocksEachEpoch())
	parentHash := genesis.Hash()
	for i := 0; i <= int(xutil.CalcBlocksEachEpoch()); i++ {

		nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(i))))[:]
		privateKey, err := crypto.GenerateKey()
		if nil != err {
			fmt.Printf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		root := crypto.Keccak256Hash([]byte(nodeId.String()))

		blockNum := big.NewInt(int64(i + 1))

		header := &types.Header{
			ParentHash:  parentHash,
			Coinbase:    sender,
			Root:        root,
			TxHash:      types.EmptyRootHash,
			ReceiptHash: types.EmptyRootHash,
			Number:      blockNum,
			Time:        big.NewInt(int64(121321213 * i)),
			Extra:       make([]byte, 97),
			Nonce:       types.EncodeNonce(nonce),
		}

		curr_Hash := header.Hash()

		if err := sndb.NewBlock(blockNum, parentHash, curr_Hash); nil != err {
			t.Errorf("Failed to snapshotDB New Block, err: %v", err)
			return
		}

		// Create Staking
		if i == 0 {

			validatorQueue := make(staking.ValidatorQueue, 101)

			for j := 0; j < 101; j++ {
				var index int
				if j >= len(balanceStr) {
					index = j % (len(balanceStr) - 1)
				}

				balance, _ := new(big.Int).SetString(balanceStr[index], 10)

				mrand.Seed(time.Now().UnixNano())

				weight := mrand.Intn(1000000000)

				ii := mrand.Intn(len(chaList))

				balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

				privateKey, err := crypto.GenerateKey()
				if nil != err {
					t.Errorf("Failed to generate random NodeId private key: %v", err)
					return
				}

				nodeId := discover.PubkeyID(&privateKey.PublicKey)

				privateKey, err = crypto.GenerateKey()
				if nil != err {
					t.Errorf("Failed to generate random Address private key: %v", err)
					return
				}

				addr := crypto.PubkeyToAddress(privateKey.PublicKey)

				canTmp := &staking.Candidate{
					NodeId:          nodeId,
					StakingAddress:  sender,
					BenefitAddress:  addr,
					StakingBlockNum: uint64(1 + i),
					StakingTxIndex:  uint32(index),
					Shares:          balance,
					ProgramVersion:  xutil.CalcVersion(initProgramVersion),
					// Prevent null pointer initialization
					Released:           common.Big0,
					ReleasedHes:        common.Big0,
					RestrictingPlan:    common.Big0,
					RestrictingPlanHes: common.Big0,

					Description: staking.Description{
						NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
						ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
						Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
						Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
					},
				}

				canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

				stakingDB.SetCanPowerStore(curr_Hash, canAddr, canTmp)
				stakingDB.SetCandidateStore(curr_Hash, canAddr, canTmp)

				v := &staking.Validator{
					NodeAddress: canAddr,
					NodeId:      canTmp.NodeId,
					StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
						fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
					ValidatorTerm: 0,
				}
				validatorQueue[j] = v
			}

			epoch_Arr := &staking.Validator_array{
				Start: 1,
				End:   xutil.CalcBlocksEachEpoch(),
				Arr:   validatorQueue,
			}
			// start := old_verifierArr.End + 1
			//	end := old_verifierArr.End + xutil.CalcBlocksEachEpoch()

			curr_Arr := &staking.Validator_array{
				Start: 1,
				End:   xutil.ConsensusSize(),
				Arr:   validatorQueue,
			}

			// add Current Validators And Epoch Validators
			setVerifierList(curr_Hash, epoch_Arr)

			setRoundValList(curr_Hash, curr_Arr)

		} else {

			var index int
			if i >= len(balanceStr) {
				index = i % (len(balanceStr) - 1)
			}

			balance, _ := new(big.Int).SetString(balanceStr[index], 10)

			mrand.Seed(time.Now().UnixNano())

			weight := mrand.Intn(1000000000)

			ii := mrand.Intn(len(chaList))

			balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

			privateKey, err := crypto.GenerateKey()
			if nil != err {
				t.Errorf("Failed to generate random NodeId private key: %v", err)
				return
			}

			nodeId := discover.PubkeyID(&privateKey.PublicKey)

			privateKey, err = crypto.GenerateKey()
			if nil != err {
				t.Errorf("Failed to generate random Address private key: %v", err)
				return
			}

			addr := crypto.PubkeyToAddress(privateKey.PublicKey)

			canTmp := &staking.Candidate{
				NodeId:          nodeId,
				StakingAddress:  sender,
				BenefitAddress:  addr,
				StakingBlockNum: uint64(i + 1),
				StakingTxIndex:  uint32(index),
				Shares:          balance,
				ProgramVersion:  xutil.CalcVersion(initProgramVersion),
				// Prevent null pointer initialization
				Released:           common.Big0,
				ReleasedHes:        common.Big0,
				RestrictingPlan:    common.Big0,
				RestrictingPlanHes: common.Big0,

				Description: staking.Description{
					NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
					ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
					Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
					Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
				},
			}

			canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

			stakingDB.SetCanPowerStore(curr_Hash, canAddr, canTmp)
			stakingDB.SetCandidateStore(curr_Hash, canAddr, canTmp)

			// build vrf
			if i+1 == 229 {
				// build ancestor nonces
				_, nonces := build_vrf_Nonce()
				if enValue, err := rlp.EncodeToBytes(nonces); nil != err {
					t.Error("Storage previous nonce failed", "num", i+1, "Hash", curr_Hash.Hex(), "err", err)
					return
				} else {
					sndb.Put(curr_Hash, xcom.NonceStorageKey, enValue)
				}
			}

			// TODO Must be this
			if xutil.IsElection(header.Number.Uint64()) {
				err = plugin.StakingInstance().Election(curr_Hash, header)
				if nil != err {
					t.Errorf("Failed to Election, num:%d, Hash: %s, err: %v", header.Number.Uint64(), header.Hash().Hex(), err)
					return
				}
			}
		}

		// SnapshotDB  Commit
		if err := sndb.Commit(curr_Hash); nil != err {
			t.Errorf("Failed to snapshotDB Commit, err: %v", err)
			return
		}

		if i+1 == switchNum {
			headerMap[i+1] = header
		}
		parentHash = curr_Hash
	}

	// build genesis VerifierList

	// Request th opening/creation of an ephemeral database and ensure it's not persisted
	ctx := node.NewServiceContext(&node.Config{DataDir: ""}, nil, new(event.TypeMux), nil)
	config := &eth.Config{
	}
	hDB, _ := eth.CreateDB(ctx, config, "historydata")
	plugin.STAKING_DB = &plugin.StakingDB{
		HistoryDB:  hDB,
	}

	blockSwitch := types.NewBlock(headerMap[switchNum], nil, nil)
	//blockElection := types.NewBlock(headerMap[electionNum], nil, nil)
	err = plugin.StakingInstance().Confirmed(blockSwitch)
	if nil != err {
		return
	}

	/**
	Start GetVerifierList
	*/
	validatorExQueue, err := plugin.StakingInstance().GetHistoryValidatorList(blockHash2, headerMap[switchNum].Number.Uint64(), plugin.CurrentRound, plugin.QueryStartIrr)
	if nil != err {
		t.Errorf("Failed to GetHistoryValidatorList by QueryStartNotIrr, err: %v", err)
		return
	}

	validatorExArr, _ := json.Marshal(validatorExQueue)
	t.Log("GetHistoryValidatorList by QueryStartNotIrr:", string(validatorExArr))


}

func TestStakingPlugin_GetVerifierList(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// New VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	// build vrf proof
	// build ancestor nonces
	_, nonces := build_vrf_Nonce()
	enValue, err := rlp.EncodeToBytes(nonces)
	if nil != err {
		t.Error("Failed to rlp vrf nonces", "err", err)
		return
	}

	// new block
	privateKey, err := crypto.GenerateKey()
	if nil != err {
		t.Errorf("Failed to generate random Address private key: %v", err)
		return
	}
	nodeId := discover.PubkeyID(&privateKey.PublicKey)
	currentHash := crypto.Keccak256Hash([]byte(nodeId.String()))
	currentNumber := big.NewInt(1)

	// build genesis veriferList and validatorList
	validatorQueue := make(staking.ValidatorQueue, xcom.EpochValidatorNum())

	for j := 0; j < 1000; j++ {
		var index int = j % 25

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(j),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(j) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(j) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		// Store Candidate power
		powerKey := staking.TallyPowerKey(canTmp.Shares, canTmp.StakingBlockNum, canTmp.StakingTxIndex, canTmp.ProgramVersion)
		if err := sndb.PutBaseDB(powerKey, canAddr.Bytes()); nil != err {
			t.Errorf("Failed to Store Candidate Power: PutBaseDB failed. error:%s", err.Error())
			return
		}

		// Store Candidate info
		canKey := staking.CandidateKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canTmp); nil != err {
			t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canKey, val); nil != err {
				t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		if j < int(xcom.EpochValidatorNum()) {
			v := &staking.Validator{
				NodeAddress: canAddr,
				NodeId:      canTmp.NodeId,
				BlsPubKey:   canTmp.BlsPubKey,
				StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
					fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
				ValidatorTerm: 0,
			}
			validatorQueue[j] = v
		}

	}

	/**
	*******
	build genesis epoch validators
	*******
	*/
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}

	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
		return
	}

	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	/**
	*******
	build genesis curr round validators
	*******
	*/
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
		return
	}

	xcom.PrintObject("Test round", validatorQueue[:xcom.ConsValidatorNum()])
	roundArr, err := rlp.EncodeToBytes(validatorQueue[:xcom.ConsValidatorNum()])
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Current Round validator
	if err := sndb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), roundArr); nil != err {
		t.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	// Store vrf nonces
	if err := sndb.PutBaseDB(handler.NonceStorageKey, enValue); nil != err {
		t.Errorf("Failed to Store Current Vrf nonces : PutBaseDB failed. error:%s", err.Error())
		return
	}

	// SetCurrent to snapshotDB
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50
	preNum1 := new(big.Int).Sub(currentNumber, big.NewInt(1))
	if err := sndb.SetCurrent(currentHash, *preNum1, *preNum1); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error()))
	}

	/**
	Start GetVerifierList
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50

	nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
	header := &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	fmt.Println("Start GetVerifierLIst CurrBlockHash", currentHash.Hex())

	/**
	Start GetVerifierList
	*/
	validatorExQueue, err := StakingInstance().GetVerifierList(currentHash, currentNumber.Uint64(), QueryStartNotIrr)
	if nil != err {
		t.Errorf("Failed to GetVerifierList by QueryStartNotIrr, err: %v", err)
		return
	}

	validatorExArr, _ := json.Marshal(validatorExQueue)
	t.Log("GetVerifierList by QueryStartNotIrr:", string(validatorExArr))

	validatorExQueue, err = StakingInstance().GetVerifierList(currentHash, currentNumber.Uint64(), QueryStartIrr)
	if nil != err {
		t.Errorf("Failed to GetVerifierList by QueryStartIrr, err: %v", err)
		return
	}
	validatorExArr, _ = json.Marshal(validatorExQueue)
	t.Log("GetVerifierList by QueryStartIrr:", string(validatorExArr))

}

func TestStakingPlugin_GetHistoryVerifierList(t *testing.T) {

	//defer plugin.ClearStakingPlugin()
	//defer plugin.ClearGovPlugin()

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	xcom.NewVrfHandler(genesis.Hash().Bytes())

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	stakingDB := staking.NewStakingDB()

	//
	headerMap := make(map[int]*types.Header, 0)
	switchNum := int(xutil.ConsensusSize() - xcom.ElectionDistance())
	parentHash := genesis.Hash()
	for i := 0; i <= int(xutil.ConsensusSize()); i++ {

		nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(i))))[:]
		privateKey, err := crypto.GenerateKey()
		if nil != err {
			fmt.Printf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		root := crypto.Keccak256Hash([]byte(nodeId.String()))

		blockNum := big.NewInt(int64(i + 1))

		header := &types.Header{
			ParentHash:  parentHash,
			Coinbase:    sender,
			Root:        root,
			TxHash:      types.EmptyRootHash,
			ReceiptHash: types.EmptyRootHash,
			Number:      blockNum,
			Time:        big.NewInt(int64(121321213 * i)),
			Extra:       make([]byte, 97),
			Nonce:       types.EncodeNonce(nonce),
		}

		curr_Hash := header.Hash()

		if err := sndb.NewBlock(blockNum, parentHash, curr_Hash); nil != err {
			t.Errorf("Failed to snapshotDB New Block, err: % v", err)
			return
		}

		// Create Staking
		if i == 0 {

			validatorQueue := make(staking.ValidatorQueue, 101)

			for j := 0; j < 101; j++ {
				var index int
				if j >= len(balanceStr) {
					index = j % (len(balanceStr) - 1)
				}

				balance, _ := new(big.Int).SetString(balanceStr[index], 10)

				mrand.Seed(time.Now().UnixNano())

				weight := mrand.Intn(1000000000)

				ii := mrand.Intn(len(chaList))

				balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

				privateKey, err := crypto.GenerateKey()
				if nil != err {
					t.Errorf("Failed to generate random NodeId private key: %v", err)
					return
				}

				nodeId := discover.PubkeyID(&privateKey.PublicKey)

				privateKey, err = crypto.GenerateKey()
				if nil != err {
					t.Errorf("Failed to generate random Address private key: %v", err)
					return
				}

				addr := crypto.PubkeyToAddress(privateKey.PublicKey)

				canTmp := &staking.Candidate{
					NodeId:          nodeId,
					StakingAddress:  sender,
					BenefitAddress:  addr,
					StakingBlockNum: uint64(1 + i),
					StakingTxIndex:  uint32(index),
					Shares:          balance,
					ProgramVersion:  xutil.CalcVersion(initProgramVersion),
					// Prevent null pointer initialization
					Released:           common.Big0,
					ReleasedHes:        common.Big0,
					RestrictingPlan:    common.Big0,
					RestrictingPlanHes: common.Big0,

					Description: staking.Description{
						NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
						ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
						Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
						Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
					},
				}

				canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

				stakingDB.SetCanPowerStore(curr_Hash, canAddr, canTmp)
				stakingDB.SetCandidateStore(curr_Hash, canAddr, canTmp)

				v := &staking.Validator{
					NodeAddress: canAddr,
					NodeId:      canTmp.NodeId,
					StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
						fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
					ValidatorTerm: 0,
				}
				validatorQueue[j] = v
			}

			epoch_Arr := &staking.Validator_array{
				Start: 1,
				End:   xutil.CalcBlocksEachEpoch(),
				Arr:   validatorQueue,
			}
			// start := old_verifierArr.End + 1
			//	end := old_verifierArr.End + xutil.CalcBlocksEachEpoch()

			curr_Arr := &staking.Validator_array{
				Start: 1,
				End:   xutil.ConsensusSize(),
				Arr:   validatorQueue,
			}

			// add Current Validators And Epoch Validators
			setVerifierList(curr_Hash, epoch_Arr)

			setRoundValList(curr_Hash, curr_Arr)

		} else {

			var index int
			if i >= len(balanceStr) {
				index = i % (len(balanceStr) - 1)
			}

			balance, _ := new(big.Int).SetString(balanceStr[index], 10)

			mrand.Seed(time.Now().UnixNano())

			weight := mrand.Intn(1000000000)

			ii := mrand.Intn(len(chaList))

			balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

			privateKey, err := crypto.GenerateKey()
			if nil != err {
				t.Errorf("Failed to generate random NodeId private key: %v", err)
				return
			}

			nodeId := discover.PubkeyID(&privateKey.PublicKey)

			privateKey, err = crypto.GenerateKey()
			if nil != err {
				t.Errorf("Failed to generate random Address private key: %v", err)
				return
			}

			addr := crypto.PubkeyToAddress(privateKey.PublicKey)

			canTmp := &staking.Candidate{
				NodeId:          nodeId,
				StakingAddress:  sender,
				BenefitAddress:  addr,
				StakingBlockNum: uint64(i + 1),
				StakingTxIndex:  uint32(index),
				Shares:          balance,
				ProgramVersion:  xutil.CalcVersion(initProgramVersion),
				// Prevent null pointer initialization
				Released:           common.Big0,
				ReleasedHes:        common.Big0,
				RestrictingPlan:    common.Big0,
				RestrictingPlanHes: common.Big0,

				Description: staking.Description{
					NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
					ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
					Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
					Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
				},
			}

			canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

			stakingDB.SetCanPowerStore(curr_Hash, canAddr, canTmp)
			stakingDB.SetCandidateStore(curr_Hash, canAddr, canTmp)

			// build vrf
			if i+1 == 229 {
				// build ancestor nonces
				_, nonces := build_vrf_Nonce()
				if enValue, err := rlp.EncodeToBytes(nonces); nil != err {
					t.Error("Storage previous nonce failed", "num", i+1, "Hash", curr_Hash.Hex(), "err", err)
					return
				} else {
					sndb.Put(curr_Hash, xcom.NonceStorageKey, enValue)
				}
			}

			// TODO Must be this
			if xutil.IsElection(header.Number.Uint64()) {
				err = plugin.StakingInstance().Election(curr_Hash, header)
				if nil != err {
					t.Errorf("Failed to Election, num:%d, Hash: %s, err: %v", header.Number.Uint64(), header.Hash().Hex(), err)
					return
				}
			}
		}

		// SnapshotDB  Commit
		if err := sndb.Commit(curr_Hash); nil != err {
			t.Errorf("Failed to snapshotDB Commit, err: %v", err)
			return
		}

		if i+1 == switchNum {
			headerMap[i+1] = header
		}
		parentHash = curr_Hash
	}

	// build genesis VerifierList

	// Request th opening/creation of an ephemeral database and ensure it's not persisted
	ctx := node.NewServiceContext(&node.Config{DataDir: ""}, nil, new(event.TypeMux), nil)
	config := &eth.Config{
	}
	hDB, _ := eth.CreateDB(ctx, config, "historydata")
	plugin.STAKING_DB = &plugin.StakingDB{
		HistoryDB:  hDB,
	}

	blockSwitch := types.NewBlock(headerMap[switchNum], nil, nil)
	//blockElection := types.NewBlock(headerMap[electionNum], nil, nil)
	err = plugin.StakingInstance().Confirmed(blockSwitch)
	if nil != err {
		return
	}

	/**
	Start GetVerifierList
	*/
	validatorExQueue, err := plugin.StakingInstance().GetHistoryVerifierList(blockHash2, uint64(switchNum), plugin.QueryStartNotIrr)
	if nil != err {
		t.Errorf("Failed to GetHistoryVerifierList by QueryStartNotIrr, err: %v", err)
		return
	}

	validatorExArr, _ := json.Marshal(validatorExQueue)
	t.Log("GetHistoryVerifierList by QueryStartNotIrr:", string(validatorExArr))


}

func TestStakingPlugin_ListCurrentValidatorID(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// New VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	// build vrf proof
	// build ancestor nonces
	_, nonces := build_vrf_Nonce()
	enValue, err := rlp.EncodeToBytes(nonces)
	if nil != err {
		t.Error("Failed to rlp vrf nonces", "err", err)
		return
	}

	// new block
	privateKey, err := crypto.GenerateKey()
	if nil != err {
		t.Errorf("Failed to generate random Address private key: %v", err)
		return
	}
	nodeId := discover.PubkeyID(&privateKey.PublicKey)
	currentHash := crypto.Keccak256Hash([]byte(nodeId.String()))
	currentNumber := big.NewInt(1)

	// build genesis veriferList and validatorList
	validatorQueue := make(staking.ValidatorQueue, xcom.EpochValidatorNum())

	for j := 0; j < 1000; j++ {
		var index int = j % 25

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(j),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(j) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(j) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		// Store Candidate power
		powerKey := staking.TallyPowerKey(canTmp.Shares, canTmp.StakingBlockNum, canTmp.StakingTxIndex, canTmp.ProgramVersion)
		if err := sndb.PutBaseDB(powerKey, canAddr.Bytes()); nil != err {
			t.Errorf("Failed to Store Candidate Power: PutBaseDB failed. error:%s", err.Error())
			return
		}

		// Store Candidate info
		canKey := staking.CandidateKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canTmp); nil != err {
			t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canKey, val); nil != err {
				t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		if j < int(xcom.EpochValidatorNum()) {
			v := &staking.Validator{
				NodeAddress: canAddr,
				NodeId:      canTmp.NodeId,
				BlsPubKey:   canTmp.BlsPubKey,
				StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
					fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
				ValidatorTerm: 0,
			}
			validatorQueue[j] = v
		}

	}

	/**
	*******
	build genesis epoch validators
	*******
	*/
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}

	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
		return
	}

	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	/**
	*******
	build genesis curr round validators
	*******
	*/
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
		return
	}

	xcom.PrintObject("Test round", validatorQueue[:xcom.ConsValidatorNum()])
	roundArr, err := rlp.EncodeToBytes(validatorQueue[:xcom.ConsValidatorNum()])
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Current Round validator
	if err := sndb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), roundArr); nil != err {
		t.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	// Store vrf nonces
	if err := sndb.PutBaseDB(handler.NonceStorageKey, enValue); nil != err {
		t.Errorf("Failed to Store Current Vrf nonces : PutBaseDB failed. error:%s", err.Error())
		return
	}

	// SetCurrent to snapshotDB
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50
	preNum1 := new(big.Int).Sub(currentNumber, big.NewInt(1))
	if err := sndb.SetCurrent(currentHash, *preNum1, *preNum1); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error()))
	}

	/**
	Start ListCurrentValidaotorID
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50

	nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
	header := &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	fmt.Println("Start ListCurrentValidatorID CurrBlockHash", currentHash.Hex())

	/**
	Start  ListCurrentValidatorID
	*/
	validatorIdQueue, err := StakingInstance().ListCurrentValidatorID(currentHash, currentNumber.Uint64())
	if nil != err {
		t.Errorf("Failed to ListCurrentValidatorID, err: %v", err)
		return
	}

	validatorIdArr, _ := json.Marshal(validatorIdQueue)
	t.Log("ListCurrentValidatorID:", string(validatorIdArr))

}

func TestStakingPlugin_ListVerifierNodeID(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// New VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	// build vrf proof
	// build ancestor nonces
	_, nonces := build_vrf_Nonce()
	enValue, err := rlp.EncodeToBytes(nonces)
	if nil != err {
		t.Error("Failed to rlp vrf nonces", "err", err)
		return
	}

	// new block
	privateKey, err := crypto.GenerateKey()
	if nil != err {
		t.Errorf("Failed to generate random Address private key: %v", err)
		return
	}
	nodeId := discover.PubkeyID(&privateKey.PublicKey)
	currentHash := crypto.Keccak256Hash([]byte(nodeId.String()))
	currentNumber := big.NewInt(1)

	// build genesis veriferList and validatorList
	validatorQueue := make(staking.ValidatorQueue, xcom.EpochValidatorNum())

	for j := 0; j < 1000; j++ {
		var index int = j % 25

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(j),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(j) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(j) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		// Store Candidate power
		powerKey := staking.TallyPowerKey(canTmp.Shares, canTmp.StakingBlockNum, canTmp.StakingTxIndex, canTmp.ProgramVersion)
		if err := sndb.PutBaseDB(powerKey, canAddr.Bytes()); nil != err {
			t.Errorf("Failed to Store Candidate Power: PutBaseDB failed. error:%s", err.Error())
			return
		}

		// Store Candidate info
		canKey := staking.CandidateKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canTmp); nil != err {
			t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canKey, val); nil != err {
				t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		if j < int(xcom.EpochValidatorNum()) {
			v := &staking.Validator{
				NodeAddress: canAddr,
				NodeId:      canTmp.NodeId,
				BlsPubKey:   canTmp.BlsPubKey,
				StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
					fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
				ValidatorTerm: 0,
			}
			validatorQueue[j] = v
		}

	}

	/**
	*******
	build genesis epoch validators
	*******
	*/
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}

	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
		return
	}

	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	/**
	*******
	build genesis curr round validators
	*******
	*/
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
		return
	}

	xcom.PrintObject("Test round", validatorQueue[:xcom.ConsValidatorNum()])
	roundArr, err := rlp.EncodeToBytes(validatorQueue[:xcom.ConsValidatorNum()])
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Current Round validator
	if err := sndb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), roundArr); nil != err {
		t.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	// Store vrf nonces
	if err := sndb.PutBaseDB(handler.NonceStorageKey, enValue); nil != err {
		t.Errorf("Failed to Store Current Vrf nonces : PutBaseDB failed. error:%s", err.Error())
		return
	}

	// SetCurrent to snapshotDB
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50
	preNum1 := new(big.Int).Sub(currentNumber, big.NewInt(1))
	if err := sndb.SetCurrent(currentHash, *preNum1, *preNum1); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error()))
	}

	/**
	Start ListVerifierNodeId
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50

	nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
	header := &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	fmt.Println("Start ListVerifierNodeId CurrBlockHash", currentHash.Hex())

	/**
	Start  ListVerifierNodeID
	*/
	validatorIdQueue, err := StakingInstance().ListVerifierNodeID(currentHash, currentNumber.Uint64())
	if nil != err {
		t.Errorf("Failed to ListVerifierNodeID, err: %v", err)
		return
	}

	validatorIdArr, _ := json.Marshal(validatorIdQueue)
	t.Log("ListVerifierNodeID:", string(validatorIdArr))
}

func TestStakingPlugin_IsCandidate(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	handler.NewVrfHandler(genesis.Hash().Bytes())

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	nodeIdArr := make([]discover.NodeID, 0)

	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}

		if i < 20 {
			nodeIdArr = append(nodeIdArr, canTmp.NodeId)
		}
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}
	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	/**
	Start  IsCandidate
	*/
	for i, nodeId := range nodeIdArr {
		yes, err := StakingInstance().IsCandidate(blockHash2, nodeId, QueryStartNotIrr)
		if nil != err {
			t.Errorf("Failed to IsCandidate, index: %d, err: %v", i, err)
			return
		}
		if !yes {
			t.Logf("The NodeId is not a Id of Candidate, nodeId: %s", nodeId.String())
		}
	}
}

func TestStakingPlugin_IsCurrValidator(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// New VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	// build vrf proof
	// build ancestor nonces
	_, nonces := build_vrf_Nonce()
	enValue, err := rlp.EncodeToBytes(nonces)
	if nil != err {
		t.Error("Failed to rlp vrf nonces", "err", err)
		return
	}

	// new block
	privateKey, err := crypto.GenerateKey()
	if nil != err {
		t.Errorf("Failed to generate random Address private key: %v", err)
		return
	}
	nodeId := discover.PubkeyID(&privateKey.PublicKey)
	currentHash := crypto.Keccak256Hash([]byte(nodeId.String()))
	currentNumber := big.NewInt(1)

	// build genesis veriferList and validatorList
	validatorQueue := make(staking.ValidatorQueue, xcom.EpochValidatorNum())

	for j := 0; j < 1000; j++ {
		var index int = j % 25

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(j),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(j) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(j) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		// Store Candidate power
		powerKey := staking.TallyPowerKey(canTmp.Shares, canTmp.StakingBlockNum, canTmp.StakingTxIndex, canTmp.ProgramVersion)
		if err := sndb.PutBaseDB(powerKey, canAddr.Bytes()); nil != err {
			t.Errorf("Failed to Store Candidate Power: PutBaseDB failed. error:%s", err.Error())
			return
		}

		// Store Candidate info
		canKey := staking.CandidateKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canTmp); nil != err {
			t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canKey, val); nil != err {
				t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		if j < int(xcom.EpochValidatorNum()) {
			v := &staking.Validator{
				NodeAddress: canAddr,
				NodeId:      canTmp.NodeId,
				BlsPubKey:   canTmp.BlsPubKey,
				StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
					fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
				ValidatorTerm: 0,
			}
			validatorQueue[j] = v
		}

	}

	/**
	*******
	build genesis epoch validators
	*******
	*/
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}

	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
		return
	}

	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	/**
	*******
	build genesis curr round validators
	*******
	*/
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
		return
	}

	xcom.PrintObject("Test round", validatorQueue[:xcom.ConsValidatorNum()])
	roundArr, err := rlp.EncodeToBytes(validatorQueue[:xcom.ConsValidatorNum()])
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Current Round validator
	if err := sndb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), roundArr); nil != err {
		t.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	// Store vrf nonces
	if err := sndb.PutBaseDB(handler.NonceStorageKey, enValue); nil != err {
		t.Errorf("Failed to Store Current Vrf nonces : PutBaseDB failed. error:%s", err.Error())
		return
	}

	// SetCurrent to snapshotDB
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50
	preNum1 := new(big.Int).Sub(currentNumber, big.NewInt(1))
	if err := sndb.SetCurrent(currentHash, *preNum1, *preNum1); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error()))
	}

	/**
	Start IsCurrValidator
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50

	nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
	header := &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	fmt.Println("Start IsCurrValidator CurrBlockHash", currentHash.Hex())

	/**
	Start  IsCurrValidator
	*/
	for i, nodeId := range nodeIdArr {
		yes, err := StakingInstance().IsCurrValidator(currentHash, currentNumber.Uint64(), nodeId, QueryStartNotIrr)
		if nil != err {
			t.Errorf("Failed to IsCurrValidator, index: %d, err: %v", i, err)
			return
		}
		if !yes {
			t.Logf("The NodeId is not a Id of current round validator, nodeId: %s", nodeId.String())
		}
	}

}

func TestStakingPlugin_IsCurrVerifier(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// New VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	// build vrf proof
	// build ancestor nonces
	_, nonces := build_vrf_Nonce()
	enValue, err := rlp.EncodeToBytes(nonces)
	if nil != err {
		t.Error("Failed to rlp vrf nonces", "err", err)
		return
	}

	// new block
	privateKey, err := crypto.GenerateKey()
	if nil != err {
		t.Errorf("Failed to generate random Address private key: %v", err)
		return
	}
	nodeId := discover.PubkeyID(&privateKey.PublicKey)
	currentHash := crypto.Keccak256Hash([]byte(nodeId.String()))
	currentNumber := big.NewInt(1)

	// build genesis veriferList and validatorList
	validatorQueue := make(staking.ValidatorQueue, xcom.EpochValidatorNum())

	for j := 0; j < 1000; j++ {
		var index int = j % 25

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(j),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(j) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(j) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		// Store Candidate power
		powerKey := staking.TallyPowerKey(canTmp.Shares, canTmp.StakingBlockNum, canTmp.StakingTxIndex, canTmp.ProgramVersion)
		if err := sndb.PutBaseDB(powerKey, canAddr.Bytes()); nil != err {
			t.Errorf("Failed to Store Candidate Power: PutBaseDB failed. error:%s", err.Error())
			return
		}

		// Store Candidate info
		canKey := staking.CandidateKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canTmp); nil != err {
			t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canKey, val); nil != err {
				t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		if j < int(xcom.EpochValidatorNum()) {
			v := &staking.Validator{
				NodeAddress: canAddr,
				NodeId:      canTmp.NodeId,
				BlsPubKey:   canTmp.BlsPubKey,
				StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
					fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
				ValidatorTerm: 0,
			}
			validatorQueue[j] = v
		}

	}

	/**
	*******
	build genesis epoch validators
	*******
	*/
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}

	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
		return
	}

	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	/**
	*******
	build genesis curr round validators
	*******
	*/
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
		return
	}

	xcom.PrintObject("Test round", validatorQueue[:xcom.ConsValidatorNum()])
	roundArr, err := rlp.EncodeToBytes(validatorQueue[:xcom.ConsValidatorNum()])
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Current Round validator
	if err := sndb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), roundArr); nil != err {
		t.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	// Store vrf nonces
	if err := sndb.PutBaseDB(handler.NonceStorageKey, enValue); nil != err {
		t.Errorf("Failed to Store Current Vrf nonces : PutBaseDB failed. error:%s", err.Error())
		return
	}

	// SetCurrent to snapshotDB
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50
	preNum1 := new(big.Int).Sub(currentNumber, big.NewInt(1))
	if err := sndb.SetCurrent(currentHash, *preNum1, *preNum1); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error()))
	}

	/**
	Start IsCurrVerifier
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50

	nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
	header := &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	fmt.Println("Start IsCurrVerifier CurrBlockHash", currentHash.Hex())
	/**
	Start  IsCurrVerifier
	*/
	for i, nodeId := range nodeIdArr {
		yes, err := StakingInstance().IsCurrVerifier(currentHash, currentNumber.Uint64(), nodeId, QueryStartNotIrr)
		if nil != err {
			t.Errorf("Failed to IsCurrVerifier, index: %d, err: %v", i, err)
			return
		}
		if !yes {
			t.Logf("The NodeId is not a Id of Epoch validator, nodeId: %s", nodeId.String())
		}
	}
}

// for consensus
func TestStakingPlugin_GetLastNumber(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// New VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	// build vrf proof
	// build ancestor nonces
	_, nonces := build_vrf_Nonce()
	enValue, err := rlp.EncodeToBytes(nonces)
	if nil != err {
		t.Error("Failed to rlp vrf nonces", "err", err)
		return
	}

	// new block
	privateKey, err := crypto.GenerateKey()
	if nil != err {
		t.Errorf("Failed to generate random Address private key: %v", err)
		return
	}
	nodeId := discover.PubkeyID(&privateKey.PublicKey)
	currentHash := crypto.Keccak256Hash([]byte(nodeId.String()))
	currentNumber := big.NewInt(1)

	// build genesis veriferList and validatorList
	validatorQueue := make(staking.ValidatorQueue, xcom.EpochValidatorNum())

	for j := 0; j < 1000; j++ {
		var index int = j % 25

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(j),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(j) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(j) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		// Store Candidate power
		powerKey := staking.TallyPowerKey(canTmp.Shares, canTmp.StakingBlockNum, canTmp.StakingTxIndex, canTmp.ProgramVersion)
		if err := sndb.PutBaseDB(powerKey, canAddr.Bytes()); nil != err {
			t.Errorf("Failed to Store Candidate Power: PutBaseDB failed. error:%s", err.Error())
			return
		}

		// Store Candidate info
		canKey := staking.CandidateKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canTmp); nil != err {
			t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canKey, val); nil != err {
				t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		if j < int(xcom.EpochValidatorNum()) {
			v := &staking.Validator{
				NodeAddress: canAddr,
				NodeId:      canTmp.NodeId,
				BlsPubKey:   canTmp.BlsPubKey,
				StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
					fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
				ValidatorTerm: 0,
			}
			validatorQueue[j] = v
		}

	}

	/**
	*******
	build genesis epoch validators
	*******
	*/
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}

	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
		return
	}

	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	/**
	*******
	build genesis curr round validators
	*******
	*/
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
		return
	}

	xcom.PrintObject("Test round", validatorQueue[:xcom.ConsValidatorNum()])
	roundArr, err := rlp.EncodeToBytes(validatorQueue[:xcom.ConsValidatorNum()])
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Current Round validator
	if err := sndb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), roundArr); nil != err {
		t.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	// Store vrf nonces
	if err := sndb.PutBaseDB(handler.NonceStorageKey, enValue); nil != err {
		t.Errorf("Failed to Store Current Vrf nonces : PutBaseDB failed. error:%s", err.Error())
		return
	}

	// SetCurrent to snapshotDB
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50
	preNum1 := new(big.Int).Sub(currentNumber, big.NewInt(1))
	if err := sndb.SetCurrent(currentHash, *preNum1, *preNum1); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error()))
	}

	/**
	Start GetLastNumber
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50

	nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
	header := &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	fmt.Println("Start GetLastNumber CurrBlockHash", currentHash.Hex())

	/**
	Start  GetLastNumber
	*/
	endNumber := StakingInstance().GetLastNumber(currentNumber.Uint64())

	t.Log("GetLastNumber the endNumber is:", endNumber)

}

func TestStakingPlugin_GetValidator(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// New VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	// build vrf proof
	// build ancestor nonces
	_, nonces := build_vrf_Nonce()
	enValue, err := rlp.EncodeToBytes(nonces)
	if nil != err {
		t.Error("Failed to rlp vrf nonces", "err", err)
		return
	}

	// new block
	privateKey, err := crypto.GenerateKey()
	if nil != err {
		t.Errorf("Failed to generate random Address private key: %v", err)
		return
	}
	nodeId := discover.PubkeyID(&privateKey.PublicKey)
	currentHash := crypto.Keccak256Hash([]byte(nodeId.String()))
	currentNumber := big.NewInt(1)

	// build genesis veriferList and validatorList
	validatorQueue := make(staking.ValidatorQueue, xcom.EpochValidatorNum())

	for j := 0; j < 1000; j++ {
		var index int = j % 25

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(j),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(j) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(j) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		// Store Candidate power
		powerKey := staking.TallyPowerKey(canTmp.Shares, canTmp.StakingBlockNum, canTmp.StakingTxIndex, canTmp.ProgramVersion)
		if err := sndb.PutBaseDB(powerKey, canAddr.Bytes()); nil != err {
			t.Errorf("Failed to Store Candidate Power: PutBaseDB failed. error:%s", err.Error())
			return
		}

		// Store Candidate info
		canKey := staking.CandidateKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canTmp); nil != err {
			t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canKey, val); nil != err {
				t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		if j < int(xcom.EpochValidatorNum()) {
			v := &staking.Validator{
				NodeAddress: canAddr,
				NodeId:      canTmp.NodeId,
				BlsPubKey:   canTmp.BlsPubKey,
				StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
					fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
				ValidatorTerm: 0,
			}
			validatorQueue[j] = v
		}

	}

	/**
	*******
	build genesis epoch validators
	*******
	*/
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}

	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
		return
	}

	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	/**
	*******
	build genesis curr round validators
	*******
	*/
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
		return
	}

	xcom.PrintObject("Test round", validatorQueue[:xcom.ConsValidatorNum()])
	roundArr, err := rlp.EncodeToBytes(validatorQueue[:xcom.ConsValidatorNum()])
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Current Round validator
	if err := sndb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), roundArr); nil != err {
		t.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	// Store vrf nonces
	if err := sndb.PutBaseDB(handler.NonceStorageKey, enValue); nil != err {
		t.Errorf("Failed to Store Current Vrf nonces : PutBaseDB failed. error:%s", err.Error())
		return
	}

	// SetCurrent to snapshotDB
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50
	preNum1 := new(big.Int).Sub(currentNumber, big.NewInt(1))
	if err := sndb.SetCurrent(currentHash, *preNum1, *preNum1); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error()))
	}

	/**
	Start GetValidator
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50

	nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
	header := &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	fmt.Println("Start GetValidator CurrBlockHash", currentHash.Hex())

	/**
	Start  GetValidator
	*/
	valArr, err := StakingInstance().GetValidator(currentNumber.Uint64())
	if nil != err {
		t.Errorf("Failed to GetValidator, err: %v", err)
	}

	arrJson, _ := json.Marshal(valArr)
	t.Log("GetValidator the validators is:", string(arrJson))

}

func TestStakingPlugin_IsCandidateNode(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// New VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	// build vrf proof
	// build ancestor nonces
	_, nonces := build_vrf_Nonce()
	enValue, err := rlp.EncodeToBytes(nonces)
	if nil != err {
		t.Error("Failed to rlp vrf nonces", "err", err)
		return
	}

	// new block
	privateKey, err := crypto.GenerateKey()
	if nil != err {
		t.Errorf("Failed to generate random Address private key: %v", err)
		return
	}
	nodeId := discover.PubkeyID(&privateKey.PublicKey)
	currentHash := crypto.Keccak256Hash([]byte(nodeId.String()))
	currentNumber := big.NewInt(1)

	// build genesis veriferList and validatorList
	validatorQueue := make(staking.ValidatorQueue, xcom.EpochValidatorNum())

	for j := 0; j < 1000; j++ {
		var index int = j % 25

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(j),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(j) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(j) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		// Store Candidate power
		powerKey := staking.TallyPowerKey(canTmp.Shares, canTmp.StakingBlockNum, canTmp.StakingTxIndex, canTmp.ProgramVersion)
		if err := sndb.PutBaseDB(powerKey, canAddr.Bytes()); nil != err {
			t.Errorf("Failed to Store Candidate Power: PutBaseDB failed. error:%s", err.Error())
			return
		}

		// Store Candidate info
		canKey := staking.CandidateKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canTmp); nil != err {
			t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canKey, val); nil != err {
				t.Errorf("Failed to Store Candidate info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		if j < int(xcom.EpochValidatorNum()) {
			v := &staking.Validator{
				NodeAddress: canAddr,
				NodeId:      canTmp.NodeId,
				BlsPubKey:   canTmp.BlsPubKey,
				StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
					fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
				ValidatorTerm: 0,
			}
			validatorQueue[j] = v
		}

	}

	/**
	*******
	build genesis epoch validators
	*******
	*/
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}

	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
		return
	}

	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	/**
	*******
	build genesis curr round validators
	*******
	*/
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
		return
	}

	xcom.PrintObject("Test round", validatorQueue[:xcom.ConsValidatorNum()])
	roundArr, err := rlp.EncodeToBytes(validatorQueue[:xcom.ConsValidatorNum()])
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Current Round validator
	if err := sndb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), roundArr); nil != err {
		t.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	// Store vrf nonces
	if err := sndb.PutBaseDB(handler.NonceStorageKey, enValue); nil != err {
		t.Errorf("Failed to Store Current Vrf nonces : PutBaseDB failed. error:%s", err.Error())
		return
	}

	// SetCurrent to snapshotDB
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50
	preNum1 := new(big.Int).Sub(currentNumber, big.NewInt(1))
	if err := sndb.SetCurrent(currentHash, *preNum1, *preNum1); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error()))
	}

	/**
	Start IsCandidateNode
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50

	nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
	header := &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	fmt.Println("Start IsCandidateNode CurrBlockHash", currentHash.Hex())
	/**
	Start  IsCandidateNode
	*/
	yes := StakingInstance().IsCandidateNode(nodeIdArr[0])

	t.Log("IsCandidateNode the flag is:", yes)

}

func TestStakingPlugin_ProbabilityElection(t *testing.T) {

	newChainState()

	curve := elliptic.P256()
	vqList := make(staking.ValidatorQueue, 0)
	preNonces := make([][]byte, 0)
	currentNonce := crypto.Keccak256([]byte(string("nonce")))
	for i := 0; i < int(xcom.EpochValidatorNum()); i++ {
		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		privKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
		nodeId := discover.PubkeyID(&privKey.PublicKey)
		addr := crypto.PubkeyToAddress(privKey.PublicKey)
		mrand.Seed(time.Now().UnixNano())
		stakingWeight := [staking.SWeightItem]string{}
		stakingWeight[0] = strconv.Itoa(mrand.Intn(5) + 1)
		v1 := new(big.Int).SetInt64(time.Now().UnixNano())
		v1.Mul(v1, new(big.Int).SetInt64(1e18))
		v1.Add(v1, new(big.Int).SetInt64(int64(mrand.Intn(1000))))
		stakingWeight[1] = v1.Text(10)
		stakingWeight[2] = strconv.Itoa(mrand.Intn(230))
		stakingWeight[3] = strconv.Itoa(mrand.Intn(1000))
		v := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        nodeId,
			BlsPubKey:     *blsKey.GetPublicKey(),
			StakingWeight: stakingWeight,
			ValidatorTerm: 1,
		}
		vqList = append(vqList, v)
		preNonces = append(preNonces, crypto.Keccak256(common.Int64ToBytes(time.Now().UnixNano() + int64(i)))[:])
		time.Sleep(time.Microsecond * 10)
	}
	for index, v := range vqList {
		t.Log("Generate Validator", "addr", hex.EncodeToString(v.NodeAddress.Bytes()), "stakingWeight", v.StakingWeight, "nonce", hex.EncodeToString(preNonces[index]))
	}
	result, err := probabilityElection(vqList, int(xcom.ShiftValidatorNum()), currentNonce, preNonces)
	if nil != err {
		t.Fatal("Failed to ProbabilityElection, err:", err)
		return
	}
	t.Log("election success", result)
	for _, v := range result {
		t.Log("Validator", "addr", hex.EncodeToString(v.NodeAddress.Bytes()), "stakingWeight", v.StakingWeight)
	}
}

/**
Expand test cases
*/

func Test_IteratorCandidate(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		//t.Log("Create Staking num:", index)

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			BlsPubKey:       *blsKey.GetPublicKey(),
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	// commit
	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	stakingDB := staking.NewStakingDB()

	iter := stakingDB.IteratorCandidatePowerByBlockHash(blockHash2, 0)
	if err := iter.Error(); nil != err {
		t.Error("Get iter err", err)
		return
	}
	defer iter.Release()

	queue := make(staking.CandidateQueue, 0)

	for iter.Valid(); iter.Next(); {
		addrSuffix := iter.Value()
		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash2, addrSuffix)
		if nil != err {
			t.Errorf("Failed to Iterator Candidate info, err: %v", err)
			return
		}

		val := fmt.Sprint(can.ProgramVersion) + "_" + can.Shares.String() + "_" + fmt.Sprint(can.StakingBlockNum) + "_" + fmt.Sprint(can.StakingTxIndex)
		t.Log("Val:", val)

		queue = append(queue, can)
	}

	arrJson, _ := json.Marshal(queue)
	t.Log("CandidateList:", string(arrJson))
	t.Log("Candidate queue length:", len(queue))
}
