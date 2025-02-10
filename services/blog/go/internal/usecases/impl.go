package usecases

import "github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"

type impl struct {
	articleRepository article.Repository
}

func NewImpl(
	articleRepository article.Repository,
) *impl {
	return &impl{
		articleRepository: articleRepository,
	}
}
