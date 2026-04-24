package dto

import (
	getgroups "testum-engine/app/internal/service/use_case/lecturer/get_groups"
	gettest "testum-engine/app/internal/service/use_case/lecturer/get_test"
	gettestresult "testum-engine/app/internal/service/use_case/lecturer/get_test_result"
	gettests "testum-engine/app/internal/service/use_case/lecturer/get_tests"
)

//
// GET TESTS
//

func ToGetTests(res gettests.GetTestsResponse) GetTestsResponse {
	tests := make([]TestInfo, 0, len(res.Tests))

	for _, t := range res.Tests {
		tests = append(tests, TestInfo{
			ID:               t.ID,
			Name:             t.Name,
			CntQuestions:     t.CntQuestions,
			CntHardQuestions: t.CntHardQuestions,
			Groups:           t.Groups,
			DateCreated:      t.DateCreated,
		})
	}

	return GetTestsResponse{Tests: tests}
}

//
// GET TEST
//

func ToGetTest(res gettest.GetTestResponse) GetTestResponse {
	hard := make([]Task, 0, len(res.HardTasks))
	base := make([]Task, 0, len(res.BaseTasks))
	groups := make([]Group, 0, len(res.Groups))

	for _, g := range res.Groups {
		groups = append(groups, Group{
			GroupName:    g.GroupName,
			MembersCount: g.MembersCount,
		})
	}

	for _, t := range res.HardTasks {
		hard = append(hard, Task{
			Text:    t.Text,
			Image:   t.ImageURL,
			IsHard:  t.IsHard,
			Answers: toAnswers(t.Answers),
		})
	}

	for _, t := range res.BaseTasks {
		base = append(base, Task{
			Text:    t.Text,
			Image:   t.ImageURL,
			IsHard:  t.IsHard,
			Answers: toAnswers(t.Answers),
		})
	}

	return GetTestResponse{
		ID:               res.ID,
		Name:             res.Name,
		CntQuestions:     res.CntQuestions,
		CntHardQuestions: res.CntHardQuestions,
		Groups:           groups,
		HardTasks:        hard,
		BaseTasks:        base,
	}
}

func toAnswers(in []gettest.Answer) []Answer {
	out := make([]Answer, 0, len(in))

	for _, a := range in {
		out = append(out, Answer{
			Text:      a.Text,
			Image:     a.ImageURL,
			IsCorrect: a.IsCorrect,
		})
	}

	return out
}

//
// GET GROUPS
//

func ToGetGroups(res getgroups.GetGroupsResponse) GetGroupsResponse {
	groups := make([]Group, 0, len(res.Groups))

	for _, g := range res.Groups {
		groups = append(groups, Group{
			GroupName:    g.GroupName,
			MembersCount: g.MembersCount,
		})
	}

	return GetGroupsResponse{Groups: groups}
}

//
// GET RESULT
//

func ToGetTestResult(res gettestresult.GetTestResultResponse) GetTestResultResponse {
	results := make([]StudentResult, 0, len(res.Results))

	for _, r := range res.Results {
		results = append(results, StudentResult{
			StudentID: r.UserID,
			Name:      r.Name,
			Email:     r.Mail,
			Score:     r.SuccessRate,
			Mark:      r.Mark,
		})
	}

	return GetTestResultResponse{Results: results}
}
