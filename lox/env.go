package lox

import (
	"fmt"

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
