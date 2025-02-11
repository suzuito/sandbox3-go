package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type pageGETAdmin struct {
	ComponentHeader     componentHeader
	ComponentCommonHead componentCommonHead
	Breadcrumbs         breadcrumbs
}

func (t *impl) pageGETAdmin(ctx *gin.Context) {
	obj := pageGETArticles{
		ComponentHeader: componentHeader{
			IsAdmin: ctxGetAdmin(ctx),
		},
		ComponentCommonHead: componentCommonHead{
			Title:              "管理画面",
			GoogleTagManagerID: t.googleTagManagerID,
		},
		Breadcrumbs: breadcrumbs{
			{
				Path: "/",
				URL:  newPageURL(t.siteOrigin, "/"),
				Name: "トップページ",
			},
			{
				Name:   "管理画面",
				NoLink: true,
			},
		},
	}

	ctx.HTML(
		http.StatusOK,
		"page_admin.html",
		&obj,
	)
}
