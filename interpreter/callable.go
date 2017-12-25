package interpreter

type Callable interface {
	Arity() int
	Call(itp *Interpreter, args []interface{}) interface{}
}
