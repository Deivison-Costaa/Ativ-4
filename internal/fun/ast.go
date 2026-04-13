package fun

import (
	"fmt"
	"strings"
)

type Exp interface {
	expNode()
	String() string
}

type Const struct {
	Value int
}

func (*Const) expNode() {}

func (c *Const) String() string {
	return fmt.Sprintf("%d", c.Value)
}

type Var struct {
	Name string
}

func (*Var) expNode() {}

func (v *Var) String() string { return v.Name }

type Call struct {
	Name string
	Args []Exp
}

func (*Call) expNode() {}

func (c *Call) String() string {
	if len(c.Args) == 0 {
		return fmt.Sprintf("%s()", c.Name)
	}

	parts := make([]string, 0, len(c.Args))
	for _, arg := range c.Args {
		parts = append(parts, arg.String())
	}

	return fmt.Sprintf("%s(%s)", c.Name, strings.Join(parts, ", "))
}

type Operator string

const (
	OpSoma       Operator = "+"
	OpSub        Operator = "-"
	OpMult       Operator = "*"
	OpDiv        Operator = "/"
	OpMod        Operator = "%"
	OpMenor      Operator = "<"
	OpMaior      Operator = ">"
	OpMenorIgual Operator = "<="
	OpMaiorIgual Operator = ">="
	OpDiferente  Operator = "!="
	OpIgualIgual Operator = "=="
	OpAnd        Operator = "and"
	OpOr         Operator = "or"
)

type OpBin struct {
	Op    Operator
	Left  Exp
	Right Exp
}

func (*OpBin) expNode() {}

func (b *OpBin) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left.String(), string(b.Op), b.Right.String())
}

type Cmd interface {
	cmdNode()
}

type IfCmd struct {
	Cond Exp
	Then []Cmd
	Else []Cmd
}

func (*IfCmd) cmdNode() {}

type WhileCmd struct {
	Cond Exp
	Body []Cmd
}

func (*WhileCmd) cmdNode() {}

type AtribCmd struct {
	Name string
	Exp  Exp
}

func (*AtribCmd) cmdNode() {}

type Decl interface {
	declNode()
}

type VarDecl struct {
	Name string
	Exp  Exp
}

func (*VarDecl) declNode() {}

type FunDecl struct {
	Name   string
	Params []string
	Locals []*VarDecl
	Cmds   []Cmd
	Result Exp
}

func (*FunDecl) declNode() {}

type Programa struct {
	Decls  []Decl
	Cmds   []Cmd
	Result Exp
}

func TreeString(e Exp) string {
	var b strings.Builder
	b.WriteString(treeLabel(e))
	b.WriteString("\n")
	writeTreeChildren(&b, e, "")
	return strings.TrimRight(b.String(), "\n")
}

func writeTreeChildren(b *strings.Builder, e Exp, prefix string) {
	switch n := e.(type) {
	case *OpBin:
		children := []Exp{n.Left, n.Right}
		for i, child := range children {
			writeTreeNode(b, child, prefix, i == len(children)-1)
		}
	case *Call:
		for i, child := range n.Args {
			writeTreeNode(b, child, prefix, i == len(n.Args)-1)
		}
	}
}

func writeTreeNode(b *strings.Builder, e Exp, prefix string, isTail bool) {
	b.WriteString(prefix)
	if isTail {
		b.WriteString("`-- ")
	} else {
		b.WriteString("|-- ")
	}
	b.WriteString(treeLabel(e))
	b.WriteString("\n")

	nextPrefix := prefix
	if isTail {
		nextPrefix += "    "
	} else {
		nextPrefix += "|   "
	}
	writeTreeChildren(b, e, nextPrefix)
}

func treeLabel(e Exp) string {
	switch n := e.(type) {
	case *Const:
		return fmt.Sprintf("%d", n.Value)
	case *OpBin:
		return string(n.Op)
	case *Var:
		return n.Name
	case *Call:
		return n.Name + "()"
	default:
		return "<?>"
	}
}
