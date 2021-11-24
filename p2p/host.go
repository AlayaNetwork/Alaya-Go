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
	"context"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/p2p/pubsub"
	"github.com/libp2p/go-libp2p-core/connmgr"
	"sync"
)

type Host struct {
	node    *enode.Node
	network *Network
	streams map[enode.ID]pubsub.Stream
	sync.Mutex
}

func NewHost(localNode *enode.Node, network *Network) *Host {
	host := &Host{
		node:    localNode,
		streams: make(map[enode.ID]pubsub.Stream),
		network: network,
	}
	return host
}

func (h *Host) ID() *enode.Node {
	return h.node
}

func (h *Host) Network() pubsub.Network {
	return h.network
}

func (h *Host) Connect(ctx context.Context, pi enode.ID) error {
	return nil
}

func (h *Host) SetStreamHandler(pid pubsub.ProtocolID, handler pubsub.StreamHandler) {

}

func (h *Host) SetStreamHandlerMatch(pubsub.ProtocolID, func(string) bool, pubsub.StreamHandler) {

}

func (h *Host) RemoveStreamHandler(pid pubsub.ProtocolID) {

}

func (h *Host) NewStream(ctx context.Context, nodeId enode.ID, pids ...pubsub.ProtocolID) (pubsub.Stream, error) {
	h.Lock()
	defer h.Unlock()
	return h.streams[nodeId], nil
}

func (h *Host) SetStream(nodeId enode.ID, stream pubsub.Stream) {
	h.Lock()
	defer h.Unlock()
	h.streams[nodeId] = stream
}

func (h *Host) Close() error {
	return nil
}

func (h *Host) ConnManager() connmgr.ConnManager {
	return nil
}

func (h *Host) NotifyAll(conn pubsub.Conn) {
	h.network.NotifyAll(conn)
}

func (h *Host) AddConn(p enode.ID, conn pubsub.Conn) {
	h.network.SetConn(p, conn)
}
