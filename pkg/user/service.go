package user

import (
	"context"
	"errors"

	"github.com/fbngrm/bid-tracker/pkg/bid"
)

type Service struct{}

func (s *Service) AddBidToUser(ctx context.Context, b *bid.Bid) error {
	return errors.New("not implemented")
}
