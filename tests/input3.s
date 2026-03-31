# Código gerado pelo compilador Cmd

.section .text
.globl _start

_start:
  mov $3, %rax
  push %rax
  mov $5, %rax
  pop %rbx
  imul %rbx, %rax
  push %rax
  mov $7, %rax
  pop %rbx
  add %rbx, %rax

  call imprime_num
  call sair

.include "runtime/runtime.s"
