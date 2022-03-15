package message

import (
	"fmt"
	"io"
	"math"
	math_bits "math/bits"

	proto "github.com/gogo/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

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

func (TraceEvent_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 0}
}

type TraceEvent struct {
	Type                 *TraceEvent_Type             `protobuf:"varint,1,opt,name=type,enum=pubsub.pb.TraceEvent_Type" json:"type,omitempty"`
	PeerID               []byte                       `protobuf:"bytes,2,opt,name=peerID" json:"peerID,omitempty"`
	Timestamp            *int64                       `protobuf:"varint,3,opt,name=timestamp" json:"timestamp,omitempty"`
	PublishMessage       *TraceEvent_PublishMessage   `protobuf:"bytes,4,opt,name=publishMessage" json:"publishMessage,omitempty"`
	RejectMessage        *TraceEvent_RejectMessage    `protobuf:"bytes,5,opt,name=rejectMessage" json:"rejectMessage,omitempty"`
	DuplicateMessage     *TraceEvent_DuplicateMessage `protobuf:"bytes,6,opt,name=duplicateMessage" json:"duplicateMessage,omitempty"`
	DeliverMessage       *TraceEvent_DeliverMessage   `protobuf:"bytes,7,opt,name=deliverMessage" json:"deliverMessage,omitempty"`
	AddPeer              *TraceEvent_AddPeer          `protobuf:"bytes,8,opt,name=addPeer" json:"addPeer,omitempty"`
	RemovePeer           *TraceEvent_RemovePeer       `protobuf:"bytes,9,opt,name=removePeer" json:"removePeer,omitempty"`
	RecvRPC              *TraceEvent_RecvRPC          `protobuf:"bytes,10,opt,name=recvRPC" json:"recvRPC,omitempty"`
	SendRPC              *TraceEvent_SendRPC          `protobuf:"bytes,11,opt,name=sendRPC" json:"sendRPC,omitempty"`
	DropRPC              *TraceEvent_DropRPC          `protobuf:"bytes,12,opt,name=dropRPC" json:"dropRPC,omitempty"`
	Join                 *TraceEvent_Join             `protobuf:"bytes,13,opt,name=join" json:"join,omitempty"`
	Leave                *TraceEvent_Leave            `protobuf:"bytes,14,opt,name=leave" json:"leave,omitempty"`
	Graft                *TraceEvent_Graft            `protobuf:"bytes,15,opt,name=graft" json:"graft,omitempty"`
	Prune                *TraceEvent_Prune            `protobuf:"bytes,16,opt,name=prune" json:"prune,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                     `json:"-"`
	XXX_unrecognized     []byte                       `json:"-"`
	XXX_sizecache        int32                        `json:"-"`
}

func (m *TraceEvent) Reset()         { *m = TraceEvent{} }
func (m *TraceEvent) String() string { return proto.CompactTextString(m) }
func (*TraceEvent) ProtoMessage()    {}
func (*TraceEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0}
}
func (m *TraceEvent) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent.Merge(m, src)
}
func (m *TraceEvent) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent proto.InternalMessageInfo

func (m *TraceEvent) GetType() TraceEvent_Type {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return TraceEvent_PUBLISH_MESSAGE
}

func (m *TraceEvent) GetPeerID() []byte {
	if m != nil {
		return m.PeerID
	}
	return nil
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
	MessageID            []byte   `protobuf:"bytes,1,opt,name=messageID" json:"messageID,omitempty"`
	Topic                *string  `protobuf:"bytes,2,opt,name=topic" json:"topic,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TraceEvent_PublishMessage) Reset()         { *m = TraceEvent_PublishMessage{} }
func (m *TraceEvent_PublishMessage) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_PublishMessage) ProtoMessage()    {}
func (*TraceEvent_PublishMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 0}
}
func (m *TraceEvent_PublishMessage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_PublishMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_PublishMessage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_PublishMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_PublishMessage.Merge(m, src)
}
func (m *TraceEvent_PublishMessage) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_PublishMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_PublishMessage.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_PublishMessage proto.InternalMessageInfo

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
	MessageID            []byte   `protobuf:"bytes,1,opt,name=messageID" json:"messageID,omitempty"`
	ReceivedFrom         []byte   `protobuf:"bytes,2,opt,name=receivedFrom" json:"receivedFrom,omitempty"`
	Reason               *string  `protobuf:"bytes,3,opt,name=reason" json:"reason,omitempty"`
	Topic                *string  `protobuf:"bytes,4,opt,name=topic" json:"topic,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TraceEvent_RejectMessage) Reset()         { *m = TraceEvent_RejectMessage{} }
func (m *TraceEvent_RejectMessage) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_RejectMessage) ProtoMessage()    {}
func (*TraceEvent_RejectMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 1}
}
func (m *TraceEvent_RejectMessage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_RejectMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_RejectMessage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_RejectMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_RejectMessage.Merge(m, src)
}
func (m *TraceEvent_RejectMessage) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_RejectMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_RejectMessage.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_RejectMessage proto.InternalMessageInfo

func (m *TraceEvent_RejectMessage) GetMessageID() []byte {
	if m != nil {
		return m.MessageID
	}
	return nil
}

func (m *TraceEvent_RejectMessage) GetReceivedFrom() []byte {
	if m != nil {
		return m.ReceivedFrom
	}
	return nil
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
	MessageID            []byte   `protobuf:"bytes,1,opt,name=messageID" json:"messageID,omitempty"`
	ReceivedFrom         []byte   `protobuf:"bytes,2,opt,name=receivedFrom" json:"receivedFrom,omitempty"`
	Topic                *string  `protobuf:"bytes,3,opt,name=topic" json:"topic,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TraceEvent_DuplicateMessage) Reset()         { *m = TraceEvent_DuplicateMessage{} }
func (m *TraceEvent_DuplicateMessage) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_DuplicateMessage) ProtoMessage()    {}
func (*TraceEvent_DuplicateMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 2}
}
func (m *TraceEvent_DuplicateMessage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_DuplicateMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_DuplicateMessage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_DuplicateMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_DuplicateMessage.Merge(m, src)
}
func (m *TraceEvent_DuplicateMessage) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_DuplicateMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_DuplicateMessage.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_DuplicateMessage proto.InternalMessageInfo

func (m *TraceEvent_DuplicateMessage) GetMessageID() []byte {
	if m != nil {
		return m.MessageID
	}
	return nil
}

func (m *TraceEvent_DuplicateMessage) GetReceivedFrom() []byte {
	if m != nil {
		return m.ReceivedFrom
	}
	return nil
}

func (m *TraceEvent_DuplicateMessage) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

type TraceEvent_DeliverMessage struct {
	MessageID            []byte   `protobuf:"bytes,1,opt,name=messageID" json:"messageID,omitempty"`
	Topic                *string  `protobuf:"bytes,2,opt,name=topic" json:"topic,omitempty"`
	ReceivedFrom         []byte   `protobuf:"bytes,3,opt,name=receivedFrom" json:"receivedFrom,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TraceEvent_DeliverMessage) Reset()         { *m = TraceEvent_DeliverMessage{} }
func (m *TraceEvent_DeliverMessage) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_DeliverMessage) ProtoMessage()    {}
func (*TraceEvent_DeliverMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 3}
}
func (m *TraceEvent_DeliverMessage) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_DeliverMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_DeliverMessage.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_DeliverMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_DeliverMessage.Merge(m, src)
}
func (m *TraceEvent_DeliverMessage) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_DeliverMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_DeliverMessage.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_DeliverMessage proto.InternalMessageInfo

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

func (m *TraceEvent_DeliverMessage) GetReceivedFrom() []byte {
	if m != nil {
		return m.ReceivedFrom
	}
	return nil
}

type TraceEvent_AddPeer struct {
	PeerID               []byte   `protobuf:"bytes,1,opt,name=peerID" json:"peerID,omitempty"`
	Proto                *string  `protobuf:"bytes,2,opt,name=proto" json:"proto,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TraceEvent_AddPeer) Reset()         { *m = TraceEvent_AddPeer{} }
func (m *TraceEvent_AddPeer) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_AddPeer) ProtoMessage()    {}
func (*TraceEvent_AddPeer) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 4}
}
func (m *TraceEvent_AddPeer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_AddPeer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_AddPeer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_AddPeer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_AddPeer.Merge(m, src)
}
func (m *TraceEvent_AddPeer) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_AddPeer) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_AddPeer.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_AddPeer proto.InternalMessageInfo

