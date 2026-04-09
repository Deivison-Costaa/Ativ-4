package interpreter

import (
	"testing"

	"ativ10/internal/lexer"
	"ativ10/internal/parser"
	"ativ10/internal/semantic"
)

func evalProgram(t *testing.T, input string) int {
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
	value, err := EvalPrograma(prog)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	return value
}

func evalProgramErr(t *testing.T, input string) error {
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
	_, err = EvalPrograma(prog)
	return err
}

func TestEvalProgramFunctionWithLocalVar(t *testing.T) {
	input := `fun abs(x) {
  var y = 0;
  if x < 0 {
    y = 0 - x;
  } else {
    y = x;
  }
  return y;
}
main {
  return abs(8) + abs(0 - 3);
}`
	if got := evalProgram(t, input); got != 11 {
		t.Fatalf("eval: got %d want 11", got)
	}
}

func TestEvalProgramFunctionCallingFunction(t *testing.T) {
	input := `fun dup(x) {
  return x + x;
}
fun quad(x) {
  return dup(dup(x));
}
main {
  return quad(7);
}`
	if got := evalProgram(t, input); got != 28 {
		t.Fatalf("eval: got %d want 28", got)
	}
}

func TestEvalProgramFunctionNoParams(t *testing.T) {
	input := `fun seed() {
  return 21;
}
main {
  return seed() + seed();
}`
	if got := evalProgram(t, input); got != 42 {
		t.Fatalf("eval: got %d want 42", got)
	}
}

func TestEvalProgramShadowGlobalByParamAndLocal(t *testing.T) {
	input := `var x = 100;
fun f(x) {
  var y = x + 1;
  return y;
}
main {
  return f(41) + x;
}`
	if got := evalProgram(t, input); got != 142 {
		t.Fatalf("eval: got %d want 142", got)
	}
}

func TestEvalProgramRecursiveFactorial(t *testing.T) {
	input := `fun fact(n) {
  if n < 2 {
    n = 1;
  } else {
    n = n * fact(n - 1);
  }
  return n;
}
main {
  return fact(5);
}`
	if got := evalProgram(t, input); got != 120 {
		t.Fatalf("eval: got %d want 120", got)
	}
}

func TestEvalBooleanOperators(t *testing.T) {
	input := `main {
  return true and not false or false;
}`
	if got := evalProgram(t, input); got != 1 {
		t.Fatalf("eval: got %d want 1", got)
	}
}

func TestEvalBooleanShortCircuit(t *testing.T) {
	input := `var x = 0;
fun touch() {
  x = x + 1;
  return true;
}
main {
  if true or touch() {
  } else {
  }
  if false and touch() {
  } else {
  }
  return x;
}`
	if got := evalProgram(t, input); got != 0 {
		t.Fatalf("eval: got %d want 0", got)
	}
}

func TestEvalModulo(t *testing.T) {
	input := `main {
  return 10 % 3;
}`
	if got := evalProgram(t, input); got != 1 {
		t.Fatalf("eval: got %d want 1", got)
	}
}

func TestEvalNewComparisons(t *testing.T) {
	input := `main {
  return (3 <= 3) and (5 >= 4) and (7 != 2);
}`
	if got := evalProgram(t, input); got != 1 {
		t.Fatalf("eval: got %d want 1", got)
	}
}

func TestEvalCompoundAssignments(t *testing.T) {
	input := `var x = 19;
main {
  x += 5;
  x -= 2;
  x *= 3;
  x /= 4;
  x %= 5;
  return x;
}`
	if got := evalProgram(t, input); got != 1 {
		t.Fatalf("eval: got %d want 1", got)
	}
}

func TestEvalModuloByZeroError(t *testing.T) {
	err := evalProgramErr(t, `main { return 10 % 0; }`)
	if err == nil {
		t.Fatalf("expected execution error")
	}
	if err.Error() != "Erro de execucao: divisao por zero" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestEvalCompoundAssignmentOnLocal(t *testing.T) {
	input := `fun f(x) {
  var y = x;
  y += 5;
  y %= 4;
  return y;
}
main {
  return f(5);
}`
	if got := evalProgram(t, input); got != 2 {
		t.Fatalf("eval: got %d want 2", got)
	}
}

func TestComparisonValuesRemainNormalized(t *testing.T) {
	input := `main {
  return (1 != 1) + (2 <= 3) + (4 >= 2);
}`
	if got := evalProgram(t, input); got != 2 {
		t.Fatalf("eval: got %d want 2", got)
	}
}
