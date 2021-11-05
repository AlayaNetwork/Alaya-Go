package pubsub

import (
	"context"
	"github.com/AlayaNetwork/Alaya-Go/p2p"
	"time"

	"github.com/libp2p/go-libp2p-core/connmgr"

	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"
)

// Connectedness signals the capacity for a connection with a given node.
// It is used to signal to services and other peers whether a node is reachable.
type Connectedness int

const (
	// NotConnected means no connection to peer, and no extra information (default)
	NotConnected Connectedness = iota

	// Connected means has an open, live connection to peer
	Connected

	// CanConnect means recently connected to peer, terminated gracefully
	CanConnect

	// CannotConnect means recently attempted connecting but failed to connect.
	// (should signal "made effort, failed")
	CannotConnect
)

type ProtocolID string

// Host is an object participating in a p2p network, which
// implements protocols or provides services. It handles
// requests like a Server, and issues requests like a Client.
// It is called Host because it is both Server and Client (and Peer
// may be confusing).
type Host interface {
	// ID returns the (local) peer.ID associated with this Host
	ID() *enode.Node

	// Peerstore returns the Host's repository of Peer Addresses and Keys.
	Peerstore() Peerstore

	// Networks returns the Network interface of the Host
	Network() Network

	// Connect ensures there is a connection between this host and the peer with
	// given peer.ID. Connect will absorb the addresses in pi into its internal
	// peerstore. If there is not an active connection, Connect will issue a
	// h.Network.Dial, and block until a connection is open, or an error is
	// returned. // TODO: Relay + NAT.
	// addConsensusNode
	Connect(ctx context.Context, pi enode.ID) error

	// SetStreamHandler sets the protocol handler on the Host's Mux.
	// This is equivalent to:
	//   host.Mux().SetHandler(proto, handler)
	// (Threadsafe)
	SetStreamHandler(pid ProtocolID, handler StreamHandler)

	// SetStreamHandlerMatch sets the protocol handler on the Host's Mux
	// using a matching function for protocol selection.
	SetStreamHandlerMatch(ProtocolID, func(string) bool, StreamHandler)

	// RemoveStreamHandler removes a handler on the mux that was set by
	// SetStreamHandler
	RemoveStreamHandler(pid ProtocolID)

	// NewStream opens a new stream to given peer p, and writes a p2p/protocol
	// header with given ProtocolID. If there is no connection to p, attempts
	// to create one. If ProtocolID is "", writes no header.
	// (Threadsafe)
	NewStream(ctx context.Context, p enode.ID, pids ...ProtocolID) (Stream, error)

	// Close shuts down the host, its Network, and services.
	Close() error

	// ConnManager returns this hosts connection manager
	ConnManager() connmgr.ConnManager
}

// StreamHandler is the type of function used to listen for
// streams opened by the remote side.
type StreamHandler func(Stream)

// Network is the interface used to connect to the outside world.
// It dials and listens for connections. it uses a Swarm to pool
// connections (see swarm pkg, and peerstream.Swarm). Connections
// are encrypted with a TLS-like protocol.
type Network interface {
	// ConnsToPeer returns the connections in this Netowrk for given peer.
	ConnsToPeer(p enode.ID) []Conn

	// Connectedness returns a state signaling connection capabilities
	Connectedness(enode.ID) Connectedness

	// Notify/StopNotify register and unregister a notifiee for signals
	Notify(Notifiee)

	// Peers returns the peers connected
	Peers() []enode.ID
}

type Notifiee interface {
	Connected(Network, Conn) // called when a connection opened
}

// Stream represents a bidirectional channel between two agents in
// a libp2p network. "agent" is as granular as desired, potentially
// being a "request -> reply" pair, or whole protocols.
//
// Streams are backed by a multiplexer underneath the hood.
type Stream interface {
	Protocol() ProtocolID

	// Conn returns the connection this stream is part of.
	Conn() Conn

	ReadWriter() p2p.MsgReadWriter //
}

// Conn is a connection to a remote peer. It multiplexes streams.
// Usually there is no need to use a Conn directly, but it may
// be useful to get information about the peer on the other side:
//  stream.Conn().RemotePeer()
type Conn interface {
	// ID returns an identifier that uniquely identifies this Conn within this
	// host, during this run. Connection IDs may repeat across restarts.
	ID() string

	// GetStreams returns all open streams over this conn.
	GetStreams() []Stream

	// Stat stores metadata pertaining to this conn.
	Stat() Stat

	// RemotePeer returns the peer ID of the remote peer.
	RemotePeer() *enode.Node
}

// Stat stores metadata pertaining to a given Stream/Conn.
type Stat struct {
	// Direction specifies whether this is an inbound or an outbound connection.
	Direction Direction
	// Opened is the timestamp when this connection was opened.
	Opened time.Time
	// Transient indicates that this connection is transient and may be closed soon.
	Transient bool
	// Extra stores additional metadata about this connection.
	Extra map[interface{}]interface{}
}

// Direction represents which peer in a stream initiated a connection.
type Direction int

const (
	// DirUnknown is the default direction.
	DirUnknown Direction = iota
	// DirInbound is for when the remote peer initiated a connection.
	DirInbound
	// DirOutbound is for when the local peer initiated a connection.
	DirOutbound
)