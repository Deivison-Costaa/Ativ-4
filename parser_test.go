package main

import "testing"

func parseArith(input string) (Exp, error) {
	p, err := NewParser(NewLexer(input))
	if err != nil {
		return nil, err
	}
	exp, err := p.parseExp()
	if err != nil {
		return nil, err
	}
	if p.cur.Type != TokenEOF {
		return nil, p.syntaxError()
	}
	return exp, nil
}

func parseFun(input string) (*Programa, error) {
	p, err := NewParser(NewLexer(input))
	if err != nil {
		return nil, err
	}
	return p.ParseProgram()
}

func TestParseCallExpression(t *testing.T) {
	exp, err := parseArith("soma(x, 2)")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	call, ok := exp.(*Call)
	if !ok {
		t.Fatalf("got %T want *Call", exp)
	}
	if call.Name != "soma" || len(call.Args) != 2 {
		t.Fatalf("got %+v", call)
	}
	if call.String() != "soma(x, 2)" {
		t.Fatalf("string: got %q", call.String())
	}
}

func TestParseVarVsCall(t *testing.T) {
	exp, err := parseArith("foo + bar(3)")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	bin, ok := exp.(*OpBin)
	if !ok {
		t.Fatalf("got %T want *OpBin", exp)
	}
	if _, ok := bin.Left.(*Var); !ok {
		t.Fatalf("left: got %T want *Var", bin.Left)
	}
	if _, ok := bin.Right.(*Call); !ok {
		t.Fatalf("right: got %T want *Call", bin.Right)
	}
}

func TestParseProgramWithFunctionAndMain(t *testing.T) {
	input := `var x = 10;
fun inc(n) {
  return n + 1;
}
main {
  return inc(x);
}`
	prog, err := parseFun(input)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(prog.Decls) != 2 {
		t.Fatalf("decls: got %d want 2", len(prog.Decls))
	}
	if _, ok := prog.Decls[0].(*VarDecl); !ok {
		t.Fatalf("decl 0: got %T want *VarDecl", prog.Decls[0])
	}
	fn, ok := prog.Decls[1].(*FunDecl)
	if !ok {
		t.Fatalf("decl 1: got %T want *FunDecl", prog.Decls[1])
	}
	if fn.Name != "inc" || len(fn.Params) != 1 || fn.Params[0] != "n" {
		t.Fatalf("function parsed incorrectly: %+v", fn)
	}
}

func TestParseFunctionNoParams(t *testing.T) {
	input := `fun answer() {
  return 42;
}
main {
  return answer();
}`
	prog, err := parseFun(input)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	fn := prog.Decls[0].(*FunDecl)
	if len(fn.Params) != 0 {
		t.Fatalf("params: got %d want 0", len(fn.Params))
	}
	call, ok := prog.Result.(*Call)
	if !ok || call.Name != "answer" || len(call.Args) != 0 {
		t.Fatalf("result: got %#v", prog.Result)
	}
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
	prog, err := parseFun(input)
	if err != nil {
		t.Fatalf("parse unexpected err: %v", err)
	}
	if err := CheckProgram(prog); err != nil {
		t.Fatalf("semantic err: %v", err)
	}
	v, err := EvalPrograma(prog)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 11 {
		t.Fatalf("eval: got %d want 11", v)
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
	prog, err := parseFun(input)
	if err != nil {
		t.Fatalf("parse unexpected err: %v", err)
	}
	if err := CheckProgram(prog); err != nil {
		t.Fatalf("semantic err: %v", err)
	}
	v, err := EvalPrograma(prog)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 28 {
		t.Fatalf("eval: got %d want 28", v)
	}
}

func TestEvalProgramFunctionNoParams(t *testing.T) {
	input := `fun seed() {
  return 21;
}
main {
  return seed() + seed();
}`
	prog, err := parseFun(input)
	if err != nil {
		t.Fatalf("parse unexpected err: %v", err)
	}
	if err := CheckProgram(prog); err != nil {
		t.Fatalf("semantic err: %v", err)
	}
	v, err := EvalPrograma(prog)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 42 {
		t.Fatalf("eval: got %d want 42", v)
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
	prog, err := parseFun(input)
	if err != nil {
		t.Fatalf("parse unexpected err: %v", err)
	}
	if err := CheckProgram(prog); err != nil {
		t.Fatalf("semantic err: %v", err)
	}
	v, err := EvalPrograma(prog)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 142 {
		t.Fatalf("eval: got %d want 142", v)
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
	prog, err := parseFun(input)
	if err != nil {
		t.Fatalf("parse unexpected err: %v", err)
	}
	if err := CheckProgram(prog); err != nil {
		t.Fatalf("semantic err: %v", err)
	}
	v, err := EvalPrograma(prog)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 120 {
		t.Fatalf("eval: got %d want 120", v)
	}
}

func TestSemanticErrorUndeclaredFunction(t *testing.T) {
	prog, err := parseFun(`main { return foo(); }`)
	if err != nil {
		t.Fatalf("parse unexpected err: %v", err)
	}
	err = CheckProgram(prog)
	if err == nil {
		t.Fatalf("expected semantic error")
	}
	if err.Error() != "Erro semantico: funcao 'foo' nao declarada" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestSemanticErrorWrongArity(t *testing.T) {
	prog, err := parseFun(`fun inc(x) { return x + 1; } main { return inc(1, 2); }`)
	if err != nil {
		t.Fatalf("parse unexpected err: %v", err)
	}
	err = CheckProgram(prog)
	if err == nil {
		t.Fatalf("expected semantic error")
	}
	if err.Error() != "Erro semantico: funcao 'inc' esperava 1 parametros, recebeu 2" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestSemanticErrorVariableUsedAsFunction(t *testing.T) {
	prog, err := parseFun(`var foo = 1; main { return foo(); }`)
	if err != nil {
		t.Fatalf("parse unexpected err: %v", err)
	}
	err = CheckProgram(prog)
	if err == nil {
		t.Fatalf("expected semantic error")
	}
	if err.Error() != "Erro semantico: funcao 'foo' nao declarada" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestSemanticErrorDuplicateLocal(t *testing.T) {
	prog, err := parseFun(`fun f(x) { var x = 1; return x; } main { return f(3); }`)
	if err != nil {
		t.Fatalf("parse unexpected err: %v", err)
	}
	err = CheckProgram(prog)
	if err == nil {
		t.Fatalf("expected semantic error")
	}
	if err.Error() != "Erro semantico: simbolo 'x' ja declarado" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestSemanticErrorCallToLaterFunctionRejected(t *testing.T) {
	prog, err := parseFun(`fun a() { return b(); } fun b() { return 1; } main { return a(); }`)
	if err != nil {
		t.Fatalf("parse unexpected err: %v", err)
	}
	err = CheckProgram(prog)
	if err == nil {
		t.Fatalf("expected semantic error")
	}
	if err.Error() != "Erro semantico: funcao 'b' nao declarada" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}
