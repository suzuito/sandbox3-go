// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlcgo

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Article struct {
	ID          uuid.UUID
	Title       string
	PublishedAt pgtype.Timestamp
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
	DeletedAt   pgtype.Timestamp
}

type ArticlesSearchIndex struct {
	ArticleID   uuid.UUID
	TagIds      []uuid.UUID
	PublishedAt pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
}

type File struct {
	ID        uuid.UUID
	Name      string
	Type      int16
	MediaType pgtype.Text
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	DeletedAt pgtype.Timestamp
}

type FileThumbnail struct {
	ID        uuid.UUID
	FileID    uuid.UUID
	MediaType pgtype.Text
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	DeletedAt pgtype.Timestamp
}

type RelArticlesTag struct {
	ArticleID uuid.UUID
	TagID     uuid.UUID
}

type Tag struct {
	ID        uuid.UUID
	Name      string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
	DeletedAt pgtype.Timestamp
}
