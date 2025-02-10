package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (t *impl) pageNoRoute(ctx *gin.Context) {
	statusCode := http.StatusNotFound
	ctx.HTML(
		statusCode,
		"page_error.html",
		pageError{
			ComponentCommonHead: componentCommonHead{
				Title:              "error",
				GoogleTagManagerID: t.googleTagManagerID,
			},
			ComponentHeader: componentHeader{
				ctxGetAdmin(ctx),
			},
			Message: fmt.Sprintf("error %d", statusCode),
		},
	)
}
