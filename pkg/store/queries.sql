-- name: Create :one
INSERT INTO roach_migrations (id, key, filename) VALUES ($1, $2, $3) RETURNING *;

-- name: List :many
SELECT * FROM roach_migrations WHERE key = $1 ORDER BY created_at;

-- name: ListForUpdate :many
SELECT * FROM roach_migrations WHERE key = $1 ORDER BY created_at FOR UPDATE;

-- name: Update :one
UPDATE roach_migrations SET completed = $1, failed = $2, fail_reason = $3 WHERE id = $4 RETURNING *;