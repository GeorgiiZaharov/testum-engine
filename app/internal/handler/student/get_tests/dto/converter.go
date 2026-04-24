package dto

import (
	"testum-engine/app/internal/service/use_case/student/get_active_test"
	"testum-engine/app/internal/service/use_case/student/get_finished_test"
)

// Конвертер для активных тестов
func ConvertActiveTestsToDTO(activeTests []getactivetest.StudentActiveTest) []StudentActiveTest {
	var result []StudentActiveTest
	for _, at := range activeTests {
		result = append(result, StudentActiveTest{
			ID:               at.ID,
			Name:             at.Name,
			LecturerName:     at.LecturerName,
			CntQuestions:     at.CntQuestions,
			CntHardQuestions: at.CntHardQuestions,
			DateStart:        at.DateStart,
		})
	}
	return result
}

// Конвертер для завершённых тестов
func ConvertFinishedTestsToDTO(finishedTests []getfinishedtest.StudentFinishTest) []StudentFinishTest {
	var result []StudentFinishTest
	for _, ft := range finishedTests {
		result = append(result, StudentFinishTest{
			ID:               ft.ID,
			Name:             ft.Name,
			LecturerName:     ft.LecturerName,
			CntQuestions:     ft.CntQuestions,
			CntHardQuestions: ft.CntHardQuestions,
			Mark:             ft.Mark,
			SuccessRate:      ft.SuccessRate,
			DateStart:        ft.DateStart,
			DateEnd:          ft.DateEnd,
		})
	}
	return result
}
