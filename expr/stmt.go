package expr

import (
	"github.com/zhiruchen/lox-go/token"
)

type StmtVisitor interface {
	VisitorExpressionStmtExpr(expr *Expression) interface{}
	VisitorPrintStmtExpr(expr *Print) interface{}
	VisitorReturnStmtExpr(expr *Return) interface{}
	VisitorVarStmtExpr(expr *Var) interface{}
	VisitorWhileStmtExpr(expr *While) interface{}
	VisitorBlockStmtExpr(expr *Block) interface{}
	VisitorIFStmtExpr(expr *IF) interface{}
	VisitorFunStmtExpr(expr *Function) interface{}
}

type Stmt interface {
	Accept(v Visitor) interface{}
}

type Expression struct {
	Expression Expr
}

type Function struct {
	Name       *token.Token
	Parameters []*token.Token
	Body       []Stmt
}

type IF struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

type Print struct {
	Print Expr
}

type Return struct {
	Keyword *token.Token
	Value   Expr
}

type Var struct {
	Name        *token.Token
	Initializer Expr
}

type While struct {
	Condition Expr
	Body      Stmt
}

type Block struct {
	Statements []Stmt
}

func (st *Expression) Accept(v Visitor) interface{} {
	return v.VisitorExpressionStmtExpr(st)
}

func (st *Function) Accept(v Visitor) interface{} {
	return v.VisitorFunStmtExpr(st)
}

func (st *IF) Accept(v Visitor) interface{} {
	return v.VisitorIFStmtExpr(st)
}

func (st *Print) Accept(v Visitor) interface{} {
	return v.VisitorPrintStmtExpr(st)
}

func (st *Return) Accept(v Visitor) interface{} {
	return v.VisitorReturnStmtExpr(st)
}

func (st *Var) Accept(v Visitor) interface{} {
	return v.VisitorVarStmtExpr(st)
}

func (st *While) Accept(v Visitor) interface{} {
	return v.VisitorWhileStmtExpr(st)
}

func (st *Block) Accept(v Visitor) interface{} {
	return v.VisitorBlockStmtExpr(st)
}

func NewPrintStmt(e Expr) *Print {
	return &Print{Print: e}
}

func NewReturnStmt(keyword *token.Token, value Expr) *Return {
	return &Return{Keyword: keyword, Value: value}
}

func NewExpressionStmt(e Expr) *Expression {
	return &Expression{Expression: e}
}

func NewIFStmt(cond Expr, thenBranch, elseBranch Stmt) *IF {
	return &IF{Condition: cond, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func NewFunctionStmt(name *token.Token, params []*token.Token, body []Stmt) *Function {
	return &Function{Name: name, Parameters: params, Body: body}
}

func NewVarStmt(name *token.Token, e Expr) *Var {
	return &Var{Name: name, Initializer: e}
}

func NewWhileStmt(cond Expr, body Stmt) *While {
	return &While{Condition: cond, Body: body}
}

func NewBlockStmt(stmts []Stmt) *Block {
	return &Block{Statements: stmts}
}
