package gobayeux

import (
	"fmt"
	"sync"
)

type subscriptionsMap struct {
	lock sync.RWMutex
	subs map[Channel]chan []Message
}

func newSubscriptionsMap() *subscriptionsMap {
	return &subscriptionsMap{subs: make(map[Channel]chan []Message)}
}

func (sm *subscriptionsMap) Add(channel Channel, ms chan []Message) error {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	if _, ok := sm.subs[channel]; !ok {
		sm.subs[channel] = ms
		return nil
	}
	return fmt.Errorf("channel '%s' already subscribed", channel)
}

func (sm *subscriptionsMap) Remove(channel Channel) {
	sm.lock.Lock()
	defer sm.lock.Unlock()
	delete(sm.subs, channel)
}

func (sm *subscriptionsMap) Get(channel Channel) (chan []Message, error) {
	sm.lock.RLock()
	defer sm.lock.RUnlock()
	ms, ok := sm.subs[channel]
	if !ok {
		for key, ms := range sm.subs {
			if key.HasWildcard() && key.Match(channel){
				return ms, nil
			}
		}
		return nil, fmt.Errorf("channel '%s' has no subscriptions", channel)
	}
	return ms, nil
}
