# Código gerado pelo compilador Cmd

.section .bss
  .lcomm a, 8
  .lcomm b, 8
  .lcomm c, 8
  .lcomm delta, 8

.section .text
.globl _start

_start:
  # a = ...
  mov $1, %rax
  mov %rax, a
  # b = ...
  mov $2, %rax
  mov %rax, b
  # c = ...
  mov $3, %rax
  mov %rax, c
  # delta = ...
  mov c, %rax
  push %rax
  mov a, %rax
  push %rax
  mov $4, %rax
  pop %rbx
  imul %rbx, %rax
  pop %rbx
  imul %rbx, %rax
  push %rax
  mov b, %rax
  push %rax
  mov b, %rax
  pop %rbx
  imul %rbx, %rax
  pop %rbx
  sub %rbx, %rax
  mov %rax, delta
  mov $0, %rax
  push %rax
  mov delta, %rax
  pop %rbx
  xor %rcx, %rcx
  cmp %rbx, %rax
  setl %cl
  mov %rcx, %rax
  cmp $0, %rax
  jz Lfalso0
  mov delta, %rax
  push %rax
  mov $0, %rax
  pop %rbx
  sub %rbx, %rax
  mov %rax, delta
  jmp Lfim0
Lfalso0:
  mov delta, %rax
  mov %rax, delta
Lfim0:
  mov delta, %rax

  call imprime_num
  call sair

.include "runtime/runtime.s"