func (m *TraceEvent_AddPeer) GetPeerID() []byte {
	if m != nil {
		return m.PeerID
	}
	return nil
}

func (m *TraceEvent_AddPeer) GetProto() string {
	if m != nil && m.Proto != nil {
		return *m.Proto
	}
	return ""
}

type TraceEvent_RemovePeer struct {
	PeerID               []byte   `protobuf:"bytes,1,opt,name=peerID" json:"peerID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TraceEvent_RemovePeer) Reset()         { *m = TraceEvent_RemovePeer{} }
func (m *TraceEvent_RemovePeer) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_RemovePeer) ProtoMessage()    {}
func (*TraceEvent_RemovePeer) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 5}
}
func (m *TraceEvent_RemovePeer) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_RemovePeer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_RemovePeer.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_RemovePeer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_RemovePeer.Merge(m, src)
}
func (m *TraceEvent_RemovePeer) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_RemovePeer) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_RemovePeer.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_RemovePeer proto.InternalMessageInfo

func (m *TraceEvent_RemovePeer) GetPeerID() []byte {
	if m != nil {
		return m.PeerID
	}
	return nil
}

type TraceEvent_RecvRPC struct {
	ReceivedFrom         []byte              `protobuf:"bytes,1,opt,name=receivedFrom" json:"receivedFrom,omitempty"`
	Meta                 *TraceEvent_RPCMeta `protobuf:"bytes,2,opt,name=meta" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *TraceEvent_RecvRPC) Reset()         { *m = TraceEvent_RecvRPC{} }
func (m *TraceEvent_RecvRPC) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_RecvRPC) ProtoMessage()    {}
func (*TraceEvent_RecvRPC) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 6}
}
func (m *TraceEvent_RecvRPC) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_RecvRPC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_RecvRPC.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_RecvRPC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_RecvRPC.Merge(m, src)
}
func (m *TraceEvent_RecvRPC) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_RecvRPC) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_RecvRPC.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_RecvRPC proto.InternalMessageInfo

func (m *TraceEvent_RecvRPC) GetReceivedFrom() []byte {
	if m != nil {
		return m.ReceivedFrom
	}
	return nil
}

func (m *TraceEvent_RecvRPC) GetMeta() *TraceEvent_RPCMeta {
	if m != nil {
		return m.Meta
	}
	return nil
}

type TraceEvent_SendRPC struct {
	SendTo               []byte              `protobuf:"bytes,1,opt,name=sendTo" json:"sendTo,omitempty"`
	Meta                 *TraceEvent_RPCMeta `protobuf:"bytes,2,opt,name=meta" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *TraceEvent_SendRPC) Reset()         { *m = TraceEvent_SendRPC{} }
func (m *TraceEvent_SendRPC) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_SendRPC) ProtoMessage()    {}
func (*TraceEvent_SendRPC) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 7}
}
func (m *TraceEvent_SendRPC) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_SendRPC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_SendRPC.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_SendRPC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_SendRPC.Merge(m, src)
}
func (m *TraceEvent_SendRPC) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_SendRPC) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_SendRPC.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_SendRPC proto.InternalMessageInfo

func (m *TraceEvent_SendRPC) GetSendTo() []byte {
	if m != nil {
		return m.SendTo
	}
	return nil
}

func (m *TraceEvent_SendRPC) GetMeta() *TraceEvent_RPCMeta {
	if m != nil {
		return m.Meta
	}
	return nil
}

type TraceEvent_DropRPC struct {
	SendTo               []byte              `protobuf:"bytes,1,opt,name=sendTo" json:"sendTo,omitempty"`
	Meta                 *TraceEvent_RPCMeta `protobuf:"bytes,2,opt,name=meta" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *TraceEvent_DropRPC) Reset()         { *m = TraceEvent_DropRPC{} }
func (m *TraceEvent_DropRPC) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_DropRPC) ProtoMessage()    {}
func (*TraceEvent_DropRPC) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 8}
}
func (m *TraceEvent_DropRPC) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_DropRPC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_DropRPC.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_DropRPC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_DropRPC.Merge(m, src)
}
func (m *TraceEvent_DropRPC) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_DropRPC) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_DropRPC.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_DropRPC proto.InternalMessageInfo

func (m *TraceEvent_DropRPC) GetSendTo() []byte {
	if m != nil {
		return m.SendTo
	}
	return nil
}

func (m *TraceEvent_DropRPC) GetMeta() *TraceEvent_RPCMeta {
	if m != nil {
		return m.Meta
	}
	return nil
}

type TraceEvent_Join struct {
	Topic                *string  `protobuf:"bytes,1,opt,name=topic" json:"topic,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TraceEvent_Join) Reset()         { *m = TraceEvent_Join{} }
func (m *TraceEvent_Join) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_Join) ProtoMessage()    {}
func (*TraceEvent_Join) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 9}
}
func (m *TraceEvent_Join) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_Join) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_Join.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_Join) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_Join.Merge(m, src)
}
func (m *TraceEvent_Join) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_Join) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_Join.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_Join proto.InternalMessageInfo

func (m *TraceEvent_Join) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

type TraceEvent_Leave struct {
	Topic                *string  `protobuf:"bytes,2,opt,name=topic" json:"topic,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TraceEvent_Leave) Reset()         { *m = TraceEvent_Leave{} }
func (m *TraceEvent_Leave) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_Leave) ProtoMessage()    {}
func (*TraceEvent_Leave) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 10}
}
func (m *TraceEvent_Leave) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_Leave) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_Leave.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_Leave) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_Leave.Merge(m, src)
}
func (m *TraceEvent_Leave) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_Leave) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_Leave.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_Leave proto.InternalMessageInfo

func (m *TraceEvent_Leave) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

type TraceEvent_Graft struct {
	PeerID               []byte   `protobuf:"bytes,1,opt,name=peerID" json:"peerID,omitempty"`
	Topic                *string  `protobuf:"bytes,2,opt,name=topic" json:"topic,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TraceEvent_Graft) Reset()         { *m = TraceEvent_Graft{} }
func (m *TraceEvent_Graft) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_Graft) ProtoMessage()    {}
func (*TraceEvent_Graft) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 11}
}
func (m *TraceEvent_Graft) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_Graft) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_Graft.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_Graft) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_Graft.Merge(m, src)
}
func (m *TraceEvent_Graft) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_Graft) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_Graft.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_Graft proto.InternalMessageInfo

func (m *TraceEvent_Graft) GetPeerID() []byte {
	if m != nil {
		return m.PeerID
	}
	return nil
}

func (m *TraceEvent_Graft) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

type TraceEvent_Prune struct {
	PeerID               []byte   `protobuf:"bytes,1,opt,name=peerID" json:"peerID,omitempty"`
	Topic                *string  `protobuf:"bytes,2,opt,name=topic" json:"topic,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TraceEvent_Prune) Reset()         { *m = TraceEvent_Prune{} }
func (m *TraceEvent_Prune) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_Prune) ProtoMessage()    {}
func (*TraceEvent_Prune) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 12}
}
func (m *TraceEvent_Prune) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_Prune) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_Prune.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_Prune) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_Prune.Merge(m, src)
}
func (m *TraceEvent_Prune) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_Prune) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_Prune.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_Prune proto.InternalMessageInfo

func (m *TraceEvent_Prune) GetPeerID() []byte {
	if m != nil {
		return m.PeerID
	}
	return nil
}

func (m *TraceEvent_Prune) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

