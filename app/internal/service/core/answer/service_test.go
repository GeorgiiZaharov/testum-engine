package answer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newService() CheckService {
	return NewCheckService()
}

// =========================
// HAPPY PATH
// =========================
func Test_Check_AllCorrect(t *testing.T) {
	s := newService()

	student := []TaskAnswer{
		{TaskID: 1, SelectedOptions: []int{1, 2}},
		{TaskID: 2, SelectedOptions: []int{3}},
	}

	trueAns := []TaskAnswer{
		{TaskID: 1, SelectedOptions: []int{2, 1}}, // порядок другой
		{TaskID: 2, SelectedOptions: []int{3}},
	}

	res, err := s.Check(student, trueAns)

	assert.NoError(t, err)
	assert.Equal(t, 2, res.TrueCnt)
	assert.Equal(t, 2, res.Total)
}

func Test_Check_PartiallyCorrect(t *testing.T) {
	s := newService()

	student := []TaskAnswer{
		{TaskID: 1, SelectedOptions: []int{1, 2}},
		{TaskID: 2, SelectedOptions: []int{3}},
	}

	trueAns := []TaskAnswer{
		{TaskID: 1, SelectedOptions: []int{1, 2}},
		{TaskID: 2, SelectedOptions: []int{4}}, // ошибка
	}

	res, err := s.Check(student, trueAns)

	assert.NoError(t, err)
	assert.Equal(t, 1, res.TrueCnt)
	assert.Equal(t, 2, res.Total)
}

func Test_Check_NoneCorrect(t *testing.T) {
	s := newService()

	student := []TaskAnswer{
		{TaskID: 1, SelectedOptions: []int{1}},
	}

	trueAns := []TaskAnswer{
		{TaskID: 1, SelectedOptions: []int{2}},
	}

	res, err := s.Check(student, trueAns)

	assert.NoError(t, err)
	assert.Equal(t, 0, res.TrueCnt)
	assert.Equal(t, 1, res.Total)
}

// =========================
// ORDER INSENSITIVE
// =========================
func Test_Check_OrderDoesNotMatter(t *testing.T) {
	s := newService()

	student := []TaskAnswer{
		{TaskID: 1, SelectedOptions: []int{3, 1, 2}},
	}

	trueAns := []TaskAnswer{
		{TaskID: 1, SelectedOptions: []int{1, 2, 3}},
	}

	res, err := s.Check(student, trueAns)

	assert.NoError(t, err)
	assert.Equal(t, 1, res.TrueCnt)
}

// =========================
// ERRORS
// =========================
func Test_Check_TaskCountMismatch(t *testing.T) {
	s := newService()

	student := []TaskAnswer{
		{TaskID: 1, SelectedOptions: []int{1}},
	}

	trueAns := []TaskAnswer{
		{TaskID: 1, SelectedOptions: []int{1}},
		{TaskID: 2, SelectedOptions: []int{2}},
	}

	_, err := s.Check(student, trueAns)

	assert.ErrorIs(t, err, ErrTaskMismatch)
}

func Test_Check_TaskIDMismatch(t *testing.T) {
	s := newService()

	student := []TaskAnswer{
		{TaskID: 999, SelectedOptions: []int{1}},
	}

	trueAns := []TaskAnswer{
		{TaskID: 1, SelectedOptions: []int{1}},
	}

	_, err := s.Check(student, trueAns)

	assert.ErrorIs(t, err, ErrTaskMismatch)
}

// =========================
// EDGE CASES
// =========================
func Test_Check_EmptySlices(t *testing.T) {
	s := newService()

	res, err := s.Check([]TaskAnswer{}, []TaskAnswer{})

	assert.NoError(t, err)
	assert.Equal(t, 0, res.TrueCnt)
	assert.Equal(t, 0, res.Total)
}

func Test_Check_EmptyOptions(t *testing.T) {
	s := newService()

	student := []TaskAnswer{
		{TaskID: 1, SelectedOptions: []int{}},
	}

	trueAns := []TaskAnswer{
		{TaskID: 1, SelectedOptions: []int{}},
	}

	res, err := s.Check(student, trueAns)

	assert.NoError(t, err)
	assert.Equal(t, 1, res.TrueCnt)
}

