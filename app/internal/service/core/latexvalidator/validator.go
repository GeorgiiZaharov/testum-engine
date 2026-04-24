package latexvalidator

type Validator struct {
	dict *Dictionary
}

func New() *Validator {
	return &Validator{
		dict: NewDictionary(),
	}
}

func (v *Validator) Validate(lines []string) []ValidationError {

	var errs []ValidationError

	for i, line := range lines {
		lineNum := i + 1

		errs = append(errs, checkBrackets(line, lineNum)...)
		errs = append(errs, checkCommands(line, lineNum, v.dict)...)
	}

	return errs
}
