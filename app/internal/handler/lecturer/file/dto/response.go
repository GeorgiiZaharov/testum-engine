package dto

type FormatError struct {
	Error string `json:"error"`
}

type ValidationError struct {
	Line  int    `json:"line"`
	Error string `json:"error"`
}

type UploadTestResponse struct {
	FormatErrors     []FormatError     `json:"format_errors"`
	ValidationErrors []ValidationError `json:"validation_errors"`
	TestID           *int              `json:"test_id"`
	Success          bool              `json:"success"`
}

type UploadPictureResponse struct {
	Url string `json:"url"`
}
