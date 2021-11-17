// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"fmt"
	"math/big"

	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"

	"github.com/AlayaNetwork/Alaya-Go/common"
	"github.com/AlayaNetwork/Alaya-Go/crypto/bls"
)

// Genesis hashes to enforce below configs on.
var (
	TestnetGenesisHash  = common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d")
	AlayanetGenesisHash = common.HexToHash("0xfb787fede6752e1a5ad85d2c6fc140454759be5e69d86d2425ceac22c23bd419")
)

var TrustedCheckpoints = map[common.Hash]*TrustedCheckpoint{
	//MainnetGenesisHash: MainnetTrustedCheckpoint,
	//TestnetGenesisHash: TestnetTrustedCheckpoint,
}

var (
	initialTestnetConsensusNodes = []initNode{
		{
			"enode://b7f1f7757a900cce7ce4caf8663ecf871205763ac201c65f9551d5b841731a9cd9550bc05f3a16fbc2ef589c9faeef74d4500b60d76047939e2ba7fa4a5915aa@127.0.0.1:16789",
			"f1735bac863706b49809a4e635fe0c2e224aef5ad549f18ba3f2f6b61c0c9d0005f12d497a301ba26a8aaf009c90e4198301875002984c5cd9bd614cd2fbcb81c57f6355a8400d56c20804e1dfb34782c1f2eadda82c8b226aa4a71bfa60be8c",
		},
		{
			"enode://3b2acc72f673173a97295728f6d9f93a8d75d1c615455f3ec3fdc4471e707e54935d89ba4736082db4b618b05de203fe828a8182b0f1ba09d495ad9a8ddb418b@127.0.0.1:16789",
			"bd0d378e9d87e552d6f3842b38e30776e45eb14af3462822cbadf5eea492477dd764d10cc73c521e05fa90c6146dd70a6ae1bd25f473147b91ddd65e5077ead12a1010de2714ef7977067df4b519ab2f7d50db7c2e150dc2d5cb2bb6cc30e485",
		},
		{
			"enode://184fb0464cd84a28e6c9aec564c058113b6c93ae80eccf1dc50be0481fd27bbfea3dfef492898987aeabba07fd0b5b7048a88163da348b729b0b56b6184f6e6c@127.0.0.1:16789",
			"fd89eb74d9277b8a02eb32f2144bf571d7badbca85ea7ea1158f47f268e36119a8e83608b4b972a5fd332bdc4f0306149390ed0afa4f6bae3e18f80a6062feff0a53cd265c38a3bff43afdf93c06dd8cdd0928a57039aaf59712623cd412a38d",
		},
		{
			"enode://31df231e05e089ee517f577a4ded288d210fc9d313c4ed862a7a758e4a74d8ae4b84de69be837371120a521199bce5ecf8665a56cd1efc2c8b665dd563a3f590@127.0.0.1:16789",
			"1c6792a7868006106547d48fb77e2941240352fd5def1f447b38c88585655e02bdf91209b48db698866a92b59914d00cacbb4206f99ece4e78fb82a12c732039d5a17e23b9980cd26421f953c3017199dbd41b62f13227d69123f431fdae3506",
		},
		{
			"enode://50ed65eb0180f771bbef453f0750bb324b626e6283b10b5203cae28e2e50744ba4bae51242f5d5dd17ce8d5bb71c50599f0432fd64b6a64d8026de181d8c5ae3@127.0.0.1:16789",
			"0e82ee8823646e871d319a5e418c5845af4efebbffc873ace0b27ff1391c720fafb98a118447cca88efa2653cb1d8a1443fd38079bd8385e6291a52776158882360657b3c401c1a1dc2fb2cfc083247f454ea160fc6aadecdbc1b44ee7aa6f8c",
		},
		{
			"enode://487a68847885a7b198b24920c1a5addc82d7c62fb43a503d7a4701a9fd72deed5e4f53116d8cca139d66731a2dc7537d5234d400188710864b83a7ccd65ab74d@127.0.0.1:16789",
			"a41e72db25f7e4d2af738e6170d5d84c72c69a185be75f55a62fb5ad0b931061f8f32c13d851f96d0b819f67d079c205b12d73507fa156967c8e6f385c2ab5b373f6f4a5e3e86bd95fe9857001a1feabcdf9507b45f446a7066213563a8d7a81",
		},
		{
			"enode://081da72f53a7cc4eaeace9f8a2960c1c01bf606153fd5acf47f540e0f2a1ebe5aed36fc8bf10263308b18845fe0079e027c85af31e03f995b6834d88f8461a19@127.0.0.1:16789",
			"77d174270903955fd8b38de00c05c17b078031722e1fc414188251a5da34aa8025dc3cb75bebedf5cd8eb91d25bee213a5fd78f12ed77892bd2f2e4a21595a413280e37d1d1726bf9fa62e042307de89ee6e243304ef768bac3a9d466dc3a406",
		},
	}

	initialAlayaConsensusNodes = []initNode{
		{
			"enode://ba0f7995ee0cf98e18e82daf9560f2007eb46a510be0b3bcc6ec7ab6f4e7f611ce783036916aa0a0af95ddc1ad691df1fe84cfab926f358aec1a786e79d9eb7b@fdn1.92fd.alaya.network:16789",
			"ccb421751a3ae082df3c134cbcc0f9e36c0ef2ddb641b1cfa4f647abe522bc5db1c432eb94368c7624a340a9021980084c0247240830f7da28d5a290c6b1e018d2f7f1b4c4a72b0dca4503f7fc6d43b63bee3534d2fb592d20f97be3e3a6048e",
		},
		{
			"enode://b273bd13e4d82793ff5a874ef514934db333bcb3b12330ecf91e831ef6e536db61b4179bea8886095e931c1c0e749367eaf84266e083fa55c48c9fe1d50e26df@fdn2.7edb.alaya.network:16789",
			"c7e05f20d53715719bd22b5e4cec453157ae5e81b3edf737e33c37b0839c630a77129e12550ae1571d1bec399b12a619ed987307435304d4b169e460929516be8d298a7493afc7449359ce76731938bb57b550578843c3ae6ba41e3d63ffd88a",
		},
		{
			"enode://45cc53b5279503b6476bb2e5bcc8f0d3f1e7f6526ade9f27f9fea528f9bda9f5c6b8a5379235363b158188add776d14a5baebe55e8e8d57cdaefcea05307aa95@fdn3.3ef0.alaya.network:16789",
			"2d317f94fa7583de8cd0d72069f3b104bcd6d8c1bc24265d5e66218d31eb4db19a2ddef75b74090fc34b951bfe4dac1535f91c5819569997131d1de688c1d3a2b7ae4ec3bfccf5a7d6af0412e01e1c6937311232ecc0a1c5bce951e84d425900",
		},
		{
			"enode://695055cfccec7536b3f1e12208af7742a8aa30f0d7ed916e557a7cd50f3bedc9044e3916ce6d4d43aff4ae09e1df0e06842f7eb9dfd226fe49f4e0a6c1ad7a0a@fdn4.550d.alaya.network:16789",
			"8f8d77b2e4d2e2e36000f56a26562fa5a3780b3263ca0e1f0533b9a14a3530d55da4dd84fe0a8713639f54cf0e7564141f1ce2779feee022e25a29bc20f2eed9c2367bba535e92f1057a1601b368a3cd2366a929e2c5f74d5e5bf45096063d18",
		},
		{
			"enode://bc93ef4c138c4e010a57ba7c4565faaf07f5cdbe4bff58f92f7aa3e9a10e830dcb4a6f63fe04e61872c5ea1f98a302cfb6cb0004538db0cf3c7c8d697dc5bfee@fdn5.b1d5.alaya.network:16789",
			"dd2be8a4a50cda5c8cbe6cc8d38aff0e18991749fff68018f86c4c46fe7d48f301a216648136fae5d27c6db128761910488f2ae886390384fffed9068fc18132c60e2874d210a1ad890eca1f99051e3aa9d7ef748df57ed9a2b93bb370e9fd80",
		},
		{
			"enode://a860a56adaa9d3bd66494973d44b54bdfdadaeabdf0377b40e9dd99235a1d764dfc6791a19e744bef1d5e1962cc8de8afbf6202b922e2a411783a3ed9d622c2b@fdn6.5hj4.alaya.network:16789",
			"ce809302eb36897b937874ad84bc1e506c543252510f1ee20ab6578457b321667f8faf18fb1143b7e1af418e56033603062dc809705e220a39032e86f3b49ada23224e28798d4b50b82b80f03d609c555acce12ddb7b3f7259d23c769b708506",
		},
		{
			"enode://25af23c768bb57bbb5b72e349cf23bbb371e7359a3c0436cc3c22f28edbfa3429a511cd1f05783f4b385f84cd0649884fd36b8d3018b0a108ed7e7b189f41566@fdn7.5wf6.alaya.network:16789",
			"05fee124fdb890a4c795142228c8812308138640fb3a173324af4e8e2df13a6b5a93af19edeb730ec884104e5dd4f70e60e578470a730fe5b3a99dd52590304603df7bc189bdac4d556e736d5659ad2ffd20a14cd5fbf567952fb4b9a678118d",
		},
	}

	// AlayaChainConfig is the chain parameters to run a node on the main network.
	AlayaChainConfig = &ChainConfig{
		ChainID:     big.NewInt(201018),
		AddressHRP:  "atp",
		EmptyBlock:  "on",
		EIP155Block: big.NewInt(1),
		Cbft: &CbftConfig{
			InitialNodes:  ConvertNodeUrl(initialAlayaConsensusNodes),
			Amount:        10,
			ValidatorMode: "ppos",
			Period:        20000,
		},
		GenesisVersion: GenesisVersion,
	}

	// TestnetChainConfig is the chain parameters to run a node on the test network.
	TestnetChainConfig = &ChainConfig{
		ChainID:     big.NewInt(104),
		AddressHRP:  "atp",
		EmptyBlock:  "on",
		EIP155Block: big.NewInt(1),
		Cbft: &CbftConfig{
			InitialNodes:  ConvertNodeUrl(initialTestnetConsensusNodes),
			Amount:        10,
			ValidatorMode: "ppos",
			Period:        20000,
		},
		GenesisVersion: GenesisVersion,
	}

	GrapeChainConfig = &ChainConfig{
		AddressHRP:  "atp",
		ChainID:     big.NewInt(304),
		EmptyBlock:  "on",
		EIP155Block: big.NewInt(3),
		Cbft: &CbftConfig{
			Period: 3,
		},
		GenesisVersion: GenesisVersion,
	}

	// AllEthashProtocolChanges contains every protocol change (EIPs) introduced
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllEthashProtocolChanges = &ChainConfig{big.NewInt(1337), "atp", "", big.NewInt(0), big.NewInt(0), nil, nil, GenesisVersion}

	TestChainConfig = &ChainConfig{big.NewInt(1), "atp", "", big.NewInt(0), big.NewInt(0), nil, new(CbftConfig), GenesisVersion}
)

