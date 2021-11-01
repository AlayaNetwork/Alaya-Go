package message

import (
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"

	proto "github.com/gogo/protobuf/proto"
)

type TraceEvent_Type int32

const (
	TraceEvent_PUBLISH_MESSAGE   TraceEvent_Type = 0
	TraceEvent_REJECT_MESSAGE    TraceEvent_Type = 1
	TraceEvent_DUPLICATE_MESSAGE TraceEvent_Type = 2
	TraceEvent_DELIVER_MESSAGE   TraceEvent_Type = 3
	TraceEvent_ADD_PEER          TraceEvent_Type = 4
	TraceEvent_REMOVE_PEER       TraceEvent_Type = 5
	TraceEvent_RECV_RPC          TraceEvent_Type = 6
	TraceEvent_SEND_RPC          TraceEvent_Type = 7
	TraceEvent_DROP_RPC          TraceEvent_Type = 8
	TraceEvent_JOIN              TraceEvent_Type = 9
	TraceEvent_LEAVE             TraceEvent_Type = 10
	TraceEvent_GRAFT             TraceEvent_Type = 11
	TraceEvent_PRUNE             TraceEvent_Type = 12
)

var TraceEvent_Type_name = map[int32]string{
	0:  "PUBLISH_MESSAGE",
	1:  "REJECT_MESSAGE",
	2:  "DUPLICATE_MESSAGE",
	3:  "DELIVER_MESSAGE",
	4:  "ADD_PEER",
	5:  "REMOVE_PEER",
	6:  "RECV_RPC",
	7:  "SEND_RPC",
	8:  "DROP_RPC",
	9:  "JOIN",
	10: "LEAVE",
	11: "GRAFT",
	12: "PRUNE",
}

var TraceEvent_Type_value = map[string]int32{
	"PUBLISH_MESSAGE":   0,
	"REJECT_MESSAGE":    1,
	"DUPLICATE_MESSAGE": 2,
	"DELIVER_MESSAGE":   3,
	"ADD_PEER":          4,
	"REMOVE_PEER":       5,
	"RECV_RPC":          6,
	"SEND_RPC":          7,
	"DROP_RPC":          8,
	"JOIN":              9,
	"LEAVE":             10,
	"GRAFT":             11,
	"PRUNE":             12,
}

func (x TraceEvent_Type) Enum() *TraceEvent_Type {
	p := new(TraceEvent_Type)
	*p = x
	return p
}

func (x TraceEvent_Type) String() string {
	return proto.EnumName(TraceEvent_Type_name, int32(x))
}

func (x *TraceEvent_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(TraceEvent_Type_value, data, "TraceEvent_Type")
	if err != nil {
		return err
	}
	*x = TraceEvent_Type(value)
	return nil
}

type TraceEvent struct {
	Type             *TraceEvent_Type             `protobuf:"varint,1,opt,name=type,enum=pubsub.pb.TraceEvent_Type" json:"type,omitempty"`
	PeerID           enode.ID                     `protobuf:"bytes,2,opt,name=peerID" json:"peerID,omitempty"`
	Timestamp        *int64                       `protobuf:"varint,3,opt,name=timestamp" json:"timestamp,omitempty"`
	PublishMessage   *TraceEvent_PublishMessage   `protobuf:"bytes,4,opt,name=publishMessage" json:"publishMessage,omitempty"`
	RejectMessage    *TraceEvent_RejectMessage    `protobuf:"bytes,5,opt,name=rejectMessage" json:"rejectMessage,omitempty"`
	DuplicateMessage *TraceEvent_DuplicateMessage `protobuf:"bytes,6,opt,name=duplicateMessage" json:"duplicateMessage,omitempty"`
	DeliverMessage   *TraceEvent_DeliverMessage   `protobuf:"bytes,7,opt,name=deliverMessage" json:"deliverMessage,omitempty"`
	AddPeer          *TraceEvent_AddPeer          `protobuf:"bytes,8,opt,name=addPeer" json:"addPeer,omitempty"`
	RemovePeer       *TraceEvent_RemovePeer       `protobuf:"bytes,9,opt,name=removePeer" json:"removePeer,omitempty"`
	RecvRPC          *TraceEvent_RecvRPC          `protobuf:"bytes,10,opt,name=recvRPC" json:"recvRPC,omitempty"`
	SendRPC          *TraceEvent_SendRPC          `protobuf:"bytes,11,opt,name=sendRPC" json:"sendRPC,omitempty"`
	DropRPC          *TraceEvent_DropRPC          `protobuf:"bytes,12,opt,name=dropRPC" json:"dropRPC,omitempty"`
	Join             *TraceEvent_Join             `protobuf:"bytes,13,opt,name=join" json:"join,omitempty"`
	Leave            *TraceEvent_Leave            `protobuf:"bytes,14,opt,name=leave" json:"leave,omitempty"`
	Graft            *TraceEvent_Graft            `protobuf:"bytes,15,opt,name=graft" json:"graft,omitempty"`
	Prune            *TraceEvent_Prune            `protobuf:"bytes,16,opt,name=prune" json:"prune,omitempty"`
}

