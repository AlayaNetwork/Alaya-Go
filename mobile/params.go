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

// Contains all the wrappers from the params package.

package platon

import (
	"encoding/json"

	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"

	"github.com/AlayaNetwork/Alaya-Go/core"
	"github.com/AlayaNetwork/Alaya-Go/params"
)

// TestnetGenesis returns the JSON spec to use for the Alpha test network.
func TestnetGenesis() string {
	enc, err := json.Marshal(core.DefaultTestnetGenesisBlock())
	if err != nil {
		panic(err)
	}
	return string(enc)
}

// FoundationBootnodes returns the enode URLs of the P2P bootstrap nodes operated
// by the foundation running the V5 discovery protocol.
func FoundationBootnodes() *Enodes {
	nodes := &Enodes{nodes: make([]*enode.Node, len(params.DiscoveryV5Bootnodes))}
	for i, url := range params.DiscoveryV5Bootnodes {
		var err error
		nodes.nodes[i], err = enode.Parse(enode.ValidSchemes, url)
		if err != nil {
			panic("invalid node URL: " + err.Error())
		}
	}
	return nodes
}
