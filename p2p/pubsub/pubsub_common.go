package pubsub

import (
	"context"
	"github.com/AlayaNetwork/Alaya-Go/p2p/pubsub/message"

	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"

	"github.com/AlayaNetwork/Alaya-Go/log"

	"github.com/gogo/protobuf/proto"
)

// get the initial RPC containing all of our subscriptions to send to new peers
func (p *PubSub) getHelloPacket() *RPC {
	var rpc RPC

	subscriptions := make(map[string]bool)

	for t := range p.mySubs {
		subscriptions[t] = true
	}

	for t := range p.myRelays {
		subscriptions[t] = true
	}

	for t := range subscriptions {
		as := &message.RPC_SubOpts{
			Topicid:   proto.String(t),
			Subscribe: proto.Bool(true),
		}
		rpc.Subscriptions = append(rpc.Subscriptions, as)
	}
	return &rpc
}

func (p *PubSub) handleNewStream(s Stream) {
	peer := s.Conn().RemotePeer()

	p.inboundStreamsMx.Lock()
	_, dup := p.inboundStreams[peer.ID()]
	if dup {
		log.Debug("duplicate inbound stream , resetting other stream", "from", peer.ID().TerminalString())
	}
	p.inboundStreams[peer.ID()] = s
	p.inboundStreamsMx.Unlock()

	defer func() {
		p.inboundStreamsMx.Lock()
		if p.inboundStreams[peer.ID()] == s {
			delete(p.inboundStreams, peer.ID())
		}
		p.inboundStreamsMx.Unlock()
	}()

	for {
		rpc := new(RPC)
		if err := s.Read(&rpc.RPC); err != nil {
			log.Error("Read message error", "id", peer.ID().TerminalString(), "err", err)
			p.notifyPeerDead(peer.ID())
			s.Close(err)
			return
		}
		rpc.from = peer
		select {
		case p.incoming <- rpc:
		case <-p.ctx.Done():
			s.Close(nil)
			return
		}
	}
}

func (p *PubSub) notifyPeerDead(pid enode.ID) {
	p.peerDeadPrioLk.RLock()
	p.peerDeadMx.Lock()
	p.peerDeadPend[pid] = struct{}{}
	p.peerDeadMx.Unlock()
	p.peerDeadPrioLk.RUnlock()

	select {
	case p.peerDead <- struct{}{}:
	default:
	}
}

func (p *PubSub) handleNewPeer(ctx context.Context, pid enode.ID, outgoing <-chan *RPC) {
	s, err := p.host.NewStream(p.ctx, pid, p.rt.Protocols()...)
	if err != nil || s == nil {
		log.Debug("opening new stream to peer: ", "id", pid.TerminalString(), "err", err)

		select {
		case p.newPeerError <- pid:
		case <-ctx.Done():
		}

		return
	}

	go p.host.StreamHandler(s.Protocol())(s)
	go p.handleSendingMessages(ctx, s, outgoing)
	select {
	case p.newPeerStream <- s:
	case <-ctx.Done():
	}
}

func (p *PubSub) handleSendingMessages(ctx context.Context, s Stream, outgoing <-chan *RPC) {
	for {
		select {
		case rpc, ok := <-outgoing:
			if !ok {
				return
			}

			if !message.IsEmpty(&rpc.RPC) {
				message.Filling(&rpc.RPC)
				if err := s.Write(&rpc.RPC); err != nil {
					log.Error("Send message fail", "id", s.Conn().ID(), "err", err)
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func rpcWithSubs(subs ...*message.RPC_SubOpts) *RPC {
	return &RPC{
		RPC: message.RPC{
			Subscriptions: subs,
			Publish:       make([]*message.Message, 0),
			Control:       &message.ControlMessage{},
		},
	}
}

func rpcWithMessages(msgs ...*message.Message) *RPC {
	return &RPC{RPC: message.RPC{Publish: msgs, Subscriptions: make([]*message.RPC_SubOpts, 0), Control: &message.ControlMessage{}}}
}

func rpcWithControl(msgs []*message.Message,
	ihave []*message.ControlIHave,
	iwant []*message.ControlIWant,
	graft []*message.ControlGraft,
	prune []*message.ControlPrune) *RPC {
	return &RPC{
		RPC: message.RPC{
			Subscriptions: make([]*message.RPC_SubOpts, 0),
			Publish:       msgs,
			Control: &message.ControlMessage{
				Ihave: ihave,
				Iwant: iwant,
				Graft: graft,
				Prune: prune,
			},
		},
	}
}

func copyRPC(rpc *RPC) *RPC {
	res := new(RPC)
	*res = *rpc
	if rpc.Control != nil {
		res.Control = new(message.ControlMessage)
		*res.Control = *rpc.Control
	}
	return res
}
