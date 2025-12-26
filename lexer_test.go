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
		if tok.Type == "" {
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
	_, err := tokensFrom("(3 + a)")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "Erro lexico na posicao 5" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}
