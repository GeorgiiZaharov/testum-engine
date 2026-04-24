package uploadtest

type UploadTestRequest struct {
	UserID           int
	File             []byte
	FileName         string
	IgnoreValidation bool
}

type UploadTestResponse struct {
	FormatErrors     []FormatError
	ValidationErrors []ValidationError
	TestID           *int
	Success          bool
}

type FormatError struct {
	Error string
}

type ValidationError struct {
	Line  int
	Error string
}
