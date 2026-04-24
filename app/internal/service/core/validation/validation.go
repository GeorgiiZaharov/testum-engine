package validation

import (
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type Parser struct {
	states     States
	repeatLine bool

	Test TestFromFile
	Task *Task

	Errors []FormatError

	logger *zap.Logger

	lineNum     int
	currentLine string
}

func (p *Parser) reset() {
	p.lineNum = 0
	p.currentLine = ""
	p.Errors = nil
	p.Test = TestFromFile{}
	p.states.Reset()
}
func NewParser(logger *zap.Logger) *Parser {
	return &Parser{
		states: NewStates(),
		logger: logger,
	}
}

func (p *Parser) Validate(lines []string) (*TestFromFile, []FormatError, error) {
	p.reset()
	p.states.TransitionTo(StateTitle)

	for _, rawLine := range lines {
		p.lineNum++
		p.currentLine = strings.TrimSpace(rawLine)

		if p.currentLine == "" {
			continue
		}

		p.repeatLine = true

		for p.repeatLine {
			p.repeatLine = false

			switch p.states.Current() {

			case StateTitle:
				p.handleTitle()

			case StateHardCount:
				p.handleHardCount()

			case StateQuestionStart:
				p.handleQuestionStart()

			case StateQuestionBody:
				p.handleQuestionBody()

			case StateAnswerStart:
				p.handleAnswerStart()

			case StateAnswerBody:
				p.handleAnswerBody()
			}
		}
	}

	p.finishTask()

	if p.states.Current() == StateTitle {
		p.addError("не указано количество сложных вопросов")
	}

	if p.states.Current() == StateHardCount {
		p.addError("тест не содержит вопросов")
	}

	if len(p.Errors) > 0 {
		return nil, p.Errors, nil
	}

	return &p.Test, nil, nil
}

// =========================
// TITLE + COUNT
// =========================
func (p *Parser) handleTitle() {
	if _, err := strconv.Atoi(p.currentLine); err == nil {
		if p.Test.Name == "" {
			p.addError("не указано название теста")
		}
		p.states.TransitionTo(StateHardCount)
		p.repeatLine = true
		return
	}

	if p.Test.Name == "" {
		p.Test.Name = p.currentLine
	} else {
		p.Test.Name += " " + p.currentLine
	}
}

func (p *Parser) handleHardCount() {
	count, err := strconv.Atoi(p.currentLine)
	if err != nil {
		p.states.TransitionTo(StateQuestionStart)
		p.repeatLine = true
		return
	}
	p.Test.HardCount = count
}

// =========================
// QUESTION
// =========================
func (p *Parser) handleQuestionStart() {
	p.startNewTask()
	p.states.TransitionTo(StateQuestionBody)
}

func (p *Parser) handleQuestionBody() {
	// новый вопрос
	if strings.HasPrefix(p.currentLine, "#") {
		p.addError("два вопроса подряд без ответов")
		p.finishTask()
		p.repeatLine = true
		p.states.TransitionTo(StateQuestionStart)
		return
	}

	// ответы
	if isAnswer(p.currentLine) {
		p.states.TransitionTo(StateAnswerStart)
		p.repeatLine = true
		return
	}

	// изображение
	if img := extractImage(p.currentLine); img != nil {
		if p.Task.ImageURL != nil {
			p.addError("у вопроса уже есть изображение")
			return
		}
		p.Task.ImageURL = img
		return
	}

	p.Task.Text += " " + p.currentLine
}

// =========================
// ANSWERS
// =========================
func (p *Parser) handleAnswerStart() {
	text := strings.TrimSpace(p.currentLine[1:])

	answer := Answer{
		Text:      text,
		IsCorrect: strings.HasPrefix(p.currentLine, "+"),
	}

	p.Task.Answers = append(p.Task.Answers, answer)
	p.states.TransitionTo(StateAnswerBody)
}

func (p *Parser) handleAnswerBody() {
	// новый вопрос
	if strings.HasPrefix(p.currentLine, "#") {
		p.finishTask()
		p.states.TransitionTo(StateQuestionStart)
		p.repeatLine = true
		return
	}

	// новый ответ
	if isAnswer(p.currentLine) {
		p.states.TransitionTo(StateAnswerStart)
		p.repeatLine = true
		return
	}

	// изображение
	if img := extractImage(p.currentLine); img != nil {
		last := &p.Task.Answers[len(p.Task.Answers)-1]

		if last.ImageURL != nil {
			p.addError("у ответа уже есть изображение")
			return
		}

		last.ImageURL = img
		return
	}

	last := &p.Task.Answers[len(p.Task.Answers)-1]
	last.Text += " " + p.currentLine
}

// =========================
// HELPERS
// =========================
func (p *Parser) startNewTask() {
	line := strings.TrimSpace(strings.TrimPrefix(p.currentLine, "#"))

	p.Task = &Task{
		IsHard: len(p.Test.Tasks) < p.Test.HardCount,
	}

	if img := extractImage(line); img != nil {
		p.Task.ImageURL = img
	} else {
		p.Task.Text = line
	}
}

func (p *Parser) finishTask() {
	if p.Task == nil {
		return
	}

	if len(p.Task.Answers) == 0 {
		p.addError("вопрос без ответов")
	}

	correct := 0
	for _, a := range p.Task.Answers {
		if a.IsCorrect {
			correct++
		}
	}

	if correct == 0 {
		p.addError("нет правильного ответа")
	}

	p.Test.Tasks = append(p.Test.Tasks, *p.Task)
	p.Task = nil
}

func (p *Parser) addError(msg string) {
	errMsg := fmt.Sprintf("строка %d: %s", p.lineNum, msg)

	p.Errors = append(p.Errors, FormatError{
		Error: errMsg,
	})

	p.logger.Warn("validation error",
		zap.Int("line", p.lineNum),
		zap.String("error", msg),
	)
}

func isAnswer(line string) bool {
	return strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-")
}
