package main

import "testing"

// parseArith analisa apenas uma expressão aritmética/comparação (sem programa completo).
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

// parseCmd analisa um programa Cmd completo.
func parseCmd(input string) (*Programa, error) {
	p, err := NewParser(NewLexer(input))
	if err != nil {
		return nil, err
	}
	return p.ParseProgram()
}

// --- Testes de expressões aritméticas ---

func TestParseConst(t *testing.T) {
	exp, err := parseArith("333")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if exp.String() != "333" {
		t.Fatalf("tree: got %q want %q", exp.String(), "333")
	}
}

func TestParseTreeAndEval(t *testing.T) {
	input := "(33 + (912 * 11))"
	exp, err := parseArith(input)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if exp.String() != input {
		t.Fatalf("tree: got %q want %q", exp.String(), input)
	}
	if TreeString(exp) != "+\n|-- 33\n`-- *\n    |-- 912\n    `-- 11" {
		t.Fatalf("visual tree: got %q", TreeString(exp))
	}
}

func TestParseWhitespace(t *testing.T) {
	input := "(  3\t+\n(4+5) )"
	exp, err := parseArith(input)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if exp.String() != "(3 + (4 + 5))" {
		t.Fatalf("tree: got %q want %q", exp.String(), "(3 + (4 + 5))")
	}
}

