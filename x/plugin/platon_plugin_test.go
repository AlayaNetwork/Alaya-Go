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
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/params"

	"github.com/AlayaNetwork/Alaya-Go/x/gov"

	"github.com/AlayaNetwork/Alaya-Go/crypto/bls"
	"github.com/AlayaNetwork/Alaya-Go/log"

	//	"github.com/AlayaNetwork/Alaya-Go/core/state"

	"github.com/AlayaNetwork/Alaya-Go/common/mock"

	"github.com/AlayaNetwork/Alaya-Go/common"
	cvm "github.com/AlayaNetwork/Alaya-Go/common/vm"
	"github.com/AlayaNetwork/Alaya-Go/core/snapshotdb"

	//	"github.com/AlayaNetwork/Alaya-Go/core/state"
	"github.com/AlayaNetwork/Alaya-Go/core/types"
	//	"github.com/AlayaNetwork/Alaya-Go/core/vm"
	"github.com/AlayaNetwork/Alaya-Go/crypto"
	"github.com/AlayaNetwork/Alaya-Go/rlp"
	"github.com/AlayaNetwork/Alaya-Go/x/restricting"
	"github.com/AlayaNetwork/Alaya-Go/x/staking"
	"github.com/AlayaNetwork/Alaya-Go/x/xcom"
	"github.com/AlayaNetwork/Alaya-Go/x/xutil"
)

func init() {
	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	bls.Init(bls.BLS12_381)
}

var currentTestGenesisVersion = params.FORKVERSION_0_17_0

