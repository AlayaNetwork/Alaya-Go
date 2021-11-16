package p2p

import (
	"context"
	"sync"
	"time"

	"github.com/AlayaNetwork/Alaya-Go/log"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
	"github.com/AlayaNetwork/Alaya-Go/p2p/enr"
)

// SubscribeTopic to  the given topic.A given validator and subscription handler is
// used to handle messages from the subnet. The base protobuf message is used to initialize new messages for decoding.
func (srv *Server) SubscribeTopic(ctx context.Context, topic string) {

	ticker := time.NewTicker(time.Second * 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				if srv.running {
					continue
				}

				// Check   there are enough peers
				if !srv.validPeersExist(topic) {
					log.Debug("No peers found subscribed  gossip topic . Searching network for peers subscribed to the topic.", "topic", topic)
					_, err := srv.FindPeersWithTopic(
						ctx,
						topic,
						srv.Config.MinimumPeersPerTopic,
					)
					if err != nil {
						log.Error("Could not search for peers", "err", err)
						return
					}
				}

			}
		}
	}()
}

// find if we have peers who are subscribed to the same subnet
func (srv *Server) validPeersExist(subnetTopic string) bool {
	numOfPeers := srv.pubSubServer.pubsub.ListPeers(subnetTopic)
	return len(numOfPeers) >= srv.Config.MinimumPeersPerTopic
}

/*
func (srv *Server) topicToIdx(topic string) uint64 {
	return 0
}*/

// FindPeersWithTopic performs a network search for peers
// subscribed to a particular subnet. Then we try to connect
// with those peers. This method will block until the required amount of
// peers are found, the method only exits in the event of context timeouts.
func (srv *Server) FindPeersWithTopic(ctx context.Context, topic string, threshold int) (bool, error) {

	if srv.DiscV5 == nil {
		// return if discovery isn't set
		return false, nil
	}

	//topic += s.Encoding().ProtocolSuffix()
	iterator := srv.DiscV5.RandomNodes()
	iterator = filterNodes(ctx, iterator, srv.filterPeerForTopic(topic))

	currNum := len(srv.pubSubServer.pubsub.ListPeers(topic))
	wg := new(sync.WaitGroup)
	for {
		if err := ctx.Err(); err != nil {
			return false, err
		}
		if currNum >= threshold {
			break
		}
		nodes := enode.ReadNodes(iterator, int(srv.Config.MinimumPeersInTopicSearch))

		for i, _ := range nodes {
			wg.Add(1)
			srv.AddConsensusPeerWithDone(nodes[i], func() {
				wg.Done()
			})
		}
		// Wait for all dials to be completed.
		wg.Wait()
		currNum = len(srv.pubSubServer.pubsub.ListPeers(topic))
	}
	return true, nil
}

// returns a method with filters peers specifically for a particular attestation subnet.
func (srv *Server) filterPeerForTopic(topic string) func(node *enode.Node) bool {
	return func(node *enode.Node) bool {
		if !srv.filterPeer(node) {
			return false
		}
		srv.topicSubscriberMu.RLock()
		defer srv.topicSubscriberMu.RUnlock()
		nodes, ok := srv.topicSubscriber[topic]
		if !ok {
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
	// do not dial nodes with their tcp ports not set
	if err := node.Record().Load(enr.WithEntry("tcp", new(enr.TCP))); err != nil {
		if !enr.IsNotFound(err) {
			log.Error("Could not retrieve tcp port", err)
		}
		return false
	}
	return true
}

/*
func (srv *Server) RefreshTopicENR() {
	// return early if discv5 isnt running
	if srv.DiscV5 == nil {
		return
	}
	bitTopics := bitfield.NewBitvector64()
	committees := srv.GetAllTopics()
	for _, idx := range committees {
		bitTopics.SetBitAt(idx, true)
	}
	currentBitTopics, err := getBitvector(srv.Config.TopicKey, srv.DiscV5.Self().Record())
	if err != nil {
		log.Error("Could not retrieve att bitfield", "err", err)
		return
	}

	if bytes.Equal(bitTopics, currentBitTopics) {
		// return early if bitfields haven't changed
		return
	}
	srv.updateTopicRecord(bitTopics)
	// ping all peers to inform them of new metadata
	//todo not sure ping here
	//s.pingPeers()
}

func (srv *Server) updateTopicRecord(bitVAtt bitfield.Bitvector64) {
	entry := enr.WithEntry(srv.Config.TopicKey, &bitVAtt)
	srv.DiscV5.LocalNode().Set(entry)
}

func (srv *Server) GetAllTopics() []uint64 {
	return nil
}

// Reads the attestation subnets entry from a node's ENR and determines
// the committee indices of the attestation subnets the node is subscribed to.
func (srv *Server) getTopics(record *enr.Record) ([]uint64, error) {
	bitV, err := getBitvector(srv.Config.TopicKey, record)
	if err != nil {
		return nil, err
	}
	if len(bitV) != byteCount(int(srv.Config.TopicCount)) {
		return []uint64{}, errors.Errorf("invalid bitvector provided, it has a size of %d", len(bitV))
	}
	var topicIdxs []uint64
	for i := uint64(0); i < srv.Config.TopicCount; i++ {
		if bitV.BitAt(i) {
			topicIdxs = append(topicIdxs, i)
		}
	}
	return topicIdxs, nil
}

// Parses the attestation subnets ENR entry in a node and extracts its value
// as a bitvector for further manipulation.
func getBitvector(key string, record *enr.Record) (bitfield.Bitvector64, error) {
	bitV := bitfield.NewBitvector64()
	entry := enr.WithEntry(key, &bitV)
	err := record.Load(entry)
	if err != nil {
		return nil, err
	}
	return bitV, nil
}

// Determines the number of bytes that are used
// to represent the provided number of bits.
func byteCount(bitCount int) int {
	numOfBytes := bitCount / 8
	if bitCount%8 != 0 {
		numOfBytes++
	}
	return numOfBytes
}*/