func TestSyntaxError(t *testing.T) {
	_, err := parseArith("(3 + )")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "Erro sintatico na posicao 5" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestSyntaxErrorMissingClosingParen(t *testing.T) {
	_, err := parseArith("(3 + 4")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "Erro sintatico na posicao 6" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestSyntaxErrorExtraTokensAfterExpression(t *testing.T) {
	_, err := parseArith("(3 + 4) 5")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "Erro sintatico na posicao 8" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestSyntaxErrorUnknownOperator(t *testing.T) {
	_, err := parseArith("(3 ^ 4)")
	if err == nil {
		t.Fatalf("expected error")
	}
	if err.Error() != "Erro lexico na posicao 3" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestPrecedenciaMultAntesSoma(t *testing.T) {
	exp, err := parseArith("7 + 5 * 3")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if exp.String() != "(7 + (5 * 3))" {
		t.Fatalf("tree: got %q want %q", exp.String(), "(7 + (5 * 3))")
	}
	v, err := evalExpEnv(exp, nil)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 22 {
		t.Fatalf("eval: got %d want %d", v, 22)
	}
}

func TestAssociatividadeEsquerda(t *testing.T) {
	exp, err := parseArith("10 - 8 - 2")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	v, err := evalExpEnv(exp, nil)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 0 {
		t.Fatalf("eval: got %d want %d", v, 0)
	}
}

func TestParentesesMudamPrecedencia(t *testing.T) {
	exp, err := parseArith("(7 + 5) * 3")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	v, err := evalExpEnv(exp, nil)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 36 {
		t.Fatalf("eval: got %d want %d", v, 36)
	}
}

// --- Testes de operadores de comparação ---

func TestComparacaoMenor(t *testing.T) {
	exp, err := parseArith("3 < 5")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	v, err := evalExpEnv(exp, nil)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 1 {
		t.Fatalf("eval: got %d want 1", v)
	}
}

func TestComparacaoMenorFalso(t *testing.T) {
	exp, err := parseArith("5 < 3")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	v, err := evalExpEnv(exp, nil)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 0 {
		t.Fatalf("eval: got %d want 0", v)
	}
}

func TestComparacaoMaior(t *testing.T) {
	exp, err := parseArith("7 > 2")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	v, err := evalExpEnv(exp, nil)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 1 {
		t.Fatalf("eval: got %d want 1", v)
	}
}

func TestComparacaoIgualdade(t *testing.T) {
	exp, err := parseArith("4 == 4")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	v, err := evalExpEnv(exp, nil)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 1 {
		t.Fatalf("eval: got %d want 1", v)
	}
}

func TestComparacaoPrecedencia(t *testing.T) {
	// 1 + 2 < 3 + 1 => (1+2) < (3+1) => 3 < 4 => 1
	exp, err := parseArith("1 + 2 < 3 + 1")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	v, err := evalExpEnv(exp, nil)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 1 {
		t.Fatalf("eval: got %d want 1", v)
	}
}

// --- Testes de programas Cmd ---

func TestCmd_ProgramaSemDeclaracoes(t *testing.T) {
	prog, err := parseCmd("{ return 7 + 5 * 3; }")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(prog.Decls) != 0 {
		t.Fatalf("decls: got %d want 0", len(prog.Decls))
	}
	v, err := EvalPrograma(prog)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 22 {
		t.Fatalf("eval: got %d want %d", v, 22)
	}
}

func TestCmd_ProgramaComDeclaracao(t *testing.T) {
	prog, err := parseCmd("x = 30;\n{ return x + 10; }")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(prog.Decls) != 1 {
		t.Fatalf("decls: got %d want 1", len(prog.Decls))
	}
	v, err := EvalPrograma(prog)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 40 {
		t.Fatalf("eval: got %d want %d", v, 40)
	}
}

func TestCmd_If(t *testing.T) {
	input := `x = 5;
{
  if x < 10 {
    x = x + 1;
  } else {
    x = x - 1;
  }
  return x;
}`
	prog, err := parseCmd(input)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	v, err := EvalPrograma(prog)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 6 {
		t.Fatalf("eval: got %d want %d", v, 6)
	}
}

func TestCmd_While(t *testing.T) {
	input := `n = 1;
soma = 0;
{
  while n < 5 {
    soma = soma + n;
    n = n + 1;
  }
  return soma;
}`
	prog, err := parseCmd(input)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	v, err := EvalPrograma(prog)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	// soma = 1+2+3+4 = 10
	if v != 10 {
		t.Fatalf("eval: got %d want %d", v, 10)
	}
}

func TestCmd_ErroVariavelNaoDeclaradaEmDecl(t *testing.T) {
	prog, err := parseCmd("x = 7 + y;\n{ return x; }")
	if err != nil {
		t.Fatalf("unexpected parse err: %v", err)
	}
	err = CheckProgram(prog)
	if err == nil {
		t.Fatalf("expected semantic error")
	}
	if err.Error() != "Erro semantico: variavel 'y' nao declarada" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestCmd_ErroVariavelNaoDeclaradaEmResultado(t *testing.T) {
	prog, err := parseCmd("{ return x; }")
	if err != nil {
		t.Fatalf("unexpected parse err: %v", err)
	}
	err = CheckProgram(prog)
	if err == nil {
		t.Fatalf("expected semantic error")
	}
	if err.Error() != "Erro semantico: variavel 'x' nao declarada" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestCmd_ErroAtribVariavelNaoDeclarada(t *testing.T) {
	input := `{
  x = 5;
  return x;
}`
	prog, err := parseCmd(input)
	if err != nil {
		t.Fatalf("unexpected parse err: %v", err)
	}
	err = CheckProgram(prog)
	if err == nil {
		t.Fatalf("expected semantic error")
	}
	if err.Error() != "Erro semantico: variavel 'x' nao declarada" {
		t.Fatalf("unexpected error: %q", err.Error())
	}
}

func TestCmd_Delta(t *testing.T) {
	// |b^2 - 4ac| com a=1, b=2, c=3 => |4 - 12| = 8
	input := `a = 1;
b = 2;
c = 3;
delta = b * b - 4 * a * c;
{
  if delta < 0 {
    delta = 0 - delta;
  } else {
    delta = delta;
  }
  return delta;
}`
	prog, err := parseCmd(input)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if err := CheckProgram(prog); err != nil {
		t.Fatalf("semantic err: %v", err)
	}
	v, err := EvalPrograma(prog)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 8 {
		t.Fatalf("eval: got %d want %d", v, 8)
	}
}

func TestCmd_SomaWhile(t *testing.T) {
	// soma de 1..9 = 45
	input := `n = 1;
m = 10;
soma = 0;
{
  while n < m {
    soma = soma + n;
    n = n + 1;
  }
  return soma;
}`
	prog, err := parseCmd(input)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if err := CheckProgram(prog); err != nil {
		t.Fatalf("semantic err: %v", err)
	}
	v, err := EvalPrograma(prog)
	if err != nil {
		t.Fatalf("eval unexpected err: %v", err)
	}
	if v != 45 {
		t.Fatalf("eval: got %d want %d", v, 45)
	}
}
