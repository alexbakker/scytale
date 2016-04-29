package main

import "strings"

func stripChar(s, c string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(c, r) < 0 {
			return r
		}
		return -1
	}, s)
}
