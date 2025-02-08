package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/playwright-community/playwright-go"
)

func WriteHTML(
	t *testing.T,
	res playwright.Response,
) {
	fname := fmt.Sprintf("e2ehtmls/blog/%s.html", t.Name())

	body, err := res.Body()
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(filepath.Dir(fname), 0750); err != nil {
		panic(err)
	}

	if err := os.WriteFile(fname, body, 0644); err != nil {
		panic(err)
	}
}
