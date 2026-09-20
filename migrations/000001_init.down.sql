ALTER TABLE users ALTER COLUMN habits DROP NOT NULL;

-- migrate -path migrations -database "postgres://postgres@localhost:5432/postgres?sslmode=disable" down 1