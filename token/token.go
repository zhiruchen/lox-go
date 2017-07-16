package token

import "fmt"

// Token lexeme
type Token struct {
	TokenType Type
	Lexeme    string
	Literal   interface{}
	Line      int
}

// ToString token的字符串表示
func (tk *Token) ToString() string {
	return fmt.Sprintf("type: %d lexeme: %s literal: %v", int(tk.TokenType), tk.Lexeme, tk.Literal)
}
