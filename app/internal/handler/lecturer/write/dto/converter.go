package dto

import (
	deleteattempt "testum-engine/app/internal/service/use_case/lecturer/delete_attempt"
	deletetest "testum-engine/app/internal/service/use_case/lecturer/delete_test"
	giveaccess "testum-engine/app/internal/service/use_case/lecturer/give_access"
	takeaccess "testum-engine/app/internal/service/use_case/lecturer/take_access"
)

func ToDeleteTestResponse(res deletetest.DeleteTestResponse) DeleteTestResponse {
	return DeleteTestResponse{
		Success: res.Success,
	}
}

func ToAccessResponseGive(res giveaccess.GiveAccessResponse) AccessResponse {
	return AccessResponse{
		Success: res.Success,
	}
}

func ToAccessResponseTake(res takeaccess.TakeAccessResponse) AccessResponse {
	return AccessResponse{
		Success: res.Success,
	}
}

func FromUseCaseResponse(resp deleteattempt.DeleteAttemptResponse) DeleteAttemptResponse {
	return DeleteAttemptResponse{
		Success: resp.Success,
	}
}
