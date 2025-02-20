package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/admin"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/infra/rdb/sqlc/sqlcgo"
)

func (t *impl) ReadSession(ctx context.Context, id admin.LoginSessionID) (*admin.LoginSession, error) {
	queries := sqlcgo.New(t.conn)

	r, err := queries.ReadAdminLoginSessionByID(ctx, sqlcgo.ReadAdminLoginSessionByIDParams{
		ID: id.UUID(),
		ExpiredAt: pgtype.Timestamp{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		return nil, terrors.Errorf("failed to read admin login session by id: %w", err)
	}

	return r.ToAdminLoginSession(), nil
}

func (t *impl) CreateSession(ctx context.Context, id admin.LoginSessionID) (*admin.LoginSession, error) {
	tx, err := t.conn.Begin(ctx)
	if err != nil {
		return nil, terrors.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	now := time.Now()
	createdAt := pgtype.Timestamp{
		Time:  now,
		Valid: true,
	}
	expiredAt := pgtype.Timestamp{
		Time:  now.Add(time.Hour * 24),
		Valid: true,
	}

	queries := sqlcgo.New(tx)
	if err := queries.CreateAdminLoginSession(ctx, sqlcgo.CreateAdminLoginSessionParams{
		ID: id.UUID(), CreatedAt: createdAt, ExpiredAt: expiredAt,
	}); err != nil {
		return nil, terrors.Errorf("failed to create admin login session: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, terrors.Errorf("failed to commit: %w", err)
	}

	return &admin.LoginSession{
		ID:        id,
		ExpiredAt: expiredAt.Time,
	}, nil
}
