package user

import (
	"context"
	"fmt"

	"github.com/fbngrm/bid-tracker/pkg/bid"
	"github.com/google/uuid"
)

type Service struct {
	store *store
}

func NewService() *Service {
	return &Service{
		store: newStore(),
	}
}

func (s *Service) CreateUser(ctx context.Context) (*User, error) {
	u, err := NewUser()
	if err != nil {
		return nil, fmt.Errorf("could not crate new user")
	}
	s.store.register(ctx, u)
	return u, nil
}

func (s *Service) AddBidToUser(ctx context.Context, b *bid.Bid) error {
	err := s.store.write(ctx, b)
	if err != nil {
		return fmt.Errorf("could not write bid to store for user [%s]", b.ID)
	}
	return nil
}

func (s *Service) GetBidsForUser(ctx context.Context, userID uuid.UUID) ([]*bid.Bid, error) {
	bids, err := s.store.readBids(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("could not read bids from store for user [%s]", userID)
	}
	return bids, nil
}
