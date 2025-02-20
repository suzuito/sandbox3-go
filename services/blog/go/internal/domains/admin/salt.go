package admin

import "context"

type Salt string

type SaltRepository interface {
	Get(ctx context.Context) (Salt, error)
}
