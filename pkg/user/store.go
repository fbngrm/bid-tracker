package user

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/fbngrm/bid-tracker/pkg/bid"
	"github.com/google/uuid"
)

type store struct {
	sync.RWMutex
	users map[uuid.UUID]*User
}

func newStore() *store {
	return &store{
		users: make(map[uuid.UUID]*User),
	}
}

// Idempotent, already registered users are ignored.
func (s *store) register(ctx context.Context, u *User) error {
	s.Lock()
	defer s.Unlock()

	if u == nil {
		return errors.New("could not register, user is nil")
	}

	if _, ok := s.users[u.ID]; ok {
		return nil
	}
	s.users[u.ID] = u
	return nil
}

func (s *store) write(ctx context.Context, b *bid.Bid) error {
	s.Lock()
	defer s.Unlock()

	if b == nil {
		return errors.New("could not write bid, bid is nil")
	}

	u, ok := s.users[b.UserID]
	if !ok {
		return fmt.Errorf("could not read user [%s], not registered", b.UserID)
	}

	u.bids = append(u.bids, b)
	s.users[b.UserID] = u

	return nil
}

func (s *store) readBids(ctx context.Context, userID uuid.UUID) ([]*bid.Bid, error) {
	s.RLock()
	defer s.RUnlock()

	u, ok := s.users[userID]
	if !ok {
		return nil, fmt.Errorf("could not read user [%s], not registered", userID)
	}

	bids := make([]*bid.Bid, len(u.bids))
	for i, b := range u.bids {
		bids[i] = b.Copy()
	}

	return bids, nil
}
