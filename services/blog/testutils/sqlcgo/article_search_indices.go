package sqlcgo

import "context"

func (q *Queries) UpsertAllArticleSearchIndices(ctx context.Context) error {
	const limit = 100
	offset := int64(0)
	for i := 0; ; i++ {
		ids, err := q.UpsertArticleSearchIndices(ctx, UpsertArticleSearchIndicesParams{
			Limit:  limit,
			Offset: int32(offset),
		})
		if err != nil {
			return err
		}

		if int32(len(ids)) < limit {
			return nil
		}

		offset += limit * int64(i)
	}
}
