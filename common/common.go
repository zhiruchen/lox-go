package common

import (
	"errors"

	"github.com/zhiruchen/lox-go/token"
)

var (
	// ErrFileNotFound file not found
	ErrFileNotFound = errors.New("no such file")
)

// ConditionalExp 三元表达式
func ConditionalExp(condition bool, v1, v2 token.Type) token.Type {
	if condition {
		return v1
	}
	return v2
}
