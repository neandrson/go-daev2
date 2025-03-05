package config

import (
	"fmt"
	"os"
	"time"
)

const (
	errMsgFmt = "The %s environment variable is not set or has an incorrect value."
)

type Config struct {
	Addtime time.Duration
	Subtime time.Duration
	Multime time.Duration
	Divtime time.Duration
}

func NewConfigFromEnv() (*Config, error) {
	at, err := time.ParseDuration(os.Getenv("TIME_ADDITION_MS") + "ms")
	if err != nil || at < 0 {
		return nil, fmt.Errorf(errMsgFmt, "TIME_ADDITION_MS")
	}
	st, err := time.ParseDuration(os.Getenv("TIME_SUBTRACTION_MS") + "ms")
	if err != nil || st < 0 {
		return nil, fmt.Errorf(errMsgFmt, "TIME_SUBTRACTION_MS")
	}
	mt, err := time.ParseDuration(os.Getenv("TIME_MULTIPLICATIONS_MS") + "ms")
	if err != nil || mt < 0 {
		return nil, fmt.Errorf(errMsgFmt, "TIME_MULTIPLICATIONS_MS")
	}
	dt, err := time.ParseDuration(os.Getenv("TIME_DIVISIONS_MS") + "ms")
	if err != nil || dt < 0 {
		return nil, fmt.Errorf(errMsgFmt, "TIME_DIVISIONS_MS")
	}

	cfg := Config{
		Addtime: at,
		Subtime: st,
		Multime: mt,
		Divtime: dt,
	}

	return &cfg, nil
}
