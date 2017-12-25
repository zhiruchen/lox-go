package interpreter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zhiruchen/lox-go/expr"
	"github.com/zhiruchen/lox-go/lox"
	"github.com/zhiruchen/lox-go/token"
)

// Interpreter the lox lang interpreter
type Interpreter struct {
	env     *Env
	globals *Env
}

func NewInterpreter() *Interpreter {
	globals := NewEnv()
	globals.Define("clock", &CLock{})

	return &Interpreter{env: globals, globals: globals}
}

// Interpret 运行解释器
func (itp *Interpreter) Interpret(statements []expr.Stmt) {
	//fmt.Println(itp.stringify(itp.evaluate(exp)))

	for _, statement := range statements {
		itp.execute(statement)
	}
}

func (itp *Interpreter) GetGlobalEnv() *Env {
	return itp.globals
}

func (itp *Interpreter) execute(stmt expr.Stmt) {
	stmt.Accept(itp)
}

func (itp *Interpreter) VisitorBlockStmtExpr(expr *expr.Block) interface{} {
	itp.executeBlock(expr.Statements, NewEnvWithEnclosing(itp.env))
	return nil
}

func (itp *Interpreter) ExecuteBlock(statements []expr.Stmt, env *Env) {
	itp.executeBlock(statements, env)
}

func (itp *Interpreter) executeBlock(statements []expr.Stmt, env *Env) {
	previous := itp.env
	defer func() {
		itp.env = previous
	}()

	itp.env = env
	for _, s := range statements {
		itp.execute(s)
	}
}

func (itp *Interpreter) VisitorBinaryExpr(exp *expr.Binary) interface{} {
	left, right := itp.evaluate(exp.Left), itp.evaluate(exp.Right)

	switch exp.Operator.TokenType {
	case token.Minus:
		v1, v2 := itp.checkNumberOperands(*exp.Operator, left, right)
		return v1 - v2
	case token.Slash:
		v1, v2 := itp.checkNumberOperands(*exp.Operator, left, right)
		return v1 / v2
	case token.Star:
		v1, v2 := itp.checkNumberOperands(*exp.Operator, left, right)
		return v1 * v2
	case token.Plus:
		v, ok := left.(float64)
		v1, ok1 := right.(float64)
		if ok && ok1 {
			return v + v1
		}
		v2, ok2 := left.(string)
		v3, ok3 := right.(string)
		if ok2 && ok3 {
			return v2 + v3
		}
		panic(fmt.Sprintf("%s Operands must be two numbers or two strings!", exp.Operator.Lexeme))
	case token.Greater:
		v1, v2 := itp.checkNumberOperands(*exp.Operator, left, right)
		return v1 > v2
	case token.GreaterEqual:
		v1, v2 := itp.checkNumberOperands(*exp.Operator, left, right)
		return v1 >= v2
	case token.Less:
		v1, v2 := itp.checkNumberOperands(*exp.Operator, left, right)
		return v1 < v2
	case token.LessEqual:
		v1, v2 := itp.checkNumberOperands(*exp.Operator, left, right)
		return v1 <= v2
	case token.EqualEqual:
		return itp.isEqual(left, right)
	case token.BangEqual:
		return !itp.isEqual(left, right)
	default:
		return nil
	}
}

func (itp *Interpreter) VisitorCallExpr(expr *expr.Call) interface{} {
	callee := itp.evaluate(expr.Callee)

	var arguments []interface{}
	for _, arg := range expr.Arguments {
		arguments = append(arguments, itp.evaluate(arg))
	}

	function, ok := callee.(Callable)
	if !ok {
		panic(lox.NewRuntimeError(expr.Paren, "Can only call functions and classes"))
	}

	if len(arguments) != function.Arity() {
		panic(lox.RuntimeError{expr.Paren, fmt.Sprintf("Expected %d arguments but got %d", function.Arity(), len(arguments))})
	}

	return function.Call(itp, arguments)
}

func (itp *Interpreter) VisitorGroupingExpr(exp *expr.Grouping) interface{} {
	return itp.evaluate(exp.Expression)
}

