package ast

type Expression interface {
	Accept(visitor Visitor) interface{}
}

type Visitor interface {
	VisitBinaryExpression(expression Binary) interface{}
	VisitGroupingExpression(expression Grouping) interface{}
	VisitLiteralExpression(expression Literal) interface{}
	VisitUnaryExpression(expression Unary) interface{}
	VisitVariableExpression(expression Variable) interface{}
	VisitAssignExpression(expression Assign) interface{}
	VisitLogicalExpression(expression Logical) interface{}
	VisitCallExpression(expression Call) interface{}
}

type Binary struct {
	Left     Expression
	Operator Token
	Right    Expression
}

func (b Binary) Accept(visitor Visitor) interface{} {
	return visitor.VisitBinaryExpression(b)
}

type Grouping struct {
	Expression Expression
}

func (g Grouping) Accept(visitor Visitor) interface{} {
	return visitor.VisitGroupingExpression(g)
}

type Literal struct {
	Value any
}

func (l Literal) Accept(visitor Visitor) interface{} {
	return visitor.VisitLiteralExpression(l)
}

type Unary struct {
	Operator Token
	Right    Expression
}

func (u Unary) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnaryExpression(u)
}

type Variable struct {
	Name Token
}

func (v Variable) Accept(visitor Visitor) interface{} {
	return visitor.VisitVariableExpression(v)
}

type Assign struct {
	Name  Token
	Value Expression
}

func (a Assign) Accept(visitor Visitor) interface{} {
	return visitor.VisitAssignExpression(a)
}

type Logical struct {
	Left     Expression
	Operator Token
	Right    Expression
}

func (l Logical) Accept(visitor Visitor) interface{} {
	return visitor.VisitLogicalExpression(l)
}

type Call struct {
	Callee    Expression
	Paren     Token
	Arguments []Expression
}

func (c Call) Accept(visitor Visitor) interface{} {
	return visitor.VisitCallExpression(c)
}
