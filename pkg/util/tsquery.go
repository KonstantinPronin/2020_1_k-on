package util

import (
	"fmt"
	"regexp"
	"strings"
)

func PlainTextToQuery(text string) string {
	var query string
	text = strings.TrimSpace(text)
	space := regexp.MustCompile(`\s+`)
	text = space.ReplaceAllString(text, " ")
	words := strings.Split(text, " ")

	for _, str := range words {
		if query == "" {
			query = str
			continue
		}
		query = fmt.Sprintf("%s | %s", query, str)
	}

	return query
}
