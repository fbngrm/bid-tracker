package item

import (
	"errors"
	"fmt"

	"github.com/fbngrm/bid-tracker/pkg/bid"
)

type auction struct {
	i       *Item // for now we use item.ID as identifier
	highest *bid.Bid
	bids    []*bid.Bid
}

func (e auction) addBid(b *bid.Bid) error {
	if e.i == nil {
		return fmt.Errorf("could not add bid [%s] to entry, item is nil", b.ID.String())
	}
	if b == nil {
		return errors.New("could not add bid, bid is nil")
	}

	// first bid is always the highest
	if e.highest == nil {
		e.highest = b
	}

	// add the bid
	// likely, performance could be optimized by keeping track of length and capacity and utilizing the underlaying array
	// manually. though, this is error prone when not done carefully so I consider it out of scope of this challenge
	e.bids = append(e.bids, b)

	// earlier bids with the same amount are staying highest
	if b.Amount > e.highest.Amount {
		e.highest = b
	}
	return nil
}
