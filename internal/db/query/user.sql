-- name: getUserByEmail :one
SELECT id, name, email, password_hash
FROM users
WHERE email = $1 AND is_active = true;

-- name: getUserByID :one
SELECT id, name, email, password_hash
FROM users
WHERE id = $1 AND is_active = true;

-- name: createUser :one
INSERT INTO users (name, email, password_hash, is_active)
VALUES ($1, $2, $3, true)
RETURNING id, name, email, password_hash;