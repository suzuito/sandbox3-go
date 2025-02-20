package usecases

import (
	"context"

	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/admin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
)

type ArticleUsecase interface {
	// Create(ctx context.Context) (domains.ArticleID, error)
	FindArticles(
		ctx context.Context,
		cond *article.FindConditions,
	) (articles article.Articles, next *article.FindConditions, prev *article.FindConditions, err error)
}

type LoginUsecase interface {
	LoginAsAdmin(ctx context.Context, inputPassword admin.PasswordAsPlainText) (*admin.LoginSession, error)
}
