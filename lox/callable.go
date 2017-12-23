package lox

import (
	"github.com/zhiruchen/lox-go/interpreter"
)

type Callable interface {
	Arity() int
	Call(interpreter *interpreter.Interpreter, args []interface{}) interface{}
}