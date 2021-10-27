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


package pubsub

import (
	"errors"
	"github.com/AlayaNetwork/Alaya-Go/p2p"
	"sync"
)

// Server manages all pubsub peers.
type Server struct {
	lock    sync.Mutex // protects running
	running bool
	Pb     *PubSub
}

// Start starts running the server.
// Servers can not be re-used after stopping.
func (srv *Server) Start() (err error) {
	srv.lock.Lock()
	defer srv.lock.Unlock()
	if srv.running {
		return errors.New("server already running")
	}
	srv.running = true
	return nil
}

// run is the main loop of the server.
func (s *Server) run() {

}

// After the node is successfully connected and the message belongs
// to the cbft.pubsub protocol message, the method is called.
func (s *Server) Handle(p *p2p.Peer, rw p2p.MsgReadWriter) error {
	return nil
}
