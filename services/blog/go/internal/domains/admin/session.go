package admin

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type LoginSessionID uuid.UUID

func (t LoginSessionID) UUID() uuid.UUID {
	return uuid.UUID(t)
}

func NewLoginSessionID() LoginSessionID {
	return LoginSessionID(uuid.New())
}

type LoginSession struct {
	ID        LoginSessionID
	ExpiredAt time.Time
}

type SessionRepository interface {
	CreateSession(ctx context.Context, id LoginSessionID) (*LoginSession, error)
	ReadSession(ctx context.Context, id LoginSessionID) (*LoginSession, error)
}
