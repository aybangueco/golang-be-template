-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password, user_type)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, first_name, last_name, email, user_type, created_at, updated_at;