package fun

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
	TokenMod          TokenType = "Mod"
	TokenMenor        TokenType = "Menor"
	TokenMaior        TokenType = "Maior"
	TokenMenorIgual   TokenType = "MenorIgual"
	TokenMaiorIgual   TokenType = "MaiorIgual"
	TokenDiferente    TokenType = "Diferente"
	TokenIgual        TokenType = "Igual"
	TokenIgualIgual   TokenType = "IgualIgual"
	TokenSomaIgual    TokenType = "SomaIgual"
	TokenSubIgual     TokenType = "SubIgual"
	TokenMultIgual    TokenType = "MultIgual"
	TokenDivIgual     TokenType = "DivIgual"
	TokenModIgual     TokenType = "ModIgual"
	TokenPontoVirgula TokenType = "PontoVirgula"
	TokenIdent        TokenType = "Ident"
	TokenFun          TokenType = "Fun"
	TokenVar          TokenType = "Var"
	TokenMain         TokenType = "Main"
	TokenIf           TokenType = "If"
	TokenElse         TokenType = "Else"
	TokenWhile        TokenType = "While"
	TokenReturn       TokenType = "Return"
	TokenTrue         TokenType = "True"
	TokenFalse        TokenType = "False"
	TokenAnd          TokenType = "And"
	TokenOr           TokenType = "Or"
	TokenNot          TokenType = "Not"
)

type Token struct {
	Type   TokenType
	Lexeme string
	Pos    int
}
