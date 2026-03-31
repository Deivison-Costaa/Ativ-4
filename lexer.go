package main

import "fmt"

var keywords = map[string]TokenType{
	"if":     TokenIf,
	"else":   TokenElse,
	"while":  TokenWhile,
	"return": TokenReturn,
}

type Lexer struct {
	input string
	pos   int
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input, pos: 0}
}

func (l *Lexer) peekNext() byte {
	if l.pos+1 < len(l.input) {
		return l.input[l.pos+1]
	}
	return 0
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
		case '{':
			l.pos++
			return Token{Type: TokenChaveEsq, Lexeme: "{", Pos: start}, nil
		case '}':
			l.pos++
			return Token{Type: TokenChaveDir, Lexeme: "}", Pos: start}, nil
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
		case '<':
			l.pos++
			return Token{Type: TokenMenor, Lexeme: "<", Pos: start}, nil
		case '>':
			l.pos++
			return Token{Type: TokenMaior, Lexeme: ">", Pos: start}, nil
		case '=':
			if l.peekNext() == '=' {
				l.pos += 2
				return Token{Type: TokenIgualIgual, Lexeme: "==", Pos: start}, nil
			}
			l.pos++
			return Token{Type: TokenIgual, Lexeme: "=", Pos: start}, nil
		case ';':
			l.pos++
			return Token{Type: TokenPontoVirgula, Lexeme: ";", Pos: start}, nil
		default:
			if isDigit(c) {
				for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
					l.pos++
				}
				if l.pos < len(l.input) && isLetter(l.input[l.pos]) {
					return Token{}, fmt.Errorf("Erro lexico na posicao %d", l.pos)
				}
				lex := l.input[start:l.pos]
				return Token{Type: TokenNumero, Lexeme: lex, Pos: start}, nil
			}
			if isLetter(c) {
				for l.pos < len(l.input) && isLetterOrDigit(l.input[l.pos]) {
					l.pos++
				}
				lex := l.input[start:l.pos]
				if kwType, ok := keywords[lex]; ok {
					return Token{Type: kwType, Lexeme: lex, Pos: start}, nil
				}
				return Token{Type: TokenIdent, Lexeme: lex, Pos: start}, nil
			}
			return Token{}, fmt.Errorf("Erro lexico na posicao %d", start)
		}
	}

	return Token{Type: TokenEOF, Lexeme: "", Pos: l.pos}, nil
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

func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isLetterOrDigit(c byte) bool {
	return isLetter(c) || isDigit(c)
}
