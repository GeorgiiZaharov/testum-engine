package dto

import (
	getlecturers "testum-engine/app/internal/service/use_case/admin/get_lecturers"
)

func ToGetLecturersResponse(
	res getlecturers.GetLecturersResponse,
) GetLecturersResponse {
	lecturers := make([]Lecturer, 0, len(res.Lecturers))

	for _, lecturer := range res.Lecturers {
		lecturers = append(lecturers, Lecturer{
			ID:           lecturer.ID,
			Login:        lecturer.Login,
			Mail:         lecturer.Mail,
			Name:         lecturer.Name,
			Group:        lecturer.Group,
			DateCreated:  lecturer.DateCreated,
			DateModified: lecturer.DateModified,
		})
	}

	return GetLecturersResponse{
		Lecturers: lecturers,
	}
}
