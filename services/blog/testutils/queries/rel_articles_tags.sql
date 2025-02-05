-- name: CreateRelArticlesTags :copyfrom
INSERT INTO rel_articles_tags(article_id, tag_id)
VALUES ($1,$2);
