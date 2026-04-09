package interpreter

import (
	"fmt"

	"ativ10/internal/fun"
)

type evalState struct {
	globals map[string]int
	funcs   map[string]*fun.FunDecl
}

func EvalPrograma(prog *fun.Programa) (int, error) {
	st := &evalState{
		globals: map[string]int{},
		funcs:   map[string]*fun.FunDecl{},
	}

	for _, decl := range prog.Decls {
		switch d := decl.(type) {
		case *fun.VarDecl:
			v, err := st.evalExp(d.Exp, nil)
			if err != nil {
				return 0, err
			}
			st.globals[d.Name] = v
		case *fun.FunDecl:
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

func (st *evalState) evalCmd(cmd fun.Cmd, local map[string]int) error {
	switch c := cmd.(type) {
	case *fun.IfCmd:
		v, err := st.evalExp(c.Cond, local)
		if err != nil {
			return err
		}
		branch := c.Else
		if v != 0 {
			branch = c.Then
		}
		for _, sub := range branch {
			if err := st.evalCmd(sub, local); err != nil {
				return err
			}
		}
	case *fun.WhileCmd:
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
	case *fun.AtribCmd:
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

func (st *evalState) evalExp(exp fun.Exp, local map[string]int) (int, error) {
	switch e := exp.(type) {
	case *fun.Const:
		return e.Value, nil
	case *fun.Var:
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
	case *fun.Call:
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
	case *fun.OpBin:
		switch e.Op {
		case fun.OpAnd:
			l, err := st.evalExp(e.Left, local)
			if err != nil {
				return 0, err
			}
			if l == 0 {
				return 0, nil
			}
			r, err := st.evalExp(e.Right, local)
			if err != nil {
				return 0, err
			}
			if r != 0 {
				return 1, nil
			}
			return 0, nil
		case fun.OpOr:
			l, err := st.evalExp(e.Left, local)
			if err != nil {
				return 0, err
			}
			if l != 0 {
				return 1, nil
			}
			r, err := st.evalExp(e.Right, local)
			if err != nil {
				return 0, err
			}
			if r != 0 {
				return 1, nil
			}
			return 0, nil
		}

		r, err := st.evalExp(e.Right, local)
		if err != nil {
			return 0, err
		}
		l, err := st.evalExp(e.Left, local)
		if err != nil {
			return 0, err
		}
		switch e.Op {
		case fun.OpSoma:
			return l + r, nil
		case fun.OpSub:
			return l - r, nil
		case fun.OpMult:
			return l * r, nil
		case fun.OpDiv:
			if r == 0 {
				return 0, fmt.Errorf("Erro de execucao: divisao por zero")
			}
			return l / r, nil
		case fun.OpMod:
			if r == 0 {
				return 0, fmt.Errorf("Erro de execucao: divisao por zero")
			}
			return l % r, nil
		case fun.OpMenor:
			if l < r {
				return 1, nil
			}
			return 0, nil
		case fun.OpMaior:
			if l > r {
				return 1, nil
			}
			return 0, nil
		case fun.OpMenorIgual:
			if l <= r {
				return 1, nil
			}
			return 0, nil
		case fun.OpMaiorIgual:
			if l >= r {
				return 1, nil
			}
			return 0, nil
		case fun.OpIgualIgual:
			if l == r {
				return 1, nil
			}
			return 0, nil
		case fun.OpDiferente:
			if l != r {
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
