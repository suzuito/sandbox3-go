package repositories

import (
	"github.com/jackc/pgx/v5"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/admin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
)

type impl struct {
	conn *pgx.Conn
}

var (
	_ article.Repository      = &impl{}
	_ admin.SessionRepository = &impl{}
)

func NewImpl(conn *pgx.Conn) *impl {
	return &impl{
		conn: conn,
	}
}
