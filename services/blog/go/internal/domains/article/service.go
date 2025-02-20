package article

import (
	"context"

	"github.com/suzuito/sandbox2-common-go/libs/terrors"
)

type Service struct {
	articleRepository Repository
}

func NewService(
	articleRepository Repository,
) *Service {
	return &Service{
		articleRepository: articleRepository,
	}
}

func (t *Service) FindArticles(
	ctx context.Context,
	conds *FindConditions,
) (Articles, *FindConditions, *FindConditions, error) {
	articleIDs, err := t.articleRepository.FindArticles(ctx, conds)
	if err != nil {
		return nil, nil, nil, terrors.Wrap(err)
	}

	if len(articleIDs) <= 0 {
		return Articles{}, nil, nil, nil
	}

	articles, err := t.articleRepository.ReadArticles(ctx, articleIDs)
	if err != nil {
		return nil, nil, nil, terrors.Wrap(err)
	}

	if len(articleIDs) != len(articles) {
		return nil, nil, nil, terrors.Errorf("some article ids are not found")
	}

	var next *FindConditions
	if len(articles) >= int(conds.Count) {
		next = conds.Next()
	}

	var prev *FindConditions
	if conds.Page > 0 {
		prev = conds.Prev()
	}

	return articles, next, prev, nil
}

func (t *Service) CreateArticle(
	ctx context.Context,
) (*Article, error) {
	art := Article{
		ID: NewID(),
	}

	if err := t.articleRepository.CreateArticle(ctx, &art); err != nil {
		return nil, terrors.Wrap(err)
	}

	return &art, nil
}
