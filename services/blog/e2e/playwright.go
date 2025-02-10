package main

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/require"
)

func Count(
	t *testing.T,
	l playwright.Locator,
) int {
	c, err := l.Count()
	require.NoError(t, err)
	return c
}
