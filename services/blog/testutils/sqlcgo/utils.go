package sqlcgo

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func NewPgTypeFromTime(t time.Time) pgtype.Timestamp {
	return NewPgTypeFromTimePtr(&t)
}

func NewPgTypeFromTimePtr(t *time.Time) pgtype.Timestamp {
	tm := pgtype.Timestamp{}

	if t == nil {
		tm.Valid = false
		return tm
	}

	tm.Valid = true
	tm.Time = *t

	return tm
}
