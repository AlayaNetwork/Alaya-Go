package p2p

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/AlayaNetwork/Alaya-Go/log"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
)

// DiscoverTopic to  the given topic.A given validator and subscription handler is
// used to handle messages from the subnet. The base protobuf message is used to initialize new messages for decoding.
func (srv *Server) DiscoverTopic(ctx context.Context, topic string) {

	ticker := time.NewTicker(time.Second * 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				if !srv.running {
					continue
				}
				// Check   there are enough peers
				if srv.topicWithPubsub(topic) {
					peers := srv.pubSubServer.PubSub().ListPeers(topic)
					peerNeedFind := srv.Config.MinimumPeersPerTopic - len(peers)
					if peerNeedFind <= 0 {
						continue
					}
					copyNodes, err := srv.getNotConnectNode(topic)
					if err != nil {
						log.Error("discover topic fail", "err", err)
						return
					}

					if len(copyNodes) == 0 {
						log.Debug("all peers are found,no need searching network", "topic", topic, "peers", len(peers))
						continue
					}
					if peerNeedFind > len(copyNodes) {
						peerNeedFind = len(copyNodes)
					}
					log.Debug("not enough nodes in this topic,searching network", "topic", topic, "peers", len(peers), "remainNodes", len(copyNodes), "peerNeedFind", peerNeedFind, "minimumPeersPerTopic", srv.Config.MinimumPeersPerTopic)
					if err := srv.FindPeersWithTopic(ctx, topic, copyNodes, peerNeedFind); err != nil {
						log.Debug("Could not search for peers", "err", err)
						return
					}
				} else {
					copyNodes, err := srv.getNotConnectNode(topic)
					if err != nil {
						log.Error("discover topic fail", "err", err)
						return
					}
					if len(copyNodes) == 0 {
						log.Debug("all peers are found,no need searching network", "topic", topic)
						continue
					}
					srv.topicSubscriberMu.RLock()
					nodes := srv.topicSubscriber[topic]
					srv.topicSubscriberMu.RUnlock()
					peerNeedFind := srv.MinimumPeersPerTopic - len(nodes) + len(copyNodes)
					if peerNeedFind <= 0 {
						continue
					}
					if peerNeedFind > len(copyNodes) {
						peerNeedFind = len(copyNodes)
					}
					log.Debug("not enough nodes in this topic,searching network", "topic", topic, "remainNodes", len(copyNodes), "peerNeedFind", peerNeedFind, "minimumPeersPerTopic", srv.Config.MinimumPeersPerTopic)
					if err := srv.FindPeersWithTopic(ctx, topic, copyNodes, peerNeedFind); err != nil {
						log.Debug("Could not search for peers", "err", err)
						return
					}
				}
			}
		}
	}()
}

func (srv *Server) topicWithPubSub(topic string) bool {
	topics := srv.pubSubServer.PubSub().GetTopics()
	for _, s := range topics {
		if s == topic {
			return true
		}
	}
	return false
}

func (srv *Server) getNotConnectNode(topic string) ([]*enode.Node, error) {
	srv.topicSubscriberMu.RLock()
	nodes, ok := srv.topicSubscriber[topic]
	if !ok {
		srv.topicSubscriberMu.RUnlock()
		return nil, fmt.Errorf("the topic %s should discover can't find", topic)
	}
	copyNodes := make([]*enode.Node, len(nodes))
	copy(copyNodes, nodes)
	srv.topicSubscriberMu.RUnlock()

	currentConnectPeer := make(map[enode.ID]struct{})
	srv.doPeerOp(func(m map[enode.ID]*Peer) {
		for id, _ := range m {
			currentConnectPeer[id] = struct{}{}
		}
	})
	// 找到没有连接的节点
	for i := 0; i < len(copyNodes); {
		if copyNodes[i].ID() == srv.localnode.ID() {
			copyNodes = append(copyNodes[:i], copyNodes[i+1:]...)
			continue
		}
		if _, ok := currentConnectPeer[copyNodes[i].ID()]; ok {
			copyNodes = append(copyNodes[:i], copyNodes[i+1:]...)
			continue
		}
		i++
	}
	return copyNodes, nil
}

// FindPeersWithTopic performs a network search for peers
// subscribed to a particular subnet. Then we try to connect
// with those peers. This method will block until the required amount of
// peers are found, the method only exits in the event of context timeouts.
func (srv *Server) FindPeersWithTopic(ctx context.Context, topic string, nodes []*enode.Node, threshold int) error {

	if srv.ntab == nil {
		// return if discovery isn't set
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	wg := new(sync.WaitGroup)
	indexs := rand.New(rand.NewSource(time.Now().UnixNano())).Perm(len(nodes))

	if threshold > 6 {
		threshold = threshold / 2
	}
	var dialShouldReTry int

	//持续发现节点,并且过滤掉dial失败的节点
	for i := 0; i < len(indexs); i++ {
		wg.Add(1)
		srv.AddConsensusPeerWithDone(nodes[indexs[i]], func(err error) {
			if err != nil {
				dialShouldReTry++
			}
			wg.Done()
		})
		threshold--
		if threshold == 0 {
			wg.Wait()
			if dialShouldReTry > 0 {
				threshold += dialShouldReTry
				dialShouldReTry = 0
			} else {
				break
			}
		}
	}
	currNum := len(srv.pubSubServer.PubSub().ListPeers(topic))

	log.Trace("Searching network for peers subscribed to the topic done.", "topic", topic, "peers", currNum, "dialShouldReTry", threshold)

	return nil
}

// returns a method with filters peers specifically for a particular attestation subnet.
func (srv *Server) filterPeerForTopic(nodes []enode.ID) func(node *enode.Node) bool {
	return func(node *enode.Node) bool {
		if !srv.filterPeer(node) {
			return false
		}

		for _, peer := range nodes {
			if peer == node.ID() {
				return true
			}
		}
		return false
	}
}

// filterPeer validates each node that we retrieve from our dht. We
// try to ascertain that the peer can be a valid protocol peer.
// Validity Conditions:
// 1) The local node is still actively looking for peers to
//    connect to.
// 2) Peer has a valid IP and TCP port set in their enr.
// 3) Peer hasn't been marked as 'bad'
// 4) Peer is not currently active or connected.
// 5) Peer is ready to receive incoming connections.
// 6) Peer's fork digest in their ENR matches that of
// 	  our localnodes.
func (srv *Server) filterPeer(node *enode.Node) bool {
	// Ignore nil node entries passed in.
	if node == nil {
		return false
	}
	// ignore nodes with no ip address stored.
	if node.IP() == nil {
		return false
	}
	if node.ID() == srv.localnode.ID() {
		return false
	}
	// do not dial nodes with their tcp ports not set
	/*if err := node.Record().Load(enr.WithEntry("tcp", new(enr.TCP))); err != nil {
		if !enr.IsNotFound(err) {
			log.Error("Could not retrieve tcp port", err)
		}
		return false
	}*/
	return true
}
