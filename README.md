# Atividade 07 — EC1 (Expressões Constantes 1) — Compilador

Este projeto implementa um **compilador completo** para a linguagem **EC1**, gerando código **assembly x86-64** para Linux.

## Componentes

- **Análise léxica** (scanner) — produz tokens
- **Análise sintática descendente recursiva** (parser) — produz uma **árvore de sintaxe abstrata (AST)**
- **Geração de código** — gera assembly x86-64 usando pilha para armazenamento de valores intermediários

## Linguagem EC1

Um programa EC1 é uma expressão aritmética com literais inteiros e operadores `+`, `-`, `*`, `/`.

A gramática é:

```text
<programa> ::= <expressao>
<expressao> ::= <literal-inteiro> | '(' <expressao> <operador> <expressao> ')'
<operador> ::= '+' | '-' | '*' | '/'
<literal-inteiro> ::= <digito>+
```

Exemplos de programas EC1:
```
333
(6 * 7)
(3 + (4 + (11 + 7)))
(33 + (912 * 11))
((427 / 7) + (11 * (231 + 5)))
```

## Como usar

### Compilar o compilador EC1

```bash
go build -o ec1 *.go
```

### Compilar um programa EC1

```bash
# Gera arquivo .s a partir do .ec1
./ec1 programa.ec1

# Ou especifica arquivo de saída
./ec1 -o saida.s programa.ec1

# Ou lê do stdin e imprime na stdout
echo "(6 * 7)" | ./ec1 -o saida.s -
```

### Montar e executar o programa

```bash
# Montar (a partir do diretório onde está o runtime.s)
as -o programa.o programa.s

# Linkar
ld -o programa programa.o

# Executar
./programa
```

### Exemplo completo

```bash
# Compilar EC1 → Assembly
./ec1 tests/input.ec1

# Montar e linkar
as -o tests/input.o tests/input.s
ld -o tests/input tests/input.o

# Executar
./tests/input
# Saída: 10065
```

## Geração de Código

O compilador usa o seguinte esquema de tradução:

### Constante
```asm
mov $valor, %rax
```

### Operação Binária
1. Gerar código para operando direito
2. `push %rax` (salvar resultado)
3. Gerar código para operando esquerdo
4. `pop %rbx` (recuperar operando direito)
5. Executar operação

Operações:
- Soma: `add %rbx, %rax`
- Subtração: `sub %rbx, %rax`
- Multiplicação: `imul %rbx, %rax`
- Divisão: `cqo` + `idiv %rbx`

### Modelo de Saída

```asm
#
# Código gerado pelo compilador EC1
#

.section .text
.globl _start

_start:
  # código da expressão aqui
  # resultado final fica em %rax

  call imprime_num
  call sair

.include "runtime.s"
```

## Arquivos

- `token.go` - Definição dos tokens
- `lexer.go` - Analisador léxico
- `ast.go` - Definição da AST
- `parser.go` - Analisador sintático
- `codegen.go` - Gerador de código assembly
- `main.go` - Programa principal
- `runtime.s` - Funções auxiliares (imprime_num, sair)
- `run_tests.sh` - Script de testes automatizados

## Testes

### Executar testes automatizados

```bash
bash run_tests.sh
```

### Executar testes unitários do Go

```bash
go test ./...
```

## Erros

O programa detecta e reporta:

- **Erros léxicos**: `Erro lexico na posicao X`
- **Erros sintáticos**: `Erro sintatico na posicao X`

## Requisitos

- Go 1.22+
- GNU Assembler (`as`)
- GNU Linker (`ld`)
- Linux x86-64
