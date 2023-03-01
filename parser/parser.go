package parser

import (
	"jota/ast"
	"jota/errors"
)

type Parser struct {
	Tokens       []ast.Token
	current      int
	ErrorHandler errors.ErrorHandler
}

func NewParser(tokens []ast.Token, errorHandler errors.ErrorHandler) *Parser {
	return &Parser{Tokens: tokens, current: 0, ErrorHandler: errorHandler}
}

func (p *Parser) Parse() []ast.Statement {
	var statements []ast.Statement
	for !p.isAtEnd() {
		statement := p.declaration()
		statements = append(statements, statement)
	}
	return statements
}

func (p *Parser) expression() ast.Expression {
	return p.assignment()
}

func (p *Parser) declaration() ast.Statement {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(ParseError); ok {
				p.synchronize()
			}
		}
	}()

	if p.match(ast.FUNCTION) {
		return p.function("function")
	}

	if p.match(ast.VARIABLE) {
		return p.variableDeclaration()
	}
	return p.statement()
}

func (p *Parser) function(kind string) *ast.FunctionStatement {
	name := p.consume(ast.IDENTIFIER, "a "+kind+" is expected")
	p.consume(ast.LEFT_BRACKET, "expected '(' after "+kind+" name")

	var parameters []ast.Token

	if !p.check(ast.RIGHT_BRACKET) {
		for {
			if len(parameters) >= 255 {
				errors.Err(p.peek(), "you are not allowed to have more than 255 parameters", p.ErrorHandler)
			}

			param := p.consume(ast.IDENTIFIER, "expected a parameter name")
			parameters = append(parameters, param)

			if !p.match(ast.COMMA) {
				break
			}
		}

		for p.match(ast.COMMA) {
			if len(parameters) > 255 {
				errors.Err(p.peek(), "you are not allowed to have more than 255 parameters", p.ErrorHandler)
			}
			parameters = append(parameters, p.consume(ast.IDENTIFIER, "expected a parameter name"))
		}
	}
	p.consume(ast.RIGHT_BRACKET, "expected ')' after parameters")

	p.consume(ast.LEFT_BRACE, "expected '{' before "+kind+" body")
	body := p.block()
	return &ast.FunctionStatement{Name: name, Params: parameters, Body: body}
}

func (p *Parser) statement() ast.Statement {
	if p.match(ast.IF) {
		return p.ifStatement()
	}

	if p.match(ast.WHILE) {
		return p.whileStatement()
	}

	if p.match(ast.FOR) {
		return p.forStatement()
	}

	if p.match(ast.PRINT) {
		return p.printStatement()
	}

	if p.match(ast.LEFT_BRACE) {
		return &ast.BlockStatement{Statements: p.block()}
	}

	return p.expressionStatement()
}

