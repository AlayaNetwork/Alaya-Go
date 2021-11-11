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
	"errors"
	"github.com/AlayaNetwork/Alaya-Go/log"
	"github.com/AlayaNetwork/Alaya-Go/p2p"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/p2p/pubsub"
	"github.com/AlayaNetwork/Alaya-Go/rlp"
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

var (
	ErrExistsTopic    = errors.New("topic already exists")
	ErrNotExistsTopic = errors.New("topic does not exist")
)

// Group consensus message
type RGMsg struct {
	Code uint64
	Data interface{}
}

type PubSub struct {
	pss *p2p.PubSubServer

	// Messages of all topics subscribed are sent out from this channel uniformly
	msgCh chan *RGMsg
	// All topics subscribed
	topics map[string]*pubsub.Topic
	// The set of topics we are subscribed to
	mySubs map[string]*pubsub.Subscription
	sync.Mutex
}

// Protocol.Run()
func (ps *PubSub) handler(peer *p2p.Peer, rw p2p.MsgReadWriter) error {

	errCh := ps.pss.NewConn(peer, rw)

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

func NewPubSub(server *p2p.PubSubServer) *PubSub {
	return &PubSub{
		pss:    server,
		msgCh:  make(chan *RGMsg),
		topics: make(map[string]*pubsub.Topic),
		mySubs: make(map[string]*pubsub.Subscription),
	}
}

//Subscribe subscribe a topic
func (ps *PubSub) Subscribe(topic string) error {
	ps.Lock()
	defer ps.Unlock()
	if _, ok := ps.mySubs[topic]; ok {
		return ErrExistsTopic
	}
	t, err := ps.pss.PubSub().Join(topic)
	if err != nil {
		return err
	}
	subscription, err := t.Subscribe()
	if err != nil {
		return err
	}
	ps.topics[topic] = t
	ps.mySubs[topic] = subscription

	go ps.listen(subscription)
	return nil
}

func (ps *PubSub) listen(s *pubsub.Subscription) {
	for {
		msg, err := s.Next(context.Background())
		if err != nil {
			if err != pubsub.ErrSubscriptionCancelled {
				ps.Cancel(s.Topic())
			}
			log.Error("Failed to listen to topic message", "error", err)
			return
		}
		if msg != nil {
			var msgData RGMsg
			if err := rlp.DecodeBytes(msg.Data, &msgData); err != nil {
				log.Error("Failed to parse topic message", "error", err)
				ps.Cancel(s.Topic())
				return
			}
			ps.msgCh <- &msgData
		}
	}
}

//UnSubscribe a topic
func (ps *PubSub) Cancel(topic string) error {
	ps.Lock()
	defer ps.Unlock()
	sb := ps.mySubs[topic]
	if sb != nil {
		sb.Cancel()
		delete(ps.mySubs, topic)
		delete(ps.topics, topic)
	}
	return nil
}

func (ps *PubSub) Publish(topic string, data *RGMsg) error {
	ps.Lock()
	defer ps.Unlock()
	t := ps.topics[topic]
	if t == nil {
		return ErrNotExistsTopic
	}
	env, err := rlp.EncodeToBytes(data)
	if err != nil {
		return err
	}
	return t.Publish(context.Background(), env)
}

func (ps *PubSub) Receive() *RGMsg {
	return <-ps.msgCh
}
