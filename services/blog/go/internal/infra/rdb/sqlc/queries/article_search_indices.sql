-- name: UpseartArticleSearchIndices :exec
WITH t1 AS (
  SELECT
    articles.id AS article_id,
    array_agg(rel_articles_tags.tag_id) AS tag_ids,
    articles.published_at AS published_at
  FROM articles
  LEFT JOIN rel_articles_tags ON articles.id = rel_articles_tags.article_id
  WHERE
    articles.deleted_at IS NULL
  GROUP BY articles.id, articles.published_at
  ORDER BY articles.created_at
  LIMIT 1000 OFFSET 0
)
INSERT INTO articles_search_indices(article_id, tag_ids, published_at)
SELECT article_id, tag_ids, published_at FROM t1
ON CONFLICT(article_id)
DO UPDATE SET
  tag_ids = EXCLUDED.tag_ids,
  published_at = EXCLUDED.published_at
;
