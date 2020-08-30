package main

import (
    "net/url"
    "strings"
)

func getSub(u string) string {
    p, err := url.Parse(u)
    check(err)

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