func (itp *Interpreter) VisitorLiteralExpr(exp *expr.Literal) interface{} {
	return exp.Value
}

func (itp *Interpreter) VisitorLogicalExpr(expr *expr.Logical) interface{} {
	left := itp.evaluate(expr.Left)

	if expr.Operator.TokenType == token.OR {
		if itp.isTruthy(left) {
			return left
		}
	} else {
		if !itp.isTruthy(left) {
			return left
		}
	}

	return itp.evaluate(expr.Right)
}

func (itp *Interpreter) VisitorUnaryExpr(exp *expr.Unary) interface{} {
	right := itp.evaluate(exp.Right)

	switch exp.Operator.TokenType {
	case token.Bang:
		return !itp.isTruthy(right)
	case token.Minus:
		v := itp.checkNumberOperand(*exp.Operator, right)
		return 0 - v
	}
	return nil
}

func (itp *Interpreter) VisitorVariableExpr(exp *expr.Variable) interface{} {
	return itp.env.Get(exp.Name)
}

func (itp *Interpreter) VisitorExpressionStmtExpr(expr *expr.Expression) interface{} {
	itp.evaluate(expr.Expression)
	return nil
}

func (itp *Interpreter) VisitorFunStmtExpr(stmt *expr.Function) interface{} {
	function := NewFunction(stmt)
	itp.env.Define(stmt.Name.Lexeme, function)
	return nil
}

func (itp *Interpreter) VisitorIFStmtExpr(stmt *expr.IF) interface{} {
	if itp.isTruthy(itp.evaluate(stmt.Condition)) {
		itp.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		itp.execute(stmt.ElseBranch)
	}

	return nil
}

func (itp *Interpreter) VisitorPrintStmtExpr(expr *expr.Print) interface{} {
	value := itp.evaluate(expr.Print)
	fmt.Printf("%s\n", itp.stringify(value))
	return nil
}

func (itp *Interpreter) VisitorVarStmtExpr(expr *expr.Var) interface{} {
	var value interface{}

	if expr.Initializer != nil {
		value = itp.evaluate(expr.Initializer)
	}

	itp.env.Define(expr.Name.Lexeme, value)
	return nil
}

func (itp *Interpreter) VisitorWhileStmtExpr(stmt *expr.While) interface{} {
	for itp.isTruthy(itp.evaluate(stmt.Condition)) {
		itp.execute(stmt.Body)
	}

	return nil
}

func (itp *Interpreter) VisitorAssignExpr(expr *expr.Assign) interface{} {
	value := itp.evaluate(expr.Value)

	itp.env.Assign(expr.Name, value)
	return value
}

func (itp *Interpreter) checkNumberOperand(operator token.Token, obj interface{}) float64 {
	v, ok := obj.(float64)
	if !ok {
		panic(fmt.Sprintf("%s Operand must be a number.", operator.Lexeme))
	}
	return v
}

func (itp *Interpreter) checkNumberOperands(operator token.Token, left, right interface{}) (float64, float64) {
	v1, ok := left.(float64)
	v2, ok1 := right.(float64)
	if !ok || !ok1 {
		panic(fmt.Sprintf("%s Operands must be numbers!", operator.Lexeme))
	}
	return v1, v2
}

func (itp *Interpreter) isEqual(left, right interface{}) bool {
	if left == nil && right == nil {
		return true
	}

	if left == nil {
		return false
	}

	return left == right
}

func (itp *Interpreter) isTruthy(obj interface{}) bool {
	if obj == nil {
		return false
	}

	v, ok := obj.(bool)
	if ok {
		return v
	}
	return true
}
func (itp *Interpreter) evaluate(exp expr.Expr) interface{} {
	return exp.Accept(itp)
}

func (itp *Interpreter) stringify(obj interface{}) string {
	if obj == nil {
		return "nil"
	}

	v, ok := obj.(float64)
	if ok {
		text := strconv.FormatFloat(v, 'f', 6, 64)
		if strings.HasSuffix(text, ".0") {
			text = text[0 : len(text)-2]
			return text
		}
	}
	return fmt.Sprintf("%v", obj)
}
