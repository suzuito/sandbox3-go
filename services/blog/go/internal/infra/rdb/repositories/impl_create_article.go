package repositories

import (
	"context"
	"fmt"

	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
)

func (t *impl) CreateArticle(ctx context.Context) (article.ID, error) {
	return article.ID{}, fmt.Errorf("not impl")
}
