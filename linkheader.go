package lh

import (
	"regexp"
	"strings"
)

var (
	linkexp  = regexp.MustCompile(`<[^>]*>\s*(\s*;\s*[^\(\)<>@,;:"\/\[\]\?={} \t]+=(([^\(\)<>@,;:"\/\[\]\?={} \t]+)|("[^"]*")))*(,|$)`)
	paramexp = regexp.MustCompile(`[^\(\)<>@,;:"\/\[\]\?={} \t]+=(([^\(\)<>@,;:"\/\[\]\?={} \t]+)|("[^"]*"))`)
)

func unquote(s string) string {
	if len(s) < 2 {
		return s
	}
	if s[0] != '"' {
		return s
	}
	if s[len(s)-1] != '"' {
		return s
	}
	return s[1 : len(s)-1]
}

func brack(s string) string {
	if len(s) > 0 && s[0] == '<' {
		return s
	}
	if len(s) > 0 && s[len(s)-1] == '>' {
		return s
	}
	return "<" + s + ">"
}

func debrack(s string) string {
	if s[0] != '<' {
		return s
	}
	if len(s) == 1 {
		return s[1:]
	}
	s = s[1:]

	if s[len(s)-1] != '>' {
		return s
	}
	return s[:len(s)-1]
}

// ParseHeader takes a the value of a Link header (e.g. usually using req.Header.Get("Link")) and returns a map where link values are keyed based on `rel` values.
func ParseHeader(header string) map[string]map[string]string {
	links := map[string]map[string]string{}
	if len(header) == 0 {
		return links
	}
	matches := linkexp.FindAllString(header, -1)
	for i := range matches {
		match := matches[i]
		parts := strings.Split(match, ">")
		href := ""
		if len(parts) >= 2 {
			href = debrack(parts[0])
			params := paramexp.FindAllString(parts[1], -1)
			// log.Printf("Params: %v\n", params)
			values := map[string]string{}
			for p := range params {
				param := strings.TrimSpace(params[p])
				extras := strings.Split(param, "=")
				if len(extras) > 1 {
					typ := extras[0]
					val := unquote(extras[1])
					if len(val) > 0 {
						values[typ] = val
					}
				}
			}
			// reorder and bind to the rel value
			if len(values) > 0 {
				if len(values["rel"]) > 0 {
					rel := values["rel"]
					links[rel] = map[string]string{}
					links[rel]["href"] = href
					for k, v := range values {
						links[rel][k] = v
					}
				}
			}
		}
	}
	return links
}

// AddLink returns a Link header with multiple values by adding new links to an existing Link header (which can also be empty). It only supports two parameters, the <link> and the rel="" value.
func AddLink(oldHeader, link string, params map[string]string) string {
	if len(link) == 0 || len(params) == 0 {
		return oldHeader
	}
	// find rel value
	rel := params["rel"]
	if len(rel) == 0 {
		return oldHeader
	}
	links := ParseHeader(oldHeader)
	links[rel] = params
	links[rel]["href"] = link
	newHeader := ""
	for lk, lv := range links {
		if len(newHeader) > 0 {
			newHeader += ", "
		}
		newHeader += `<` + links[lk]["href"] + `>`
		delete(lv, "href")
		for pk, pv := range lv {
			newHeader += `; ` + pk + `="` + pv + `"`
		}
	}
	return newHeader
}
