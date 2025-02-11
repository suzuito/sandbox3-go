package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
)

type pageGETAdminArticles struct {
	ComponentHeader     componentHeader
	ComponentCommonHead componentCommonHead
	Breadcrumbs         breadcrumbs
}

func (t *impl) pageGETAdminArticles(ctx *gin.Context) {
	conds := article.NewFindConditionsFromQuery(ctx.Request.URL.Query())

	articles, next, prev, err := t.articleUsecase.FindArticles(ctx, conds)
	if err != nil {
		t.pageError(ctx, err)
		return
	}

	obj := pageGETArticles{
		ComponentHeader: componentHeader{
			IsAdmin: ctxGetAdmin(ctx),
		},
		ComponentCommonHead: componentCommonHead{
			Title:              "記事一覧",
			Meta:               nil,
			GoogleTagManagerID: t.googleTagManagerID,
		},
		Breadcrumbs: breadcrumbs{
			{
				Path: "/",
				URL:  newPageURL(t.siteOrigin, "/"),
				Name: "トップページ",
			},
			{
				Name:   "記事一覧",
				NoLink: true,
			},
		},
		Articles:                  articles,
		ComponentArticleListPager: componentArticleListPager{},
	}

	if next != nil {
		obj.ComponentArticleListPager.NextURL = next.URL()
	}
	if prev != nil {
		obj.ComponentArticleListPager.PrevURL = prev.URL()
	}

	ctx.HTML(
		http.StatusOK,
		"page_articles.html",
		&obj,
	)
}
