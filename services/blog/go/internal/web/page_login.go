package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/admin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/errors/stderror"
)

type pageGETLogin struct {
	ComponentHeader     componentHeader
	ComponentCommonHead componentCommonHead
}

func (t *impl) pageGETLogin(ctx *gin.Context) {
	ctx.HTML(
		http.StatusOK,
		"page_login.html",
		pageIndex{
			ComponentHeader: componentHeader{
				IsAdmin: ctxGetAdmin(ctx),
			},
			ComponentCommonHead: componentCommonHead{
				Title:              SiteName + " | ログイン",
				GoogleTagManagerID: t.googleTagManagerID,
			},
		},
	)
}

func (t *impl) pagePOSTLogin(ctx *gin.Context) {
	password := admin.PasswordAsPlainText(ctx.PostForm("password"))
	if password == "" {
		t.pageError(ctx, stderror.NewBadRequest("password is empty"))
		return
	}

	session, err := t.loginUsecase.LoginAsAdmin(ctx, password)
	if err != nil {
		t.pageError(ctx, err)
		return
	}
	ctx.SetCookie(
		"session",
		session.ID.UUID().String(),
		86400,
		"/",
		"",
		false,
		false,
	)

	ctx.Redirect(http.StatusFound, "/")
}

func (t *impl) pagePOSTLogout(ctx *gin.Context) {
	ctx.SetCookie(
		"session",
		"",
		0,
		"/",
		"",
		false,
		false,
	)

	ctx.Redirect(http.StatusFound, "/")
}
