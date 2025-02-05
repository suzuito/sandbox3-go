-- name: CreateArticles :copyfrom
INSERT INTO articles(id, title, published_at, created_at, updated_at, deleted_at)
VALUES ($1, $2, $3, $4, $5, $6);

