package helper

import "regexp"

func IsDivarLink(link string) bool {
	divarRegex := regexp.MustCompile(`^https://(?:[a-zA-Z0-9-]+\.)?divar\.ir(/?|/.+)`)

	return divarRegex.MatchString(link)
}
