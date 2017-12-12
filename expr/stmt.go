package expr

import "github.com/zhiruchen/lox-go/token"

type StmtVisitor interface {
	VisitorExpressionStmtExpr(expr *Expression) interface{}
	VisitorPrintStmtExpr(expr *Print) interface{}
	VisitorVarStmtExpr(expr *Var) interface{}
	VisitorBlockStmtExpr(expr *Block) interface{}
}

type Stmt interface {
	Accept(v Visitor) interface{}
}

type Expression struct {
	Expression Expr
}

type Print struct {
	Print Expr
}

type Var struct {
	Name        *token.Token
	Initializer Expr
}

type Block struct {
	Statements []Stmt
}

func (st *Expression) Accept(v Visitor) interface{} {
	return v.VisitorExpressionStmtExpr(st)
}

func (st *Print) Accept(v Visitor) interface{} {
	return v.VisitorPrintStmtExpr(st)
}

func (st *Var) Accept(v Visitor) interface{} {
	return v.VisitorVarStmtExpr(st)
}

func (st *Block) Accept(v Visitor) interface{} {
	return v.VisitorBlockStmtExpr(st)
}

func NewPrintStmt(e Expr) *Print {
	return &Print{Print: e}
}

func NewExpressionStmt(e Expr) *Expression {
	return &Expression{Expression: e}
}

func NewVarStmt(name *token.Token, e Expr) *Var {
	return &Var{Name: name, Initializer: e}
}

func NewBlockStmt(stmts []Stmt) *Block {
	return &Block{Statements: stmts}
}
