package web

import (
	"github.com/gin-gonic/gin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains"
)

type pageArticles struct {
	ComponentHeader           componentHeader
	ComponentCommonHead       componentCommonHead
	Articles                  domains.Articles
	ComponentArticleListPager componentArticleListPager
	Breadcrumbs               breadcrumbs
}

func (t *impl) pageArticles(ctx *gin.Context) {
}
