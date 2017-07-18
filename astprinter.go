package main

import (
	"fmt"

	"strings"

	"github.com/zhiruchen/lox-go/expr"
	"github.com/zhiruchen/lox-go/token"
)

type AstPrinter struct{}

func (ast *AstPrinter) print(e expr.Expr) string {
	v := e.Accept(ast)
	v1, ok := v.(string)
	if ok {
		return v1
	}
	return ""
}

func (ast *AstPrinter) VisitorBinaryExpr(e *expr.Binary) interface{} {
	return ast.parenthesize(e.Operator.Lexeme, e.Left, e.Right)
}

func (ast *AstPrinter) VisitorGroupingExpr(e *expr.Grouping) interface{} {
	return ast.parenthesize("group", e.Expression)
}

func (ast *AstPrinter) VisitorLiteralExpr(e *expr.Literal) interface{} {
	return fmt.Sprintf("%v", e.Value)
}

func (ast *AstPrinter) VisitorUnaryExpr(e *expr.Unary) interface{} {
	return ast.parenthesize(e.Operator.Lexeme, e.Right)
}

func (ast *AstPrinter) parenthesize(name string, exprs ...expr.Expr) string {
	var ll = []string{}

	ll = append(ll, fmt.Sprintf("(%s", name))
	for _, ex := range exprs {
		ll = append(ll, " ")
		v, ok := ex.Accept(ast).(string)
		if ok {
			ll = append(ll, v)
		} else {
			ll = append(ll, "")
		}
	}
	ll = append(ll, ")")
	return strings.Join(ll, "")
}

func main() {
	exp := expr.NewBinary(
		expr.NewUnary(
			&token.Token{token.Minus, "-", nil, 1},
			expr.NewLiteral(123),
		),
		expr.NewGrouping(expr.NewLiteral(45.67)),
		&token.Token{},
	)

	ast := &AstPrinter{}
	fmt.Println(ast.print(exp))
}
