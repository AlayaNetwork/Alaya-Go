package pubsub

import (
	"errors"

	"github.com/AlayaNetwork/Alaya-Go/p2p/pubsub/message"

	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
)

// ErrTooManySubscriptions may be returned by a SubscriptionFilter to signal that there are too many
// subscriptions to process.
var ErrTooManySubscriptions = errors.New("too many subscriptions")

// SubscriptionFilter is a function that tells us whether we are interested in allowing and tracking
// subscriptions for a given topic.
//
// The filter is consulted whenever a subscription notification is received by another peer; if the
// filter returns false, then the notification is ignored.
//
// The filter is also consulted when joining topics; if the filter returns false, then the Join
// operation will result in an error.
type SubscriptionFilter interface {
	// CanSubscribe returns true if the topic is of interest and we can subscribe to it
	CanSubscribe(topic string) bool

	// FilterIncomingSubscriptions is invoked for all RPCs containing subscription notifications.
	// It should filter only the subscriptions of interest and my return an error if (for instance)
	// there are too many subscriptions.
	FilterIncomingSubscriptions(enode.ID, []*message.RPC_SubOpts) ([]*message.RPC_SubOpts, error)
}
