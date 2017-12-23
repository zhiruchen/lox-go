package lox

import (
	"github.com/zhiruchen/lox-go/expr"
	itp "github.com/zhiruchen/lox-go/interpreter"
)

type Function struct {
	declaration *expr.Function
}

func NewFunction (declaration *expr.Function) *Function {
	return &Function{declaration:declaration}
}

func (f *Function) Arity() int {
	return len(f.declaration.Parameters)
}

func (f *Function) Call(itp *itp.Interpreter, args []interface{}) interface{} {
	env := NewEnvWithEnclosing(itp.Globals)

	for i, param := range f.declaration.Parameters {
		env.Define(param.Lexeme, args[i])
	}

	itp.ExecuteBlock(f.declaration.Body, env)
	return nil
}

func (f *Function) String() string {
	return "<fn " + f.declaration.Name.Lexeme +">"
}