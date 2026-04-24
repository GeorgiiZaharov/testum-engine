package dto

import (
	gettestresult "testum-engine/app/internal/service/use_case/student/get_test_result"
	"time"
)

func ToResponse(resp gettestresult.GetTestResultResponse) GetTestResultResponse {
	out := GetTestResultResponse{
		Mark:        resp.Mark,
		SuccessRate: resp.SuccessRate,
		DateStart:   resp.DateStart.Format(time.RFC3339),
	}

	if resp.DateEnd != nil {
		t := resp.DateEnd.Format(time.RFC3339)
		out.DateEnd = &t
	}

	return out
}
