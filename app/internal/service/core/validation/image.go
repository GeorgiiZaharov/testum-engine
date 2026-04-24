package validation

import (
	"regexp"
	"strings"
)

var imageLineRegex = regexp.MustCompile(`^(https?://[^\s]+\.(png|jpg|jpeg|gif|webp))$`)

// возвращает ссылку или nil
func extractImage(line string) *string {
	line = strings.TrimSpace(line)

	if imageLineRegex.MatchString(line) {
		return &line
	}

	return nil
}
