// Package listing contains the core domain model for publishable listings.
package listing

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrEmptyID indicates that an entity ID is empty.
	ErrEmptyID = errors.New("empty id")

	// ErrEmptyTitle indicates that a listing title is empty.
	ErrEmptyTitle = errors.New("empty title")

	// ErrEmptyDescription indicates that a listing description is empty.
	ErrEmptyDescription = errors.New("empty description")

	// ErrInvalidPrice indicates that a listing price is invalid.
	ErrInvalidPrice = errors.New("invalid price")

	// ErrInvalidStatus indicates that a status transition is invalid.
	ErrInvalidStatus = errors.New("invalid status")
)

// ID identifies a listing.
type ID string

// ProductID identifies a product.
type ProductID string

// ZoneID identifies a delivery or publication zone.
type ZoneID string

// Status describes the current lifecycle state of a listing.
type Status string

const (
	// StatusDraft marks a listing that is still incomplete.
	StatusDraft Status = "draft"

	// StatusReady marks a listing that is ready for export.
	StatusReady Status = "ready"

	// StatusExported marks a listing that has already been exported.
	StatusExported Status = "exported"

	// StatusArchived marks a listing that is no longer active.
	StatusArchived Status = "archived"
)

// Mix describes the technical characteristics of a concrete product.
type Mix struct {
	Grade           string
	Class           string
	Mobility        string
	FrostResistance string
	WaterResistance string
}

// Product describes a concrete product being offered.
type Product struct {
	ID          ProductID
	Name        string
	Mix         Mix
	MinVolumeM3 float64
}

// Zone describes where a listing is relevant.
type Zone struct {
	ID       ZoneID
	Region   string
	City     string
	District string
	Address  string
	Delivery bool
}

// Price describes listing price data.
type Price struct {
	Amount   int64
	Currency string
}

// Image describes a listing image.
type Image struct {
	URL string
}

// Text contains listing title and description.
type Text struct {
	Title       string
	Description string
}

// Listing describes a publishable listing variant.
type Listing struct {
	id        ID
	productID ProductID
	zoneID    ZoneID
	text      Text
	price     Price
	images    []Image
	status    Status
	createdAt time.Time
	updatedAt time.Time
}

// NewProduct creates a new product.
func NewProduct(id ProductID, name string, mix Mix, minVolumeM3 float64) (Product, error) {
	if strings.TrimSpace(string(id)) == "" {
		return Product{}, ErrEmptyID
	}

	if strings.TrimSpace(name) == "" {
		return Product{}, errors.New("empty product name")
	}

	if minVolumeM3 <= 0 {
		return Product{}, errors.New("invalid min volume")
	}

	return Product{
		ID:          id,
		Name:        strings.TrimSpace(name),
		Mix:         mix,
		MinVolumeM3: minVolumeM3,
	}, nil
}

// NewZone creates a new zone.
func NewZone(id ZoneID, region, city, district, address string, delivery bool) (Zone, error) {
	if strings.TrimSpace(string(id)) == "" {
		return Zone{}, ErrEmptyID
	}

	if strings.TrimSpace(region) == "" && strings.TrimSpace(city) == "" {
		return Zone{}, errors.New("empty zone location")
	}

	return Zone{
		ID:       id,
		Region:   strings.TrimSpace(region),
		City:     strings.TrimSpace(city),
		District: strings.TrimSpace(district),
		Address:  strings.TrimSpace(address),
		Delivery: delivery,
	}, nil
}

