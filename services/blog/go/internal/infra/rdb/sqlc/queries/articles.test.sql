-- name: CreateArticlesForTest :copyfrom
INSERT INTO articles (id, title, published, published_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?);
