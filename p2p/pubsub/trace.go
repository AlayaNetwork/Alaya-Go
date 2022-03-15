package pubsub

import (
	"time"

	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"

	"github.com/AlayaNetwork/Alaya-Go/p2p/pubsub/message"
)

// EventTracer is a generic event tracer interface.
// This is a high level tracing interface which delivers tracing events, as defined by the protobuf
// schema in pb/trace.proto.
type EventTracer interface {
	Trace(evt *message.TraceEvent)
}

// RawTracer is a low level tracing interface that allows an application to trace the internal
// operation of the pubsub subsystem.
//
// Note that the tracers are invoked synchronously, which means that application tracers must
// take care to not block or modify arguments.
//
// Warning: this interface is not fixed, we may be adding new methods as necessitated by the system
// in the future.
type RawTracer interface {
	// AddPeer is invoked when a new peer is added.
	AddPeer(p *enode.Node, proto ProtocolID)
	// RemovePeer is invoked when a peer is removed.
	RemovePeer(p enode.ID)
	// Join is invoked when a new topic is joined
	Join(topic string)
	// Leave is invoked when a topic is abandoned
	Leave(topic string)
	// Graft is invoked when a new peer is grafted on the mesh (gossipsub)
	Graft(p enode.ID, topic string)
	// Prune is invoked when a peer is pruned from the message (gossipsub)
	Prune(p enode.ID, topic string)
	// ValidateMessage is invoked when a message first enters the validation pipeline.
	ValidateMessage(msg *Message)
	// DeliverMessage is invoked when a message is delivered
	DeliverMessage(msg *Message)
	// RejectMessage is invoked when a message is Rejected or Ignored.
	// The reason argument can be one of the named strings Reject*.
	RejectMessage(msg *Message, reason string)
	// DuplicateMessage is invoked when a duplicate message is dropped.
	DuplicateMessage(msg *Message)
	// ThrottlePeer is invoked when a peer is throttled by the peer gater.
	ThrottlePeer(p enode.ID)
	// RecvRPC is invoked when an incoming RPC is received.
	RecvRPC(rpc *RPC)
	// SendRPC is invoked when a RPC is sent.
	SendRPC(rpc *RPC, p enode.ID)
	// DropRPC is invoked when an outbound RPC is dropped, typically because of a queue full.
	DropRPC(rpc *RPC, p enode.ID)
	// UndeliverableMessage is invoked when the consumer of Subscribe is not reading messages fast enough and
	// the pressure release mechanism trigger, dropping messages.
	UndeliverableMessage(msg *Message)
}

// pubsub tracer details
type pubsubTracer struct {
	tracer EventTracer
	raw    []RawTracer
	pid    enode.ID
	msgID  MsgIdFunction
}

func (t *pubsubTracer) PublishMessage(msg *Message) {
	if t == nil {
		return
	}

	if t.tracer == nil {
		return
	}

	now := time.Now().UnixNano()
	evt := &message.TraceEvent{
		Type:      message.TraceEvent_PUBLISH_MESSAGE.Enum(),
		PeerID:    t.pid.Bytes(),
		Timestamp: &now,
		PublishMessage: &message.TraceEvent_PublishMessage{
			MessageID: []byte(t.msgID(msg.Message)),
			Topic:     msg.Message.Topic,
		},
	}

	t.tracer.Trace(evt)
}

func (t *pubsubTracer) ValidateMessage(msg *Message) {
	if t == nil {
		return
	}

	if msg.ReceivedFrom.ID() != t.pid {
		for _, tr := range t.raw {
			tr.ValidateMessage(msg)
		}
	}
}

func (t *pubsubTracer) RejectMessage(msg *Message, reason string) {
	if t == nil {
		return
	}

	if msg.ReceivedFrom.ID() != t.pid {
		for _, tr := range t.raw {
			tr.RejectMessage(msg, reason)
		}
	}

	if t.tracer == nil {
		return
	}

	now := time.Now().UnixNano()
	evt := &message.TraceEvent{
		Type:      message.TraceEvent_REJECT_MESSAGE.Enum(),
		PeerID:    t.pid.Bytes(),
		Timestamp: &now,
		RejectMessage: &message.TraceEvent_RejectMessage{
			MessageID:    []byte(t.msgID(msg.Message)),
			ReceivedFrom: msg.ReceivedFrom.ID().Bytes(),
			Reason:       &reason,
			Topic:        msg.Topic,
		},
	}

	t.tracer.Trace(evt)
}

func (t *pubsubTracer) DuplicateMessage(msg *Message) {
	if t == nil {
		return
	}

	if msg.ReceivedFrom.ID() != t.pid {
		for _, tr := range t.raw {
			tr.DuplicateMessage(msg)
		}
	}

	if t.tracer == nil {
		return
	}

	now := time.Now().UnixNano()
	evt := &message.TraceEvent{
		Type:      message.TraceEvent_DUPLICATE_MESSAGE.Enum(),
		PeerID:    t.pid.Bytes(),
		Timestamp: &now,
		DuplicateMessage: &message.TraceEvent_DuplicateMessage{
			MessageID:    []byte(t.msgID(msg.Message)),
			ReceivedFrom: msg.ReceivedFrom.ID().Bytes(),
			Topic:        msg.Topic,
		},
	}

	t.tracer.Trace(evt)
}