var (
	nodeIdArr = []enode.IDv0{
		enode.MustHexIDv0("5a942bc607d970259e203f5110887d6105cc787f7433c16ce28390fb39f1e67897b0fb445710cc836b89ed7f951c57a1f26a0940ca308d630448b5bd391a8aa6"),
		enode.MustHexIDv0("c453d29394e613e85999129b8fb93146d584d5a0be16f7d13fd1f44de2d01bae104878eba8e8f6b8d2c162b5a35d5939d38851f856e56186471dd7de57e9bfa9"),
		enode.MustHexIDv0("2c1733caf5c23086612a309f5ee8e76ca45455351f7cf069bcde59c07175607325cf2bf2485daa0fbf1f9cdee6eea246e5e00b9a0d0bfed0f02b37f3b0c70490"),
		enode.MustHexIDv0("e7edfb4f9c3e1fe0288ddcf0894535214fa03acea941c7360ccf90e86460aefa118ba9f2573921349c392cd1b5d4db90b4795ab353df3c915b2e8481d241ec57"),

		enode.MustHexIDv0("3a06953a2d5d45b29167bef58208f1287225bdd2591260af29ae1300aeed362e9b548369dfc1659abbef403c9b3b07a8a194040e966acd6e5b6d55aa2df7c1d8"),
		enode.MustHexIDv0("fd06314e027c3812bd0d1cf0ce1b5742d21d1ae5a397da6e7eed463ad1172c268092c2b3de52a204aabb3a6048be48f4880ce54ff3116a3843d4087d219db054"),
		enode.MustHexIDv0("811eb49e3127389065f41aac395d15e1e9968555f43913447ebb358705a63b2de37ab890f06854034a2dd171daf873adf8647498200a54cf376fcbe07d12ecd8"),
		enode.MustHexIDv0("b3d3667793ea2c2a77848b89bed514cd6fd7d685af4ee9d2482b6c58f8b3dd371cf8a41aa638e45ce420df323dfff6ed041213c343066348b4e1b39bd1396f48"),

		enode.MustHexIDv0("0x248af08a775ff63a47a5970e4928bcccd1a8cef984fd4142ea7f89cd13015bdab9ca4a8c5e1070dc00fa81a047542f53ca596f553c4acfb7abe75a8fb5019057"),
		enode.MustHexIDv0("0xfd790ff5dc48baccb9418ce5cfac6a10c3646f20a3fe32d9502c4edce3a77fa90bfee0361d8a72093b7994f8cbc28ee537bdda2b634c5966b1a9253d9d270145"),
		enode.MustHexIDv0("0x56d243db84a521cb204f582ee84bca7f4af29437dd447a6e36d17f4853888e05343844bd64294b99b835ca7f72ef5b1325ef1c89b0c5c2744154cdadf7c4e9fa"),
		enode.MustHexIDv0("0x8796a6fcefd9037d8433e3a959ff8f3c4552a482ce727b00a90bfd1ec365ce2faa33e19aa6a172b5c186b51f5a875b5acd35063171f0d9501a9c8f1c98513825"),
		enode.MustHexIDv0("0x547b876036165d66274ce31692165c8acb6f140a65cab0e0e12f1f09d1c7d8d53decf997830919e4f5cacb2df1adfe914c53d22e3ab284730b78f5c63a273b8c"),
		enode.MustHexIDv0("0x9fdbeb873bea2557752eabd2c96419b8a700b680716081472601ddf7498f0db9b8a40797b677f2fac541031f742c2bbd110ff264ae3400bf177c456a76a93d42"),
		enode.MustHexIDv0("0xc553783799bfef7c34a84b2737f2c77f8f2c5cfedc3fd7af2d944da6ece90aa94cf621e6de5c4495881fbfc9beec655ffb10e39cb4ca9be7768d284409040f32"),
		enode.MustHexIDv0("0x75ad2ee8ca77619c3ba0ddcec5dab1375fe4fa90bab9e751caef3996ce082dfed32fe4c137401ee05e501c079b2e4400397b09de14b08b09c9e7f9698e9e4f0a"),
		enode.MustHexIDv0("0xdb18af9be2af9dff2347c3d06db4b1bada0598d099a210275251b68fa7b5a863d47fcdd382cc4b3ea01e5b55e9dd0bdbce654133b7f58928ce74629d5e68b974"),
		enode.MustHexIDv0("0x472d19e5e9888368c02f24ebbbe0f2132096e7183d213ab65d96b8c03205f88398924af8876f3c615e08aa0f9a26c38911fda26d51c602c8d4f8f3cb866808d7"),
		enode.MustHexIDv0("4f1f036e5e18cc812347d5073cbec2a8da7930de323063c39b0d4413a396e088bfa90e8c28174313d8d82e9a14bc0884b13a48fc28e619e44c48a49b4fd9f107"),
		enode.MustHexIDv0("f18c596232d637409c6295abb1e720db99ffc12363a1eb8123d6f54af80423a5edd06f91115115a1dca1377e97b9031e2ddb864d34d9b3491d6fa07e8d9b951b"),
		enode.MustHexIDv0("7a8f7a28ac1c4eaf98b2be890f372e5abc58ebe6d3aab47aedcb0076e34eb42882e926676ebab327a4ef4e2ea5c4296e9c7bc0991360cb44f52672631012db1b"),
		enode.MustHexIDv0("9eeb448babf9e93449e831b91f98d9cbc0c2324fe8c43baac69d090717454f3f930713084713fe3a9f01e4ca59b80a0f2b41dbd6d531f414650bab0363e3691a"),
		enode.MustHexIDv0("cc1d7314c15e30dc5587f675eb5f803b1a2d88bfe76cec591cec1ff678bc6abce98f40054325bdcb44fb83174f27d38a54fbce4846af8f027b333868bc5144a4"),
		enode.MustHexIDv0("e4d99694be2fc8a53d8c2446f947aec1c7de3ee26f7cd43f4f6f77371f56f11156218dec32b51ddce470e97127624d330bb7a3237ba5f0d87d2d3166faf1035e"),
		enode.MustHexIDv0("9c61f59f70296b6d494e7230888e58f19b13c5c6c85562e57e1fe02d0ff872b4957238c73559d017c8770b999891056aa6329dbf628bc19028d8f4d35ec35823"),
	}

	addrArr = []common.Address{

		common.MustBech32ToAddress("atx1avltgjnqmy6alefayfry3cd9rpguduawy506sh"),
		common.MustBech32ToAddress("atx1rkdnqnnsl5shqm7e00897dpey33h3pcnh2yalf"),
		common.MustBech32ToAddress("atx184w6gavcetzpyytxja005ynq8rmjeagl6t8sat"),
		common.MustBech32ToAddress("atx1erk3dpm9u9cfutnsqskfsrgkvc533r4pevzmq4"),

		common.MustBech32ToAddress("atx1a4g8npqllsa5ffkw8y2p3lxxvql2955yml3mmh"),
		common.MustBech32ToAddress("atx1s5f554lz0agjvdlxkwkz9epftv9lr8m60pk84t"),
		common.MustBech32ToAddress("atx1snputz9gzhyg9cz9hn5alq35q5df2szcpf050e"),
		common.MustBech32ToAddress("atx1aqkmnvq0vve9xglf29qtkl3v2tdnqzxxfpww3q"),

		common.MustBech32ToAddress("atx17va4mfrudm9uv88s03ec0t7xau8297rxfktajz"),
		common.MustBech32ToAddress("atx19e0mfauw87umrzvdulta3keagnrzqs97hl5e9c"),
		common.MustBech32ToAddress("atx19pw0sn4ru9m7rlyl894whjfjngy02xa47pxmu8"),
		common.MustBech32ToAddress("atx1jxllmjyr9xham97ldl5jel20eduj0tkdj6sc4q"),
		common.MustBech32ToAddress("atx1tzmzll6sg6hjy5h3lrkttse59tdrjnmj7uy5kk"),
		common.MustBech32ToAddress("atx13mq3dsgas527sg3v40ztasr23qx9rkffa49lau"),
		common.MustBech32ToAddress("atx1xe8vht0ycdd7utu2qnuzpxatyd453g66untgue"),
		common.MustBech32ToAddress("atx1y6ykcw22rcfqjh5z9edssr5wlgzscuuv9exapw"),
		common.MustBech32ToAddress("atx1x9p98qjv66mmeutp8j4sqyndvpm00guf4dqw9f"),
		common.MustBech32ToAddress("atx124z0qh23u306vjt6lmq0rfwkg5cmyxlq88k9ex"),
		common.MustBech32ToAddress("atx18k5rp7kj56vrm9ydwf3t9t0haffmj5a7n07533"),
		common.MustBech32ToAddress("atx1s9d8jyxqxhe0h9z3ek35n95h3pzfcg5g8t9pwk"),
		common.MustBech32ToAddress("atx1fnw5ncy9slyzf3mzne73ysuskux3q46qhy3y5s"),
		common.MustBech32ToAddress("atx16pqmt742fdepysd92yrlarceecd68ela9sywpn"),
		common.MustBech32ToAddress("atx1e0zc8hkmhe44rwrqxmqyqkttkzsznxnnata8hl"),
		common.MustBech32ToAddress("atx1rs9y2zd6gm02gam44k9jpgvl8x9cypjzcp0wey"),
		common.MustBech32ToAddress("atx1amssl39r4vee7knc3uac966hww8swhkw8n4p88"),
	}

	priKeyArr = []*ecdsa.PrivateKey{
		crypto.HexMustToECDSA("1191dc5317d5930beb77848f416ee023921fa4452f4d783384f35352409c0ad0"),
		crypto.HexMustToECDSA("544b084ebd8a3be88d4817e7015468617407e66ca6de578fcb7e315006ef0d3d"),
		crypto.HexMustToECDSA("d4f11f439304bffa8f014f7d4d5171f12ed84491af948d788dae75c14619773b"),
		crypto.HexMustToECDSA("3115afdf65a417bf830dfce94ad93e73b04114ee4c42bfd14ac6077711b86534"),

		crypto.HexMustToECDSA("8c56e4a0d8bb1f82b94231d535c499fdcbf6e06221acf669d5a964f5bb974903"),
		crypto.HexMustToECDSA("76b3c8c6e9756c7470d5eb4727ee30bdbb5365af523429875d575bf50ad00c7f"),
		crypto.HexMustToECDSA("287c96b9490dd5785ce3005b510b4d3f5bd6ecbbaa27af8d830d9f65224a63aa"),
		crypto.HexMustToECDSA("1e236662904246bfd5c0839d88ebf362269b1695a2a732440af5bdea925d4d37"),

		crypto.HexMustToECDSA("343d10559147d42e1632b4e932aeae36e360d3e0083b9d8d30bb8cc9bb6923c1"),
		crypto.HexMustToECDSA("15439211a0e25c58d7985e11138ce60f675e5243e2b4387fadbd6a0c85755791"),
		crypto.HexMustToECDSA("4a931cfc05fd33b3f3b0f3d910b4358b4cfeac6e1f13b3461a56945ab0de8d96"),
		crypto.HexMustToECDSA("72c8e5bc83fd79debd0af75dab09617198c5f06656ef24009bf7e9a944750bd2"),
		crypto.HexMustToECDSA("d58b015ad107166bd648ba3fb15672e4958f8df668d85acacda7a2ed6f855683"),
		crypto.HexMustToECDSA("1fa19b3862cb9ec584da03d56a84766abdc03cbb3a5e07645531563c1fe2ede6"),
		crypto.HexMustToECDSA("a2be5c2766e9eeed2575448364313cfa91caeb1f1fd03cdbe6f9cee1ded2bffa"),
		crypto.HexMustToECDSA("7da86d7aca8b5dbec9d0bd3c0c2e91552f504df3a42a6e4493992b251bc6c438"),
		crypto.HexMustToECDSA("ed46c6521237ffba7626c67574f8e29d2941ef4bdef561e6d2b4bc877f7c4745"),
		crypto.HexMustToECDSA("b5f8a8bff108a3e674eef019121bdb1c1e0c14857888ff4052954db5700520c3"),
		crypto.HexMustToECDSA("548ceef29a39093e48ef65bc98b210320dedd79ca40acebeb573f8eb72018aac"),
		crypto.HexMustToECDSA("73a2bd8694f883ff5f11551c04303ff7180ae6ef1b89170a67ace10d04c7c3e2"),
		crypto.HexMustToECDSA("996e2bb9c1371e50125fb8b1d0e6f9c46148dfb8b01d9edd6e8b5ec1a6241316"),
		crypto.HexMustToECDSA("51c977a01d5517406fcce2bf7bbb44c67e6b876641a5dac6d2fc26b2f6a97001"),
		crypto.HexMustToECDSA("41d4ce3f8b18fc7ccb4bb0e9514e0863d0c0bd4bb26e9fba3c2a384189c2000b"),
		crypto.HexMustToECDSA("3653b25ba39e59d12a3f45f0fb324b8588db839de4bafd9b938315c356a37051"),
		crypto.HexMustToECDSA("e066f9c4daabcc354162165f8aa161c0bc1cede1b0d14a269f63f6d6bdb1ec5d"),
	}

	blockNumber = big.NewInt(1)
	blockHash   = common.HexToHash("9d4fb5346abcf593ad80a0d3d5a371b22c962418ad34189d5b1b39065668d663")

	blockNumber2 = big.NewInt(2)
	blockHash2   = common.HexToHash("c95876b92443d652d7eb7d7a9c0e2c58a95e934c0c1197978c5445180cc60980")

	blockNumber3 = big.NewInt(3)
	blockHash3   = common.HexToHash("3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e")

	lastBlockNumber uint64
	lastBlockHash   common.Hash
	lastHeader      types.Header

	sender        = common.MustBech32ToAddress("atx1pmhjxvfqeccm87kzpkkr08djgvpp5535gxm6s5")
	anotherSender = common.MustBech32ToAddress("atx1pmhjxvfqeccm87kzpkkr08djgvpp55344s00dx")
	sndb          = snapshotdb.Instance()

	// serial use only
	senderBalance = "9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999"

	txHashArr = []common.Hash{
		common.HexToHash("0x00000000000000000000000000000000000000886d5ba2d3dfb2e2f6a1814f22"),
		common.HexToHash("0x000000000000000000000000000000005249b59609286f2fa91a2abc8555e887"),
		common.HexToHash("0x000000008dba388834e2515c4d9ccb02a48bae177e73959330e55067211c2456"),
		common.HexToHash("0x0000000000000000000000000000000000009a715a765a72b8a289156f9543c9"),
		common.HexToHash("0x0000e1b4a5508c11772b61f463657585c33b577019e4a23bd359c018a4e306d1"),
		common.HexToHash("0x00fd854f940e2d2af8e74c33e640ea6f75c1d9ee49b816b8a4647611d0c91863"),
		common.HexToHash("0x0000000000001038575739a53385cfe42321585a56050e18f8ea2b3e8dc21966"),
		common.HexToHash("0x0000000000000000000000000000000000000048f3b312dc8d081e1186abe8c2"),
		common.HexToHash("0x000000000000000000000000f5bd37579e7ca954eba8fbe7a65646250e92ab7d"),
		common.HexToHash("0x00000000000000000000000000000000000000001d65a5a69fed6ddb0cb58dff"),
		common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000d2"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000f2e8b2706c9e"),
		common.HexToHash("0x00000000000000000000000000e22a393898aac376b079e0894e8e2be6024d03"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000483570dd0679860a"),
		common.HexToHash("0x000000000000000000000000000000000000007fc9e1dc435b5d0064ac50fd4e"),
		common.HexToHash("0x00000000000000000000000000cbeb8f4d51969d7eb70a4f6e8505950d870df7"),
		common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000b4"),
		common.HexToHash("0x000000008fd2abdf28d87efb2c7fa2d37618c8dba97059376d6a58007bee3d8b"),
		common.HexToHash("0x0000000000000000000000003566f3a0adf49d90e610ef3d3548b5a72b1fe199"),
		common.HexToHash("0x00000000000054fa3d19eb57e98aa1dd69d216722054d8539ede4b89c5b77ee9"),
	}

	initProgramVersion = uint32(0<<16 | 8<<8 | 0) // 65536, version: 0.8.0
	promoteVersion     = params.CodeVersion()

	balanceStr = []string{
		"90000000000000000000000000",
		"600000000000000000000000000",
		"13000000000000000000000000",
		"11000000000000000000000000",
		"10000000000000000000000000",
		"48790000000000000000000000",
		"18000000000000000000000000",
		"10000000000000000000000000",
		"10000000000000000000000000",
		"700000000000000000000000000",
		"55500000000000000000000000",
		"90000000000000000000000000",
		"600000000000000000000000000",
		"13000000000000000000000000",
		"11000000000000000000000000",
		"10000000000000000000000000",
		"48790000000000000000000000",
		"18000000000000000000000000",
		"10000000000000000000000000",
		"10000000000000000000000000",
		"700000000000000000000000000",
		"55500000000000000000000000",
		"10000000000000000000000000",
		"700000000000000000000000000",
		"55500000000000000000000000",
	}

	nodeNameArr = []string{
		"PlatON",
		"Gavin",
		"Emma",
		"Kally",
		"Juzhen",
		"Baidu",
		"Alibaba",
		"Tencent",
		"ming",
		"hong",
		"gang",
		"guang",
		"hua",
		"PlatON_2",
		"Gavin_2",
		"Emma_2",
		"Kally_2",
		"Juzhen_2",
		"Baidu_2",
		"Alibaba_2",
		"Tencent_2",
		"ming_2",
		"hong_2",
		"gang_2",
		"guang_2",
	}

	chaList = []string{"A", "a", "B", "b", "C", "c", "D", "d", "E", "e", "F", "f", "G", "g", "H", "h", "J", "j", "K", "k", "M", "m",
		"N", "n", "P", "p", "Q", "q", "R", "r", "S", "s", "T", "t", "U", "u", "V", "v", "W", "w", "X", "x", "Y", "y", "Z", "z"}

	specialCharList = []string{
		"☄", "★", "☎", "☻", "♨", "✠", "❝", "♚", "♘", "✎", "♞", "✩", "✪", "❦", "❥", "❣", "웃", "❂", "Ⓞ", "▶", "◙", "⊕", "◌", "⅓", "∭",
		"∮", "╳", "㏒", "㏕", "‱", "㎏", "❶", "Ň", "🅱", "🅾", "𝖋", "𝕻", "𝕼", "𝕽", "お", "な", "ぬ", "㊎", "㊞", "㊮", "✘"}
)

