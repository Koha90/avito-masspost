package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/koha90/avito-masspost/internal/listing"
)

// ProductRepository stores products in PostgreSQL.
type ProductRepository struct {
	db *sql.DB
}

// Products returns a product repository.
func (r *Repository) Products() *ProductRepository {
	return &ProductRepository{db: r.db}
}

// Save creates or updates a product.
func (r *ProductRepository) Save(ctx context.Context, product listing.Product) error {
	const q = `
		INSERT INTO products (
			id,
			name,
			grade,
			class,
			mobility,
			frost_resistance,
			water_resistance,
			min_volume_m3,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET
			name = EXCLUDED.name,
			grade = EXCLUDED.grade,
			class = EXCLUDED.class,
			mobility = EXCLUDED.mobility,
			frost_resistance = EXCLUDED.frost_resistance,
			water_resistance = EXCLUDED.water_resistance,
			min_volume_m3 = EXCLUDED.min_volume_m3,
			updated_at = NOW()
	`

	mix := product.Mix()

	_, err := r.db.ExecContext(
		ctx,
		q,
		string(product.ID()),
		product.Name(),
		mix.Grade,
		mix.Class,
		mix.Mobility,
		mix.FrostResistance,
		mix.WaterResistance,
		product.MinVolumeM3(),
	)
	if err != nil {
		return fmt.Errorf("save product: %w", err)
	}

	return nil
}

// ByID returns a product by its ID.
func (r *ProductRepository) ByID(ctx context.Context, id listing.ProductID) (listing.Product, error) {
	const q = `
		SELECT
			id,
			name,
			grade,
			class,
			mobility,
			frost_resistance,
			water_resistance,
			min_volume_m3
		FROM products
		WHERE id = $1
	`

	var row productRow
	if err := r.db.QueryRowContext(ctx, q, string(id)).Scan(
		&row.ID,
		&row.Name,
		&row.Grade,
		&row.Class,
		&row.Mobility,
		&row.FrostResistance,
		&row.WaterResistance,
		&row.MinVolumeM3,
	); err != nil {
		if err == sql.ErrNoRows {
			return listing.Product{}, fmt.Errorf("product by id %q: %w", id, err)
		}

		return listing.Product{}, fmt.Errorf("query product by id %q: %w", id, err)
	}

	product, err := row.Product()
	if err != nil {
		return listing.Product{}, fmt.Errorf("build product: %w", err)
	}

	return product, nil
}

// List returns all products.
func (r *ProductRepository) List(ctx context.Context) ([]listing.Product, error) {
	const q = `
		SELECT
			id,
			name,
			grade,
			class,
			mobility,
			frost_resistance,
			water_resistance,
			min_volume_m3
		FROM products
		ORDER BY name, id
	`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	out := make([]listing.Product, 0)
	for rows.Next() {
		var row productRow
		if err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.Grade,
			&row.Class,
			&row.Mobility,
			&row.FrostResistance,
			&row.WaterResistance,
			&row.MinVolumeM3,
		); err != nil {
			return nil, fmt.Errorf("scan product: %w", err)
		}

		product, err := row.Product()
		if err != nil {
			return nil, fmt.Errorf("build product: %w", err)
		}

		out = append(out, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate products: %w", err)
	}

	return out, nil
}
