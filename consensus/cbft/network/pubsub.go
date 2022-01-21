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
	"sync"

	"github.com/AlayaNetwork/Alaya-Go/core/cbfttypes"
	"github.com/AlayaNetwork/Alaya-Go/event"

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
	onReceive   receiveCallback

	// All topics subscribed
	topics      map[string]*pubsub.Topic
	topicCtx    map[string]context.Context
	topicCancel map[string]context.CancelFunc
	// The set of topics we are subscribed to
	mySubs map[string]*pubsub.Subscription
	sync.Mutex

	quit chan struct{}
}

// Protocol.Run()
func (ps *PubSub) handler(peer *p2p.Peer, rw p2p.MsgReadWriter) error {
	log.Debug("Start PubSub's processors", "id", peer.ID().TerminalString())
	errCh := ps.pss.NewConn(peer, rw)
	defer ps.pss.Host().DisConn(peer.ID())

	handlerErr := <-errCh
	log.Info("pubsub's handler ends", "id", peer.ID().TerminalString(), "err", handlerErr)

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
				return ps.pss.GetPeerInfo(id)
			},
		},
	}
}

func NewPubSub(server *p2p.PubSubServer) *PubSub {
	return &PubSub{
		pss:         server,
		topics:      make(map[string]*pubsub.Topic),
		topicCtx:    make(map[string]context.Context),
		topicCancel: make(map[string]context.CancelFunc),
		mySubs:      make(map[string]*pubsub.Subscription),
		quit:        make(chan struct{}),
	}
}

func (ps *PubSub) Start(config ctypes.Config, get getByIDFunc, onReceive receiveCallback, eventMux *event.TypeMux) {
	ps.config = config
	ps.getPeerById = get
	ps.onReceive = onReceive

	go ps.watching(eventMux)
}

func (ps *PubSub) watching(eventMux *event.TypeMux) {
	events := eventMux.Subscribe(cbfttypes.GroupTopicEvent{}, cbfttypes.ExpiredGroupTopicEvent{})
	defer events.Unsubscribe()

	for {
		select {
		case ev := <-events.Chan():
			if ev == nil {
				continue
			}
			switch data := ev.Data.(type) {
			case cbfttypes.GroupTopicEvent:
				log.Trace("Received GroupTopicEvent", "topic", data.Topic)
				// 需要订阅主题（发现节点并接收topic消息）
				if data.PubSub {
					ps.Subscribe(data.Topic)
				} else { // 只需要发现节点，不需要接收对应topic消息
					ps.DiscoverTopic(data.Topic)
				}
			case cbfttypes.ExpiredGroupTopicEvent:
				log.Trace("Received ExpiredGroupTopicEvent", "topic", data.Topic)
				ps.Cancel(data.Topic)
			default:
				log.Error("Received unexcepted event")
			}
		case <-ps.quit:
			return
		}
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
	ctx, cancel := context.WithCancel(context.Background())
	ps.pss.DiscoverTopic(ctx, topic)

	subscription, err := t.Subscribe()
	if err != nil {
		return err
	}
	ps.topics[topic] = t
	ps.mySubs[topic] = subscription
	ps.topicCtx[topic] = ctx
	ps.topicCancel[topic] = cancel

	go ps.listen(subscription)
	return nil
}

// 发现topic对应的节点
func (ps *PubSub) DiscoverTopic(topic string) error {
	ps.Lock()
	defer ps.Unlock()
	if _, ok := ps.topicCtx[topic]; ok {
		return ErrExistsTopic
	}
	ctx, cancel := context.WithCancel(context.Background())
	ps.pss.DiscoverTopic(ctx, topic)

	ps.topicCtx[topic] = ctx
	ps.topicCancel[topic] = cancel
	return nil
}

func (ps *PubSub) listen(s *pubsub.Subscription) {
	for {
		subMsg, err := s.Next(context.Background())
		if err != nil {
			if err != pubsub.ErrSubscriptionCancelled {
				ps.Cancel(s.Topic())
				log.Error("Failed to listen to topic message", "topic", s.Topic(), "error", err)
			}
			return
		}
		if subMsg != nil {
			var gmsg GMsg
			if err := rlp.DecodeBytes(subMsg.Data, &gmsg); err != nil {
				log.Error("Failed to parse topic message", "topic", s.Topic(), "error", err)
				ps.Cancel(s.Topic())
				return
			}
			msg := p2p.Msg{
				Code:    gmsg.Code,
				Size:    uint32(len(subMsg.Data)),
				Payload: bytes.NewReader(common.CopyBytes(gmsg.Data)),
			}
			if ps.pss.Host().ID().ID() == subMsg.From {
				log.Trace("Receive a message from myself", "fromId", subMsg.From.TerminalString(), "topic", s.Topic(), "msgCode", gmsg.Code)
				continue
			}
			fromPeer, err := ps.getPeerById(subMsg.ReceivedFrom.ID().TerminalString())
			if err != nil {
				log.Error("Failed to execute getPeerById", "receivedFrom", subMsg.ReceivedFrom.ID().TerminalString(), "topic", s.Topic(), "err", err)
			} else {
				log.Trace("Receive a message", "topic", s.Topic(), "receivedFrom", fromPeer.ID().TerminalString(), "msgFrom", subMsg.From.TerminalString(), "msgCode", gmsg.Code)
				ps.onReceive(fromPeer, &msg)
			}
		}
	}
}

//UnSubscribe a topic
func (ps *PubSub) Cancel(topic string) error {
	ps.Lock()
	defer ps.Unlock()
	sb, ok := ps.mySubs[topic]
	if ok && sb != nil {
		sb.Cancel()
		delete(ps.mySubs, topic)
	}
	t, ok := ps.topics[topic]
	if ok && t != nil {
		t.Close()
		delete(ps.topics, topic)
	}

	if cancel, ok := ps.topicCancel[topic]; ok {
		cancel()
		delete(ps.topicCtx, topic)
		delete(ps.topicCancel, topic)
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

	return t.Publish(ps.topicCtx[topic], env)
}

func (ps *PubSub) Stop() {
	close(ps.quit)
}

func (ps *PubSub) GetAllPubSubStatus() *pubsub.Status {
	return ps.pss.GetAllPubSubStatus()
}