// NewListing creates a new listing in draft status.
func NewListing(
	id ID,
	productID ProductID,
	zoneID ZoneID,
	text Text,
	price Price,
	images []Image,
	now time.Time,
) (*Listing, error) {
	if strings.TrimSpace(string(id)) == "" {
		return nil, ErrEmptyID
	}

	if strings.TrimSpace(string(productID)) == "" {
		return nil, errors.New("empty product id")
	}

	if strings.TrimSpace(string(zoneID)) == "" {
		return nil, errors.New("empty zone id")
	}

	if err := validateText(text); err != nil {
		return nil, err
	}

	if err := validatePrice(price); err != nil {
		return nil, err
	}

	out := make([]Image, 0, len(images))
	for _, image := range images {
		if strings.TrimSpace(image.URL) == "" {
			continue
		}

		out = append(out, Image{URL: strings.TrimSpace(image.URL)})
	}

	return &Listing{
		id:        id,
		productID: productID,
		zoneID:    zoneID,
		text: Text{
			Title:       strings.TrimSpace(text.Title),
			Description: strings.TrimSpace(text.Description),
		},
		price: Price{
			Amount:   price.Amount,
			Currency: normalizeCurrency(price.Currency),
		},
		images:    out,
		status:    StatusDraft,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ID returns the listing ID.
func (l *Listing) ID() ID {
	return l.id
}

// ProductID returns the related product ID.
func (l *Listing) ProductID() ProductID {
	return l.productID
}

// ZoneID returns the related zone ID.
func (l *Listing) ZoneID() ZoneID {
	return l.zoneID
}

// Title returns the listing title.
func (l *Listing) Title() string {
	return l.text.Title
}

// Description returns the listing description.
func (l *Listing) Description() string {
	return l.text.Description
}

// Text returns full listing text.
func (l *Listing) Text() Text {
	return l.text
}

// Price returns the listing price.
func (l *Listing) Price() Price {
	return l.price
}

// Images returns listing images.
func (l *Listing) Images() []Image {
	out := make([]Image, len(l.images))
	copy(out, l.images)

	return out
}

// Status returns the listing status.
func (l *Listing) Status() Status {
	return l.status
}

// CreatedAt returns the creation time.
func (l *Listing) CreatedAt() time.Time {
	return l.createdAt
}

// UpdatedAt returns the last update time.
func (l *Listing) UpdatedAt() time.Time {
	return l.updatedAt
}

// Rewrite replace the listing title and description.
func (l *Listing) Rewrite(text Text, now time.Time) error {
	if err := validateText(text); err != nil {
		return err
	}

	l.text = Text{
		Title:       strings.TrimSpace(text.Title),
		Description: strings.TrimSpace(text.Description),
	}
	l.updatedAt = now

	if l.status == StatusReady || l.status == StatusExported {
		l.status = StatusDraft
	}

	return nil
}

// ChangePrice replaces the listing price.
func (l *Listing) ChangePrice(price Price, now time.Time) error {
	if err := validatePrice(price); err != nil {
		return err
	}

	l.price = Price{
		Amount:   price.Amount,
		Currency: normalizeCurrency(price.Currency),
	}
	l.updatedAt = now

	if l.status == StatusExported {
		l.status = StatusDraft
	}

	return nil
}

// ReplaceImages replaces the listing images.
func (l *Listing) ReplaceImages(images []Image, now time.Time) {
	out := make([]Image, 0, len(images))
	for _, image := range images {
		if strings.TrimSpace(image.URL) == "" {
			continue
		}

		out = append(out, Image{URL: strings.TrimSpace(image.URL)})
	}

	l.images = out
	l.updatedAt = now
}

// Ready marks the listing as ready for export.
func (l *Listing) Ready(now time.Time) error {
	if err := validateText(l.text); err != nil {
		return err
	}

	if err := validatePrice(l.price); err != nil {
		return err
	}

	if l.status == StatusArchived {
		return fmt.Errorf("%w: archived listing cannot become ready", ErrInvalidStatus)
	}

	l.status = StatusReady
	l.updatedAt = now

	return nil
}

// Exported marks the listing as exported.
func (l *Listing) Exported(now time.Time) error {
	if l.status != StatusReady {
		return fmt.Errorf("%w: listing must be ready before export", ErrInvalidStatus)
	}

	l.status = StatusExported
	l.updatedAt = now

	return nil
}

// Archive marks the listing as archived.
func (l *Listing) Archive(now time.Time) {
	l.status = StatusArchived
	l.updatedAt = now
}

func validateText(text Text) error {
	if strings.TrimSpace(text.Title) == "" {
		return ErrEmptyTitle
	}

	if strings.TrimSpace(text.Description) == "" {
		return ErrEmptyDescription
	}

	return nil
}

func validatePrice(price Price) error {
	if price.Amount <= 0 {
		return ErrInvalidPrice
	}

	return nil
}

func normalizeCurrency(v string) string {
	v = strings.TrimSpace(strings.ToUpper(v))
	if v == "" {
		return "RUB"
	}

	return v
}
