package bid

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type BidService struct {
	itemService ItemService
	userService UserService
}

func NewBidService(is ItemService, us UserService) *BidService {
	return &BidService{
		itemService: is,
		userService: us,
	}
}

func (s *BidService) CreateBid(itemID, userID uuid.UUID, amount float32, t time.Time) (*Bid, error) {
	b, err := NewBid(itemID, userID, amount, t)
	if err != nil {
		return nil, fmt.Errorf("could not create new bid for user [%s] for item [%s]: %w", userID, itemID, err)
	}

	err = s.itemService.PlaceBidForItem(b)
	if err != nil {
		return nil, fmt.Errorf("could not place bid for user [%s] for item [%s]: %w", userID, itemID, err)
	}

	// in a future event-based version, we would emit a 'BidCreated' event and consume it
	// in a user service that would take care or adding it to the user
	err = s.userService.AddBidToUser(b)
	if err != nil {
		return nil, fmt.Errorf("could not add bid to user [%s] for item [%s]: %w", userID, itemID, err)
	}
	return b, nil
}
