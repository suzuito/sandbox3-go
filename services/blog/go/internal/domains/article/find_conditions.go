package article

import (
	"net/url"
	"strconv"
	"time"
)

const (
	defaultFindConditionPage  = 0
	defaultFindConditionCount = 10
)

type FindConditions struct {
	TagName          *string
	PublishedAtRange FindConditionRange

	Page  uint16
	Count uint16
}

func (t *FindConditions) Next() *FindConditions {
	return &FindConditions{
		TagName: t.TagName,
		PublishedAtRange: FindConditionRange{
			Since: t.PublishedAtRange.Since,
			Until: t.PublishedAtRange.Until,
		},
		Page:  t.Page + 1,
		Count: t.Count,
	}
}

func (t *FindConditions) Prev() *FindConditions {
	return &FindConditions{
		TagName: t.TagName,
		PublishedAtRange: FindConditionRange{
			Since: t.PublishedAtRange.Since,
			Until: t.PublishedAtRange.Until,
		},
		Page:  t.Page - 1,
		Count: t.Count,
	}
}

func (t *FindConditions) Offset() uint64 {
	return uint64(t.Page) * uint64(t.Count)
}

func (t *FindConditions) ParseQuery(q url.Values) {
	if v := q.Get("tag"); v != "" {
		t.TagName = &v
	}
	if v, err := strconv.ParseUint(q.Get("page"), 10, 64); err == nil {
		t.Page = uint16(v)
	}
	if v, err := strconv.ParseUint(q.Get("limit"), 10, 64); err == nil {
		t.Count = uint16(v)
	}
	if v, err := strconv.ParseInt(q.Get("since"), 10, 64); err == nil {
		since := time.Unix(v, 0)
		t.PublishedAtRange.Since = &since
	}
	if v, err := strconv.ParseInt(q.Get("until"), 10, 64); err == nil {
		until := time.Unix(v, 0)
		t.PublishedAtRange.Until = &until
	}
}

func (t *FindConditions) Query() url.Values {
	q := url.Values{}

	if t.TagName != nil {
		q.Set("tag", *t.TagName)
	}
	q.Set("page", strconv.FormatUint(uint64(t.Page), 10))
	if t.Count != defaultFindConditionCount {
		q.Set("limit", strconv.FormatUint(uint64(t.Count), 10))
	}
	if t.PublishedAtRange.Since != nil {
		q.Set("since", strconv.FormatInt(t.PublishedAtRange.Since.Unix(), 10))
	}
	if t.PublishedAtRange.Until != nil {
		q.Set("until", strconv.FormatInt(t.PublishedAtRange.Until.Unix(), 10))
	}

	return q
}

func (t *FindConditions) URL() *url.URL {
	u, _ := url.Parse("")
	u.RawQuery = t.Query().Encode()
	return u
}

func newDefaultFindConditions() *FindConditions {
	fd := FindConditions{
		PublishedAtRange: FindConditionRange{},
		Page:             defaultFindConditionPage,
		Count:            defaultFindConditionCount,
	}
	return &fd
}

func NewFindConditionsFromQuery(q url.Values) *FindConditions {
	fd := newDefaultFindConditions()
	fd.ParseQuery(q)
	return fd
}

type FindConditionRange struct {
	Since *time.Time
	Until *time.Time
}

func (t *FindConditionRange) IsUsed() bool {
	return t.Since != nil || t.Until != nil
}
