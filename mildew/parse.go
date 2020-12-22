package mildew

import (
	"net/url"
	"regexp"
	"strings"
)

// GetRoots returns all unique root domains from Subs field
func (mw *Mildew) GetRoots() (roots []string) {
	seenRoot := make(map[string]bool)
	for _, v := range mw.Subs.Slice() {
		root := extractRoot(v)
		if _, ok := seenRoot[root]; ok {
			continue
		}
		seenRoot[root] = true
		roots = append(roots, root)
	}
	return roots
}

// regex extract dotmil domain only, case-insensitive
var dotmilRx = regexp.MustCompile(`(?i)((?:([a-z0-9]\.|[a-z0-9][a-z0-9\-]{0,61}[a-z0-9])\.)+)(mil)\.?`)

func urlToSub(u string) string {
	p, err := url.Parse(u)
	if err != nil {
		return ""
	}

	sub := cleanDomain(p.Hostname())
	return sub
}

func cleanDomain(d string) string {
	d = strings.ToLower(d)
	return d
}

func isDotmil(d string) bool {
	return strings.HasSuffix(d, "mil")
}

func extractRoot(d string) string {
	split := strings.Split(d, ".")
	split = split[len(split)-2:]
	root := strings.Join(split, ".")
	return root
}
