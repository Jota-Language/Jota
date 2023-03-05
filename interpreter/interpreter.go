package interpreter

import (
	"fmt"
	"jota/ast"
	"jota/environment"
	"jota/errors"
	"math"
	"strings"
	"time"
)

type Interpreter struct {
	Globals      *environment.Environment
	Environment  *environment.Environment
	ErrorHandler errors.ErrorHandler
}

func NewInterpreter(errorHandler errors.ErrorHandler) *Interpreter {
	globals := environment.NewEnvironment(nil)

	// TODO: put these in a separate file!
	globals.Define("clock", &BuiltInFunction{
		ArityNumber: 0,
		NativeLogic: func(interpreter *Interpreter, arguments []any) any {
			return float64(time.Now().UnixNano()) / 1e9 // Returns the elapsed time in seconds.
		},
	})
	globals.Define("milliseconds", &BuiltInFunction{
		ArityNumber: 1,
		NativeLogic: func(interpreter *Interpreter, arguments []any) any {
			value, ok := arguments[0].(float64)
			if !ok {
				return nil
			}
			return float64(math.Round(value*1000*100) / 100)
		},
	})
	globals.Define("stringify", &BuiltInFunction{
		ArityNumber: 1,
		NativeLogic: func(interpreter *Interpreter, arguments []any) any {
			return fmt.Sprint(arguments[0])
		},
	})
	globals.Define("type", &BuiltInFunction{
		ArityNumber: 1,
		NativeLogic: func(interpreter *Interpreter, arguments []any) any {
			if fmt.Sprintf("%T", arguments[0]) == "float64" {
				return "number"
			} else {
				return fmt.Sprintf("%T", arguments[0])
			}
		},
	})

	return &Interpreter{
		Globals:      globals,
		Environment:  globals,
		ErrorHandler: errorHandler,
	}
}

func (i *Interpreter) Interpret(statements []ast.Statement) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(errors.RuntimeError); ok {
				errors.RuntimeErr(e, i.ErrorHandler)
				return
			}
		}
	}()

	for _, statement := range statements {
		i.execute(statement)
	}
	return
}

