package article_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/suzuito/sandbox2-common-go/libs/utils"
	"github.com/suzuito/sandbox3-go/services/blog/go/internal/domains/article"
)

func TestParseQuery(t *testing.T) {
	testCases := []struct {
		desc     string
		inputQ   url.Values
		expected article.FindConditions
	}{
		{
			desc:     "ok - no query",
			inputQ:   url.Values{},
			expected: article.FindConditions{},
		},
		{
			desc: "ok",
			inputQ: url.Values{
				"tag":   []string{"hoge"},
				"page":  []string{"111"},
				"limit": []string{"222"},
				"since": []string{"333"},
				"until": []string{"444"},
			},
			expected: article.FindConditions{
				TagName: utils.Ptr("hoge"),
				Page:    111,
				Count:   222,
				PublishedAtRange: article.FindConditionRange{
					Since: utils.Ptr(time.Unix(333, 0)),
					Until: utils.Ptr(time.Unix(444, 0)),
				},
			},
		},
		{
			desc: "ok - unparsed values are parsed as zero",
			inputQ: url.Values{
				"page":  []string{"not int"},
				"limit": []string{"not int"},
				"since": []string{"not int"},
				"until": []string{"not int"},
			},
			expected: article.FindConditions{},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			fd := article.FindConditions{}
			fd.ParseQuery(tC.inputQ)
			assert.Equal(t, tC.expected, fd)
		})
	}
}

func TestQuery(t *testing.T) {
	testCases := []struct {
		desc     string
		input    article.FindConditions
		expected url.Values
	}{
		{
			desc: "ok - default page and count are ignored in query string",
			input: article.FindConditions{
				Page:  0,
				Count: 10,
			},
			expected: url.Values{},
		},
		{
			desc: "ok",
			input: article.FindConditions{
				TagName: utils.Ptr("hoge"),
				Page:    111,
				Count:   222,
				PublishedAtRange: article.FindConditionRange{
					Since: utils.Ptr(time.Unix(333, 0)),
					Until: utils.Ptr(time.Unix(444, 0)),
				},
			},
			expected: url.Values{
				"tag":   []string{"hoge"},
				"page":  []string{"111"},
				"limit": []string{"222"},
				"since": []string{"333"},
				"until": []string{"444"},
			},
		},
		/*
			{
				desc: "ok",
				inputQ: url.Values{
					"tag":   []string{"hoge"},
					"page":  []string{"111"},
					"limit": []string{"222"},
					"since": []string{"333"},
					"until": []string{"444"},
				},
				expected: article.FindConditions{
					TagName: utils.Ptr("hoge"),
					Page:    111,
					Count:   222,
					PublishedAtRange: article.FindConditionRange{
						Since: utils.Ptr(time.Unix(333, 0)),
						Until: utils.Ptr(time.Unix(444, 0)),
					},
				},
			},
			{
				desc: "ok - unparsed values are parsed as zero",
				inputQ: url.Values{
					"page":  []string{"not int"},
					"limit": []string{"not int"},
					"since": []string{"not int"},
					"until": []string{"not int"},
				},
				expected: article.FindConditions{},
			},
		*/
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := tC.input.Query()
			assert.Equal(t, tC.expected, actual)
		})
	}
}
