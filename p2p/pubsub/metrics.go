package pubsub

import "github.com/AlayaNetwork/Alaya-Go/metrics"

var (
	notSubscribeCounter = metrics.NewRegisteredCounter("p2p/pubSub/notSub/count", nil)
	duplicateMessageCounter = metrics.NewRegisteredCounter("p2p/pubSub/duplicateMessage/count", nil)
	messageCounter = metrics.NewRegisteredCounter("p2p/pubSub/message/count", nil)
)
