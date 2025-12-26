package main

type TokenType string

const (
	TokenNumero   TokenType = "Numero"
	TokenParenEsq TokenType = "ParenEsq"
	TokenParenDir TokenType = "ParenDir"
	TokenSoma     TokenType = "Soma"
	TokenSub      TokenType = "Sub"
	TokenMult     TokenType = "Mult"
	TokenDiv      TokenType = "Div"
)

type Token struct {
	Type   TokenType
	Lexeme string
	Pos    int
}
