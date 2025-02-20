package admin

import (
	"context"

	domainadmin "github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/admin"
)

type passwordRepository struct {
	passwordHash domainadmin.PasswordAsHash
}

func (t *passwordRepository) Get(_ context.Context) (domainadmin.PasswordAsHash, error) {
	return t.passwordHash, nil
}

func NewPasswordRepository(passwordHash domainadmin.PasswordAsHash) *passwordRepository {
	return &passwordRepository{
		passwordHash: passwordHash,
	}
}
