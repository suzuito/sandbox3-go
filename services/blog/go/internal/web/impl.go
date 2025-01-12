package web

import "github.com/gin-gonic/gin"

type impl struct {
}

func (t *impl) SetEngine(e *gin.Engine) {}

func New() *impl {
	return &impl{}
}
