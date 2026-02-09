package main

import (
	"fmt"
	"strings"
)

type CodeGenerator struct {
	builder strings.Builder
}

func NewCodeGenerator() *CodeGenerator {
	return &CodeGenerator{}
}

func (cg *CodeGenerator) Generate(exp Exp) string {
	cg.builder.Reset()

	// Cabeçalho do assembly
	cg.emit("#")
	cg.emit("# Código gerado pelo compilador EC1")
	cg.emit("#")
	cg.emit("")
	cg.emit(".section .text")
	cg.emit(".globl _start")
	cg.emit("")
	cg.emit("_start:")

	// Gera código para a expressão
	cg.genExp(exp)

	// Chamadas para imprimir e sair
	cg.emit("")
	cg.emit("  call imprime_num")
	cg.emit("  call sair")
	cg.emit("")
	cg.emit(".include \"runtime.s\"")

	return cg.builder.String()
}

func (cg *CodeGenerator) emit(s string) {
	cg.builder.WriteString(s)
	cg.builder.WriteString("\n")
}

func (cg *CodeGenerator) genExp(exp Exp) {
	switch e := exp.(type) {
	case *Const:
		cg.genConst(e)
	case *OpBin:
		cg.genOpBin(e)
	}
}

func (cg *CodeGenerator) genConst(c *Const) {
	cg.emit(fmt.Sprintf("  mov $%d, %%rax", c.Value))
}

func (cg *CodeGenerator) genOpBin(op *OpBin) {
	// Esquema de tradução com ordem invertida (direito primeiro):
	// 1. Gerar código para operando direito
	// 2. push %rax (salvar resultado do direito)
	// 3. Gerar código para operando esquerdo
	// 4. pop %rbx (recuperar resultado do direito em RBX)
	// 5. Executar operação (RAX = RAX op RBX)

	// 1. Gera código para operando direito
	cg.genExp(op.Right)

	// 2. Salva resultado na pilha
	cg.emit("  push %rax")

	// 3. Gera código para operando esquerdo
	cg.genExp(op.Left)

	// 4. Recupera operando direito em RBX
	cg.emit("  pop %rbx")

	// 5. Executa a operação
	switch op.Op {
	case OpSoma:
		cg.emit("  add %rbx, %rax")
	case OpSub:
		cg.emit("  sub %rbx, %rax")
	case OpMult:
		cg.emit("  imul %rbx, %rax")
	case OpDiv:
		// Para divisão: RAX / RBX
		// cqo estende RAX para RDX:RAX (sinal)
		// idiv %rbx divide RDX:RAX por RBX, quociente em RAX
		cg.emit("  cqo")
		cg.emit("  idiv %rbx")
	}
}
