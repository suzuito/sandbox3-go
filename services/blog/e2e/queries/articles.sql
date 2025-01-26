-- name: CreateArticlesForTest :copyfrom
INSERT INTO articles (id, title, published, published_at, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6);
