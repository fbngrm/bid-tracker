package api

import (
	"context"
	"time"

	"github.com/fbngrm/bid-tracker/pkg/bid"
	"github.com/fbngrm/bid-tracker/pkg/item"
	"github.com/google/uuid"
)

type BidService interface {
	CreateBid(ctx context.Context, itemID, userID uuid.UUID, amount float32, t time.Time) (*bid.Bid, error)
}

type ItemService interface {
	GetHighestBidForItem(ctx context.Context, itemID uuid.UUID) (*bid.Bid, error)
	GetItemsForBids(ctx context.Context, bids []*bid.Bid) ([]*item.Item, error)
	GetBidsForItem(ctx context.Context, itemId uuid.UUID) ([]*bid.Bid, error)
}

type UserService interface {
	GetBidsForUser(ctx context.Context, userID uuid.UUID) ([]*bid.Bid, error)
}