// TrustedCheckpoint represents a set of post-processed trie roots (CHT and
// BloomTrie) associated with the appropriate section index and head hash. It is
// used to start light syncing from this checkpoint and avoid downloading the
// entire header chain while still being able to securely access old headers/logs.
type TrustedCheckpoint struct {
	Name         string      `json:"-"`
	SectionIndex uint64      `json:"sectionIndex"`
	SectionHead  common.Hash `json:"sectionHead"`
	CHTRoot      common.Hash `json:"chtRoot"`
	BloomRoot    common.Hash `json:"bloomRoot"`
}

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	ChainID     *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection
	AddressHRP  string   `json:"addressHRP"`
	EmptyBlock  string   `json:"emptyBlock"`
	EIP155Block *big.Int `json:"eip155Block,omitempty"` // EIP155 HF block
	EWASMBlock  *big.Int `json:"ewasmBlock,omitempty"`  // EWASM switch block (nil = no fork, 0 = already activated)
	// Various consensus engines
	Clique *CliqueConfig `json:"clique,omitempty"`
	Cbft   *CbftConfig   `json:"cbft,omitempty"`

	GenesisVersion uint32 `json:"genesisVersion"`
}

type CbftNode struct {
	Node      *enode.Node   `json:"node"`
	BlsPubKey bls.PublicKey `json:"blsPubKey"`
}

