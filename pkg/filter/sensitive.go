package filter

import "strings"

var blocklist = []string{
	"fuck", "shit", "stupid", "dick", "porn", "sex",
	"version", "health", "api", "admin", "login", "logout",
	"root", "system", "config", "shorten", "redirect",
}

func IsSensitive(shortCode string) bool {
	lower := strings.ToLower(shortCode)
	for _, w := range blocklist {
		if strings.Contains(lower, w) {
			return true
		}
	}
	return false
}
