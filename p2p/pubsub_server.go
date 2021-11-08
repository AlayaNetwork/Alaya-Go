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
	"errors"
	"fmt"
	ctypes "github.com/AlayaNetwork/Alaya-Go/consensus/cbft/types"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/p2p/pubsub"
	"sync"
	"time"
)

const (

	// CbftProtocolName is protocol name of CBFT.pubsub.
	CbftPubSubProtocolName = "cbft.pubsub"

	// CbftProtocolVersion is protocol version of CBFT.pubsub.
	CbftPubSubProtocolVersion = 1

	// CbftProtocolLength are the number of implemented message corresponding to cbft.pubsub protocol versions.
	CbftPubSubProtocolLength = 10

	// Msg code of PubSub's message
	PubSubMsgCode = 0xff
)

type PubSubServer struct {
	lock      sync.Mutex // protects running
	running   bool
	ps        *pubsub.PubSub
	host      *Host
	peerMsgCh chan *ctypes.MsgInfo
}

// Start starts running the server.
// Servers can not be re-used after stopping.
func (srv *PubSubServer) Start() (err error) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if srv.running {
		return errors.New("server already running")
	}
	srv.running = true
	return nil
}

// run is the main loop of the server.
func (s *PubSubServer) run() {

}

// After the node is successfully connected and the message belongs
// to the cbft.pubsub protocol message, the method is called.
func (s *PubSubServer) Handle(p *Peer, rw MsgReadWriter) error {

	return nil
}

// watchingAddTopicValidatorEvent
func (s *PubSubServer) watchingAddTopicValidatorEvent() {

}

// watchingAddTopicValidatorEvent
func (s *PubSubServer) watchingRemoveTopicValidatorEvent() {

}

// PublishMsg
func (s *PubSubServer) PublishMsg(topic string) error {
	if topic == "" {
		return fmt.Errorf("topic is nil")
	}
	return nil
}

func (s *PubSubServer) ReadTopicLoop() {

}

// Protocol.NodeInfo()
func (s *PubSubServer) NodeInfo() interface{} {
	return nil
}

// Protocol.Run()
func (s *PubSubServer) handler(peer *Peer, rw MsgReadWriter) error {
	direction := pubsub.DirOutbound
	if peer.Inbound() {
		direction = pubsub.DirInbound
	}
	stat := pubsub.Stat{
		Direction: direction,
		Opened:    time.Now(),
		Extra:     make(map[interface{}]interface{}),
	}
	conn := &Conn{
		remote: peer.Node(),
		stat:   stat,
	}
	conn.streams.m = make(map[pubsub.Stream]struct{})
	stream := &Stream{
		rw:   rw,
		conn: conn,
	}
	conn.streams.m[stream] = struct{}{}
	s.host.SetStream(peer.ID(), stream)
	s.host.network.NotifyAll(conn)

	return nil
}

//Protocols implemented the Protocols method and returned basic information about the CBFT.pubsub protocol.
func (s *PubSubServer) Protocols() []Protocol {
	return []Protocol{
		{
			Name:    CbftPubSubProtocolName,
			Version: CbftPubSubProtocolVersion,
			Length:  CbftPubSubProtocolLength,
			Run: func(p *Peer, rw MsgReadWriter) error {
				return s.handler(p, rw)
			},
			NodeInfo: func() interface{} {
				return s.NodeInfo()
			},
			PeerInfo: func(id enode.ID) interface{} {
				return nil
			},
		},
	}
}

func NewPubSubServer(node *enode.Node, server *Server) *PubSubServer {
	network := NewNetwork(server)
	host := &Host{
		node:    node,
		streams: make(map[enode.ID]pubsub.Stream),
		network: network,
	}
	pubsub, err := pubsub.NewGossipSub(context.Background(), host)
	if err != nil {
		panic("Failed to NewGossipSub: " + err.Error())
	}

	return &PubSubServer{
		ps:        pubsub,
		host:      host,
		peerMsgCh: make(chan *ctypes.MsgInfo),
	}
}
