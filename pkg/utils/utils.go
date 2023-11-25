package utils

import (
	"fmt"
	"log/slog"
	"os"
)

func GetEnvDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	} else {
		return defaultValue
	}
}

func WithErr(err error) slog.Attr {
	return slog.String("err", err.Error())
}

func DoWithAttempts[v any](f func() (v, error), attempts int) (res v, err error) {
	for i := 0; i < attempts; i += 1 {
		v, err := f()
		if err == nil {
			return v, nil
		}
	}
	return res, fmt.Errorf("attempts exeeded: %v", err)
}
