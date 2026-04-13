package semantic

import (
	"testing"

	"ativ10/internal/lexer"
	"ativ10/internal/parser"
)

func checkProgram(t *testing.T, input string) error {
	t.Helper()

	p, err := parser.NewParser(lexer.NewLexer(input))
	if err != nil {
		t.Fatalf("parse setup unexpected err: %v", err)
	}
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse unexpected err: %v", err)
	}

	return CheckProgram(prog)
}

func TestSemanticErrorUndeclaredFunction(t *testing.T) {
	err := checkProgram(t, `main { return foo(); }`)
	if err == nil {
		t.Fatalf("expected semantic error")
	}
	if err.Error() != "Erro semantico: funcao 'foo' nao declarada" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestSemanticErrorWrongArity(t *testing.T) {
	err := checkProgram(t, `fun inc(x) { return x + 1; } main { return inc(1, 2); }`)
	if err == nil {
		t.Fatalf("expected semantic error")
	}
	if err.Error() != "Erro semantico: funcao 'inc' esperava 1 parametros, recebeu 2" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestSemanticErrorVariableUsedAsFunction(t *testing.T) {
	err := checkProgram(t, `var foo = 1; main { return foo(); }`)
	if err == nil {
		t.Fatalf("expected semantic error")
	}
	if err.Error() != "Erro semantico: funcao 'foo' nao declarada" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestSemanticErrorDuplicateLocal(t *testing.T) {
	err := checkProgram(t, `fun f(x) { var x = 1; return x; } main { return f(3); }`)
	if err == nil {
		t.Fatalf("expected semantic error")
	}
	if err.Error() != "Erro semantico: simbolo 'x' ja declarado" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestSemanticErrorCallToLaterFunctionRejected(t *testing.T) {
	err := checkProgram(t, `fun a() { return b(); } fun b() { return 1; } main { return a(); }`)
	if err == nil {
		t.Fatalf("expected semantic error")
	}
	if err.Error() != "Erro semantico: funcao 'b' nao declarada" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}