func (p *Parser) block() []ast.Statement {
	var statements []ast.Statement
	for !p.check(ast.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	p.consume(ast.RIGHT_BRACE, "expected '}' after a block")
	return statements
}

func (p *Parser) ifStatement() ast.Statement {
	p.consume(ast.LEFT_BRACKET, "expected '(' after an 'if' statement")
	condition := p.expression()
	p.consume(ast.RIGHT_BRACKET, "expected ')' after an 'if' condition")

	thenBranch := p.statement()

	var elseBranch ast.Statement = nil
	if p.match(ast.ELSE) {
		elseBranch = p.statement()
	}

	return &ast.IfStatement{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}
}

func (p *Parser) whileStatement() ast.Statement {
	p.consume(ast.LEFT_BRACKET, "expected '(' after a 'while' statement")
	condition := p.expression()
	p.consume(ast.RIGHT_BRACKET, "expected ')' after a 'while' condition")
	body := p.statement()
	return &ast.WhileStatement{Condition: condition, Body: body}
}

func (p *Parser) forStatement() ast.Statement {
	p.consume(ast.LEFT_BRACKET, "expected '(' after a 'for' statement")

	var initializer ast.Statement
	if p.match(ast.SEMICOLON) {
		initializer = nil
	} else if p.match(ast.VARIABLE) {
		initializer = p.variableDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition ast.Expression = nil
	if !p.check(ast.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(ast.SEMICOLON, "expected ';' after a 'for' loop condition")

	var increment ast.Expression = nil
	if !p.check(ast.RIGHT_BRACKET) {
		increment = p.expression()
	}
	p.consume(ast.RIGHT_BRACKET, "expected ')' after 'for' loop conditions")

	body := p.statement()

	if increment != nil {
		body = &ast.BlockStatement{
			Statements: []ast.Statement{
				body, &ast.ExpressionStatement{Expression: increment},
			},
		}
	}

	if condition == nil {
		condition = &ast.Literal{Value: true}
	}
	body = &ast.WhileStatement{Condition: condition, Body: body}

	if initializer != nil {
		body = &ast.BlockStatement{
			Statements: []ast.Statement{
				initializer, body,
			},
		}
	}

	return body
}

func (p *Parser) printStatement() ast.Statement {
	value := p.expression()
	p.consume(ast.SEMICOLON, "expected ';' after a value")
	return &ast.PrintStatement{Expression: value}
}

func (p *Parser) expressionStatement() ast.Statement {
	expression := p.expression()
	p.consume(ast.SEMICOLON, "expected ';' after an expression")
	return &ast.ExpressionStatement{Expression: expression}
}

func (p *Parser) assignment() ast.Expression {
	expression := p.or()

	if p.match(ast.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if variable, ok := expression.(*ast.Variable); ok {
			return &ast.Assign{Name: variable.Name, Value: value}
		}

		errors.Err(equals, "invalid assignment target", p.ErrorHandler)
	}

	return expression
}

func (p *Parser) or() ast.Expression {
	expression := p.and()

	for p.match(ast.OR) {
		operator := p.previous()
		right := p.and()
		expression = &ast.Logical{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

func (p *Parser) and() ast.Expression {
	expression := p.equality()

	for p.match(ast.AND) {
		operator := p.previous()
		right := p.equality()
		expression = &ast.Logical{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

func (p *Parser) variableDeclaration() ast.Statement {
	name := p.consume(ast.IDENTIFIER, "expected a variable name")

	var initializer ast.Expression
	if p.match(ast.EQUAL) {
		initializer = p.expression()
	}
	p.consume(ast.SEMICOLON, "expected ';' after a variable declaration")
	return &ast.VariableStatement{Name: name, Initializer: initializer}
}

func (p *Parser) equality() ast.Expression {
	expression := p.comparison()

	for p.match(ast.BANG_EQUAL, ast.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expression = &ast.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

func (p *Parser) comparison() ast.Expression {
	expression := p.term()

	for p.match(ast.GREATER, ast.GREATER_EQUAL, ast.LESS, ast.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expression = &ast.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

func (p *Parser) term() ast.Expression {
	expression := p.factor()

	for p.match(ast.MINUS, ast.PLUS, ast.CARET) {
		operator := p.previous()
		right := p.factor()
		expression = &ast.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

func (p *Parser) factor() ast.Expression {
	expression := p.unary()

	for p.match(ast.SLASH, ast.ASTERISK) {
		operator := p.previous()
		right := p.unary()
		expression = &ast.Binary{Left: expression, Operator: operator, Right: right}
	}

	return expression
}

func (p *Parser) unary() ast.Expression {
	if p.match(ast.BANG, ast.MINUS, ast.DECREMENT, ast.INCREMENT) {
		operator := p.previous()
		right := p.unary()
		return &ast.Unary{Operator: operator, Right: right}
	}

	return p.call()
}

func (p *Parser) finishCall(callee ast.Expression) ast.Expression {
	var arguments []ast.Expression

	if !p.check(ast.RIGHT_BRACKET) {
		for {
			if len(arguments) >= 255 {
				errors.Err(p.peek(), "can't have more than 255 arguments", p.ErrorHandler)
			}
			arguments = append(arguments, p.expression())
			if !p.match(ast.COMMA) {
				break
			}
		}
	}

	paren := p.consume(ast.RIGHT_BRACKET, "expected ')' after arguments")
	return &ast.Call{Callee: callee, Paren: paren, Arguments: arguments}
}

func (p *Parser) call() ast.Expression {
	expression := p.primary()

	for {
		if p.match(ast.LEFT_BRACKET) {
			expression = p.finishCall(expression)
		} else {
			break
		}
	}

	return expression
}

func (p *Parser) primary() ast.Expression {
	if p.match(ast.FALSE) {
		return &ast.Literal{Value: false}
	}
	if p.match(ast.TRUE) {
		return &ast.Literal{Value: true}
	}
	if p.match(ast.NIL) {
		return &ast.Literal{Value: nil}
	}

	if p.match(ast.NUMBER, ast.STRING) {
		return &ast.Literal{Value: p.previous().Literal}
	}

	if p.match(ast.IDENTIFIER) {
		return &ast.Variable{Name: p.previous()}
	}

	if p.match(ast.LEFT_BRACKET) {
		expression := p.expression()
		p.consume(ast.RIGHT_BRACKET, "expected ')' after expression")
		return &ast.Grouping{Expression: expression}
	}

	panic(p.throwError(p.peek(), "expected an expression"))
}

func (p *Parser) consume(tokentype ast.Type, message string) ast.Token {
	if p.check(tokentype) {
		return p.advance()
	}

	panic(p.throwError(p.peek(), message))
}

func (p *Parser) throwError(token ast.Token, message string) ParseError {
	errors.Err(token, message, p.ErrorHandler)
	return ParseError{}
}

type ParseError struct{}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == ast.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case ast.CLASS, ast.FUNCTION, ast.VARIABLE, ast.FOR, ast.IF, ast.WHILE, ast.PRINT, ast.RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) match(types ...ast.Type) bool {
	for i := range types {
		if p.check(types[i]) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(tokentype ast.Type) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokentype
}

func (p *Parser) advance() ast.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == ast.EOF
}

func (p *Parser) peek() ast.Token {
	return p.Tokens[p.current]
}

func (p *Parser) previous() ast.Token {
	return p.Tokens[p.current-1]
}
