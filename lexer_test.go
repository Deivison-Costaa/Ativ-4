package main

import "testing"

func tokensFrom(input string) ([]Token, error) {
	l := NewLexer(input)
	var out []Token
	for {
		tok, err := l.NextToken()
		if err != nil {
			return nil, err
		}
		if tok.Type == TokenEOF {
			return out, nil
		}
		out = append(out, tok)
	}
}

func TestExampleTokens(t *testing.T) {
	input := "(33 + (912 * 11))"
	toks, err := tokensFrom(input)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	exp := []Token{
		{Type: TokenParenEsq, Lexeme: "(", Pos: 0},
		{Type: TokenNumero, Lexeme: "33", Pos: 1},
		{Type: TokenSoma, Lexeme: "+", Pos: 4},
		{Type: TokenParenEsq, Lexeme: "(", Pos: 6},
		{Type: TokenNumero, Lexeme: "912", Pos: 7},
		{Type: TokenMult, Lexeme: "*", Pos: 11},
		{Type: TokenNumero, Lexeme: "11", Pos: 13},
		{Type: TokenParenDir, Lexeme: ")", Pos: 15},
		{Type: TokenParenDir, Lexeme: ")", Pos: 16},
	}

	if len(toks) != len(exp) {
		t.Fatalf("token count: got %d want %d", len(toks), len(exp))
	}
	for i := range exp {
		if toks[i] != exp[i] {
			t.Fatalf("token %d: got %+v want %+v", i, toks[i], exp[i])
		}
	}
}

func TestWhitespace(t *testing.T) {
	input := "(  3\t+\n(4+5) )"
	toks, err := tokensFrom(input)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	exp := []Token{
		{Type: TokenParenEsq, Lexeme: "(", Pos: 0},
		{Type: TokenNumero, Lexeme: "3", Pos: 3},
		{Type: TokenSoma, Lexeme: "+", Pos: 5},
		{Type: TokenParenEsq, Lexeme: "(", Pos: 7},
		{Type: TokenNumero, Lexeme: "4", Pos: 8},
		{Type: TokenSoma, Lexeme: "+", Pos: 9},
		{Type: TokenNumero, Lexeme: "5", Pos: 10},
		{Type: TokenParenDir, Lexeme: ")", Pos: 11},
		{Type: TokenParenDir, Lexeme: ")", Pos: 13},
	}

	if len(toks) != len(exp) {
		t.Fatalf("token count: got %d want %d", len(toks), len(exp))
	}
	for i := range exp {
		if toks[i] != exp[i] {
			t.Fatalf("token %d: got %+v want %+v", i, toks[i], exp[i])
		}
	}
}

func TestLexicalErrorPosition(t *testing.T) {
	_, err := tokensFrom("(3 + @)")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "Erro lexico na posicao 5" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestEmptyInput(t *testing.T) {
	toks, err := tokensFrom("")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(toks) != 0 {
		t.Fatalf("token count: got %d want %d", len(toks), 0)
	}
}

func TestInvalidCharacterAtStart(t *testing.T) {
	_, err := tokensFrom("@")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "Erro lexico na posicao 0" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestNumberTokenPosition(t *testing.T) {
	toks, err := tokensFrom("  123")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(toks) != 1 {
		t.Fatalf("token count: got %d want %d", len(toks), 1)
	}
	if toks[0].Type != TokenNumero || toks[0].Lexeme != "123" || toks[0].Pos != 2 {
		t.Fatalf("token: got %+v want %+v", toks[0], Token{Type: TokenNumero, Lexeme: "123", Pos: 2})
	}
}

func TestIdentToken(t *testing.T) {
	cases := []struct {
		input  string
		lexeme string
		pos    int
	}{
		{"abc", "abc", 0},
		{"x", "x", 0},
		{"Foo123", "Foo123", 0},
	}
	for _, c := range cases {
		toks, err := tokensFrom(c.input)
		if err != nil {
			t.Fatalf("input %q: unexpected err: %v", c.input, err)
		}
		if len(toks) != 1 {
			t.Fatalf("input %q: token count: got %d want 1", c.input, len(toks))
		}
		want := Token{Type: TokenIdent, Lexeme: c.lexeme, Pos: c.pos}
		if toks[0] != want {
			t.Fatalf("input %q: got %+v want %+v", c.input, toks[0], want)
		}
	}
}

func TestIgualToken(t *testing.T) {
	toks, err := tokensFrom("=")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(toks) != 1 {
		t.Fatalf("token count: got %d want 1", len(toks))
	}
	want := Token{Type: TokenIgual, Lexeme: "=", Pos: 0}
	if toks[0] != want {
		t.Fatalf("got %+v want %+v", toks[0], want)
	}
}

func TestIgualIgualToken(t *testing.T) {
	toks, err := tokensFrom("==")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(toks) != 1 {
		t.Fatalf("token count: got %d want 1", len(toks))
	}
	want := Token{Type: TokenIgualIgual, Lexeme: "==", Pos: 0}
	if toks[0] != want {
		t.Fatalf("got %+v want %+v", toks[0], want)
	}
}

func TestPontoVirgulaToken(t *testing.T) {
	toks, err := tokensFrom(";")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(toks) != 1 {
		t.Fatalf("token count: got %d want 1", len(toks))
	}
	want := Token{Type: TokenPontoVirgula, Lexeme: ";", Pos: 0}
	if toks[0] != want {
		t.Fatalf("got %+v want %+v", toks[0], want)
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

func TestChaveTokens(t *testing.T) {
	toks, err := tokensFrom("{}")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(toks) != 2 {
		t.Fatalf("token count: got %d want 2", len(toks))
	}
	if toks[0].Type != TokenChaveEsq || toks[1].Type != TokenChaveDir {
		t.Fatalf("got %+v %+v, want ChaveEsq ChaveDir", toks[0], toks[1])
	}
}

func TestComparadorTokens(t *testing.T) {
	cases := []struct {
		input string
		typ   TokenType
		lex   string
	}{
		{"<", TokenMenor, "<"},
		{">", TokenMaior, ">"},
		{"==", TokenIgualIgual, "=="},
	}
	for _, c := range cases {
		toks, err := tokensFrom(c.input)
		if err != nil {
			t.Fatalf("input %q: unexpected err: %v", c.input, err)
		}
		if len(toks) != 1 {
			t.Fatalf("input %q: token count: got %d want 1", c.input, len(toks))
		}
		want := Token{Type: c.typ, Lexeme: c.lex, Pos: 0}
		if toks[0] != want {
			t.Fatalf("input %q: got %+v want %+v", c.input, toks[0], want)
		}
	}
}

func TestKeywordTokens(t *testing.T) {
	cases := []struct {
		input string
		typ   TokenType
	}{
		{"if", TokenIf},
		{"else", TokenElse},
		{"while", TokenWhile},
		{"return", TokenReturn},
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
