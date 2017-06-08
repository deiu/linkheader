# linkparser
Golang tool to parse and set Link headers

[![Build Status](https://api.travis-ci.org/deiu/linkparser.svg?branch=master)](https://travis-ci.org/deiu/linkparser)
[![Coverage Status](https://coveralls.io/repos/github/deiu/linkparser/badge.svg?branch=master)](https://coveralls.io/github/deiu/linkparser?branch=master)


## Install

`go get -u github.com/deiu/linkparser`

## Usage

### Getting values from a Link header

```golang
// the header value usually comes from an http.Request object, using
// req.Header.Get("Link")
header := `<http://example.org>; rel="foo", <http://test.com>; rel="bar"`

links = ParseHeader(header)

println(links["foo"]) // -> http://example.org
println(links["bar"]) // -> http://test.com
```

### Setting values to a Link header

Golang's normal behavior when setting Link headers is a bit weird. For instance, doing something like `request.Header.Add("Link", headerValue)` several times with different values seems to append those values to the `Link` header slice. However, when reading back from that header with `request.Header.Get("Link")`, it only returns the first element of the slice, which made me write my own function to set multiple values to a Link header.

```golang
// the header value usually comes from an http.Request object, using
// req.Header.Get("Link")
header := req.Header.Get("Link")
println(header) // -> ""

req.Header.Set("Link", AddLink(header, "http://example.org", "foo"))
println(req.Header.Get("Link")) // -> <http://example.org>; rel="foo"

header = req.Header.Get("Link")
req.Header.Set("Link", AddLink(header, "http://test.com", "bar"))
println(req.Header.Get("Link"))
// -> <http://example.org>; rel="foo", <http://test.com>; rel="bar"
```
