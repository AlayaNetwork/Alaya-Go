package pubsub

import (
	"bytes"
	"context"
	crand "crypto/rand"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enr"
	"github.com/libp2p/go-libp2p-core/connmgr"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type TestNetwork struct {
	sync.RWMutex

	conns struct {
		sync.RWMutex
		m map[enode.ID][]Conn
	}
}

func (n *TestNetwork) ConnsToPeer(p enode.ID) []Conn {
	return nil
}

func (n *TestNetwork) Connectedness(enode.ID) Connectedness {
	return NotConnected
}

func (n *TestNetwork) Notify(Notifiee) {

}

func (n *TestNetwork) Peers() []enode.ID {
	return nil
}

type TestHost struct {
	network     *TestNetwork
	Node        *enode.Node
	handlers    map[ProtocolID]StreamHandler
	handlerLock sync.Mutex
}

func NewTestHost() *TestHost {
	var p enode.ID
	crand.Read(p[:])
	netw := &TestNetwork{}
	netw.conns.m = make(map[enode.ID][]Conn)
	h := &TestHost{
		network:     netw,
		Node:        enode.SignNull(new(enr.Record), p),
		handlers:    make(map[ProtocolID]StreamHandler),
		handlerLock: sync.Mutex{},
	}
	return h
}

func (h *TestHost) ID() *enode.Node {
	return h.Node
}

func (h *TestHost) Network() Network {
	return h.network
}

func (h *TestHost) Connect(ctx context.Context, pi enode.ID) error {
	return nil
}

func (h *TestHost) SetStreamHandler(pid ProtocolID, handler StreamHandler) {
	h.handlerLock.Lock()
	defer h.handlerLock.Unlock()
	h.handlers[pid] = handler
}

func (h *TestHost) SetStreamHandlerMatch(ProtocolID, func(string) bool, StreamHandler) {
}

func (h *TestHost) RemoveStreamHandler(pid ProtocolID) {
	h.handlerLock.Lock()
	defer h.handlerLock.Unlock()
	delete(h.handlers, pid)
}

func (h *TestHost) NewStream(ctx context.Context, p enode.ID, pids ...ProtocolID) (Stream, error) {
	return nil, nil
}

func (h *TestHost) Close() error {
	return nil
}

func (h *TestHost) ConnManager() connmgr.ConnManager {
	return nil
}

func checkMessageRouting(t *testing.T, topic string, pubs []*PubSub, subs []*Subscription) {
	data := make([]byte, 16)
	rand.Read(data)

	for _, p := range pubs {
		err := p.Publish(topic, data)
		if err != nil {
			t.Fatal(err)
		}

		for _, s := range subs {
			assertReceive(t, s, data)
		}
	}
}

func getNetHosts(t *testing.T, ctx context.Context, n int) []Host {
	var out []Host

	for i := 0; i < n; i++ {
		out = append(out, NewTestHost())
	}

	return out
}

func connect(t *testing.T, a, b Host) {
	err := b.Connect(context.Background(), a.ID().ID())
	if err != nil {
		t.Fatal(err)
	}
}

func sparseConnect(t *testing.T, hosts []Host) {
	connectSome(t, hosts, 3)
}

func denseConnect(t *testing.T, hosts []Host) {
	connectSome(t, hosts, 10)
}

func connectSome(t *testing.T, hosts []Host, d int) {
	for i, a := range hosts {
		for j := 0; j < d; j++ {
			n := rand.Intn(len(hosts))
			if n == i {
				j--
				continue
			}

			b := hosts[n]

			connect(t, a, b)
		}
	}
}

func connectAll(t *testing.T, hosts []Host) {
	for i, a := range hosts {
		for j, b := range hosts {
			if i == j {
				continue
			}

			connect(t, a, b)
		}
	}
}

func getPubsub(ctx context.Context, h Host, opts ...Option) *PubSub {
	ps, err := NewGossipSub(ctx, h, opts...)
	if err != nil {
		panic(err)
	}
	return ps
}

func getPubsubs(ctx context.Context, hs []Host, opts ...Option) []*PubSub {
	var psubs []*PubSub
	for _, h := range hs {
		psubs = append(psubs, getPubsub(ctx, h, opts...))
	}
	return psubs
}

func assertReceive(t *testing.T, ch *Subscription, exp []byte) {
	select {
	case msg := <-ch.ch:
		if !bytes.Equal(msg.GetData(), exp) {
			t.Fatalf("got wrong message, expected %s but got %s", string(exp), string(msg.GetData()))
		}
	case <-time.After(time.Second * 5):
		t.Logf("%#v\n", ch)
		t.Fatal("timed out waiting for message of: ", string(exp))
	}
}

func assertNeverReceives(t *testing.T, ch *Subscription, timeout time.Duration) {
	select {
	case msg := <-ch.ch:
		t.Logf("%#v\n", ch)
		t.Fatal("got unexpected message: ", string(msg.GetData()))
	case <-time.After(timeout):
	}
}

func assertPeerLists(t *testing.T, hosts []Host, ps *PubSub, has ...int) {
	peers := ps.ListPeers("")
	set := make(map[enode.ID]struct{})
	for _, p := range peers {
		set[p] = struct{}{}
	}

	for _, h := range has {
		if _, ok := set[hosts[h].ID().ID()]; !ok {
			t.Fatal("expected to have connection to peer: ", h)
		}
	}
}

// See https://github.com/libp2p/go-libp2p-pubsub/issues/426
func TestPubSubRemovesBlacklistedPeer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	hosts := getNetHosts(t, ctx, 2)

	bl := NewMapBlacklist()

	psubs0 := getPubsub(ctx, hosts[0])
	psubs1 := getPubsub(ctx, hosts[1], WithBlacklist(bl))
	connect(t, hosts[0], hosts[1])

	// Bad peer is blacklisted after it has connected.
	// Calling p.BlacklistPeer directly does the right thing but we should also clean
	// up the peer if it has been added the the blacklist by another means.
	bl.Add(hosts[0].ID().ID())

	_, err := psubs0.Subscribe("test")
	if err != nil {
		t.Fatal(err)
	}

	sub1, err := psubs1.Subscribe("test")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(time.Millisecond * 100)

	psubs0.Publish("test", []byte("message"))

	wctx, cancel2 := context.WithTimeout(ctx, 1*time.Second)
	defer cancel2()

	_, _ = sub1.Next(wctx)

	// Explicitly cancel context so PubSub cleans up peer channels.
	// Issue 426 reports a panic due to a peer channel being closed twice.
	cancel()
	time.Sleep(time.Millisecond * 100)
}
