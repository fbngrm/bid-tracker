package item

import (
	"fmt"

	"github.com/fbngrm/bid-tracker/pkg/bid"
)

type ItemService struct {
	// for now we keep a tight coupling here
	store *ItemStore
}

func NewService() *ItemService {
	return &ItemService{
		store: &ItemStore{},
	}
}

// Note, we must escape logging of potentially harmful user input using %q formatting directive.
func (s *ItemService) CreateItem(name string) (*Item, error) {
	i, err := newItem(name)
	if err != nil {
		return nil, fmt.Errorf("could not create new item [%q]: %w", name, err)
	}
	err = s.store.register(i)
	if err != nil {
		return nil, fmt.Errorf("could not register item [%q]: %w", name, err)
	}
	return i, nil
}

func (s *ItemService) PlaceBid(b *bid.Bid) error {
	err := s.store.write(b)
	if err != nil {
		return fmt.Errorf("could not write bid [%q] to item store: %w", b.ID.String(), err)
	}
	return nil
}
