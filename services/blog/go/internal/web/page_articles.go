package web

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/tag"
)

type pageGETArticles struct {
	ComponentHeader           componentHeader
	ComponentCommonHead       componentCommonHead
	Articles                  article.Articles
	ComponentArticleListPager componentArticleListPager
	Breadcrumbs               breadcrumbs
}

// TODO 後でsandbox2-common-goへ移す
func DefaultQueryAsInt64(ctx *gin.Context, key string, dflt int64) int64 {
	value, exists := ctx.GetQuery(key)
	if !exists {
		return dflt
	}

	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return dflt
	}

	return i
}

// TODO 後でsandbox2-common-goへ移す
func DefaultQueryAsUint64(ctx *gin.Context, key string, dflt uint64) uint64 {
	value, exists := ctx.GetQuery(key)
	if !exists {
		return dflt
	}

	i, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return dflt
	}

	return i
}

func DefaultQueryAsUint16(ctx *gin.Context, key string, dflt uint16) uint16 {
	value, exists := ctx.GetQuery(key)
	if !exists {
		return dflt
	}

	i, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return dflt
	}

	if i > math.MaxUint16 {
		return dflt
	}

	return uint16(i)
}

func (t *impl) pageGETArticles(ctx *gin.Context) {
	var tagID *tag.ID
	tagIDAsString := ctx.DefaultQuery("tagId", "")
	if tagIDAsString != "" {
		tid, err := tag.NewIDFromString(tagIDAsString)
		if err == nil {
			ttid := tag.ID(tid)
			tagID = &ttid
		}
	}

	conds := article.FindConditions{
		TagID: tagID,
		Page:  DefaultQueryAsUint16(ctx, "page", 0),
		Count: uint16(DefaultQueryAsUint64(ctx, "limit", 10)),
		PublishedAtRange: *article.NewFindConditionRangeFromTimestamp(
			DefaultQueryAsInt64(ctx, "since", -1),
			DefaultQueryAsInt64(ctx, "until", -1),
		),
	}
	articles, next, err := t.articleUsecase.FindArticles(ctx, &conds)
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
		obj.ComponentArticleListPager.NextPage = &next.Page
	}
	if conds.Page >= 1 {
		prevPage := conds.Page - 1
		obj.ComponentArticleListPager.PrevPage = &prevPage
	}

	ctx.HTML(
		http.StatusOK,
		"page_articles.html",
		&obj,
	)
}
