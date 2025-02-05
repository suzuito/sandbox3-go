package sqlcgo

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (t *CreateTagsParams) String() string {
	return fmt.Sprintf("%+v", *t)
}

type CreateTagsParamsList []CreateTagsParams

func NewCreateTagsParamsListAtRandom(
	beginI int,
	beginPublishedAt time.Time,
	n int,
) CreateTagsParamsList {
	r := CreateTagsParamsList{}

	for i := range n {
		j := beginI + i
		r = append(r, CreateTagsParams{
			ID:        uuid.New(),
			Name:      fmt.Sprintf("テストタグ%d", j),
			CreatedAt: NewPgTypeFromTime(beginPublishedAt.Add(time.Hour * time.Duration(24*i))),
			UpdatedAt: NewPgTypeFromTime(beginPublishedAt.Add(time.Hour * time.Duration(24*i))),
			DeletedAt: NewPgTypeFromTimePtr(nil),
		})
	}

	return r
}
