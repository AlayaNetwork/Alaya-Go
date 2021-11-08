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

func (h *Host) ID() *enode.Node {
	return h.node
}

func (h *Host) Peerstore() pubsub.Peerstore {
	return nil
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
