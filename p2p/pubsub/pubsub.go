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


package pubsub

import (
	"fmt"
	"github.com/AlayaNetwork/Alaya-Go/p2p"
	"github.com/AlayaNetwork/Alaya-Go/p2p/discover"
)

const (

	// CbftPubSubProtocolName is protocol name of CBFT.PubSub
	CbftPubSubProtocolName = "cbft.pubsub"

	// CbftPubSubProtocolVersion is protocol version of CBFT.PubSub
	CbftPubSubProtocolVersion = 1

	// CbftPubSubProtocolLength are the number of implemented message corresponding to cbft.pubsub protocol versions.
	CbftPubSubProtocolLength = 40
)

// PubSub is the implementation of the pubsub system.
type PubSub struct {
	// atomic counter for seqnos
	// NOTE: Must be declared at the top of the struct as we perform atomic
	// operations on this field.
	//
	// See: https://golang.org/pkg/sync/atomic/#pkg-note-BUG
	counter uint64

	// maxMessageSize is the maximum message size; it applies globally to all
	// topics.
	maxMessageSize int

	// size of the outbound message channel that we maintain for each peer
	peerOutboundQueueSize int
}

// Protocols return consensus engine to provide protocol information.
func (p *PubSub) Protocols() []p2p.Protocol {
	return []p2p.Protocol{
		{
			Name:    CbftPubSubProtocolName,
			Version: CbftPubSubProtocolVersion,
			Length:  CbftPubSubProtocolLength,
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				return p.handler(p, rw)
			},
			NodeInfo: func() interface{} {
				return p.NodeInfo()
			},
			PeerInfo: func(id discover.NodeID) interface{} {
				if p, err := p.peers.get(fmt.Sprintf("%x", id[:8])); err == nil {
					return p.Info()
				}
				return nil
			},
		},
	}
}