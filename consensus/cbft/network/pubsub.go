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
	"fmt"
	"github.com/AlayaNetwork/Alaya-Go/p2p/pubsub"
)

type PubSub struct {
	pubsub pubsub.PubSubServer

	topicChan	chan string
	topics		[]string
	// The set of topics we are subscribed to
	mySubs		map[string]map[*pubsub.Subscription]struct{}
}

//Subscribe subscribe a topic
func (ps *PubSub) Subscribe(topic string) error {
	if err := ps.pubsub.PublishMsg(topic); err != nil {
		return err
	}
	ps.topics = append(ps.topics, topic)
	return nil
}

//UnSubscribe a topic
func (ps *PubSub) Cancel(topic string) error {
	sm := ps.mySubs[topic]
	if sm != nil {
		for s, _ := range sm {
			if s != nil {
				s.Cancel()
			}
		}
		return nil
	}

	return fmt.Errorf("Can not find","topic", topic)
}