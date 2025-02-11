package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
)

type pageGETArticles struct {
	ComponentHeader           componentHeader
	ComponentCommonHead       componentCommonHead
	Articles                  article.Articles
	ComponentArticleListPager componentArticleListPager
	Breadcrumbs               breadcrumbs
}

// TODO 後でsandbox2-common-goへ移す
func GetQueryAsUint[T uint | uint8 | uint16 | uint32 | uint64](ctx *gin.Context, key string) (T, bool) {
	value, exists := ctx.GetQuery(key)
	if !exists {
		return 0, false
	}

	i, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, false
	}

	return T(i), true
}

// TODO 後でsandbox2-common-goへ移す
func GetQueryAsInt[T int | int8 | int16 | int32 | int64](ctx *gin.Context, key string) (T, bool) {
	value, exists := ctx.GetQuery(key)
	if !exists {
		return 0, false
	}

	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, false
	}

	return T(i), true
}

// TODO 後でsandbox2-common-goへ移す
func GetQueryAsTimestamp(ctx *gin.Context, key string) (time.Time, bool) {
	t, exists := GetQueryAsInt[int64](ctx, key)
	if !exists {
		return time.Unix(0, 0), false
	}

	return time.Unix(t, 0), true
}

func (t *impl) pageGETArticles(ctx *gin.Context) {
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
