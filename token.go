package main

type TokenType string

const (
	TokenEOF          TokenType = "EOF"
	TokenNumero       TokenType = "Numero"
	TokenParenEsq     TokenType = "ParenEsq"
	TokenParenDir     TokenType = "ParenDir"
	TokenChaveEsq     TokenType = "ChaveEsq"
	TokenChaveDir     TokenType = "ChaveDir"
	TokenVirgula      TokenType = "Virgula"
	TokenSoma         TokenType = "Soma"
	TokenSub          TokenType = "Sub"
	TokenMult         TokenType = "Mult"
	TokenDiv          TokenType = "Div"
	TokenMenor        TokenType = "Menor"
	TokenMaior        TokenType = "Maior"
	TokenIgual        TokenType = "Igual"
	TokenIgualIgual   TokenType = "IgualIgual"
	TokenPontoVirgula TokenType = "PontoVirgula"
	TokenIdent        TokenType = "Ident"
	TokenFun          TokenType = "Fun"
	TokenVar          TokenType = "Var"
	TokenMain         TokenType = "Main"
	TokenIf           TokenType = "If"
	TokenElse         TokenType = "Else"
	TokenWhile        TokenType = "While"
	TokenReturn       TokenType = "Return"
)

type Token struct {
	Type   TokenType
	Lexeme string
	Pos    int
}