func (m *TraceEvent) Reset()         { *m = TraceEvent{} }
func (m *TraceEvent) String() string { return "" }
func (*TraceEvent) ProtoMessage()    {}

func (m *TraceEvent) GetType() TraceEvent_Type {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return TraceEvent_PUBLISH_MESSAGE
}

func (m *TraceEvent) GetPeerID() enode.ID {
	if m != nil {
		return m.PeerID
	}
	return enode.ZeroID
}

func (m *TraceEvent) GetTimestamp() int64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func (m *TraceEvent) GetPublishMessage() *TraceEvent_PublishMessage {
	if m != nil {
		return m.PublishMessage
	}
	return nil
}

func (m *TraceEvent) GetRejectMessage() *TraceEvent_RejectMessage {
	if m != nil {
		return m.RejectMessage
	}
	return nil
}

func (m *TraceEvent) GetDuplicateMessage() *TraceEvent_DuplicateMessage {
	if m != nil {
		return m.DuplicateMessage
	}
	return nil
}

func (m *TraceEvent) GetDeliverMessage() *TraceEvent_DeliverMessage {
	if m != nil {
		return m.DeliverMessage
	}
	return nil
}

func (m *TraceEvent) GetAddPeer() *TraceEvent_AddPeer {
	if m != nil {
		return m.AddPeer
	}
	return nil
}

func (m *TraceEvent) GetRemovePeer() *TraceEvent_RemovePeer {
	if m != nil {
		return m.RemovePeer
	}
	return nil
}

func (m *TraceEvent) GetRecvRPC() *TraceEvent_RecvRPC {
	if m != nil {
		return m.RecvRPC
	}
	return nil
}

func (m *TraceEvent) GetSendRPC() *TraceEvent_SendRPC {
	if m != nil {
		return m.SendRPC
	}
	return nil
}

func (m *TraceEvent) GetDropRPC() *TraceEvent_DropRPC {
	if m != nil {
		return m.DropRPC
	}
	return nil
}

func (m *TraceEvent) GetJoin() *TraceEvent_Join {
	if m != nil {
		return m.Join
	}
	return nil
}

func (m *TraceEvent) GetLeave() *TraceEvent_Leave {
	if m != nil {
		return m.Leave
	}
	return nil
}

func (m *TraceEvent) GetGraft() *TraceEvent_Graft {
	if m != nil {
		return m.Graft
	}
	return nil
}

func (m *TraceEvent) GetPrune() *TraceEvent_Prune {
	if m != nil {
		return m.Prune
	}
	return nil
}

type TraceEvent_PublishMessage struct {
	MessageID []byte  `protobuf:"bytes,1,opt,name=messageID" json:"messageID,omitempty"`
	Topic     *string `protobuf:"bytes,2,opt,name=topic" json:"topic,omitempty"`
}

func (m *TraceEvent_PublishMessage) Reset()         { *m = TraceEvent_PublishMessage{} }
func (m *TraceEvent_PublishMessage) String() string { return "" }
func (*TraceEvent_PublishMessage) ProtoMessage()    {}

