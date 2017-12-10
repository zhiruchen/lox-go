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
	env *lox.Env
}

func NewInterpreter() *Interpreter {
	return &Interpreter{env: lox.NewEnv()}
}

// Interpret 运行解释器
func (itp *Interpreter) Interpret(statements []expr.Stmt) {
	//fmt.Println(itp.stringify(itp.evaluate(exp)))

	for _, statement := range statements {
		itp.execute(statement)
	}
}

func (itp *Interpreter) execute(stmt expr.Stmt) {
	stmt.Accept(itp)
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

func (itp *Interpreter) VisitorGroupingExpr(exp *expr.Grouping) interface{} {
	return itp.evaluate(exp.Expression)
}

func (itp *Interpreter) VisitorLiteralExpr(exp *expr.Literal) interface{} {
	return exp.Value
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
