package pubsub

import (
	"bytes"
	"context"
	crand "crypto/rand"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enr"
	"github.com/AlayaNetwork/Alaya-Go/rlp"
	"github.com/libp2p/go-libp2p-core/connmgr"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type TestStream struct {
	conn  Conn
	read  chan []byte
	write chan []byte
}

func (s *TestStream) Protocol() ProtocolID {
	return GossipSubID_v11
}

// Conn returns the connection this stream is part of.
func (s *TestStream) Conn() Conn {
	return s.conn
}

func (s *TestStream) Read(data interface{}) error {
	outData := <-s.read
	return rlp.DecodeBytes(outData, data)
}

func (s *TestStream) Write(data interface{}) error {
	enVal, err := rlp.EncodeToBytes(data)
	if err != nil {
		return err
	}
	s.write <- enVal
	return nil
}

func (s *TestStream) Close(err error) {
}

type TestConn struct {
	remote *enode.Node
	stat   Stat

	streams struct {
		sync.Mutex
		m map[Stream]struct{}
	}
}

func (c *TestConn) ID() string {
	return c.remote.ID().String()
}

func (c *TestConn) GetStreams() []Stream {
	c.streams.Lock()
	defer c.streams.Unlock()
	streams := make([]Stream, 0, len(c.streams.m))
	for s := range c.streams.m {
		streams = append(streams, s)
	}
	return streams
}

func (c *TestConn) Stat() Stat {
	return c.stat
}

func (c *TestConn) RemotePeer() *enode.Node {
	return c.remote
}

func (c *TestConn) Close() error {
	c.streams.Lock()
	defer c.streams.Unlock()
	for s := range c.streams.m {
		s.Close(nil)
	}
	c.streams.m = nil
	return nil
}

type TestNetwork struct {
	sync.RWMutex
	m map[Notifiee]struct{}

	conns struct {
		sync.RWMutex
		m map[enode.ID][]Conn
	}
}

func (n *TestNetwork) ConnsToPeer(p enode.ID) []Conn {
	n.conns.Lock()
	defer n.conns.Unlock()
	return n.conns.m[p]
}

func (n *TestNetwork) Connectedness(enode.ID) Connectedness {
	return NotConnected
}

func (n *TestNetwork) Notify(nf Notifiee) {
	n.Lock()
	n.m[nf] = struct{}{}
	n.Unlock()
}

// notifyAll sends a signal to all Notifiees
func (n *TestNetwork) NotifyAll(conn Conn) {
	var wg sync.WaitGroup

	n.RLock()
	wg.Add(len(n.m))
	for f := range n.m {
		go func(f Notifiee) {
			defer wg.Done()
			f.Connected(n, conn)
		}(f)
	}

	wg.Wait()
	n.RUnlock()
}

func (n *TestNetwork) Peers() []enode.ID {
	n.conns.Lock()
	defer n.conns.Unlock()
	var peers []enode.ID
	for pid := range n.conns.m {
		peers = append(peers, pid)
	}
	return peers
}

func (n *TestNetwork) Close() error {
	n.conns.Lock()
	defer n.conns.Unlock()
	for _, c := range n.conns.m {
		if err := c[0].Close(); err != nil {
			return err
		}
	}
	return nil
}

func (n *TestNetwork) SetConn(p enode.ID, conn Conn) {
	n.conns.RLock()
	defer n.conns.RUnlock()
	conns := make([]Conn, 0, 1)
	conns = append(conns, conn)
	n.conns.m[p] = conns
}

func (n *TestNetwork) Conns() []Conn {
	n.conns.RLock()
	defer n.conns.RUnlock()
	var connList []Conn
	for _, cs := range n.conns.m {
		connList = append(connList, cs...)
	}
	return connList
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
	netw := &TestNetwork{
		m: make(map[Notifiee]struct{}),
	}
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

func (h *TestHost) StreamHandler(pid ProtocolID) StreamHandler {
	h.handlerLock.Lock()
	defer h.handlerLock.Unlock()
	return h.handlers[pid]
}

func (h *TestHost) RemoveStreamHandler(pid ProtocolID) {
	h.handlerLock.Lock()
	defer h.handlerLock.Unlock()
	delete(h.handlers, pid)
}

func (h *TestHost) NewStream(ctx context.Context, p enode.ID, pids ...ProtocolID) (Stream, error) {
	return h.newStream(p, make(chan []byte), make(chan []byte))
}

func (h *TestHost) newStream(p enode.ID, read chan []byte, write chan []byte) (Stream, error) {
	if conns := h.network.ConnsToPeer(p); len(conns) > 0 {
		if streams := conns[0].GetStreams(); len(streams) > 0 {
			return streams[0], nil
		}
	}
	conn := &TestConn{
		remote: enode.SignNull(new(enr.Record), p),
		stat:   Stat{},
	}
	conn.streams.m = make(map[Stream]struct{})

	stream := &TestStream{
		conn:  conn,
		read:  read,
		write: write,
	}
	conn.streams.m[stream] = struct{}{}

	h.network.SetConn(p, conn)
	return stream, nil
}

func (h *TestHost) Close() error {
	return h.network.Close()
}

func (h *TestHost) ConnManager() connmgr.ConnManager {
	return &connmgr.NullConnMgr{}
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
	hostA := a.(*TestHost)
	hostB := b.(*TestHost)
	chan1 := make(chan []byte, 100)
	chan2 := make(chan []byte, 100)
	if stream, err := hostA.newStream(b.ID().ID(), chan2, chan1); err != nil {
		t.Fatal(err)
	} else {
		hostA.network.NotifyAll(stream.Conn())
	}

	if stream, err := hostB.newStream(a.ID().ID(), chan1, chan2); err != nil {
		t.Fatal(err)
	} else {
		hostB.network.NotifyAll(stream.Conn())
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

//See https://github.com/libp2p/go-libp2p-pubsub/issues/426
func TestPubSubRemovesBlacklistedPeer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	hosts := getNetHosts(t, ctx, 2)

	bl := NewMapBlacklist()

	psubs0 := getGossipsub(ctx, hosts[0])
	psubs1 := getGossipsub(ctx, hosts[1], WithBlacklist(bl))
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
