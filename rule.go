package robots

import (
	"bytes"
	"net/url"
	"strings"
)

type Rule struct {
	Path    string
	Allowed bool
}

var endChar string = url.PathEscape("$")

// AppliesTo takes on string argument and returns a bool, true if path mactches the rule
func (r *Rule) AppliesTo(inPath string) bool {
	unescapedURL := decodePath(inPath)
	escapedURL, _ := encodePath(unescapedURL)
	endIndex := strings.Index(r.Path, endChar)
	hasStrictEnd := (endIndex > 0) && endIndex == (len(r.Path)-1)

	if (r.Path == "*") || (strings.Index(r.Path, escapedURL) == 0) {
		return true
	}

	if hasStrictEnd && r.Path == escapedURL+endChar {
		return true
	}

	if !hasStrictEnd && strings.Index(escapedURL, r.Path) == 0 {
		return true
	}

	// no globs defined
	if strings.Index(r.Path, "*") < 0 {
		return false
	}

	parts := strings.Split(r.Path, "*")

	// Some more shit to fix here
	lastMatchedIndex := 0
	lastMatchedString := escapedURL + endChar
	for _, part := range parts {
		matchedIndex := strings.Index(lastMatchedString, part)

		// no match, so lets go away
		if matchedIndex < 0 {
			return false
		}

		// matched before the last match
		if matchedIndex < lastMatchedIndex {
			return false
		}

		// update the last matched index
		lastMatchedIndex = matchedIndex
		lastMatchedString = substringFrom(lastMatchedString, lastMatchedIndex+len(part))

		// stuff in a forward slash, helps with * separated matches
		if len(lastMatchedString) > 0 && lastMatchedString[0] != '/' {
			lastMatchedString = "/" + lastMatchedString
		}
	}

	if hasStrictEnd && lastMatchedString != "" {
		return false
	}

	return true
}

func (r *Rule) String() string {
	buffer := bytes.NewBuffer([]byte{})

	if r.Allowed {
		buffer.WriteString("Allow:")
	} else {
		buffer.WriteString("Disallow:")
	}

	buffer.WriteString(r.Path)

	return buffer.String()
}
