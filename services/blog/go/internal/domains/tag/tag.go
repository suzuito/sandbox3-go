package tag

import (
	"github.com/google/uuid"
)

type Tags []*Tag

type IDs []ID

type ID uuid.UUID

func (t ID) String() string {
	return uuid.UUID(t).String()
}

func NewIDFromUUID(id uuid.UUID) ID {
	return ID(id)
}

type Tag struct {
	ID   ID
	Name string `validate:"required"`
}

/*
func (t *Tag) Validate() error {
	if err := validate.Struct(t); err != nil {
		return terrors.Errorf("tag is invalid: %w", err)
	}

	return nil
}
*/
