package article

import "context"

type Repository interface {
	CreateArticle(ctx context.Context, article *Article) error
	UpdateArticle(ctx context.Context, article *Article) error
	ReadArticles(ctx context.Context, ids IDs) (Articles, error)
	FindArticles(ctx context.Context, conds *FindConditions) (IDs, error)
}
