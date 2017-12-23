package builtin

import (
	"time"

	"github.com/zhiruchen/lox-go/interpreter"
)

type Lock struct {}

func (l *Lock) Arity() int {
	return 0
}

func (l *Lock) Call(interpreter *interpreter.Interpreter, args []interface{}) interface{} {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

