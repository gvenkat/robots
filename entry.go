package robots

import (
	"bytes"
	"fmt"
	"strings"
)

type Entry struct {
	UserAgents []string
	Rules      []*Rule
	CrawlDelay string
}

func NewEntry() *Entry {
	return &Entry{
		UserAgents: make([]string, 0, 100),
		Rules:      make([]*Rule, 0, 100),
	}
}

func (e *Entry) String() string {
	buffer := bytes.NewBuffer([]byte{})

	buffer.WriteString("AGENTS:\n")

	for index, ua := range e.UserAgents {
		buffer.WriteString(fmt.Sprintf("\t%d: %s\n", index, ua))
	}

	buffer.WriteString("\n\nRULES:\n\n")

	for index, rule := range e.Rules {
		buffer.WriteString(fmt.Sprintf("\t%d: %s\n", index, rule))
	}

	buffer.WriteString(
		fmt.Sprintf("\n\nCRAWL DELAY: %s\n", e.CrawlDelay))

	return buffer.String()
}

func (e *Entry) AppliesTo(ua string) bool {
	agent := strings.Split(ua, "/")[0]

	for _, a := range e.UserAgents {
		catchAll := a == "*"

		if catchAll || (strings.Index(agent, a) == 0) {
			return true
		}
	}

	return false
}

func (e *Entry) CanFetch(url string) bool {
	return e.Allowance(url)
}

func (e *Entry) Allowance(url string) bool {
	for _, rule := range e.Rules {
		if rule.AppliesTo(url) {
			return rule.Allowed
		}
	}
	return true
}

func (e *Entry) SetCrawlDelay(delay string) {
	e.CrawlDelay = delay
}

func (e *Entry) AddAgent(agent string) {
	e.UserAgents = append(e.UserAgents, agent)
}

func (e *Entry) AddRule(operation, path string) *Rule {
	rule := &Rule{Path: path, Allowed: operation == "allow"}
	e.Rules = append(e.Rules, rule)
	return rule
}