type initNode struct {
	Enode     string
	BlsPubkey string
}

type CbftConfig struct {
	Period        uint64     `json:"period,omitempty"`        // Number of seconds between blocks to enforce
	Amount        uint32     `json:"amount,omitempty"`        //The maximum number of blocks generated per cycle
	InitialNodes  []CbftNode `json:"initialNodes,omitempty"`  //Genesis consensus node
	ValidatorMode string     `json:"validatorMode,omitempty"` //Validator mode for easy testing
	GroupValidatorsLimit uint32	`json:"GroupValidatorsLimit,omitempty"` //Max validators per group
	CoordinatorLimit uint32	`json:"CoordinatorLimit,omitempty"` //Coordinators Limit C0>C1>C2...
}

// CliqueConfig is the consensus engine configs for proof-of-authority based sealing.
type CliqueConfig struct {
	Period uint64 `json:"period"` // Number of seconds between blocks to enforce
	Epoch  uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
}

// String implements the stringer interface, returning the consensus engine details.
func (c *CliqueConfig) String() string {
	return "clique"
}

// String implements the fmt.Stringer interface.
func (c *ChainConfig) String() string {
	var engine interface{}
	switch {
	case c.Clique != nil:
		engine = c.Clique
	case c.Cbft != nil:
		engine = c.Cbft
	default:
		engine = "unknown"
	}
	return fmt.Sprintf("{ChainID: %v EIP155: %v Engine: %v}",
		c.ChainID,
		c.EIP155Block,
		engine,
	)
}

