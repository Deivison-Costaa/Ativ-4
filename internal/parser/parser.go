package parser

import (
	"fmt"
	"strconv"

	"ativ10/internal/fun"
	"ativ10/internal/lexer"
)

type Parser struct {
	lex  *lexer.Lexer
	cur  fun.Token
	peek fun.Token
}

func NewParser(lex *lexer.Lexer) (*Parser, error) {
	p := &Parser{lex: lex}

	first, err := p.lex.NextToken()
	if err != nil {
		return nil, err
	}
	second, err := p.lex.NextToken()
	if err != nil {
		return nil, err
	}

	p.cur = first
	p.peek = second
	return p, nil
}

func (p *Parser) advance() error {
	p.cur = p.peek
	tok, err := p.lex.NextToken()
	if err != nil {
		return err
	}
	p.peek = tok
	return nil
}

func (p *Parser) syntaxError() error {
	return fmt.Errorf("Erro sintatico na posicao %d", p.cur.Pos)
}

func (p *Parser) consume(tt fun.TokenType) error {
	if p.cur.Type != tt {
		return p.syntaxError()
	}
	return p.advance()
}

func (p *Parser) ParseProgram() (*fun.Programa, error) {
	var decls []fun.Decl
	for p.cur.Type == fun.TokenVar || p.cur.Type == fun.TokenFun {
		decl, err := p.parseDecl()
		if err != nil {
			return nil, err
		}
		decls = append(decls, decl)
	}
	if err := p.consume(fun.TokenMain); err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenChaveEsq); err != nil {
		return nil, err
	}
	cmds, err := p.parseCmdList()
	if err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenReturn); err != nil {
		return nil, err
	}
	result, err := p.parseExp()
	if err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenPontoVirgula); err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenChaveDir); err != nil {
		return nil, err
	}
	if p.cur.Type != fun.TokenEOF {
		return nil, p.syntaxError()
	}
	return &fun.Programa{Decls: decls, Cmds: cmds, Result: result}, nil
}

func (p *Parser) parseDecl() (fun.Decl, error) {
	switch p.cur.Type {
	case fun.TokenVar:
		return p.parseVarDecl()
	case fun.TokenFun:
		return p.parseFunDecl()
	default:
		return nil, p.syntaxError()
	}
}

func (p *Parser) parseVarDecl() (*fun.VarDecl, error) {
	if err := p.consume(fun.TokenVar); err != nil {
		return nil, err
	}
	name := p.cur.Lexeme
	if err := p.consume(fun.TokenIdent); err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenIgual); err != nil {
		return nil, err
	}
	exp, err := p.parseExp()
	if err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenPontoVirgula); err != nil {
		return nil, err
	}
	return &fun.VarDecl{Name: name, Exp: exp}, nil
}

func (p *Parser) parseFunDecl() (*fun.FunDecl, error) {
	if err := p.consume(fun.TokenFun); err != nil {
		return nil, err
	}
	name := p.cur.Lexeme
	if err := p.consume(fun.TokenIdent); err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenParenEsq); err != nil {
		return nil, err
	}
	params, err := p.parseIdentList()
	if err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenParenDir); err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenChaveEsq); err != nil {
		return nil, err
	}

	var locals []*fun.VarDecl
	for p.cur.Type == fun.TokenVar {
		decl, err := p.parseVarDecl()
		if err != nil {
			return nil, err
		}
		locals = append(locals, decl)
	}

	cmds, err := p.parseCmdList()
	if err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenReturn); err != nil {
		return nil, err
	}
	result, err := p.parseExp()
	if err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenPontoVirgula); err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenChaveDir); err != nil {
		return nil, err
	}

	return &fun.FunDecl{Name: name, Params: params, Locals: locals, Cmds: cmds, Result: result}, nil
}

func (p *Parser) parseIdentList() ([]string, error) {
	if p.cur.Type == fun.TokenParenDir {
		return nil, nil
	}

	var names []string
	for {
		if p.cur.Type != fun.TokenIdent {
			return nil, p.syntaxError()
		}
		names = append(names, p.cur.Lexeme)
		if err := p.consume(fun.TokenIdent); err != nil {
			return nil, err
		}
		if p.cur.Type != fun.TokenVirgula {
			return names, nil
		}
		if err := p.consume(fun.TokenVirgula); err != nil {
			return nil, err
		}
	}
}

