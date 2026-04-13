package lexer

import (
	"fmt"

	"ativ10/internal/fun"
)

var keywords = map[string]fun.TokenType{
	"fun":    fun.TokenFun,
	"var":    fun.TokenVar,
	"main":   fun.TokenMain,
	"if":     fun.TokenIf,
	"else":   fun.TokenElse,
	"while":  fun.TokenWhile,
	"return": fun.TokenReturn,
	"true":   fun.TokenTrue,
	"false":  fun.TokenFalse,
	"and":    fun.TokenAnd,
	"or":     fun.TokenOr,
	"not":    fun.TokenNot,
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

func (l *Lexer) NextToken() (fun.Token, error) {
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
			return fun.Token{Type: fun.TokenParenEsq, Lexeme: "(", Pos: start}, nil
		case ')':
			l.pos++
			return fun.Token{Type: fun.TokenParenDir, Lexeme: ")", Pos: start}, nil
		case '{':
			l.pos++
			return fun.Token{Type: fun.TokenChaveEsq, Lexeme: "{", Pos: start}, nil
		case '}':
			l.pos++
			return fun.Token{Type: fun.TokenChaveDir, Lexeme: "}", Pos: start}, nil
		case ',':
			l.pos++
			return fun.Token{Type: fun.TokenVirgula, Lexeme: ",", Pos: start}, nil
		case '+':
			if l.peekNext() == '=' {
				l.pos += 2
				return fun.Token{Type: fun.TokenSomaIgual, Lexeme: "+=", Pos: start}, nil
			}
			l.pos++
			return fun.Token{Type: fun.TokenSoma, Lexeme: "+", Pos: start}, nil
		case '-':
			if l.peekNext() == '=' {
				l.pos += 2
				return fun.Token{Type: fun.TokenSubIgual, Lexeme: "-=", Pos: start}, nil
			}
			l.pos++
			return fun.Token{Type: fun.TokenSub, Lexeme: "-", Pos: start}, nil
		case '*':
			if l.peekNext() == '=' {
				l.pos += 2
				return fun.Token{Type: fun.TokenMultIgual, Lexeme: "*=", Pos: start}, nil
			}
			l.pos++
			return fun.Token{Type: fun.TokenMult, Lexeme: "*", Pos: start}, nil
		case '/':
			if l.peekNext() == '=' {
				l.pos += 2
				return fun.Token{Type: fun.TokenDivIgual, Lexeme: "/=", Pos: start}, nil
			}
			l.pos++
			return fun.Token{Type: fun.TokenDiv, Lexeme: "/", Pos: start}, nil
		case '%':
			if l.peekNext() == '=' {
				l.pos += 2
				return fun.Token{Type: fun.TokenModIgual, Lexeme: "%=", Pos: start}, nil
			}
			l.pos++
			return fun.Token{Type: fun.TokenMod, Lexeme: "%", Pos: start}, nil
		case '<':
			if l.peekNext() == '=' {
				l.pos += 2
				return fun.Token{Type: fun.TokenMenorIgual, Lexeme: "<=", Pos: start}, nil
			}
			l.pos++
			return fun.Token{Type: fun.TokenMenor, Lexeme: "<", Pos: start}, nil
		case '>':
			if l.peekNext() == '=' {
				l.pos += 2
				return fun.Token{Type: fun.TokenMaiorIgual, Lexeme: ">=", Pos: start}, nil
			}
			l.pos++
			return fun.Token{Type: fun.TokenMaior, Lexeme: ">", Pos: start}, nil
		case '!':
			if l.peekNext() == '=' {
				l.pos += 2
				return fun.Token{Type: fun.TokenDiferente, Lexeme: "!=", Pos: start}, nil
			}
			return fun.Token{}, fmt.Errorf("Erro lexico na posicao %d", start)
		case '=':
			if l.peekNext() == '=' {
				l.pos += 2
				return fun.Token{Type: fun.TokenIgualIgual, Lexeme: "==", Pos: start}, nil
			}
			l.pos++
			return fun.Token{Type: fun.TokenIgual, Lexeme: "=", Pos: start}, nil
		case ';':
			l.pos++
			return fun.Token{Type: fun.TokenPontoVirgula, Lexeme: ";", Pos: start}, nil
		default:
			if isDigit(c) {
				for l.pos < len(l.input) && isDigit(l.input[l.pos]) {
					l.pos++
				}
				if l.pos < len(l.input) && isLetter(l.input[l.pos]) {
					return fun.Token{}, fmt.Errorf("Erro lexico na posicao %d", l.pos)
				}
				lex := l.input[start:l.pos]
				return fun.Token{Type: fun.TokenNumero, Lexeme: lex, Pos: start}, nil
			}
			if isLetter(c) {
				for l.pos < len(l.input) && isLetterOrDigit(l.input[l.pos]) {
					l.pos++
				}
				lex := l.input[start:l.pos]
				if kwType, ok := keywords[lex]; ok {
					return fun.Token{Type: kwType, Lexeme: lex, Pos: start}, nil
				}
				return fun.Token{Type: fun.TokenIdent, Lexeme: lex, Pos: start}, nil
			}
			return fun.Token{}, fmt.Errorf("Erro lexico na posicao %d", start)
		}
	}

	return fun.Token{Type: fun.TokenEOF, Lexeme: "", Pos: l.pos}, nil
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