func TestVersion(t *testing.T) {

	t.Log("the version is:", promoteVersion)
}

func newEvm(blockNumber *big.Int, blockHash common.Hash, state xcom.StateDB) {
	if nil == state {
		state, _, _ = newChainState()
	}
	//evm := &vm.EVM{
	//	StateDB: state,
	//}
	//context := vm.Context{
	//	BlockNumber: blockNumber,
	//	BlockHash:   blockHash,
	//}
	//	evm.Context = context

	//set a default active version

	gov.AddActiveVersion(initProgramVersion, 0, state)

	return
}

func newPlugins() {
	GovPluginInstance()
	StakingInstance()
	SlashInstance()
	RestrictingInstance()
	RewardMgrInstance()

	snapshotdb.Instance()
}

func newChainState() (xcom.StateDB, *types.Block, error) {
	//	testGenesis := new(types.Block)
	chain := mock.NewChain()
	//	var state *state.StateDB

	sBalance, _ := new(big.Int).SetString(senderBalance, 10)
	chain.StateDB.AddBalance(sender, sBalance)
	for i, addr := range addrArr {
		amount, _ := new(big.Int).SetString(balanceStr[len(addrArr)-1-i], 10)
		amount = new(big.Int).Mul(common.Big257, amount)
		chain.StateDB.AddBalance(addr, amount)
	}
	return chain.StateDB, chain.Genesis, nil
}

