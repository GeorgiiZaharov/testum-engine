package dto

import (
	getbasetasks "testum-engine/app/internal/service/use_case/student/get_base_tasks"
	gethardtasks "testum-engine/app/internal/service/use_case/student/get_hard_tasks"
	postbaseanswers "testum-engine/app/internal/service/use_case/student/post_base_answers"
	posthardanswers "testum-engine/app/internal/service/use_case/student/post_hard_answers"
)

func ToPostHardAnswersModel(in PostAnswersRequest) []posthardanswers.TaskAnswer {
	out := make([]posthardanswers.TaskAnswer, 0, len(in.Answers))

	for _, a := range in.Answers {
		out = append(out, posthardanswers.TaskAnswer{
			TaskID:  a.TaskID,
			Options: a.Options,
		})
	}

	return out
}

func ToPostBaseAnswersModel(in PostAnswersRequest) []postbaseanswers.TaskAnswer {
	out := make([]postbaseanswers.TaskAnswer, 0, len(in.Answers))

	for _, a := range in.Answers {
		out = append(out, postbaseanswers.TaskAnswer{
			TaskID:  a.TaskID,
			Options: a.Options,
		})
	}

	return out
}

func ToGetHardTasksResponse(
	res gethardtasks.GetHardTasksResponse,
) GetTasksResponse {
	return GetTasksResponse{
		Tasks: convertHardTasks(res.HardTasks),
	}
}

func ToGetBaseTasksResponse(
	res getbasetasks.GetBaseTasksResponse,
) GetTasksResponse {
	return GetTasksResponse{
		Tasks: convertBaseTasks(res.BaseTasks),
	}
}

func ToPostHardAnswersResponse(
	res posthardanswers.PostHardAnswersResponse,
) PostHardAnswersResponse {
	return PostHardAnswersResponse{
		IsAllCorrect: res.IsAllCorrect,
	}
}

func ToPostBaseAnswersResponse(
	res postbaseanswers.PostBaseAnswersResponse,
) PostBaseAnswersResponse {
	return PostBaseAnswersResponse{
		Success: res.Success,
	}
}

func convertHardTasks(
	in []gethardtasks.Task,
) []TaskResponse {
	out := make([]TaskResponse, 0, len(in))

	for _, t := range in {
		out = append(out, TaskResponse{
			Text:     t.Text,
			ImageURL: t.ImageURL,
			IsHard:   t.IsHard,
			Answers:  convertHardAnswers(t.Answers),
		})
	}

	return out
}

func convertBaseTasks(
	in []getbasetasks.Task,
) []TaskResponse {
	out := make([]TaskResponse, 0, len(in))

	for _, t := range in {
		out = append(out, TaskResponse{
			Text:     t.Text,
			ImageURL: t.ImageURL,
			IsHard:   t.IsHard,
			Answers:  convertBaseAnswers(t.Answers),
		})
	}

	return out
}

func convertHardAnswers(
	in []gethardtasks.Answer,
) []AnswerResponse {
	out := make([]AnswerResponse, 0, len(in))

	for _, a := range in {
		out = append(out, AnswerResponse{
			Text:     a.Text,
			ImageURL: a.ImageURL,
		})
	}

	return out
}

func convertBaseAnswers(
	in []getbasetasks.Answer,
) []AnswerResponse {
	out := make([]AnswerResponse, 0, len(in))

	for _, a := range in {
		out = append(out, AnswerResponse{
			Text:     a.Text,
			ImageURL: a.ImageURL,
		})
	}

	return out
}
