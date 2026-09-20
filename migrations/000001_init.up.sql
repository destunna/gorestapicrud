ALTER TABLE users ALTER COLUMN habits SET NOT NULL;

-- migrate -path migrations -database "postgres://postgres@localhost:5432/postgres?sslmode=disable" up 1