package utils

import (
	"fmt"
	"strings"
)

// MustTrue
func MustTrue(flat bool, msg string) {
	if !flat {
		panic(msg)
	}
}

func MustNotError(err error) {
	MustTrue(err == nil, fmt.Sprintf("Shouldn't get error. %v", err))
}

func MustNotNil(v interface{}) {
	MustTrue(v == nil, "Shouldn't be nil")
}

func MustNotBlank(s string) {
	MustTrue(s != "" && strings.TrimSpace(s) != "", "Shouldn't be blank.")
}
