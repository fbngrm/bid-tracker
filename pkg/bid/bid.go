package bid

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	ID        uuid.UUID // uuid RFC 4122
	ItemID    uuid.UUID // todo, should we store an item pointer here instead?
	UserID    uuid.UUID
	Amount    float32
	Timestamp time.Time
}

// We want to know when a bid was created to be able to replay it later in a event-sourcing
// based architecture. Conveniently, we can get a timestamp for free from gRPC and use it here.
func NewBid(itemID, userID uuid.UUID, amount float32, t time.Time) (*Bid, error) {
	// depending on the datastore we want to use, we might want to let it generate
	// the IDs for us. since we don't have one yet, we do it ourselves.
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("could not generate id: %w", err)
	}
	return &Bid{
		ID:        id,
		ItemID:    itemID,
		UserID:    userID,
		Amount:    amount,
		Timestamp: t,
	}, nil
}

func (b *Bid) Copy() *Bid {
	v := *b
	return &Bid{
		ID:        v.ID,
		ItemID:    v.ItemID,
		UserID:    v.UserID,
		Amount:    v.Amount,
		Timestamp: v.Timestamp,
	}
}