type TraceEvent_RPCMeta struct {
	Messages             []*TraceEvent_MessageMeta `protobuf:"bytes,1,rep,name=messages" json:"messages,omitempty"`
	Subscription         []*TraceEvent_SubMeta     `protobuf:"bytes,2,rep,name=subscription" json:"subscription,omitempty"`
	Control              *TraceEvent_ControlMeta   `protobuf:"bytes,3,opt,name=control" json:"control,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *TraceEvent_RPCMeta) Reset()         { *m = TraceEvent_RPCMeta{} }
func (m *TraceEvent_RPCMeta) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_RPCMeta) ProtoMessage()    {}
func (*TraceEvent_RPCMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 13}
}
func (m *TraceEvent_RPCMeta) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_RPCMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_RPCMeta.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_RPCMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_RPCMeta.Merge(m, src)
}
func (m *TraceEvent_RPCMeta) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_RPCMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_RPCMeta.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_RPCMeta proto.InternalMessageInfo

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
	MessageID            []byte   `protobuf:"bytes,1,opt,name=messageID" json:"messageID,omitempty"`
	Topic                *string  `protobuf:"bytes,2,opt,name=topic" json:"topic,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TraceEvent_MessageMeta) Reset()         { *m = TraceEvent_MessageMeta{} }
func (m *TraceEvent_MessageMeta) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_MessageMeta) ProtoMessage()    {}
func (*TraceEvent_MessageMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 14}
}
func (m *TraceEvent_MessageMeta) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_MessageMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_MessageMeta.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_MessageMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_MessageMeta.Merge(m, src)
}
func (m *TraceEvent_MessageMeta) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_MessageMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_MessageMeta.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_MessageMeta proto.InternalMessageInfo

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
	Subscribe            *bool    `protobuf:"varint,1,opt,name=subscribe" json:"subscribe,omitempty"`
	Topic                *string  `protobuf:"bytes,2,opt,name=topic" json:"topic,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TraceEvent_SubMeta) Reset()         { *m = TraceEvent_SubMeta{} }
func (m *TraceEvent_SubMeta) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_SubMeta) ProtoMessage()    {}
func (*TraceEvent_SubMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 15}
}
func (m *TraceEvent_SubMeta) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_SubMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_SubMeta.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_SubMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_SubMeta.Merge(m, src)
}
func (m *TraceEvent_SubMeta) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_SubMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_SubMeta.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_SubMeta proto.InternalMessageInfo

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
	Ihave                []*TraceEvent_ControlIHaveMeta `protobuf:"bytes,1,rep,name=ihave" json:"ihave,omitempty"`
	Iwant                []*TraceEvent_ControlIWantMeta `protobuf:"bytes,2,rep,name=iwant" json:"iwant,omitempty"`
	Graft                []*TraceEvent_ControlGraftMeta `protobuf:"bytes,3,rep,name=graft" json:"graft,omitempty"`
	Prune                []*TraceEvent_ControlPruneMeta `protobuf:"bytes,4,rep,name=prune" json:"prune,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                       `json:"-"`
	XXX_unrecognized     []byte                         `json:"-"`
	XXX_sizecache        int32                          `json:"-"`
}

func (m *TraceEvent_ControlMeta) Reset()         { *m = TraceEvent_ControlMeta{} }
func (m *TraceEvent_ControlMeta) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_ControlMeta) ProtoMessage()    {}
func (*TraceEvent_ControlMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 16}
}
func (m *TraceEvent_ControlMeta) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_ControlMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_ControlMeta.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_ControlMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_ControlMeta.Merge(m, src)
}
func (m *TraceEvent_ControlMeta) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_ControlMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_ControlMeta.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_ControlMeta proto.InternalMessageInfo

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
	Topic                *string  `protobuf:"bytes,1,opt,name=topic" json:"topic,omitempty"`
	MessageIDs           [][]byte `protobuf:"bytes,2,rep,name=messageIDs" json:"messageIDs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TraceEvent_ControlIHaveMeta) Reset()         { *m = TraceEvent_ControlIHaveMeta{} }
func (m *TraceEvent_ControlIHaveMeta) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_ControlIHaveMeta) ProtoMessage()    {}
func (*TraceEvent_ControlIHaveMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 17}
}
func (m *TraceEvent_ControlIHaveMeta) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_ControlIHaveMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_ControlIHaveMeta.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_ControlIHaveMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_ControlIHaveMeta.Merge(m, src)
}
func (m *TraceEvent_ControlIHaveMeta) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_ControlIHaveMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_ControlIHaveMeta.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_ControlIHaveMeta proto.InternalMessageInfo

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
	MessageIDs           [][]byte `protobuf:"bytes,1,rep,name=messageIDs" json:"messageIDs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TraceEvent_ControlIWantMeta) Reset()         { *m = TraceEvent_ControlIWantMeta{} }
func (m *TraceEvent_ControlIWantMeta) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_ControlIWantMeta) ProtoMessage()    {}
func (*TraceEvent_ControlIWantMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 18}
}
func (m *TraceEvent_ControlIWantMeta) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_ControlIWantMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_ControlIWantMeta.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_ControlIWantMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_ControlIWantMeta.Merge(m, src)
}
func (m *TraceEvent_ControlIWantMeta) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_ControlIWantMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_ControlIWantMeta.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_ControlIWantMeta proto.InternalMessageInfo

func (m *TraceEvent_ControlIWantMeta) GetMessageIDs() [][]byte {
	if m != nil {
		return m.MessageIDs
	}
	return nil
}

