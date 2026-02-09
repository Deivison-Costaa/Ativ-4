package main

import "testing"

func parse(input string) (Exp, error) {
	p, err := NewParser(NewLexer(input))
	if err != nil {
		return nil, err
	}
	return p.ParseProgram()
}

func TestParseConst(t *testing.T) {
	exp, err := parse("333")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if exp.String() != "333" {
		t.Fatalf("tree: got %q want %q", exp.String(), "333")
	}
	v, err := exp.Eval()
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 333 {
		t.Fatalf("eval: got %d want %d", v, 333)
	}
}

func TestParseTreeAndEval(t *testing.T) {
	input := "(33 + (912 * 11))"
	exp, err := parse(input)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if exp.String() != input {
		t.Fatalf("tree: got %q want %q", exp.String(), input)
	}
	if TreeString(exp) != "+\n|-- 33\n`-- *\n    |-- 912\n    `-- 11" {
		t.Fatalf("visual tree: got %q", TreeString(exp))
	}
	v, err := exp.Eval()
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 10065 {
		t.Fatalf("eval: got %d want %d", v, 10065)
	}
}

func TestParseWhitespace(t *testing.T) {
	input := "(  3\t+\n(4+5) )"
	exp, err := parse(input)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if exp.String() != "(3 + (4 + 5))" {
		t.Fatalf("tree: got %q want %q", exp.String(), "(3 + (4 + 5))")
	}
	v, err := exp.Eval()
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 12 {
		t.Fatalf("eval: got %d want %d", v, 12)
	}
}

func TestSyntaxError(t *testing.T) {
	_, err := parse("(3 + )")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "Erro sintatico na posicao 5" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestSyntaxErrorMissingClosingParen(t *testing.T) {
	_, err := parse("(3 + 4")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "Erro sintatico na posicao 6" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestSyntaxErrorExtraTokensAfterExpression(t *testing.T) {
	_, err := parse("(3 + 4) 5")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "Erro sintatico na posicao 8" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestSyntaxErrorUnknownOperator(t *testing.T) {
	_, err := parse("(3 ^ 4)")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "Erro lexico na posicao 3" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestParseBiggerExpressionEval(t *testing.T) {
	input := "(((1000 + (72 * (55 - 17))) - (840 / (14 - 7))) + ((96 / 8) * (30 + 12)))"
	exp, err := parse(input)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if exp.String() != input {
		t.Fatalf("tree: got %q want %q", exp.String(), input)
	}
	v, err := exp.Eval()
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 4120 {
		t.Fatalf("eval: got %d want %d", v, 4120)
	}
}

func TestEvalDivisionByZero(t *testing.T) {
	exp, err := parse("(1 / 0)")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	_, err = exp.Eval()
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "Erro de execucao: divisao por zero" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}