func (m *TraceEvent_PublishMessage) GetMessageID() []byte {
	if m != nil {
		return m.MessageID
	}
	return nil
}

func (m *TraceEvent_PublishMessage) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

type TraceEvent_RejectMessage struct {
	MessageID    []byte   `protobuf:"bytes,1,opt,name=messageID" json:"messageID,omitempty"`
	ReceivedFrom enode.ID `protobuf:"bytes,2,opt,name=receivedFrom" json:"receivedFrom,omitempty"`
	Reason       *string  `protobuf:"bytes,3,opt,name=reason" json:"reason,omitempty"`
	Topic        *string  `protobuf:"bytes,4,opt,name=topic" json:"topic,omitempty"`
}

func (m *TraceEvent_RejectMessage) Reset()         { *m = TraceEvent_RejectMessage{} }
func (m *TraceEvent_RejectMessage) String() string { return "" }
func (*TraceEvent_RejectMessage) ProtoMessage()    {}

func (m *TraceEvent_RejectMessage) GetMessageID() []byte {
	if m != nil {
		return m.MessageID
	}
	return nil
}

func (m *TraceEvent_RejectMessage) GetReceivedFrom() enode.ID {
	if m != nil {
		return m.ReceivedFrom
	}
	return enode.ZeroID
}

func (m *TraceEvent_RejectMessage) GetReason() string {
	if m != nil && m.Reason != nil {
		return *m.Reason
	}
	return ""
}

func (m *TraceEvent_RejectMessage) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

type TraceEvent_DuplicateMessage struct {
	MessageID    []byte   `protobuf:"bytes,1,opt,name=messageID" json:"messageID,omitempty"`
	ReceivedFrom enode.ID `protobuf:"bytes,2,opt,name=receivedFrom" json:"receivedFrom,omitempty"`
	Topic        *string  `protobuf:"bytes,3,opt,name=topic" json:"topic,omitempty"`
}

func (m *TraceEvent_DuplicateMessage) Reset()         { *m = TraceEvent_DuplicateMessage{} }
func (m *TraceEvent_DuplicateMessage) String() string { return "" }
func (*TraceEvent_DuplicateMessage) ProtoMessage()    {}

func (m *TraceEvent_DuplicateMessage) GetMessageID() []byte {
	if m != nil {
		return m.MessageID
	}
	return nil
}

func (m *TraceEvent_DuplicateMessage) GetReceivedFrom() enode.ID {
	if m != nil {
		return m.ReceivedFrom
	}
	return enode.ZeroID
}

func (m *TraceEvent_DuplicateMessage) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

type TraceEvent_DeliverMessage struct {
	MessageID    []byte   `protobuf:"bytes,1,opt,name=messageID" json:"messageID,omitempty"`
	Topic        *string  `protobuf:"bytes,2,opt,name=topic" json:"topic,omitempty"`
	ReceivedFrom enode.ID `protobuf:"bytes,3,opt,name=receivedFrom" json:"receivedFrom,omitempty"`
}

func (m *TraceEvent_DeliverMessage) Reset()         { *m = TraceEvent_DeliverMessage{} }
func (m *TraceEvent_DeliverMessage) String() string { return "" }
func (*TraceEvent_DeliverMessage) ProtoMessage()    {}

func (m *TraceEvent_DeliverMessage) GetMessageID() []byte {
	if m != nil {
		return m.MessageID
	}
	return nil
}

func (m *TraceEvent_DeliverMessage) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

func (m *TraceEvent_DeliverMessage) GetReceivedFrom() enode.ID {
	if m != nil {
		return m.ReceivedFrom
	}
	return enode.ZeroID
}

type TraceEvent_AddPeer struct {
	PeerID enode.ID `json:"peerID,omitempty"`
	Proto  *string  `json:"proto,omitempty"`
}

func (m *TraceEvent_AddPeer) Reset()         { *m = TraceEvent_AddPeer{} }
func (m *TraceEvent_AddPeer) String() string { return "" }
func (*TraceEvent_AddPeer) ProtoMessage()    {}

