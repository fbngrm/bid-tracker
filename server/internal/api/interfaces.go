package api

import (
	"context"
	"time"

	auctionv1 "github.com/fbngrm/bid-tracker/gen/proto/go/auction/v1"
	"github.com/fbngrm/bid-tracker/pkg/bid"
	"github.com/google/uuid"
)

type BidGetter interface {
	GetBidsForItem(ctx context.Context, in *auctionv1.GetBidsRequest) (*auctionv1.Bids, error)
}

type BidService interface {
	CreateBid(ctx context.Context, itemID, userID uuid.UUID, amount float32, t time.Time) (*bid.Bid, error)
}

type ItemService interface {
	GetHighestBidForItem(ctx context.Context, itemID uuid.UUID) (*bid.Bid, error)
	GetItemsForUserBids(ctx context.Context, in *auctionv1.GetItemsForUserBidsRequest) (*auctionv1.Items, error)
}