func (p *Parser) parseExpList() ([]fun.Exp, error) {
	if p.cur.Type == fun.TokenParenDir {
		return nil, nil
	}

	var exps []fun.Exp
	for {
		exp, err := p.parseExp()
		if err != nil {
			return nil, err
		}
		exps = append(exps, exp)
		if p.cur.Type != fun.TokenVirgula {
			return exps, nil
		}
		if err := p.consume(fun.TokenVirgula); err != nil {
			return nil, err
		}
	}
}

func (p *Parser) parseCmdList() ([]fun.Cmd, error) {
	var cmds []fun.Cmd
	for p.cur.Type == fun.TokenIf || p.cur.Type == fun.TokenWhile || p.cur.Type == fun.TokenIdent {
		cmd, err := p.parseCmd()
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}

func (p *Parser) parseCmd() (fun.Cmd, error) {
	switch p.cur.Type {
	case fun.TokenIf:
		return p.parseIfCmd()
	case fun.TokenWhile:
		return p.parseWhileCmd()
	case fun.TokenIdent:
		return p.parseAtribCmd()
	default:
		return nil, p.syntaxError()
	}
}

func (p *Parser) parseIfCmd() (*fun.IfCmd, error) {
	if err := p.consume(fun.TokenIf); err != nil {
		return nil, err
	}
	cond, err := p.parseExp()
	if err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenChaveEsq); err != nil {
		return nil, err
	}
	then, err := p.parseCmdList()
	if err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenChaveDir); err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenElse); err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenChaveEsq); err != nil {
		return nil, err
	}
	els, err := p.parseCmdList()
	if err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenChaveDir); err != nil {
		return nil, err
	}
	return &fun.IfCmd{Cond: cond, Then: then, Else: els}, nil
}

func (p *Parser) parseWhileCmd() (*fun.WhileCmd, error) {
	if err := p.consume(fun.TokenWhile); err != nil {
		return nil, err
	}
	cond, err := p.parseExp()
	if err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenChaveEsq); err != nil {
		return nil, err
	}
	body, err := p.parseCmdList()
	if err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenChaveDir); err != nil {
		return nil, err
	}
	return &fun.WhileCmd{Cond: cond, Body: body}, nil
}

func (p *Parser) parseAtribCmd() (*fun.AtribCmd, error) {
	name := p.cur.Lexeme
	if err := p.consume(fun.TokenIdent); err != nil {
		return nil, err
	}

	var op fun.Operator
	switch p.cur.Type {
	case fun.TokenIgual:
	case fun.TokenSomaIgual:
		op = fun.OpSoma
	case fun.TokenSubIgual:
		op = fun.OpSub
	case fun.TokenMultIgual:
		op = fun.OpMult
	case fun.TokenDivIgual:
		op = fun.OpDiv
	case fun.TokenModIgual:
		op = fun.OpMod
	default:
		return nil, p.syntaxError()
	}
	if err := p.advance(); err != nil {
		return nil, err
	}

	exp, err := p.parseExp()
	if err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenPontoVirgula); err != nil {
		return nil, err
	}

	if op != "" {
		exp = &fun.OpBin{
			Op:    op,
			Left:  &fun.Var{Name: name},
			Right: exp,
		}
	}
	return &fun.AtribCmd{Name: name, Exp: exp}, nil
}

func (p *Parser) parseExp() (fun.Exp, error) {
	return p.parseOr()
}

func (p *Parser) parseOr() (fun.Exp, error) {
	esq, err := p.parseAnd()
	if err != nil {
		return nil, err
	}

	for p.cur.Type == fun.TokenOr {
		tok := p.cur.Type
		if err := p.advance(); err != nil {
			return nil, err
		}
		dir, err := p.parseAnd()
		if err != nil {
			return nil, err
		}

		var op fun.Operator
		switch tok {
		case fun.TokenOr:
			op = fun.OpOr
		}
		esq = &fun.OpBin{Op: op, Left: esq, Right: dir}
	}

	return esq, nil
}

func (p *Parser) parseAnd() (fun.Exp, error) {
	esq, err := p.parseNot()
	if err != nil {
		return nil, err
	}

	for p.cur.Type == fun.TokenAnd {
		if err := p.advance(); err != nil {
			return nil, err
		}
		dir, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		esq = &fun.OpBin{Op: fun.OpAnd, Left: esq, Right: dir}
	}

	return esq, nil
}

func (p *Parser) parseNot() (fun.Exp, error) {
	if p.cur.Type == fun.TokenNot {
		if err := p.advance(); err != nil {
			return nil, err
		}
		exp, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return &fun.OpBin{
			Op:    fun.OpIgualIgual,
			Left:  exp,
			Right: &fun.Const{Value: 0},
		}, nil
	}

	return p.parseComparison()
}

