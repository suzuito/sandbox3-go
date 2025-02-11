package main

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/require"
)

func RequireElementExists(
	t *testing.T,
	l playwright.Locator,
) {
	c, err := l.Count()
	require.NoError(t, err)

	if c <= 0 {
		require.Fail(t, "element doesn't exist")
	}
}

func RequireElementNotExists(
	t *testing.T,
	l playwright.Locator,
) {
	c, err := l.Count()
	require.NoError(t, err)

	if c > 0 {
		require.Fail(t, "element exists")
	}
}

func requireHeader(
	t *testing.T,
	page playwright.Page,
	expectedAdminLinkExists bool,
) {
	locHeader := page.Locator(`[data-e2e-val="header"]`)
	RequireElementExists(t, locHeader)

	locLinkToAdmin := locHeader.Locator(`[data-e2e-val="link-to-admin"]`)
	if expectedAdminLinkExists {
		RequireElementExists(t, locLinkToAdmin)
	} else {
		RequireElementNotExists(t, locLinkToAdmin)
	}

	locFooter := page.Locator(`[data-e2e-val="footer"]`)
	RequireElementExists(t, locFooter)
}

func RequireElementsCount(t *testing.T, expectedElementCount int, loc playwright.Locator) []playwright.Locator {
	c, err := loc.Count()
	require.NoError(t, err)
	require.Equal(t, expectedElementCount, c)

	elems, err := loc.All()
	require.NoError(t, err)
	return elems
}

func RequireElementInnerText(t *testing.T, expected string, loc playwright.Locator) {
	c, err := loc.Count()
	require.NoError(t, err)
	require.Greater(t, c, 0)

	txt, err := loc.InnerText()
	require.NoError(t, err)
	require.Equal(t, expected, txt)
}

func RequireElementHasAttribute(t *testing.T, expected string, loc playwright.Locator, attr string) {
	v, err := loc.GetAttribute(attr)
	require.NoError(t, err)
	require.Equal(t, expected, v)
}
