package scanner

import (
	"strconv"

	"github.com/zhiruchen/lox-go/common"
	"github.com/zhiruchen/lox-go/lox"
	"github.com/zhiruchen/lox-go/token"
)

// Scanner lox scanner
type Scanner struct {
	source   string
	runes    []rune
	tokens   []*token.Token
	start    int
	current  int
	line     int
	keywords map[string]token.Type
}

// NewScanner a new scanner
func NewScanner(source string) *Scanner {
	return &Scanner{
		source: source,
		runes:  []rune(source),
		tokens: []*token.Token{},
		line:   1,
		keywords: map[string]token.Type{
			"and":    token.And,
			"class":  token.Class,
			"else":   token.Else,
			"false":  token.False,
			"for":    token.For,
			"fun":    token.Fun,
			"if":     token.If,
			"nil":    token.Nil,
			"or":     token.OR,
			"print":  token.Print,
			"super":  token.Super,
			"this":   token.This,
			"return": token.Return,
			"true":   token.True,
			"var":    token.Var,
			"while":  token.While,
		},
	}
}

// ScanTokens 返回扫描到的token列表
func (scan *Scanner) ScanTokens() []*token.Token {
	for !scan.isAtEnd() {
		scan.start = scan.current
		scan.scanToken()
	}
	scan.tokens = append(
		scan.tokens,
		&token.Token{
			TokenType: token.Eof,
			Lexeme:    "",
			Literal:   nil,
			Line:      scan.line,
		},
	)
	return scan.tokens
}

func (scan *Scanner) isAtEnd() bool {
	return scan.current >= len(scan.runes)
}

func (scan *Scanner) scanToken() {
	c := scan.advance()
	switch c {
	case '(':
		scan.addToken(token.LeftParen, nil)
	case ')':
		scan.addToken(token.RightParen, nil)
	case '{':
		scan.addToken(token.LeftBrace, nil)
	case '}':
		scan.addToken(token.RightBrace, nil)
	case ',':
		scan.addToken(token.Comma, nil)
	case '.':
		scan.addToken(token.Dot, nil)
	case '-':
		scan.addToken(token.Minus, nil)
	case '+':
		scan.addToken(token.Plus, nil)
	case ';':
		scan.addToken(token.Semicolon, nil)
	case '*':
		scan.addToken(token.Star, nil)
	case '!':
		scan.addToken(common.ConditionalExp(scan.match('='), token.BangEqual, token.Bang), nil)
	case '=':
		scan.addToken(common.ConditionalExp(scan.match('='), token.EqualEqual, token.Equal), nil)
	case '<':
		scan.addToken(common.ConditionalExp(scan.match('='), token.LessEqual, token.Less), nil)
	case '>':
		scan.addToken(common.ConditionalExp(scan.match('='), token.GreaterEqual, token.Greater), nil)
	case '/':
		if scan.match('/') {
			scan.skipLineComment()
		} else if scan.match('*') {
			scan.skipBlockComment()
		} else {
			scan.addToken(token.Slash, nil)
		}
	case ' ', '\r', '\t': // 自动 break
	case '\n':
		scan.line++
	case '"':
		scan.getStr()
	default:
		if isDigits(c) {
			scan.getNumber()
		} else if isAlpha(c) {
			scan.getIdentifier()
		} else {
			lox.LineError(scan.line, "Unexpected token!")
		}
	}
}

func (scan *Scanner) advance() rune {
	scan.current++
	return scan.runes[scan.current-1]
}

func (scan *Scanner) addToken(tokenType token.Type, literal interface{}) {
	text := string(scan.runes[scan.start:scan.current])
	scan.tokens = append(scan.tokens, &token.Token{TokenType: tokenType, Lexeme: text, Literal: literal, Line: scan.line})
}

func (scan *Scanner) match(expected rune) bool {
	if scan.isAtEnd() {
		return false
	}

	if scan.runes[scan.current] != expected {
		return false
	}

	scan.current++
	return true
}

func (scan *Scanner) peek() rune {
	if scan.current >= len(scan.runes) {
		return '\000' // https://stackoverflow.com/questions/38007361/is-there-anyway-to-create-null-terminated-string-in-go
	}
	return scan.runes[scan.current]
}

func (scan *Scanner) peekNext() rune {
	if (scan.current + 1) >= len(scan.runes) {
		return '\000'
	}
	return scan.runes[scan.current+1]
}

func (scan *Scanner) skipLineComment() {
	for scan.peek() != '\n' && !scan.isAtEnd() {
		scan.advance()
	}
}

// skipBlockComment 跳过块注释
// https://github.com/munificent/wren/blob/master/src/vm/wren_compiler.c#L660
func (scan *Scanner) skipBlockComment() {
	var nesting = 1
	for nesting > 0 {
		if scan.isAtEnd() {
			lox.LineError(scan.line, "Unterminated block comment!")
			return
		}

		if scan.peek() == '\n' {
			scan.line++
		}

		if scan.peek() == '/' && scan.peekNext() == '*' {
			scan.advance()
			scan.advance()
			nesting++
			continue
		}

		if scan.peek() == '*' && scan.peekNext() == '/' {
			scan.advance()
			scan.advance()
			nesting--
			continue
		}

		scan.advance()
	}
}

func (scan *Scanner) getStr() {
	for scan.peek() != '"' && !scan.isAtEnd() {
		if scan.peek() == '\n' {
			scan.line++
		}
		scan.advance()
	}
	if scan.isAtEnd() {
		lox.LineError(scan.line, "Unterminated string")
		return
	}

	scan.advance()

	value := string(scan.runes[scan.start+1 : scan.current-1])
	scan.addToken(token.String, value)
}

func (scan *Scanner) getNumber() {
	for isDigits(scan.peek()) {
		scan.advance()
	}

	if scan.peek() == '.' && isDigits(scan.peekNext()) {
		scan.advance()

		for isDigits(scan.peek()) {
			scan.advance()
		}
	}

	text := string(scan.runes[scan.start:scan.current])
	number, _ := strconv.ParseFloat(text, 64)
	scan.addToken(token.Number, number)
}

func (scan *Scanner) getIdentifier() {
	for isAlphaNumberic(scan.peek()) {
		scan.advance()
	}

	text := string(scan.runes[scan.start:scan.current])
	tokenType, ok := scan.keywords[text]
	if !ok {
		tokenType = token.Identifier
	}
	scan.addToken(tokenType, nil)
}

func isDigits(c rune) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '_'
}

func isAlphaNumberic(c rune) bool {
	return isAlpha(c) || isDigits(c)
}
