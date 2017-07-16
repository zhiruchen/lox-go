package lox

import (
	"fmt"

	"github.com/zhiruchen/lox-go/token"
)

// Lox the lox lang
type Lox struct {
	HasError       bool
	HasRuntimeErro bool
}

func runFile(path string) error {
	return nil
}

func runPromt() error {
	return nil
}

// LineError lox line error
func LineError(line int, message string) {
	report(line, "", message)
}

// TokenError lox token error
func TokenError(tk token.Token, message string) {
	if tk.TokenType == token.Eof {
		report(tk.Line, " is at end", message)
	} else {
		report(tk.Line, " at '"+tk.Lexeme+"'", message)
	}
}

func report(line int, where string, message string) {
	fmt.Printf("[line %d] where: %s: %s\n", line, where, message)
}
