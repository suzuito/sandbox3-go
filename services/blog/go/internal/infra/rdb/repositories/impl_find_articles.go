package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/tag"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/infra/rdb/sqlc/sqlcgo"
)

func (t *impl) FindArticles(ctx context.Context, conds *article.FindConditions) (article.IDs, error) {
	queries := sqlcgo.New(t.conn)
	var tagID *tag.ID
	if conds.TagName != nil {
		id, err := queries.ReadTagIDByName(ctx, *conds.TagName)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, terrors.Errorf("failed to read tag id by name: %w", err)
			}

			id = uuid.New() // tagNameに一致するtagIDは存在していないので、新規作成する
		}

		tid := tag.ID(id)
		tagID = &tid
	}

	args := []any{
		conds.Count,
		conds.Offset(),
	}
	sqlWhereClauses := []string{
		"1 = 1",
	}

	if conds.ExcludeDraft {
		sqlWhereClauses = append(sqlWhereClauses, "published_at IS NOT NULL")
	}

	if conds.PublishedAtRange.IsUsed() {
		if conds.PublishedAtRange.Since != nil {
			args = append(args, conds.PublishedAtRange.Since)
			sqlWhereClauses = append(sqlWhereClauses, fmt.Sprintf("published_at >= $%d", len(args)))
		}
		if conds.PublishedAtRange.Until != nil {
			args = append(args, conds.PublishedAtRange.Until)
			sqlWhereClauses = append(sqlWhereClauses, fmt.Sprintf("published_at <= $%d", len(args)))
		}
	}

	if tagID != nil {
		args = append(args, tagID)
		sqlWhereClauses = append(sqlWhereClauses, fmt.Sprintf("$%d = ANY (tag_ids)", len(args)))
	}

	sql := fmt.Sprintf(`
	SELECT article_id FROM articles_search_indices
	WHERE %s
	ORDER BY published_at DESC, updated_at DESC
	LIMIT $1
	OFFSET $2`, strings.Join(sqlWhereClauses, " AND "))
	rows, err := t.conn.Query(
		ctx,
		sql,
		args...,
	)
	if err != nil {
		return nil, terrors.Errorf("failed to sql query: %w", err)
	}

	ids, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return nil, terrors.Errorf("failed to scan rows: %w", err)
	}

	return article.NewIDsFromUUIDs(ids), nil
}