func build_staking_data_more(block uint64) {

	no := int64(block)
	header := types.Header{
		Number: big.NewInt(no),
	}
	hash := header.Hash()

	stakingDB := staking.NewStakingDB()
	sndb.NewBlock(big.NewInt(int64(block)), lastBlockHash, hash)
	// MOCK

	validatorArr := make(staking.ValidatorQueue, 0)

	// build  more data
	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		rand.Seed(time.Now().UnixNano())

		weight := rand.Intn(1000000000)

		ii := rand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		randBuildFunc := func() (enode.IDv0, common.Address, error) {
			privateKey, err := crypto.GenerateKey()
			if nil != err {
				fmt.Printf("Failed to generate random NodeId private key: %v", err)
				return enode.IDv0{}, common.ZeroAddr, err
			}

			nodeId := enode.PublicKeyToIDv0(&privateKey.PublicKey)

			privateKey, err = crypto.GenerateKey()
			if nil != err {
				fmt.Printf("Failed to generate random Address private key: %v", err)
				return enode.IDv0{}, common.ZeroAddr, err
			}

			addr := crypto.PubkeyToAddress(privateKey.PublicKey)

			return nodeId, addr, nil
		}

		var nodeId enode.IDv0
		var addr common.Address

		if i < 25 {
			nodeId = nodeIdArr[i]
			ar, _ := xutil.NodeId2Addr(nodeId)
			addr = common.Address(ar)
		} else {
			id, ar, err := randBuildFunc()
			if nil != err {
				return
			}
			nodeId = id
			addr = ar
		}

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		canTmp := &staking.Candidate{}

		var blsKeyHex bls.PublicKeyHex
		b, _ := blsKey.GetPublicKey().MarshalText()
		if err := blsKeyHex.UnmarshalText(b); nil != err {
			log.Error("Failed to blsKeyHex.UnmarshalText", "err", err)
			return
		}

		canBase := &staking.CandidateBase{
			NodeId:          nodeId,
			BlsPubKey:       blsKeyHex,
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(i + 1),
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canMutable := &staking.CandidateMutable{
			Shares: balance,
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,
		}

		canTmp.CandidateBase = canBase
		canTmp.CandidateMutable = canMutable

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		stakingDB.SetCanPowerStore(hash, canAddr, canTmp)
		stakingDB.SetCandidateStore(hash, canAddr, canTmp)

		v := &staking.Validator{
			NodeAddress:     canAddr,
			NodeId:          canTmp.NodeId,
			BlsPubKey:       canTmp.BlsPubKey,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			Shares:          canTmp.Shares,
			StakingBlockNum: canTmp.StakingBlockNum,
			StakingTxIndex:  canTmp.StakingTxIndex,
			ValidatorTerm:   0,
		}
		validatorArr = append(validatorArr, v)
	}

	queue := validatorArr[:25]

	epoch_Arr := &staking.ValidatorArray{
		//Start: ((block-1)/22000)*22000 + 1,
		//End:   ((block-1)/22000)*22000 + 22000,
		Start: ((block-1)/xutil.CalcBlocksEachEpoch(currentTestGenesisVersion))*xutil.CalcBlocksEachEpoch(currentTestGenesisVersion) + 1,
		End:   ((block-1)/xutil.CalcBlocksEachEpoch(currentTestGenesisVersion))*xutil.CalcBlocksEachEpoch(currentTestGenesisVersion) + xutil.CalcBlocksEachEpoch(currentTestGenesisVersion),
		Arr:   queue,
	}

	pre_Arr := &staking.ValidatorArray{
		Start: 0,
		End:   0,
		Arr:   queue,
	}

	curr_Arr := &staking.ValidatorArray{
		//Start: ((block-1)/250)*250 + 1,
		//End:   ((block-1)/250)*250 + 250,
		Start: ((block-1)/xcom.ConsensusSize(currentTestGenesisVersion))*xcom.ConsensusSize(currentTestGenesisVersion) + 1,
		End:   ((block-1)/xcom.ConsensusSize(currentTestGenesisVersion))*xcom.ConsensusSize(currentTestGenesisVersion) + xcom.ConsensusSize(currentTestGenesisVersion),
		Arr:   queue,
	}

	setVerifierList(hash, epoch_Arr)
	setRoundValList(hash, pre_Arr)
	setRoundValList(hash, curr_Arr)

	lastBlockHash = hash
	lastBlockNumber = block
	lastHeader = header
}

