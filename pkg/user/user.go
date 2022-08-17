package user

import (
	"fmt"

	"github.com/fbngrm/bid-tracker/pkg/bid"
	"github.com/google/uuid"
)

type User struct {
	ID   uuid.UUID
	bids []*bid.Bid
}

func NewUser() (*User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("could not generate id: %w", err)
	}
	return &User{
		ID:   id,
		bids: make([]*bid.Bid, 0),
	}, nil
}
