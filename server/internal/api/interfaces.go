package api

import (
	"context"
	"time"

	auctionv1 "github.com/fbngrm/bid-tracker/gen/proto/go/auction/v1"
	"github.com/fbngrm/bid-tracker/pkg/bid"
	"github.com/google/uuid"
)

type Bidder interface {
	CreateBid(itemID, userID uuid.UUID, amount float32, t time.Time) (*bid.Bid, error)
}

type BidGetter interface {
	GetHighestBidForItem(ctx context.Context, in *auctionv1.GetHighestBidRequest) (*auctionv1.Bid, error)
	GetBidsForItem(ctx context.Context, in *auctionv1.GetBidsRequest) (*auctionv1.Bids, error)
}

type ItemGetter interface {
	GetItemsForUserBids(ctx context.Context, in *auctionv1.GetItemsForUserBidsRequest) (*auctionv1.Items, error)
}
