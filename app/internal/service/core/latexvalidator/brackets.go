package latexvalidator

import "fmt"

func checkBrackets(line string, lineNum int) []ValidationError {
	var errs []ValidationError

	type bracket struct {
		char rune
	}

	var stack []bracket

	pairs := map[rune]rune{
		')': '(',
		'}': '{',
	}

	for _, ch := range line {
		switch ch {
		case '(', '{':
			stack = append(stack, bracket{char: ch})

		case ')', '}':
			if len(stack) == 0 {
				errs = append(errs, ValidationError{
					Line:  lineNum,
					Error: fmt.Sprintf("лишняя закрывающая скобка '%c'", ch),
				})
				continue
			}

			last := stack[len(stack)-1]

			if last.char != pairs[ch] {
				errs = append(errs, ValidationError{
					Line:  lineNum,
					Error: fmt.Sprintf("неправильный порядок скобок '%c' и '%c'", last.char, ch),
				})
				stack = stack[:len(stack)-1]
				continue
			}

			stack = stack[:len(stack)-1]
		}
	}

	for _, b := range stack {
		errs = append(errs, ValidationError{
			Line:  lineNum,
			Error: fmt.Sprintf("не хватает закрывающей скобки '%c'", b.char),
		})
	}

	return errs
}
