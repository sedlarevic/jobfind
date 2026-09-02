package preprocessing

import "strings"

func Normalize(s string) string {
	s = strings.ToLower(s)

	return s
}
