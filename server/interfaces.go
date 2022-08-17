package server

import (
	"context"

	auctionv1 "github.com/fbngrm/bid-tracker/gen/proto/go/auction/v1"
)

type AuctionServiceClient interface {
	// Creates a new bid for an item.
	CreateBid(ctx context.Context, in *auctionv1.CreateBidRequest) (*auctionv1.Bid, error)
	// Get the highest bid for an item.
	GetHighestBid(ctx context.Context, in *auctionv1.GetHighestBidRequest) (*auctionv1.Bid, error)
	// Get all bids for an item.
	GetBids(ctx context.Context, in *auctionv1.GetBidsRequest) (*auctionv1.Bids, error)
	// Get all items a user holds bids for.
	GetItemsForUserBids(ctx context.Context, in *auctionv1.GetItemsForUserBidsRequest) (*auctionv1.Items, error)
}
