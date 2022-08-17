package item

import "fmt"

type ItemService struct {
	store *ItemStore
}

func NewService() *ItemService {
	return &ItemService{
		store: NewStore(),
	}
}

// we must escape logging of potentially harmful user input
func (s *ItemService) CreateItem(name string) (*Item, error) {
	i, err := NewItem(name)
	if err != nil {
		return nil, fmt.Errorf("could not create new item [%q]: %w", name, err)
	}
	err = s.store.Register(i)
	if err != nil {
		return nil, fmt.Errorf("could not register item [%q]: %w", name, err)
	}
	return i, nil
}
