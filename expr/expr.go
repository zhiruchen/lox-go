package expr

import (
	"github.com/zhiruchen/lox-go/token"
)

type Visitor interface {
	VisitorBinaryExpr(expr *Binary) interface{}
	VisitorGroupingExpr(expr *Grouping) interface{}
	VisitorLiteralExpr(expr *Literal) interface{}
	VisitorUnaryExpr(expr *Unary) interface{}
	VisitorVariableExpr(expr *Variable) interface{}
	StmtVisitor
}

type Expr interface {
	Accept(v Visitor) interface{}
}

type Binary struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func NewBinary(left, right Expr, operator *token.Token) *Binary {
	return &Binary{Left: left, Operator: operator, Right: right}
}

func (bin *Binary) Accept(v Visitor) interface{} {
	return v.VisitorBinaryExpr(bin)
}

type Grouping struct {
	Expression Expr
}

func NewGrouping(expr Expr) *Grouping {
	return &Grouping{Expression: expr}
}

func (g *Grouping) Accept(v Visitor) interface{} {
	return v.VisitorGroupingExpr(g)
}

type Literal struct {
	Value interface{}
}

func NewLiteral(v interface{}) *Literal {
	return &Literal{Value: v}
}

func (l *Literal) Accept(v Visitor) interface{} {
	return v.VisitorLiteralExpr(l)
}

type Unary struct {
	Operator *token.Token
	Right    Expr
}

func NewUnary(operator *token.Token, right Expr) *Unary {
	return &Unary{Operator: operator, Right: right}
}

func (u *Unary) Accept(v Visitor) interface{} {
	return v.VisitorUnaryExpr(u)
}

type Variable struct {
	Name *token.Token
}

func NewVariable(name *token.Token) *Variable {
	return &Variable{Name: name}
}

func (v *Variable) Accept(vt Visitor) interface{} {
	return vt.VisitorVariableExpr(v)
}
