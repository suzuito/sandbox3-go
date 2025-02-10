package repositories

import "github.com/jackc/pgx/v5"

type impl struct {
	conn *pgx.Conn
}

func NewImpl(conn *pgx.Conn) *impl {
	return &impl{
		conn: conn,
	}
}
