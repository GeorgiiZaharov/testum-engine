package latexvalidator

import (
	"fmt"
	"regexp"
)

var commandRegex = regexp.MustCompile(`\\([a-zA-Z_]+)`)

func checkCommands(line string, lineNum int, dict *Dictionary) []ValidationError {
	var errs []ValidationError

	matches := commandRegex.FindAllStringSubmatch(line, -1)

	for _, m := range matches {
		cmd := m[1]

		if dict.IsValid(cmd) {
			continue
		}

		suggestion := dict.Suggest(cmd)

		if suggestion != "" && suggestion != cmd {
			errs = append(errs, ValidationError{
				Line: lineNum,
				Error: fmt.Sprintf(
					"возможно вы имели в виду %s вместо %s",
					suggestion,
					cmd,
				),
			})
		}
	}

	return errs
}
