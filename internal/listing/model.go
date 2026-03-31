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

	// ErrEmptyName indicates that a name is empty.
	ErrEmptyName = errors.New("empty name")

	// ErrEmptyTitle indicates that a listing title is empty.
	ErrEmptyTitle = errors.New("empty title")

	// ErrEmptyDescription indicates that a listing description is empty.
	ErrEmptyDescription = errors.New("empty description")

	// ErrEmptyLocation indicates that a zone location is empty.
	ErrEmptyLocation = errors.New("empty location")

	// ErrInvalidPrice indicates that a price is invalid for publication.
	ErrInvalidPrice = errors.New("invalid price")

	// ErrInvalidVolume indicates that a volume is invalid.
	ErrInvalidVolume = errors.New("invalid volume")

	// ErrInvalidStatus indicates that a status transition or stored status is invalid.
	ErrInvalidStatus = errors.New("invalid status")
)

// ID identifies a listing.
type ID string

// ProductID identifies a product.
type ProductID string

// ZoneID identifies a zone.
type ZoneID string

// Status describes the lifecycle state of a listing.
type Status string

const (
	// StatusDraft marks a listing that is still being prepared.
	StatusDraft Status = "draft"

	// StatusReady marks a listing that is ready for export.
	StatusReady Status = "ready"

	// StatusExported marks a listing that has been exported.
	StatusExported Status = "exported"

	// StatusArchived marks a listing that is no longer active.
	StatusArchived Status = "archived"
)

