package parser

import (
	"testing"

	"ativ10/internal/fun"
	"ativ10/internal/lexer"
)

func parseArith(input string) (fun.Exp, error) {
	p, err := NewParser(lexer.NewLexer(input))
	if err != nil {
		return nil, err
	}
	exp, err := p.parseExp()
	if err != nil {
		return nil, err
	}
	if p.cur.Type != fun.TokenEOF {
		return nil, p.syntaxError()
	}
	return exp, nil
}

func parseFun(input string) (*fun.Programa, error) {
	p, err := NewParser(lexer.NewLexer(input))
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
	call, ok := exp.(*fun.Call)
	if !ok {
		t.Fatalf("got %T want *fun.Call", exp)
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
	bin, ok := exp.(*fun.OpBin)
	if !ok {
		t.Fatalf("got %T want *fun.OpBin", exp)
	}
	if _, ok := bin.Left.(*fun.Var); !ok {
		t.Fatalf("left: got %T want *fun.Var", bin.Left)
	}
	if _, ok := bin.Right.(*fun.Call); !ok {
		t.Fatalf("right: got %T want *fun.Call", bin.Right)
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
	if _, ok := prog.Decls[0].(*fun.VarDecl); !ok {
		t.Fatalf("decl 0: got %T want *fun.VarDecl", prog.Decls[0])
	}
	fn, ok := prog.Decls[1].(*fun.FunDecl)
	if !ok {
		t.Fatalf("decl 1: got %T want *fun.FunDecl", prog.Decls[1])
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
	fn := prog.Decls[0].(*fun.FunDecl)
	if len(fn.Params) != 0 {
		t.Fatalf("params: got %d want 0", len(fn.Params))
	}
	call, ok := prog.Result.(*fun.Call)
	if !ok || call.Name != "answer" || len(call.Args) != 0 {
		t.Fatalf("result: got %#v", prog.Result)
	}
}

func TestParseBooleanOperatorPrecedence(t *testing.T) {
	exp, err := parseArith("true and false or true")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	root, ok := exp.(*fun.OpBin)
	if !ok || root.Op != fun.OpOr {
		t.Fatalf("root: got %#v want OpOr", exp)
	}
	left, ok := root.Left.(*fun.OpBin)
	if !ok || left.Op != fun.OpAnd {
		t.Fatalf("left: got %#v want OpAnd", root.Left)
	}
	if c, ok := left.Left.(*fun.Const); !ok || c.Value != 1 {
		t.Fatalf("left.left: got %#v want true literal", left.Left)
	}
	if c, ok := left.Right.(*fun.Const); !ok || c.Value != 0 {
		t.Fatalf("left.right: got %#v want false literal", left.Right)
	}
	if c, ok := root.Right.(*fun.Const); !ok || c.Value != 1 {
		t.Fatalf("right: got %#v want true literal", root.Right)
	}
}

func TestParseNotBindsOverComparisonGroup(t *testing.T) {
	exp, err := parseArith("not 1 == 1")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	root, ok := exp.(*fun.OpBin)
	if !ok || root.Op != fun.OpIgualIgual {
		t.Fatalf("root: got %#v want equality from not desugaring", exp)
	}
	comp, ok := root.Left.(*fun.OpBin)
	if !ok || comp.Op != fun.OpIgualIgual {
		t.Fatalf("root.Left: got %#v want comparison expression", root.Left)
	}
	if c, ok := root.Right.(*fun.Const); !ok || c.Value != 0 {
		t.Fatalf("root.Right: got %#v want zero literal", root.Right)
	}
}

func TestParseNotFalseAndTrue(t *testing.T) {
	exp, err := parseArith("not false and true")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	root, ok := exp.(*fun.OpBin)
	if !ok || root.Op != fun.OpAnd {
		t.Fatalf("root: got %#v want OpAnd", exp)
	}
	left, ok := root.Left.(*fun.OpBin)
	if !ok || left.Op != fun.OpIgualIgual {
		t.Fatalf("root.Left: got %#v want equality from not desugaring", root.Left)
	}
	if c, ok := root.Right.(*fun.Const); !ok || c.Value != 1 {
		t.Fatalf("root.Right: got %#v want true literal", root.Right)
	}
}

func TestParseModuloHasMultiplicativePrecedence(t *testing.T) {
	exp, err := parseArith("2 + 9 % 4 * 3")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	root, ok := exp.(*fun.OpBin)
	if !ok || root.Op != fun.OpSoma {
		t.Fatalf("root: got %#v want OpSoma", exp)
	}
	right, ok := root.Right.(*fun.OpBin)
	if !ok || right.Op != fun.OpMult {
		t.Fatalf("right: got %#v want OpMult", root.Right)
	}
	left, ok := right.Left.(*fun.OpBin)
	if !ok || left.Op != fun.OpMod {
		t.Fatalf("right.Left: got %#v want OpMod", right.Left)
	}
}

func TestParseNewComparisonOperators(t *testing.T) {
	exp, err := parseArith("1 <= 2 and 3 >= 2 and 4 != 5")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	root, ok := exp.(*fun.OpBin)
	if !ok || root.Op != fun.OpAnd {
		t.Fatalf("root: got %#v want OpAnd", exp)
	}
	rightCmp, ok := root.Right.(*fun.OpBin)
	if !ok || rightCmp.Op != fun.OpDiferente {
		t.Fatalf("root.Right: got %#v want OpDiferente", root.Right)
	}
	leftAnd, ok := root.Left.(*fun.OpBin)
	if !ok || leftAnd.Op != fun.OpAnd {
		t.Fatalf("root.Left: got %#v want OpAnd", root.Left)
	}
	if cmp, ok := leftAnd.Left.(*fun.OpBin); !ok || cmp.Op != fun.OpMenorIgual {
		t.Fatalf("left cmp: got %#v want OpMenorIgual", leftAnd.Left)
	}
	if cmp, ok := leftAnd.Right.(*fun.OpBin); !ok || cmp.Op != fun.OpMaiorIgual {
		t.Fatalf("right cmp: got %#v want OpMaiorIgual", leftAnd.Right)
	}
}

func TestParseCompoundAssignmentDesugars(t *testing.T) {
	prog, err := parseFun(`var x = 10; main { x %= 4; return x; }`)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	cmd, ok := prog.Cmds[0].(*fun.AtribCmd)
	if !ok {
		t.Fatalf("cmd: got %T want *fun.AtribCmd", prog.Cmds[0])
	}
	bin, ok := cmd.Exp.(*fun.OpBin)
	if !ok || bin.Op != fun.OpMod {
		t.Fatalf("cmd.Exp: got %#v want OpMod", cmd.Exp)
	}
	if v, ok := bin.Left.(*fun.Var); !ok || v.Name != "x" {
		t.Fatalf("bin.Left: got %#v want Var(x)", bin.Left)
	}
	if c, ok := bin.Right.(*fun.Const); !ok || c.Value != 4 {
		t.Fatalf("bin.Right: got %#v want Const(4)", bin.Right)
	}
}
