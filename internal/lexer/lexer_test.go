package lexer

import (
	"testing"

	"ativ10/internal/fun"
)

func tokensFrom(input string) ([]fun.Token, error) {
	l := NewLexer(input)
	var out []fun.Token
	for {
		tok, err := l.NextToken()
		if err != nil {
			return nil, err
		}
		if tok.Type == fun.TokenEOF {
			return out, nil
		}
		out = append(out, tok)
	}
}

func TestKeywordAndCommaTokens(t *testing.T) {
	input := "fun soma(x, y) { var z = x + y; return z; } main { return soma(1, 2); }"
	toks, err := tokensFrom(input)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	wantTypes := []fun.TokenType{
		fun.TokenFun, fun.TokenIdent, fun.TokenParenEsq, fun.TokenIdent, fun.TokenVirgula, fun.TokenIdent, fun.TokenParenDir,
		fun.TokenChaveEsq, fun.TokenVar, fun.TokenIdent, fun.TokenIgual, fun.TokenIdent, fun.TokenSoma, fun.TokenIdent,
		fun.TokenPontoVirgula, fun.TokenReturn, fun.TokenIdent, fun.TokenPontoVirgula, fun.TokenChaveDir,
		fun.TokenMain, fun.TokenChaveEsq, fun.TokenReturn, fun.TokenIdent, fun.TokenParenEsq, fun.TokenNumero,
		fun.TokenVirgula, fun.TokenNumero, fun.TokenParenDir, fun.TokenPontoVirgula, fun.TokenChaveDir,
	}

	if len(toks) != len(wantTypes) {
		t.Fatalf("token count: got %d want %d", len(toks), len(wantTypes))
	}
	for i, want := range wantTypes {
		if toks[i].Type != want {
			t.Fatalf("token %d: got %v want %v", i, toks[i].Type, want)
		}
	}
}

func TestKeywordTokens(t *testing.T) {
	cases := []struct {
		input string
		typ   fun.TokenType
	}{
		{"fun", fun.TokenFun},
		{"var", fun.TokenVar},
		{"main", fun.TokenMain},
		{"if", fun.TokenIf},
		{"else", fun.TokenElse},
		{"while", fun.TokenWhile},
		{"return", fun.TokenReturn},
		{"true", fun.TokenTrue},
		{"false", fun.TokenFalse},
		{"and", fun.TokenAnd},
		{"or", fun.TokenOr},
		{"not", fun.TokenNot},
	}
	for _, c := range cases {
		toks, err := tokensFrom(c.input)
		if err != nil {
			t.Fatalf("input %q: unexpected err: %v", c.input, err)
		}
		if len(toks) != 1 {
			t.Fatalf("input %q: token count: got %d want 1", c.input, len(toks))
		}
		if toks[0].Type != c.typ {
			t.Fatalf("input %q: got type %v want %v", c.input, toks[0].Type, c.typ)
		}
	}
}

func TestVirgulaToken(t *testing.T) {
	toks, err := tokensFrom(",")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(toks) != 1 {
		t.Fatalf("token count: got %d want 1", len(toks))
	}
	want := fun.Token{Type: fun.TokenVirgula, Lexeme: ",", Pos: 0}
	if toks[0] != want {
		t.Fatalf("got %+v want %+v", toks[0], want)
	}
}

func TestCompositeOperatorTokens(t *testing.T) {
	input := "% <= >= != += -= *= /= %="
	toks, err := tokensFrom(input)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	wantTypes := []fun.TokenType{
		fun.TokenMod,
		fun.TokenMenorIgual,
		fun.TokenMaiorIgual,
		fun.TokenDiferente,
		fun.TokenSomaIgual,
		fun.TokenSubIgual,
		fun.TokenMultIgual,
		fun.TokenDivIgual,
		fun.TokenModIgual,
	}

	if len(toks) != len(wantTypes) {
		t.Fatalf("token count: got %d want %d", len(toks), len(wantTypes))
	}
	for i, want := range wantTypes {
		if toks[i].Type != want {
			t.Fatalf("token %d: got %v want %v", i, toks[i].Type, want)
		}
	}
}

func TestLexicalErrorDigitLetraSequence(t *testing.T) {
	_, err := tokensFrom("237axy")
	if err == nil {
		t.Fatalf("expected lexical error")
	}
	if err.Error() != "Erro lexico na posicao 3" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestWhitespace(t *testing.T) {
	input := "fun f(\n x,\t2y)"
	_, err := tokensFrom(input)
	if err == nil {
		t.Fatalf("expected lexical error")
	}
	if err.Error() != "Erro lexico na posicao 12" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}
