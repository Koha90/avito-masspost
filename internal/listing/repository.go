package listing

import "context"

// Repository provides access to listing storage.
type Repository interface {
	Save(ctx context.Context, listing *Listing) error
	ByID(ctx context.Context, id ID) (*Listing, error)
	List(ctx context.Context, filter Filter) ([]*Listing, error)
}

// ProductRepository provides access to product storage.
type ProductRepository interface {
	Save(ctx context.Context, product Product) error
	ByID(ctx context.Context, id ProductID) (Product, error)
	List(ctx context.Context) ([]Product, error)
}

// ZoneRepository provides access to product storage.
type ZoneRepository interface {
	Save(ctx context.Context, zone Zone) error
	ByID(ctx context.Context, id ZoneID) (Zone, error)
	List(ctx context.Context) ([]Zone, error)
}

// Filter describes listing selection criteria.
type Filter struct {
	ProductID ProductID
	ZoneID    ZoneID
	Status    Status
}
