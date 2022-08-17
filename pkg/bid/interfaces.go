package bid

type ItemService interface {
	PlaceBidForItem(b *Bid) error
}

type UserService interface {
	AddBidToUser(b *Bid) error
}
