package parser

import (
	"fmt"

	"github.com/zhiruchen/lox-go/expr"
	"github.com/zhiruchen/lox-go/lox"
	"github.com/zhiruchen/lox-go/token"
)

type Parser struct {
	tokens  []*token.Token
	current int
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() []expr.Stmt {

	statements := make([]expr.Stmt, 0)
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) declaration() expr.Stmt {
	if p.match(token.Var) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() expr.Stmt {
	name := p.consume(token.Identifier, "Expect variable name.")

	var initializer expr.Expr
	if p.match(token.Equal) {
		initializer = p.expression()
	}

	p.consume(token.Semicolon, `Expect ';' after variable declaration.`)
	return expr.NewVarStmt(name, initializer)
}

func (p *Parser) statement() expr.Stmt {
	if p.match(token.Print) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() *expr.Print {
	value := p.expression()
	p.consume(token.Semicolon, `Expect ":" after value.`)
	return expr.NewPrintStmt(value)
}

func (p *Parser) expressionStatement() *expr.Expression {
	value := p.expression()
	p.consume(token.Semicolon, `Expect ":" after value.`)
	return expr.NewExpressionStmt(value)
}

func (p *Parser) expression() expr.Expr {
	//return p.equality()
	return p.assignment()
}

func (p *Parser) assignment() expr.Expr {
	exp := p.equality()

	if p.match(token.Equal) {
		equals := p.previous()
		value := p.assignment()

		if v, ok := exp.(*expr.Variable); ok {
			name := v.Name
			return expr.NewAssign(name, value)
		}

		lox.TokenError(equals, "Invalid Assignment target.")
	}

	return exp
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

	if p.match(token.Identifier) {
		return expr.NewVariable(p.previous())
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

func (p *Parser) error() {

}
