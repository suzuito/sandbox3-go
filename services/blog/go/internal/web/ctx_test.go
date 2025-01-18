package web

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCtxAdmin(t *testing.T) {
	ctx := gin.CreateTestContextOnly(
		httptest.NewRecorder(),
		gin.New(),
	)
	assert.False(t, ctxGetAdmin(ctx))
	ctxSetAdmin(ctx)
	assert.True(t, ctxGetAdmin(ctx))
}
