package usecases

import (
	"context"

	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/admin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/errors/repoerror"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/errors/stderror"
)

func (t *impl) LoginAsAdmin(ctx context.Context, inputPassword admin.PasswordAsPlainText) (*admin.LoginSession, error) {
	return t.adminService.Login(ctx, inputPassword)
}

func (t *impl) AuthnAdmin(ctx context.Context, id admin.LoginSessionID) (*admin.LoginSession, error) {
	session, err := t.adminService.Authn(ctx, id)
	switch {
	case repoerror.Has(err, repoerror.CodeNotFound):
		return nil, terrors.Wrap(stderror.NewUnauthorized("session is not found"))
	case err != nil:
		return nil, terrors.Wrap(err)
	}

	return session, nil
}
