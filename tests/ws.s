# Código gerado pelo compilador EV

.section .bss
  .lcomm x, 8

.section .text
.globl _start

_start:
  # x = ...
  mov $5, %rax
  push %rax
  mov $4, %rax
  pop %rbx
  add %rbx, %rax
  push %rax
  mov $3, %rax
  pop %rbx
  add %rbx, %rax
  mov %rax, x
  mov $2, %rax
  push %rax
  mov x, %rax
  pop %rbx
  imul %rbx, %rax

  call imprime_num
  call sair

.include "runtime/runtime.s"