func (t *pubsubTracer) DeliverMessage(msg *Message) {
	if t == nil {
		return
	}

	if msg.ReceivedFrom.ID() != t.pid {
		for _, tr := range t.raw {
			tr.DeliverMessage(msg)
		}
	}

	if t.tracer == nil {
		return
	}

	now := time.Now().UnixNano()
	evt := &message.TraceEvent{
		Type:      message.TraceEvent_DELIVER_MESSAGE.Enum(),
		PeerID:    t.pid.Bytes(),
		Timestamp: &now,
		DeliverMessage: &message.TraceEvent_DeliverMessage{
			MessageID:    []byte(t.msgID(msg.Message)),
			Topic:        msg.Topic,
			ReceivedFrom: msg.ReceivedFrom.ID().Bytes(),
		},
	}

	t.tracer.Trace(evt)
}

func (t *pubsubTracer) AddPeer(p *enode.Node, proto ProtocolID) {
	if t == nil {
		return
	}

	for _, tr := range t.raw {
		tr.AddPeer(p, proto)
	}

	if t.tracer == nil {
		return
	}

	protoStr := string(proto)
	now := time.Now().UnixNano()
	evt := &message.TraceEvent{
		Type:      message.TraceEvent_ADD_PEER.Enum(),
		PeerID:    t.pid.Bytes(),
		Timestamp: &now,
		AddPeer: &message.TraceEvent_AddPeer{
			PeerID: p.ID().Bytes(),
			Proto:  &protoStr,
		},
	}

	t.tracer.Trace(evt)
}

func (t *pubsubTracer) RemovePeer(p enode.ID) {
	if t == nil {
		return
	}

	for _, tr := range t.raw {
		tr.RemovePeer(p)
	}

	if t.tracer == nil {
		return
	}

	now := time.Now().UnixNano()
	evt := &message.TraceEvent{
		Type:      message.TraceEvent_REMOVE_PEER.Enum(),
		PeerID:    t.pid.Bytes(),
		Timestamp: &now,
		RemovePeer: &message.TraceEvent_RemovePeer{
			PeerID: p.Bytes(),
		},
	}

	t.tracer.Trace(evt)
}

func (t *pubsubTracer) RecvRPC(rpc *RPC) {
	if t == nil {
		return
	}

	for _, tr := range t.raw {
		tr.RecvRPC(rpc)
	}

	if t.tracer == nil {
		return
	}

	now := time.Now().UnixNano()
	evt := &message.TraceEvent{
		Type:      message.TraceEvent_RECV_RPC.Enum(),
		PeerID:    t.pid.Bytes(),
		Timestamp: &now,
		RecvRPC: &message.TraceEvent_RecvRPC{
			ReceivedFrom: rpc.from.ID().Bytes(),
			Meta:         t.traceRPCMeta(rpc),
		},
	}

	t.tracer.Trace(evt)
}

func (t *pubsubTracer) SendRPC(rpc *RPC, p enode.ID) {
	if t == nil {
		return
	}

	for _, tr := range t.raw {
		tr.SendRPC(rpc, p)
	}

	if t.tracer == nil {
		return
	}

	now := time.Now().UnixNano()
	evt := &message.TraceEvent{
		Type:      message.TraceEvent_SEND_RPC.Enum(),
		PeerID:    t.pid.Bytes(),
		Timestamp: &now,
		SendRPC: &message.TraceEvent_SendRPC{
			SendTo: p.Bytes(),
			Meta:   t.traceRPCMeta(rpc),
		},
	}

	t.tracer.Trace(evt)
}

func (t *pubsubTracer) DropRPC(rpc *RPC, p enode.ID) {
	if t == nil {
		return
	}

	for _, tr := range t.raw {
		tr.DropRPC(rpc, p)
	}

	if t.tracer == nil {
		return
	}

	now := time.Now().UnixNano()
	evt := &message.TraceEvent{
		Type:      message.TraceEvent_DROP_RPC.Enum(),
		PeerID:    t.pid.Bytes(),
		Timestamp: &now,
		DropRPC: &message.TraceEvent_DropRPC{
			SendTo: p.Bytes(),
			Meta:   t.traceRPCMeta(rpc),
		},
	}

	t.tracer.Trace(evt)
}

func (t *pubsubTracer) UndeliverableMessage(msg *Message) {
	if t == nil {
		return
	}

	for _, tr := range t.raw {
		tr.UndeliverableMessage(msg)
	}
}

