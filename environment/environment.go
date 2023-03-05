package environment

import (
	"jota/ast"
	"jota/errors"
)

type Environment struct {
	Enclosing *Environment
	Values    map[string]any
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{Enclosing: enclosing, Values: make(map[string]any)}
}

func (e *Environment) Define(name string, value any) {
	e.Values[name] = value
}

func (e *Environment) Get(name ast.Token) any {
	if value, ok := e.Values[name.Lexeme]; ok {
		return value
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	panic(errors.RuntimeError{Token: name, Message: "Undefined variable '" + name.Lexeme + "'"})
}

func (e *Environment) Assign(name ast.Token, value any) {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Define(name.Lexeme, value)
		return
	}

	if e.Enclosing != nil {
		e.Enclosing.Assign(name, value)
		return
	}

	panic(errors.RuntimeError{Token: name, Message: "Undefined variable '" + name.Lexeme + "'"})
}
