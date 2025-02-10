-- name: ReadTagIDByName :one
SELECT id FROM tags WHERE name = $1 AND deleted_at IS NULL;
