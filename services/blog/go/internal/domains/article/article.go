package article

import (
	"time"

	"github.com/google/uuid"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
)

type IDs []ID

func (t IDs) ToUUIDs() []uuid.UUID {
	r := make([]uuid.UUID, 0, len(t))
	for _, id := range t {
		r = append(r, uuid.UUID(id))
	}
	return r
}

func NewIDsFromUUIDs(ids []uuid.UUID) IDs {
	r := make(IDs, 0, len(ids))
	for _, id := range ids {
		r = append(r, ID(id))
	}
	return r
}

type ID uuid.UUID

type Article struct {
	ID          ID
	Title       string
	PublishedAt *time.Time
}

func (t *Article) ValidateAsDraft() error {
	if err := validate.Struct(t); err != nil {
		return terrors.Wrap(err)
	}

	return nil
}

func (t *Article) ValidateAsPublished() error {
	if err := validate.Struct(t); err != nil {
		return terrors.Errorf("article is invalid: %w", err)
	}

	if t.Title == "" {
		return terrors.Errorf("title is required as published article")
	}

	return nil
}

func (t *Article) IsDraft() bool {
	return t.PublishedAt == nil
}
