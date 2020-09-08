package robots

import (
	"strings"
	"testing"
)

func TestEmptyFile(t *testing.T) {
	empty := strings.NewReader("")
	parser, err := FromReader(empty)

	if err != nil {
		t.Errorf("failed to parse and empty string")
	}

	if len(parser.Entries) > 0 {
		t.Errorf("can not have positive number of entries")
	}
}

func TestCatchAllUserAgent(t *testing.T) {
	example := `
User-agent: *
Disallow: */s?k=*&rh=n*p_*p_*p_
Disallow: /dp/product-availability/
Disallow: /gp/cart
# globs
Allow: /wishlist/universal*
Allow: /wishlist/vendor-button*
Allow: /wishlist/get-button*
# globs and termination
Disallow: /hz/help/contact/*/message/$
`

	robot, err := FromReader(strings.NewReader(example))

	if err != nil {
		t.Errorf("parsing the catch all example failed")
	}

	if len(robot.Entries) != 1 {
		t.Errorf("must exactly have one entries")
	}

	if len(robot.GetDisallowedPaths("*")) == 0 {
		t.Errorf("get disallowed path must have registered disallowed paths")
	}

	if len(robot.GetAllowedPaths("*")) == 0 {
		t.Errorf("get allowed path must have registered allowed paths")
	}

	allowedUrls := []string{
		"/some/blah/blah",
		"/wishlist/get-buttonhello/doodlbar",
		"/hz/help/contact/moodlebar/message/hello",
	}

	disallowedUrls := []string{
		"/dp/product-availability/",
		"/dp/product-availability/hello/foo/bar",
		"/gp/cart",
		"/hz/help/contact/foobar/message/",
		"/hz/help/contact/message/",
	}

	for _, url := range allowedUrls {
		if !robot.CanFetch("*", url) {
			t.Errorf("expecting %s to be allowed", url)
		}
	}

	for _, url := range disallowedUrls {
		if robot.CanFetch("*", url) {
			t.Errorf("expecting %s to be disallowed", url)
		}
	}
}

func TestAllowAndDisallowAll(t *testing.T) {
	example := `
User-agent: *
Disallow: / 
	`
	r1, _ := FromReader(strings.NewReader(example))

	if !r1.DisallowAll {
		t.Errorf("disallow all not set")
	}

	example2 := `
User-agent: *
Allow: / 
	`

	r2, _ := FromReader(strings.NewReader(example2))

	if !r2.AllowAll {
		t.Errorf("disallow all not set")
	}

}

func TestEncodeUrls(t *testing.T) {
	if decodePath("/s?*rh=n%3A1380045031") != "/s?*rh=n:1380045031" {
		t.Errorf("decode path failed")
	}
}
