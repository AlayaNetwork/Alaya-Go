package p2p

import (
	"context"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/p2p/pubsub"
)

type PubSubServer struct {
	p2pServer *Server
	pubSub    *pubsub.PubSub
	host      *Host
}

func NewPubSubServer(localNode *enode.Node, p2pServer *Server) *PubSubServer {
	network := NewNetwork(p2pServer)
	host := NewHost(localNode, network)
	gossipSub, err := pubsub.NewGossipSub(context.Background(), host)
	if err != nil {
		panic("Failed to NewGossipSub: " + err.Error())
	}

	return &PubSubServer{
		p2pServer: p2pServer,
		pubSub:    gossipSub,
		host:      host,
	}
}

func (pss *PubSubServer) Host() *Host {
	return pss.host
}

func (pss *PubSubServer) PubSub() *pubsub.PubSub {
	return pss.pubSub
}

func (pss *PubSubServer) NewConn(peer *Peer, rw MsgReadWriter) chan error {
	conn := NewConn(peer.Node(), peer.Inbound())

	// Wait for the connection to exit
	errCh := make(chan error)

	stream := NewStream(conn, rw, errCh, "")
	conn.SetStream(stream)

	pss.Host().SetStream(peer.ID(), stream)
	pss.Host().NotifyAll(conn)
	return errCh
}
