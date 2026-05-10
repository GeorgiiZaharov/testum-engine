package uploadpicture

type UploadPictureRequest struct {
	UserID   int
	File     []byte
	FileName string
}

type UploadPictureResponse struct {
	URL     string
	Success bool
}
