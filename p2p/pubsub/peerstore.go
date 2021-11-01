package pubsub

import (
	"io"
	"math"
	"time"

	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
)

// Permanent TTLs (distinct so we can distinguish between them, constant as they
// are, in fact, permanent)
const (
	// PermanentAddrTTL is the ttl for a "permanent address" (e.g. bootstrap nodes).
	PermanentAddrTTL = math.MaxInt64 - iota

	// ConnectedAddrTTL is the ttl used for the addresses of a peer to whom
	// we're connected directly. This is basically permanent, as we will
	// clear them + re-add under a TempAddrTTL after disconnecting.
	ConnectedAddrTTL
)

// Peerstore provides a threadsafe store of Peer related
// information.
type Peerstore interface {
	io.Closer

	// AddAddrs gives this AddrBook addresses to use, with a given ttl
	// (time-to-live), after which the address is no longer valid.
	// If the manager has a longer TTL, the operation is a no-op for that address
	AddAddrs(p enode.ID, ttl time.Duration)

	// PrivKey returns the private key of a peer, if known. Generally this might only be our own
	// private key
	//PrivKey(id enode.ID) *ecdsa.PrivateKey
}
