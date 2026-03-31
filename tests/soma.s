# Código gerado pelo compilador Cmd

.section .bss
  .lcomm n, 8
  .lcomm m, 8
  .lcomm soma, 8

.section .text
.globl _start

_start:
  # n = ...
  mov $1, %rax
  mov %rax, n
  # m = ...
  mov $10, %rax
  mov %rax, m
  # soma = ...
  mov $0, %rax
  mov %rax, soma
Linicio0:
  mov m, %rax
  push %rax
  mov n, %rax
  pop %rbx
  xor %rcx, %rcx
  cmp %rbx, %rax
  setl %cl
  mov %rcx, %rax
  cmp $0, %rax
  jz Lfim0
  mov n, %rax
  push %rax
  mov soma, %rax
  pop %rbx
  add %rbx, %rax
  mov %rax, soma
  mov $1, %rax
  push %rax
  mov n, %rax
  pop %rbx
  add %rbx, %rax
  mov %rax, n
  jmp Linicio0
Lfim0:
  mov soma, %rax

  call imprime_num
  call sair

.include "runtime/runtime.s"
