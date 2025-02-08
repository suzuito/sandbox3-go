package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/suzuito/sandbox2-common-go/libs/terrors"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
)

func (t *impl) FindArticles(ctx context.Context, conds *article.FindConditions) (article.IDs, error) {
	args := []any{
		conds.Count,
		conds.Offset(),
	}
	sqlWhereClauses := []string{
		"1 = 1",
		"published_at IS NOT NULL",
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

	if conds.TagID != nil {
		args = append(args, conds.TagID)
		sqlWhereClauses = append(sqlWhereClauses, fmt.Sprintf("$%d = ANY (tag_ids)", len(args)))
	}

	sql := fmt.Sprintf(`
	SELECT article_id FROM articles_search_indices
	WHERE %s
	ORDER BY published_at DESC
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
