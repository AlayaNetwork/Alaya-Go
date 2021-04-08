// Copyright 2015 The go-ethereum Authors
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

// MainnetBootnodes are the enode URLs of the P2P bootstrap nodes running on
// the main PlatON network.
var MainnetBootnodes = []string{
	"enode://81dd24640878badc06de82a82fdb0fe55f24d27877144261d81b1f39faa6686ac8f5d2489dbf97cd44d583b4b00976a8f92845378084d25c7a8bae671a543983@ms1.bfa6.alaya.network:16789",
	"enode://32d628cfd32d3f464666792f4fa0bf097c723045f8fe415a8015f1c3cbd0a1bba23e7c76defac277967101dc73e3bd8fc255febb7a52d77c1018ed0cbf8d3ad4@ms2.6cc3.alaya.network:16789",
	"enode://3b2ca03c94a2b8f36b88983d8666947dd08e15347980f95c395b36a4e69218c902894e9e2e92c5a2e0fe8b5c137732d2df40a118766245fdac88c480eb120c18@ms3.cd41.alaya.network:16789",
	"enode://7b5323a73e9cbffd1e6d9178f0b1d55e92649aa71ebe55a0a9c577d374a9ae21ee4980aef2a3214b6e16aa9928ee48df65a382bd2d7ec19f7b87e6d993654d17@ms4.1fda.alaya.network:16789",
	"enode://e6123b585a8e030b42d873d7d09b68847d1f3bba86fab84490fc29acf332a94682a8f8e1518ca857fc75391d62eaf2117703dfeed386b4e0926bf017b5cae445@ms5.ee7a.alaya.network:16789",
	"enode://ab2f7cdf347d4ca26f4fdf5657d7b669464c5712cddc42609ad2060691226187815f0ce87f4dca2cac3ee618d4beeeba9618dbd31c54f97af21d16b7cbf0dccd@ms6.63a8.alaya.network:16789",
	"enode://4c5a092156c43d5aa3dc71f9dc11d304d7631d393725b09e574577c583759e58ddc245e38f993cc0f32fe873cc782bfc1c62fbd49097eec3278b240de785800b@ms7.66dc.alaya.network:16789",
}

// TestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the test network.
var TestnetBootnodes = []string{
	"enode://3fec5e5982a0b32a25168dae575c4705ab8509f266947cb8b16b62ac9eafb78d3e7efce2c31bac447edce3446a12b71383a41dcbdbe80fa856d8739b0214ff35@127.0.0.1:16789",
}

// DiscoveryV5Bootnodes are the enode URLs of the P2P bootstrap nodes for the
// experimental RLPx v5 topic-discovery network.
var DiscoveryV5Bootnodes = []string{}

// AlayanetBootnodes are the enode URLs of the P2P bootstrap nodes running on the test network.
var AlayanetBootnodes = []string{
	"enode://6928c02520ff4e1784d49b4987eee9852dd7b1552f89836292de1002869a78697a9e697d96001e079b5c662135fd5234ae31b3a94e53283d56caf7652f6d6e90@seed1.a72b.alaya.network:16789",
	"enode://3fdefc6e19e46cb05700b947ff8261087706697fe6054ddf925a261af2780084252bf7b2f6cf652e1fec4d64d3e9539ca9408a1b5ead5fe82b555d95cf143fb2@seed2.afc7.alaya.network:16789",
	"enode://3fe92730eb9b53a2e58a9be11a1707c346432ee0c531c24a22d4bb2d0d4a9b4ef04e23988b4fa5d91a790c7f821e9983ae71b03903a3d75bfcce156b060cf99b@seed3.ccf5.alaya.network:16789",
	"enode://1972a5a7d75010e199eac231ab513e564cad5f0e88331a53606b7d55220803c1816d3b0d06ca9b0e10389264f4fade77c46814dd44df502599d3f0a286160498@seed4.5e92.alaya.network:16789",
	"enode://02dc695641f5cada2c685e3bf3dca0218a9dc7a5d5ce8165a2f5bee40d002d18ec6d899abaac1472d88b71e49691019766abd177b8d5d94f72f6f6dc842fded2@seed5.10a1.alaya.network:16789",
	"enode://49648f184dab8acf0927238452e1f7e8f0e86135dcd148baa0a2d22cd931ed9770f726fccf13a3976bbdb738a880a074dc21d811037190267cd9c8c1378a6043@seed6.8d4b.alaya.network:16789",
	"enode://5670e1b34fe39da46ebdd9c3377053c5214c7b3e6b371d31fcc381f788414a38d44cf844ad71305eb1d0d8afddee8eccafb4d30b33b54ca002db47c4864ba080@seed7.bdac.alaya.network:16789",
}

// AlayaTestnetBootnodes are the enode URLs of the P2P bootstrap nodes running on the test network.
var AlayaTestnetBootnodes = []string{}
