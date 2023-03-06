package scanner

import (
	"jota/ast"
	"jota/errors"
	"strconv"
)

var (
	keywords = map[string]ast.Type{
		"class":    ast.CLASS,
		"else":     ast.ELSE,
		"false":    ast.FALSE,
		"for":      ast.FOR,
		"function": ast.FUNCTION,
		"if":       ast.IF,
		"nil":      ast.NIL,
		"print":    ast.PRINT,
		"return":   ast.RETURN,
		"super":    ast.SUPER,
		"this":     ast.THIS,
		"true":     ast.TRUE,
		"assign":   ast.VARIABLE,
		"while":    ast.WHILE,
	}
)

type Scanner struct {
	source string
	tokens []ast.Token

	errorHandler errors.ErrorHandler

	start, current, line int
}

func CreateScanner(source string, errorHandler errors.ErrorHandler) *Scanner {
	return &Scanner{source: source, start: 0, current: 0, line: 1, errorHandler: errorHandler}
}

func (s *Scanner) ScanTokens() []ast.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, ast.Token{Type: ast.EOF, Lexeme: "", Literal: nil, Line: s.line})
	return s.tokens
}

func (s *Scanner) scanToken() {
	char := s.advance()

	switch char {
	case '(':
		s.addToken(ast.LEFT_BRACKET)
	case ')':
		s.addToken(ast.RIGHT_BRACKET)
	case '{':
		s.addToken(ast.LEFT_BRACE)
	case '}':
		s.addToken(ast.RIGHT_BRACE)
	case ',':
		s.addToken(ast.COMMA)
	case '.':
		s.addToken(ast.DOT)
	case '-':
		if s.match('-') {
			s.addToken(ast.DECREMENT)
		} else {
			s.addToken(ast.MINUS)
		}
	case '+':
		if s.match('+') {
			s.addToken(ast.INCREMENT)
		} else {
			s.addToken(ast.PLUS)
		}
	case '^':
		s.addToken(ast.CARET)
	case ';':
		s.addToken(ast.SEMICOLON)
	case '*':
		s.addToken(ast.ASTERISK)
	case '&':
		if s.match('&') {
			s.addToken(ast.AND)
		} else {
			errors.ErrWithoutToken(s.line, "unexpected character found", s.errorHandler)
		}
	case '|':
		if s.match('|') {
			s.addToken(ast.OR)
		} else {
			errors.ErrWithoutToken(s.line, "unexpected character found", s.errorHandler)
		}
	case '!':
		if s.match('=') {
			s.addToken(ast.BANG_EQUAL)
		} else {
			s.addToken(ast.BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(ast.EQUAL_EQUAL)
		} else {
			s.addToken(ast.EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(ast.LESS_EQUAL)
		} else {
			s.addToken(ast.LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(ast.GREATER_EQUAL)
		} else {
			s.addToken(ast.GREATER)
		}
	case '/':
		s.addToken(ast.SLASH)
	case '%':
		s.addToken(ast.PERCENT)
	case '#':
		for s.peek() != '\n' && !s.isAtEnd() {
			s.advance()
		}
	case ' ', '\r', '\t':
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if isDigit(char) {
			s.number()
		} else if s.isAlpha(char) {
			s.identifier()
		} else {
			errors.ErrWithoutToken(s.line, "unexpected character found", s.errorHandler)
		}
	}
}

func (s *Scanner) identifier() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	token, found := keywords[text]
	if !found {
		token = ast.IDENTIFIER
	}
	s.addToken(token)
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	num, err := strconv.ParseFloat(s.source[s.start:s.current], 64)

	if err != nil {
		errors.ErrWithoutToken(s.line, "scanner has an issue parsing a number", s.errorHandler)
	}

	s.addTokenWithLiteral(ast.NUMBER, num)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		errors.ErrWithoutToken(s.line, "unterminated string", s.errorHandler)
		return
	}

	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addTokenWithLiteral(ast.STRING, value)
}

func (s *Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) isAlpha(char byte) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char == '_')
}

func (s *Scanner) isAlphaNumeric(char byte) bool {
	return s.isAlpha(char) || isDigit(char)
}

func isDigit(char byte) bool {
	return char >= '0' && char <= '9'
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) addToken(tokenType ast.Type) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) addTokenWithLiteral(tokenType ast.Type, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, ast.Token{Type: tokenType, Lexeme: text, Literal: literal, Line: s.line})
}
