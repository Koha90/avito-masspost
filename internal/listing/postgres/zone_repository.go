package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/koha90/avito-masspost/internal/listing"
)

// ZoneRepository stores zones in PostgreSQL.
type ZoneRepository struct {
	db *sql.DB
}

// Zones returns a zone repository.
func (r *Repository) Zones() *ZoneRepository {
	return &ZoneRepository{db: r.db}
}

// Save creates or updates a zone.
func (r *ZoneRepository) Save(ctx context.Context, zone listing.Zone) error {
	const q = `
		INSERT INTO zones (
			id,
			region,
			city,
			district,
			address,
			delivery,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET
			region = EXCLUDED.region,
			city = EXCLUDED.city,
			district = EXCLUDED.district,
			address = EXCLUDED.address,
			delivery = EXCLUDED.delivery,
			updated_at = NOW()
	`

	_, err := r.db.ExecContext(
		ctx,
		q,
		string(zone.ID()),
		zone.Region(),
		zone.City(),
		zone.District(),
		zone.Address(),
		zone.Delivery(),
	)
	if err != nil {
		return fmt.Errorf("save zone: %w", err)
	}

	return nil
}

// ByID returns a zone by its ID.
func (r *ZoneRepository) ByID(ctx context.Context, id listing.ZoneID) (listing.Zone, error) {
	const q = `
		SELECT
			id,
			region,
			city,
			district,
			address,
			delivery
		FROM zones
		WHERE id = $1
	`

	var row zoneRow
	if err := r.db.QueryRowContext(ctx, q, string(id)).Scan(
		&row.ID,
		&row.Region,
		&row.City,
		&row.District,
		&row.Address,
		&row.Delivery,
	); err != nil {
		if err == sql.ErrNoRows {
			return listing.Zone{}, fmt.Errorf("zone by id %q: %w", id, err)
		}

		return listing.Zone{}, fmt.Errorf("query zone by id %q: %w", id, err)
	}

	zone, err := row.Zone()
	if err != nil {
		return listing.Zone{}, fmt.Errorf("build zone: %w", err)
	}

	return zone, nil
}

// List returns all zones.
func (r *ZoneRepository) List(ctx context.Context) ([]listing.Zone, error) {
	const q = `
		SELECT
			id,
			region,
			city,
			district,
			address,
			delivery
		FROM zones
		ORDER BY region, city, district, id
	`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query zones: %w", err)
	}
	defer rows.Close()

	out := make([]listing.Zone, 0)
	for rows.Next() {
		var row zoneRow
		if err := rows.Scan(
			&row.ID,
			&row.Region,
			&row.City,
			&row.District,
			&row.Address,
			&row.Delivery,
		); err != nil {
			return nil, fmt.Errorf("scan zone: %w", err)
		}

		zone, err := row.Zone()
		if err != nil {
			return nil, fmt.Errorf("build zone: %w", err)
		}

		out = append(out, zone)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate zones: %w", err)
	}

	return out, nil
}
