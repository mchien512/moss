-- go/internal/link/repository/db/queries/entry_links.sql

-- 1. Insert a new link between two entries
-- Returns the inserted row (so SQLC can map it to an EntryLink struct).
-- name: CreateEntryLink :one
INSERT INTO entry_links (
    source_entry_id,
    target_entry_id,
    user_id,
    created_at
) VALUES (
             $1,  -- source_entry_id
             $2,  -- target_entry_id
             $3,  -- user_id (who created/owns this link)
             CURRENT_TIMESTAMP
         )
RETURNING source_entry_id, target_entry_id, user_id, created_at;

-- 2. Delete a link (unlink two entries)
-- name: DeleteEntryLink :exec
DELETE FROM entry_links
WHERE source_entry_id = $1
  AND target_entry_id = $2;

-- 3. List all links where a given entry is the “source”
-- (i.e. all outgoing links from entry X)
-- name: ListLinksBySource :many
SELECT
    source_entry_id,
    target_entry_id,
    user_id,
    created_at
FROM entry_links
WHERE source_entry_id = $1
ORDER BY created_at;

-- 4. List all links where a given entry is the “target”
-- (i.e. all incoming/backlinks to entry X)
-- name: ListLinksByTarget :many
SELECT
    source_entry_id,
    target_entry_id,
    user_id,
    created_at
FROM entry_links
WHERE target_entry_id = $1
ORDER BY created_at;

-- 5. Count how many outgoing links a given entry has
-- (useful for setting “link_count” in your proto if you want outgoing count)
-- name: CountLinksBySource :one
SELECT COUNT(*) AS count
FROM entry_links
WHERE source_entry_id = $1;

-- 6. Count how many incoming links a given entry has
-- (useful for backlink counts)
-- name: CountLinksByTarget :one
SELECT COUNT(*) AS count
FROM entry_links
WHERE target_entry_id = $1;

-- 7. (Optional) List the actual Entry rows that a given source is linked to,
--     with pagination parameters (page size + offset). This is if you want to
--     fetch full Entry data in one go. Adjust the SELECT columns as needed.
-- name: ListLinkedEntries :many
SELECT
    e.id,
    e.user_id,
    e.title,
    e.content,
    e.growth_stage,
    e.created_at,
    e.updated_at
FROM entries AS e
         JOIN entry_links AS l
              ON l.target_entry_id = e.id
WHERE l.source_entry_id = $1
ORDER BY e.created_at
LIMIT $2      -- page_size
    OFFSET $3;    -- offset (page_token converted to integer)

-- 8. (Optional) List the actual Entry rows that link *into* a given entry,
--     with pagination. Adjust as above.
-- name: ListBacklinkedEntries :many
SELECT
    e.id,
    e.user_id,
    e.title,
    e.content,
    e.growth_stage,
    e.created_at,
    e.updated_at
FROM entries AS e
         JOIN entry_links AS l
              ON l.source_entry_id = e.id
WHERE l.target_entry_id = $1
ORDER BY e.created_at
LIMIT $2      -- page_size
    OFFSET $3;    -- offset (page_token)
