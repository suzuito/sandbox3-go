package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/stderror"
)

type pageError struct {
	ComponentCommonHead componentCommonHead
	ComponentHeader     componentHeader
	Message             string
}

func (t *impl) pageError(ctx *gin.Context, err error) {
	defer ctx.Abort()

	t.logger.Error(err.Error(), "err", err)

	code := stderror.ToCode(err)
	httpStatusCode := code.HTTPStatusCode()

	ctx.HTML(
		httpStatusCode,
		"page_error.html",
		pageError{
			ComponentCommonHead: componentCommonHead{
				Title:              "error",
				GoogleTagManagerID: t.googleTagManagerID,
			},
			ComponentHeader: componentHeader{
				ctxGetAdmin(ctx),
			},
			Message: fmt.Sprintf("error %d", httpStatusCode),
		},
	)
}
