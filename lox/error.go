package lox

import (
	"github.com/zhiruchen/lox-go/token"
)

type RuntimeError struct {
	Tk  *token.Token
	Msg string
}

func NewRuntimeError(tk *token.Token, msg string) *RuntimeError {
	return &RuntimeError{Tk: tk, Msg: msg}
}