func build_staking_data(genesisHash common.Hash) {
	stakingDB := staking.NewStakingDB()
	sndb.NewBlock(big.NewInt(1), genesisHash, blockHash)
	// MOCK

	validatorArr := make(staking.ValidatorQueue, 0)

	count := 0
	// build  more data
	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		rand.Seed(time.Now().UnixNano())

		weight := rand.Intn(1000000000)

		ii := rand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		randBuildFunc := func() (enode.IDv0, common.Address, error) {
			privateKey, err := crypto.GenerateKey()
			if nil != err {
				fmt.Printf("Failed to generate random NodeId private key: %v", err)
				return enode.IDv0{}, common.ZeroAddr, err
			}

			nodeId := enode.PublicKeyToIDv0(&privateKey.PublicKey)

			privateKey, err = crypto.GenerateKey()
			if nil != err {
				fmt.Printf("Failed to generate random Address private key: %v", err)
				return enode.IDv0{}, common.ZeroAddr, err
			}

			addr := crypto.PubkeyToAddress(privateKey.PublicKey)

			return nodeId, addr, nil
		}

		var nodeId enode.IDv0
		var addr common.Address

		if i < 25 {
			nodeId = nodeIdArr[i]
			ar, _ := xutil.NodeId2Addr(nodeId)
			addr = common.Address(ar)
		} else {
			id, ar, err := randBuildFunc()
			if nil != err {
				return
			}
			nodeId = id
			addr = ar
		}

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()

		canTmp := &staking.Candidate{}

		var blsKeyHex bls.PublicKeyHex
		b, _ := blsKey.GetPublicKey().MarshalText()
		if err := blsKeyHex.UnmarshalText(b); nil != err {
			log.Error("Failed to blsKeyHex.UnmarshalText", "err", err)
			return
		}

		canBase := &staking.CandidateBase{
			NodeId:          nodeId,
			BlsPubKey:       blsKeyHex,
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(i + 1),
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canMutable := &staking.CandidateMutable{
			Shares: balance,
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,
		}

		canTmp.CandidateBase = canBase
		canTmp.CandidateMutable = canMutable

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		err := stakingDB.SetCanPowerStore(blockHash, canAddr, canTmp)
		if nil != err {
			fmt.Printf("Failed to SetCanPowerStore: %v", err)
			return
		}
		err = stakingDB.SetCandidateStore(blockHash, canAddr, canTmp)
		if nil != err {
			fmt.Printf("Failed to SetCandidateStore: %v", err)
			return
		}

		v := &staking.Validator{
			NodeAddress:     canAddr,
			NodeId:          canTmp.NodeId,
			BlsPubKey:       canTmp.BlsPubKey,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			Shares:          canTmp.Shares,
			StakingBlockNum: canTmp.StakingBlockNum,
			StakingTxIndex:  canTmp.StakingTxIndex,

			ValidatorTerm: 0,
		}
		validatorArr = append(validatorArr, v)
		count++
	}

	fmt.Printf("build staking  data count: %d \n", count)
	queue := validatorArr[:25]

	epoch_Arr := &staking.ValidatorArray{
		Start: 1,
		End:   uint64(xutil.CalcBlocksEachEpoch(currentTestGenesisVersion)),
		Arr:   queue,
	}

	pre_Arr := &staking.ValidatorArray{
		Start: 0,
		End:   0,
		Arr:   queue,
	}

	curr_Arr := &staking.ValidatorArray{
		Start: 1,
		End:   uint64(xcom.ConsensusSize(currentTestGenesisVersion)),
		Arr:   queue,
	}

	setVerifierList(blockHash, epoch_Arr)
	setRoundValList(blockHash, pre_Arr)
	setRoundValList(blockHash, curr_Arr)

	lastBlockHash = blockHash
	lastBlockNumber = blockNumber.Uint64()
	lastHeader = types.Header{
		Number: blockNumber,
	}

}

