package item

import (
	"fmt"

	"github.com/fbngrm/bid-tracker/pkg/bid"
	"github.com/google/uuid"
)

// Note, Item is not safe for concurrent use.
type Item struct {
	ID      uuid.UUID // uuid RFC 4122
	name    string
	highest *bid.Bid
	bids    []*bid.Bid
}

func NewItem(name string) (*Item, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("could not generate id [%q]: %w", name, err)
	}
	return &Item{
		ID:   id,
		name: name,
		bids: make([]*bid.Bid, 0),
	}, nil
}

func (i *Item) Bid(b *bid.Bid) {
	if b == nil {
		return
	}
	if i == nil {
		return
	}

	// first bid is always the highest
	if i.highest == nil {
		i.highest = b
	}

	// add the bid
	// likely, performance could be optimized by keeping track of
	// length and capacity and utilizing the underlaying array manually
	// this is error prone when not done properly so I consider it out
	// of scope of this challenge
	i.bids = append(i.bids, b)

	// earlier bids with the same amount are staying highest
	if b.Amount > i.highest.Amount {
		i.highest = b
	}
}
