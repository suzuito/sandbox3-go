package usecases

import (
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/admin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
)

type impl struct {
	articleRepository article.Repository

	adminService *admin.Service
}

func NewImpl(
	articleRepository article.Repository,
	adminService *admin.Service,
) *impl {
	return &impl{
		articleRepository: articleRepository,
		adminService:      adminService,
	}
}
