package admin

import "context"

type PasswordRepository interface {
	Get(ctx context.Context) (PasswordAsHash, error)
}

type SaltRepository interface {
	Get(ctx context.Context) (Salt, error)
}

type SessionRepository interface {
	CreateSession(ctx context.Context, id LoginSessionID) (*LoginSession, error)
	ReadSession(ctx context.Context, id LoginSessionID) (*LoginSession, error)
}
