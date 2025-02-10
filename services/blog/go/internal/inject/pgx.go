package inject

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
)

type DBEnvironment interface {
	DBURI() string
}

func NewPgxConn(ctx context.Context, env DBEnvironment) (*pgx.Conn, error) {
	conf, err := pgx.ParseConfig(env.DBURI())
	if err != nil {
		return nil, terrors.Errorf("failed to pgx.ParseConfig: %w", err)
	}

	conn, err := pgx.ConnectConfig(ctx, conf)
	if err != nil {
		return nil, terrors.Errorf("failed to pgx.ConnectConfig: %w", err)
	}

	return conn, nil
}
