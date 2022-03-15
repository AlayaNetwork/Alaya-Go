package pubsub

import (
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
)

type PubSubNotif PubSub

func (p *PubSubNotif) OpenedStream(n Network, s Stream) {
}

func (p *PubSubNotif) ClosedStream(n Network, s Stream) {
}

func (p *PubSubNotif) Connected(n Network, c Conn) {
	// ignore transient connections
	if c.Stat().Transient {
		return
	}

	go func() {
		p.newPeersPrioLk.RLock()
		p.newPeersMx.Lock()
		p.newPeersPend[c.RemotePeer().ID()] = struct{}{}
		p.newPeersMx.Unlock()
		p.newPeersPrioLk.RUnlock()

		select {
		case p.newPeers <- struct{}{}:
		default:
		}
	}()
}

func (p *PubSubNotif) Initialize() {
	isTransient := func(pid enode.ID) bool {
		for _, c := range p.host.Network().ConnsToPeer(pid) {
			if !c.Stat().Transient {
				return false
			}
		}

		return true
	}

	p.newPeersPrioLk.RLock()
	p.newPeersMx.Lock()
	for _, pid := range p.host.Network().Peers() {
		if isTransient(pid) {
			continue
		}

		p.newPeersPend[pid] = struct{}{}
	}
	p.newPeersMx.Unlock()
	p.newPeersPrioLk.RUnlock()

	select {
	case p.newPeers <- struct{}{}:
	default:
	}
}
