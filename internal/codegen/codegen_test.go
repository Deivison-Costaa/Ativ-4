package codegen

import (
	"strings"
	"testing"

	"ativ10/internal/lexer"
	"ativ10/internal/parser"
	"ativ10/internal/semantic"
)

func generateCode(t *testing.T, input string) string {
	t.Helper()

	p, err := parser.NewParser(lexer.NewLexer(input))
	if err != nil {
		t.Fatalf("parse setup unexpected err: %v", err)
	}
	prog, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("parse unexpected err: %v", err)
	}
	if err := semantic.CheckProgram(prog); err != nil {
		t.Fatalf("semantic err: %v", err)
	}

	return NewCodeGenerator().Generate(prog)
}

func TestCodegenFunctionCallPushesArgsInReverseOrder(t *testing.T) {
	code := generateCode(t, `fun soma(x, y) { return x + y; } main { return soma(11, 202); }`)

	idx202 := strings.Index(code, "  mov $202, %rax\n  push %rax\n  mov $11, %rax\n  push %rax\n  call soma")
	if idx202 < 0 {
		t.Fatalf("reverse-order push sequence not found in code:\n%s", code)
	}
	if !strings.Contains(code, "  add $16, %rsp") {
		t.Fatalf("stack cleanup after call not found:\n%s", code)
	}
}

func TestCodegenFunctionPrologueAndOffsets(t *testing.T) {
	code := generateCode(t, `fun square(x) { var y = x * x; return y; } main { return square(7); }`)

	wants := []string{
		"square:",
		"  push %rbp",
		"  sub $8, %rsp",
		"  mov %rsp, %rbp",
		"  mov 24(%rbp), %rax",
		"  mov %rax, 0(%rbp)",
		"  add $8, %rsp",
		"  pop %rbp",
		"  ret",
	}
	for _, want := range wants {
		if !strings.Contains(code, want) {
			t.Fatalf("code missing %q:\n%s", want, code)
		}
	}
}

func TestCodegenUsesRipRelativeGlobals(t *testing.T) {
	code := generateCode(t, `var x = 10; fun inc(n) { return n + x; } main { return inc(5); }`)

	if !strings.Contains(code, "  mov %rax, x(%rip)") {
		t.Fatalf("global store not using RIP-relative addressing:\n%s", code)
	}
	if !strings.Contains(code, "  mov x(%rip), %rax") {
		t.Fatalf("global load not using RIP-relative addressing:\n%s", code)
	}
}

func TestCodegenBooleanShortCircuit(t *testing.T) {
	code := generateCode(t, `fun touch() { return true; } main { return false and touch() or true; }`)

	if !strings.Contains(code, "  jz Lfalse") {
		t.Fatalf("missing short-circuit false branch in code:\n%s", code)
	}
	if !strings.Contains(code, "  jnz Ltrue") {
		t.Fatalf("missing short-circuit true branch in code:\n%s", code)
	}
	if !strings.Contains(code, "  call touch") {
		t.Fatalf("expected generated code to include rhs function call:\n%s", code)
	}
}

func TestCodegenModuloUsesRdx(t *testing.T) {
	code := generateCode(t, `main { return 10 % 3; }`)

	if !strings.Contains(code, "  idiv %rbx") {
		t.Fatalf("missing idiv for modulo:\n%s", code)
	}
	if !strings.Contains(code, "  mov %rdx, %rax") {
		t.Fatalf("missing remainder move for modulo:\n%s", code)
	}
}

func TestCodegenNewComparisons(t *testing.T) {
	code := generateCode(t, `main { return (3 <= 4) + (5 >= 1) + (2 != 9); }`)

	for _, want := range []string{"  setle %cl", "  setge %cl", "  setne %cl"} {
		if !strings.Contains(code, want) {
			t.Fatalf("code missing %q:\n%s", want, code)
		}
	}
}
