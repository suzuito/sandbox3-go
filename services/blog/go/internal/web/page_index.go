package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type pageIndex struct {
	ComponentHeader     componentHeader
	ComponentCommonHead componentCommonHead
}

func (t *impl) pageIndex(ctx *gin.Context) {
	ctx.HTML(
		http.StatusOK,
		"page_index.html",
		pageIndex{
			ComponentHeader: componentHeader{
				IsAdmin: ctxGetAdmin(ctx),
			},
			ComponentCommonHead: componentCommonHead{
				Title: SiteName,
				Meta: &siteMetaData{
					OGP: ogpData{
						Title:       SiteName,
						Description: "とあるエンジニアの日記",
						Locale:      "ja_JP",
						Type:        "website",
						URL:         newPageURLFromContext(ctx, t.siteOrigin),
						SiteName:    SiteName,
						Image:       newPageURL(t.siteOrigin, "images/avatar.jpg"),
					},
				},
				GoogleTagManagerID: t.googleTagManagerID,
			},
		},
	)
}