func (t *pubsubTracer) traceRPCMeta(rpc *RPC) *message.TraceEvent_RPCMeta {
	rpcMeta := new(message.TraceEvent_RPCMeta)

	var msgs []*message.TraceEvent_MessageMeta
	for _, m := range rpc.Publish {
		msgs = append(msgs, &message.TraceEvent_MessageMeta{
			MessageID: []byte(t.msgID(m)),
			Topic:     m.Topic,
		})
	}
	rpcMeta.Messages = msgs

	var subs []*message.TraceEvent_SubMeta
	for _, sub := range rpc.Subscriptions {
		subs = append(subs, &message.TraceEvent_SubMeta{
			Subscribe: sub.Subscribe,
			Topic:     sub.Topicid,
		})
	}
	rpcMeta.Subscription = subs

	if rpc.Control != nil {
		var ihave []*message.TraceEvent_ControlIHaveMeta
		for _, ctl := range rpc.Control.Ihave {
			var mids [][]byte
			for _, mid := range ctl.MessageIDs {
				mids = append(mids, []byte(mid))
			}
			ihave = append(ihave, &message.TraceEvent_ControlIHaveMeta{
				Topic:      ctl.TopicID,
				MessageIDs: mids,
			})
		}

		var iwant []*message.TraceEvent_ControlIWantMeta
		for _, ctl := range rpc.Control.Iwant {
			var mids [][]byte
			for _, mid := range ctl.MessageIDs {
				mids = append(mids, []byte(mid))
			}
			iwant = append(iwant, &message.TraceEvent_ControlIWantMeta{
				MessageIDs: mids,
			})
		}

		var graft []*message.TraceEvent_ControlGraftMeta
		for _, ctl := range rpc.Control.Graft {
			graft = append(graft, &message.TraceEvent_ControlGraftMeta{
				Topic: ctl.TopicID,
			})
		}

		var prune []*message.TraceEvent_ControlPruneMeta
		for _, ctl := range rpc.Control.Prune {
			peers := make([][]byte, 0, len(ctl.Peers))
			for _, pi := range ctl.Peers {
				peers = append(peers, pi.PeerID.Bytes())
			}
			prune = append(prune, &message.TraceEvent_ControlPruneMeta{
				Topic: ctl.TopicID,
				Peers: peers,
			})
		}

		rpcMeta.Control = &message.TraceEvent_ControlMeta{
			Ihave: ihave,
			Iwant: iwant,
			Graft: graft,
			Prune: prune,
		}
	}

	return rpcMeta
}

func (t *pubsubTracer) Join(topic string) {
	if t == nil {
		return
	}

	for _, tr := range t.raw {
		tr.Join(topic)
	}

	if t.tracer == nil {
		return
	}

	now := time.Now().UnixNano()
	evt := &message.TraceEvent{
		Type:      message.TraceEvent_JOIN.Enum(),
		PeerID:    t.pid.Bytes(),
		Timestamp: &now,
		Join: &message.TraceEvent_Join{
			Topic: &topic,
		},
	}

	t.tracer.Trace(evt)
}

func (t *pubsubTracer) Leave(topic string) {
	if t == nil {
		return
	}

	for _, tr := range t.raw {
		tr.Leave(topic)
	}

	if t.tracer == nil {
		return
	}

	now := time.Now().UnixNano()
	evt := &message.TraceEvent{
		Type:      message.TraceEvent_LEAVE.Enum(),
		PeerID:    t.pid.Bytes(),
		Timestamp: &now,
		Leave: &message.TraceEvent_Leave{
			Topic: &topic,
		},
	}

	t.tracer.Trace(evt)
}

func (t *pubsubTracer) Graft(p enode.ID, topic string) {
	if t == nil {
		return
	}

	for _, tr := range t.raw {
		tr.Graft(p, topic)
	}

	if t.tracer == nil {
		return
	}

	now := time.Now().UnixNano()
	evt := &message.TraceEvent{
		Type:      message.TraceEvent_GRAFT.Enum(),
		PeerID:    t.pid.Bytes(),
		Timestamp: &now,
		Graft: &message.TraceEvent_Graft{
			PeerID: p.Bytes(),
			Topic:  &topic,
		},
	}

	t.tracer.Trace(evt)
}

func (t *pubsubTracer) Prune(p enode.ID, topic string) {
	if t == nil {
		return
	}

	for _, tr := range t.raw {
		tr.Prune(p, topic)
	}

	if t.tracer == nil {
		return
	}

	now := time.Now().UnixNano()
	evt := &message.TraceEvent{
		Type:      message.TraceEvent_PRUNE.Enum(),
		PeerID:    t.pid.Bytes(),
		Timestamp: &now,
		Prune: &message.TraceEvent_Prune{
			PeerID: p.Bytes(),
			Topic:  &topic,
		},
	}

	t.tracer.Trace(evt)
}

func (t *pubsubTracer) ThrottlePeer(p enode.ID) {
	if t == nil {
		return
	}

	for _, tr := range t.raw {
		tr.ThrottlePeer(p)
	}
}
