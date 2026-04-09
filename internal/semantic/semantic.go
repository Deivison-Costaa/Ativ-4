package semantic

import (
	"fmt"

	"ativ10/internal/fun"
)

type SymbolKind string

const (
	SymbolGlobalVar SymbolKind = "global_var"
	SymbolFunction  SymbolKind = "function"
	SymbolLocalVar  SymbolKind = "local_var"
)

type Symbol struct {
	Kind  SymbolKind
	Arity int
}

func CheckProgram(prog *fun.Programa) error {
	globals := map[string]Symbol{}

	for _, decl := range prog.Decls {
		switch d := decl.(type) {
		case *fun.VarDecl:
			if _, exists := globals[d.Name]; exists {
				return fmt.Errorf("Erro semantico: simbolo '%s' ja declarado", d.Name)
			}
			if err := checkExp(d.Exp, globals, nil); err != nil {
				return err
			}
			globals[d.Name] = Symbol{Kind: SymbolGlobalVar}
		case *fun.FunDecl:
			if _, exists := globals[d.Name]; exists {
				return fmt.Errorf("Erro semantico: simbolo '%s' ja declarado", d.Name)
			}
			globals[d.Name] = Symbol{Kind: SymbolFunction, Arity: len(d.Params)}
			if err := checkFunction(d, globals); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Erro semantico: declaracao desconhecida")
		}
	}

	for _, cmd := range prog.Cmds {
		if err := checkCmd(cmd, globals, nil); err != nil {
			return err
		}
	}

	return checkExp(prog.Result, globals, nil)
}

func checkFunction(fn *fun.FunDecl, globals map[string]Symbol) error {
	locals := map[string]Symbol{}

	for _, param := range fn.Params {
		if _, exists := locals[param]; exists {
			return fmt.Errorf("Erro semantico: simbolo '%s' ja declarado", param)
		}
		locals[param] = Symbol{Kind: SymbolLocalVar}
	}

	for _, decl := range fn.Locals {
		if _, exists := locals[decl.Name]; exists {
			return fmt.Errorf("Erro semantico: simbolo '%s' ja declarado", decl.Name)
		}
		if err := checkExp(decl.Exp, globals, locals); err != nil {
			return err
		}
		locals[decl.Name] = Symbol{Kind: SymbolLocalVar}
	}

	for _, cmd := range fn.Cmds {
		if err := checkCmd(cmd, globals, locals); err != nil {
			return err
		}
	}

	return checkExp(fn.Result, globals, locals)
}

func checkCmd(cmd fun.Cmd, globals map[string]Symbol, locals map[string]Symbol) error {
	switch c := cmd.(type) {
	case *fun.IfCmd:
		if err := checkExp(c.Cond, globals, locals); err != nil {
			return err
		}
		for _, sub := range c.Then {
			if err := checkCmd(sub, globals, locals); err != nil {
				return err
			}
		}
		for _, sub := range c.Else {
			if err := checkCmd(sub, globals, locals); err != nil {
				return err
			}
		}
	case *fun.WhileCmd:
		if err := checkExp(c.Cond, globals, locals); err != nil {
			return err
		}
		for _, sub := range c.Body {
			if err := checkCmd(sub, globals, locals); err != nil {
				return err
			}
		}
	case *fun.AtribCmd:
		if err := checkExp(c.Exp, globals, locals); err != nil {
			return err
		}
		if err := checkVarRef(c.Name, globals, locals); err != nil {
			return err
		}
	}
	return nil
}

func checkVarRef(name string, globals map[string]Symbol, locals map[string]Symbol) error {
	if locals != nil {
		if sym, ok := locals[name]; ok && sym.Kind == SymbolLocalVar {
			return nil
		}
	}
	sym, ok := globals[name]
	if !ok || sym.Kind != SymbolGlobalVar {
		return fmt.Errorf("Erro semantico: variavel '%s' nao declarada", name)
	}
	return nil
}

func checkExp(exp fun.Exp, globals map[string]Symbol, locals map[string]Symbol) error {
	switch e := exp.(type) {
	case *fun.Var:
		return checkVarRef(e.Name, globals, locals)
	case *fun.Call:
		sym, ok := globals[e.Name]
		if !ok || sym.Kind != SymbolFunction {
			return fmt.Errorf("Erro semantico: funcao '%s' nao declarada", e.Name)
		}
		if len(e.Args) != sym.Arity {
			return fmt.Errorf("Erro semantico: funcao '%s' esperava %d parametros, recebeu %d", e.Name, sym.Arity, len(e.Args))
		}
		for _, arg := range e.Args {
			if err := checkExp(arg, globals, locals); err != nil {
				return err
			}
		}
	case *fun.OpBin:
		if err := checkExp(e.Left, globals, locals); err != nil {
			return err
		}
		if err := checkExp(e.Right, globals, locals); err != nil {
			return err
		}
	}
	return nil
}
