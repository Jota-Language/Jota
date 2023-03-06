// Package meant only for debugging, out of date (does not implement Visitor anymore)

package ast

// import (
// 	"fmt"
// 	"strings"
// )

// type AstPrinter struct{}

// func (ap AstPrinter) Print(expression Expression) interface{} {
//     return expression.Accept(ap)
// }

// func (ap AstPrinter) VisitBinaryExpression(expression *Binary) interface{} {
// 	return ap.parenthesize(expression.Operator.Lexeme, expression.Left, expression.Right)
// }

// func (ap AstPrinter) VisitGroupingExpression(expression *Grouping) interface{} {
// 	return ap.parenthesize("group", expression.Expression)
// }

// func (ap AstPrinter) VisitLiteralExpression(expression *Literal) interface{} {
// 	if expression.Value == nil {
// 		return "nil"
// 	}
// 	return fmt.Sprint(expression.Value)
// }

// func (ap AstPrinter) VisitUnaryExpression(expression *Unary) interface{} {
// 	return ap.parenthesize(expression.Operator.Lexeme, expression.Right)
// }

// func (ap AstPrinter) parenthesize(name string, expressions ...Expression) interface{} {
//     var builder strings.Builder

//     builder.WriteString("(")
//     builder.WriteString(name)
//     for _, expression := range expressions {
//         builder.WriteString(" ")
//         builder.WriteString(expression.Accept(ap).(string))
//     }
//     builder.WriteString(")")
//     return builder.String()
// }
