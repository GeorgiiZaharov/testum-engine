package answer

import (
	"reflect"
	"sort"
)

type CheckService interface {
	Check(studentAnswers []TaskAnswer, trueAnswers []TaskAnswer) (CheckResult, error)
}

type service struct{}

func NewCheckService() CheckService {
	return &service{}
}

func (s *service) Check(
	studentAnswers []TaskAnswer,
	trueAnswers []TaskAnswer,
) (CheckResult, error) {

	if len(studentAnswers) != len(trueAnswers) {
		return CheckResult{}, ErrTaskMismatch
	}

	trueMap := make(map[int][]int, len(trueAnswers))

	for _, t := range trueAnswers {
		trueMap[t.TaskID] = normalize(t.SelectedOptions)
	}

	trueCnt := 0

	for _, student := range studentAnswers {
		trueOpts, ok := trueMap[student.TaskID]
		if !ok {
			return CheckResult{}, ErrTaskMismatch
		}

		studentOpts := normalize(student.SelectedOptions)

		if equal(studentOpts, trueOpts) {
			trueCnt++
		}
	}

	return CheckResult{
		TrueCnt: trueCnt,
		Total:   len(trueAnswers),
	}, nil
}

func normalize(opts []int) []int {
	cp := make([]int, len(opts))
	copy(cp, opts)

	sort.Ints(cp)
	return cp
}

func equal(a, b []int) bool {
	return reflect.DeepEqual(a, b)
}
