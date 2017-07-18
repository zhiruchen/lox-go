package main

import (
	"fmt"

	"strings"

	"github.com/zhiruchen/lox-go/expr"
	"github.com/zhiruchen/lox-go/token"
)

type rpnPrinter struct{}

func (rpn *rpnPrinter) print(exp expr.Expr) string {
	v := exp.Accept(rpn)
	v1, ok := v.(string)
	if ok {
		return v1
	}
	return ""
}

func (rpn *rpnPrinter) VisitorBinaryExpr(e *expr.Binary) interface{} {
	return rpn.rpn(e.Operator.Lexeme, e.Left, e.Right)
}

func (rpn *rpnPrinter) VisitorGroupingExpr(e *expr.Grouping) interface{} {
	return rpn.rpn("group", e.Expression)
}

func (rpn *rpnPrinter) VisitorLiteralExpr(e *expr.Literal) interface{} {
	return fmt.Sprintf("%v", e.Value)
}

func (rpn *rpnPrinter) VisitorUnaryExpr(e *expr.Unary) interface{} {
	return rpn.rpn(e.Operator.Lexeme, e.Right)
}

func (rpn *rpnPrinter) rpn(name string, exprs ...expr.Expr) string {
	var ll = []string{}

	for _, exp := range exprs {
		v, ok := exp.Accept(rpn).(string)
		if ok {
			ll = append(ll, v)
		} else {
			ll = append(ll, "")
		}
		ll = append(ll, " ")
	}
	ll = append(ll, fmt.Sprintf("%s", name))

	return strings.Join(ll, "")
}

func main() {
	left := expr.NewBinary(
		expr.NewLiteral(1),
		expr.NewLiteral(2),
		&token.Token{TokenType: token.Plus, Lexeme: "+", Literal: nil, Line: 1},
	)

	right := expr.NewBinary(
		expr.NewLiteral(4),
		expr.NewLiteral(3),
		&token.Token{TokenType: token.Minus, Lexeme: "-", Literal: nil, Line: 1},
	)

	exp := expr.NewBinary(
		left,
		right,
		&token.Token{TokenType: token.Star, Lexeme: "*", Literal: nil, Line: 1},
	)

	rpn := &rpnPrinter{}
	fmt.Println(rpn.print(exp)) // 1 2 + 4 3 - *
}
