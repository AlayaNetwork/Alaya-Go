package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/p2p/pubsub/message"
	"os"
	"testing"
	"time"
)

func testWithTracer(t *testing.T, tracer EventTracer) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	hosts := getNetHosts(t, ctx, 20)
	psubs := getGossipsubs(ctx, hosts,
		WithEventTracer(tracer),
		// to bootstrap from star topology
		WithPeerExchange(true),
		// to exercise the score paths in the tracer
		WithPeerScore(
			&PeerScoreParams{
				TopicScoreCap:    100,
				AppSpecificScore: func(enode.ID) float64 { return 0 },
				DecayInterval:    time.Second,
				DecayToZero:      0.01,
			},
			&PeerScoreThresholds{
				GossipThreshold:             -1,
				PublishThreshold:            -2,
				GraylistThreshold:           -3,
				OpportunisticGraftThreshold: 1,
			}))

	// add a validator that rejects some messages to exercise those code paths in the tracer
	for _, ps := range psubs {
		ps.RegisterTopicValidator("test", func(ctx context.Context, p enode.ID, msg *Message) bool {
			if string(msg.Data) == "invalid!" {
				return false
			} else {
				return true
			}
		})
	}

	// build the star
	for i := 1; i < 20; i++ {
		connect(t, hosts[0], hosts[i])
	}
	connectSome(t, hosts, 6)

	time.Sleep(time.Second)

	// build the mesh
	var subs []*Subscription
	for _, ps := range psubs {
		sub, err := ps.Subscribe("test")
		if err != nil {
			t.Fatal(err)
		}
		go func(sub *Subscription) {
			for {
				_, err := sub.Next(ctx)
				if err != nil {
					return
				}
			}
		}(sub)
		subs = append(subs, sub)
	}

	// wait for the mesh to build
	time.Sleep(5 * time.Second)

	// publish some messages
	for i := 0; i < 20; i++ {
		if i%7 == 0 {
			psubs[i].Publish("test", []byte("invalid!"))
		} else {
			msg := []byte(fmt.Sprintf("message %d", i))
			psubs[i].Publish("test", msg)
		}
	}

	// wait a bit for propagation and call it day
	time.Sleep(time.Second)

	// close all subscriptions to get some leave events
	for _, sub := range subs {
		sub.Cancel()
	}

	// wait for the leave to take effect
	time.Sleep(time.Second)
}

type traceStats struct {
	publish, reject, duplicate, deliver, add, remove, recv, send, drop, join, leave, graft, prune int
}

func (t *traceStats) process(evt *message.TraceEvent) {
	//fmt.Printf("process event %s\n", evt.GetType())
	switch evt.GetType() {
	case message.TraceEvent_PUBLISH_MESSAGE:
		t.publish++
	case message.TraceEvent_REJECT_MESSAGE:
		t.reject++
	case message.TraceEvent_DUPLICATE_MESSAGE:
		t.duplicate++
	case message.TraceEvent_DELIVER_MESSAGE:
		t.deliver++
	case message.TraceEvent_ADD_PEER:
		t.add++
	case message.TraceEvent_REMOVE_PEER:
		t.remove++
	case message.TraceEvent_RECV_RPC:
		t.recv++
	case message.TraceEvent_SEND_RPC:
		t.send++
	case message.TraceEvent_DROP_RPC:
		t.drop++
	case message.TraceEvent_JOIN:
		t.join++
	case message.TraceEvent_LEAVE:
		t.leave++
	case message.TraceEvent_GRAFT:
		t.graft++
	case message.TraceEvent_PRUNE:
		t.prune++
	}
}

func (ts *traceStats) check(t *testing.T) {
	if ts.publish == 0 {
		t.Fatal("expected non-zero count")
	}
	if ts.duplicate == 0 {
		t.Fatal("expected non-zero count")
	}
	if ts.deliver == 0 {
		t.Fatal("expected non-zero count")
	}
	if ts.reject == 0 {
		t.Fatal("expected non-zero count")
	}
	if ts.add == 0 {
		t.Fatal("expected non-zero count")
	}
	if ts.recv == 0 {
		t.Fatal("expected non-zero count")
	}
	if ts.send == 0 {
		t.Fatal("expected non-zero count")
	}
	if ts.join == 0 {
		t.Fatal("expected non-zero count")
	}
	if ts.leave == 0 {
		t.Fatal("expected non-zero count")
	}
	if ts.graft == 0 {
		t.Fatal("expected non-zero count")
	}
	if ts.prune == 0 {
		t.Fatal("expected non-zero count")
	}
}

func TestJSONTracer(t *testing.T) {
	tracer, err := NewJSONTracer("D:\\trace.out.json")
	if err != nil {
		t.Fatal(err)
	}

	testWithTracer(t, tracer)
	time.Sleep(time.Second)
	tracer.Close()

	var stats traceStats
	var evt message.TraceEvent

	f, err := os.Open("D:\\trace.out.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	for {
		evt.Reset()
		err := dec.Decode(&evt)
		if err != nil {
			break
		}

		stats.process(&evt)
	}

	stats.check(t)
}

//
//func TestPBTracer(t *testing.T) {
//	tracer, err := NewPBTracer("/tmp/trace.out.pb")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	testWithTracer(t, tracer)
//	time.Sleep(time.Second)
//	tracer.Close()
//
//	var stats traceStats
//	var evt message.TraceEvent
//
//	f, err := os.Open("/tmp/trace.out.pb")
//	if err != nil {
//		t.Fatal(err)
//	}
//	defer f.Close()
//
//	r := protoio.NewDelimitedReader(f, 1<<20)
//	for {
//		evt.Reset()
//		err := r.ReadMsg(&evt)
//		if err != nil {
//			break
//		}
//
//		stats.process(&evt)
//	}
//
//	stats.check(t)
//}

//type mockRemoteTracer struct {
//	mx sync.Mutex
//	ts traceStats
//}
//
//func (mrt *mockRemoteTracer) handleStream(s network.Stream) {
//	defer s.Close()
//
//	gzr, err := gzip.NewReader(s)
//	if err != nil {
//		panic(err)
//	}
//
//	r := protoio.NewDelimitedReader(gzr, 1<<24)
//
//	var batch message.TraceEventBatch
//	for {
//		batch.Reset()
//		err := r.ReadMsg(&batch)
//		if err != nil {
//			if err != io.EOF {
//				s.Reset()
//			}
//			return
//		}
//
//		mrt.mx.Lock()
//		for _, evt := range batch.GetBatch() {
//			mrt.ts.process(evt)
//		}
//		mrt.mx.Unlock()
//	}
//}
//
//func (mrt *mockRemoteTracer) check(t *testing.T) {
//	mrt.mx.Lock()
//	defer mrt.mx.Unlock()
//	mrt.ts.check(t)
//}
//
//func TestRemoteTracer(t *testing.T) {
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//
//	h1 := NewTestHost()
//	h2 := NewTestHost()
//
//	mrt := &mockRemoteTracer{}
//	h1.SetStreamHandler(RemoteTracerProtoID, mrt.handleStream)
//
//	tracer, err := NewRemoteTracer(ctx, h2, peer.AddrInfo{ID: h1.ID(), Addrs: h1.Addrs()})
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	testWithTracer(t, tracer)
//	time.Sleep(time.Second)
//	tracer.Close()
//
//	mrt.check(t)
//}
