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

package p2p

import (
	"fmt"
	"github.com/AlayaNetwork/Alaya-Go/log"
	"github.com/AlayaNetwork/Alaya-Go/p2p/pubsub"
	"sync/atomic"
)

const (

	// Msg code of PubSub's message
	PubSubMsgCode = 0xff
)

// Stream is the stream type used by pubsub. In general
type Stream struct {
	conn     pubsub.Conn
	rw       MsgReadWriter
	protocol atomic.Value
	errCh    chan error
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
func (s *Stream) Protocol() pubsub.ProtocolID {
	// Ignore type error. It means that the protocol is unset.
	p, _ := s.protocol.Load().(pubsub.ProtocolID)
	return p
}

// SetProtocol sets the protocol for this stream.
//
// This doesn't actually *do* anything other than record the fact that we're
// speaking the given protocol over this stream. It's still up to the user to
// negotiate the protocol. This is usually done by the Host.
func (s *Stream) SetProtocol(p pubsub.ProtocolID) {
	s.protocol.Store(p)
}

func (s *Stream) Read(data interface{}) error {
	msg, err := s.rw.ReadMsg()
	if err != nil {
		log.Debug("Failed to read PubSub message", "err", err)
		return err
	}

	if err := msg.Decode(data); err != nil {
		log.Error("Decode PubSub message fail", "err", err)
		return err
	}
	return nil
}

func (s *Stream) Write(data interface{}) error {
	if err := Send(s.rw, PubSubMsgCode, data); err != nil {
		log.Error("Failed to send PubSub message", "err", err)
		return err
	}
	return nil
}

func (s *Stream) Close(err error) {
	s.errCh <- err
}

func NewStream(conn pubsub.Conn, rw MsgReadWriter, errCh chan error, id pubsub.ProtocolID) *Stream {
	s := &Stream{
		conn:  conn,
		rw:    rw,
		errCh: errCh,
	}
	s.SetProtocol(id)
	return s
}
