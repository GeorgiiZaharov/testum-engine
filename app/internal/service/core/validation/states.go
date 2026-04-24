package validation

type State int

const (
	StateStart State = iota
	StateTitle
	StateHardCount

	// ===== ВОПРОС =====
	StateQuestionStart
	StateQuestionBody

	// ===== ОТВЕТ =====
	StateAnswerStart
	StateAnswerBody
)

type States struct {
	states []State
}

func NewStates() States {
	return States{
		states: []State{StateStart},
	}
}

func (s *States) Current() State {
	return s.states[len(s.states)-1]
}

func (s *States) TransitionTo(state State) {
	s.states = append(s.states, state)
}
func (s *States) Reset() {
	s.states = []State{StateStart}
}
