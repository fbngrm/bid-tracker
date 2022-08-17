package bid

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	ID        uuid.UUID // uuid RFC 4122
	ItemID    uuid.UUID // todo, should we store an item pointer here instead?
	userID    uuid.UUID
	Amount    float32
	timestamp time.Time
}

// We want to know when a bid was created to be able to replay later.
// We get a timestamp for free from gRPC (wow, that rhymes :)
func NewBid(itemID, userID uuid.UUID, amount float32, t time.Time) (*Bid, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("could not generate id: %w", err)
	}
	return &Bid{
		ID:        id,
		ItemID:    itemID,
		userID:    userID,
		Amount:    amount,
		timestamp: t,
	}, nil
}
