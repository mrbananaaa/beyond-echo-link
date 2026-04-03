-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  username VARCHAR(21) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  lookup_id VARCHAR(12) UNIQUE NOT NULL,
  bio VARCHAR(255),
  profile_picture VARCHAR(255),
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
