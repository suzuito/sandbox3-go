package admin

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"

	"github.com/suzuito/sandbox2-common-go/libs/terrors"
)

type PasswordAsPlainText string

func (t PasswordAsPlainText) CombineSalt(s Salt) string {
	return fmt.Sprintf("%s.%s", t, s)
}

func (t PasswordAsPlainText) String() string {
	return "****"
}

type PasswordAsHash string

type HashFunc func(p PasswordAsPlainText, s Salt) (PasswordAsHash, error)

func HashFuncV1(p PasswordAsPlainText, s Salt) (PasswordAsHash, error) {
	h := md5.New()
	if _, err := io.WriteString(h, p.CombineSalt(s)); err != nil {
		return "", terrors.Errorf("failed to create hash of password")
	}
	sum := h.Sum(nil)
	return PasswordAsHash(fmt.Sprintf("%x", sum)), nil
}

type PasswordRepository interface {
	Get(ctx context.Context) (PasswordAsHash, error)
}