type TraceEvent_ControlGraftMeta struct {
	Topic                *string  `protobuf:"bytes,1,opt,name=topic" json:"topic,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TraceEvent_ControlGraftMeta) Reset()         { *m = TraceEvent_ControlGraftMeta{} }
func (m *TraceEvent_ControlGraftMeta) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_ControlGraftMeta) ProtoMessage()    {}
func (*TraceEvent_ControlGraftMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 19}
}
func (m *TraceEvent_ControlGraftMeta) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_ControlGraftMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_ControlGraftMeta.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_ControlGraftMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_ControlGraftMeta.Merge(m, src)
}
func (m *TraceEvent_ControlGraftMeta) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_ControlGraftMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_ControlGraftMeta.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_ControlGraftMeta proto.InternalMessageInfo

func (m *TraceEvent_ControlGraftMeta) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

type TraceEvent_ControlPruneMeta struct {
	Topic                *string  `protobuf:"bytes,1,opt,name=topic" json:"topic,omitempty"`
	Peers                [][]byte `protobuf:"bytes,2,rep,name=peers" json:"peers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TraceEvent_ControlPruneMeta) Reset()         { *m = TraceEvent_ControlPruneMeta{} }
func (m *TraceEvent_ControlPruneMeta) String() string { return proto.CompactTextString(m) }
func (*TraceEvent_ControlPruneMeta) ProtoMessage()    {}
func (*TraceEvent_ControlPruneMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{0, 20}
}
func (m *TraceEvent_ControlPruneMeta) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEvent_ControlPruneMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEvent_ControlPruneMeta.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEvent_ControlPruneMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEvent_ControlPruneMeta.Merge(m, src)
}
func (m *TraceEvent_ControlPruneMeta) XXX_Size() int {
	return m.Size()
}
func (m *TraceEvent_ControlPruneMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEvent_ControlPruneMeta.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEvent_ControlPruneMeta proto.InternalMessageInfo

func (m *TraceEvent_ControlPruneMeta) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

func (m *TraceEvent_ControlPruneMeta) GetPeers() [][]byte {
	if m != nil {
		return m.Peers
	}
	return nil
}

type TraceEventBatch struct {
	Batch                []*TraceEvent `protobuf:"bytes,1,rep,name=batch" json:"batch,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *TraceEventBatch) Reset()         { *m = TraceEventBatch{} }
func (m *TraceEventBatch) String() string { return proto.CompactTextString(m) }
func (*TraceEventBatch) ProtoMessage()    {}
func (*TraceEventBatch) Descriptor() ([]byte, []int) {
	return fileDescriptor_0571941a1d628a80, []int{1}
}
func (m *TraceEventBatch) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TraceEventBatch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TraceEventBatch.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TraceEventBatch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TraceEventBatch.Merge(m, src)
}
func (m *TraceEventBatch) XXX_Size() int {
	return m.Size()
}
func (m *TraceEventBatch) XXX_DiscardUnknown() {
	xxx_messageInfo_TraceEventBatch.DiscardUnknown(m)
}

var xxx_messageInfo_TraceEventBatch proto.InternalMessageInfo

func (m *TraceEventBatch) GetBatch() []*TraceEvent {
	if m != nil {
		return m.Batch
	}
	return nil
}

func init() {
	proto.RegisterEnum("pubsub.pb.TraceEvent_Type", TraceEvent_Type_name, TraceEvent_Type_value)
	proto.RegisterType((*TraceEvent)(nil), "pubsub.pb.TraceEvent")
	proto.RegisterType((*TraceEvent_PublishMessage)(nil), "pubsub.pb.TraceEvent.PublishMessage")
	proto.RegisterType((*TraceEvent_RejectMessage)(nil), "pubsub.pb.TraceEvent.RejectMessage")
	proto.RegisterType((*TraceEvent_DuplicateMessage)(nil), "pubsub.pb.TraceEvent.DuplicateMessage")
	proto.RegisterType((*TraceEvent_DeliverMessage)(nil), "pubsub.pb.TraceEvent.DeliverMessage")
	proto.RegisterType((*TraceEvent_AddPeer)(nil), "pubsub.pb.TraceEvent.AddPeer")
	proto.RegisterType((*TraceEvent_RemovePeer)(nil), "pubsub.pb.TraceEvent.RemovePeer")
	proto.RegisterType((*TraceEvent_RecvRPC)(nil), "pubsub.pb.TraceEvent.RecvRPC")
	proto.RegisterType((*TraceEvent_SendRPC)(nil), "pubsub.pb.TraceEvent.SendRPC")
	proto.RegisterType((*TraceEvent_DropRPC)(nil), "pubsub.pb.TraceEvent.DropRPC")
	proto.RegisterType((*TraceEvent_Join)(nil), "pubsub.pb.TraceEvent.Join")
	proto.RegisterType((*TraceEvent_Leave)(nil), "pubsub.pb.TraceEvent.Leave")
	proto.RegisterType((*TraceEvent_Graft)(nil), "pubsub.pb.TraceEvent.Graft")
	proto.RegisterType((*TraceEvent_Prune)(nil), "pubsub.pb.TraceEvent.Prune")
	proto.RegisterType((*TraceEvent_RPCMeta)(nil), "pubsub.pb.TraceEvent.RPCMeta")
	proto.RegisterType((*TraceEvent_MessageMeta)(nil), "pubsub.pb.TraceEvent.MessageMeta")
	proto.RegisterType((*TraceEvent_SubMeta)(nil), "pubsub.pb.TraceEvent.SubMeta")
	proto.RegisterType((*TraceEvent_ControlMeta)(nil), "pubsub.pb.TraceEvent.ControlMeta")
	proto.RegisterType((*TraceEvent_ControlIHaveMeta)(nil), "pubsub.pb.TraceEvent.ControlIHaveMeta")
	proto.RegisterType((*TraceEvent_ControlIWantMeta)(nil), "pubsub.pb.TraceEvent.ControlIWantMeta")
	proto.RegisterType((*TraceEvent_ControlGraftMeta)(nil), "pubsub.pb.TraceEvent.ControlGraftMeta")
	proto.RegisterType((*TraceEvent_ControlPruneMeta)(nil), "pubsub.pb.TraceEvent.ControlPruneMeta")
	proto.RegisterType((*TraceEventBatch)(nil), "pubsub.pb.TraceEventBatch")
}

func init() { proto.RegisterFile("trace.proto", fileDescriptor_0571941a1d628a80) }

var fileDescriptor_0571941a1d628a80 = []byte{
	// 999 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x96, 0x51, 0x6f, 0xda, 0x56,
	0x14, 0xc7, 0xe7, 0x00, 0x01, 0x0e, 0x84, 0x78, 0x77, 0x6d, 0x65, 0xb1, 0x36, 0x62, 0x59, 0x55,
	0x21, 0x4d, 0x42, 0x6a, 0xa4, 0xa9, 0x0f, 0x6b, 0xab, 0x11, 0xec, 0x26, 0x44, 0x24, 0xb1, 0x0e,
	0x24, 0x7b, 0xcc, 0x0c, 0xdc, 0x35, 0x8e, 0xc0, 0xb6, 0xec, 0x0b, 0x53, 0x9f, 0xf6, 0xb4, 0xef,
	0xd6, 0xb7, 0xed, 0x23, 0x54, 0xf9, 0x24, 0xd3, 0xbd, 0xd7, 0x36, 0x36, 0xd8, 0xb4, 0x8b, 0xfa,
	0xe6, 0x73, 0xf3, 0xff, 0x9d, 0x7b, 0xce, 0xbd, 0xe7, 0x7f, 0x03, 0xd4, 0x98, 0x6f, 0x4d, 0x68,
	0xc7, 0xf3, 0x5d, 0xe6, 0x92, 0xaa, 0xb7, 0x18, 0x07, 0x8b, 0x71, 0xc7, 0x1b, 0x1f, 0x7e, 0x7a,
	0x02, 0x30, 0xe2, 0x7f, 0x32, 0x96, 0xd4, 0x61, 0xa4, 0x03, 0x45, 0xf6, 0xc1, 0xa3, 0x9a, 0xd2,
	0x52, 0xda, 0x8d, 0xa3, 0x66, 0x27, 0x16, 0x76, 0x56, 0xa2, 0xce, 0xe8, 0x83, 0x47, 0x51, 0xe8,
	0xc8, 0x13, 0xd8, 0xf5, 0x28, 0xf5, 0xfb, 0xba, 0xb6, 0xd3, 0x52, 0xda, 0x75, 0x0c, 0x23, 0xf2,
	0x14, 0xaa, 0xcc, 0x9e, 0xd3, 0x80, 0x59, 0x73, 0x4f, 0x2b, 0xb4, 0x94, 0x76, 0x01, 0x57, 0x0b,
	0x64, 0x00, 0x0d, 0x6f, 0x31, 0x9e, 0xd9, 0xc1, 0xed, 0x39, 0x0d, 0x02, 0xeb, 0x3d, 0xd5, 0x8a,
	0x2d, 0xa5, 0x5d, 0x3b, 0x7a, 0x9e, 0xbd, 0x9f, 0x99, 0xd2, 0xe2, 0x1a, 0x4b, 0xfa, 0xb0, 0xe7,
	0xd3, 0x3b, 0x3a, 0x61, 0x51, 0xb2, 0x92, 0x48, 0xf6, 0x63, 0x76, 0x32, 0x4c, 0x4a, 0x31, 0x4d,
	0x12, 0x04, 0x75, 0xba, 0xf0, 0x66, 0xf6, 0xc4, 0x62, 0x34, 0xca, 0xb6, 0x2b, 0xb2, 0xbd, 0xc8,
	0xce, 0xa6, 0xaf, 0xa9, 0x71, 0x83, 0xe7, 0xcd, 0x4e, 0xe9, 0xcc, 0x5e, 0x52, 0x3f, 0xca, 0x58,
	0xde, 0xd6, 0xac, 0x9e, 0xd2, 0xe2, 0x1a, 0x4b, 0x5e, 0x41, 0xd9, 0x9a, 0x4e, 0x4d, 0x4a, 0x7d,
	0xad, 0x22, 0xd2, 0x3c, 0xcb, 0x4e, 0xd3, 0x95, 0x22, 0x8c, 0xd4, 0xe4, 0x57, 0x00, 0x9f, 0xce,
	0xdd, 0x25, 0x15, 0x6c, 0x55, 0xb0, 0xad, 0xbc, 0x23, 0x8a, 0x74, 0x98, 0x60, 0xf8, 0xd6, 0x3e,
	0x9d, 0x2c, 0xd1, 0xec, 0x69, 0xb0, 0x6d, 0x6b, 0x94, 0x22, 0x8c, 0xd4, 0x1c, 0x0c, 0xa8, 0x33,
	0xe5, 0x60, 0x6d, 0x1b, 0x38, 0x94, 0x22, 0x8c, 0xd4, 0x1c, 0x9c, 0xfa, 0xae, 0xc7, 0xc1, 0xfa,
	0x36, 0x50, 0x97, 0x22, 0x8c, 0xd4, 0x7c, 0x8c, 0xef, 0x5c, 0xdb, 0xd1, 0xf6, 0x04, 0x95, 0x33,
	0xc6, 0x67, 0xae, 0xed, 0xa0, 0xd0, 0x91, 0x97, 0x50, 0x9a, 0x51, 0x6b, 0x49, 0xb5, 0x86, 0x00,
	0xbe, 0xcf, 0x06, 0x06, 0x5c, 0x82, 0x52, 0xc9, 0x91, 0xf7, 0xbe, 0xf5, 0x07, 0xd3, 0xf6, 0xb7,
	0x21, 0x27, 0x5c, 0x82, 0x52, 0xc9, 0x11, 0xcf, 0x5f, 0x38, 0x54, 0x53, 0xb7, 0x21, 0x26, 0x97,
	0xa0, 0x54, 0x36, 0x75, 0x68, 0xa4, 0xa7, 0x9f, 0x3b, 0x6b, 0x2e, 0x3f, 0xfb, 0xba, 0xb0, 0x69,
	0x1d, 0x57, 0x0b, 0xe4, 0x11, 0x94, 0x98, 0xeb, 0xd9, 0x13, 0x61, 0xc7, 0x2a, 0xca, 0xa0, 0xf9,
	0x17, 0xec, 0xa5, 0xc6, 0xfe, 0x33, 0x49, 0x0e, 0xa1, 0xee, 0xd3, 0x09, 0xb5, 0x97, 0x74, 0xfa,
	0xce, 0x77, 0xe7, 0xa1, 0xb5, 0x53, 0x6b, 0xdc, 0xf8, 0x3e, 0xb5, 0x02, 0xd7, 0x11, 0xee, 0xae,
	0x62, 0x18, 0xad, 0x0a, 0x28, 0x26, 0x0b, 0xb8, 0x03, 0x75, 0xdd, 0x29, 0x5f, 0xa1, 0x86, 0x78,
	0xaf, 0x42, 0x72, 0xaf, 0x5b, 0x68, 0xa4, 0x3d, 0xf4, 0x90, 0x23, 0xdb, 0xd8, 0xbf, 0xb0, 0xb9,
	0x7f, 0xf3, 0x15, 0x94, 0x43, 0x9b, 0x25, 0xde, 0x41, 0x25, 0xf5, 0x0e, 0x3e, 0xe2, 0x57, 0xee,
	0x32, 0x37, 0x4a, 0x2e, 0x82, 0xe6, 0x73, 0x80, 0x95, 0xc7, 0xf2, 0xd8, 0xe6, 0xef, 0x50, 0x0e,
	0xad, 0xb4, 0x51, 0x8d, 0x92, 0x71, 0x1a, 0x2f, 0xa1, 0x38, 0xa7, 0xcc, 0x12, 0x3b, 0xe5, 0x7b,
	0xd3, 0xec, 0x9d, 0x53, 0x66, 0xa1, 0x90, 0x36, 0x47, 0x50, 0x0e, 0x3d, 0xc7, 0x8b, 0xe0, 0xae,
	0x1b, 0xb9, 0x51, 0x11, 0x32, 0x7a, 0x60, 0xd6, 0xd0, 0x90, 0x5f, 0x33, 0xeb, 0x53, 0x28, 0x72,
	0xc3, 0xae, 0xae, 0x4b, 0x49, 0x5e, 0xfa, 0x33, 0x28, 0x09, 0x77, 0xe6, 0x18, 0xe0, 0x67, 0x28,
	0x09, 0x27, 0x6e, 0xbb, 0xa7, 0x6c, 0x4c, 0xb8, 0xf1, 0x7f, 0x62, 0x1f, 0x15, 0x28, 0x87, 0xc5,
	0x93, 0x37, 0x50, 0x09, 0x47, 0x2d, 0xd0, 0x94, 0x56, 0xa1, 0x5d, 0x3b, 0xfa, 0x21, 0xbb, 0xdb,
	0x70, 0x58, 0x45, 0xc7, 0x31, 0x42, 0xba, 0x50, 0x0f, 0x16, 0xe3, 0x60, 0xe2, 0xdb, 0x1e, 0xb3,
	0x5d, 0x47, 0xdb, 0x11, 0x29, 0xf2, 0xde, 0xcf, 0xc5, 0x58, 0xe0, 0x29, 0x84, 0xfc, 0x02, 0xe5,
	0x89, 0xeb, 0x30, 0xdf, 0x9d, 0x89, 0x21, 0xce, 0x2d, 0xa0, 0x27, 0x45, 0x22, 0x43, 0x44, 0x34,
	0xbb, 0x50, 0x4b, 0x14, 0xf6, 0xa0, 0xc7, 0xe7, 0x0d, 0x94, 0xc3, 0xc2, 0x38, 0x1e, 0x96, 0x36,
	0x96, 0x3f, 0x31, 0x2a, 0xb8, 0x5a, 0xc8, 0xc1, 0xff, 0xde, 0x81, 0x5a, 0xa2, 0x34, 0xf2, 0x1a,
	0x4a, 0xf6, 0x2d, 0x7f, 0xaa, 0xe5, 0x69, 0xbe, 0xd8, 0xda, 0x4c, 0xff, 0xd4, 0x5a, 0xca, 0x23,
	0x95, 0x90, 0xa0, 0xff, 0xb4, 0x1c, 0x16, 0x1e, 0xe4, 0x67, 0xe8, 0xdf, 0x2c, 0x87, 0x85, 0x34,
	0x87, 0x38, 0x2d, 0xdf, 0xfc, 0xc2, 0x17, 0xd0, 0x62, 0xe0, 0x24, 0x2d, 0x9f, 0xff, 0xd7, 0xd1,
	0xf3, 0x5f, 0xfc, 0x02, 0x5a, 0xcc, 0x9d, 0xa4, 0xe5, 0x7f, 0x82, 0x53, 0x50, 0xd7, 0x9b, 0xca,
	0xf6, 0x02, 0x39, 0x00, 0x88, 0xef, 0x24, 0x10, 0x8d, 0xd6, 0x31, 0xb1, 0xd2, 0x3c, 0x5a, 0x65,
	0x8a, 0x1a, 0x5c, 0x63, 0x94, 0x0d, 0xa6, 0x1d, 0x33, 0x71, 0x5b, 0x39, 0x4e, 0x7c, 0x1b, 0x2b,
	0xe3, 0x16, 0x72, 0xea, 0xe4, 0x6f, 0x23, 0xa5, 0x7e, 0x54, 0xa2, 0x0c, 0x0e, 0xff, 0x51, 0xa0,
	0xc8, 0x7f, 0x60, 0x92, 0xef, 0x60, 0xdf, 0xbc, 0x3a, 0x1e, 0xf4, 0x87, 0xa7, 0x37, 0xe7, 0xc6,
	0x70, 0xd8, 0x3d, 0x31, 0xd4, 0x6f, 0x08, 0x81, 0x06, 0x1a, 0x67, 0x46, 0x6f, 0x14, 0xaf, 0x29,
	0xe4, 0x31, 0x7c, 0xab, 0x5f, 0x99, 0x83, 0x7e, 0xaf, 0x3b, 0x32, 0xe2, 0xe5, 0x1d, 0xce, 0xeb,
	0xc6, 0xa0, 0x7f, 0x6d, 0x60, 0xbc, 0x58, 0x20, 0x75, 0xa8, 0x74, 0x75, 0xfd, 0xc6, 0x34, 0x0c,
	0x54, 0x8b, 0x64, 0x1f, 0x6a, 0x68, 0x9c, 0x5f, 0x5e, 0x1b, 0x72, 0xa1, 0xc4, 0xff, 0x8c, 0x46,
	0xef, 0xfa, 0x06, 0xcd, 0x9e, 0xba, 0xcb, 0xa3, 0xa1, 0x71, 0xa1, 0x8b, 0xa8, 0xcc, 0x23, 0x1d,
	0x2f, 0x4d, 0x11, 0x55, 0x48, 0x05, 0x8a, 0x67, 0x97, 0xfd, 0x0b, 0xb5, 0x4a, 0xaa, 0x50, 0x1a,
	0x18, 0xdd, 0x6b, 0x43, 0x05, 0xfe, 0x79, 0x82, 0xdd, 0x77, 0x23, 0xb5, 0xc6, 0x3f, 0x4d, 0xbc,
	0xba, 0x30, 0xd4, 0xfa, 0xe1, 0x5b, 0xd8, 0x5f, 0xdd, 0xef, 0xb1, 0xc5, 0x26, 0xb7, 0xe4, 0x27,
	0x28, 0x8d, 0xf9, 0x47, 0x38, 0xc4, 0x8f, 0x33, 0x47, 0x01, 0xa5, 0xe6, 0xb8, 0xfe, 0xf1, 0xfe,
	0x40, 0xf9, 0xf7, 0xfe, 0x40, 0xf9, 0x74, 0x7f, 0xa0, 0xfc, 0x17, 0x00, 0x00, 0xff, 0xff, 0xdb,
	0x3a, 0x1c, 0xe4, 0xc9, 0x0b, 0x00, 0x00,
}

func (m *TraceEvent) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Prune != nil {
		{
			size, err := m.Prune.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1
		i--
		dAtA[i] = 0x82
	}
	if m.Graft != nil {
		{
			size, err := m.Graft.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x7a
	}
	if m.Leave != nil {
		{
			size, err := m.Leave.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x72
	}
	if m.Join != nil {
		{
			size, err := m.Join.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x6a
	}
	if m.DropRPC != nil {
		{
			size, err := m.DropRPC.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x62
	}
	if m.SendRPC != nil {
		{
			size, err := m.SendRPC.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x5a
	}
	if m.RecvRPC != nil {
		{
			size, err := m.RecvRPC.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x52
	}
	if m.RemovePeer != nil {
		{
			size, err := m.RemovePeer.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x4a
	}
	if m.AddPeer != nil {
		{
			size, err := m.AddPeer.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x42
	}
	if m.DeliverMessage != nil {
		{
			size, err := m.DeliverMessage.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x3a
	}
	if m.DuplicateMessage != nil {
		{
			size, err := m.DuplicateMessage.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x32
	}
	if m.RejectMessage != nil {
		{
			size, err := m.RejectMessage.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.PublishMessage != nil {
		{
			size, err := m.PublishMessage.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if m.Timestamp != nil {
		i = encodeVarintTrace(dAtA, i, uint64(*m.Timestamp))
		i--
		dAtA[i] = 0x18
	}
	if m.PeerID != nil {
		i -= len(m.PeerID)
		copy(dAtA[i:], m.PeerID)
		i = encodeVarintTrace(dAtA, i, uint64(len(m.PeerID)))
		i--
		dAtA[i] = 0x12
	}
	if m.Type != nil {
		i = encodeVarintTrace(dAtA, i, uint64(*m.Type))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_PublishMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_PublishMessage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_PublishMessage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Topic != nil {
		i -= len(*m.Topic)
		copy(dAtA[i:], *m.Topic)
		i = encodeVarintTrace(dAtA, i, uint64(len(*m.Topic)))
		i--
		dAtA[i] = 0x12
	}
	if m.MessageID != nil {
		i -= len(m.MessageID)
		copy(dAtA[i:], m.MessageID)
		i = encodeVarintTrace(dAtA, i, uint64(len(m.MessageID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_RejectMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_RejectMessage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_RejectMessage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Topic != nil {
		i -= len(*m.Topic)
		copy(dAtA[i:], *m.Topic)
		i = encodeVarintTrace(dAtA, i, uint64(len(*m.Topic)))
		i--
		dAtA[i] = 0x22
	}
	if m.Reason != nil {
		i -= len(*m.Reason)
		copy(dAtA[i:], *m.Reason)
		i = encodeVarintTrace(dAtA, i, uint64(len(*m.Reason)))
		i--
		dAtA[i] = 0x1a
	}
	if m.ReceivedFrom != nil {
		i -= len(m.ReceivedFrom)
		copy(dAtA[i:], m.ReceivedFrom)
		i = encodeVarintTrace(dAtA, i, uint64(len(m.ReceivedFrom)))
		i--
		dAtA[i] = 0x12
	}
	if m.MessageID != nil {
		i -= len(m.MessageID)
		copy(dAtA[i:], m.MessageID)
		i = encodeVarintTrace(dAtA, i, uint64(len(m.MessageID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_DuplicateMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_DuplicateMessage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_DuplicateMessage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Topic != nil {
		i -= len(*m.Topic)
		copy(dAtA[i:], *m.Topic)
		i = encodeVarintTrace(dAtA, i, uint64(len(*m.Topic)))
		i--
		dAtA[i] = 0x1a
	}
	if m.ReceivedFrom != nil {
		i -= len(m.ReceivedFrom)
		copy(dAtA[i:], m.ReceivedFrom)
		i = encodeVarintTrace(dAtA, i, uint64(len(m.ReceivedFrom)))
		i--
		dAtA[i] = 0x12
	}
	if m.MessageID != nil {
		i -= len(m.MessageID)
		copy(dAtA[i:], m.MessageID)
		i = encodeVarintTrace(dAtA, i, uint64(len(m.MessageID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_DeliverMessage) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_DeliverMessage) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_DeliverMessage) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.ReceivedFrom != nil {
		i -= len(m.ReceivedFrom)
		copy(dAtA[i:], m.ReceivedFrom)
		i = encodeVarintTrace(dAtA, i, uint64(len(m.ReceivedFrom)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Topic != nil {
		i -= len(*m.Topic)
		copy(dAtA[i:], *m.Topic)
		i = encodeVarintTrace(dAtA, i, uint64(len(*m.Topic)))
		i--
		dAtA[i] = 0x12
	}
	if m.MessageID != nil {
		i -= len(m.MessageID)
		copy(dAtA[i:], m.MessageID)
		i = encodeVarintTrace(dAtA, i, uint64(len(m.MessageID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_AddPeer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_AddPeer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_AddPeer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Proto != nil {
		i -= len(*m.Proto)
		copy(dAtA[i:], *m.Proto)
		i = encodeVarintTrace(dAtA, i, uint64(len(*m.Proto)))
		i--
		dAtA[i] = 0x12
	}
	if m.PeerID != nil {
		i -= len(m.PeerID)
		copy(dAtA[i:], m.PeerID)
		i = encodeVarintTrace(dAtA, i, uint64(len(m.PeerID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_RemovePeer) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_RemovePeer) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_RemovePeer) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.PeerID != nil {
		i -= len(m.PeerID)
		copy(dAtA[i:], m.PeerID)
		i = encodeVarintTrace(dAtA, i, uint64(len(m.PeerID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_RecvRPC) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_RecvRPC) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_RecvRPC) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Meta != nil {
		{
			size, err := m.Meta.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.ReceivedFrom != nil {
		i -= len(m.ReceivedFrom)
		copy(dAtA[i:], m.ReceivedFrom)
		i = encodeVarintTrace(dAtA, i, uint64(len(m.ReceivedFrom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_SendRPC) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_SendRPC) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_SendRPC) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Meta != nil {
		{
			size, err := m.Meta.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.SendTo != nil {
		i -= len(m.SendTo)
		copy(dAtA[i:], m.SendTo)
		i = encodeVarintTrace(dAtA, i, uint64(len(m.SendTo)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_DropRPC) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_DropRPC) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_DropRPC) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Meta != nil {
		{
			size, err := m.Meta.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.SendTo != nil {
		i -= len(m.SendTo)
		copy(dAtA[i:], m.SendTo)
		i = encodeVarintTrace(dAtA, i, uint64(len(m.SendTo)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_Join) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_Join) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_Join) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Topic != nil {
		i -= len(*m.Topic)
		copy(dAtA[i:], *m.Topic)
		i = encodeVarintTrace(dAtA, i, uint64(len(*m.Topic)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_Leave) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_Leave) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_Leave) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Topic != nil {
		i -= len(*m.Topic)
		copy(dAtA[i:], *m.Topic)
		i = encodeVarintTrace(dAtA, i, uint64(len(*m.Topic)))
		i--
		dAtA[i] = 0x12
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_Graft) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_Graft) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_Graft) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Topic != nil {
		i -= len(*m.Topic)
		copy(dAtA[i:], *m.Topic)
		i = encodeVarintTrace(dAtA, i, uint64(len(*m.Topic)))
		i--
		dAtA[i] = 0x12
	}
	if m.PeerID != nil {
		i -= len(m.PeerID)
		copy(dAtA[i:], m.PeerID)
		i = encodeVarintTrace(dAtA, i, uint64(len(m.PeerID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_Prune) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_Prune) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_Prune) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Topic != nil {
		i -= len(*m.Topic)
		copy(dAtA[i:], *m.Topic)
		i = encodeVarintTrace(dAtA, i, uint64(len(*m.Topic)))
		i--
		dAtA[i] = 0x12
	}
	if m.PeerID != nil {
		i -= len(m.PeerID)
		copy(dAtA[i:], m.PeerID)
		i = encodeVarintTrace(dAtA, i, uint64(len(m.PeerID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_RPCMeta) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_RPCMeta) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_RPCMeta) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Control != nil {
		{
			size, err := m.Control.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintTrace(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Subscription) > 0 {
		for iNdEx := len(m.Subscription) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Subscription[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTrace(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Messages) > 0 {
		for iNdEx := len(m.Messages) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Messages[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTrace(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_MessageMeta) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_MessageMeta) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_MessageMeta) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Topic != nil {
		i -= len(*m.Topic)
		copy(dAtA[i:], *m.Topic)
		i = encodeVarintTrace(dAtA, i, uint64(len(*m.Topic)))
		i--
		dAtA[i] = 0x12
	}
	if m.MessageID != nil {
		i -= len(m.MessageID)
		copy(dAtA[i:], m.MessageID)
		i = encodeVarintTrace(dAtA, i, uint64(len(m.MessageID)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_SubMeta) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_SubMeta) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_SubMeta) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Topic != nil {
		i -= len(*m.Topic)
		copy(dAtA[i:], *m.Topic)
		i = encodeVarintTrace(dAtA, i, uint64(len(*m.Topic)))
		i--
		dAtA[i] = 0x12
	}
	if m.Subscribe != nil {
		i--
		if *m.Subscribe {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_ControlMeta) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_ControlMeta) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_ControlMeta) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Prune) > 0 {
		for iNdEx := len(m.Prune) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Prune[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTrace(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Graft) > 0 {
		for iNdEx := len(m.Graft) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Graft[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTrace(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x1a
		}
	}
	if len(m.Iwant) > 0 {
		for iNdEx := len(m.Iwant) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Iwant[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTrace(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Ihave) > 0 {
		for iNdEx := len(m.Ihave) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Ihave[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTrace(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_ControlIHaveMeta) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_ControlIHaveMeta) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_ControlIHaveMeta) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.MessageIDs) > 0 {
		for iNdEx := len(m.MessageIDs) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.MessageIDs[iNdEx])
			copy(dAtA[i:], m.MessageIDs[iNdEx])
			i = encodeVarintTrace(dAtA, i, uint64(len(m.MessageIDs[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Topic != nil {
		i -= len(*m.Topic)
		copy(dAtA[i:], *m.Topic)
		i = encodeVarintTrace(dAtA, i, uint64(len(*m.Topic)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_ControlIWantMeta) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_ControlIWantMeta) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_ControlIWantMeta) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.MessageIDs) > 0 {
		for iNdEx := len(m.MessageIDs) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.MessageIDs[iNdEx])
			copy(dAtA[i:], m.MessageIDs[iNdEx])
			i = encodeVarintTrace(dAtA, i, uint64(len(m.MessageIDs[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_ControlGraftMeta) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_ControlGraftMeta) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_ControlGraftMeta) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if m.Topic != nil {
		i -= len(*m.Topic)
		copy(dAtA[i:], *m.Topic)
		i = encodeVarintTrace(dAtA, i, uint64(len(*m.Topic)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TraceEvent_ControlPruneMeta) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEvent_ControlPruneMeta) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEvent_ControlPruneMeta) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Peers) > 0 {
		for iNdEx := len(m.Peers) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Peers[iNdEx])
			copy(dAtA[i:], m.Peers[iNdEx])
			i = encodeVarintTrace(dAtA, i, uint64(len(m.Peers[iNdEx])))
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Topic != nil {
		i -= len(*m.Topic)
		copy(dAtA[i:], *m.Topic)
		i = encodeVarintTrace(dAtA, i, uint64(len(*m.Topic)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *TraceEventBatch) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TraceEventBatch) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TraceEventBatch) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.XXX_unrecognized != nil {
		i -= len(m.XXX_unrecognized)
		copy(dAtA[i:], m.XXX_unrecognized)
	}
	if len(m.Batch) > 0 {
		for iNdEx := len(m.Batch) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Batch[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTrace(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintTrace(dAtA []byte, offset int, v uint64) int {
	offset -= sovTrace(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TraceEvent) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Type != nil {
		n += 1 + sovTrace(uint64(*m.Type))
	}
	if m.PeerID != nil {
		l = len(m.PeerID)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Timestamp != nil {
		n += 1 + sovTrace(uint64(*m.Timestamp))
	}
	if m.PublishMessage != nil {
		l = m.PublishMessage.Size()
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.RejectMessage != nil {
		l = m.RejectMessage.Size()
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.DuplicateMessage != nil {
		l = m.DuplicateMessage.Size()
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.DeliverMessage != nil {
		l = m.DeliverMessage.Size()
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.AddPeer != nil {
		l = m.AddPeer.Size()
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.RemovePeer != nil {
		l = m.RemovePeer.Size()
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.RecvRPC != nil {
		l = m.RecvRPC.Size()
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.SendRPC != nil {
		l = m.SendRPC.Size()
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.DropRPC != nil {
		l = m.DropRPC.Size()
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Join != nil {
		l = m.Join.Size()
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Leave != nil {
		l = m.Leave.Size()
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Graft != nil {
		l = m.Graft.Size()
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Prune != nil {
		l = m.Prune.Size()
		n += 2 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_PublishMessage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MessageID != nil {
		l = len(m.MessageID)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Topic != nil {
		l = len(*m.Topic)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_RejectMessage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MessageID != nil {
		l = len(m.MessageID)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.ReceivedFrom != nil {
		l = len(m.ReceivedFrom)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Reason != nil {
		l = len(*m.Reason)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Topic != nil {
		l = len(*m.Topic)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_DuplicateMessage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MessageID != nil {
		l = len(m.MessageID)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.ReceivedFrom != nil {
		l = len(m.ReceivedFrom)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Topic != nil {
		l = len(*m.Topic)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_DeliverMessage) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MessageID != nil {
		l = len(m.MessageID)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Topic != nil {
		l = len(*m.Topic)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.ReceivedFrom != nil {
		l = len(m.ReceivedFrom)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_AddPeer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PeerID != nil {
		l = len(m.PeerID)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Proto != nil {
		l = len(*m.Proto)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_RemovePeer) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PeerID != nil {
		l = len(m.PeerID)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_RecvRPC) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.ReceivedFrom != nil {
		l = len(m.ReceivedFrom)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Meta != nil {
		l = m.Meta.Size()
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_SendRPC) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SendTo != nil {
		l = len(m.SendTo)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Meta != nil {
		l = m.Meta.Size()
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_DropRPC) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.SendTo != nil {
		l = len(m.SendTo)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Meta != nil {
		l = m.Meta.Size()
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_Join) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Topic != nil {
		l = len(*m.Topic)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_Leave) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Topic != nil {
		l = len(*m.Topic)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_Graft) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PeerID != nil {
		l = len(m.PeerID)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Topic != nil {
		l = len(*m.Topic)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_Prune) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PeerID != nil {
		l = len(m.PeerID)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Topic != nil {
		l = len(*m.Topic)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_RPCMeta) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Messages) > 0 {
		for _, e := range m.Messages {
			l = e.Size()
			n += 1 + l + sovTrace(uint64(l))
		}
	}
	if len(m.Subscription) > 0 {
		for _, e := range m.Subscription {
			l = e.Size()
			n += 1 + l + sovTrace(uint64(l))
		}
	}
	if m.Control != nil {
		l = m.Control.Size()
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_MessageMeta) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MessageID != nil {
		l = len(m.MessageID)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.Topic != nil {
		l = len(*m.Topic)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_SubMeta) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Subscribe != nil {
		n += 2
	}
	if m.Topic != nil {
		l = len(*m.Topic)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_ControlMeta) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Ihave) > 0 {
		for _, e := range m.Ihave {
			l = e.Size()
			n += 1 + l + sovTrace(uint64(l))
		}
	}
	if len(m.Iwant) > 0 {
		for _, e := range m.Iwant {
			l = e.Size()
			n += 1 + l + sovTrace(uint64(l))
		}
	}
	if len(m.Graft) > 0 {
		for _, e := range m.Graft {
			l = e.Size()
			n += 1 + l + sovTrace(uint64(l))
		}
	}
	if len(m.Prune) > 0 {
		for _, e := range m.Prune {
			l = e.Size()
			n += 1 + l + sovTrace(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_ControlIHaveMeta) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Topic != nil {
		l = len(*m.Topic)
		n += 1 + l + sovTrace(uint64(l))
	}
	if len(m.MessageIDs) > 0 {
		for _, b := range m.MessageIDs {
			l = len(b)
			n += 1 + l + sovTrace(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_ControlIWantMeta) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.MessageIDs) > 0 {
		for _, b := range m.MessageIDs {
			l = len(b)
			n += 1 + l + sovTrace(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_ControlGraftMeta) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Topic != nil {
		l = len(*m.Topic)
		n += 1 + l + sovTrace(uint64(l))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEvent_ControlPruneMeta) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Topic != nil {
		l = len(*m.Topic)
		n += 1 + l + sovTrace(uint64(l))
	}
	if len(m.Peers) > 0 {
		for _, b := range m.Peers {
			l = len(b)
			n += 1 + l + sovTrace(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func (m *TraceEventBatch) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Batch) > 0 {
		for _, e := range m.Batch {
			l = e.Size()
			n += 1 + l + sovTrace(uint64(l))
		}
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovTrace(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTrace(x uint64) (n int) {
	return sovTrace(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TraceEvent) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TraceEvent: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TraceEvent: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var v TraceEvent_Type
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= TraceEvent_Type(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Type = &v
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PeerID", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PeerID = append(m.PeerID[:0], dAtA[iNdEx:postIndex]...)
			if m.PeerID == nil {
				m.PeerID = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var v int64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Timestamp = &v
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PublishMessage", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.PublishMessage == nil {
				m.PublishMessage = &TraceEvent_PublishMessage{}
			}
			if err := m.PublishMessage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RejectMessage", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.RejectMessage == nil {
				m.RejectMessage = &TraceEvent_RejectMessage{}
			}
			if err := m.RejectMessage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DuplicateMessage", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.DuplicateMessage == nil {
				m.DuplicateMessage = &TraceEvent_DuplicateMessage{}
			}
			if err := m.DuplicateMessage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DeliverMessage", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.DeliverMessage == nil {
				m.DeliverMessage = &TraceEvent_DeliverMessage{}
			}
			if err := m.DeliverMessage.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AddPeer", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.AddPeer == nil {
				m.AddPeer = &TraceEvent_AddPeer{}
			}
			if err := m.AddPeer.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RemovePeer", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.RemovePeer == nil {
				m.RemovePeer = &TraceEvent_RemovePeer{}
			}
			if err := m.RemovePeer.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RecvRPC", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.RecvRPC == nil {
				m.RecvRPC = &TraceEvent_RecvRPC{}
			}
			if err := m.RecvRPC.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SendRPC", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.SendRPC == nil {
				m.SendRPC = &TraceEvent_SendRPC{}
			}
			if err := m.SendRPC.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DropRPC", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.DropRPC == nil {
				m.DropRPC = &TraceEvent_DropRPC{}
			}
			if err := m.DropRPC.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Join", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Join == nil {
				m.Join = &TraceEvent_Join{}
			}
			if err := m.Join.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 14:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Leave", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Leave == nil {
				m.Leave = &TraceEvent_Leave{}
			}
			if err := m.Leave.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 15:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Graft", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Graft == nil {
				m.Graft = &TraceEvent_Graft{}
			}
			if err := m.Graft.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 16:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Prune", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Prune == nil {
				m.Prune = &TraceEvent_Prune{}
			}
			if err := m.Prune.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_PublishMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: PublishMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PublishMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MessageID", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MessageID = append(m.MessageID[:0], dAtA[iNdEx:postIndex]...)
			if m.MessageID == nil {
				m.MessageID = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Topic = &s
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_RejectMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RejectMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RejectMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MessageID", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MessageID = append(m.MessageID[:0], dAtA[iNdEx:postIndex]...)
			if m.MessageID == nil {
				m.MessageID = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReceivedFrom", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReceivedFrom = append(m.ReceivedFrom[:0], dAtA[iNdEx:postIndex]...)
			if m.ReceivedFrom == nil {
				m.ReceivedFrom = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Reason", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Reason = &s
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Topic = &s
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_DuplicateMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DuplicateMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DuplicateMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MessageID", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MessageID = append(m.MessageID[:0], dAtA[iNdEx:postIndex]...)
			if m.MessageID == nil {
				m.MessageID = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReceivedFrom", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReceivedFrom = append(m.ReceivedFrom[:0], dAtA[iNdEx:postIndex]...)
			if m.ReceivedFrom == nil {
				m.ReceivedFrom = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Topic = &s
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_DeliverMessage) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DeliverMessage: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DeliverMessage: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MessageID", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MessageID = append(m.MessageID[:0], dAtA[iNdEx:postIndex]...)
			if m.MessageID == nil {
				m.MessageID = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Topic = &s
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReceivedFrom", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReceivedFrom = append(m.ReceivedFrom[:0], dAtA[iNdEx:postIndex]...)
			if m.ReceivedFrom == nil {
				m.ReceivedFrom = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_AddPeer) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AddPeer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AddPeer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PeerID", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PeerID = append(m.PeerID[:0], dAtA[iNdEx:postIndex]...)
			if m.PeerID == nil {
				m.PeerID = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Proto", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Proto = &s
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_RemovePeer) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RemovePeer: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RemovePeer: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PeerID", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PeerID = append(m.PeerID[:0], dAtA[iNdEx:postIndex]...)
			if m.PeerID == nil {
				m.PeerID = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_RecvRPC) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RecvRPC: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RecvRPC: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ReceivedFrom", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ReceivedFrom = append(m.ReceivedFrom[:0], dAtA[iNdEx:postIndex]...)
			if m.ReceivedFrom == nil {
				m.ReceivedFrom = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Meta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Meta == nil {
				m.Meta = &TraceEvent_RPCMeta{}
			}
			if err := m.Meta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_SendRPC) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SendRPC: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SendRPC: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SendTo", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SendTo = append(m.SendTo[:0], dAtA[iNdEx:postIndex]...)
			if m.SendTo == nil {
				m.SendTo = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Meta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Meta == nil {
				m.Meta = &TraceEvent_RPCMeta{}
			}
			if err := m.Meta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_DropRPC) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: DropRPC: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DropRPC: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SendTo", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SendTo = append(m.SendTo[:0], dAtA[iNdEx:postIndex]...)
			if m.SendTo == nil {
				m.SendTo = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Meta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Meta == nil {
				m.Meta = &TraceEvent_RPCMeta{}
			}
			if err := m.Meta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_Join) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Join: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Join: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Topic = &s
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_Leave) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Leave: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Leave: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Topic = &s
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_Graft) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Graft: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Graft: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PeerID", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PeerID = append(m.PeerID[:0], dAtA[iNdEx:postIndex]...)
			if m.PeerID == nil {
				m.PeerID = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Topic = &s
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_Prune) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Prune: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Prune: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PeerID", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PeerID = append(m.PeerID[:0], dAtA[iNdEx:postIndex]...)
			if m.PeerID == nil {
				m.PeerID = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Topic = &s
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_RPCMeta) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: RPCMeta: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: RPCMeta: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Messages", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Messages = append(m.Messages, &TraceEvent_MessageMeta{})
			if err := m.Messages[len(m.Messages)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Subscription", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Subscription = append(m.Subscription, &TraceEvent_SubMeta{})
			if err := m.Subscription[len(m.Subscription)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Control", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Control == nil {
				m.Control = &TraceEvent_ControlMeta{}
			}
			if err := m.Control.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_MessageMeta) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MessageMeta: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MessageMeta: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MessageID", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MessageID = append(m.MessageID[:0], dAtA[iNdEx:postIndex]...)
			if m.MessageID == nil {
				m.MessageID = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Topic = &s
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_SubMeta) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: SubMeta: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SubMeta: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Subscribe", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			b := bool(v != 0)
			m.Subscribe = &b
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Topic = &s
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_ControlMeta) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ControlMeta: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ControlMeta: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ihave", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Ihave = append(m.Ihave, &TraceEvent_ControlIHaveMeta{})
			if err := m.Ihave[len(m.Ihave)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Iwant", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Iwant = append(m.Iwant, &TraceEvent_ControlIWantMeta{})
			if err := m.Iwant[len(m.Iwant)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Graft", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Graft = append(m.Graft, &TraceEvent_ControlGraftMeta{})
			if err := m.Graft[len(m.Graft)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Prune", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Prune = append(m.Prune, &TraceEvent_ControlPruneMeta{})
			if err := m.Prune[len(m.Prune)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_ControlIHaveMeta) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ControlIHaveMeta: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ControlIHaveMeta: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Topic = &s
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MessageIDs", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MessageIDs = append(m.MessageIDs, make([]byte, postIndex-iNdEx))
			copy(m.MessageIDs[len(m.MessageIDs)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_ControlIWantMeta) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ControlIWantMeta: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ControlIWantMeta: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MessageIDs", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MessageIDs = append(m.MessageIDs, make([]byte, postIndex-iNdEx))
			copy(m.MessageIDs[len(m.MessageIDs)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_ControlGraftMeta) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ControlGraftMeta: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ControlGraftMeta: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Topic = &s
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEvent_ControlPruneMeta) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ControlPruneMeta: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ControlPruneMeta: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Topic", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			s := string(dAtA[iNdEx:postIndex])
			m.Topic = &s
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Peers", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Peers = append(m.Peers, make([]byte, postIndex-iNdEx))
			copy(m.Peers[len(m.Peers)-1], dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *TraceEventBatch) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TraceEventBatch: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TraceEventBatch: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Batch", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthTrace
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTrace
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Batch = append(m.Batch, &TraceEvent{})
			if err := m.Batch[len(m.Batch)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTrace(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTrace
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTrace(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTrace
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTrace
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTrace
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTrace
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTrace
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTrace        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTrace          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTrace = fmt.Errorf("proto: unexpected end of group")
)
