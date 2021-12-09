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
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/AlayaNetwork/Alaya-Go/common"
	ctypes "github.com/AlayaNetwork/Alaya-Go/consensus/cbft/types"
	"github.com/AlayaNetwork/Alaya-Go/log"
	"github.com/AlayaNetwork/Alaya-Go/p2p"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/p2p/pubsub"
	"github.com/AlayaNetwork/Alaya-Go/rlp"
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
type GMsg struct {
	Code uint64
	Data []byte
}

type PubSub struct {
	pss         *p2p.PubSubServer
	config      ctypes.Config
	getPeerById getByIDFunc // Used to get peer by ID.

	onReceive receiveCallback

	// All topics subscribed
	topics   map[string]*pubsub.Topic
	topicCtx map[string]context.Context
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

func (ps *PubSub) Config() *ctypes.Config {
	return &ps.config
}

func (ps *PubSub) NodeInfo() interface{} {
	cfg := ps.Config()
	return &NodeInfo{
		Config: *cfg,
	}
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
				if p, err := ps.getPeerById(fmt.Sprintf("%x", id[:8])); err == nil {
					return p.Info()
				}
				return nil
			},
		},
	}
}

func NewPubSub(server *p2p.PubSubServer) *PubSub {
	return &PubSub{
		pss:    server,
		topics: make(map[string]*pubsub.Topic),
		mySubs: make(map[string]*pubsub.Subscription),
	}
}

func (ps *PubSub) Init(config ctypes.Config, get getByIDFunc, onReceive receiveCallback) {
	ps.config = config
	ps.getPeerById = get
	ps.onReceive = onReceive
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
	ctx := context.Background()
	ps.pss.DiscoverTopic(ctx, topic)

	subscription, err := t.Subscribe()
	if err != nil {
		return err
	}
	ps.topics[topic] = t
	ps.mySubs[topic] = subscription
	ps.topicCtx[topic] = ctx

	go ps.listen(subscription)
	return nil
}

func (ps *PubSub) listen(s *pubsub.Subscription) {
	for {
		subMsg, err := s.Next(context.Background())
		if err != nil {
			if err != pubsub.ErrSubscriptionCancelled {
				ps.Cancel(s.Topic())
			}
			log.Error("Failed to listen to topic message", "error", err)
			return
		}
		if subMsg != nil {
			var gmsg GMsg
			if err := rlp.DecodeBytes(subMsg.Data, &gmsg); err != nil {
				log.Error("Failed to parse topic message", "error", err)
				ps.Cancel(s.Topic())
				return
			}
			msg := p2p.Msg{
				Code:    gmsg.Code,
				Size:    uint32(len(subMsg.Data)),
				Payload: bytes.NewReader(common.CopyBytes(gmsg.Data)),
			}
			if ps.pss.Host().ID().ID() == subMsg.From {
				log.Trace("Receive a message from myself", "fromId", subMsg.From.TerminalString())
				continue
			}
			fromPeer, err := ps.getPeerById(subMsg.ReceivedFrom.ID().TerminalString())
			if err != nil {
				log.Error("Failed to execute getPeerById", "err", err)
			} else {
				ps.onReceive(fromPeer, &msg)
			}
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
	if ctx, ok := ps.topicCtx[topic]; ok {
		ctx.Done()
		delete(ps.topicCtx, topic)
	}
	return nil
}

func (ps *PubSub) Publish(topic string, code uint64, data interface{}) error {
	ps.Lock()
	defer ps.Unlock()
	t := ps.topics[topic]
	if t == nil {
		return ErrNotExistsTopic
	}
	dataEnv, err := rlp.EncodeToBytes(data)
	if err != nil {
		return err
	}
	gmsg := &GMsg{
		Code: code,
		Data: dataEnv,
	}
	env, err := rlp.EncodeToBytes(gmsg)
	if err != nil {
		return err
	}

	return t.Publish(context.Background(), env)
}
