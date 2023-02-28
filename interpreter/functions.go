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
}

func (f *Function) Call(interpreter *Interpreter, arguments []any) any {
	env := environment.NewEnvironmentWithEnclosing(interpreter.Environment)
	for i, param := range f.Declaration.Params {
		env.Define(param.Lexeme, arguments[i])
	}

	interpreter.executeBlock(f.Declaration.Body, env)
	return nil
}

func (f *Function) Arity() any {
	return len(f.Declaration.Params)
}

func (f *Function) toString() string {
	return "<fn " + f.Declaration.Name.Lexeme + ">"
}

type BuiltInFunction struct {
	ArityNumber int
	NativeLogic func(interpreter *Interpreter, arguments []any) any
	String      string
}

func (bif *BuiltInFunction) Call(interpreter *Interpreter, arguments []any) any {
	return bif.NativeLogic(interpreter, arguments)
}
func (bif *BuiltInFunction) Arity() int {
	return bif.ArityNumber
}
func (bif *BuiltInFunction) toString() string {
	return "<native fn>"
}
