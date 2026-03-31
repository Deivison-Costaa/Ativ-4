package main

import (
	"fmt"
	"strings"
)

// --- Expressões ---

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

type Operator string

const (
	OpSoma      Operator = "+"
	OpSub       Operator = "-"
	OpMult      Operator = "*"
	OpDiv       Operator = "/"
	OpMenor     Operator = "<"
	OpMaior     Operator = ">"
	OpIgualIgual Operator = "=="
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

// --- Comandos ---

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

// --- Declaração e Programa ---

type Decl struct {
	Name string
	Exp  Exp
}

type Programa struct {
	Decls  []*Decl
	Cmds   []Cmd
	Result Exp
}

// --- Avaliação (intérprete) ---

func EvalPrograma(prog *Programa) (int, error) {
	env := map[string]int{}
	for _, decl := range prog.Decls {
		v, err := evalExpEnv(decl.Exp, env)
		if err != nil {
			return 0, err
		}
		env[decl.Name] = v
	}
	for _, cmd := range prog.Cmds {
		if err := evalCmd(cmd, env); err != nil {
			return 0, err
		}
	}
	return evalExpEnv(prog.Result, env)
}

func evalCmd(cmd Cmd, env map[string]int) error {
	switch c := cmd.(type) {
	case *IfCmd:
		v, err := evalExpEnv(c.Cond, env)
		if err != nil {
			return err
		}
		var branch []Cmd
		if v != 0 {
			branch = c.Then
		} else {
			branch = c.Else
		}
		for _, sub := range branch {
			if err := evalCmd(sub, env); err != nil {
				return err
			}
		}
	case *WhileCmd:
		for {
			v, err := evalExpEnv(c.Cond, env)
			if err != nil {
				return err
			}
			if v == 0 {
				break
			}
			for _, sub := range c.Body {
				if err := evalCmd(sub, env); err != nil {
					return err
				}
			}
		}
	case *AtribCmd:
		v, err := evalExpEnv(c.Exp, env)
		if err != nil {
			return err
		}
		env[c.Name] = v
	}
	return nil
}

func evalExpEnv(exp Exp, env map[string]int) (int, error) {
	switch e := exp.(type) {
	case *Const:
		return e.Value, nil
	case *Var:
		v, ok := env[e.Name]
		if !ok {
			return 0, fmt.Errorf("Erro de execucao: variavel '%s' nao definida", e.Name)
		}
		return v, nil
	case *OpBin:
		l, err := evalExpEnv(e.Left, env)
		if err != nil {
			return 0, err
		}
		r, err := evalExpEnv(e.Right, env)
		if err != nil {
			return 0, err
		}
		switch e.Op {
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
		case OpMenor:
			if l < r {
				return 1, nil
			}
			return 0, nil
		case OpMaior:
			if l > r {
				return 1, nil
			}
			return 0, nil
		case OpIgualIgual:
			if l == r {
				return 1, nil
			}
			return 0, nil
		default:
			return 0, fmt.Errorf("Erro de execucao: operador desconhecido")
		}
	default:
		return 0, fmt.Errorf("Erro de execucao: no desconhecido")
	}
}

// --- Impressão da árvore ---

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
	case *Var:
		return n.Name
	default:
		return "<?>"
	}
}
