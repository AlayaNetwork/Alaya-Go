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
	"errors"
	"fmt"
	ctypes "github.com/AlayaNetwork/Alaya-Go/consensus/cbft/types"
	"github.com/AlayaNetwork/Alaya-Go/log"
	"github.com/AlayaNetwork/Alaya-Go/p2p"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"sync"
)

const (

	// CbftProtocolName is protocol name of CBFT.pubsub.
	CbftPubSubProtocolName = "cbft.pubsub"

	// CbftProtocolVersion is protocol version of CBFT.pubsub.
	CbftPubSubProtocolVersion = 1

	// CbftProtocolLength are the number of implemented message corresponding to cbft.pubsub protocol versions.
	CbftPubSubProtocolLength = 10

)

type SubServer struct {
	lock		sync.Mutex // protects running
	running 	bool
	ps			*PubSub
	peerMsgCh	chan *ctypes.MsgInfo
}

var (
	subServerOnce sync.Once
	svr           *SubServer
)

func SubServerInstance() *SubServer {
	subServerOnce.Do(func() {
		log.Info("Init SubServer ...")
		svr = NewPubSubServer()
	})
	return svr
}

// Start starts running the server.
// Servers can not be re-used after stopping.
func (srv *SubServer) Start() (err error) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if srv.running {
		return errors.New("server already running")
	}
	srv.running = true
	return nil
}

// run is the main loop of the server.
func (s *SubServer) run() {

}

// After the node is successfully connected and the message belongs
// to the cbft.pubsub protocol message, the method is called.
func (s *SubServer) Handle(p *p2p.Peer, rw p2p.MsgReadWriter) error {
	return nil
}

// watchingAddTopicValidatorEvent
func (s *SubServer) watchingAddTopicValidatorEvent() {

}

// watchingAddTopicValidatorEvent
func (s *SubServer) watchingRemoveTopicValidatorEvent() {

}

// PublishMsg
func (s *SubServer) PublishMsg(topic string) error  {
	if topic == "" {
		return fmt.Errorf("topic is nil")
	}
	return nil
}

func (s *SubServer) ReadTopicLoop()  {

}

// Protocol.NodeInfo()
func (s *SubServer) NodeInfo() interface{}  {
	return nil
}

// Protocol.Run()
func (s *SubServer) handler(peer *p2p.Peer, rw p2p.MsgReadWriter) error  {
	return nil
}

//Protocols implemented the Protocols method and returned basic information about the CBFT.pubsub protocol.
func (s *SubServer) Protocols() []p2p.Protocol {
	return []p2p.Protocol{
		{
			Name:    CbftPubSubProtocolName,
			Version: CbftPubSubProtocolVersion,
			Length:  CbftPubSubProtocolLength,
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
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

func NewPubSubServer() *SubServer {
	return &SubServer{}
}