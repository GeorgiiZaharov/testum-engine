package result

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"testum-engine/app/internal/service/core/answer"
)

func newService() CalculationService {
	return NewCalculationService()
}

// =========================
// EDGE CASE
// =========================

func TestCalc_TotalZero(t *testing.T) {
	s := newService()

	input := answer.CheckResult{
		TrueCnt: 0,
		Total:   0,
	}

	res := s.Calc(input)

	assert.Equal(t, 2, res.Mark)
	assert.Equal(t, 0.0, res.SuccessRate)
}

// =========================
// 100% -> EXCELLENT
// =========================

func TestCalc_100Percent(t *testing.T) {
	s := newService()

	input := answer.CheckResult{
		TrueCnt: 10,
		Total:   10,
	}

	res := s.Calc(input)

	assert.Equal(t, 5, res.Mark)
	assert.Equal(t, 100.0, res.SuccessRate)
}

// =========================
// 80% boundary -> GOOD
// =========================

func TestCalc_80Percent(t *testing.T) {
	s := newService()

	input := answer.CheckResult{
		TrueCnt: 8,
		Total:   10,
	}

	res := s.Calc(input)

	assert.Equal(t, 4, res.Mark)
	assert.Equal(t, 80.0, res.SuccessRate)
}

// =========================
// 79% -> SATISFACTORY
// =========================

func TestCalc_79Percent(t *testing.T) {
	s := newService()

	input := answer.CheckResult{
		TrueCnt: 79,
		Total:   100,
	}

	res := s.Calc(input)

	assert.Equal(t, 3, res.Mark)
	assert.Equal(t, 79.0, res.SuccessRate)
}

// =========================
// 50% boundary -> SATISFACTORY
// =========================

func TestCalc_50Percent(t *testing.T) {
	s := newService()

	input := answer.CheckResult{
		TrueCnt: 5,
		Total:   10,
	}

	res := s.Calc(input)

	assert.Equal(t, 3, res.Mark)
	assert.Equal(t, 50.0, res.SuccessRate)
}

// =========================
// BELOW 50% -> BAD
// =========================

func TestCalc_Below50Percent(t *testing.T) {
	s := newService()

	input := answer.CheckResult{
		TrueCnt: 4,
		Total:   10,
	}

	res := s.Calc(input)

	assert.Equal(t, 2, res.Mark)
	assert.Equal(t, 40.0, res.SuccessRate)
}

// =========================
// COMBINED SCENARIOS
// =========================

func TestCalc_MixedScenario(t *testing.T) {
	s := newService()

	input := answer.CheckResult{
		TrueCnt: 7,
		Total:   10,
	}

	res := s.Calc(input)

	assert.Equal(t, 3, res.Mark)
	assert.InDelta(t, 70.0, res.SuccessRate, 0.0001)
}

func TestCalc_ZeroTrueAnswers(t *testing.T) {
	s := newService()

	input := answer.CheckResult{
		TrueCnt: 0,
		Total:   10,
	}

	res := s.Calc(input)

	assert.Equal(t, 2, res.Mark)
	assert.Equal(t, 0.0, res.SuccessRate)
}
