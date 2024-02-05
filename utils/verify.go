package utils

import (
	"errors"
	"strings"
)

func NotEmpty(s string) error {
	if strings.TrimSpace(s) == "" {
		return errors.New("cannot be empty")
	}
	return nil
}

func NotContainSpace(s string) error {
	if strings.Contains(s, " ") {
		return errors.New("cannot contain space")
	}
	return nil
}

func IsClusterToken(s string) error {
	const prefix = "pds-g^KU_"
	if strings.HasPrefix(s, prefix) {
		return nil
	}
	return errors.New("not a valid cluster token")
}

func Unique(items []string) func(string) error {
	return func(s string) error {
		for _, item := range items {
			if item == s {
				return errors.New("already exists")
			}
		}
		return nil
	}
}
