-- name: CreateEntry :one
INSERT INTO entries (
    id,
    user_id,
    title,
    content,
    mood,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING *;

-- name: GetEntryByID :one
SELECT * FROM entries
WHERE id = $1;

-- name: ListEntriesByUser :many
SELECT * FROM entries
WHERE user_id = $1
ORDER BY created_at;

-- name: ListEntriesByUserSince :many
SELECT * FROM entries
WHERE user_id = $1 AND updated_at > $2
ORDER BY updated_at;

-- name: UpdateEntry :one
UPDATE entries
SET title = $2,
    content = $3,
    mood = $4,
    updated_at = $5
WHERE id = $1
    RETURNING *;

-- name: DeleteEntry :exec
DELETE FROM entries
WHERE id = $1;
