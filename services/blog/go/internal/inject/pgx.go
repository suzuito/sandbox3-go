package inject

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
)

func NewPgxConn(ctx context.Context, env *Environment) (*pgx.Conn, error) {
	runtimeParams := map[string]string{}
	if env.Env == "loc" {
		runtimeParams["sslmode"] = "disable"
	}

	conf := pgx.ConnConfig{
		Config: pgconn.Config{
			Host:          env.DBHost,
			Port:          env.DBPort,
			User:          env.DBUser,
			Password:      env.DBPassword,
			Database:      env.DBName,
			RuntimeParams: runtimeParams,
		},
	}

	conn, err := pgx.ConnectConfig(ctx, &conf)
	if err != nil {
		return nil, terrors.Errorf("failed to pgx.ConnectConfig: %w", err)
	}

	return conn, nil
}
