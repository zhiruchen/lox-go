package interpreter

import (
	"github.com/zhiruchen/lox-go/token"
)

type Env struct {
	Enclosing *Env
	values    map[string]interface{}
}

func NewEnv() *Env {
	return &Env{
		values: make(map[string]interface{}),
	}
}

// NewEnvWithEnclosing  new env with enclosing
func NewEnvWithEnclosing(enclosing *Env) *Env {
	return &Env{
		Enclosing: enclosing,
		values:    make(map[string]interface{}),
	}
}

func (env *Env) Define(name string, value interface{}) {
	env.values[name] = value
}

func (env *Env) Get(name *token.Token) interface{} {
	v, ok := env.values[name.Lexeme]
	if ok {
		return v
	}

	if env.Enclosing != nil {
		return env.Enclosing.Get(name)
	}

	panic("Undefined variable '" + name.Lexeme + "'.")
}

func (env *Env) Assign(name *token.Token, value interface{}) {
	if _, ok := env.values[name.Lexeme]; ok {
		env.values[name.Lexeme] = value
		return
	}

	if env.Enclosing != nil {
		env.Enclosing.Assign(name, value)
		return
	}

	panic("Undefined variable" + name.Lexeme + ".")
}
