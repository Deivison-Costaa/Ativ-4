package main

import (
	"fmt"
)

type Lexer struct {
	input string
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0}
}

func (l *Lexer) NextToken() (Token, error) {
	for l.pos < len(l.input) {
		c := l.input[l.pos]
		if isWhitespace(c) {
			l.pos++
			continue
		}

		start := l.pos

		switch c {
		case '(':
			l.pos++
			return Token{Type: TokenParenEsq, Lexeme: "(", Pos: start}, nil
		case ')':
			l.pos++
			return Token{Type: TokenParenDir, Lexeme: ")", Pos: start}, nil
		case '+':
			l.pos++
			return Token{Type: TokenSoma, Lexeme: "+", Pos: start}, nil
		case '-':
			l.pos++
			return Token{Type: TokenSub, Lexeme: "-", Pos: start}, nil
		case '*':
			l.pos++
			return Token{Type: TokenMult, Lexeme: "*", Pos: start}, nil
		case '/':
			l.pos++
			return Token{Type: TokenDiv, Lexeme: "/", Pos: start}, nil
		default:
			if isDigit(c) {
				for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
					l.pos++
				}
				lex := l.input[start:l.pos]
				return Token{Type: TokenNumero, Lexeme: lex, Pos: start}, nil
			}
			return Token{}, fmt.Errorf("Erro lexico na posicao %d", start)
		}
	}

	return Token{}, nil
}

func isWhitespace(c byte) bool {
	switch c {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}
