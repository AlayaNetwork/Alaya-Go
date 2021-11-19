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

package network

import (
	ctypes "github.com/AlayaNetwork/Alaya-Go/consensus/cbft/types"
	"github.com/AlayaNetwork/Alaya-Go/crypto"
	"github.com/AlayaNetwork/Alaya-Go/node"
	"github.com/AlayaNetwork/Alaya-Go/p2p"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/params"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func makePubSub(handlerMsg func(p *peer, msg *p2p.Msg) error) *PubSub {
	sk, err := crypto.GenerateKey()
	if nil != err {
		panic(err)
	}
	p2pServer := &p2p.Server{Config: node.DefaultConfig.P2P}
	p2pServer.PrivateKey = sk
	p2pServer.NoDiscovery = true
	p2pServer.ListenAddr = ""
	if err := p2pServer.Start(); err != nil {
		panic(err)
	}
	localNode := enode.NewV4(&sk.PublicKey, nil, 0, 0)
	psServer := p2p.NewPubSubServer(localNode, p2pServer)
	pubSub := NewPubSub(psServer)

	pubSub.Init(ctypes.Config{Sys: params.AlayaChainConfig.Cbft, Option: nil}, func(id string) (p *peer, err error) {
		return nil, nil
	}, handlerMsg)
	return pubSub
}

type TestPSRW struct {
	writeMsgChan chan p2p.Msg
	readMsgChan  chan p2p.Msg
}

func (rw *TestPSRW) ReadMsg() (p2p.Msg, error) {
	return <-rw.readMsgChan, nil
}

func (rw *TestPSRW) WriteMsg(msg p2p.Msg) error {
	rw.writeMsgChan <- msg
	return nil
}

type TestMsg struct {
	Title string
}

func TestPubSubPublish(t *testing.T) {
	expect := []*TestMsg{
		{
			Title: "n1_msg",
		},
		{
			Title: "n2_msg",
		},
	}

	wg := sync.WaitGroup{}
	wg.Add(len(expect))
	pubSub1 := makePubSub(func(p *peer, msg *p2p.Msg) error {
		var tm TestMsg
		err := msg.Decode(&tm)
		assert.Equal(t, tm.Title, expect[1].Title)
		wg.Done()
		t.Log("pubSub1 receive:", "data", tm, "err", err)
		return nil
	})
	pubSub2 := makePubSub(func(p *peer, msg *p2p.Msg) error {
		var tm TestMsg
		err := msg.Decode(&tm)
		assert.Equal(t, tm.Title, expect[0].Title)
		wg.Done()
		t.Log("pubSub2 receive:", "data", tm, "err", err)
		return nil
	})

	n1Chan := make(chan p2p.Msg)
	n2Chan := make(chan p2p.Msg)
	trw1 := &TestPSRW{
		writeMsgChan: n1Chan,
		readMsgChan:  n2Chan,
	}
	trw2 := &TestPSRW{
		writeMsgChan: n2Chan,
		readMsgChan:  n1Chan,
	}

	// connect peer
	// peer n1
	go func() {
		newPeer := p2p.NewPeer(pubSub2.pss.Host().ID().ID(), "n2", nil)
		t.Log("newPeer", "id", newPeer.ID().TerminalString(), "localId", pubSub1.pss.Host().ID().ID().TerminalString(), "name", "n1")
		if err := pubSub1.handler(newPeer, trw1); err != nil {
			t.Fatal(err)
		}
	}()

	// peer n2
	go func() {
		newPeer := p2p.NewPeer(pubSub1.pss.Host().ID().ID(), "n1", nil)
		t.Log("newPeer", "id", newPeer.ID().TerminalString(), "localId", pubSub2.pss.Host().ID().ID().TerminalString(), "name", "n2")
		if err := pubSub2.handler(newPeer, trw2); err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(time.Millisecond * 6)

	// Topics of interest for node registration.
	// Send messages between nodes
	topic := "test"
	go func() {
		if err := pubSub1.Subscribe(topic); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Millisecond * 3)
		pubSub1.Publish(topic, uint64(1), expect[0])
	}()
	go func() {
		if err := pubSub2.Subscribe(topic); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Millisecond * 3)
		pubSub2.Publish(topic, uint64(2), expect[1])
	}()
	wg.Wait()
	pubSub1.Cancel(topic)
	pubSub2.Cancel(topic)
}

