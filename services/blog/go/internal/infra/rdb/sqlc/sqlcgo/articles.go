package sqlcgo

import (
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/tag"
)

type ReadArticlesByIDsRows []ReadArticlesByIDsRow

func (t ReadArticlesByIDsRows) ToArticles() article.Articles {
	articles := article.Articles{}
	for _, row := range t {
		articles = append(articles, row.ToArticle())
	}
	return articles
}

func (t *ReadArticlesByIDsRow) ToArticle() *article.Article {
	a := article.Article{
		ID:    article.ID(t.ID),
		Title: t.Title,
		Tags:  make(tag.Tags, 0, len(t.TagIds)),
	}

	for i, tagID := range t.TagIds {
		a.Tags = append(a.Tags, &tag.Tag{
			ID:   tag.NewIDFromUUID(tagID),
			Name: t.TagNames[i],
		})
	}

	if t.PublishedAt.Valid {
		a.PublishedAt = &t.PublishedAt.Time
	}

	return &a
}
