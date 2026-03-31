CREATE TABLE products (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    grade TEXT NOT NULL DEFAULT '',
    class TEXT NOT NULL DEFAULT '',
    mobility TEXT NOT NULL DEFAULT '',
    frost_resistance TEXT NOT NULL DEFAULT '',
    water_resistance TEXT NOT NULL DEFAULT '',
    min_volume_m3 DOUBLE PRECISION NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE zones (
    id TEXT PRIMARY KEY,
    region TEXT NOT NULL DEFAULT '',
    city TEXT NOT NULL DEFAULT '',
    district TEXT NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    delivery BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE listings (
    id TEXT PRIMARY KEY,
    product_id TEXT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    zone_id TEXT NOT NULL REFERENCES zones(id) ON DELETE RESTRICT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    price_amount BIGINT NOT NULL,
    price_currency TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX listings_product_id_idx ON listings(product_id);
CREATE INDEX listings_zone_id_idx ON listings(zone_id);
CREATE INDEX listings_status_idx ON listings(status);

CREATE TABLE listing_images (
    listing_id TEXT NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    position INT NOT NULL,
    url TEXT NOT NULL,
    PRIMARY KEY (listing_id, position)
);

CREATE INDEX listing_images_listing_id_idx ON listing_images(listing_id);
