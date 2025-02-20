package admin

import (
	"context"

	domainadmin "github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/admin"
)

type saltRepository struct {
	salt domainadmin.Salt
}

func (t *saltRepository) Get(_ context.Context) (domainadmin.Salt, error) {
	return t.salt, nil
}

func NewSaltRepository(salt domainadmin.Salt) *saltRepository {
	return &saltRepository{
		salt: salt,
	}
}
