package parser

import (
	"fmt"

	"github.com/zhiruchen/lox-go/expr"
	"github.com/zhiruchen/lox-go/lox"
	"github.com/zhiruchen/lox-go/token"
)

type ErrFunc func(tk *token.Token, msg string)

type Parser struct {
	tokens  []*token.Token
	current int
	errFunc ErrFunc
}

func NewParser(tokens []*token.Token, errFunc ErrFunc) *Parser {
	return &Parser{tokens: tokens, errFunc: errFunc}
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

	if p.match(token.Fun) {
		return p.function("function")
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

func (p *Parser) function(kind string) *expr.Function {
	name := p.consume(token.Identifier, "expect " +kind+ "name.")
	p.consume(token.LeftParen, "Expect `(` after "+kind+" name.")

	var params []*token.Token
	if !p.check(token.RightParen) {
		params = append(params, p.consume(token.Identifier, "Expect parameter name."))

		for p.match(token.Comma) {
			if len(params) > 8 {
				p.errFunc(p.peek(), "Cannot have more than 8 parameters.")
			}

			params = append(params, p.consume(token.Identifier, "Expect parameter name."))
		}
	}
	p.consume(token.RightParen, "Expect `)` after parameters")

	p.consume(token.LeftBrace, "Expect `{` before "+kind+ " body.")
	body := p.block()
	return expr.NewFunctionStmt(name, params, body)
}


func (p *Parser) statement() expr.Stmt {
	if p.match(token.For) {
		return p.forStatement()
	}

	if p.match(token.If) {
		return p.ifStatement()
	}

	if p.match(token.Print) {
		return p.printStatement()
	}

	if p.match(token.While) {
		return p.whileStatement()
	}

	if p.match(token.LeftBrace) {
		return expr.NewBlockStmt(p.block())
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() expr.Stmt {
	p.consume(token.LeftParen, "Expect `(` after for.")

	var initializer expr.Stmt

	if p.match(token.Semicolon) {
		initializer = nil
	} else if p.match(token.Var) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var cond expr.Expr
	if !p.check(token.Semicolon) {
		cond = p.expression()
	}
	p.consume(token.Semicolon, "Expect `;` after loop condition.")

	var increment expr.Expr
	if !p.check(token.RightParen) {
		increment = p.expression()
	}
	p.consume(token.RightParen, "Expect `)` after for clauses.")
	body := p.statement()

	if increment != nil {
		body = expr.NewBlockStmt([]expr.Stmt{body, expr.NewExpressionStmt(increment)})
	}

	if cond == nil {
		cond = expr.NewLiteral(true)
	}
	body = expr.NewWhileStmt(cond, body)

	if initializer != nil {
		body = expr.NewBlockStmt([]expr.Stmt{initializer, body})
	}

	return body
}

func (p *Parser) ifStatement() expr.Stmt {
	p.consume(token.LeftParen, `expect "(" after 'if'`)
	cond := p.expression()
	p.consume(token.RightParen, `expect ")" after if condition `)

	thenBranch := p.statement()
	var elseBranch expr.Stmt

	if p.match(token.Else) {
		elseBranch = p.statement()
	}
	return expr.NewIFStmt(cond, thenBranch, elseBranch)
}

func (p *Parser) printStatement() *expr.Print {
	value := p.expression()
	p.consume(token.Semicolon, `Expect ":" after value.`)
	return expr.NewPrintStmt(value)
}

func (p *Parser) whileStatement() expr.Stmt {
	p.consume(token.LeftParen, `expect "(" after 'while'.`)
	cond := p.expression()
	p.consume(token.RightParen, `expect ")" after condition.`)

	body := p.statement()

	return expr.NewWhileStmt(cond, body)
}

func (p *Parser) expressionStatement() *expr.Expression {
	value := p.expression()
	p.consume(token.Semicolon, `Expect ":" after value.`)
	return expr.NewExpressionStmt(value)
}

func (p *Parser) block() []expr.Stmt {
	statements := make([]expr.Stmt, 0)

	for !p.check(token.RightBrace) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(token.RightBrace, `Expect "}" after block!`)
	return statements
}

func (p *Parser) expression() expr.Expr {
	//return p.equality()
	return p.assignment()
}

func (p *Parser) assignment() expr.Expr {
	exp := p.or()
	//exp := p.equality()

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

func (p *Parser) or() expr.Expr {
	exp := p.and()

	for p.match(token.OR) {
		op := p.previous()
		right := p.and()

		exp = expr.NewLogical(exp, right, op)
	}

	return exp
}

func (p *Parser) and() expr.Expr {
	exp := p.equality()

	for p.match(token.And) {
		op := p.previous()
		right := p.equality()

		exp = expr.NewLogical(exp, right, op)
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

	return p.call()
}

func (p *Parser) call() expr.Expr {
	exp := p.primary()

	for {
		if p.match(token.LeftParen) {
			exp = p.finishCall(exp)
		} else {
			break
		}
	}

	return exp
}

func (p *Parser) finishCall(callee expr.Expr) expr.Expr {
	var arguments []expr.Expr

	if !p.check(token.RightParen) {
		if len(arguments) > 8 {
			p.errFunc(p.peek(), "Cannot have more than 8 arguments")
		}

		arguments = append(arguments, p.expression())

		for p.match(token.Comma) {
			arguments = append(arguments, p.expression())
		}
	}

	paren := p.consume(token.RightParen, "Expect `)` after arguments!")
	return expr.NewCall(callee, paren, arguments)
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
