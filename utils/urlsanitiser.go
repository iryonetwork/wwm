package utils

import (
	"strings"
)

// URLSanitize is a function sanitizing URL.
type URLSanitize func(url string) string

// WhitelistURLSanitizer returns URL Sanitizer function that replaces with '*' all the keywords in URL not matching provided whitelist.
func WhitelistURLSanitize(keyWordWhitelist []string) URLSanitize {
	whitelistMap := make(map[string]struct{})

	for _, v := range keyWordWhitelist {
		whitelistMap[v] = struct{}{}
	}

	return func(url string) string {
		s := strings.Split(url, "/")
		for k := range s {
			if _, ok := whitelistMap[s[k]]; s[k] != "" && !ok {
				s[k] = "*"
			}
		}
		return strings.Join(s, "/")
	}
}
