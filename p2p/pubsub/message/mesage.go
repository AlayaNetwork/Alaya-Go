package message

import (
	"errors"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/rlp"
)

type RPC struct {
	Subscriptions []*RPC_SubOpts
	Publish       []*Message
	Control       *ControlMessage
}

func IsEmpty(rpc *RPC) bool {
	if rpc != nil {
		if rpc.Subscriptions != nil || rpc.Publish != nil || rpc.Control != nil {
			return false
		}
	}
	return true
}

func Filling(rpc *RPC) {
	if rpc != nil {
		if rpc.Subscriptions == nil {
			rpc.Subscriptions = make([]*RPC_SubOpts, 0)
		}
		if rpc.Publish == nil {
			rpc.Publish = make([]*Message, 0)
		}
		if rpc.Control == nil {
			rpc.Control = &ControlMessage{}
		}
	}
}

func (m *RPC) Reset() {
	m.Subscriptions = make([]*RPC_SubOpts, 0)
	m.Publish = make([]*Message, 0)
	m.Control = &ControlMessage{}
}

func (m *RPC) Size() (n int) {
	if m == nil {
		return 0
	}
	if m.Subscriptions != nil {
		for _, v := range m.Subscriptions {
			n += v.Size()
		}
	}
	if m.Publish != nil {
		for _, v := range m.Publish {
			n += v.Size()
		}
	}
	if m.Control != nil {
		n += m.Control.Size()
	}
	return n
}

func (m *RPC) GetSubscriptions() []*RPC_SubOpts {
	if m != nil {
		return m.Subscriptions
	}
	return nil
}

func (m *RPC) GetPublish() []*Message {
	if m != nil {
		return m.Publish
	}
	return nil
}

func (m *RPC) GetControl() *ControlMessage {
	if m != nil {
		return m.Control
	}
	return nil
}

type RPC_SubOpts struct {
	Subscribe *bool
	Topicid   *string
}

func (m *RPC_SubOpts) GetSubscribe() bool {
	if m != nil && m.Subscribe != nil {
		return *m.Subscribe
	}
	return false
}

func (m *RPC_SubOpts) GetTopicid() string {
	if m != nil && m.Topicid != nil {
		return *m.Topicid
	}
	return ""
}

func (m *RPC_SubOpts) Size() (n int) {
	if m == nil {
		return 0
	}
	if m.Subscribe != nil {
		n += 1
	}
	if m.Topicid != nil {
		n += 1 + (len(*m.Topicid))
	}
	return n
}

type Message struct {
	From      enode.ID
	Data      []byte
	Seqno     []byte
	Topic     *string
	Signature []byte
	Key       []byte
}

func (m *Message) GetFrom() enode.ID {
	if m != nil {
		return m.From
	}
	return enode.ZeroID
}

func (m *Message) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *Message) GetSeqno() []byte {
	if m != nil {
		return m.Seqno
	}
	return nil
}

func (m *Message) GetTopic() string {
	if m != nil && m.Topic != nil {
		return *m.Topic
	}
	return ""
}

func (m *Message) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *Message) GetKey() []byte {
	if m != nil {
		return m.Key
	}
	return nil
}

func (m *Message) Size() (n int) {
	if m == nil {
		return 0
	}
	n += len(m.From)
	if m.Data != nil {
		n += len(m.Data)
	}
	if m.Seqno != nil {
		n += len(m.Seqno)
	}
	if m.Topic != nil {
		n += len(*m.Topic)
	}
	if m.Signature != nil {
		n += len(m.Signature)
	}
	if m.Key != nil {
		n += len(m.Key)
	}
	return n
}

type ControlMessage struct {
	Ihave []*ControlIHave
	Iwant []*ControlIWant
	Graft []*ControlGraft
	Prune []*ControlPrune
}

func (m *ControlMessage) GetIhave() []*ControlIHave {
	if m != nil {
		return m.Ihave
	}
	return nil
}

func (m *ControlMessage) GetIwant() []*ControlIWant {
	if m != nil {
		return m.Iwant
	}
	return nil
}

