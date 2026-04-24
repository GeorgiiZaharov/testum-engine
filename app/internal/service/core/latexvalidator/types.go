package latexvalidator

type ValidationError struct {
	Line  int
	Error string
}
