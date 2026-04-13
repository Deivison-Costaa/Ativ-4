package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ativ10/internal/codegen"
	"ativ10/internal/fun"
	"ativ10/internal/interpreter"
	"ativ10/internal/lexer"
	"ativ10/internal/parser"
	"ativ10/internal/semantic"
)

func parseFixture(t *testing.T, parts ...string) (string, *fun.Programa) {
	t.Helper()

	pathParts := append([]string{"..", "..", "testdata"}, parts...)
	path := filepath.Join(pathParts...)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read fixture %q: %v", path, err)
	}

	p, err := parser.NewParser(lexer.NewLexer(string(data)))
	if err != nil {
		t.Fatalf("parse setup %q: %v", path, err)
	}
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse %q: %v", path, err)
	}

	return path, prog
}

func TestPositiveFixturesInterpret(t *testing.T) {
	cases := map[string]int{
		"abs.fun":           11,
		"bool_func.fun":     1,
		"bool_if.fun":       1,
		"bool_short_or.fun": 0,
		"bool_while.fun":    3,
		"chain.fun":         28,
		"compare.fun":       1,
		"compound.fun":      1,
		"combined.fun":      1,
		"fact.fun":          120,
		"mod.fun":           1,
		"noargs.fun":        42,
		"shadow.fun":        142,
	}

	for name, want := range cases {
		t.Run(name, func(t *testing.T) {
			path, prog := parseFixture(t, "fun", "positive", name)
			if err := semantic.CheckProgram(prog); err != nil {
				t.Fatalf("semantic %q: %v", path, err)
			}
			got, err := interpreter.EvalPrograma(prog)
			if err != nil {
				t.Fatalf("interpret %q: %v", path, err)
			}
			if got != want {
				t.Fatalf("fixture %q: got %d want %d", path, got, want)
			}
		})
	}
}

func TestPositiveFixturesGenerateAssembly(t *testing.T) {
	fixtures := []string{
		"abs.fun",
		"bool_func.fun",
		"bool_if.fun",
		"bool_short_or.fun",
		"bool_while.fun",
		"chain.fun",
		"compare.fun",
		"compound.fun",
		"combined.fun",
		"fact.fun",
		"mod.fun",
		"noargs.fun",
		"shadow.fun",
	}

	for _, name := range fixtures {
		t.Run(name, func(t *testing.T) {
			path, prog := parseFixture(t, "fun", "positive", name)
			if err := semantic.CheckProgram(prog); err != nil {
				t.Fatalf("semantic %q: %v", path, err)
			}
			code := codegen.NewCodeGenerator().Generate(prog)
			if !strings.Contains(code, ".globl _start") {
				t.Fatalf("generated code for %q missing entrypoint:\n%s", path, code)
			}
			if !strings.Contains(code, ".include \"runtime/runtime.s\"") {
				t.Fatalf("generated code for %q missing runtime include:\n%s", path, code)
			}
		})
	}
}

func TestNegativeFixturesFailSemantic(t *testing.T) {
	cases := map[string]string{
		"err_fun_arity.fun": "Erro semantico: funcao 'inc' esperava 1 parametros, recebeu 2",
		"err_fun_undef.fun": "Erro semantico: funcao 'foo' nao declarada",
	}

	for name, want := range cases {
		t.Run(name, func(t *testing.T) {
			path, prog := parseFixture(t, "fun", "negative", name)
			err := semantic.CheckProgram(prog)
			if err == nil {
				t.Fatalf("expected semantic error for %q", path)
			}
			if !strings.Contains(err.Error(), want) {
				t.Fatalf("semantic error for %q: got %q want substring %q", path, err.Error(), want)
			}
		})
	}
}
