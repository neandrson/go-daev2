package timeout

import (
	"context"
	"time"
)

type Timeout struct {
	Timer  *time.Timer
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewTimeout(duration time.Duration) *Timeout {
	ctx, cancel := context.WithCancel(context.Background())
	return &Timeout{
		Timer:  time.NewTimer(duration),
		Ctx:    ctx,
		Cancel: cancel,
	}
}
