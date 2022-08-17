package events

import (
	"errors"
	"fmt"
	"sync"

	"github.com/fbngrm/bid-tracker/pkg/bid"
	"github.com/fbngrm/bid-tracker/pkg/item"
	"github.com/google/uuid"
)

var eventStore *EventStore

// We use pessimistic locking as a concurrency model for our event store.
// Since we expect many concurrent and conflicting writes (bids), we cannot
// guarantee data consistency using an optimistic locking approach.
// For performance reasons, we choose locking over channel synchronization.

// Note, the mutex must not be copied after first use.
// In other words, use pointer receivers for EventStore methods.
type EventStore struct {
	// a read/write mutex is favourable over mutex in this scenario since
	// it can be held by an arbitrary number of readers or a single writer
	// whereas a "regular" mutex can be held by a single reader or writer only
	sync.RWMutex
	events map[uuid.UUID]*item.Item
}

// We want to make sure only a single instance of EventStore exists.
func NewStore() *EventStore {
	if eventStore == nil {
		eventStore = &EventStore{}
	}
	return eventStore
}

// Idempotent, already registered items are ignored.
func (es *EventStore) Register(i *item.Item) error {
	es.Lock()
	defer es.Unlock()

	if i == nil {
		return errors.New("could not register, item is nil")
	}

	if _, ok := es.events[i.ID]; ok {
		return nil
	}
	es.events[i.ID] = i
	return nil
}

// Fails if the bid's item is not registered.
func (es *EventStore) Write(b *bid.Bid) error {
	es.Lock()
	defer es.Unlock()

	if b == nil {
		return errors.New("could not bid, bid is nil")
	}

	item, ok := es.events[b.ID]
	if !ok {
		return fmt.Errorf("could not bid [%s]: item not registered [%s]", b.ID, b.ItemID)
	}

	// not concurrency safe but we locked the store
	item.Bid(b)

	return nil
}