func (m *ControlMessage) GetGraft() []*ControlGraft {
	if m != nil {
		return m.Graft
	}
	return nil
}

func (m *ControlMessage) GetPrune() []*ControlPrune {
	if m != nil {
		return m.Prune
	}
	return nil
}

func (m *ControlMessage) Size() (n int) {
	if m == nil {
		return 0
	}
	if m.Ihave != nil {
		for _, v := range m.Ihave {
			n += v.Size()
		}
	}
	if m.Iwant != nil {
		for _, v := range m.Iwant {
			n += v.Size()
		}
	}
	if m.Graft != nil {
		for _, v := range m.Graft {
			n += v.Size()
		}
	}
	if m.Prune != nil {
		for _, v := range m.Prune {
			n += v.Size()
		}
	}
	return n
}

func (m *ControlMessage) Marshal() ([]byte, error) {
	if m != nil {
		return rlp.EncodeToBytes(m)
	}
	return nil, errors.New("serialized object is empty")
}

type ControlIHave struct {
	TopicID    *string
	MessageIDs []string
}

func (m *ControlIHave) GetTopicID() string {
	if m != nil && m.TopicID != nil {
		return *m.TopicID
	}
	return ""
}

func (m *ControlIHave) GetMessageIDs() []string {
	if m != nil {
		return m.MessageIDs
	}
	return nil
}

func (m *ControlIHave) Size() (n int) {
	if m == nil {
		return 0
	}
	if m.TopicID != nil {
		n += len(*m.TopicID)
	}
	if m.MessageIDs != nil {
		for _, mid := range m.MessageIDs {
			n += len(mid)
		}
	}
	return n
}

type ControlIWant struct {
	MessageIDs []string
}

func (m *ControlIWant) GetMessageIDs() []string {
	if m != nil {
		return m.MessageIDs
	}
	return nil
}

func (m *ControlIWant) Size() (n int) {
	if m == nil {
		return 0
	}
	if m.MessageIDs != nil {
		for _, mid := range m.MessageIDs {
			n += len(mid)
		}
	}
	return n
}

type ControlGraft struct {
	TopicID *string
}

func (m *ControlGraft) GetTopicID() string {
	if m != nil && m.TopicID != nil {
		return *m.TopicID
	}
	return ""
}

func (m *ControlGraft) Size() (n int) {
	if m == nil {
		return 0
	}
	if m.TopicID != nil {
		n += len(*m.TopicID)
	}
	return n
}

type ControlPrune struct {
	TopicID *string
	Peers   []*PeerInfo
	Backoff *uint64
}

func (m *ControlPrune) GetTopicID() string {
	if m != nil && m.TopicID != nil {
		return *m.TopicID
	}
	return ""
}

func (m *ControlPrune) GetPeers() []*PeerInfo {
	if m != nil {
		return m.Peers
	}
	return nil
}

func (m *ControlPrune) GetBackoff() uint64 {
	if m != nil && m.Backoff != nil {
		return *m.Backoff
	}
	return 0
}

func (m *ControlPrune) Size() (n int) {
	if m == nil {
		return 0
	}
	if m.TopicID != nil {
		n += len(*m.TopicID)
	}
	if m.Backoff != nil {
		n += 8
	}
	if m.Peers != nil {
		for _, p := range m.Peers {
			n += p.Size()
		}
	}
	return n
}

type PeerInfo struct {
	PeerID           enode.ID
	SignedPeerRecord []byte
}

func (m *PeerInfo) GetPeerID() enode.ID {
	if m != nil {
		return m.PeerID
	}
	return enode.ZeroID
}

func (m *PeerInfo) GetSignedPeerRecord() []byte {
	if m != nil {
		return m.SignedPeerRecord
	}
	return nil
}

func (m *PeerInfo) Size() (n int) {
	if m == nil {
		return 0
	}
	n += len(m.PeerID)
	if m.SignedPeerRecord != nil {
		n += len(m.SignedPeerRecord)
	}
	return n
}