func (i *Interpreter) VisitIfStatement(statement ast.IfStatement) any {
	if i.isTruthy(i.evaluate(statement.Condition)) {
		i.execute(statement.ThenBranch)
	} else if statement.ElseBranch != nil {
		i.execute(statement.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitExpressionStatement(statement ast.ExpressionStatement) any {
	return i.evaluate(statement.Expression)
}

func (i *Interpreter) VisitPrintStatement(statement ast.PrintStatement) any {
	value := i.evaluate(statement.Expression)
	fmt.Println(i.stringify(value))
	return nil
}

func (i *Interpreter) VisitVariableStatement(statement ast.VariableStatement) any {
	var value any
	if statement.Initializer != nil {
		value = i.evaluate(statement.Initializer)
	}

	i.Environment.Define(statement.Name.Lexeme, value)
	return nil
}

func (i *Interpreter) VisitWhileStatement(statement ast.WhileStatement) any {
	for i.isTruthy(i.evaluate(statement.Condition)) {
		i.execute(statement.Body)
	}
	return nil
}

func (i *Interpreter) VisitBlockStatement(statement ast.BlockStatement) any {
	i.executeBlock(statement.Statements, environment.NewEnvironment(i.Environment))
	return nil
}

func (i *Interpreter) VisitFunctionStatement(statement ast.FunctionStatement) any {
	function := Function{Declaration: statement, Closure: i.Environment}
	i.Environment.Define(statement.Name.Lexeme, function)
	return nil
}

func (i *Interpreter) VisitReturnStatement(statement ast.ReturnStatement) any {
	var value any
	if statement.Value != nil {
		value = i.evaluate(statement.Value)
	}

	panic(Return{Value: value})
}

type Return struct {
	Value any
}

func (i *Interpreter) VisitLiteralExpression(expression ast.Literal) any {
	return expression.Value
}

func (i *Interpreter) VisitLogicalExpression(expression ast.Logical) any {
	left := i.evaluate(expression.Left)

	if expression.Operator.Type == ast.OR {
		if i.isTruthy(left) {
			return left
		}
	} else {
		if !i.isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expression.Right)
}

func (i *Interpreter) VisitGroupingExpression(expression ast.Grouping) any {
	return i.evaluate(expression.Expression)
}

func (i *Interpreter) VisitCallExpression(expression ast.Call) any {
	callee := i.evaluate(expression.Callee)

	arguments := make([]any, len(expression.Arguments))
	for index, argument := range expression.Arguments {
		arguments[index] = i.evaluate(argument)
	}

	function, ok := callee.(Callable)
	if !ok {
		panic(errors.RuntimeError{Token: expression.Paren, Message: "can only call functions and classes"})
	}

	if len(arguments) != function.Arity() {
		panic(errors.RuntimeError{Token: expression.Paren, Message: "expected " + fmt.Sprint(function.Arity()) + " arguments but got " + fmt.Sprint(len(arguments))})
	}

	return function.Call(i, arguments)
}

func (i *Interpreter) VisitUnaryExpression(expression ast.Unary) any {
	right := i.evaluate(expression.Right)

	switch expression.Operator.Type {
	case ast.MINUS:
		i.checkNumberOperand(expression.Operator, right)
		return -right.(float64)
	case ast.BANG:
		return !i.isTruthy(right)
	case ast.INCREMENT:
		i.checkNumberOperand(expression.Operator, right)
		return right.(float64) + 1
	case ast.DECREMENT:
		i.checkNumberOperand(expression.Operator, right)
		return right.(float64) - 1
	}

	return nil
}

func (i *Interpreter) VisitVariableExpression(expression ast.Variable) any {
	return i.Environment.Values[expression.Name.Lexeme]
}

func (i *Interpreter) VisitBinaryExpression(expression ast.Binary) any {
	left := i.evaluate(expression.Left)
	right := i.evaluate(expression.Right)

	switch expression.Operator.Type {
	case ast.MINUS:
		i.checkNumberOperands(expression.Operator, left, right)
		return left.(float64) - right.(float64)
	case ast.SLASH:
		i.checkNumberOperands(expression.Operator, left, right)
		return left.(float64) / right.(float64)
	case ast.PERCENT:
		i.checkNumberOperands(expression.Operator, left, right)
		return math.Mod(left.(float64), right.(float64))
	case ast.ASTERISK:
		i.checkNumberOperands(expression.Operator, left, right)
		return left.(float64) * right.(float64)
	case ast.CARET:
		i.checkNumberOperands(expression.Operator, left, right)
		return math.Pow(left.(float64), right.(float64))
	case ast.PLUS:
		if lf, lok := left.(float64); lok {
			if rf, rok := right.(float64); rok {
				return lf + rf
			}
		}
		if lf, lok := left.(string); lok {
			if rf, rok := right.(string); rok {
				return lf + rf
			}
		}
		panic(errors.RuntimeError{Token: expression.Operator, Message: "operands must be either two numbers or two strings"})
	case ast.GREATER:
		i.checkNumberOperands(expression.Operator, left, right)
		return left.(float64) > right.(float64)
	case ast.GREATER_EQUAL:
		i.checkNumberOperands(expression.Operator, left, right)
		return left.(float64) >= right.(float64)
	case ast.LESS:
		i.checkNumberOperands(expression.Operator, left, right)
		return left.(float64) < right.(float64)
	case ast.LESS_EQUAL:
		i.checkNumberOperands(expression.Operator, left, right)
		return left.(float64) <= right.(float64)
	case ast.EQUAL_EQUAL:
		return i.isEqual(left, right)
	case ast.BANG_EQUAL:
		return !i.isEqual(left, right)
	}

	return nil
}

func (i *Interpreter) VisitAssignExpression(expression ast.Assign) any {
	value := i.evaluate(expression.Value)
	i.Environment.Assign(expression.Name, value)
	return value
}

func (i *Interpreter) evaluate(expression ast.Expression) any {
	return expression.Accept(i)
}

func (i *Interpreter) execute(statement ast.Statement) any {
	return statement.Accept(i)
}

func (i *Interpreter) executeBlock(statements []ast.Statement, environment *environment.Environment) {
	previous := i.Environment
	defer func() {
		i.Environment = previous
	}()

	i.Environment = environment
	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) isTruthy(object any) bool {
	if object == nil {
		return false
	}
	if boolean, ok := object.(bool); ok {
		return boolean
	}

	return true
}

func (i *Interpreter) isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	return a == b
}

func (i *Interpreter) checkNumberOperand(operator ast.Token, operand any) {
	if _, ok := operand.(float64); ok {
		return
	}

	panic(errors.RuntimeError{Token: operator, Message: "operand must be a number"})
}

func (i *Interpreter) checkNumberOperands(operator ast.Token, left, right any) {
	if _, lok := left.(float64); lok {
		if _, rok := right.(float64); rok {
			return
		}
	}

	panic(errors.RuntimeError{Token: operator, Message: "operands must be numbers"})
}

func (i *Interpreter) stringify(object any) string {
	if object == nil {
		return "nil"
	}

	if number, ok := object.(float64); ok {
		text := fmt.Sprint(number)
		if strings.HasSuffix(text, ".0") {
			text = text[0 : len(text)-2]
		}
		return text
	}

	return fmt.Sprint(object)
}
