package repositories

import (
	"context"

	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/infra/rdb/sqlc/sqlcgo"
)

func (t *impl) ReadArticles(ctx context.Context, ids article.IDs) (article.Articles, error) {
	queries := sqlcgo.New(t.conn)

	rows, err := queries.ReadArticlesByIDs(ctx, ids.ToUUIDs())
	if err != nil {
		return nil, terrors.Errorf("failed to read articles: %w", err)
	}

	return sqlcgo.ReadArticlesByIDsRows(rows).ToArticles(), nil
}
