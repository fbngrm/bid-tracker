package item

import (
	"context"
	"errors"
	"fmt"

	"github.com/fbngrm/bid-tracker/pkg/bid"
	"github.com/google/uuid"
)

// Note, Item is not safe for concurrent use and is intended to be used via a store only.
// Thus, we encapsulate fields and methods on package boundary, that need to be protected.
// TODO: provide serialziation accessible from outside the package.
type Item struct {
	ID      uuid.UUID // uuid RFC 4122
	name    string
	highest *bid.Bid
	bids    []*bid.Bid
}

func newItem(name string) (*Item, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		// note, escape name when logging
		return nil, fmt.Errorf("could not generate id [%q]: %w", name, err)
	}
	return &Item{
		ID:   id,
		name: name,
		bids: make([]*bid.Bid, 0),
	}, nil
}

func (i *Item) addBid(ctx context.Context, b *bid.Bid) error {
	if b == nil {
		return errors.New("could not add bid, bid is nil")
	}
	if i == nil {
		return fmt.Errorf("could not add bid [%s], item is nil", b.ID.String())
	}

	// first bid is always the highest
	if i.highest == nil {
		i.highest = b
	}

	// add the bid
	// likely, performance could be optimized by keeping track of length and capacity and utilizing the underlaying array
	// manually. though, this is error prone when not done carefully so I consider it out of scope of this challenge
	i.bids = append(i.bids, b)

	// earlier bids with the same amount are staying highest
	if b.Amount > i.highest.Amount {
		i.highest = b
	}
	return nil
}
