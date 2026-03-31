# Código gerado pelo compilador Cmd

.section .bss
  .lcomm x, 8
  .lcomm y, 8

.section .text
.globl _start

_start:
  # x = ...
  mov $12, %rax
  push %rax
  mov $4, %rax
  push %rax
  mov $7, %rax
  pop %rbx
  add %rbx, %rax
  pop %rbx
  imul %rbx, %rax
  mov %rax, x
  # y = ...
  mov $11, %rax
  push %rax
  mov $3, %rax
  push %rax
  mov x, %rax
  pop %rbx
  imul %rbx, %rax
  pop %rbx
  add %rbx, %rax
  mov %rax, y
  mov $13, %rax
  push %rax
  mov y, %rax
  pop %rbx
  imul %rbx, %rax
  push %rax
  mov $11, %rax
  push %rax
  mov x, %rax
  pop %rbx
  imul %rbx, %rax
  push %rax
  mov y, %rax
  push %rax
  mov x, %rax
  pop %rbx
  imul %rbx, %rax
  pop %rbx
  add %rbx, %rax
  pop %rbx
  add %rbx, %rax

  call imprime_num
  call sair

.include "runtime/runtime.s"
