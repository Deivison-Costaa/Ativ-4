# Código gerado pelo compilador Cmd

.section .bss
  .lcomm a, 8
  .lcomm b, 8
  .lcomm r, 8

.section .text
.globl _start

_start:
  # a = ...
  mov $18, %rax
  mov %rax, a
  # b = ...
  mov $12, %rax
  mov %rax, b
  # r = ...
  mov $0, %rax
  mov %rax, r
  mov a, %rax
  mov %rax, r
Linicio0:
  mov b, %rax
  push %rax
  mov $1, %rax
  push %rax
  mov r, %rax
  pop %rbx
  add %rbx, %rax
  pop %rbx
  xor %rcx, %rcx
  cmp %rbx, %rax
  setg %cl
  mov %rcx, %rax
  cmp $0, %rax
  jz Lfim0
  mov b, %rax
  push %rax
  mov r, %rax
  pop %rbx
  sub %rbx, %rax
  mov %rax, r
  jmp Linicio0
Lfim0:
Linicio1:
  mov $0, %rax
  push %rax
  mov r, %rax
  pop %rbx
  xor %rcx, %rcx
  cmp %rbx, %rax
  setg %cl
  mov %rcx, %rax
  cmp $0, %rax
  jz Lfim1
  mov b, %rax
  mov %rax, a
  mov r, %rax
  mov %rax, b
  mov a, %rax
  mov %rax, r
Linicio2:
  mov b, %rax
  push %rax
  mov $1, %rax
  push %rax
  mov r, %rax
  pop %rbx
  add %rbx, %rax
  pop %rbx
  xor %rcx, %rcx
  cmp %rbx, %rax
  setg %cl
  mov %rcx, %rax
  cmp $0, %rax
  jz Lfim2
  mov b, %rax
  push %rax
  mov r, %rax
  pop %rbx
  sub %rbx, %rax
  mov %rax, r
  jmp Linicio2
Lfim2:
  jmp Linicio1
Lfim1:
  mov b, %rax

  call imprime_num
  call sair

.include "runtime/runtime.s"
