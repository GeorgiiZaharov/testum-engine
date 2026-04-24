package result

import "testum-engine/app/internal/service/core/answer"

type CalculationService interface {
	Calc(res answer.CheckResult) CalcResult
}

type service struct{}

func NewCalculationService() CalculationService {
	return &service{}
}

func (s *service) Calc(res answer.CheckResult) CalcResult {
	if res.Total == 0 {
		return CalcResult{
			Mark:        2,
			SuccessRate: 0,
		}
	}

	successRate := float64(res.TrueCnt) / float64(res.Total) * 100

	return CalcResult{
		Mark:        calcMark(successRate),
		SuccessRate: successRate,
	}
}

func calcMark(rate float64) int {
	switch {
	case rate == 100:
		return 5
	case rate >= 80:
		return 4
	case rate >= 50:
		return 3
	default:
		return 2
	}
}