func (m *TraceEvent_AddPeer) GetPeerID() enode.ID {
	if m != nil {
		return m.PeerID
	}
	return enode.ZeroID
}

func (m *TraceEvent_AddPeer) GetProto() string {
	if m != nil && m.Proto != nil {
		return *m.Proto
	}
	return ""
}

type TraceEvent_RemovePeer struct {
	PeerID enode.ID `json:"peerID,omitempty"`
}

func (m *TraceEvent_RemovePeer) Reset()         { *m = TraceEvent_RemovePeer{} }
func (m *TraceEvent_RemovePeer) String() string { return "" }
func (*TraceEvent_RemovePeer) ProtoMessage()    {}

func (m *TraceEvent_RemovePeer) GetPeerID() enode.ID {
	if m != nil {
		return m.PeerID
	}
	return enode.ZeroID
}

type TraceEvent_RecvRPC struct {
	ReceivedFrom enode.ID            `json:"receivedFrom,omitempty"`
	Meta         *TraceEvent_RPCMeta `json:"meta,omitempty"`
}

func (m *TraceEvent_RecvRPC) Reset()         { *m = TraceEvent_RecvRPC{} }
func (m *TraceEvent_RecvRPC) String() string { return "" }
func (*TraceEvent_RecvRPC) ProtoMessage()    {}

func (m *TraceEvent_RecvRPC) GetReceivedFrom() enode.ID {
	if m != nil {
		return m.ReceivedFrom
	}
	return enode.ZeroID
}

func (m *TraceEvent_RecvRPC) GetMeta() *TraceEvent_RPCMeta {
	if m != nil {
		return m.Meta
	}
	return nil
}

type TraceEvent_SendRPC struct {
	SendTo enode.ID            `protobuf:"bytes,1,opt,name=sendTo" json:"sendTo,omitempty"`
	Meta   *TraceEvent_RPCMeta `protobuf:"bytes,2,opt,name=meta" json:"meta,omitempty"`
}

func (m *TraceEvent_SendRPC) Reset()         { *m = TraceEvent_SendRPC{} }
func (m *TraceEvent_SendRPC) String() string { return "" }
func (*TraceEvent_SendRPC) ProtoMessage()    {}

func (m *TraceEvent_SendRPC) GetSendTo() enode.ID {
	if m != nil {
		return m.SendTo
	}
	return enode.ZeroID
}

func (m *TraceEvent_SendRPC) GetMeta() *TraceEvent_RPCMeta {
	if m != nil {
		return m.Meta
	}
	return nil
}

type TraceEvent_DropRPC struct {
	SendTo enode.ID            `json:"sendTo,omitempty"`
	Meta   *TraceEvent_RPCMeta `json:"meta,omitempty"`
}

func (m *TraceEvent_DropRPC) Reset()         { *m = TraceEvent_DropRPC{} }
func (m *TraceEvent_DropRPC) String() string { return "" }
func (*TraceEvent_DropRPC) ProtoMessage()    {}

func (m *TraceEvent_DropRPC) GetSendTo() enode.ID {
	if m != nil {
		return m.SendTo
	}
	return enode.ZeroID
}

func (m *TraceEvent_DropRPC) GetMeta() *TraceEvent_RPCMeta {
	if m != nil {
		return m.Meta
	}
	return nil
}

type TraceEvent_Join struct {
	Topic *string `protobuf:"bytes,1,opt,name=topic" json:"topic,omitempty"`
}

func (m *TraceEvent_Join) Reset()         { *m = TraceEvent_Join{} }
func (m *TraceEvent_Join) String() string { return "" }
func (*TraceEvent_Join) ProtoMessage()    {}

func (m *TraceEvent_Join) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

type TraceEvent_Leave struct {
	Topic *string `protobuf:"bytes,2,opt,name=topic" json:"topic,omitempty"`
}

func (m *TraceEvent_Leave) Reset()         { *m = TraceEvent_Leave{} }
func (m *TraceEvent_Leave) String() string { return "" }
func (*TraceEvent_Leave) ProtoMessage()    {}

