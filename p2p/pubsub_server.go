package p2p

import (
	"context"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	pubsub2 "github.com/AlayaNetwork/Alaya-Go/p2p/pubsub"
)

type PubSubServer struct {
	p2pServer *Server
	pubsub    *pubsub2.PubSub
}

func NewPubSubServer(localNode *enode.Node, p2pServer *Server) *PubSubServer {
	network := NewNetwork(p2pServer)
	host := NewHost(localNode, network)
	gossipSub, err := pubsub2.NewGossipSub(context.Background(), host)
	if err != nil {
		panic("Failed to NewGossipSub: " + err.Error())
	}
	return &PubSubServer{
		p2pServer: p2pServer,
		pubsub:    gossipSub,
	}
}
