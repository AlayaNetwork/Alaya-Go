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

package p2p

import (
	"context"

	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/p2p/pubsub"
)

type PubSubServer struct {
	p2pServer *Server
	pubSub    *pubsub.PubSub
	host      *Host
}

func NewPubSubServer(ctx context.Context, localNode *enode.Node, p2pServer *Server) *PubSubServer {
	network := NewNetwork(p2pServer.Peers)
	host := NewHost(localNode, network)
	// If the tracer proxy address is specified, the remote tracer function will be enabled.
	tracers := make([]pubsub.Option, 0)
	if p2pServer.PubSubTraceHost != "" {
		remoteTracer, _ := pubsub.NewRemoteTracer(ctx, p2pServer.PubSubTraceHost)
		tracers = append(tracers, pubsub.WithEventTracer(remoteTracer))
	}
	gossipSub, err := pubsub.NewGossipSub(ctx, host, tracers...)
	if err != nil {
		panic("Failed to NewGossipSub: " + err.Error())
	}

	return &PubSubServer{
		p2pServer: p2pServer,
		pubSub:    gossipSub,
		host:      host,
	}
}

func (pss *PubSubServer) Host() *Host {
	return pss.host
}

func (pss *PubSubServer) PubSub() *pubsub.PubSub {
	return pss.pubSub
}

func (pss *PubSubServer) NewConn(peer *Peer, rw MsgReadWriter) chan error {
	conn := NewConn(peer.Node(), peer.Inbound())

	// Wait for the connection to exit
	errCh := make(chan error)

	stream := NewStream(conn, rw, errCh, pubsub.GossipSubID_v11)
	conn.SetStream(stream)

	pss.Host().SetStream(peer.ID(), stream)
	pss.Host().AddConn(peer.ID(), conn)
	pss.Host().NotifyAll(conn)
	return errCh
}

func (pss *PubSubServer) DiscoverTopic(ctx context.Context, topic string) {
	pss.p2pServer.DiscoverTopic(ctx, topic)
}

func (pss *PubSubServer) GetAllPubSubStatus() *pubsub.Status {
	return pss.pubSub.GetAllPubSubStatus()
}

func (pss *PubSubServer) GetPeerInfo(nodeId enode.ID) *pubsub.PeerInfo {
	return pss.pubSub.GetPeerInfo(nodeId)
}
