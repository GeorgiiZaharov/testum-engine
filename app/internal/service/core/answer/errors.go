package answer

import "errors"

var (
	// несоответствие task_id или количества
	ErrTaskMismatch = errors.New("task ids mismatch between student answers and true answers")
)
