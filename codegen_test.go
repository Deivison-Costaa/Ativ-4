package main

import (
	"strings"
	"testing"
)

func TestCodegenFunctionCallPushesArgsInReverseOrder(t *testing.T) {
	prog, err := parseFun(`fun soma(x, y) { return x + y; } main { return soma(11, 202); }`)
	if err != nil {
		t.Fatalf("parse unexpected err: %v", err)
	}
	if err := CheckProgram(prog); err != nil {
		t.Fatalf("semantic err: %v", err)
	}
	code := NewCodeGenerator().Generate(prog)

	idx202 := strings.Index(code, "  mov $202, %rax\n  push %rax\n  mov $11, %rax\n  push %rax\n  call soma")
	if idx202 < 0 {
		t.Fatalf("reverse-order push sequence not found in code:\n%s", code)
	}
	if !strings.Contains(code, "  add $16, %rsp") {
		t.Fatalf("stack cleanup after call not found:\n%s", code)
	}
}

func TestCodegenFunctionPrologueAndOffsets(t *testing.T) {
	prog, err := parseFun(`fun square(x) { var y = x * x; return y; } main { return square(7); }`)
	if err != nil {
		t.Fatalf("parse unexpected err: %v", err)
	}
	if err := CheckProgram(prog); err != nil {
		t.Fatalf("semantic err: %v", err)
	}
	code := NewCodeGenerator().Generate(prog)

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
	prog, err := parseFun(`var x = 10; fun inc(n) { return n + x; } main { return inc(5); }`)
	if err != nil {
		t.Fatalf("parse unexpected err: %v", err)
	}
	if err := CheckProgram(prog); err != nil {
		t.Fatalf("semantic err: %v", err)
	}
	code := NewCodeGenerator().Generate(prog)

	if !strings.Contains(code, "  mov %rax, x(%rip)") {
		t.Fatalf("global store not using RIP-relative addressing:\n%s", code)
	}
	if !strings.Contains(code, "  mov x(%rip), %rax") {
		t.Fatalf("global load not using RIP-relative addressing:\n%s", code)
	}
}
