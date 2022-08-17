package api

import (
	"context"

	auctionv1 "github.com/fbngrm/bid-tracker/gen/proto/go/auction/v1"
)

type Api struct {
	bidService Bidder
}

func NewApi(b Bidder) *Api {
	return &Api{
		bidService: b,
	}
}

// the error handling is overly simplified here due to time constraints.
// internal errors must not be exposed to the outside world!
func (a *Api) CreateBid(ctx context.Context, in *auctionv1.CreateBidRequest) (*auctionv1.Bid, error) {
	// we need to take special care here since all request input is potentially harmful

	// Id     int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	// ItemId int64 `protobuf:"varint,2,opt,name=item_id,json=itemId,proto3" json:"item_id,omitempty"`
	// UserId int64 `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`

	// // varint encoding, 4 Bytes only until 2038
	// Timestamp *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`

	// thus, we cannot just assert a type to the req.Material
	return nil, nil
}
