# linkheader
Golang tool to parse and set Link headers

[![Build Status](https://api.travis-ci.org/deiu/linkparser.svg?branch=master)](https://travis-ci.org/deiu/linkparser)
[![Coverage Status](https://coveralls.io/repos/github/deiu/linkparser/badge.svg?branch=master)](https://coveralls.io/github/deiu/linkparser?branch=master)


## Install

`go get -u github.com/deiu/linkparser`

## Usage

### Getting values from a Link header

The `ParseHeader()` function returns a map where link values are keyed based on `rel` values.

```golang
// the header value usually comes from an http.Request object, using
// req.Header.Get("Link")
header := `<http://example.org>; rel="foo", <http://test.com>; rel="bar"`

links = ParseHeader(header)

println(links["foo"]) // -> http://example.org
println(links["bar"]) // -> http://test.com
```

### Setting values to a Link header

Golang's normal behavior when setting Link headers is a bit weird. For instance, doing something like `request.Header.Add("Link", headerValue)` several times with different values seems to append those values to the `Link` header slice. However, when reading back from that header with `request.Header.Get("Link")`, it only returns the first element of the slice. This behavior forced me to write my own function `AddLink()` to set multiple values to a Link header.

```golang
// the header value usually comes from an http.Request object, using
// req.Header.Get("Link")
oldHeader := req.Header.Get("Link")
println(oldHeader) // -> ""

newHeader := AddLink(oldHeader, "http://example.org", "foo")
req.Header.Set("Link", newHeader)

oldHeader = req.Header.Get("Link")
println(oldHeader) // -> <http://example.org>; rel="foo"

newHeader = AddLink(oldHeader, "http://test.com", "bar")
req.Header.Set("Link", newHeader)

println(req.Header.Get("Link"))
// -> <http://example.org>; rel="foo", <http://test.com>; rel="bar"
```
