-- name: ReadArticlesByIDs :many
SELECT id, title, published_at FROM articles WHERE id = ANY($1::uuid[]) AND deleted_at IS NULL;
