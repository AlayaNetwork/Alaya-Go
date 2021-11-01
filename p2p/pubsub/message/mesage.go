package message

import "github.com/AlayaNetwork/Alaya-Go/p2p/enode"

type RPC struct {
	Subscriptions []*RPC_SubOpts
	Publish       []*Message
	Control       *ControlMessage
}

func (m *RPC) Size() (n int) {
	return 0
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
	return 0
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
	return 0
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
	return 0
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
	return 0
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
	return 0
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
	return 0
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
