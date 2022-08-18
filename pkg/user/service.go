package user

import (
	"context"
	"fmt"

	"github.com/fbngrm/bid-tracker/pkg/bid"
	"github.com/google/uuid"
)

type Service struct {
	store *userStore
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