func (m *TraceEvent_Leave) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

type TraceEvent_Graft struct {
	PeerID enode.ID `protobuf:"bytes,1,opt,name=peerID" json:"peerID,omitempty"`
	Topic  *string  `protobuf:"bytes,2,opt,name=topic" json:"topic,omitempty"`
}

func (m *TraceEvent_Graft) Reset()         { *m = TraceEvent_Graft{} }
func (m *TraceEvent_Graft) String() string { return "" }

func (m *TraceEvent_Graft) GetPeerID() enode.ID {
	if m != nil {
		return m.PeerID
	}
	return enode.ZeroID
}

func (m *TraceEvent_Graft) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

type TraceEvent_Prune struct {
	PeerID enode.ID `json:"peerID,omitempty"`
	Topic  *string  `json:"topic,omitempty"`
}

func (m *TraceEvent_Prune) Reset()         { *m = TraceEvent_Prune{} }
func (m *TraceEvent_Prune) String() string { return "" }

func (m *TraceEvent_Prune) GetPeerID() enode.ID {
	if m != nil {
		return m.PeerID
	}
	return enode.ZeroID
}

func (m *TraceEvent_Prune) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

type TraceEvent_RPCMeta struct {
	Messages     []*TraceEvent_MessageMeta `protobuf:"bytes,1,rep,name=messages" json:"messages,omitempty"`
	Subscription []*TraceEvent_SubMeta     `protobuf:"bytes,2,rep,name=subscription" json:"subscription,omitempty"`
	Control      *TraceEvent_ControlMeta   `protobuf:"bytes,3,opt,name=control" json:"control,omitempty"`
}

func (m *TraceEvent_RPCMeta) Reset()         { *m = TraceEvent_RPCMeta{} }
func (m *TraceEvent_RPCMeta) String() string { return "" }

func (m *TraceEvent_RPCMeta) GetMessages() []*TraceEvent_MessageMeta {
	if m != nil {
		return m.Messages
	}
	return nil
}

func (m *TraceEvent_RPCMeta) GetSubscription() []*TraceEvent_SubMeta {
	if m != nil {
		return m.Subscription
	}
	return nil
}

func (m *TraceEvent_RPCMeta) GetControl() *TraceEvent_ControlMeta {
	if m != nil {
		return m.Control
	}
	return nil
}

type TraceEvent_MessageMeta struct {
	MessageID []byte  `protobuf:"bytes,1,opt,name=messageID" json:"messageID,omitempty"`
	Topic     *string `protobuf:"bytes,2,opt,name=topic" json:"topic,omitempty"`
}

func (m *TraceEvent_MessageMeta) Reset()         { *m = TraceEvent_MessageMeta{} }
func (m *TraceEvent_MessageMeta) String() string { return "" }

func (m *TraceEvent_MessageMeta) GetMessageID() []byte {
	if m != nil {
		return m.MessageID
	}
	return nil
}

func (m *TraceEvent_MessageMeta) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

type TraceEvent_SubMeta struct {
	Subscribe *bool   `protobuf:"varint,1,opt,name=subscribe" json:"subscribe,omitempty"`
	Topic     *string `protobuf:"bytes,2,opt,name=topic" json:"topic,omitempty"`
}

func (m *TraceEvent_SubMeta) Reset()         { *m = TraceEvent_SubMeta{} }
func (m *TraceEvent_SubMeta) String() string { return "" }
func (*TraceEvent_SubMeta) ProtoMessage()    {}

func (m *TraceEvent_SubMeta) GetSubscribe() bool {
	if m != nil && m.Subscribe != nil {
		return *m.Subscribe
	}
	return false
}

func (m *TraceEvent_SubMeta) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

