package bid

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	itemService ItemService
	userService UserService
}

func NewService(is ItemService, us UserService) *Service {
	return &Service{
		itemService: is,
		userService: us,
	}
}

// Fixme, we have to copy here to protect against race condition between itemStore, user and caller.
// We should operate on copies by default instead of pointers to resolve this.
func (s *Service) CreateBid(ctx context.Context, itemID, userID uuid.UUID, amount float32, t time.Time) (*Bid, error) {
	b, err := NewBid(itemID, userID, amount, t)
	if err != nil {
		return nil, fmt.Errorf("could not create new bid for user [%s] for item [%s]: %w", userID, itemID, err)
	}

	// in a future event-based version, we would emit a 'BidCreated' event and consume it in a user service that would
	// take care or adding it to the user. also, we could run it in a go routine here to optimize performance.
	err = s.userService.AddBidToUser(ctx, b.Copy())
	if err != nil {
		return nil, fmt.Errorf("could not add bid to user [%s] for item [%s]: %w", userID, itemID, err)
	}

	err = s.itemService.PlaceBidForItem(ctx, b.Copy())
	if err != nil {
		return nil, fmt.Errorf("could not place bid for user [%s] for item [%s]: %w", userID, itemID, err)
	}

	return b, nil
}
