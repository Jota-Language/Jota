package errors

import (
	"jota/ast"
	"jota/utils"
	"log"
	"strconv"
)

type ErrorHandler struct {
	Error        bool
	RuntimeError bool
	Log          log.Logger
}

func Err(token ast.Token, message string, handler ErrorHandler) {
	switch token.Type {
	case ast.EOF:
		report(token.Line, "at end", message, handler)
	default:
		report(token.Line, "at '"+token.Lexeme+"'", message, handler)
	}
}

func ErrWithoutToken(line int, message string, handler ErrorHandler) {
	report(line, "", message, handler)
}

type RuntimeError struct {
	Token   ast.Token
	Message string
}

func RuntimeErr(error RuntimeError, handler ErrorHandler) {
	strLine := strconv.Itoa(error.Token.Line)
	handler.Log.Println(utils.Red + "(:" + strLine + ") Runtime error ->" + utils.White + " " + error.Message + utils.Reset)
	handler.RuntimeError = true
}

func report(line int, where, message string, handler ErrorHandler) {
	strLine := strconv.Itoa(line)
	handler.Log.Println(utils.Red + "(:" + strLine + ") Error at " + where + " ->" + utils.White + " " + message + utils.Reset)
	handler.Error = true
}
