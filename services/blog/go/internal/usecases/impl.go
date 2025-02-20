package usecases

import (
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/admin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
)

type impl struct {
	adminService   *admin.Service
	articleService *article.Service
}

func NewImpl(
	adminService *admin.Service,
	articleService *article.Service,
) *impl {
	return &impl{
		adminService:   adminService,
		articleService: articleService,
	}
}
