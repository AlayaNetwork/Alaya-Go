package p2p

import (
	"fmt"
	"github.com/AlayaNetwork/Alaya-Go/p2p/pubsub"
	"github.com/libp2p/go-libp2p-core/protocol"
	"sync/atomic"
)

// Stream is the stream type used by pubsub. In general
type Stream struct {
	conn   pubsub.Conn
	rw     MsgReadWriter
	protocol atomic.Value
}

func (s *Stream) String() string {
	return fmt.Sprintf(
		"<swarm.Stream local <-> %s>",
		s.conn.RemotePeer(),
	)
}

// Conn returns the Conn associated with this stream, as an network.Conn
func (s *Stream) Conn() pubsub.Conn {
	return s.conn
}

// Protocol returns the protocol negotiated on this stream (if set).
func (s *Stream) Protocol() protocol.ID {
	// Ignore type error. It means that the protocol is unset.
	p, _ := s.protocol.Load().(protocol.ID)
	return p
}

// SetProtocol sets the protocol for this stream.
//
// This doesn't actually *do* anything other than record the fact that we're
// speaking the given protocol over this stream. It's still up to the user to
// negotiate the protocol. This is usually done by the Host.
func (s *Stream) SetProtocol(p protocol.ID) {
	s.protocol.Store(p)
}

func (s *Stream) ReadWriter() MsgReadWriter {
	return s.rw
}

// newStream creates a new Stream.
func newStream(conn pubsub.Conn, rw MsgReadWriter, id protocol.ID) *Stream {
	s := &Stream{
		conn:	conn,
		rw:		rw,
	}
	s.SetProtocol(id)
	return s
}
