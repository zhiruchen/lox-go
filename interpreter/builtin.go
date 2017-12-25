package interpreter

import (
	"time"
)

type CLock struct{}

func (l *CLock) Arity() int {
	return 0
}

func (l *CLock) Call(itp Interpreter, args []interface{}) interface{} {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
