package usecases

import (
	"context"
	"fmt"

	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
)

func (t *impl) FindArticles(ctx context.Context, conds *article.FindConditions) (
	article.Articles,
	*article.FindConditions,
	error,
) {
	articleIDs, err := t.articleRepository.FindArticles(ctx, conds)
	if err != nil {
		return nil, nil, terrors.Wrap(err)
	}

	articles, err := t.articleRepository.ReadArticles(ctx, articleIDs)
	if err != nil {
		return nil, nil, terrors.Wrap(err)
	}

	fmt.Println("aaaaa")
	for _, a := range articles {
		fmt.Printf("%+v\n", a)
	}

	if len(articleIDs) != len(articles) {
		return nil, nil, terrors.Errorf("some article ids are not found")
	}

	var next *article.FindConditions
	if len(articles) >= int(conds.Count) {
		next = conds.Next()
	}

	return articles, next, nil
}
