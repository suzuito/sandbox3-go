package usecases

import (
	"context"

	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/admin"
)

func (t *impl) LoginAsAdmin(ctx context.Context, inputPassword admin.PasswordAsPlainText) (*admin.LoginSession, error) {
	return t.adminService.Login(ctx, inputPassword)
}
