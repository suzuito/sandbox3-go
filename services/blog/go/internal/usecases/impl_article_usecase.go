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
	conds.ExcludeDraft = true

	articleIDs, err := t.articleRepository.FindArticles(ctx, conds)
	if err != nil {
		return nil, nil, nil, terrors.Wrap(err)
	}

	if len(articleIDs) <= 0 {
		return article.Articles{}, nil, nil, nil
	}

	articles, err := t.articleRepository.ReadArticles(ctx, articleIDs)
	if err != nil {
		return nil, nil, nil, terrors.Wrap(err)
	}

	if len(articleIDs) != len(articles) {
		return nil, nil, nil, terrors.Errorf("some article ids are not found")
	}

	var next *article.FindConditions
	if len(articles) >= int(conds.Count) {
		next = conds.Next()
	}

	var prev *article.FindConditions
	if conds.Page > 0 {
		prev = conds.Prev()
	}

	return articles, next, prev, nil
}
