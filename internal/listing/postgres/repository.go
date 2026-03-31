// Package postgres provides PostgreSQL repositories for the listing domain.
package postgres

import "database/sql"

// Repository provides access to PostgreSQL-backend listing repositories.
type Repository struct {
	db *sql.DB
}

// NewRepository returns a new PostgreSQL repository set.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}
