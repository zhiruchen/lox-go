package lox

import (
	"github.com/zhiruchen/lox-go/token"
)

type Env struct {
	values map[string]interface{}
}

func NewEnv() *Env {
	return &Env{
		values: make(map[string]interface{}),
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

	panic("Undefined variable '" + name.Lexeme + "'.")
}

func (env *Env) Assign(name *token.Token, value interface{}) {
	if _, ok := env.values[name.Lexeme]; ok {
		env.values[name.Lexeme] = value
		return
	}

	panic("Undefined variable" + name.Lexeme + ".")
}