func buildBlockNoCommit(blockNum int) {

	no := int64(blockNum)
	header := types.Header{
		Number: big.NewInt(no),
	}
	hash := header.Hash()

	staking.NewStakingDB()
	sndb.NewBlock(big.NewInt(int64(blockNum)), lastBlockHash, hash)

	lastBlockHash = hash
	lastBlockNumber = uint64(blockNum)
	lastHeader = header
}

func build_gov_data(state xcom.StateDB) {

	//set a default active version
	gov.AddActiveVersion(initProgramVersion, 0, state)
	gov.InitGenesisGovernParam(common.ZeroHash, snapshotdb.Instance(), 2048)
}

func buildStateDB(t *testing.T) xcom.StateDB {
	chain := mock.NewChain()

	return chain.StateDB
}

type restrictingTest struct {
	restrictingInfo restricting.RestrictingInfo
	plans           []restricting.RestrictingPlan
	account         common.Address
	stateDB         xcom.StateDB
}

func buildDbRestrictingPlan(account common.Address, t *testing.T, stateDB xcom.StateDB) {

	const Epochs = 5
	var list = make([]uint64, 0)

	for epoch := 1; epoch <= Epochs; epoch++ {
		// build release account record
		releaseAccountKey := restricting.GetReleaseAccountKey(uint64(epoch), 1)
		stateDB.SetState(cvm.RestrictingContractAddr, releaseAccountKey, account.Bytes())

		// build release amount record
		releaseAmount := big.NewInt(int64(1e18))
		releaseAmountKey := restricting.GetReleaseAmountKey(uint64(epoch), account)
		stateDB.SetState(cvm.RestrictingContractAddr, releaseAmountKey, releaseAmount.Bytes())

		// build release epoch record
		releaseEpochKey := restricting.GetReleaseEpochKey(uint64(epoch))
		stateDB.SetState(cvm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(1))

		list = append(list, uint64(epoch))
	}

	// build restricting user info
	var user restricting.RestrictingInfo
	user.CachePlanAmount = big.NewInt(int64(5e18))
	user.AdvanceAmount = big.NewInt(0)
	user.NeedRelease = big.NewInt(0)
	user.ReleaseList = list

	bUser, err := rlp.EncodeToBytes(user)
	if err != nil {
		t.Fatalf("failed to rlp encode restricting info: %s", err.Error())
	}

	// build restricting account info record
	restrictingKey := restricting.GetRestrictingKey(account)
	stateDB.SetState(cvm.RestrictingContractAddr, restrictingKey, bUser)

	sBalance, _ := new(big.Int).SetString(senderBalance, 10)
	stateDB.AddBalance(sender, sBalance)
	stateDB.AddBalance(cvm.RestrictingContractAddr, big.NewInt(int64(5e18)))
}

