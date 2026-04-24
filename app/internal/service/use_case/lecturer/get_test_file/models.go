package gettestfile

import (
	"os"
	"time"
)

type TestInfo struct {
	ID               int
	Name             string
	CntQuestions     int
	CntHardQuestions int
	FileName         string
	Groups           []int
	DateCreated      time.Time
}

type GetTestFileRequest struct {
	UserID int
	TestID int
}

type GetTestFileResponse struct {
	File *os.File
}
