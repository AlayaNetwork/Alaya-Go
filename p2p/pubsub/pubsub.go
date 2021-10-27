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
	"github.com/AlayaNetwork/Alaya-Go/p2p"
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

// After the node is successfully connected and the message belongs
// to the cbft.pubsub protocol message, the method is called.
func (pubsub *PubSub) handler(p *p2p.Peer, rw p2p.MsgReadWriter) error {
	return nil
}

// NewPubSub returns a new PubSub management object.
func NewPubSub() *PubSub {
	ps := &PubSub{
		counter: 0,
		maxMessageSize: 2048,
		peerOutboundQueueSize: 32,
	}
	return ps
}