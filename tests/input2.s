# Código gerado pelo compilador Cmd

.section .bss
  .lcomm x, 8
  .lcomm y, 8

.section .text
.globl _start

_start:
  # x = ...
  mov $30, %rax
  mov %rax, x
  # y = ...
  mov $2, %rax
  push %rax
  mov x, %rax
  pop %rbx
  imul %rbx, %rax
  mov %rax, y
  mov y, %rax
  push %rax
  mov x, %rax
  pop %rbx
  add %rbx, %rax

  call imprime_num
  call sair

.include "runtime/runtime.s"
