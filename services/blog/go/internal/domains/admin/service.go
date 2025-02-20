package admin

import (
	"context"
	"fmt"

	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/errors/stderror"
)

type Service struct {
	saltRepository     SaltRepository
	passwordRepository PasswordRepository
	sessionRepository  SessionRepository
	hashFunc           HashFunc
}

func (t *Service) Login(
	ctx context.Context,
	inputPassword PasswordAsPlainText,
) (*LoginSession, error) {
	salt, err := t.saltRepository.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get salt: %w", err)
	}

	inputPasswordHash, err := t.hashFunc(inputPassword, salt)
	if err != nil {
		return nil, fmt.Errorf("failed to create password hash: %w", err)
	}

	passwordHash, err := t.passwordRepository.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get password hash: %w", err)
	}

	if inputPasswordHash != passwordHash {
		return nil, terrors.Wrap(stderror.NewUnauthorized("mismatch password"))
	}

	sessionID := NewLoginSessionID()
	session, err := t.sessionRepository.CreateSession(ctx, sessionID)
	if err != nil {
		return nil, terrors.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

func (t *Service) Authn(
	ctx context.Context,
	id LoginSessionID,
) (*LoginSession, error) {
	session, err := t.sessionRepository.ReadSession(ctx, id)
	if err != nil {
		return nil, terrors.Errorf("failed to read session: %w", err)
	}

	return session, nil
}

func NewService(
	saltRepository SaltRepository,
	passwordRepository PasswordRepository,
	sessionRepository SessionRepository,
	hashFunc HashFunc,
) *Service {
	return &Service{
		saltRepository:     saltRepository,
		passwordRepository: passwordRepository,
		sessionRepository:  sessionRepository,
		hashFunc:           hashFunc,
	}
}
