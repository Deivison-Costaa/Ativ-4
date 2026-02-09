package main

import (
	"fmt"
	"strconv"
)

type Parser struct {
	lex *Lexer
	cur Token
}

func NewParser(lex *Lexer) (*Parser, error) {
	p := &Parser{lex: lex}
	if err := p.advance(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Parser) advance() error {
	tok, err := p.lex.NextToken()
	if err != nil {
		return err
	}
	p.cur = tok
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

func (p *Parser) ParseProgram() (Exp, error) {
	exp, err := p.parseExp()
	if err != nil {
		return nil, err
	}
	if p.cur.Type != TokenEOF {
		return nil, p.syntaxError()
	}
	return exp, nil
}

func (p *Parser) parseExp() (Exp, error) {
	switch p.cur.Type {
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
		left, err := p.parseExp()
		if err != nil {
			return nil, err
		}
		op, err := p.parseOperator()
		if err != nil {
			return nil, err
		}
		right, err := p.parseExp()
		if err != nil {
			return nil, err
		}
		if err := p.consume(TokenParenDir); err != nil {
			return nil, err
		}
		return &OpBin{Op: op, Left: left, Right: right}, nil
	default:
		return nil, p.syntaxError()
	}
}

func (p *Parser) parseOperator() (Operator, error) {
	switch p.cur.Type {
	case TokenSoma:
		if err := p.advance(); err != nil {
			return "", err
		}
		return OpSoma, nil
	case TokenSub:
		if err := p.advance(); err != nil {
			return "", err
		}
		return OpSub, nil
	case TokenMult:
		if err := p.advance(); err != nil {
			return "", err
		}
		return OpMult, nil
	case TokenDiv:
		if err := p.advance(); err != nil {
			return "", err
		}
		return OpDiv, nil
	default:
		return "", p.syntaxError()
	}
}