type TraceEvent_ControlMeta struct {
	Ihave []*TraceEvent_ControlIHaveMeta `protobuf:"bytes,1,rep,name=ihave" json:"ihave,omitempty"`
	Iwant []*TraceEvent_ControlIWantMeta `protobuf:"bytes,2,rep,name=iwant" json:"iwant,omitempty"`
	Graft []*TraceEvent_ControlGraftMeta `protobuf:"bytes,3,rep,name=graft" json:"graft,omitempty"`
	Prune []*TraceEvent_ControlPruneMeta `protobuf:"bytes,4,rep,name=prune" json:"prune,omitempty"`
}

func (m *TraceEvent_ControlMeta) Reset()         { *m = TraceEvent_ControlMeta{} }
func (m *TraceEvent_ControlMeta) String() string { return "" }

func (m *TraceEvent_ControlMeta) GetIhave() []*TraceEvent_ControlIHaveMeta {
	if m != nil {
		return m.Ihave
	}
	return nil
}

func (m *TraceEvent_ControlMeta) GetIwant() []*TraceEvent_ControlIWantMeta {
	if m != nil {
		return m.Iwant
	}
	return nil
}

func (m *TraceEvent_ControlMeta) GetGraft() []*TraceEvent_ControlGraftMeta {
	if m != nil {
		return m.Graft
	}
	return nil
}

func (m *TraceEvent_ControlMeta) GetPrune() []*TraceEvent_ControlPruneMeta {
	if m != nil {
		return m.Prune
	}
	return nil
}

type TraceEvent_ControlIHaveMeta struct {
	Topic      *string  `protobuf:"bytes,1,opt,name=topic" json:"topic,omitempty"`
	MessageIDs [][]byte `protobuf:"bytes,2,rep,name=messageIDs" json:"messageIDs,omitempty"`
}

func (m *TraceEvent_ControlIHaveMeta) Reset()         { *m = TraceEvent_ControlIHaveMeta{} }
func (m *TraceEvent_ControlIHaveMeta) String() string { return "" }

func (m *TraceEvent_ControlIHaveMeta) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

func (m *TraceEvent_ControlIHaveMeta) GetMessageIDs() [][]byte {
	if m != nil {
		return m.MessageIDs
	}
	return nil
}

type TraceEvent_ControlIWantMeta struct {
	MessageIDs [][]byte `protobuf:"bytes,1,rep,name=messageIDs" json:"messageIDs,omitempty"`
}

func (m *TraceEvent_ControlIWantMeta) Reset()         { *m = TraceEvent_ControlIWantMeta{} }
func (m *TraceEvent_ControlIWantMeta) String() string { return "" }

func (m *TraceEvent_ControlIWantMeta) GetMessageIDs() [][]byte {
	if m != nil {
		return m.MessageIDs
	}
	return nil
}

type TraceEvent_ControlGraftMeta struct {
	Topic *string `protobuf:"bytes,1,opt,name=topic" json:"topic,omitempty"`
}

func (m *TraceEvent_ControlGraftMeta) Reset()         { *m = TraceEvent_ControlGraftMeta{} }
func (m *TraceEvent_ControlGraftMeta) String() string { return "" }

func (m *TraceEvent_ControlGraftMeta) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

type TraceEvent_ControlPruneMeta struct {
	Topic *string    `protobuf:"bytes,1,opt,name=topic" json:"topic,omitempty"`
	Peers []enode.ID `protobuf:"bytes,2,rep,name=peers" json:"peers,omitempty"`
}

func (m *TraceEvent_ControlPruneMeta) Reset()         { *m = TraceEvent_ControlPruneMeta{} }
func (m *TraceEvent_ControlPruneMeta) String() string { return "" }

func (m *TraceEvent_ControlPruneMeta) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

func (m *TraceEvent_ControlPruneMeta) GetPeers() []enode.ID {
	if m != nil {
		return m.Peers
	}
	return nil
}

type TraceEventBatch struct {
	Batch []*TraceEvent `protobuf:"bytes,1,rep,name=batch" json:"batch,omitempty"`
}

func (m *TraceEventBatch) Reset()         { *m = TraceEventBatch{} }
func (m *TraceEventBatch) String() string { return "" }

func (m *TraceEventBatch) GetBatch() []*TraceEvent {
	if m != nil {
		return m.Batch
	}
	return nil
}
