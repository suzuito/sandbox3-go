-- name: ReadArticlesByIDs :many
SELECT
  articles.id AS id,
  articles.title AS title,
  articles.published_at AS published_at,
  array_agg(tags.id)::uuid[] AS tag_ids,
  array_agg(tags.name)::text[] AS tag_names
FROM articles
LEFT JOIN rel_articles_tags ON articles.id = rel_articles_tags.article_id
LEFT JOIN tags ON tags.id = rel_articles_tags.tag_id AND tags.deleted_at IS NULL
WHERE articles.id = ANY($1::uuid[]) AND articles.deleted_at IS NULL
GROUP BY articles.id, articles.title, articles.published_at
;
