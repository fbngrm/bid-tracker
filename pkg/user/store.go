package user

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/fbngrm/bid-tracker/pkg/bid"
	"github.com/google/uuid"
)

type userStore struct {
	sync.RWMutex
	users map[uuid.UUID]User
}

// Idempotent, already registered users are ignored.
func (us *userStore) register(ctx context.Context, u User) {
	us.Lock()
	defer us.Unlock()

	if _, ok := us.users[u.ID]; ok {
		return
	}
	us.users[u.ID] = u
}

func (us *userStore) write(ctx context.Context, b *bid.Bid) error {
	us.Lock()
	defer us.Unlock()

	if b == nil {
		return errors.New("could not write bid, bid is nil")
	}

	u, ok := us.users[b.UserID]
	if !ok {
		return fmt.Errorf("could not read user [%s], not registered", b.UserID)
	}

	u.bids = append(u.bids, b)
	us.users[b.UserID] = u

	return nil
}

func (us *userStore) readBids(ctx context.Context, userID uuid.UUID) ([]*bid.Bid, error) {
	us.RLock()
	defer us.RUnlock()

	u, ok := us.users[userID]
	if !ok {
		return nil, fmt.Errorf("could not read user [%s], not registered", userID)
	}

	bids := make([]*bid.Bid, len(u.bids))
	for i, b := range u.bids {
		bids[i] = b.Copy()
	}

	return bids, nil
}
