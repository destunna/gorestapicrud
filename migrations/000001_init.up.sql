CREATE TABLE users (
  id SERIAL PRIMARY KEY
  full_name VARCHAR(42) NOT NULL
  age VARCHAR(3) NOT NULL
  phone_number (20) NOT NULL
  habits VARCHAR(200)
  alive BOOLEAN NOT NULL
  created_at TIMESTAMP NOTE NULL
)

-- migrate -path migrations -database "postgres://postgres@localhost:5432/postgres?sslmode=disable" up 1