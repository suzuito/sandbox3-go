package main

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
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

func assertHeader(
	t *testing.T,
	page playwright.Page,
	expectedAdminLinkExists bool,
) {
	locHeader := page.Locator(`[data-e2e-val="header"]`)
	e2ehelpers.AssertElementExists(t, locHeader)

	locLinkToAdmin := locHeader.Locator(`[data-e2e-val="link-to-admin"]`)
	if expectedAdminLinkExists {
		e2ehelpers.AssertElementExists(t, locLinkToAdmin)
	} else {
		e2ehelpers.AssertElementNotExists(t, locLinkToAdmin)
	}

	locFooter := page.Locator(`[data-e2e-val="footer"]`)
	e2ehelpers.AssertElementExists(t, locFooter)
}

func AssertElementsCount(t *testing.T, expectedElementCount int, loc playwright.Locator) []playwright.Locator {
	c, err := loc.Count()
	assert.NoError(t, err)
	require.Equal(t, expectedElementCount, c)

	elems, err := loc.All()
	assert.NoError(t, err)
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
