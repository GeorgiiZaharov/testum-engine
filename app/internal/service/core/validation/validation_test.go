package validation

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// =========================
// HELPERS
// =========================
func newParser() *Parser {
	return NewParser(zap.NewNop())
}

func hasError(errs []FormatError, substr string) bool {
	for _, e := range errs {
		if strings.Contains(e.Error, substr) {
			return true
		}
	}
	return false
}

// =========================
// TITLE
// =========================
func Test_Title_Simple(t *testing.T) {
	lines := []string{
		"My Test",
		"",
		"2",
		"# Q",
		"+ A",
	}

	p := newParser()
	test, errs, err := p.Validate(lines)

	assert.NoError(t, err)
	assert.Nil(t, errs)

	assert.Equal(t, "My Test", test.Name)
	assert.Equal(t, 2, test.HardCount)
	assert.Len(t, test.Tasks, 1)
}

func Test_Title_MultiLine(t *testing.T) {
	lines := []string{
		"My Test",
		"Name",
		"2",
		"# Q",
		"+ A",
	}

	p := newParser()
	test, errs, err := p.Validate(lines)

	assert.NoError(t, err)
	assert.Nil(t, errs)

	assert.Equal(t, "My Test Name", test.Name)
	assert.Equal(t, 2, test.HardCount)
	assert.Len(t, test.Tasks, 1)
}

func Test_Title_Without_Name_Error(t *testing.T) {
	lines := []string{
		"2",
		"# Q",
		"+ A",
	}

	p := newParser()
	_, errs, _ := p.Validate(lines)

	assert.NotEmpty(t, errs)
	assert.True(t, hasError(errs, "не указано название"))
}

// =========================
// HARD COUNT
// =========================
func Test_HardCount_Simple(t *testing.T) {
	lines := []string{
		"Test",
		"3",
		"# Q",
		"+ A",
		"# Q",
		"+ A",
	}

	p := newParser()
	test, errs, _ := p.Validate(lines)

	assert.Empty(t, errs)
	assert.Equal(t, 3, test.HardCount)
	assert.Equal(t, test.Tasks[0].Text, "Q")
	assert.Equal(t, test.Tasks[0].Answers[0].Text, "A")
}

func Test_HardCount_Invalid_Fallback(t *testing.T) {
	lines := []string{
		"Test",
		"abc",
		"# Q",
		"+ A",
	}

	p := newParser()
	test, errs, _ := p.Validate(lines)

	assert.True(t, hasError(errs, "количество сложных"))
	assert.Nil(t, test)
}

// =========================
// QUESTION FLOW
// =========================
func Test_TwoQuestionsWithoutAnswers_Error(t *testing.T) {
	lines := []string{
		"Test",
		"1",
		"# Q1",
		"# Q2",
		"+ A",
	}

	p := newParser()
	_, errs, _ := p.Validate(lines)

	assert.NotEmpty(t, errs)
	assert.True(t, hasError(errs, "два вопроса подряд"))
}

func Test_Question_With_Image(t *testing.T) {
	lines := []string{
		"Test",
		"1",
		"# Q1",
		"https://img.com/q.png",
		"+ A",
	}

	p := newParser()
	test, errs, _ := p.Validate(lines)

	assert.Nil(t, errs)
	assert.Equal(t, "https://img.com/q.png", *test.Tasks[0].ImageURL)
}

func Test_Question_Multiple_Images_Error(t *testing.T) {
	lines := []string{
		"Test",
		"1",
		"# Q1",
		"https://img.com/1.png",
		"https://img.com/2.png",
		"+ A",
	}

	p := newParser()
	_, errs, _ := p.Validate(lines)

	assert.NotEmpty(t, errs)
	assert.True(t, hasError(errs, "у вопроса уже есть изображение"))
}

// =========================
// ANSWERS
// =========================
func Test_Answers_Multiple(t *testing.T) {
	lines := []string{
		"Test",
		"1",
		"# Q",
		"+ A1",
		"- A2",
	}

	p := newParser()
	test, errs, _ := p.Validate(lines)

	assert.Nil(t, errs)
	assert.Len(t, test.Tasks[0].Answers, 2)
	assert.True(t, test.Tasks[0].Answers[0].IsCorrect)
	assert.False(t, test.Tasks[0].Answers[1].IsCorrect)
}

