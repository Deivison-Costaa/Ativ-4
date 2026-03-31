# Código gerado pelo compilador Cmd

.section .bss
  .lcomm m, 8
  .lcomm n, 8

.section .text
.globl _start

_start:
  # m = ...
  mov $10, %rax
  mov %rax, m
  # n = ...
  mov $4, %rax
  mov %rax, n
Linicio0:
  mov n, %rax
  push %rax
  mov $1, %rax
  push %rax
  mov m, %rax
  pop %rbx
  add %rbx, %rax
  pop %rbx
  xor %rcx, %rcx
  cmp %rbx, %rax
  setg %cl
  mov %rcx, %rax
  cmp $0, %rax
  jz Lfim0
  mov n, %rax
  push %rax
  mov m, %rax
  pop %rbx
  sub %rbx, %rax
  mov %rax, m
  jmp Linicio0
Lfim0:
  mov m, %rax

  call imprime_num
  call sair

.include "runtime/runtime.s"
