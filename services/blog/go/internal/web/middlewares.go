package web

import "github.com/gin-gonic/gin"

func middlewareXRobotsTag(ctx *gin.Context) {
	ctx.Header("X-Robots-Tag", "noindex")
	ctx.Next()
}

func (t *impl) middlewareAdminAuthe(ctx *gin.Context) {
	token, err := ctx.Cookie("admin_auth_token")
	if err != nil {
		return
	}
	if token != t.adminToken {
		return
	}
	ctxSetAdmin(ctx)
}
