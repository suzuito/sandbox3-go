package main

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
)

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
