package sqlcgo

import "github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/admin"

func (t *ReadAdminLoginSessionByIDRow) ToAdminLoginSession() *admin.LoginSession {
	return &admin.LoginSession{
		ID:        admin.LoginSessionID(t.ID),
		ExpiredAt: t.ExpiredAt.Time,
	}
}
