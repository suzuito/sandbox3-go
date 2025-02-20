package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/infra/rdb/sqlc/sqlcgo"
)

func (t *impl) CreateArticle(
	ctx context.Context,
	article *article.Article,
) error {
	queries := sqlcgo.New(t.conn)

	publishedAt := pgtype.Timestamp{}
	if article.PublishedAt != nil {
		publishedAt.Time = *article.PublishedAt
		publishedAt.Valid = true
	}

	if err := queries.CreateArticle(ctx, sqlcgo.CreateArticleParams{
		ID:          article.ID.UUID(),
		Title:       article.Title,
		PublishedAt: publishedAt,
	}); err != nil {
		return terrors.Wrap(err)
	}

	return nil
}
