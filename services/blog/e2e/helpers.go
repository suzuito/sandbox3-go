package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/suzuito/sandbox2-common-go/libs/e2ehelpers"
)

type PlaywrightTestCaseForSSR struct {
	Desc     string
	Setup    func(t *testing.T, testID e2ehelpers.TestID, exe *PlaywrightTestCaseForSSRExec)
	Teardown func(t *testing.T, testID e2ehelpers.TestID)
}

func (c *PlaywrightTestCaseForSSR) Run(t *testing.T) {
	if c.Setup == nil {
		panic(errors.New("setup function is required"))
	}

	err := playwright.Install()
	require.NoError(t, err)

	testID := e2ehelpers.NewTestID()
	exe := PlaywrightTestCaseForSSRExec{}
	c.Setup(t, testID, &exe)

	pw, err := playwright.Run()
	require.NoError(t, err)

	browser, err := pw.Chromium.Launch()
	require.NoError(t, err)

	page, err := browser.NewPage()
	require.NoError(t, err)

	if exe.Do != nil {
		exe.Do(t, pw, browser, page)
	}
}

type PlaywrightTestCaseForSSRExec struct {
	Do func(t *testing.T, pw *playwright.Playwright, browser playwright.Browser, page playwright.Page)
}

func AssertElementExists(t *testing.T, loc playwright.Locator) {
	c, err := loc.Count()
	require.NoError(t, err)
	assert.Greaterf(t, c, 0, "element %+v does not exist", loc)
}

func AssertElementNotExists(t *testing.T, loc playwright.Locator) {
	c, err := loc.Count()
	require.NoError(t, err)
	assert.Lessf(t, c, 1, "element %+v exists", loc)
}

type HTTPServerTestCase struct {
	Desc     string
	Setup    func(t *testing.T, testID e2ehelpers.TestID, expected *HTTPServerTestCaseExpected) *http.Request
	Teardown func(t *testing.T, testID e2ehelpers.TestID)
}

func (c *HTTPServerTestCase) Run(t *testing.T) {
	if c.Setup == nil {
		panic(errors.New("setup function is required"))
	}

	testID := e2ehelpers.NewTestID()
	expected := HTTPServerTestCaseExpected{
		Header: make(http.Header),
	}
	input := c.Setup(t, testID, &expected)

	actual, err := http.DefaultClient.Do(input)
	require.NoError(t, err)

	assert.Equal(t, expected.StatusCode, actual.StatusCode)
	for expectedKey, expectedValue := range expected.Header {
		expectedKey = http.CanonicalHeaderKey(expectedKey)
		actualValue, ok := actual.Header[expectedKey]
		if ok {
			assert.Equal(t, expectedValue, actualValue)
		} else {
			assert.Fail(t, fmt.Sprintf("expected key '%s' does not exist in headers", expectedKey))
		}
	}

	if expected.Assertions != nil {
		expected.Assertions(t)
	}
}

type HTTPServerTestCaseExpected struct {
	StatusCode int
	Header     http.Header
	Assertions func(t *testing.T)
}

func MustHTTPRequest(method, url string, body io.Reader) *http.Request {
	r, err := http.NewRequest(
		method,
		url,
		nil,
	)
	if err != nil {
		panic(err)
	}

	return r
}
