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
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/p2p/pubsub"
	"sync"
	"time"
)

type PeerManager interface {
}

type Network struct {
	sync.RWMutex
	server *Server

	m map[pubsub.Notifiee]struct{}

	conns struct {
		sync.RWMutex
		m map[enode.ID][]pubsub.Conn
	}
}

func NewNetwork(server *Server) *Network {
	n := &Network{
		RWMutex: sync.RWMutex{},
		server:  server,
		m:       make(map[pubsub.Notifiee]struct{}),
	}
	n.conns.m = make(map[enode.ID][]pubsub.Conn)
	return n
}

func (n *Network) SetConn(p enode.ID, conn pubsub.Conn) {
	n.conns.RLock()
	defer n.conns.RUnlock()
	conns := n.conns.m[p]
	if conns == nil {
		conns = make([]pubsub.Conn, 0)
	}
	conns = append(conns, conn)
	n.conns.m[p] = conns
}

func (n *Network) ConnsToPeer(p enode.ID) []pubsub.Conn {
	n.conns.RLock()
	defer n.conns.RUnlock()
	conns := n.conns.m[p]
	output := make([]pubsub.Conn, len(conns))
	for i, c := range conns {
		output[i] = c
	}
	return output
}

func (n *Network) Connectedness(id enode.ID) pubsub.Connectedness {
	for _, p := range n.server.Peers() {
		if p.ID() == id {
			return pubsub.Connected
		}
	}
	return pubsub.NotConnected
}

func (n *Network) Notify(f pubsub.Notifiee) {
	n.Lock()
	n.m[f] = struct{}{}
	n.Unlock()
}

func (n *Network) StopNotify(f pubsub.Notifiee) {
	n.Lock()
	delete(n.m, f)
	n.Unlock()
}

// notifyAll sends a signal to all Notifiees
func (n *Network) NotifyAll(conn pubsub.Conn) {
	var wg sync.WaitGroup

	n.RLock()
	wg.Add(len(n.m))
	for f := range n.m {
		go func(f pubsub.Notifiee) {
			defer wg.Done()
			f.Connected(n, conn)
		}(f)
	}

	wg.Wait()
	n.RUnlock()
}

func (n *Network) Peers() []enode.ID {
	var eids []enode.ID
	for _, p := range n.server.Peers() {
		eids = append(eids, p.ID())
	}
	return eids
}

type Conn struct {
	remote *enode.Node
	stat   pubsub.Stat

	streams struct {
		sync.Mutex
		m map[pubsub.Stream]struct{}
	}
}

func NewConn(node *enode.Node, inbound bool) *Conn {
	stat := pubsub.Stat{
		Opened: time.Now(),
		Extra:  make(map[interface{}]interface{}),
	}
	if inbound {
		stat.Direction = pubsub.DirInbound
	} else {
		stat.Direction = pubsub.DirOutbound
	}
	conn := &Conn{
		remote: node,
		stat:   stat,
	}
	conn.streams.m = make(map[pubsub.Stream]struct{})
	return conn
}

func (c *Conn) ID() string {
	return c.remote.ID().String()
}

func (c *Conn) SetStream(stream pubsub.Stream) {
	c.streams.Lock()
	defer c.streams.Unlock()
	c.streams.m[stream] = struct{}{}
}

func (c *Conn) GetStreams() []pubsub.Stream {
	c.streams.Lock()
	defer c.streams.Unlock()
	streams := make([]pubsub.Stream, 0, len(c.streams.m))
	for s := range c.streams.m {
		streams = append(streams, s)
	}
	return streams
}

func (c *Conn) Stat() pubsub.Stat {
	return c.stat
}

func (c *Conn) RemotePeer() *enode.Node {
	return c.remote
}
