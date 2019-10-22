package plugin

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"

	"github.com/PlatONnetwork/PlatON-Go/common/mock"

	"github.com/PlatONnetwork/PlatON-Go/common/vm"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"github.com/stretchr/testify/assert"
)

func initInfo(t *testing.T) (*SlashingPlugin, xcom.StateDB) {
	si := SlashInstance()
	StakingInstance()
	RestrictingInstance()
	chain := mock.NewChain(nil)
	return si, chain.StateDB
}

func buildStakingData(blockHash common.Hash, pri *ecdsa.PrivateKey, blsKey bls.SecretKey, t *testing.T, stateDb xcom.StateDB) {
	stakingDB := staking.NewStakingDB()

	sender := common.HexToAddress("0xeef233120ce31b3fac20dac379db243021a5234")

	buildDbRestrictingPlan(sender, t, stateDb)

	if nil == pri {
		sk, err := crypto.GenerateKey()
		if nil != err {
			panic(err)
		}
		pri = sk
	}

	nodeIdA := discover.PubkeyID(&pri.PublicKey)
	addrA, _ := xutil.NodeId2Addr(nodeIdA)

	nodeIdB := nodeIdArr[1]
	addrB, _ := xutil.NodeId2Addr(nodeIdB)

	nodeIdC := nodeIdArr[2]
	addrC, _ := xutil.NodeId2Addr(nodeIdC)

	c1 := &staking.Candidate{
		NodeId:             nodeIdA,
		BlsPubKey:          *blsKey.GetPublicKey(),
		StakingAddress:     sender,
		BenefitAddress:     addrArr[1],
		StakingTxIndex:     uint32(2),
		ProgramVersion:     uint32(1),
		Status:             staking.Valided,
		StakingEpoch:       uint32(1),
		StakingBlockNum:    uint64(1),
		Shares:             common.Big256,
		Released:           new(big.Int).Mul(common.Big256, new(big.Int).SetUint64(1000)),
		ReleasedHes:        common.Big32,
		RestrictingPlan:    common.Big0,
		RestrictingPlanHes: common.Big0,
		Description: staking.Description{
			ExternalId: "xxccccdddddddd",
			NodeName:   "I Am " + fmt.Sprint(1),
			Website:    "www.baidu.com",
			Details:    "this is  baidu ~~",
		},
	}

	var blsKey2 bls.SecretKey
	blsKey2.SetByCSPRNG()
	c2 := &staking.Candidate{
		NodeId:             nodeIdB,
		BlsPubKey:          *blsKey2.GetPublicKey(),
		StakingAddress:     sender,
		BenefitAddress:     addrArr[2],
		StakingTxIndex:     uint32(3),
		ProgramVersion:     uint32(1),
		Status:             staking.Valided,
		StakingEpoch:       uint32(1),
		StakingBlockNum:    uint64(1),
		Shares:             common.Big256,
		Released:           common.Big2,
		ReleasedHes:        common.Big32,
		RestrictingPlan:    common.Big0,
		RestrictingPlanHes: common.Big0,
		Description: staking.Description{
			ExternalId: "SFSFSFSFSFSFSSFS",
			NodeName:   "I Am " + fmt.Sprint(2),
			Website:    "www.JD.com",
			Details:    "this is  JD ~~",
		},
	}

	var blsKey3 bls.SecretKey
	blsKey3.SetByCSPRNG()
	c3 := &staking.Candidate{
		NodeId:             nodeIdC,
		BlsPubKey:          *blsKey3.GetPublicKey(),
		StakingAddress:     sender,
		BenefitAddress:     addrArr[3],
		StakingTxIndex:     uint32(4),
		ProgramVersion:     uint32(1),
		Status:             staking.Valided,
		StakingEpoch:       uint32(1),
		StakingBlockNum:    uint64(1),
		Shares:             common.Big256,
		Released:           common.Big2,
		ReleasedHes:        common.Big32,
		RestrictingPlan:    common.Big0,
		RestrictingPlanHes: common.Big0,
		Description: staking.Description{
			ExternalId: "FWAGGDGDGG",
			NodeName:   "I Am " + fmt.Sprint(3),
			Website:    "www.alibaba.com",
			Details:    "this is  alibaba ~~",
		},
	}

	stakingDB.SetCanPowerStore(blockHash, addrA, c1)
	stakingDB.SetCanPowerStore(blockHash, addrB, c2)
	stakingDB.SetCanPowerStore(blockHash, addrC, c3)

	stakingDB.SetCandidateStore(blockHash, addrA, c1)
	stakingDB.SetCandidateStore(blockHash, addrB, c2)
	stakingDB.SetCandidateStore(blockHash, addrC, c3)

	log.Info("addr_A", hex.EncodeToString(addrA.Bytes()), "addr_B", hex.EncodeToString(addrB.Bytes()), "addr_C", hex.EncodeToString(addrC.Bytes()))

	queue := make(staking.ValidatorQueue, 0)

	v1 := &staking.Validator{
		NodeAddress:   addrA,
		NodeId:        c1.NodeId,
		BlsPubKey:     c1.BlsPubKey,
		StakingWeight: [staking.SWeightItem]string{"1", common.Big256.String(), fmt.Sprint(c1.StakingBlockNum), fmt.Sprint(c1.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	v2 := &staking.Validator{
		NodeAddress:   addrB,
		NodeId:        c2.NodeId,
		BlsPubKey:     c2.BlsPubKey,
		StakingWeight: [staking.SWeightItem]string{"1", common.Big256.String(), fmt.Sprint(c2.StakingBlockNum), fmt.Sprint(c2.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	v3 := &staking.Validator{
		NodeAddress:   addrC,
		NodeId:        c3.NodeId,
		BlsPubKey:     c3.BlsPubKey,
		StakingWeight: [staking.SWeightItem]string{"1", common.Big256.String(), fmt.Sprint(c3.StakingBlockNum), fmt.Sprint(c3.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	queue = append(queue, v1)
	queue = append(queue, v2)
	queue = append(queue, v3)

	epochArr := &staking.ValidatorArray{
		Start: 1,
		End:   uint64(xutil.CalcBlocksEachEpoch()),
		Arr:   queue,
	}

	preArr := &staking.ValidatorArray{
		Start: 1,
		End:   xutil.ConsensusSize(),
		Arr:   queue,
	}

	curArr := &staking.ValidatorArray{
		Start: xutil.ConsensusSize() + 1,
		End:   xutil.ConsensusSize() * 2,
		Arr:   queue,
	}

	setVerifierList(blockHash, epochArr)
	setRoundValList(blockHash, preArr)
	setRoundValList(blockHash, curArr)
	stk.storeRoundValidatorAddrs(blockHash, 1, queue)
	balance, ok := new(big.Int).SetString("9999999999999999999999999999999999999999999999999", 10)
	if !ok {
		panic("set balance fail")
	}
	stateDb.AddBalance(vm.StakingContractAddr, balance)
}

func TestSlashingPlugin_BeginBlock(t *testing.T) {
	_, _, _ = newChainState()
	si, stateDB := initInfo(t)
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	startNumber := xutil.ConsensusSize()
	startNumber += xutil.ConsensusSize() - xcom.ElectionDistance() - 2
	pri, phash := buildBlock(t, int(startNumber), stateDB)
	startNumber++
	blockNumber := new(big.Int).SetInt64(int64(startNumber))
	if err := snapshotdb.Instance().NewBlock(blockNumber, phash, common.ZeroHash); err != nil {
		t.Fatal(err)
	}
	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()
	buildStakingData(common.ZeroHash, pri, blsKey, t, stateDB)

	phash = common.HexToHash("0x0a0409021f020b080a16070609071c141f19011d090b091303121e1802130406")
	if err := snapshotdb.Instance().Flush(phash, blockNumber); err != nil {
		t.Fatal(err)
	}
	if err := snapshotdb.Instance().Commit(phash); err != nil {
		t.Fatal(err)
	}
	startNumber++
	header := &types.Header{
		Number: new(big.Int).SetUint64(uint64(startNumber)),
		Extra:  make([]byte, 97),
	}
	sk, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	sign, err := crypto.Sign(header.SealHash().Bytes(), sk)
	if nil != err {
		t.Fatal(err)
	}
	copy(header.Extra[len(header.Extra)-common.ExtraSeal:], sign[:])
	if err := snapshotdb.Instance().NewBlock(header.Number, phash, common.ZeroHash); nil != err {
		t.Fatal(err)
	}
	if err := si.BeginBlock(common.ZeroHash, header, stateDB); nil != err {
		t.Fatal(err)
	}
}

func buildBlock(t *testing.T, maxNumber int, stateDb xcom.StateDB) (*ecdsa.PrivateKey, common.Hash) {
	pri, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	pri2, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	db := snapshotdb.Instance()
	sk := pri
	_, genesis, _ := newChainState()
	parentHash := genesis.Hash()
	for i := 0; i < maxNumber; i++ {
		blockNum := big.NewInt(int64(i + 1))
		if i == int(1) {
			sk = pri2
		}
		header := &types.Header{
			Number: blockNum,
			Extra:  make([]byte, 97),
		}
		sign, err := crypto.Sign(header.SealHash().Bytes(), sk)
		if nil != err {
			t.Fatal(err)
		}
		copy(header.Extra[len(header.Extra)-common.ExtraSeal:], sign[:])
		if err := db.NewBlock(blockNum, parentHash, common.ZeroHash); err != nil {
			t.Fatal(err)
		}
		if err := SlashInstance().BeginBlock(common.ZeroHash, header, stateDb); nil != err {
			t.Fatal(err)
		}
		if err := db.Flush(header.Hash(), blockNum); err != nil {
			t.Fatal(err)
		}
		if err := db.Commit(header.Hash()); err != nil {
			t.Fatal(err)
		}
		parentHash = header.Hash()
	}
	return pri, parentHash
}

func TestSlashingPlugin_Slash(t *testing.T) {
	_, genesis, _ := newChainState()
	si, stateDB := initInfo(t)
	blockNumber := new(big.Int).SetUint64(1)
	commitHash := common.HexToHash("0x0a0409021f020b080a16070609071c141f19011d090b091303121e1802130406")
	if err := snapshotdb.Instance().NewBlock(blockNumber, genesis.Hash(), common.ZeroHash); nil != err {
		t.Fatal(err)
	}
	var nodeBlsKey bls.SecretKey
	nodeBlsSkByte, err := hex.DecodeString("72fc21a19510d93d726746602344f96bf181efdd8d6d95be1a2a2de59bd59501")
	if nil != err {
		t.Fatalf("ReportDuplicateSign DecodeString byte data fail: %v", err)
	}
	nodeBlsKey.SetLittleEndian(nodeBlsSkByte)
	buildStakingData(common.ZeroHash, crypto.HexMustToECDSA("f976838ffad88eb7cb45217e0d74c71a1adcb03d550aea9c32df4cfd41e1e0ca"), nodeBlsKey, t, stateDB)
	if err := snapshotdb.Instance().Flush(commitHash, blockNumber); nil != err {
		t.Fatal(err)
	}
	if err := snapshotdb.Instance().Commit(commitHash); nil != err {
		t.Fatal(err)
	}
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	si.SetDecodeEvidenceFun(evidence.NewEvidence)
	GovPluginInstance()
	normalData := `{
           "prepareA": {
            "epoch": 1,
            "viewNumber": 1,
            "blockHash": "0x3ba21ce163d1ca763d0772b70e652259608b3f40c0641238813f1cac708b75c5",
            "blockNumber": 1,
            "blockIndex": 1,
            "blockData": "0x94ac820f54471ae9a32342f8a86e516944ec333a717241428ed997c4d3c1c8e3",
            "validateNode": {
             "index": 0,
             "address": "0x85396cdef1d2800c621361437c2439c59c934038",
             "nodeId": "c0b49363fa1c2a0d3c55cafec4955cb261a537afd4fe45ff21c7b84cba660d5157865d984c2d2a61b4df1d3d028634136d04030ed6a388b429eaa6e2bdefaed1",
             "blsPubKey": "f9b5e5b333418f5f6cb23ad092d2321c49a6fc17dfa2e5899a0fa0a6ab96bc44482552c9149f5909ec7772a902094401912576fdd78497bf57399c711566284ae2f5db3f8e611ac21dbc53cf7c1ff881ab760c0f1e5954b9cd2602b98007ef05"
            },
            "signature": "0xcd01a6ed0ee36d346fc0cf2eaa0151b775e22ffd97a8c9c5fada22f43deee2940776a3da82e2ba9ea4499037c4a3321200000000000000000000000000000000"
           },
           "prepareB": {
            "epoch": 1,
            "viewNumber": 1,
            "blockHash": "0xb337d296f74ac7063f180bbd86834377fd77459e414b87c547787a384010490e",
            "blockNumber": 1,
            "blockIndex": 1,
            "blockData": "0xc789252723b04c60fc4566abefa23aa4e9ef18d9b4ebd1b083a564700cbb8891",
            "validateNode": {
             "index": 0,
             "address": "0x85396cdef1d2800c621361437c2439c59c934038",
             "nodeId": "c0b49363fa1c2a0d3c55cafec4955cb261a537afd4fe45ff21c7b84cba660d5157865d984c2d2a61b4df1d3d028634136d04030ed6a388b429eaa6e2bdefaed1",
             "blsPubKey": "f9b5e5b333418f5f6cb23ad092d2321c49a6fc17dfa2e5899a0fa0a6ab96bc44482552c9149f5909ec7772a902094401912576fdd78497bf57399c711566284ae2f5db3f8e611ac21dbc53cf7c1ff881ab760c0f1e5954b9cd2602b98007ef05"
            },
            "signature": "0x06ebc53e4227a89c6a7f2adf978436cb829fbc47d4e6189569fb240a2c8c1f0a2b3b63fdbf905aaa3f1ffe7b0b4d7e8e00000000000000000000000000000000"
           }
          }`

	normalData2 := `{
           "prepareA": {
            "epoch": 1,
            "viewNumber": 1,
            "blockHash": "0xbb6d4b83af8667929a9cb4918bbf790a97bb136775353765388d0add3437cde6",
            "blockNumber": 1,
            "blockIndex": 1,
            "blockData": "0x45b20c5ba595be254943aa57cc80562e84f1fb3bafbf4a414e30570c93a39579",
            "validateNode": {
             "index": 0,
             "address": "0x195667cdefcad94c521bdff0bf85079761e0f8f3",
             "nodeId": "51c0559c065400151377d71acd7a17282a7c8abcfefdb11992dcecafde15e100b8e31e1a5e74834a04792d016f166c80b9923423fe280570e8131debf591d483",
             "blsPubKey": "752fe419bbdc2d2222009e450f2932657bbc2370028d396ba556a49439fe1cc11903354dcb6dac552a124e0b3db0d90edcd334d7aabda0c3f1ade12ca22372f876212ac456d549dbbd04d2c8c8fb3e33760215e114b4d60313c142f7b8bbfd87"
            },
            "signature": "0x36015fee15253487e8125b86505377d8540b1a95d1a6b13f714baa55b12bd06ec7d5755a98230cdc88858470afa8cb0000000000000000000000000000000000"
           },
           "prepareB": {
            "epoch": 1,
            "viewNumber": 1,
            "blockHash": "0xf46c45f7ebb4a999dd030b9f799198b785654293dbe41aa7e909223af0e8c4ba",
            "blockNumber": 1,
            "blockIndex": 1,
            "blockData": "0xd630e96d127f55319392f20d4fd917e3e7cba19ad366c031b9dff05e056d9420",
            "validateNode": {
             "index": 0,
             "address": "0x195667cdefcad94c521bdff0bf85079761e0f8f3",
             "nodeId": "51c0559c065400151377d71acd7a17282a7c8abcfefdb11992dcecafde15e100b8e31e1a5e74834a04792d016f166c80b9923423fe280570e8131debf591d483",
             "blsPubKey": "752fe419bbdc2d2222009e450f2932657bbc2370028d396ba556a49439fe1cc11903354dcb6dac552a124e0b3db0d90edcd334d7aabda0c3f1ade12ca22372f876212ac456d549dbbd04d2c8c8fb3e33760215e114b4d60313c142f7b8bbfd87"
            },
            "signature": "0x783892b9b766f9f4c2a1d45b1fd53ca9ea56a82e38a998939edc17bc7fd756267d3c145c03bc6c1412302cf590645d8200000000000000000000000000000000"
           }
          }`

	abnormalData := `{
           "prepareA": {
            "epoch": 1,
            "viewNumber": 1,
            "blockHash": "0x2973fa91cd5cc27cc598bf2bb72d074a2fcfd17f820135f5343401ef59909d31",
            "blockNumber": 1,
            "blockIndex": 1,
            "blockData": "0x45b20c5ba595be254943aa57cc80562e84f1fb3bafbf4a414e30570c93a39579",
            "validateNode": {
             "index": 0,
             "address": "0x195667cdefcad94c521bdff0bf85079761e0f8f3",
             "nodeId": "51c0559c065400151377d71acd7a17282a7c8abcfefdb11992dcecafde15e100b8e31e1a5e74834a04792d016f166c80b9923423fe280570e8131debf591d483",
             "blsPubKey": "752fe419bbdc2d2222009e450f2932657bbc2370028d396ba556a49439fe1cc11903354dcb6dac552a124e0b3db0d90edcd334d7aabda0c3f1ade12ca22372f876212ac456d549dbbd04d2c8c8fb3e33760215e114b4d60313c142f7b8bbfd87"
            },
            "signature": "0x36015fee15253487e8125b86505377d8540b1a95d1a6b13f714baa55b12bd06ec7d5755a98230cdc88858470afa8cb0000000000000000000000000000000000"
           },
           "prepareB": {
            "epoch": 1,
            "viewNumber": 1,
            "blockHash": "0xf46c45f7ebb4a999dd030b9f799198b785654293dbe41aa7e909223af0e8c4ba",
            "blockNumber": 1,
            "blockIndex": 1,
            "blockData": "0xd630e96d127f55319392f20d4fd917e3e7cba19ad366c031b9dff05e056d9420",
            "validateNode": {
             "index": 0,
             "address": "0x195667cdefcad94c521bdff0bf85079761e0f8f3",
             "nodeId": "51c0559c065400151377d71acd7a17282a7c8abcfefdb11992dcecafde15e100b8e31e1a5e74834a04792d016f166c80b9923423fe280570e8131debf591d483",
             "blsPubKey": "752fe419bbdc2d2222009e450f2932657bbc2370028d396ba556a49439fe1cc11903354dcb6dac552a124e0b3db0d90edcd334d7aabda0c3f1ade12ca22372f876212ac456d549dbbd04d2c8c8fb3e33760215e114b4d60313c142f7b8bbfd87"
            },
            "signature": "0x783892b9b766f9f4c2a1d45b1fd53ca9ea56a82e38a998939edc17bc7fd756267d3c145c03bc6c1412302cf590645d8200000000000000000000000000000000"
           }
          }`
	blockNumber = new(big.Int).Add(blockNumber, common.Big1)
	stakingAddr := common.HexToAddress("0x195667cdefcad94c521bdff0bf85079761e0f8f3")
	stakingNodeId, err := discover.HexID("51c0559c065400151377d71acd7a17282a7c8abcfefdb11992dcecafde15e100b8e31e1a5e74834a04792d016f166c80b9923423fe280570e8131debf591d483")
	if nil != err {
		t.Fatal(err)
	}
	var stakingBlsKey bls.SecretKey
	blsSkByte, err := hex.DecodeString("b36d4c3c3e8ee7fba3fbedcda4e0493e699cd95b68594093a8498c618680480a")
	if nil != err {
		t.Fatalf("ReportDuplicateSign DecodeString byte data fail: %v", err)
	}
	stakingBlsKey.SetLittleEndian(blsSkByte)
	can := &staking.Candidate{
		NodeId:          stakingNodeId,
		BlsPubKey:       *stakingBlsKey.GetPublicKey(),
		StakingAddress:  stakingAddr,
		BenefitAddress:  stakingAddr,
		StakingBlockNum: blockNumber.Uint64(),
		StakingTxIndex:  1,
		ProgramVersion:  xutil.CalcVersion(initProgramVersion),
		Shares:          new(big.Int).SetUint64(1000),

		Released:           common.Big256,
		ReleasedHes:        common.Big0,
		RestrictingPlan:    common.Big0,
		RestrictingPlanHes: common.Big0,
	}
	stateDB.CreateAccount(stakingAddr)
	stateDB.AddBalance(stakingAddr, new(big.Int).SetUint64(1000000000000000000))
	if err := snapshotdb.Instance().NewBlock(blockNumber, commitHash, common.ZeroHash); nil != err {
		t.Fatal(err)
	}
	if err := StakingInstance().CreateCandidate(stateDB, common.ZeroHash, blockNumber, can.Shares, 0, stakingAddr, can); nil != err {
		t.Fatal(err)
	}
	normalEvidence, err := si.DecodeEvidence(1, normalData)
	if nil != err {
		t.Fatal(err)
	}
	// Report yourself, expect failure
	err = si.Slash(normalEvidence, common.ZeroHash, blockNumber.Uint64(), stateDB, sender)
	assert.NotNil(t, err)
	if err := si.Slash(normalEvidence, common.ZeroHash, blockNumber.Uint64(), stateDB, anotherSender); nil != err {
		t.Fatal(err)
	}
	if value, err := si.CheckDuplicateSign(common.HexToAddress("0x85396cdef1d2800c621361437c2439c59c934038"), common.Big1.Uint64(), 1, stateDB); nil != err || len(value) == 0 {
		t.Fatal(err)
	}
	abnormalEvidence, err := si.DecodeEvidence(1, abnormalData)
	if nil != err {
		t.Fatal(err)
	}
	err = si.Slash(abnormalEvidence, common.ZeroHash, blockNumber.Uint64(), stateDB, anotherSender)
	assert.NotNil(t, err)

	// Repeat report, expected failure
	err = si.Slash(normalEvidence, common.ZeroHash, blockNumber.Uint64(), stateDB, anotherSender)
	assert.NotNil(t, err)

	// Report outdated evidence, expected failure
	err = si.Slash(abnormalEvidence, common.ZeroHash, new(big.Int).SetUint64(xutil.CalcBlocksEachEpoch()*uint64(xcom.EvidenceValidEpoch())*2).Uint64(), stateDB, anotherSender)
	assert.NotNil(t, err)

	normalEvidence2, err := si.DecodeEvidence(1, normalData2)
	if nil != err {
		t.Fatal(err)
	}

	// Report candidate nodes, expected failure
	err = si.Slash(normalEvidence2, common.ZeroHash, blockNumber.Uint64(), stateDB, anotherSender)
	assert.NotNil(t, err)
}

func TestSlashingPlugin_CheckMutiSign(t *testing.T) {
	si, stateDB := initInfo(t)
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	addr := common.HexToAddress("0x120b77ab712589ebd42d69003893ef962cc52832")
	if _, err := si.CheckDuplicateSign(addr, 1, 1, stateDB); nil != err {
		t.Fatal(err)
	}
}
