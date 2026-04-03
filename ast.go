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
	OpMenor      Operator = "<"
	OpMaior      Operator = ">"
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

// --- Declarações e Programa ---

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

// --- Avaliação (intérprete) ---

type evalState struct {
	globals map[string]int
	funcs   map[string]*FunDecl
}

func EvalPrograma(prog *Programa) (int, error) {
	st := &evalState{
		globals: map[string]int{},
		funcs:   map[string]*FunDecl{},
	}

	for _, decl := range prog.Decls {
		switch d := decl.(type) {
		case *VarDecl:
			v, err := st.evalExp(d.Exp, nil)
			if err != nil {
				return 0, err
			}
			st.globals[d.Name] = v
		case *FunDecl:
			st.funcs[d.Name] = d
		default:
			return 0, fmt.Errorf("Erro de execucao: declaracao desconhecida")
		}
	}

	for _, cmd := range prog.Cmds {
		if err := st.evalCmd(cmd, nil); err != nil {
			return 0, err
		}
	}

	return st.evalExp(prog.Result, nil)
}

func (st *evalState) evalCmd(cmd Cmd, local map[string]int) error {
	switch c := cmd.(type) {
	case *IfCmd:
		v, err := st.evalExp(c.Cond, local)
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
			if err := st.evalCmd(sub, local); err != nil {
				return err
			}
		}
	case *WhileCmd:
		for {
			v, err := st.evalExp(c.Cond, local)
			if err != nil {
				return err
			}
			if v == 0 {
				break
			}
			for _, sub := range c.Body {
				if err := st.evalCmd(sub, local); err != nil {
					return err
				}
			}
		}
	case *AtribCmd:
		v, err := st.evalExp(c.Exp, local)
		if err != nil {
			return err
		}
		if local != nil {
			if _, ok := local[c.Name]; ok {
				local[c.Name] = v
				return nil
			}
		}
		if _, ok := st.globals[c.Name]; ok {
			st.globals[c.Name] = v
			return nil
		}
		return fmt.Errorf("Erro de execucao: variavel '%s' nao definida", c.Name)
	}
	return nil
}

func (st *evalState) evalExp(exp Exp, local map[string]int) (int, error) {
	switch e := exp.(type) {
	case *Const:
		return e.Value, nil
	case *Var:
		if local != nil {
			if v, ok := local[e.Name]; ok {
				return v, nil
			}
		}
		v, ok := st.globals[e.Name]
		if !ok {
			return 0, fmt.Errorf("Erro de execucao: variavel '%s' nao definida", e.Name)
		}
		return v, nil
	case *Call:
		fn, ok := st.funcs[e.Name]
		if !ok {
			return 0, fmt.Errorf("Erro de execucao: funcao '%s' nao definida", e.Name)
		}
		if len(e.Args) != len(fn.Params) {
			return 0, fmt.Errorf("Erro de execucao: funcao '%s' esperava %d parametros, recebeu %d", e.Name, len(fn.Params), len(e.Args))
		}
		values := make([]int, len(e.Args))
		for i := len(e.Args) - 1; i >= 0; i-- {
			v, err := st.evalExp(e.Args[i], local)
			if err != nil {
				return 0, err
			}
			values[i] = v
		}
		fnLocal := map[string]int{}
		for i, name := range fn.Params {
			fnLocal[name] = values[i]
		}
		for _, decl := range fn.Locals {
			v, err := st.evalExp(decl.Exp, fnLocal)
			if err != nil {
				return 0, err
			}
			fnLocal[decl.Name] = v
		}
		for _, cmd := range fn.Cmds {
			if err := st.evalCmd(cmd, fnLocal); err != nil {
				return 0, err
			}
		}
		return st.evalExp(fn.Result, fnLocal)
	case *OpBin:
		r, err := st.evalExp(e.Right, local)
		if err != nil {
			return 0, err
		}
		l, err := st.evalExp(e.Left, local)
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
	switch n := e.(type) {
	case *OpBin:
		children := []Exp{n.Left, n.Right}
		for i, child := range children {
			isTail := i == len(children)-1
			writeTreeNode(b, child, prefix, isTail)
		}
	case *Call:
		for i, child := range n.Args {
			isTail := i == len(n.Args)-1
			writeTreeNode(b, child, prefix, isTail)
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