func Test_Answer_With_Image(t *testing.T) {
	lines := []string{
		"Test",
		"1",
		"# Q",
		"+ A1",
		"https://img.com/a.png",
	}

	p := newParser()
	test, errs, _ := p.Validate(lines)

	assert.Nil(t, errs)
	assert.Equal(t, "https://img.com/a.png", *test.Tasks[0].Answers[0].ImageURL)
}

func Test_Answer_Multiple_Images_Error(t *testing.T) {
	lines := []string{
		"Test",
		"1",
		"# Q",
		"+ A1",
		"https://img.com/a1.png",
		"https://img.com/a2.png",
	}

	p := newParser()
	_, errs, _ := p.Validate(lines)

	assert.NotEmpty(t, errs)
	assert.True(t, hasError(errs, "у ответа уже есть изображение"))
}

func Test_Answer_Continuation_Text(t *testing.T) {
	lines := []string{
		"Test",
		"1",
		"# Q",
		"+ A1",
		"continued",
	}

	p := newParser()
	test, errs, _ := p.Validate(lines)

	assert.Nil(t, errs)
	assert.Contains(t, test.Tasks[0].Answers[0].Text, "continued")
}

// =========================
// FINISH TASK
// =========================
func Test_NoAnswers_Error(t *testing.T) {
	lines := []string{
		"Test",
		"1",
		"# Q",
	}

	p := newParser()
	_, errs, _ := p.Validate(lines)

	assert.NotEmpty(t, errs)
	assert.True(t, hasError(errs, "вопрос без ответов"))
}

func Test_NoCorrectAnswer_Error(t *testing.T) {
	lines := []string{
		"Test",
		"1",
		"# Q",
		"- A1",
		"- A2",
	}

	p := newParser()
	_, errs, _ := p.Validate(lines)

	assert.NotEmpty(t, errs)
	assert.True(t, hasError(errs, "нет правильного ответа"))
}

// =========================
// FINAL VALIDATION
// =========================
func Test_NoQuestions_Error(t *testing.T) {
	lines := []string{
		"Test",
		"1",
	}

	p := newParser()
	_, errs, _ := p.Validate(lines)

	assert.NotEmpty(t, errs)
	assert.True(t, hasError(errs, "тест не содержит вопросов"))
}

func Test_NoHardQuestionCount_Error(t *testing.T) {
	lines := []string{
		"Test",
	}

	p := newParser()
	_, errs, _ := p.Validate(lines)

	assert.NotEmpty(t, errs)
	assert.True(t, hasError(errs, "количество сложных"))
}

// =========================
// EDGE CASES
// =========================
func Test_OnlyImageQuestion(t *testing.T) {
	lines := []string{
		"Test",
		"1",
		"# https://img.com/q.png",
		"+ A",
	}

	p := newParser()
	test, errs, _ := p.Validate(lines)

	assert.Nil(t, errs)
	assert.Equal(t, "https://img.com/q.png", *test.Tasks[0].ImageURL)
	assert.Empty(t, test.Tasks[0].Text)
}

func Test_AnswerWithoutPlusMinus_Ignored(t *testing.T) {
	lines := []string{
		"Test",
		"1",
		"# Q",
		"A1",
	}

	p := newParser()
	_, errs, _ := p.Validate(lines)

	assert.NotEmpty(t, errs)
}

// =========================
// HAPPY PATH
// =========================
func Test_HappyPath_Full(t *testing.T) {
	lines := []string{
		"My Test",
		"1",
		"# Question 1",
		"text",
		"+ Correct",
		"- Wrong",
		"# Question 2",
		"+ Correct2",
	}

	p := newParser()
	test, errs, err := p.Validate(lines)

	assert.NoError(t, err)
	assert.Nil(t, errs)

	assert.Len(t, test.Tasks, 2)
	assert.Equal(t, "My Test", test.Name)

	assert.True(t, test.Tasks[0].IsHard)
	assert.False(t, test.Tasks[1].IsHard)
}
