package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/admin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/errors/stderror"
)

func middlewareXRobotsTag(ctx *gin.Context) {
	ctx.Header("X-Robots-Tag", "noindex")
	ctx.Next()
}

func (t *impl) middlewareAdminAuthn(ctx *gin.Context) {
	sessionIDString, err := ctx.Cookie("session")
	if err != nil {
		return
	}

	sessionID, err := admin.NewLoginSessionIDFromString(sessionIDString)
	if err != nil {
		return
	}

	if _, err := t.loginUsecase.AuthnAdmin(ctx, sessionID); err != nil {
		return
	}

	ctxSetAdmin(ctx)

	ctx.Next()
}

func (t *impl) middlewareAdminOnly(ctx *gin.Context) {
	if !ctxGetAdmin(ctx) {
		t.pageError(ctx, stderror.NewNotFound("not admin"))
		return
	}

	ctx.Next()
}

func (t *impl) middlewareSetCookieSameSiteStrict(ctx *gin.Context) {
	ctx.SetSameSite(http.SameSiteStrictMode)
	ctx.Next()
}
