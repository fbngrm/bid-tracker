package bid

import (
	"context"
)

type ItemService interface {
	PlaceBidForItem(ctx context.Context, b *Bid) error
}

type UserService interface {
	AddBidToUser(ctx context.Context, b *Bid) error
}
