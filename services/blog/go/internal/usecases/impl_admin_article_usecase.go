package usecases

import (
	"context"

	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
)

func (t *impl) FindAdminArticles(ctx context.Context, conds *article.FindConditions) (
	article.Articles,
	*article.FindConditions,
	*article.FindConditions,
	error,
) {
	return t.FindArticles(ctx, conds)
}
