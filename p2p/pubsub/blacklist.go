package pubsub

import (
	"sync"
	"time"

	"github.com/AlayaNetwork/Alaya-Go/p2p/enode"

	"github.com/whyrusleeping/timecache"
)

// Blacklist is an interface for peer blacklisting.
type Blacklist interface {
	Add(id enode.ID) bool
	Contains(enode.ID) bool
}

// MapBlacklist is a blacklist implementation using a perfect map
type MapBlacklist map[enode.ID]struct{}

// NewMapBlacklist creates a new MapBlacklist
func NewMapBlacklist() Blacklist {
	return MapBlacklist(make(map[enode.ID]struct{}))
}

func (b MapBlacklist) Add(p enode.ID) bool {
	b[p] = struct{}{}
	return true
}

func (b MapBlacklist) Contains(p enode.ID) bool {
	_, ok := b[p]
	return ok
}

// TimeCachedBlacklist is a blacklist implementation using a time cache
type TimeCachedBlacklist struct {
	sync.RWMutex
	tc *timecache.TimeCache
}

// NewTimeCachedBlacklist creates a new TimeCachedBlacklist with the given expiry duration
func NewTimeCachedBlacklist(expiry time.Duration) (Blacklist, error) {
	b := &TimeCachedBlacklist{tc: timecache.NewTimeCache(expiry)}
	return b, nil
}

// Add returns a bool saying whether Add of peer was successful
func (b *TimeCachedBlacklist) Add(p enode.ID) bool {
	b.Lock()
	defer b.Unlock()
	s := p.String()
	if b.tc.Has(s) {
		return false
	}
	b.tc.Add(s)
	return true
}

func (b *TimeCachedBlacklist) Contains(p enode.ID) bool {
	b.RLock()
	defer b.RUnlock()

	return b.tc.Has(p.String())
}
