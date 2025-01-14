package web

import "github.com/gin-gonic/gin"

var ctxKeyAdmin = "admin"

func ctxSetAdmin(ctx *gin.Context) {
	ctx.Set(ctxKeyAdmin, true)
}

func ctxGetAdmin(ctx *gin.Context) bool {
	v, ok := ctx.Get(ctxKeyAdmin)
	if !ok {
		return false
	}
	vv, ok := v.(bool)
	if !ok {
		return false
	}
	return vv
}
