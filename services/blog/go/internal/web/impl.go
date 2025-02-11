package web

import (
	"log/slog"
	"net/http"
	"net/url"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/inject"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/usecases"
)

type impl struct {
	siteOrigin          url.URL
	googleTagManagerID  string
	adminToken          string
	dirPathHTMLTemplate string
	dirPathCSS          string

	logger *slog.Logger

	articleUsecase usecases.ArticleUsecase
}

func (t *impl) SetEngine(e *gin.Engine) {
	e.NoRoute(middlewareXRobotsTag, t.pageNoRoute)

	e.GET(
		"health",
		middlewareXRobotsTag,
		func(ctx *gin.Context) { ctx.JSON(http.StatusOK, struct{}{}) },
	)
	e.LoadHTMLGlob(path.Join(t.dirPathHTMLTemplate, "*.html"))
	e.Static("css", t.dirPathCSS)
	e.Use(t.middlewareAdminAuthn)

	e.GET("", t.pageIndex)
	{
		gArticles := e.Group("articles")
		gArticles.GET("", t.pageGETArticles)
	}

	{
		gAdmin := e.Group("admin")
		gAdmin.Use(middlewareXRobotsTag)
		gAdmin.GET("", t.pageGETAdmin)
		{
			gAdminArticles := gAdmin.Group("articles")
			gAdminArticles.GET("", t.pageGETAdminArticles)
		}
	}
}

func New(
	env *inject.Environment,
	logger *slog.Logger,
	articleUsecase usecases.ArticleUsecase,
) (*impl, error) {
	urlSiteOrigin, err := url.Parse(env.SiteOrigin)
	if err != nil {
		return nil, terrors.Errorf("failed to url.Parse: %w", err)
	}

	return &impl{
		siteOrigin:          *urlSiteOrigin,
		googleTagManagerID:  env.GoogleTagManagerID,
		adminToken:          env.AdminToken,
		dirPathHTMLTemplate: env.DirPathHTMLTemplate,
		dirPathCSS:          env.DirPathCSS,
		logger:              logger,
		articleUsecase:      articleUsecase,
	}, nil
}
