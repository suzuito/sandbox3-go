package web

import (
	"github.com/gin-gonic/gin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/admin"
)

func SetSessionCookie(ctx *gin.Context, session *admin.LoginSession) {
	ctx.SetCookie(
		"session",
		session.ID.UUID().String(),
		86400,
		"/",
		"",
		false,
		false,
	)
}
