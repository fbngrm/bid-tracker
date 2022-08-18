package item

import (
	"context"
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

// Note, the mutex must not be copied after first use. In other words, use
// pointer receivers for itemStore methods. itemStore is encapsulated on
// package boundary and must be accessed via a Service.
type itemStore struct {
	// a read/write mutex is favourable over mutex in this scenario since
	// it can be held by an arbitrary number of readers or a single writer.
	// whereas a "regular" mutex can be held by a single reader or writer only
	sync.RWMutex
	items map[uuid.UUID]*Item // fixme, store is never copied so it does't make sense to use pointer
}

// Idempotent, already registered items are ignored.
func (is *itemStore) register(ctx context.Context, i *Item) error {
	is.Lock()
	defer is.Unlock()

	if i == nil {
		return errors.New("could not register, item is nil")
	}

	if _, ok := is.items[i.ID]; ok {
		return nil
	}
	is.items[i.ID] = i
	return nil
}

// Note, for write operations on Item we could optimize by adding a mutex to the Item type
// and unlock the store immediately after reading an item. The Item would be locked for
// writes while the store could serve other write requests without being locked by waiting
// for the write to Item to finish. So that stores would be locked for writing only
// when a new item is registered, not while operating on an item.
// I consider it out of scope due to time constraints, even though it is a low hanging fruit
// for optimization.

// Errors if the bid's item is not registered.
func (is *itemStore) write(ctx context.Context, b *bid.Bid) error {
	is.Lock()
	defer is.Unlock()

	if b == nil {
		return errors.New("could not bid, bid is nil")
	}

	item, ok := is.items[b.ID]
	if !ok {
		return fmt.Errorf("could not read item for bid [%s], not registered [%s]", b.ID, b.ItemID)
	}

	// not concurrency safe but we locked the store. here we could lock item and unlock store
	err := item.addBid(ctx, b)
	if err != nil {
		return fmt.Errorf("could not bid [%s]: %w", b.ID, err)
	}

	return nil
}

func (is *itemStore) readHighest(ctx context.Context, itemID uuid.UUID) (*bid.Bid, error) {
	is.RLock()         // we can have concurrent reads but no writes
	defer is.RUnlock() // unlock for writes

	i, ok := is.items[itemID]
	if !ok {
		return nil, fmt.Errorf("could not read item, not registered [%s]", itemID)
	}

	if i == nil {
		return nil, fmt.Errorf("could not get highest bid for item [%s], item in store is nil", itemID.String())
	}

	return i.highest, nil
}
