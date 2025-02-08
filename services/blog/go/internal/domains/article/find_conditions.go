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

func NewFindConditionRangeFromTimestamp(
	since int64,
	until int64,
) *FindConditionRange {
	r := FindConditionRange{}

	if since > 0 {
		sinceTime := time.Unix(since, 0)
		r.Since = &sinceTime
	}

	if until > 0 {
		untilTime := time.Unix(until, 0)
		r.Until = &untilTime
	}

	return &r
}
