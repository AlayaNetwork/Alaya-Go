package pubsub

import (
	"context"
	crand "crypto/rand"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enr"
	"testing"
	"time"
)

func TestPeerGater(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var peerAId enode.ID
	crand.Read(peerAId[:])
	peerA := enode.SignNull(new(enr.Record), peerAId)

	peerAip := "1.2.3.4"

	params := NewPeerGaterParams(.1, .9, .999)
	err := params.validate()
	if err != nil {
		t.Fatal(err)
	}

	pg := newPeerGater(ctx, nil, params)
	pg.getIP = func(p enode.ID) string {
		switch p {
		case peerAId:
			return peerAip
		default:
			return "<wtf>"
		}
	}

	pg.AddPeer(peerA, "")

	status := pg.AcceptFrom(peerA)
	if status != AcceptAll {
		t.Fatal("expected AcceptAll")
	}

	msg := &Message{ReceivedFrom: peerA}

	pg.ValidateMessage(msg)
	status = pg.AcceptFrom(peerA)
	if status != AcceptAll {
		t.Fatal("expected AcceptAll")
	}

	pg.RejectMessage(msg, RejectValidationQueueFull)
	status = pg.AcceptFrom(peerA)
	if status != AcceptAll {
		t.Fatal("expected AcceptAll")
	}

	pg.RejectMessage(msg, RejectValidationThrottled)
	status = pg.AcceptFrom(peerA)
	if status != AcceptAll {
		t.Fatal("expected AcceptAll")
	}

	for i := 0; i < 100; i++ {
		pg.RejectMessage(msg, RejectValidationIgnored)
		pg.RejectMessage(msg, RejectValidationFailed)
	}

	accepted := false
	for i := 0; !accepted && i < 1000; i++ {
		status = pg.AcceptFrom(peerA)
		if status == AcceptControl {
			accepted = true
		}
	}
	if !accepted {
		t.Fatal("expected AcceptControl")
	}

	for i := 0; i < 100; i++ {
		pg.DeliverMessage(msg)
	}

	accepted = false
	for i := 0; !accepted && i < 1000; i++ {
		status = pg.AcceptFrom(peerA)
		if status == AcceptAll {
			accepted = true
		}
	}
	if !accepted {
		t.Fatal("expected to accept at least once")
	}

	for i := 0; i < 100; i++ {
		pg.decayStats()
	}

	status = pg.AcceptFrom(peerA)
	if status != AcceptAll {
		t.Fatal("expected AcceptAll")
	}

	pg.RemovePeer(peerAId)
	pg.Lock()
	_, ok := pg.peerStats[peerAId]
	pg.Unlock()
	if ok {
		t.Fatal("still have a stat record for peerA")
	}

	pg.Lock()
	_, ok = pg.ipStats[peerAip]
	pg.Unlock()
	if !ok {
		t.Fatal("expected to still have a stat record for peerA's ip")
	}

	pg.Lock()
	pg.ipStats[peerAip].expire = time.Now()
	pg.Unlock()

	time.Sleep(2 * time.Second)

	pg.Lock()
	_, ok = pg.ipStats["1.2.3.4"]
	pg.Unlock()
	if ok {
		t.Fatal("still have a stat record for peerA's ip")
	}
}
