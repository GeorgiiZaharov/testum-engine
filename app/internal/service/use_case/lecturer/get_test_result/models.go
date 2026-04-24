package gettestresult

type GetTestResultRequest struct {
	UserID int
	TestID int
	Group  string
	Year   int
}

type GetTestResultResponse struct {
	Results []StudentResult
}
