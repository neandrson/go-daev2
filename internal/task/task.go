package task

import "time"

type Task struct {
	ID            int64         `json:"id"`
	Arg1          string        `json:"arg1"`
	Arg2          string        `json:"arg2"`
	Operation     string        `json:"operation"`
	OperationTime time.Duration `json:"operation_time"`
}