func TestPubSubPublish_DifferentTopics(t *testing.T) {
	expect := []*TestMsg{
		{
			Title: "n1_msg",
		},
		{
			Title: "n2_msg",
		},
		{
			Title: "n3_msg",
		},
	}

	wg := sync.WaitGroup{}
	wg.Add(len(expect))
	pubSub1 := makePubSub(func(p *peer, msg *p2p.Msg) error {
		var tm TestMsg
		err := msg.Decode(&tm)
		if tm.Title != expect[1].Title && tm.Title != expect[2].Title {
			t.Fatal("Expected value not received")
		}
		wg.Done()
		t.Log("pubSub1 receive:", "data", tm, "err", err)
		return nil
	})
	pubSub2 := makePubSub(func(p *peer, msg *p2p.Msg) error {
		var tm TestMsg
		err := msg.Decode(&tm)
		assert.Equal(t, tm.Title, expect[0].Title)
		wg.Done()
		t.Log("pubSub2 receive:", "data", tm, "err", err)
		return nil
	})
	pubSub3 := makePubSub(func(p *peer, msg *p2p.Msg) error {
		var tm TestMsg
		err := msg.Decode(&tm)
		assert.Equal(t, tm.Title, expect[0].Title)
		wg.Done()
		t.Log("pubSub3 receive:", "data", tm, "err", err)
		return nil
	})

	n1_2Chan := make(chan p2p.Msg)
	n2_1Chan := make(chan p2p.Msg)
	trw1 := &TestPSRW{
		writeMsgChan: n1_2Chan,
		readMsgChan:  n2_1Chan,
	}
	trw2 := &TestPSRW{
		writeMsgChan: n2_1Chan,
		readMsgChan:  n1_2Chan,
	}

	n1_3Chan := make(chan p2p.Msg)
	n3_1Chan := make(chan p2p.Msg)
	trw1_3 := &TestPSRW{
		writeMsgChan: n1_3Chan,
		readMsgChan:  n3_1Chan,
	}
	trw3_1 := &TestPSRW{
		writeMsgChan: n3_1Chan,
		readMsgChan:  n1_3Chan,
	}

	// connect peer
	// peer1 <-> peer2
	go func() {
		newPeer := p2p.NewPeer(pubSub2.pss.Host().ID().ID(), "n2", nil)
		t.Log("newPeer", "id", newPeer.ID().TerminalString(), "localId", pubSub1.pss.Host().ID().ID().TerminalString())
		if err := pubSub1.handler(newPeer, trw1); err != nil {
			t.Fatal(err)
		}
	}()
	go func() {
		newPeer := p2p.NewPeer(pubSub1.pss.Host().ID().ID(), "n1", nil)
		t.Log("newPeer", "id", newPeer.ID().TerminalString(), "localId", pubSub2.pss.Host().ID().ID().TerminalString())
		if err := pubSub2.handler(newPeer, trw2); err != nil {
			t.Fatal(err)
		}
	}()

	// peer1 <-> peer3
	go func() {
		newPeer := p2p.NewPeer(pubSub3.pss.Host().ID().ID(), "n3", nil)
		t.Log("newPeer", "id", newPeer.ID().TerminalString(), "localId", pubSub1.pss.Host().ID().ID().TerminalString())
		if err := pubSub1.handler(newPeer, trw1_3); err != nil {
			t.Fatal(err)
		}
	}()
	go func() {
		newPeer := p2p.NewPeer(pubSub1.pss.Host().ID().ID(), "n1", nil)
		t.Log("newPeer", "id", newPeer.ID().TerminalString(), "localId", pubSub3.pss.Host().ID().ID().TerminalString())
		if err := pubSub3.handler(newPeer, trw3_1); err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(time.Millisecond * 6)

	topic1 := "test1"
	topic2 := "test2"
	go func() {
		if err := pubSub1.Subscribe(topic1); err != nil {
			t.Fatal(err)
		}
		if err := pubSub1.Subscribe(topic2); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Millisecond * 3)
		pubSub1.Publish(topic1, uint64(1), expect[0])
	}()
	go func() {
		if err := pubSub2.Subscribe(topic1); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Millisecond * 3)
		pubSub2.Publish(topic1, uint64(2), expect[1])
	}()
	go func() {
		if err := pubSub3.Subscribe(topic2); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Millisecond * 3)
		pubSub3.Publish(topic2, uint64(3), expect[2])
	}()
	wg.Wait()
	pubSub1.Cancel(topic1)
	pubSub1.Cancel(topic2)
	pubSub2.Cancel(topic1)
	pubSub3.Cancel(topic2)
}

