package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/koha90/avito-masspost/internal/listing"
)

// ListingRepository stores listings in PostgreSQL.
type ListingRepository struct {
	db *sql.DB
}

// Listings returns a listing repository.
func (r *Repository) Listings() *ListingRepository {
	return &ListingRepository{db: r.db}
}

// Save creates or updates a listing.
func (r *ListingRepository) Save(ctx context.Context, item *listing.Listing) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin listing save tx: %w", err)
	}
	defer tx.Rollback()

	const saveListing = `
		INSERT INTO listings (
			id,
			product_id,
			zone_id,
			title,
			description,
			price_amount,
			price_currency,
			status,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO UPDATE SET
			product_id = EXCLUDED.product_id,
			zone_id = EXCLUDED.zone_id,
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			price_amount = EXCLUDED.price_amount,
			price_currency = EXCLUDED.price_currency,
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at
	`

	text := item.Text()
	price := item.Price()

	_, err = tx.ExecContext(
		ctx,
		saveListing,
		string(item.ID()),
		string(item.ProductID()),
		string(item.ZoneID()),
		text.Title,
		text.Description,
		price.Amount,
		price.Currency,
		string(item.Status()),
		item.CreatedAt(),
		item.UpdatedAt(),
	)
	if err != nil {
		return fmt.Errorf("save listing: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM listing_images WHERE listing_id = $1`, string(item.ID())); err != nil {
		return fmt.Errorf("delete listing images: %w", err)
	}

	for i, image := range item.Images() {
		const saveImage = `
			INSERT INTO listing_images (listing_id, position, url)
			VALUES ($1, $2, $3)
		`

		if _, err := tx.ExecContext(ctx, saveImage, string(item.ID()), i, image.URL); err != nil {
			return fmt.Errorf("save listing image: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit listing save tx: %w", err)
	}

	return nil
}

// ByID returns a listing by its ID.
func (r *ListingRepository) ByID(ctx context.Context, id listing.ID) (*listing.Listing, error) {
	const q = `
		SELECT
			id,
			product_id,
			zone_id,
			title,
			description,
			price_amount,
			price_currency,
			status,
			created_at,
			updated_at
		FROM listings
		WHERE id = $1
	`

	var row listingRow
	if err := r.db.QueryRowContext(ctx, q, string(id)).Scan(
		&row.ID,
		&row.ProductID,
		&row.ZoneID,
		&row.Title,
		&row.Description,
		&row.PriceAmount,
		&row.PriceCurrency,
		&row.Status,
		&row.CreatedAt,
		&row.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("listing by id %q: %w", id, err)
		}

		return nil, fmt.Errorf("query listing by id %q: %w", id, err)
	}

	images, err := r.images(ctx, id)
	if err != nil {
		return nil, err
	}
	row.Images = images

	return row.Listing()
}

// List returns listings matching the filter.
func (r *ListingRepository) List(ctx context.Context, filter listing.Filter) ([]*listing.Listing, error) {
	var (
		args  []any
		where []string
	)

	q := `
		SELECT
			id,
			product_id,
			zone_id,
			title,
			description,
			price_amount,
			price_currency,
			status,
			created_at,
			updated_at
		FROM listings
	`

	if filter.ProductID != "" {
		args = append(args, string(filter.ProductID))
		where = append(where, fmt.Sprintf("product_id = $%d", len(args)))
	}

	if filter.ZoneID != "" {
		args = append(args, string(filter.ZoneID))
		where = append(where, fmt.Sprintf("zone_id = $%d", len(args)))
	}

	if filter.Status != "" {
		args = append(args, string(filter.Status))
		where = append(where, fmt.Sprintf("status = $%d", len(args)))
	}

	if len(where) > 0 {
		q += " WHERE " + strings.Join(where, " AND ")
	}

	q += " ORDER BY updated_at DESC, id"

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("query listings: %w", err)
	}
	defer rows.Close()

	out := make([]*listing.Listing, 0)
	for rows.Next() {
		var row listingRow
		if err := rows.Scan(
			&row.ID,
			&row.ProductID,
			&row.ZoneID,
			&row.Title,
			&row.Description,
			&row.PriceAmount,
			&row.PriceCurrency,
			&row.Status,
			&row.CreatedAt,
			&row.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan listing: %w", err)
		}

		images, err := r.images(ctx, listing.ID(row.ID))
		if err != nil {
			return nil, err
		}
		row.Images = images

		item, err := row.Listing()
		if err != nil {
			return nil, err
		}

		out = append(out, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listings: %w", err)
	}

	return out, nil
}

func (r *ListingRepository) images(ctx context.Context, id listing.ID) ([]string, error) {
	const q = `
		SELECT url
		FROM listing_images
		WHERE listing_id = $1
		ORDER BY position
	`

	rows, err := r.db.QueryContext(ctx, q, string(id))
	if err != nil {
		return nil, fmt.Errorf("query listing images: %w", err)
	}
	defer rows.Close()

	out := make([]string, 0)
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, fmt.Errorf("scan listing image: %w", err)
		}

		out = append(out, url)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listing images: %w", err)
	}

	return out, nil
}
