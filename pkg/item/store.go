package item

import (
	"errors"
	"fmt"
	"sync"

	"github.com/fbngrm/bid-tracker/pkg/bid"
	"github.com/google/uuid"
)

// We use pessimistic locking as a concurrency strategy for the item store.
// Since we expect many concurrent and conflicting writes (bids), we could
// not guarantee data integrity using an optimistic locking approach.
// For performance reasons, we choose locking over channel synchronization.

// Note, the mutex must not be copied after first use.
// In other words, use pointer receivers for ItemStore methods.
type ItemStore struct {
	// a read/write mutex is favourable over mutex in this scenario since
	// it can be held by an arbitrary number of readers or a single writer.
	// whereas a "regular" mutex can be held by a single reader or writer only
	sync.RWMutex
	items map[uuid.UUID]*Item
}

// Idempotent, already registered items are ignored.
func (es *ItemStore) register(i *Item) error {
	es.Lock()
	defer es.Unlock()

	if i == nil {
		return errors.New("could not register, item is nil")
	}

	if _, ok := es.items[i.ID]; ok {
		return nil
	}
	es.items[i.ID] = i
	return nil
}

// Errors if the bid's item is not registered.
func (es *ItemStore) write(b *bid.Bid) error {
	es.Lock()
	defer es.Unlock()

	if b == nil {
		return errors.New("could not bid, bid is nil")
	}

	item, ok := es.items[b.ID]
	if !ok {
		return fmt.Errorf("could get item for bid [%s], not registered [%s]", b.ID, b.ItemID)
	}

	// not concurrency safe but we locked the store
	err := item.addBid(b)
	if err != nil {
		return fmt.Errorf("could not bid [%s]: %w", b.ID, err)
	}

	return nil
}