// Mix describes concrete mix characteristics.
type Mix struct {
	Grade           string
	Class           string
	Mobility        string
	FrostResistance string
	WaterResistance string
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

// State contains full listing state for repository restore and persistence.
type State struct {
	ID        ID
	ProductID ProductID
	ZoneID    ZoneID
	Text      Text
	Price     Price
	Images    []Image
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Product describes a concrete product being offered.
type Product struct {
	id          ProductID
	name        string
	mix         Mix
	minVolumeM3 float64
}

// NewProduct creates a new product.
func NewProduct(id ProductID, name string, mix Mix, minVolumeM3 float64) (Product, error) {
	if isBlank(string(id)) {
		return Product{}, ErrEmptyID
	}

	if isBlank(name) {
		return Product{}, ErrEmptyName
	}

	if minVolumeM3 <= 0 {
		return Product{}, ErrInvalidVolume
	}

	return Product{
		id:          id,
		name:        strings.TrimSpace(name),
		mix:         normalizeMix(mix),
		minVolumeM3: minVolumeM3,
	}, nil
}

// ID returns the product ID.
func (p Product) ID() ProductID {
	return p.id
}

// Name returns the product name.
func (p Product) Name() string {
	return p.name
}

// Mix returns the product mix characteristics.
func (p Product) Mix() Mix {
	return p.mix
}

// MinVolumeM3 returns the minimal order volume.
func (p Product) MinVolumeM3() float64 {
	return p.minVolumeM3
}

// Zone describes where a listing is relevant.
type Zone struct {
	id       ZoneID
	region   string
	city     string
	district string
	address  string
	delivery bool
}

// NewZone creates a new zone.
func NewZone(id ZoneID, region, city, district, address string, delivery bool) (Zone, error) {
	if isBlank(string(id)) {
		return Zone{}, ErrEmptyID
	}

	region = strings.TrimSpace(region)
	city = strings.TrimSpace(city)
	district = strings.TrimSpace(district)
	address = strings.TrimSpace(address)

	if region == "" && city == "" && district == "" && address == "" {
		return Zone{}, ErrEmptyLocation
	}

	return Zone{
		id:       id,
		region:   region,
		city:     city,
		district: district,
		address:  address,
		delivery: delivery,
	}, nil
}

// ID returns the zone ID.
func (z Zone) ID() ZoneID {
	return z.id
}

// Region returns the zone region.
func (z Zone) Region() string {
	return z.region
}

// City returns the zone city.
func (z Zone) City() string {
	return z.city
}

// District returns the zone district.
func (z Zone) District() string {
	return z.district
}

// Address returns the zone address.
func (z Zone) Address() string {
	return z.address
}

// Delivery reports whether delivery is available in the zone.
func (z Zone) Delivery() bool {
	return z.delivery
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

// NewListing creates a new draft listing.
func NewListing(id ID, productID ProductID, zoneID ZoneID, now time.Time) (*Listing, error) {
	if isBlank(string(id)) {
		return nil, ErrEmptyID
	}

	if isBlank(string(productID)) {
		return nil, errors.New("empty product id")
	}

	if isBlank(string(zoneID)) {
		return nil, errors.New("empty zone id")
	}

	return &Listing{
		id:        id,
		productID: productID,
		zoneID:    zoneID,
		status:    StatusDraft,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// Restore rebuilds a listing from stored state.
func Restore(state State) (*Listing, error) {
	if isBlank(string(state.ID)) {
		return nil, ErrEmptyID
	}

	if isBlank(string(state.ProductID)) {
		return nil, errors.New("empty product id")
	}

	if isBlank(string(state.ZoneID)) {
		return nil, errors.New("empty zone id")
	}

	if !validStatus(state.Status) {
		return nil, fmt.Errorf("%w: %q", ErrInvalidStatus, state.Status)
	}

	state.Text = normalizeText(state.Text)
	state.Price = normalizePrice(state.Price)
	state.Images = normalizeImages(state.Images)

	switch state.Status {
	case StatusReady, StatusExported:
		if err := validatePublishable(state.Text, state.Price); err != nil {
			return nil, err
		}
	case StatusDraft, StatusArchived:
		if err := validateDraftPrice(state.Price); err != nil {
			return nil, err
		}
	}

	return &Listing{
		id:        state.ID,
		productID: state.ProductID,
		zoneID:    state.ZoneID,
		text:      state.Text,
		price:     state.Price,
		images:    state.Images,
		status:    state.Status,
		createdAt: state.CreatedAt,
		updatedAt: state.UpdatedAt,
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

// Text returns the full listing text.
func (l *Listing) Text() Text {
	return l.text
}

// Price returns the listing price.
func (l *Listing) Price() Price {
	return l.price
}

// Images returns the listing images.
func (l *Listing) Images() []Image {
	return copyImages(l.images)
}

// Status returns the current listing status.
func (l *Listing) Status() Status {
	return l.status
}

// CreatedAt returns the listing creation time.
func (l *Listing) CreatedAt() time.Time {
	return l.createdAt
}

// UpdatedAt returns the listing update time.
func (l *Listing) UpdatedAt() time.Time {
	return l.updatedAt
}

// Snapshot returns the full listing state.
func (l *Listing) Snapshot() State {
	return State{
		ID:        l.id,
		ProductID: l.productID,
		ZoneID:    l.zoneID,
		Text:      l.text,
		Price:     l.price,
		Images:    copyImages(l.images),
		Status:    l.status,
		CreatedAt: l.createdAt,
		UpdatedAt: l.updatedAt,
	}
}

// Rewrite replaces the listing text.
//
// Any text change moves a ready or exported listing back to draft.
func (l *Listing) Rewrite(text Text, now time.Time) {
	l.text = normalizeText(text)
	l.updatedAt = now

	if l.status == StatusReady || l.status == StatusExported {
		l.status = StatusDraft
	}
}

// ChangePrice replaces the listing price.
//
// Zero price is allowed for drafts. Negative price is rejected.
// Any price change moves a ready or exported listing back to draft.
func (l *Listing) ChangePrice(price Price, now time.Time) error {
	price = normalizePrice(price)

	if err := validateDraftPrice(price); err != nil {
		return err
	}

	l.price = price
	l.updatedAt = now

	if l.status == StatusReady || l.status == StatusExported {
		l.status = StatusDraft
	}

	return nil
}

// ReplaceImages replaces the listing images.
//
// Any image change moves a ready or exported listing back to draft.
func (l *Listing) ReplaceImages(images []Image, now time.Time) {
	l.images = normalizeImages(images)
	l.updatedAt = now

	if l.status == StatusReady || l.status == StatusExported {
		l.status = StatusDraft
	}
}

// Ready marks the listing as ready for export.
func (l *Listing) Ready(now time.Time) error {
	if l.status == StatusArchived {
		return fmt.Errorf("%w: archived listing cannot become ready", ErrInvalidStatus)
	}

	if err := validatePublishable(l.text, l.price); err != nil {
		return err
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

func validatePublishable(text Text, price Price) error {
	if isBlank(text.Title) {
		return ErrEmptyTitle
	}

	if isBlank(text.Description) {
		return ErrEmptyDescription
	}

	if price.Amount <= 0 {
		return ErrInvalidPrice
	}

	return nil
}

func validateDraftPrice(price Price) error {
	if price.Amount < 0 {
		return ErrInvalidPrice
	}

	return nil
}

func validStatus(status Status) bool {
	switch status {
	case StatusDraft, StatusReady, StatusExported, StatusArchived:
		return true
	default:
		return false
	}
}

func normalizeMix(mix Mix) Mix {
	return Mix{
		Grade:           strings.TrimSpace(mix.Grade),
		Class:           strings.TrimSpace(mix.Class),
		Mobility:        strings.TrimSpace(mix.Mobility),
		FrostResistance: strings.TrimSpace(mix.FrostResistance),
		WaterResistance: strings.TrimSpace(mix.WaterResistance),
	}
}

func normalizeText(text Text) Text {
	return Text{
		Title:       strings.TrimSpace(text.Title),
		Description: strings.TrimSpace(text.Description),
	}
}

func normalizePrice(price Price) Price {
	return Price{
		Amount:   price.Amount,
		Currency: normalizeCurrency(price.Currency),
	}
}

func normalizeCurrency(v string) string {
	v = strings.TrimSpace(strings.ToUpper(v))
	if v == "" {
		return "RUB"
	}

	return v
}

func normalizeImages(images []Image) []Image {
	out := make([]Image, 0, len(images))
	for _, image := range images {
		url := strings.TrimSpace(image.URL)
		if url == "" {
			continue
		}

		out = append(out, Image{URL: url})
	}

	return out
}

func copyImages(images []Image) []Image {
	out := make([]Image, len(images))
	copy(out, images)

	return out
}

func isBlank(v string) bool {
	return strings.TrimSpace(v) == ""
}
