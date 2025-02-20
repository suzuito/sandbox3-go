package admin

import (
	"time"

	"github.com/google/uuid"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
)

type LoginSessionID uuid.UUID

func (t LoginSessionID) UUID() uuid.UUID {
	return uuid.UUID(t)
}

func NewLoginSessionID() LoginSessionID {
	return LoginSessionID(uuid.New())
}

func NewLoginSessionIDFromString(s string) (LoginSessionID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return LoginSessionID{}, terrors.Errorf("invalid uuid: %w", err)
	}
	return LoginSessionID(id), nil
}

type LoginSession struct {
	ID        LoginSessionID
	ExpiredAt time.Time
}
