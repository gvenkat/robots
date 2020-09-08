package robots

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	// PortRegex is RE to match port in a URL
	PortRegex, _ = regexp.Compile(`:\d+`)

	// ReversePortRegex is RE to take already replaced port
	ReversePortRegex, _ = regexp.Compile(`__\d+`)
)

const FAKE_DOMAIN string = "http://www.example.com"

func hideColonsFromURL(url string) string {
	newURL := strings.ReplaceAll(url, "http:", "http_")
	newURL = strings.ReplaceAll(newURL, "https:", "http_")
	newURL = string(PortRegex.ReplaceAll([]byte(newURL), []byte("__$1")))

	return newURL
}

func restoreColonsIntoURL(url string) string {
	newURL := strings.ReplaceAll(url, "http_", "http:")
	newURL = strings.ReplaceAll(newURL, "https_", "https:")
	newURL = string(ReversePortRegex.ReplaceAll([]byte(newURL), []byte(":$1")))

	return newURL
}

func substringFrom(haystack string, fromIndex int) string {
	return haystack[fromIndex:len(haystack)]
}

func encodePath(path string) (string, error) {
	url_, err := url.Parse(path)

	if err != nil {
		return "", err
	}

	return url_.String(), nil
}

func decodePath(path string) string {
	unescapedUrl, err := url.QueryUnescape(path)

	if err != nil {
		return ""
	}

	return unescapedUrl
}