func setRoundValList(blockHash common.Hash, valArr *staking.ValidatorArray) error {

	stakeDB := staking.NewStakingDB()

	queue, err := stakeDB.GetRoundValIndexByBlockHash(blockHash)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to setRoundValList: Query round valIndex is failed", "blockHash",
			blockHash.Hex(), "Start", valArr.Start, "End", valArr.End, "err", err)
		return err
	}

	var indexQueue staking.ValArrIndexQueue

	index := &staking.ValArrIndex{
		Start: valArr.Start,
		End:   valArr.End,
	}

	if len(queue) == 0 {
		indexQueue = make(staking.ValArrIndexQueue, 0)
		_, indexQueue = indexQueue.ConstantAppend(index, RoundValIndexSize)
	} else {

		has := false
		for _, indexInfo := range queue {
			if indexInfo.Start == valArr.Start && indexInfo.End == valArr.End {
				has = true
				break
			}
		}
		indexQueue = queue
		if !has {

			shabby, queue := queue.ConstantAppend(index, RoundValIndexSize)
			indexQueue = queue
			// delete the shabby validators
			if nil != shabby {
				if err := stakeDB.DelRoundValListByBlockHash(blockHash, shabby.Start, shabby.End); nil != err {
					log.Error("Failed to setRoundValList: delete shabby validators is failed",
						"shabby start", shabby.Start, "shabby end", shabby.End, "blockHash", blockHash.Hex())
					return err
				}
			}
		}
	}

	// Store new index Arr
	if err := stakeDB.SetRoundValIndex(blockHash, indexQueue); nil != err {
		log.Error("Failed to setRoundValList: store round validators new indexArr is failed", "blockHash", blockHash.Hex())
		return err
	}

	// Store new round validator Item
	if err := stakeDB.SetRoundValList(blockHash, index.Start, index.End, valArr.Arr); nil != err {
		log.Error("Failed to setRoundValList: store new round validators is failed", "blockHash", blockHash.Hex())
		return err
	}

	return nil
}

