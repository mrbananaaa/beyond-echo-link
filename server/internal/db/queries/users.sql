-- name: CreateUser :one
INSERT INTO users (
  id, email, username, password, lookup_id, bio, profile_picture
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUserByID :one
SELECT
  id, email, username, password, lookup_id, bio, profile_picture, created_at, updated_at
FROM
  users
WHERE
  id = @id;

-- name: GetUserByLookupID :one
SELECT
  id, email, username, password, lookup_id, bio, profile_picture, created_at, updated_at
FROM
  users
WHERE
  lookup_id = @lookup_id;
