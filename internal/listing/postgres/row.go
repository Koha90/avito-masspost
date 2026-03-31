package postgres

import (
	"fmt"
	"time"

	"github.com/koha90/avito-masspost/internal/listing"
)

type productRow struct {
	ID              string
	Name            string
	Grade           string
	Class           string
	Mobility        string
	FrostResistance string
	WaterResistance string
	MinVolumeM3     float64
}

func (r productRow) Product() (listing.Product, error) {
	return listing.NewProduct(
		listing.ProductID(r.ID),
		r.Name,
		listing.Mix{
			Grade:           r.Grade,
			Class:           r.Class,
			Mobility:        r.Mobility,
			FrostResistance: r.FrostResistance,
			WaterResistance: r.WaterResistance,
		},
		r.MinVolumeM3,
	)
}

type zoneRow struct {
	ID       string
	Region   string
	City     string
	District string
	Address  string
	Delivery bool
}

func (r zoneRow) Zone() (listing.Zone, error) {
	return listing.NewZone(
		listing.ZoneID(r.ID),
		r.Region,
		r.City,
		r.District,
		r.Address,
		r.Delivery,
	)
}

type listingRow struct {
	ID            string
	ProductID     string
	ZoneID        string
	Title         string
	Description   string
	PriceAmount   int64
	PriceCurrency string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Images        []string
}

func (r listingRow) Listing() (*listing.Listing, error) {
	images := make([]listing.Image, 0, len(r.Images))
	for _, u := range r.Images {
		images = append(images, listing.Image{URL: u})
	}

	item, err := listing.Restore(
		listing.State{
			ID:        listing.ID(r.ID),
			ProductID: listing.ProductID(r.ProductID),
			ZoneID:    listing.ZoneID(r.ZoneID),
			Text: listing.Text{
				Title:       r.Title,
				Description: r.Description,
			},
			Price: listing.Price{
				Amount:   r.PriceAmount,
				Currency: r.PriceCurrency,
			},
			Images:    images,
			Status:    listing.Status(r.Status),
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("restore listing: %w", err)
	}

	return item, nil
}