func setVerifierList(blockHash common.Hash, valArr *staking.ValidatorArray) error {

	stakeDB := staking.NewStakingDB()

	queue, err := stakeDB.GetEpochValIndexByBlockHash(blockHash)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to setVerifierList: Query epoch valIndex is failed", "blockHash",
			blockHash.Hex(), "Start", valArr.Start, "End", valArr.End, "err", err)
		return err
	}

	var indexQueue staking.ValArrIndexQueue

	index := &staking.ValArrIndex{
		Start: valArr.Start,
		End:   valArr.End,
	}

	if len(queue) == 0 {
		indexQueue = make(staking.ValArrIndexQueue, 0)
		_, indexQueue = indexQueue.ConstantAppend(index, EpochValIndexSize)
	} else {

		has := false
		for _, indexInfo := range queue {
			if indexInfo.Start == valArr.Start && indexInfo.End == valArr.End {
				has = true
				break
			}
		}
		indexQueue = queue
		if !has {

			shabby, queue := queue.ConstantAppend(index, EpochValIndexSize)
			indexQueue = queue
			// delete the shabby validators
			if nil != shabby {
				if err := stakeDB.DelEpochValListByBlockHash(blockHash, shabby.Start, shabby.End); nil != err {
					log.Error("Failed to setVerifierList: delete shabby validators is failed",
						"shabby start", shabby.Start, "shabby end", shabby.End, "blockHash", blockHash.Hex())
					return err
				}
			}
		}
	}

	// Store new index Arr
	if err := stakeDB.SetEpochValIndex(blockHash, indexQueue); nil != err {
		log.Error("Failed to setVerifierList: store epoch validators new indexArr is failed", "blockHash", blockHash.Hex())
		return err
	}

	// Store new epoch validator Item
	if err := stakeDB.SetEpochValList(blockHash, index.Start, index.End, valArr.Arr); nil != err {
		log.Error("Failed to setVerifierList: store new epoch validators is failed", "blockHash", blockHash.Hex())
		return err
	}

	return nil
}
