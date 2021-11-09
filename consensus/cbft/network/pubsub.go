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

package network

import (
	"context"
	"fmt"
	"github.com/AlayaNetwork/Alaya-Go/log"
	"github.com/AlayaNetwork/Alaya-Go/p2p"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/p2p/pubsub"
)

const (

	// CbftProtocolName is protocol name of CBFT.pubsub.
	CbftPubSubProtocolName = "cbft.pubsub"

	// CbftProtocolVersion is protocol version of CBFT.pubsub.
	CbftPubSubProtocolVersion = 1

	// CbftProtocolLength are the number of implemented message corresponding to cbft.pubsub protocol versions.
	CbftPubSubProtocolLength = 10
)

type PubSub struct {
	host *p2p.Host
	ps   *pubsub.PubSub

	topicChan chan string
	topics    []string
	// The set of topics we are subscribed to
	mySubs map[string]map[*pubsub.Subscription]struct{}
}

// Protocol.Run()
func (ps *PubSub) handler(peer *p2p.Peer, rw p2p.MsgReadWriter) error {

	conn := p2p.NewConn(peer.Node(), peer.Inbound())

	// Wait for the connection to exit
	errCh := make(chan error)

	stream := p2p.NewStream(conn, rw, errCh, "")
	conn.SetStream(stream)

	ps.host.SetStream(peer.ID(), stream)
	ps.host.NotifyAll(conn)

	handlerErr := <-errCh
	log.Info("pubsub's handler ends", "err", handlerErr)
	return handlerErr
}

func (ps *PubSub) NodeInfo() interface{} {
	return nil
}

//Protocols implemented the Protocols method and returned basic information about the CBFT.pubsub protocol.
func (ps *PubSub) Protocols() []p2p.Protocol {
	return []p2p.Protocol{
		{
			Name:    CbftPubSubProtocolName,
			Version: CbftPubSubProtocolVersion,
			Length:  CbftPubSubProtocolLength,
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				return ps.handler(p, rw)
			},
			NodeInfo: func() interface{} {
				return ps.NodeInfo()
			},
			PeerInfo: func(id enode.ID) interface{} {
				return nil
			},
		},
	}
}

func NewPubSub(localNode *enode.Node, server *p2p.Server) *PubSub {
	network := p2p.NewNetwork(server)
	host := p2p.NewHost(localNode, network)
	pubsub, err := pubsub.NewGossipSub(context.Background(), host)
	if err != nil {
		panic("Failed to NewGossipSub: " + err.Error())
	}

	return &PubSub{
		ps:   pubsub,
		host: host,
	}
}

//Subscribe subscribe a topic
func (ps *PubSub) Subscribe(topic string) error {
	//if err := ps.ps.PublishMsg(topic); err != nil {
	//	return err
	//}
	ps.topics = append(ps.topics, topic)
	return nil
}

//UnSubscribe a topic
func (ps *PubSub) Cancel(topic string) error {
	sm := ps.mySubs[topic]
	if sm != nil {
		for s, _ := range sm {
			if s != nil {
				s.Cancel()
			}
		}
		return nil
	}

	return fmt.Errorf("Can not find", "topic", topic)
}
