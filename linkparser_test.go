package lp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseHeader(t *testing.T) {
	link1 := "http://foo.bar"
	rel1 := "foo"
	link2 := "http://test.baz"
	rel2 := "bar"

	header := ""
	links := ParseHeader(header)
	assert.Equal(t, 0, len(links))

	header = `<` + link1 + `>; rel="` + rel1 + `", <` + link2 + `>; rel="` + rel2 + `"`
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1])
	assert.Equal(t, link2, links[rel2])

	header = `<` + link1 + `>; rel="` + rel1 + `" <` + link2 + `>; rel="` + rel2 + `"`
	links = ParseHeader(header)
	assert.Empty(t, links[rel1])
	assert.Equal(t, link2, links[rel2])

	header = `<` + link1 + `>; rel="` + rel1 + `"; <` + link2 + `>; rel="` + rel2 + `"`
	links = ParseHeader(header)
	assert.Empty(t, links[rel1])
	assert.Equal(t, link2, links[rel2])

	header = `<` + link1 + `>; rel="` + rel1 + `", ` + link2 + `>; rel="` + rel2 + `"`
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1])
	assert.Empty(t, links[rel2])

	header = `<` + link1 + `>; rel="` + rel1 + `", <` + link2 + `; rel="` + rel2 + `"`
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1])
	assert.Empty(t, links[rel2])

	header = `<` + link1 + `>; rel="` + rel1 + `", ` + link2 + `; rel="` + rel2 + `"`
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1])
	assert.Empty(t, links[rel2])

	header = `<` + link1 + `>; rel="` + rel1 + `", <` + link2 + `>; rel=` + rel2 + ``
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1])
	assert.Equal(t, link2, links[rel2])

	header = `<` + link1 + `>; rel="` + rel1 + `", <>; rel="` + rel2 + `"`
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1])
	assert.Equal(t, "", links[rel2])

	header = `<` + link1 + `>; rel="` + rel1 + `", <` + link2 + `>; rel=""`
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1])
	assert.Empty(t, links[rel2])
	assert.Equal(t, 1, len(links))

	header = `<` + link1 + `>; rel="` + rel1 + `", ; rel="` + rel2 + `"`
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1])
	assert.Equal(t, "", links[rel2])
	assert.Equal(t, 1, len(links))

	header = `<` + link1 + `>; rel="` + rel1 + `", ; rel"` + rel2 + `"`
	links = ParseHeader(header)
	assert.Equal(t, link1, links[rel1])
	assert.Equal(t, "", links[rel2])
	assert.Equal(t, 1, len(links))
}

func TestAddLink(t *testing.T) {
	header := ""
	newH := AddLink(header, "", "foo")
	assert.Empty(t, newH)
	newH = AddLink(header, "http://example.org", "")
	assert.Empty(t, newH)

	header = ""
	newH = AddLink(header, "http://example.org", "bar")
	links := ParseHeader(newH)
	assert.Equal(t, "http://example.org", links["bar"])

	header = `<http://foo.bar>; rel="foo"`
	newH = AddLink(header, "http://example.org", "bar")
	links = ParseHeader(newH)
	assert.Equal(t, "http://foo.bar", links["foo"])
	assert.Equal(t, "http://example.org", links["bar"])
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
	header := AddLink("", "http://one.com", "foo")
	header = AddLink(header, "http://two.com", "bar")
	header = AddLink(header, "http://three.com", "baz")
	header = AddLink(header, "http://four.com", "boo")
	for i := 0; i < b.N; i++ {
		ParseHeader(header)
	}
}
