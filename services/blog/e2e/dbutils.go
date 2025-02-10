package main

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/suzuito/sandbox3-go/services/blog/testutils/sqlcgo"
)

func RunTx(ctx context.Context, conn *pgx.Conn, f func(pgx.Tx) error) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) // nolint:errcheck

	if err := f(tx); err != nil {
		return err
	} else if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

type InitDBArg struct {
	Articles        sqlcgo.CreateArticlesParamsList
	Tags            sqlcgo.CreateTagsParamsList
	RelArticlesTags sqlcgo.CreateRelArticlesTagsParamsList
}

func MustSetupDB(
	ctx context.Context,
	conn *pgx.Conn,
	arg InitDBArg,
) {
	err := RunTx(ctx, conn, func(tx pgx.Tx) error {
		q := sqlcgo.New(tx)

		if len(arg.Articles) > 0 {
			if _, err := q.CreateArticles(ctx, arg.Articles); err != nil {
				return err
			}
		}
		if len(arg.Tags) > 0 {
			if _, err := q.CreateTags(ctx, arg.Tags); err != nil {
				return err
			}
		}
		if len(arg.RelArticlesTags) > 0 {
			if _, err := q.CreateRelArticlesTags(ctx, arg.RelArticlesTags); err != nil {
				return err
			}
		}

		if err := q.UpsertAllArticleSearchIndices(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}

func MustTeardownDB(
	ctx context.Context,
	conn *pgx.Conn,
) {
	err := RunTx(ctx, conn, func(tx pgx.Tx) error {
		q := sqlcgo.New(tx)

		if err := q.DeleteRelArticlesTagsPhysically(ctx); err != nil {
			return err
		}
		if err := q.DeleteTagsPhysically(ctx); err != nil {
			return err
		}
		if err := q.DeleteArticlesPhysically(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		panic(err)
	}
}
