package usecases

import (
	"context"

	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
)

func (t *impl) FindArticles(ctx context.Context, conds *article.FindConditions) (
	article.Articles,
	*article.FindConditions,
	*article.FindConditions,
	error,
) {
	return t.articleService.FindArticles(ctx, conds)
}

func (t *impl) CreateArticle(ctx context.Context) (article.ID, error) {
	art, err := t.articleService.CreateArticle(ctx)
	if err != nil {
		return article.ID{}, terrors.Wrap(err)
	}

	return art.ID, nil
}
