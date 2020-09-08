package robots

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	// ErrInvalidURL given url is invalid
	ErrInvalidURL = errors.New("Invalid URL")
	// ErrHTTPRequest error making http request
	ErrHTTPRequest = errors.New("Making HTTP request failed")
)

const (
	// StateStart state representing a blank state or a state when complete block is read
	StateStart = iota
	// StateAgent state representing when an agent is read from the robots.txt file
	StateAgent
	// StateAllowOrDisallow represents state when actual rrules are being read
	StateAllowOrDisallow
)

// Line is just an alias for string type
type Line = string

// Lines is just a slice of lines
type Lines = []Line

// Robot is a struct representing the entire robots.txt file
type Robot struct {
	Lines        Lines
	Entries      []*Entry
	Sitemaps     []string
	DefaultEntry *Entry
	DisallowAll  bool
	AllowAll     bool
}

func (r *Robot) CanFetch(ua, url string) bool {
	if r.DisallowAll {
		return false
	}

	if r.AllowAll {
		return true
	}

	// FIXME: Still
	matchingEntries := r.GetMatchingEntries(ua)
	for _, x := range matchingEntries {
		return x.CanFetch(url)
	}

	return true
}

func (r *Robot) GetMatchingEntries(ua string) []*Entry {
	matchingEntries := make([]*Entry, 0, 5)

	for _, entry := range r.Entries {
		if entry.AppliesTo(ua) {
			matchingEntries = append(matchingEntries, entry)
		}
	}

	return matchingEntries
}

func (r *Robot) GetCrawlDelay(ua string) string {
	matchingEntries := r.GetMatchingEntries(ua)

	if len(matchingEntries) == 0 {
		return ""
	}

	return matchingEntries[0].CrawlDelay
}

func (r *Robot) GetAllowedPaths(ua string) []string {
	matchingEntries := r.GetMatchingEntries(ua)
	allowedPaths := make([]string, 0, 0)

	for _, entry := range matchingEntries {
		for _, rule := range entry.Rules {
			if rule.Allowed {
				allowedPaths = append(allowedPaths, rule.Path)
			}
		}
	}

	return allowedPaths
}

func (r *Robot) GetDisallowedPaths(ua string) []string {
	matchingEntries := r.GetMatchingEntries(ua)
	disallowedPaths := make([]string, 0, 0)

	for _, entry := range matchingEntries {
		for _, rule := range entry.Rules {
			if !rule.Allowed {
				disallowedPaths = append(disallowedPaths, rule.Path)
			}
		}
	}

	return disallowedPaths
}

func (r *Robot) String() string {
	buffer := bytes.NewBuffer([]byte{})

	buffer.WriteString(
		fmt.Sprintf("NUMBER OF ENTRIES: %d ", len(r.Entries)))

	buffer.WriteString(
		fmt.Sprintf("NUMBER OF SITEMAPS: %d ", len(r.Sitemaps)))

	for index, entry := range r.Entries {
		buffer.WriteString(fmt.Sprintf("ENTRY: %d \n", index))
		buffer.WriteString(entry.String())
	}

	return buffer.String()
}

func (r *Robot) AddSitemap(sitemap string) {
	r.Sitemaps = append(r.Sitemaps, sitemap)
}

func (r *Robot) AddEntry(entry *Entry) {
	r.Entries = append(r.Entries, entry)
}

func (r *Robot) Parse() {
	state := StateStart
	entry := NewEntry()

	for _, line := range r.Lines {
		// is line empty?
		if len(line) == 0 {
			switch state {
			case StateAgent:
				entry = NewEntry()
				state = StateStart

			case StateAllowOrDisallow:
				r.AddEntry(entry)
				entry = NewEntry()
				state = StateStart
			}

			continue
		}

		commentIndex := strings.Index(line, "#")
		if commentIndex > -1 {
			line = line[0:commentIndex]
		}

		line = hideColonsFromURL(line)
		parts := strings.Split(line, ":")
		if len(parts) != 2 {
			continue
		}

		parts[1] = restoreColonsIntoURL(parts[1])

		field := strings.TrimSpace(strings.ToLower(parts[0]))
		value := strings.TrimSpace(parts[1])

		switch field {
		case "user-agent":
			if state == StateAllowOrDisallow {
				r.AddEntry(entry)
				entry = NewEntry()
			}

			state = StateAgent

			if len(value) > 0 {
				entry.AddAgent(value)
			}

		case "allow", "disallow":

			if state != StateStart {
				state = StateAllowOrDisallow
				// simple string right now
				if len(value) > 0 {
					rule := entry.AddRule(field, value)

					// last rule wild card rule is applied here
					if entry.AppliesTo("*") && rule.Path == "/" {
						if rule.Allowed {
							r.AllowAll = true
							r.DisallowAll = false
						} else {
							r.AllowAll = false
							r.DisallowAll = true
						}
					}

				}
			}

		case "sitemap":
			r.AddSitemap(value)

		case "crawl-delay":
			if state != StateStart {
				state = StateAllowOrDisallow
				entry.SetCrawlDelay(value)
			}
		}
	}

	if state == StateAllowOrDisallow {
		r.AddEntry(entry)
	}
}

func NewRobot(lines Lines) *Robot {
	robot := &Robot{Lines: lines}
	robot.Parse()

	return robot
}

func isValidURL(requestURL string) bool {
	_, err := url.ParseRequestURI(requestURL)

	if err != nil {
		return false
	}

	return true
}

func FromFileName(filename string) (*Robot, error) {
	fh, err := os.Open(filename)

	defer fh.Close()

	if err != nil {
		return nil, err
	}

	return FromReader(fh)
}

func FromReader(rd io.Reader) (*Robot, error) {
	lines := make(Lines, 0, 200)
	scanner := bufio.NewScanner(rd)

	for scanner.Scan() {
		line := Line(strings.TrimSpace(scanner.Text()))
		lines = append(lines, line)
	}

	return NewRobot(lines), nil
}

func FromURL(requestURL string) (*Robot, error) {
	if !isValidURL(requestURL) {
		return nil, ErrInvalidURL
	}

	resp, err := http.Get(requestURL)

	defer resp.Body.Close()

	if err != nil {
		return nil, ErrHTTPRequest
	}

	return FromReader(resp.Body)
}
