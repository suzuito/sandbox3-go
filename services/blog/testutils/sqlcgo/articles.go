package sqlcgo

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func (t *CreateArticlesParams) String() string {
	return fmt.Sprintf("%+v", *t)
}

type CreateArticlesParamsList []CreateArticlesParams

func NewCreateArticlesParamsListAtRandom(
	beginI int,
	beginPublishedAt time.Time,
	n int,
) CreateArticlesParamsList {
	r := CreateArticlesParamsList{}

	for i := range n {
		publishedAt := pgtype.Timestamp{}
		if i%5 == 0 {
			publishedAt.Valid = false
		} else {
			publishedAt = NewPgTypeFromTime(beginPublishedAt.Add(time.Hour * time.Duration(24*i)))
		}
		j := beginI + i
		r = append(r, CreateArticlesParams{
			ID:          uuid.New(),
			Title:       fmt.Sprintf("テスト記事%d", j),
			PublishedAt: publishedAt,
			CreatedAt:   NewPgTypeFromTime(beginPublishedAt.Add(time.Hour * time.Duration(24*i))),
			UpdatedAt:   NewPgTypeFromTime(beginPublishedAt.Add(time.Hour * time.Duration(24*i))),
			DeletedAt:   NewPgTypeFromTimePtr(nil),
		})
	}

	return r
}
