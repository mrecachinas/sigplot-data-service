//go:build !ui

package ui

import (
	"net/http"
)

var stubHTML = []byte(`<!DOCTYPE html>
<html>
<body>
<p>SDS UI is not available in this build. Build with -tags ui to include the UI.</p>
</body>
</html>
`)

// GetFileSystem returns an empty filesystem when UI is not built in
func GetFileSystem() http.FileSystem {
	return http.Dir(".")
}

// LoadUI returns a stub HTML page when UI is not built in
func LoadUI() ([]byte, error) {
	return stubHTML, nil
}
