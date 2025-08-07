package helpers

import "strings"

func TrimSpaceToUpper(input string) string {
	return strings.TrimSpace(strings.ToUpper(input))
}
