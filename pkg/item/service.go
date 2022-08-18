package item

import (
	"context"
	"fmt"

	"github.com/fbngrm/bid-tracker/pkg/bid"
	"github.com/google/uuid"
)

// Service exposes access to the item store to the outside world.
type Service struct {
	// for now we keep a tight coupling here
	store *itemStore
}

func NewService() *Service {
	return &Service{
		store: &itemStore{},
	}
}

// Note, we must escape logging of potentially harmful user input using %q formatting directive.
func (s *Service) CreateItem(ctx context.Context, name string) (*Item, error) {
	i, err := newItem(name)
	if err != nil {
		return nil, fmt.Errorf("could not create new item [%q]: %w", name, err)
	}
	err = s.store.register(ctx, i)
	if err != nil {
		return nil, fmt.Errorf("could not register item [%q]: %w", name, err)
	}
	return i, nil
}

func (s *Service) PlaceBid(ctx context.Context, b *bid.Bid) error {
	err := s.store.write(ctx, b)
	if err != nil {
		return fmt.Errorf("could not write bid [%s] to item store: %w", b.ID.String(), err)
	}
	return nil
}

func (s *Service) GetHighestBidForItem(ctx context.Context, itemID uuid.UUID) (*bid.Bid, error) {
	b, err := s.store.readHighest(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("could not read item [%s] from store: %w", itemID.String(), err)
	}
	return b, nil
}

func (s *Service) GetAllBidsForItem(ctx context.Context, itemID uuid.UUID) ([]*bid.Bid, error) {
	bids, err := s.store.readBids(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("could not read bids for item [%s] from store: %w", itemID.String(), err)
	}
	return bids, nil
}

func (s *Service) GetItemsForBids(ctx context.Context, itemIDs []uuid.UUID) ([]*Item, error) {
	items, err := s.store.readItems(ctx, itemIDs)
	if err != nil {
		return nil, fmt.Errorf("could not read items from store: %w", err)
	}
	return items, nil
}
