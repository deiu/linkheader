package lh

import (
	"github.com/stretchr/testify/assert"

	// "log"
	"testing"
)

func TestParseHeader(t *testing.T) {
	link1 := "http://foo.bar"
	rel1 := "foo"
	title1 := "baz"
	link2 := "http://test.baz"
	rel2 := "bar"

	header := ""
	links := ParseHeader(header)
	assert.Equal(t, 0, len(links))

	header = `<` + link1 + `>; rel="` + rel1 + `"; title="` + title1 + `", <` + link2 + `>; rel="` + rel2 + `"`
	links = ParseHeader(header)
	// log.Printf("%v\n", links)
	assert.Equal(t, link1, links[rel1]["href"])
	assert.Equal(t, title1, links[rel1]["title"])
	assert.Equal(t, link2, links[rel2]["href"])

	header = `<` + link1 + `>; rel="` + rel1 + `" <` + link2 + `>; rel="` + rel2 + `"`
	links = ParseHeader(header)
	assert.Empty(t, links[rel1])
	assert.Equal(t, link2, links[rel2]["href"])

	header = `<` + link1 + `>; rel="` + rel1 + `"; <` + link2 + `>; rel="` + rel2 + `"`
	links = ParseHeader(header)
	assert.Empty(t, links[rel1])
	assert.Equal(t, link2, links[rel2]["href"])

	header = `<` + link1 + `>; rel="` + rel1 + `", ` + link2 + `>; rel="` + rel2 + `"`
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1]["href"])
	assert.Empty(t, links[rel2])

	header = `<` + link1 + `>; rel="` + rel1 + `", <` + link2 + `; rel="` + rel2 + `"`
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1]["href"])
	assert.Empty(t, links[rel2])

	header = `<` + link1 + `>; rel="` + rel1 + `", ` + link2 + `; rel="` + rel2 + `"`
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1]["href"])
	assert.Empty(t, links[rel2])

	header = `<` + link1 + `>; rel="` + rel1 + `", <` + link2 + `>; rel=` + rel2 + ``
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1]["href"])
	assert.Equal(t, link2, links[rel2]["href"])

	header = `<` + link1 + `>; rel="` + rel1 + `", <>; rel="` + rel2 + `"`
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1]["href"])
	assert.Equal(t, "", links[rel2]["href"])

	header = `<` + link1 + `>; rel="` + rel1 + `", <` + link2 + `>; rel=""`
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1]["href"])
	assert.Empty(t, links[rel2])
	assert.Equal(t, 1, len(links))

	header = `<` + link1 + `>; rel="` + rel1 + `", ; rel="` + rel2 + `"`
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1]["href"])
	assert.Equal(t, "", links[rel2]["href"])
	assert.Equal(t, 1, len(links))

	header = `<` + link1 + `>; rel="` + rel1 + `", ; rel"` + rel2 + `"`
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1]["href"])
	assert.Equal(t, "", links[rel2]["href"])
	assert.Equal(t, 1, len(links))
}

func TestAddLink(t *testing.T) {
	header := ""
	params := map[string]string{"": "foo"}
	newH := AddLink(header, "", params)
	assert.Empty(t, newH)
	newH = AddLink(header, "http://example.org", map[string]string{})
	assert.Empty(t, newH)
	newH = AddLink(header, "http://example.org", map[string]string{"foo": "bar"})
	assert.Empty(t, newH)

	header = ""
	params = map[string]string{"rel": "bar", "title": "baz"}
	newH = AddLink(header, "http://example.org", params)
	links := ParseHeader(newH)
	assert.Equal(t, "http://example.org", links["bar"]["href"])
	assert.Equal(t, "baz", links["bar"]["title"])

	header = `<http://foo.bar>; rel="foo"; title="baz"`
	params = map[string]string{"rel": "bar"}
	newH = AddLink(header, "http://example.org", params)
	links = ParseHeader(newH)
	assert.Equal(t, "http://foo.bar", links["foo"]["href"])
	assert.Equal(t, "baz", links["foo"]["title"])
	assert.Equal(t, "http://example.org", links["bar"]["href"])
}

func TestUnquote(t *testing.T) {
	s := `"`
	assert.Equal(t, s, unquote(s))
	s = `foo"`
	assert.Equal(t, s, unquote(s))
	s = `"foo`
	assert.Equal(t, s, unquote(s))
	s = `foo`
	assert.Equal(t, s, unquote(s))
	s = `"foo"`
	assert.Equal(t, "foo", unquote(s))
}

func TestRDFBrack(t *testing.T) {
	assert.Equal(t, "<test>", brack("test"))
	assert.Equal(t, "<test", brack("<test"))
	assert.Equal(t, "test>", brack("test>"))
}

func TestRDFDebrack(t *testing.T) {
	assert.Equal(t, "a", debrack("a"))
	assert.Equal(t, "test", debrack("<test>"))
	assert.Equal(t, "test", debrack("<test"))
	assert.Equal(t, "test>", debrack("test>"))
}

func BenchmarkParse(b *testing.B) {
	params := map[string]string{"rel": "bar"}
	header := AddLink("", "http://one.com", params)
	header = AddLink(header, "http://two.com", params)
	header = AddLink(header, "http://three.com", params)
	header = AddLink(header, "http://four.com", params)
	for i := 0; i < b.N; i++ {
		ParseHeader(header)
	}
}