func TestPubSubPublish_ForwardMessage(t *testing.T) {
	expect := []*TestMsg{
		{
			Title: "n1_msg",
		},
		{
			Title: "n2_msg",
		},
		{
			Title: "n3_msg",
		},
	}

	wg := sync.WaitGroup{}
	wg.Add(6)
	pubSub1 := makePubSub(func(p *peer, msg *p2p.Msg) error {
		var tm TestMsg
		err := msg.Decode(&tm)
		if tm.Title != expect[1].Title && tm.Title != expect[2].Title {
			t.Fatal("Expected value not received")
		}
		t.Log("pubSub1 receive:", "data", tm, "err", err)
		wg.Done()
		return nil
	})
	pubSub2 := makePubSub(func(p *peer, msg *p2p.Msg) error {
		var tm TestMsg
		err := msg.Decode(&tm)
		if tm.Title != expect[0].Title && tm.Title != expect[2].Title {
			t.Fatal("Expected value not received")
		}
		t.Log("pubSub2 receive:", "data", tm, "err", err)
		wg.Done()
		return nil
	})
	pubSub3 := makePubSub(func(p *peer, msg *p2p.Msg) error {
		var tm TestMsg
		err := msg.Decode(&tm)
		if tm.Title != expect[0].Title && tm.Title != expect[1].Title {
			t.Fatal("Expected value not received")
		}
		t.Log("pubSub3 receive:", "data", tm, "err", err)
		wg.Done()
		return nil
	})

	n1_2Chan := make(chan p2p.Msg)
	n2_1Chan := make(chan p2p.Msg)
	trw1 := &TestPSRW{
		writeMsgChan: n1_2Chan,
		readMsgChan:  n2_1Chan,
	}
	trw2 := &TestPSRW{
		writeMsgChan: n2_1Chan,
		readMsgChan:  n1_2Chan,
	}

	n2_3Chan := make(chan p2p.Msg)
	n3_2Chan := make(chan p2p.Msg)
	trw2_3 := &TestPSRW{
		writeMsgChan: n2_3Chan,
		readMsgChan:  n3_2Chan,
	}
	trw3_2 := &TestPSRW{
		writeMsgChan: n3_2Chan,
		readMsgChan:  n2_3Chan,
	}

	// connect peer
	// peer1 <-> peer2
	go func() {
		newPeer := p2p.NewPeer(pubSub2.pss.Host().ID().ID(), "n2", nil)
		t.Log("newPeer", "id", newPeer.ID().TerminalString(), "localId", pubSub1.pss.Host().ID().ID().TerminalString(), "n1")
		if err := pubSub1.handler(newPeer, trw1); err != nil {
			t.Fatal(err)
		}
	}()
	go func() {
		newPeer := p2p.NewPeer(pubSub1.pss.Host().ID().ID(), "n1", nil)
		t.Log("newPeer", "id", newPeer.ID().TerminalString(), "localId", pubSub2.pss.Host().ID().ID().TerminalString(), "n2")
		if err := pubSub2.handler(newPeer, trw2); err != nil {
			t.Fatal(err)
		}
	}()

	// peer2 <-> peer3
	go func() {
		newPeer := p2p.NewPeer(pubSub3.pss.Host().ID().ID(), "n3", nil)
		t.Log("newPeer", "id", newPeer.ID().TerminalString(), "localId", pubSub2.pss.Host().ID().ID().TerminalString(), "n2")
		if err := pubSub2.handler(newPeer, trw2_3); err != nil {
			t.Fatal(err)
		}
	}()
	go func() {
		newPeer := p2p.NewPeer(pubSub2.pss.Host().ID().ID(), "n2", nil)
		t.Log("newPeer", "id", newPeer.ID().TerminalString(), "localId", pubSub3.pss.Host().ID().ID().TerminalString(), "n3")
		if err := pubSub3.handler(newPeer, trw3_2); err != nil {
			t.Fatal(err)
		}
	}()
	time.Sleep(time.Millisecond * 6)

	topic := "test"
	go func() {
		if err := pubSub1.Subscribe(topic); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Millisecond * 3)
		pubSub1.Publish(topic, uint64(1), expect[0])
	}()
	go func() {
		if err := pubSub2.Subscribe(topic); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Millisecond * 3)
		pubSub2.Publish(topic, uint64(2), expect[1])
	}()
	go func() {
		if err := pubSub3.Subscribe(topic); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Millisecond * 3)
		pubSub3.Publish(topic, uint64(3), expect[2])
	}()
	wg.Wait()
	pubSub1.Cancel(topic)
	pubSub2.Cancel(topic)
	pubSub3.Cancel(topic)
}
