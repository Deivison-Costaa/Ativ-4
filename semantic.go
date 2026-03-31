package main

import "fmt"

// CheckProgram verifica erros semânticos no programa.
func CheckProgram(prog *Programa) error {
	syms := map[string]bool{}

	// Verificar declarações globais
	for _, decl := range prog.Decls {
		if err := checkExp(decl.Exp, syms); err != nil {
			return err
		}
		syms[decl.Name] = true
	}

	// Verificar comandos
	for _, cmd := range prog.Cmds {
		if err := checkCmd(cmd, syms); err != nil {
			return err
		}
	}

	// Verificar expressão de resultado
	return checkExp(prog.Result, syms)
}

func checkCmd(cmd Cmd, syms map[string]bool) error {
	switch c := cmd.(type) {
	case *IfCmd:
		if err := checkExp(c.Cond, syms); err != nil {
			return err
		}
		for _, sub := range c.Then {
			if err := checkCmd(sub, syms); err != nil {
				return err
			}
		}
		for _, sub := range c.Else {
			if err := checkCmd(sub, syms); err != nil {
				return err
			}
		}
	case *WhileCmd:
		if err := checkExp(c.Cond, syms); err != nil {
			return err
		}
		for _, sub := range c.Body {
			if err := checkCmd(sub, syms); err != nil {
				return err
			}
		}
	case *AtribCmd:
		if err := checkExp(c.Exp, syms); err != nil {
			return err
		}
		if !syms[c.Name] {
			return fmt.Errorf("Erro semantico: variavel '%s' nao declarada", c.Name)
		}
	}
	return nil
}

func checkExp(exp Exp, syms map[string]bool) error {
	switch e := exp.(type) {
	case *Var:
		if !syms[e.Name] {
			return fmt.Errorf("Erro semantico: variavel '%s' nao declarada", e.Name)
		}
	case *OpBin:
		if err := checkExp(e.Left, syms); err != nil {
			return err
		}
		if err := checkExp(e.Right, syms); err != nil {
			return err
		}
	}
	return nil
}
