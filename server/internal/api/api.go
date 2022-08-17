package api

import (
	"context"
	"fmt"

	auctionv1 "github.com/fbngrm/bid-tracker/gen/proto/go/auction/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Api struct {
	bidService Bidder
	bidGetter  BidGetter
	itemGetter ItemGetter
}

func NewApi(b Bidder, bg BidGetter, ig ItemGetter) *Api {
	return &Api{
		bidService: b,
		bidGetter:  bg,
		itemGetter: ig,
	}
}

// Note, the error handling is overly simplified here due to time constraints. Internal errors must not be exposed to the outside world.
// Also, we must escape logging of potentially harmful user input using %q formatting directive for strings.
// Ideally, we log a correlation-id from the request, e.g. a trace-id, in all error logs.
// We assume, an item and user with given id exist. In a more realistic scenario, we must assert this.
func (a *Api) CreateBid(ctx context.Context, in *auctionv1.CreateBidRequest) (*auctionv1.Bid, error) {
	itemID, err := uuid.Parse(in.Bid.ItemId)
	if err != nil {
		return nil, fmt.Errorf("could not parse uuid for item id [%q]: %w", in.Bid.ItemId, err)
	}

	userID, err := uuid.Parse(in.Bid.UserId)
	if err != nil {
		return nil, fmt.Errorf("could not parse uuid for user id [%q]: %w", in.Bid.UserId, err)
	}

	bid, err := a.bidService.CreateBid(itemID, userID, in.Bid.Amount, in.Bid.Timestamp.AsTime())
	if err != nil {
		return nil, fmt.Errorf("could not create bid from request [%q]: %w", in.Bid.UserId, err)
	}

	return &auctionv1.Bid{
		Id:        bid.ID.String(),
		ItemId:    bid.ItemID.String(),
		UserId:    bid.UserID.String(),
		Timestamp: timestamppb.New(bid.Timestamp),
	}, nil
}