func (p *Parser) parseComparison() (fun.Exp, error) {
	esq, err := p.parseExpA()
	if err != nil {
		return nil, err
	}

	for p.cur.Type == fun.TokenMenor ||
		p.cur.Type == fun.TokenMaior ||
		p.cur.Type == fun.TokenMenorIgual ||
		p.cur.Type == fun.TokenMaiorIgual ||
		p.cur.Type == fun.TokenIgualIgual ||
		p.cur.Type == fun.TokenDiferente {
		tok := p.cur.Type
		if err := p.advance(); err != nil {
			return nil, err
		}
		dir, err := p.parseExpA()
		if err != nil {
			return nil, err
		}

		var op fun.Operator
		switch tok {
		case fun.TokenMenor:
			op = fun.OpMenor
		case fun.TokenMaior:
			op = fun.OpMaior
		case fun.TokenMenorIgual:
			op = fun.OpMenorIgual
		case fun.TokenMaiorIgual:
			op = fun.OpMaiorIgual
		case fun.TokenIgualIgual:
			op = fun.OpIgualIgual
		case fun.TokenDiferente:
			op = fun.OpDiferente
		}
		esq = &fun.OpBin{Op: op, Left: esq, Right: dir}
	}
	return esq, nil
}

func (p *Parser) parseExpA() (fun.Exp, error) {
	esq, err := p.parseExpM()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == fun.TokenSoma || p.cur.Type == fun.TokenSub {
		tok := p.cur.Type
		if err := p.advance(); err != nil {
			return nil, err
		}
		dir, err := p.parseExpM()
		if err != nil {
			return nil, err
		}

		op := fun.OpSub
		if tok == fun.TokenSoma {
			op = fun.OpSoma
		}
		esq = &fun.OpBin{Op: op, Left: esq, Right: dir}
	}
	return esq, nil
}

func (p *Parser) parseExpM() (fun.Exp, error) {
	esq, err := p.parsePrim()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == fun.TokenMult || p.cur.Type == fun.TokenDiv || p.cur.Type == fun.TokenMod {
		tok := p.cur.Type
		if err := p.advance(); err != nil {
			return nil, err
		}
		dir, err := p.parsePrim()
		if err != nil {
			return nil, err
		}

		op := fun.OpDiv
		if tok == fun.TokenMult {
			op = fun.OpMult
		} else if tok == fun.TokenMod {
			op = fun.OpMod
		}
		esq = &fun.OpBin{Op: op, Left: esq, Right: dir}
	}
	return esq, nil
}

func (p *Parser) parsePrim() (fun.Exp, error) {
	switch p.cur.Type {
	case fun.TokenTrue:
		if err := p.advance(); err != nil {
			return nil, err
		}
		return &fun.Const{Value: 1}, nil
	case fun.TokenFalse:
		if err := p.advance(); err != nil {
			return nil, err
		}
		return &fun.Const{Value: 0}, nil
	case fun.TokenIdent:
		if p.peek.Type == fun.TokenParenEsq {
			return p.parseCall()
		}
		v := &fun.Var{Name: p.cur.Lexeme}
		if err := p.advance(); err != nil {
			return nil, err
		}
		return v, nil
	case fun.TokenNumero:
		v, err := strconv.Atoi(p.cur.Lexeme)
		if err != nil {
			return nil, fmt.Errorf("Erro sintatico na posicao %d", p.cur.Pos)
		}
		n := &fun.Const{Value: v}
		if err := p.advance(); err != nil {
			return nil, err
		}
		return n, nil
	case fun.TokenParenEsq:
		if err := p.advance(); err != nil {
			return nil, err
		}
		exp, err := p.parseExp()
		if err != nil {
			return nil, err
		}
		if err := p.consume(fun.TokenParenDir); err != nil {
			return nil, err
		}
		return exp, nil
	default:
		return nil, p.syntaxError()
	}
}

func (p *Parser) parseCall() (fun.Exp, error) {
	name := p.cur.Lexeme
	if err := p.consume(fun.TokenIdent); err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenParenEsq); err != nil {
		return nil, err
	}
	args, err := p.parseExpList()
	if err != nil {
		return nil, err
	}
	if err := p.consume(fun.TokenParenDir); err != nil {
		return nil, err
	}
	return &fun.Call{Name: name, Args: args}, nil
}
