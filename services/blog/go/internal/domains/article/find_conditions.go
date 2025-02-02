package article

import (
	"time"

	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/tag"
)

type FindConditions struct {
	TagID            *tag.ID
	PublishedAtRange FindConditionRange

	Page  uint16
	Count uint16
}

func (t *FindConditions) Next() *FindConditions {
	return &FindConditions{
		TagID: t.TagID,
		PublishedAtRange: FindConditionRange{
			Since: t.PublishedAtRange.Since,
			Until: t.PublishedAtRange.Until,
		},
		Page:  t.Page + 1,
		Count: t.Count,
	}
}

func (t *FindConditions) Offset() uint64 {
	return uint64(t.Page) * uint64(t.Count)
}

type FindConditionRange struct {
	Since *time.Time
	Until *time.Time
}

func (t *FindConditionRange) IsUsed() bool {
	return t.Since != nil || t.Until != nil
}
