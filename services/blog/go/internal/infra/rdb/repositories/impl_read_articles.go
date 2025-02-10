package repositories

import (
	"context"

	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/infra/rdb/sqlc/sqlcgo"
)

func (t *impl) ReadArticles(ctx context.Context, ids article.IDs) (article.Articles, error) {
	if len(ids) <= 0 {
		return nil, terrors.Errorf("ids are required")
	}

	queries := sqlcgo.New(t.conn)

	rows, err := queries.ReadArticlesByIDs(ctx, ids.ToUUIDs())
	if err != nil {
		return nil, terrors.Errorf("failed to read articles: %w", err)
	}

	if len(rows) < len(ids) {
		return nil, terrors.Errorf("some articles are not found")
	}

	articlesByID := sqlcgo.ReadArticlesByIDsRows(rows).ToArticles().GroupByID()

	articles := make(article.Articles, 0, len(articlesByID))
	for _, id := range ids {
		articlesPerID := articlesByID[id]
		articles = append(articles, articlesPerID...)
	}

	return articles, nil
}
