package dto

import uploadtest "testum-engine/app/internal/service/use_case/lecturer/upload_test"

func ToUploadResponse(res uploadtest.UploadTestResponse) UploadTestResponse {
	format := make([]FormatError, len(res.FormatErrors))
	for i, e := range res.FormatErrors {
		format[i] = FormatError{Error: e.Error}
	}

	validation := make([]ValidationError, len(res.ValidationErrors))
	for i, e := range res.ValidationErrors {
		validation[i] = ValidationError{
			Line:  e.Line,
			Error: e.Error,
		}
	}

	return UploadTestResponse{
		FormatErrors:     format,
		ValidationErrors: validation,
		TestID:           res.TestID,
		Success:          res.Success,
	}
}
