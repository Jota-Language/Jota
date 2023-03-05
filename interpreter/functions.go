package interpreter

import (
	"jota/ast"
	"jota/environment"
)

type Callable interface {
	Call(interpreter *Interpreter, arguments []any) any
	Arity() int
}

type Function struct {
	Declaration ast.FunctionStatement
    Closure *environment.Environment
}

func (f Function) Call(interpreter *Interpreter, arguments []any) any {
    var returnValue any
    defer func() {
        if r := recover(); r != nil {
            if e, ok := r.(Return); ok {
                returnValue = e.Value
                return
            }
        } 
    }()

    env := environment.NewEnvironment(f.Closure)
	for i, param := range f.Declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}
	interpreter.executeBlock(f.Declaration.Body, env)

	return returnValue
}

func (f Function) Arity() any {
	return len(f.Declaration.Params)
}

func (f Function) String() string {
	return "<fn " + f.Declaration.Name.Lexeme + ">"
}

type BuiltInFunction struct {
	ArityNumber int
	NativeLogic func(interpreter *Interpreter, arguments []any) any
}

func (bif BuiltInFunction) Call(interpreter *Interpreter, arguments []any) any {
	return bif.NativeLogic(interpreter, arguments)
}
func (bif BuiltInFunction) Arity() int {
	return bif.ArityNumber
}
func (bif BuiltInFunction) String() string {
	return "<native fn>"
}
