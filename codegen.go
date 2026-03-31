package main

import (
	"fmt"
	"strings"
)

type CodeGenerator struct {
	builder    strings.Builder
	labelCount int
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

func (cg *CodeGenerator) newLabel() int {
	n := cg.labelCount
	cg.labelCount++
	return n
}

func (cg *CodeGenerator) Generate(prog *Programa) string {
	cg.builder.Reset()
	cg.labelCount = 0

	cg.emit("# Código gerado pelo compilador Cmd")
	cg.emit("")

	// Seção BSS: variáveis globais
	if len(prog.Decls) > 0 {
		cg.emit(".section .bss")
		for _, decl := range prog.Decls {
			cg.emit(fmt.Sprintf("  .lcomm %s, 8", decl.Name))
		}
		cg.emit("")
	}

	cg.emit(".section .text")
	cg.emit(".globl _start")
	cg.emit("")
	cg.emit("_start:")

	// Inicializar variáveis globais
	for _, decl := range prog.Decls {
		cg.emit(fmt.Sprintf("  # %s = ...", decl.Name))
		cg.genExp(decl.Exp)
		cg.emit(fmt.Sprintf("  mov %%rax, %s", decl.Name))
	}

	// Comandos do bloco principal
	for _, cmd := range prog.Cmds {
		cg.genCmd(cmd)
	}

	// Expressão de resultado
	cg.genExp(prog.Result)

	cg.emit("")
	cg.emit("  call imprime_num")
	cg.emit("  call sair")
	cg.emit("")
	cg.emit(".include \"runtime/runtime.s\"")

	return cg.builder.String()
}

func (cg *CodeGenerator) emit(s string) {
	cg.builder.WriteString(s)
	cg.builder.WriteString("\n")
}

// genCmd despacha geração de código para cada tipo de comando.
func (cg *CodeGenerator) genCmd(cmd Cmd) {
	switch c := cmd.(type) {
	case *IfCmd:
		cg.genIfCmd(c)
	case *WhileCmd:
		cg.genWhileCmd(c)
	case *AtribCmd:
		cg.genAtribCmd(c)
	}
}

// genIfCmd gera código para: if E { C1 } else { C2 }
//
//	<codigo_E>
//	cmp $0, %rax
//	jz LfalsoN
//	<codigo_C1>
//	jmp LfimN
//
// LfalsoN:
//
//	<codigo_C2>
//
// LfimN:
func (cg *CodeGenerator) genIfCmd(c *IfCmd) {
	n := cg.newLabel()
	lfalso := fmt.Sprintf("Lfalso%d", n)
	lfim := fmt.Sprintf("Lfim%d", n)

	cg.genExp(c.Cond)
	cg.emit("  cmp $0, %rax")
	cg.emit(fmt.Sprintf("  jz %s", lfalso))
	for _, sub := range c.Then {
		cg.genCmd(sub)
	}
	cg.emit(fmt.Sprintf("  jmp %s", lfim))
	cg.emit(fmt.Sprintf("%s:", lfalso))
	for _, sub := range c.Else {
		cg.genCmd(sub)
	}
	cg.emit(fmt.Sprintf("%s:", lfim))
}

// genWhileCmd gera código para: while E { C }
//
// LinicioN:
//
//	<codigo_E>
//	cmp $0, %rax
//	jz LfimN
//	<codigo_C>
//	jmp LinicioN
//
// LfimN:
func (cg *CodeGenerator) genWhileCmd(c *WhileCmd) {
	n := cg.newLabel()
	linicio := fmt.Sprintf("Linicio%d", n)
	lfim := fmt.Sprintf("Lfim%d", n)

	cg.emit(fmt.Sprintf("%s:", linicio))
	cg.genExp(c.Cond)
	cg.emit("  cmp $0, %rax")
	cg.emit(fmt.Sprintf("  jz %s", lfim))
	for _, sub := range c.Body {
		cg.genCmd(sub)
	}
	cg.emit(fmt.Sprintf("  jmp %s", linicio))
	cg.emit(fmt.Sprintf("%s:", lfim))
}

// genAtribCmd gera código para: var = exp
func (cg *CodeGenerator) genAtribCmd(c *AtribCmd) {
	cg.genExp(c.Exp)
	cg.emit(fmt.Sprintf("  mov %%rax, %s", c.Name))
}

// genExp gera código para uma expressão (resultado fica em %rax).
func (cg *CodeGenerator) genExp(exp Exp) {
	switch e := exp.(type) {
	case *Const:
		cg.emit(fmt.Sprintf("  mov $%d, %%rax", e.Value))
	case *Var:
		cg.emit(fmt.Sprintf("  mov %s, %%rax", e.Name))
	case *OpBin:
		cg.genOpBin(e)
	}
}

func (cg *CodeGenerator) genOpBin(op *OpBin) {
	switch op.Op {
	case OpMenor, OpMaior, OpIgualIgual:
		cg.genComparison(op)
	default:
		// Operações aritméticas: gera direito primeiro, depois esquerdo
		cg.genExp(op.Right)
		cg.emit("  push %rax")
		cg.genExp(op.Left)
		cg.emit("  pop %rbx")
		switch op.Op {
		case OpSoma:
			cg.emit("  add %rbx, %rax")
		case OpSub:
			cg.emit("  sub %rbx, %rax")
		case OpMult:
			cg.emit("  imul %rbx, %rax")
		case OpDiv:
			cg.emit("  cqo")
			cg.emit("  idiv %rbx")
		}
	}
}

// genComparison gera código para operadores de comparação.
// Após execução, %rax contém 1 (verdadeiro) ou 0 (falso).
//
//	<codigo_right>
//	push %rax
//	<codigo_left>
//	pop %rbx
//	xor %rcx, %rcx
//	cmp %rbx, %rax   # calcula left - right, seta flags
//	set? %cl         # setz/setl/setg conforme operador
//	mov %rcx, %rax
func (cg *CodeGenerator) genComparison(op *OpBin) {
	cg.genExp(op.Right)
	cg.emit("  push %rax")
	cg.genExp(op.Left)
	cg.emit("  pop %rbx")
	cg.emit("  xor %rcx, %rcx")
	cg.emit("  cmp %rbx, %rax")
	switch op.Op {
	case OpIgualIgual:
		cg.emit("  setz %cl")
	case OpMenor:
		cg.emit("  setl %cl")
	case OpMaior:
		cg.emit("  setg %cl")
	}
	cg.emit("  mov %rcx, %rax")
}
