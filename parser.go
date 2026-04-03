package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	lex  *Lexer
	cur  Token
	peek Token
}

func NewParser(lex *Lexer) (*Parser, error) {
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

func (p *Parser) consume(tt TokenType) error {
	if p.cur.Type != tt {
		return p.syntaxError()
	}
	return p.advance()
}

// ParseProgram: <decl>* 'main' '{' <cmd>* 'return' <exp> ';' '}'
func (p *Parser) ParseProgram() (*Programa, error) {
	var decls []Decl
	for p.cur.Type == TokenVar || p.cur.Type == TokenFun {
		decl, err := p.parseDecl()
		if err != nil {
			return nil, err
		}
		decls = append(decls, decl)
	}
	if err := p.consume(TokenMain); err != nil {
		return nil, err
	}
	if err := p.consume(TokenChaveEsq); err != nil {
		return nil, err
	}
	cmds, err := p.parseCmdList()
	if err != nil {
		return nil, err
	}
	if err := p.consume(TokenReturn); err != nil {
		return nil, err
	}
	result, err := p.parseExp()
	if err != nil {
		return nil, err
	}
	if err := p.consume(TokenPontoVirgula); err != nil {
		return nil, err
	}
	if err := p.consume(TokenChaveDir); err != nil {
		return nil, err
	}
	if p.cur.Type != TokenEOF {
		return nil, p.syntaxError()
	}
	return &Programa{Decls: decls, Cmds: cmds, Result: result}, nil
}

func (p *Parser) parseDecl() (Decl, error) {
	switch p.cur.Type {
	case TokenVar:
		return p.parseVarDecl()
	case TokenFun:
		return p.parseFunDecl()
	default:
		return nil, p.syntaxError()
	}
}

// parseVarDecl: 'var' <ident> '=' <exp> ';'
func (p *Parser) parseVarDecl() (*VarDecl, error) {
	if err := p.consume(TokenVar); err != nil {
		return nil, err
	}
	name := p.cur.Lexeme
	if err := p.consume(TokenIdent); err != nil {
		return nil, err
	}
	if err := p.consume(TokenIgual); err != nil {
		return nil, err
	}
	exp, err := p.parseExp()
	if err != nil {
		return nil, err
	}
	if err := p.consume(TokenPontoVirgula); err != nil {
		return nil, err
	}
	return &VarDecl{Name: name, Exp: exp}, nil
}

// parseFunDecl: 'fun' <ident> '(' <arglist>? ')' '{' <vardecl>* <cmd>* 'return' <exp> ';' '}'
func (p *Parser) parseFunDecl() (*FunDecl, error) {
	if err := p.consume(TokenFun); err != nil {
		return nil, err
	}
	name := p.cur.Lexeme
	if err := p.consume(TokenIdent); err != nil {
		return nil, err
	}
	if err := p.consume(TokenParenEsq); err != nil {
		return nil, err
	}
	params, err := p.parseIdentList()
	if err != nil {
		return nil, err
	}
	if err := p.consume(TokenParenDir); err != nil {
		return nil, err
	}
	if err := p.consume(TokenChaveEsq); err != nil {
		return nil, err
	}
	var locals []*VarDecl
	for p.cur.Type == TokenVar {
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
	if err := p.consume(TokenReturn); err != nil {
		return nil, err
	}
	result, err := p.parseExp()
	if err != nil {
		return nil, err
	}
	if err := p.consume(TokenPontoVirgula); err != nil {
		return nil, err
	}
	if err := p.consume(TokenChaveDir); err != nil {
		return nil, err
	}
	return &FunDecl{Name: name, Params: params, Locals: locals, Cmds: cmds, Result: result}, nil
}

func (p *Parser) parseIdentList() ([]string, error) {
	if p.cur.Type == TokenParenDir {
		return nil, nil
	}

	var names []string
	for {
		if p.cur.Type != TokenIdent {
			return nil, p.syntaxError()
		}
		names = append(names, p.cur.Lexeme)
		if err := p.consume(TokenIdent); err != nil {
			return nil, err
		}
		if p.cur.Type != TokenVirgula {
			return names, nil
		}
		if err := p.consume(TokenVirgula); err != nil {
			return nil, err
		}
	}
}

func (p *Parser) parseExpList() ([]Exp, error) {
	if p.cur.Type == TokenParenDir {
		return nil, nil
	}

	var exps []Exp
	for {
		exp, err := p.parseExp()
		if err != nil {
			return nil, err
		}
		exps = append(exps, exp)
		if p.cur.Type != TokenVirgula {
			return exps, nil
		}
		if err := p.consume(TokenVirgula); err != nil {
			return nil, err
		}
	}
}

// parseCmdList: <cmd>*  (termina quando não é if/while/ident)
func (p *Parser) parseCmdList() ([]Cmd, error) {
	var cmds []Cmd
	for p.cur.Type == TokenIf || p.cur.Type == TokenWhile || p.cur.Type == TokenIdent {
		cmd, err := p.parseCmd()
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, cmd)
	}
	return cmds, nil
}

// parseCmd: <if> | <while> | <atrib>
func (p *Parser) parseCmd() (Cmd, error) {
	switch p.cur.Type {
	case TokenIf:
		return p.parseIfCmd()
	case TokenWhile:
		return p.parseWhileCmd()
	case TokenIdent:
		return p.parseAtribCmd()
	default:
		return nil, p.syntaxError()
	}
}

// parseIfCmd: 'if' <exp> '{' <cmd>* '}' 'else' '{' <cmd>* '}'
func (p *Parser) parseIfCmd() (*IfCmd, error) {
	if err := p.consume(TokenIf); err != nil {
		return nil, err
	}
	cond, err := p.parseExp()
	if err != nil {
		return nil, err
	}
	if err := p.consume(TokenChaveEsq); err != nil {
		return nil, err
	}
	then, err := p.parseCmdList()
	if err != nil {
		return nil, err
	}
	if err := p.consume(TokenChaveDir); err != nil {
		return nil, err
	}
	if err := p.consume(TokenElse); err != nil {
		return nil, err
	}
	if err := p.consume(TokenChaveEsq); err != nil {
		return nil, err
	}
	els, err := p.parseCmdList()
	if err != nil {
		return nil, err
	}
	if err := p.consume(TokenChaveDir); err != nil {
		return nil, err
	}
	return &IfCmd{Cond: cond, Then: then, Else: els}, nil
}

// parseWhileCmd: 'while' <exp> '{' <cmd>* '}'
func (p *Parser) parseWhileCmd() (*WhileCmd, error) {
	if err := p.consume(TokenWhile); err != nil {
		return nil, err
	}
	cond, err := p.parseExp()
	if err != nil {
		return nil, err
	}
	if err := p.consume(TokenChaveEsq); err != nil {
		return nil, err
	}
	body, err := p.parseCmdList()
	if err != nil {
		return nil, err
	}
	if err := p.consume(TokenChaveDir); err != nil {
		return nil, err
	}
	return &WhileCmd{Cond: cond, Body: body}, nil
}

// parseAtribCmd: <ident> '=' <exp> ';'
func (p *Parser) parseAtribCmd() (*AtribCmd, error) {
	name := p.cur.Lexeme
	if err := p.consume(TokenIdent); err != nil {
		return nil, err
	}
	if err := p.consume(TokenIgual); err != nil {
		return nil, err
	}
	exp, err := p.parseExp()
	if err != nil {
		return nil, err
	}
	if err := p.consume(TokenPontoVirgula); err != nil {
		return nil, err
	}
	return &AtribCmd{Name: name, Exp: exp}, nil
}

// parseExp: <exp_a> (('<' | '>' | '==') <exp_a>)*
func (p *Parser) parseExp() (Exp, error) {
	esq, err := p.parseExpA()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == TokenMenor || p.cur.Type == TokenMaior || p.cur.Type == TokenIgualIgual {
		tok := p.cur.Type
		if err := p.advance(); err != nil {
			return nil, err
		}
		dir, err := p.parseExpA()
		if err != nil {
			return nil, err
		}
		var op Operator
		switch tok {
		case TokenMenor:
			op = OpMenor
		case TokenMaior:
			op = OpMaior
		case TokenIgualIgual:
			op = OpIgualIgual
		}
		esq = &OpBin{Op: op, Left: esq, Right: dir}
	}
	return esq, nil
}

// parseExpA: <exp_m> (('+' | '-') <exp_m>)*
func (p *Parser) parseExpA() (Exp, error) {
	esq, err := p.parseExpM()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == TokenSoma || p.cur.Type == TokenSub {
		tok := p.cur.Type
		if err := p.advance(); err != nil {
			return nil, err
		}
		dir, err := p.parseExpM()
		if err != nil {
			return nil, err
		}
		var op Operator
		if tok == TokenSoma {
			op = OpSoma
		} else {
			op = OpSub
		}
		esq = &OpBin{Op: op, Left: esq, Right: dir}
	}
	return esq, nil
}

// parseExpM: <prim> (('*' | '/') <prim>)*
func (p *Parser) parseExpM() (Exp, error) {
	esq, err := p.parsePrim()
	if err != nil {
		return nil, err
	}
	for p.cur.Type == TokenMult || p.cur.Type == TokenDiv {
		tok := p.cur.Type
		if err := p.advance(); err != nil {
			return nil, err
		}
		dir, err := p.parsePrim()
		if err != nil {
			return nil, err
		}
		var op Operator
		if tok == TokenMult {
			op = OpMult
		} else {
			op = OpDiv
		}
		esq = &OpBin{Op: op, Left: esq, Right: dir}
	}
	return esq, nil
}

// parsePrim: <num> | <ident> | '(' <exp> ')' | <fun>
func (p *Parser) parsePrim() (Exp, error) {
	switch p.cur.Type {
	case TokenIdent:
		if p.peek.Type == TokenParenEsq {
			return p.parseCall()
		}
		v := &Var{Name: p.cur.Lexeme}
		if err := p.advance(); err != nil {
			return nil, err
		}
		return v, nil
	case TokenNumero:
		v, err := strconv.Atoi(p.cur.Lexeme)
		if err != nil {
			return nil, fmt.Errorf("Erro sintatico na posicao %d", p.cur.Pos)
		}
		n := &Const{Value: v}
		if err := p.advance(); err != nil {
			return nil, err
		}
		return n, nil
	case TokenParenEsq:
		if err := p.advance(); err != nil {
			return nil, err
		}
		exp, err := p.parseExp()
		if err != nil {
			return nil, err
		}
		if err := p.consume(TokenParenDir); err != nil {
			return nil, err
		}
		return exp, nil
	default:
		return nil, p.syntaxError()
	}
}

func (p *Parser) parseCall() (Exp, error) {
	name := p.cur.Lexeme
	if err := p.consume(TokenIdent); err != nil {
		return nil, err
	}
	if err := p.consume(TokenParenEsq); err != nil {
		return nil, err
	}
	args, err := p.parseExpList()
	if err != nil {
		return nil, err
	}
	if err := p.consume(TokenParenDir); err != nil {
		return nil, err
	}
	return &Call{Name: name, Args: args}, nil
}
