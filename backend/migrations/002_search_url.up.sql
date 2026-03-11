ALTER TABLE searches ADD COLUMN url TEXT NOT NULL DEFAULT '';
ALTER TABLE searches DROP COLUMN location;
ALTER TABLE searches DROP COLUMN property_type;
ALTER TABLE searches DROP COLUMN min_price;
ALTER TABLE searches DROP COLUMN max_price;
ALTER TABLE searches DROP COLUMN min_size_sqft;
ALTER TABLE searches DROP COLUMN max_size_sqft;
