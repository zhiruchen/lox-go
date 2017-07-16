package token

// Type token 类型
type Type int

const (
	// LeftParen 左括号
	LeftParen Type = iota
	// RightParen 右括号
	RightParen

	LeftBrace
	RightBrace

	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	Identifier
	String
	Number

	// KeyWords
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	OR
	Print
	Return
	Super
	This
	True
	Var
	While

	Eof
)
