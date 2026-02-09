package main

import (
	"fmt"
	"strings"
)

type Exp interface {
	Eval() (int, error)
	String() string
}

type Const struct {
	Value int
}

func (c *Const) Eval() (int, error) {
	return c.Value, nil
}

func (c *Const) String() string {
	return fmt.Sprintf("%d", c.Value)
}

type Operator string

const (
	OpSoma Operator = "+"
	OpSub  Operator = "-"
	OpMult Operator = "*"
	OpDiv  Operator = "/"
)

type OpBin struct {
	Op    Operator
	Left  Exp
	Right Exp
}

func (b *OpBin) Eval() (int, error) {
	l, err := b.Left.Eval()
	if err != nil {
		return 0, err
	}
	r, err := b.Right.Eval()
	if err != nil {
		return 0, err
	}

	switch b.Op {
	case OpSoma:
		return l + r, nil
	case OpSub:
		return l - r, nil
	case OpMult:
		return l * r, nil
	case OpDiv:
		if r == 0 {
			return 0, fmt.Errorf("Erro de execucao: divisao por zero")
		}
		return l / r, nil
	default:
		return 0, fmt.Errorf("Erro de execucao: operador desconhecido")
	}
}

func (b *OpBin) String() string {
	return fmt.Sprintf("(%s %s %s)", b.Left.String(), string(b.Op), b.Right.String())
}

func TreeString(e Exp) string {
	var b strings.Builder
	b.WriteString(treeLabel(e))
	b.WriteString("\n")
	writeTreeChildren(&b, e, "")
	return strings.TrimRight(b.String(), "\n")
}

func writeTreeChildren(b *strings.Builder, e Exp, prefix string) {
	op, ok := e.(*OpBin)
	if !ok {
		return
	}

	children := []Exp{op.Left, op.Right}
	for i, child := range children {
		isTail := i == len(children)-1
		writeTreeNode(b, child, prefix, isTail)
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
	default:
		return "<?>"
	}
}
