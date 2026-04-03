package main

import (
	"fmt"
	"strings"
)

type functionContext struct {
	slots      map[string]int
	localCount int
}

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

	cg.emit("# Codigo gerado pelo compilador Fun")
	cg.emit("")

	if cg.hasGlobalVars(prog) {
		cg.emit(".section .bss")
		for _, decl := range prog.Decls {
			if d, ok := decl.(*VarDecl); ok {
				cg.emit(fmt.Sprintf("  .lcomm %s, 8", d.Name))
			}
		}
		cg.emit("")
	}

	cg.emit(".section .text")
	cg.emit(".globl _start")
	cg.emit("")
	cg.emit("_start:")

	for _, decl := range prog.Decls {
		d, ok := decl.(*VarDecl)
		if !ok {
			continue
		}
		cg.emit(fmt.Sprintf("  # %s = ...", d.Name))
		cg.genExp(d.Exp, nil)
		cg.emit(fmt.Sprintf("  mov %%rax, %s(%%rip)", d.Name))
	}

	for _, cmd := range prog.Cmds {
		cg.genCmd(cmd, nil)
	}
	cg.genExp(prog.Result, nil)

	cg.emit("")
	cg.emit("  call imprime_num")
	cg.emit("  call sair")
	cg.emit("")

	for _, decl := range prog.Decls {
		if fn, ok := decl.(*FunDecl); ok {
			cg.genFunction(fn)
			cg.emit("")
		}
	}

	cg.emit(".include \"runtime/runtime.s\"")
	return cg.builder.String()
}

func (cg *CodeGenerator) hasGlobalVars(prog *Programa) bool {
	for _, decl := range prog.Decls {
		if _, ok := decl.(*VarDecl); ok {
			return true
		}
	}
	return false
}

func (cg *CodeGenerator) emit(s string) {
	cg.builder.WriteString(s)
	cg.builder.WriteString("\n")
}

func (cg *CodeGenerator) genFunction(fn *FunDecl) {
	ctx := buildFunctionContext(fn)
	cg.emit(fmt.Sprintf("%s:", fn.Name))
	cg.emit("  push %rbp")
	if ctx.localCount > 0 {
		cg.emit(fmt.Sprintf("  sub $%d, %%rsp", ctx.localCount*8))
	}
	cg.emit("  mov %rsp, %rbp")

	for _, decl := range fn.Locals {
		cg.genExp(decl.Exp, ctx)
		cg.emit(fmt.Sprintf("  mov %%rax, %d(%%rbp)", ctx.slots[decl.Name]))
	}
	for _, cmd := range fn.Cmds {
		cg.genCmd(cmd, ctx)
	}
	cg.genExp(fn.Result, ctx)

	if ctx.localCount > 0 {
		cg.emit(fmt.Sprintf("  add $%d, %%rsp", ctx.localCount*8))
	}
	cg.emit("  pop %rbp")
	cg.emit("  ret")
}

func buildFunctionContext(fn *FunDecl) *functionContext {
	slots := map[string]int{}
	for i, decl := range fn.Locals {
		slots[decl.Name] = i * 8
	}
	baseParamOffset := len(fn.Locals)*8 + 16
	for i, param := range fn.Params {
		slots[param] = baseParamOffset + i*8
	}
	return &functionContext{
		slots:      slots,
		localCount: len(fn.Locals),
	}
}

func (cg *CodeGenerator) genCmd(cmd Cmd, ctx *functionContext) {
	switch c := cmd.(type) {
	case *IfCmd:
		cg.genIfCmd(c, ctx)
	case *WhileCmd:
		cg.genWhileCmd(c, ctx)
	case *AtribCmd:
		cg.genAtribCmd(c, ctx)
	}
}

func (cg *CodeGenerator) genIfCmd(c *IfCmd, ctx *functionContext) {
	n := cg.newLabel()
	lfalso := fmt.Sprintf("Lfalso%d", n)
	lfim := fmt.Sprintf("Lfim%d", n)

	cg.genExp(c.Cond, ctx)
	cg.emit("  cmp $0, %rax")
	cg.emit(fmt.Sprintf("  jz %s", lfalso))
	for _, sub := range c.Then {
		cg.genCmd(sub, ctx)
	}
	cg.emit(fmt.Sprintf("  jmp %s", lfim))
	cg.emit(fmt.Sprintf("%s:", lfalso))
	for _, sub := range c.Else {
		cg.genCmd(sub, ctx)
	}
	cg.emit(fmt.Sprintf("%s:", lfim))
}

func (cg *CodeGenerator) genWhileCmd(c *WhileCmd, ctx *functionContext) {
	n := cg.newLabel()
	linicio := fmt.Sprintf("Linicio%d", n)
	lfim := fmt.Sprintf("Lfim%d", n)

	cg.emit(fmt.Sprintf("%s:", linicio))
	cg.genExp(c.Cond, ctx)
	cg.emit("  cmp $0, %rax")
	cg.emit(fmt.Sprintf("  jz %s", lfim))
	for _, sub := range c.Body {
		cg.genCmd(sub, ctx)
	}
	cg.emit(fmt.Sprintf("  jmp %s", linicio))
	cg.emit(fmt.Sprintf("%s:", lfim))
}

func (cg *CodeGenerator) genAtribCmd(c *AtribCmd, ctx *functionContext) {
	cg.genExp(c.Exp, ctx)
	cg.storeVar(c.Name, ctx)
}

func (cg *CodeGenerator) genExp(exp Exp, ctx *functionContext) {
	switch e := exp.(type) {
	case *Const:
		cg.emit(fmt.Sprintf("  mov $%d, %%rax", e.Value))
	case *Var:
		cg.loadVar(e.Name, ctx)
	case *Call:
		cg.genCall(e, ctx)
	case *OpBin:
		cg.genOpBin(e, ctx)
	}
}

func (cg *CodeGenerator) genCall(call *Call, ctx *functionContext) {
	for i := len(call.Args) - 1; i >= 0; i-- {
		cg.genExp(call.Args[i], ctx)
		cg.emit("  push %rax")
	}
	cg.emit(fmt.Sprintf("  call %s", call.Name))
	if len(call.Args) > 0 {
		cg.emit(fmt.Sprintf("  add $%d, %%rsp", len(call.Args)*8))
	}
}

func (cg *CodeGenerator) loadVar(name string, ctx *functionContext) {
	if ctx != nil {
		if offset, ok := ctx.slots[name]; ok {
			cg.emit(fmt.Sprintf("  mov %d(%%rbp), %%rax", offset))
			return
		}
	}
	cg.emit(fmt.Sprintf("  mov %s(%%rip), %%rax", name))
}

func (cg *CodeGenerator) storeVar(name string, ctx *functionContext) {
	if ctx != nil {
		if offset, ok := ctx.slots[name]; ok {
			cg.emit(fmt.Sprintf("  mov %%rax, %d(%%rbp)", offset))
			return
		}
	}
	cg.emit(fmt.Sprintf("  mov %%rax, %s(%%rip)", name))
}

func (cg *CodeGenerator) genOpBin(op *OpBin, ctx *functionContext) {
	switch op.Op {
	case OpMenor, OpMaior, OpIgualIgual:
		cg.genComparison(op, ctx)
	default:
		cg.genExp(op.Right, ctx)
		cg.emit("  push %rax")
		cg.genExp(op.Left, ctx)
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

func (cg *CodeGenerator) genComparison(op *OpBin, ctx *functionContext) {
	cg.genExp(op.Right, ctx)
	cg.emit("  push %rax")
	cg.genExp(op.Left, ctx)
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
