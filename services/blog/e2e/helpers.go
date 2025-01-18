package main

/*
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
*/
