package api

import (
	"context"
	"errors"
	"fmt"

	auctionv1 "github.com/fbngrm/bid-tracker/gen/proto/go/auction/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Api struct {
	bidService  BidService
	itemService ItemService
	userService UserService
}

func NewApi(b BidService, g ItemService) *Api {
	return &Api{
		bidService:  b,
		itemService: g,
	}
}

// Note, the error handling is overly simplified here due to time constraints. Internal errors must not be exposed
// to the outside world. Also, we must escape logging of potentially harmful user input using %q formatting directive
// for strings. Ideally, we log a correlation-id from the request, e.g. a trace-id, in all error logs. We assume, an
// item and user with given id exist. In a more realistic scenario, we must assert this.
func (a *Api) CreateBid(ctx context.Context, in *auctionv1.CreateBidRequest) (*auctionv1.Bid, error) {
	if in == nil {
		return nil, errors.New("could not create bid, request is nil")
	}

	itemID, err := uuid.Parse(in.Bid.ItemId)
	if err != nil {
		return nil, fmt.Errorf("could not parse uuid from item id [%q]: %w", in.Bid.ItemId, err)
	}

	userID, err := uuid.Parse(in.Bid.UserId)
	if err != nil {
		return nil, fmt.Errorf("could not parse uuid from user id [%q]: %w", in.Bid.UserId, err)
	}

	bid, err := a.bidService.CreateBid(ctx, itemID, userID, in.Bid.Amount, in.Bid.Timestamp.AsTime())
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

func (a *Api) GetHighestBidForItem(ctx context.Context, in *auctionv1.GetHighestBidRequest) (*auctionv1.Bid, error) {
	if in == nil {
		return nil, errors.New("could not get highest bid, request is nil")
	}

	itemID, err := uuid.Parse(in.ItemId)
	if err != nil {
		return nil, fmt.Errorf("could not parse uuid from item id [%q]: %w", in.ItemId, err)
	}

	bid, err := a.itemService.GetHighestBidForItem(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("could not get highest bid for item [%q]: %w", in.ItemId, err)
	}
	return &auctionv1.Bid{
		Id:        bid.ID.String(),
		ItemId:    bid.ItemID.String(),
		UserId:    bid.UserID.String(),
		Timestamp: timestamppb.New(bid.Timestamp),
	}, nil
}

func (a *Api) GetBidsForItem(ctx context.Context, in *auctionv1.GetBidsRequest) (*auctionv1.Bids, error) {
	if in == nil {
		return nil, errors.New("could not get highest bid, request is nil")
	}

	itemID, err := uuid.Parse(in.ItemId)
	if err != nil {
		return nil, fmt.Errorf("could not parse uuid from item id [%q]: %w", in.ItemId, err)
	}

	itemBids, err := a.itemService.GetBidsForItem(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("could not get bids for item [%q]: %w", in.ItemId, err)
	}

	bids := make([]*auctionv1.Bid, len(itemBids))
	for i, b := range itemBids {
		bids[i] = &auctionv1.Bid{
			Id:        b.ID.String(),
			ItemId:    b.ItemID.String(),
			UserId:    b.UserID.String(),
			Timestamp: timestamppb.New(b.Timestamp),
		}
	}
	return &auctionv1.Bids{Bids: bids}, nil
}

func (a *Api) GetItemsForUserBids(ctx context.Context, in *auctionv1.GetItemsForUserBidsRequest) (*auctionv1.Items, error) {
	if in == nil {
		return nil, errors.New("could not get highest bid, request is nil")
	}

	userID, err := uuid.Parse(in.UserId)
	if err != nil {
		return nil, fmt.Errorf("could not parse uuid from user id [%q]: %w", in.UserId, err)
	}

	bids, err := a.userService.GetBidsForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("could not get bids for user id [%q]: %w", in.UserId, err)
	}

	storeItems, err := a.itemService.GetItemsForBids(ctx, bids)
	if err != nil {
		return nil, fmt.Errorf("could not get items for bids for user id [%q]: %w", in.UserId, err)
	}

	items := make([]*auctionv1.Item, len(storeItems))
	for i, it := range storeItems {
		items[i] = &auctionv1.Item{
			Id:   it.ID.String(),
			Name: it.Name,
		}
	}
	return &auctionv1.Items{Items: items}, nil
}
