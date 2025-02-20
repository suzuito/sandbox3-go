package usecases

import (
	"context"

	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/admin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
)

type ArticleUsecase interface {
	CreateArticle(ctx context.Context) (article.ID, error)
	FindArticles(
		ctx context.Context,
		cond *article.FindConditions,
	) (articles article.Articles, next *article.FindConditions, prev *article.FindConditions, err error)
}

type LoginUsecase interface {
	LoginAsAdmin(ctx context.Context, inputPassword admin.PasswordAsPlainText) (*admin.LoginSession, error)
	AuthnAdmin(ctx context.Context, id admin.LoginSessionID) (*admin.LoginSession, error)
}
