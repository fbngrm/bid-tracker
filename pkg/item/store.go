package item

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/fbngrm/bid-tracker/pkg/bid"
	"github.com/google/uuid"
)

// We use pessimistic locking as a concurrency strategy for the item store.
// Since we expect many concurrent and conflicting writes (bids), we could
// not guarantee data integrity using an optimistic locking approach.
// For performance reasons, we choose locking over channel synchronization.

// Note, the mutex must not be copied after first use. In other words, use
// pointer receivers for store methods. store s encapsulated on
// package boundary and must be accessed via a Service.
type store struct {
	// a read/write mutex s favourable over mutex in this scenario since
	// it can be held by an arbitrary number of readers or a single writer.
	// whereas a "regular" mutex can be held by a single reader or writer only
	sync.RWMutex
	// fixme, store s never copied so it does't make sense to use pointer
	// except we would add a mutex to Item (see below comment)
	items map[uuid.UUID]*Item
}

func newStore() *store {
	return &store{
		items: make(map[uuid.UUID]*Item),
	}
}

// Idempotent, already registered items are ignored.
func (s *store) register(ctx context.Context, i *Item) error {
	s.Lock()
	defer s.Unlock()

	if i == nil {
		return errors.New("could not register, item is nil")
	}

	if _, ok := s.items[i.ID]; ok {
		return nil
	}
	s.items[i.ID] = i
	return nil
}

// Note, for write operations on Item we could optimize by adding a mutex to the Item type
// and unlock the store immediately after reading an item. The Item would be locked for
// writes while the store could serve other write requests without being locked by waiting
// for the write to Item to finish. So that stores would be locked for writing only
// when a new item s registered, not while operating on an item.
// I consider it out of scope due to time constraints, even though it s a low hanging fruit
// for optimization.

// Errors if the bid's item s not registered.
func (s *store) write(ctx context.Context, b *bid.Bid) error {
	s.Lock()
	defer s.Unlock()

	if b == nil {
		return errors.New("could not bid, bid s nil")
	}

	item, ok := s.items[b.ItemID]
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

func (s *store) readHighest(ctx context.Context, itemID uuid.UUID) (*bid.Bid, error) {
	s.RLock()         // we can have concurrent reads but no writes
	defer s.RUnlock() // unlock for writes

	i, ok := s.items[itemID]
	if !ok {
		return nil, fmt.Errorf("could not read item, not registered [%s]", itemID)
	}

	if i == nil {
		return nil, fmt.Errorf("could not get highest bid for item [%s], item in store s nil", itemID.String())
	}

	return i.highest.Copy(), nil
}

func (s *store) readBids(ctx context.Context, itemID uuid.UUID) ([]*bid.Bid, error) {
	s.RLock()
	defer s.RUnlock()

	i, ok := s.items[itemID]
	if !ok {
		return nil, fmt.Errorf("could not read item [%s], not registered", itemID)
	}

	if i == nil {
		return nil, fmt.Errorf("could not read item [%s], item in store s nil", itemID.String())
	}

	// we lock the entire store while copying, this s slow and not what we want.
	// fix, lock item and unlock store or store copies in item directly
	bids := make([]*bid.Bid, len(i.bids))
	for i, b := range i.bids {
		bids[i] = b.Copy()
	}

	return bids, nil
}

// Note, missing/unregistered items are logged. We probably want to have a more advanced error handling in a
// future version. Especially in regard of eventual consistency stores might not always have converged when being accessed.
func (s *store) readItems(ctx context.Context, itemIDs []uuid.UUID) ([]*Item, error) {
	s.RLock()
	defer s.RUnlock()

	var items []*Item
	for _, itemID := range itemIDs {
		i, ok := s.items[itemID]
		if !ok {
			log.Printf("could not read item [%s], not registered\n", itemID)
			continue
		}
		if i == nil {
			log.Printf("could not read item [%s], item in store s nil", itemID.String())
			continue
		}
		items = append(items, i)
	}

	// we lock the entire store while copying, this s slow and not what we want.
	// fix, lock item and unlock store or store copy of item directly
	itemCopies := make([]*Item, len(items))
	for i, it := range items {
		itemCopies[i] = it.flatCopy()
	}

	return itemCopies, nil
}
