package parser

import (
	"fmt"

	"github.com/zhiruchen/lox-go/expr"
	"github.com/zhiruchen/lox-go/token"
)

type Parser struct {
	tokens  []*token.Token
	current int
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() expr.Expr {
	return p.expression()
}

func (p *Parser) expression() expr.Expr {
	return p.equality()
}

func (p *Parser) equality() expr.Expr {
	exp := p.comparison()

	for p.match(token.BangEqual, token.EqualEqual) {
		op := p.previous()
		right := p.comparison()
		exp = expr.NewBinary(exp, right, op)
	}
	return exp
}

func (p *Parser) comparison() expr.Expr {
	exp := p.term()

	for p.match(token.Greater, token.GreaterEqual, token.Less, token.LessEqual) {
		op := p.previous()
		right := p.term()
		exp = expr.NewBinary(exp, right, op)
	}
	return exp
}

func (p *Parser) term() expr.Expr {
	exp := p.factor()

	for p.match(token.Minus, token.Plus) {
		op := p.previous()
		right := p.factor()
		exp = expr.NewBinary(exp, right, op)
	}

	return exp
}

func (p *Parser) factor() expr.Expr {
	exp := p.unary()

	for p.match(token.Slash, token.Star) {
		op := p.previous()
		right := p.unary()
		exp = expr.NewBinary(exp, right, op)
	}

	return exp
}

func (p *Parser) unary() expr.Expr {
	if p.match(token.Bang, token.Minus) {
		op := p.previous()
		right := p.unary()
		return expr.NewUnary(op, right)
	}

	return p.primary()
}

func (p *Parser) primary() expr.Expr {
	if p.match(token.False) {
		return expr.NewLiteral(false)
	}

	if p.match(token.True) {
		return expr.NewLiteral(true)
	}

	if p.match(token.Nil) {
		return expr.NewLiteral(nil)
	}

	if p.match(token.Number, token.String) {
		return expr.NewLiteral(p.previous().Literal)
	}

	if p.match(token.LeftParen) {
		exp := p.expression()
		p.consume(token.RightParen, "Expect ')' after expression")
		return expr.NewGrouping(exp)
	}
	panic(fmt.Sprintf("expect expression: %s", p.peek().Lexeme))
}

func (p *Parser) consume(t token.Type, msg string) *token.Token {
	if p.check(t) {
		return p.advance()
	}
	panic(msg)
}

func (p *Parser) match(tokenTypes ...token.Type) bool {
	for _, t := range tokenTypes {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t token.Type) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().TokenType == t
}

func (p *Parser) peek() *token.Token {
	return p.tokens[p.current]
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == token.Eof
}

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}