// IsEIP155 returns whether num is either equal to the EIP155 fork block or greater.
func (c *ChainConfig) IsEIP155(num *big.Int) bool {
	//	return isForked(c.EIP155Block, num)
	return true
}

// IsEWASM returns whether num represents a block number after the EWASM fork
func (c *ChainConfig) IsEWASM(num *big.Int) bool {
	return isForked(c.EWASMBlock, num)
}

// GasTable returns the gas table corresponding to the current phase (homestead or homestead reprice).
//
// The returned GasTable's fields shouldn't, under any circumstances, be changed.
func (c *ChainConfig) GasTable(num *big.Int) GasTable {
	return GasTableConstantinople
}

// CheckCompatible checks whether scheduled fork transitions have been imported
// with a mismatching chain configuration.
func (c *ChainConfig) CheckCompatible(newcfg *ChainConfig, height uint64) *ConfigCompatError {
	bhead := new(big.Int).SetUint64(height)

	// Iterate checkCompatible to find the lowest conflict.
	var lasterr *ConfigCompatError
	for {
		err := c.checkCompatible(newcfg, bhead)
		if err == nil || (lasterr != nil && err.RewindTo == lasterr.RewindTo) {
			break
		}
		lasterr = err
		bhead.SetUint64(err.RewindTo)
	}
	return lasterr
}

func (c *ChainConfig) checkCompatible(newcfg *ChainConfig, head *big.Int) *ConfigCompatError {
	if isForkIncompatible(c.EIP155Block, newcfg.EIP155Block, head) {
		return newCompatError("EIP155 fork block", c.EIP155Block, newcfg.EIP155Block)
	}
	if isForkIncompatible(c.EWASMBlock, newcfg.EWASMBlock, head) {
		return newCompatError("ewasm fork block", c.EWASMBlock, newcfg.EWASMBlock)
	}
	return nil
}

// isForkIncompatible returns true if a fork scheduled at s1 cannot be rescheduled to
// block s2 because head is already past the fork.
func isForkIncompatible(s1, s2, head *big.Int) bool {
	return (isForked(s1, head) || isForked(s2, head)) && !configNumEqual(s1, s2)
}

// isForked returns whether a fork scheduled at block s is active at the given head block.
func isForked(s, head *big.Int) bool {
	if s == nil || head == nil {
		return false
	}
	return s.Cmp(head) <= 0
}

func configNumEqual(x, y *big.Int) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return x.Cmp(y) == 0
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *big.Int
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}

func newCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{what, storedblock, newblock, 0}
	if rew != nil && rew.Sign() > 0 {
		err.RewindTo = rew.Uint64() - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
}

func ConvertNodeUrl(initialNodes []initNode) []CbftNode {
	bls.Init(bls.BLS12_381)
	NodeList := make([]CbftNode, 0, len(initialNodes))
	for _, n := range initialNodes {

		cbftNode := new(CbftNode)

		if node, err := enode.Parse(enode.ValidSchemes, n.Enode); nil == err {
			cbftNode.Node = node
		}

		if n.BlsPubkey != "" {
			var blsPk bls.PublicKey
			if err := blsPk.UnmarshalText([]byte(n.BlsPubkey)); nil == err {
				cbftNode.BlsPubKey = blsPk
			}
		}

		NodeList = append(NodeList, *cbftNode)
	}
	return NodeList
}
