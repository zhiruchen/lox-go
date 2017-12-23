package expr

import (
	"github.com/zhiruchen/lox-go/token"
)

type Visitor interface {
	VisitorBinaryExpr(expr *Binary) interface{}
	VisitorGroupingExpr(expr *Grouping) interface{}
	VisitorLiteralExpr(expr *Literal) interface{}
	VisitorLogicalExpr(expr *Logical) interface{}
	VisitorUnaryExpr(expr *Unary) interface{}
	VisitorVariableExpr(expr *Variable) interface{}
	VisitorAssignExpr(expr *Assign) interface{}
	VisitorCallExpr(expr *Call) interface{}

	StmtVisitor
}

type Expr interface {
	Accept(v Visitor) interface{}
}

type Assign struct {
	Name  *token.Token
	Value Expr
}

func NewAssign(name *token.Token, value Expr) *Assign {
	return &Assign{Name: name, Value: value}
}

func (as *Assign) Accept(v Visitor) interface{} {
	return v.VisitorAssignExpr(as)
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

type Call struct {
	Callee    Expr
	Paren     *token.Token
	Arguments []Expr
}

func NewCall(callee Expr, paren *token.Token, arguments []Expr) *Call {
	return &Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}
}

func (cl *Call) Accept(v Visitor) interface{} {
	return v.VisitorCallExpr(cl)
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

type Logical struct {
	Left     Expr
	Right    Expr
	Operator *token.Token
}

func NewLogical(left, right Expr, op *token.Token) *Logical {
	return &Logical{Left: left, Right: right, Operator: op}
}

func (l *Logical) Accept(v Visitor) interface{} {
	return v.VisitorLogicalExpr(l)
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
