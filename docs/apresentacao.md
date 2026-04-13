# Compilador Fun — Apresentação do Projeto Final

---

## Sumário

1. [Visão Geral do Projeto](#1-visão-geral-do-projeto)
2. [Arquitetura do Compilador](#2-arquitetura-do-compilador)
3. [A Linguagem Fun Base (Atividade 11)](#3-a-linguagem-fun-base-atividade-11)
   - Análise Léxica
   - Análise Sintática
   - Análise Semântica
   - Geração de Código
4. [Extensões Implementadas (Projeto Final)](#4-extensões-implementadas-projeto-final)
   - Operador de Módulo `%`
   - Novos Comparadores `<=` `>=` `!=`
   - Atribuições Compostas `+=` `-=` `*=` `/=` `%=`
   - Valores Booleanos `true` / `false`
   - Operadores Lógicos `and` / `or` / `not`
5. [Testes](#5-testes)
6. [Como Usar](#6-como-usar)

---

## 1. Visão Geral do Projeto

O projeto é um **compilador completo** para a linguagem Fun, escrito em Go, que gera assembly x86-64 para Linux.

```
programa.fun  →  [compilador]  →  programa.s  →  [as + ld]  →  executável
```

A linguagem Fun é uma linguagem imperativa simples com:
- Variáveis globais e locais
- Funções com parâmetros e retorno de valor
- Estruturas de controle `if/else` e `while`
- Expressões aritméticas e de comparação

### Estrutura do repositório

```
cmd/funcc/          → CLI do compilador
internal/
  fun/              → tokens e AST compartilhados
  lexer/            → análise léxica
  parser/           → análise sintática
  semantic/         → análise semântica
  codegen/          → geração de assembly x86-64
  interpreter/      → interpretador (usado nos testes)
  integration/      → testes end-to-end
runtime/runtime.s   → rotinas auxiliares (imprime_num, sair)
testdata/           → fixtures de teste positivos e negativos
```

---

## 2. Arquitetura do Compilador

O compilador é organizado em **4 fases sequenciais**:

```
Código-fonte
     │
     ▼
┌─────────────┐
│  LÉXICO     │  Transforma texto em tokens
└──────┬──────┘
       │ []Token
       ▼
┌─────────────┐
│  SINTÁTICO  │  Constrói a Árvore Sintática Abstrata (AST)
└──────┬──────┘
       │ *Programa (AST)
       ▼
┌─────────────┐
│  SEMÂNTICO  │  Verifica uso correto de variáveis e funções
└──────┬──────┘
       │ *Programa (validada)
       ▼
┌─────────────┐
│  CODEGEN    │  Gera assembly x86-64 AT&T
└─────────────┘
       │
       ▼
   arquivo .s
```

---

## 3. A Linguagem Fun Base (Atividade 11)

### Gramática resumida

```
programa  ::= decl* 'main' '{' cmd* 'return' exp ';' '}'
decl      ::= vardecl | fundecl
fundecl   ::= 'fun' ident '(' arglist? ')' '{' vardecl* cmd* 'return' exp ';' '}'
vardecl   ::= 'var' ident '=' exp ';'
cmd       ::= if | while | atrib
exp       ::= exp_a (('<' | '>' | '==') exp_a)*
```

---

### Análise Léxica

O lexer reconhece todos os tokens da linguagem:

| Categoria       | Exemplos                            |
|-----------------|-------------------------------------|
| Palavras-chave  | `fun` `var` `main` `if` `else` `while` `return` |
| Identificadores | `x`, `abs`, `fatorial`              |
| Números         | `0`, `42`, `1000`                   |
| Operadores      | `+` `-` `*` `/` `<` `>` `==`       |
| Delimitadores   | `(` `)` `{` `}` `,` `;`            |

**Novo em relação à linguagem Cmd:** token vírgula (`,`) e as palavras-chave `fun`, `var`, `main`.

---

### Análise Sintática

O parser é **recursivo-descendente** com 1 token de lookahead.

**Ponto importante — diferenciar variável de chamada de função:**

```
x      →  próximo token é ';' ou operador  →  referência a variável
x(...)  →  próximo token é '('             →  chamada de função
```

**Exemplo de programa válido:**

```fun
fun abs(x) {
  var y = 0;
  if x < 0 {
    y = 0 - x;
  } else {
    y = x;
  }
  return y;
}

main {
  return abs(8) + abs(0 - 3);
}
```
> Saída: `11`

---

### Análise Semântica

Três verificações principais:

#### 1. Função declarada antes do uso
```fun
main {
  return foo();   // ERRO: função 'foo' não declarada
}
```
```
erro semântico: função 'foo' não declarada
```

#### 2. Aridade (número de parâmetros)
```fun
fun inc(x) { return x + 1; }

main {
  return inc(1, 2);   // ERRO: inc espera 1, recebeu 2
}
```
```
erro semântico: função 'inc' chamada com 2 argumento(s), esperava 1
```

#### 3. Escopo local (shadowing)
```fun
var x = 100;

fun f(x) {
  var y = x + 1;   // x aqui é o parâmetro, não a global
  return y;
}

main {
  return f(41) + x;   // f(41) = 42, x global = 100 → resultado: 142
}
```
> Saída: `142`

A tabela de símbolos tem **escopo global** + **escopo local por função**. Dentro de uma função, a busca consulta primeiro o local, depois o global.

---

### Geração de Código — Funções

O compilador segue as **convenções de chamada** descritas na atividade:

#### Estrutura do registro de ativação na pilha

```
     ┌───────────┐  ← %rbp  (topo das locais)
     │  local a  │  [%rbp + 0]
     ├───────────┤
     │  local b  │  [%rbp + 8]
     ├───────────┤
     │  RBP ant. │  [%rbp + 16]
     ├───────────┤
     │  End. Ret │  [%rbp + 24]
     ├───────────┤
     │  param x  │  [%rbp + 32]
     ├───────────┤
     │  param y  │  [%rbp + 40]
     └───────────┘
```

#### Chamada de função — código gerado

Para `quad(7)` onde `quad` chama `dup` que chama `dup`:

```asm
mov $7, %rax
push %rax          # empilha parâmetro (ordem inversa)
call quad
add $8, %rsp       # limpa parâmetro da pilha após retorno
```

#### Corpo da função — prólogo e epílogo

```asm
abs:
  push %rbp            # salva frame pointer anterior
  sub $8, %rsp         # aloca espaço para variável local y
  mov %rsp, %rbp       # rbp aponta para início das locais

  # ... código da função ...

  add $8, %rsp         # libera espaço das locais
  pop %rbp             # restaura frame pointer anterior
  ret                  # retorna (resultado em %rax)
```

#### Funções encadeadas (chain)

```fun
fun dup(x) { return x + x; }

fun quad(x) { return dup(dup(x)); }

main { return quad(7); }
```
> Saída: `28`

#### Função recursiva (fatorial)

```fun
fun fact(n) {
  if n < 2 {
    n = 1;
  } else {
    n = n * fact(n - 1);
  }
  return n;
}

main { return fact(5); }
```
> Saída: `120`

#### Função sem parâmetros

```fun
fun seed() { return 21; }

main { return seed() + seed(); }
```
> Saída: `42`

---

## 4. Extensões Implementadas (Projeto Final)

O Projeto Final exige **1 extensão de complexidade média/alta** ou **pelo menos 3 extensões simples**.

Foram implementadas **5 extensões simples**, todas integradas em todas as fases do compilador (lexer → parser → semântica → codegen → interpretador → testes).

---

### Extensão 1 — Operador de Módulo `%`

**Lexer:** novo token `TokenMod` para o caractere `%`.

**Parser:** mesmo nível de precedência de `*` e `/` (multiplicativo).

**Codegen:** usa `idiv` e captura o resto em `%rdx`:

```asm
cqo
idiv %rbx
mov %rdx, %rax    # resto da divisão
```

**Exemplo:**
```fun
main {
  return 10 % 3;
}
```
> Saída: `1`

---

### Extensão 2 — Novos Comparadores `<=` `>=` `!=`

**Lexer:** novos tokens `TokenMenorIgual`, `TokenMaiorIgual`, `TokenDiferente`.

**Codegen:** usa instruções `set*` do x86-64:

| Operador | Instrução assembly |
|----------|--------------------|
| `<=`     | `setle %cl`        |
| `>=`     | `setge %cl`        |
| `!=`     | `setne %cl`        |

**Exemplo:**
```fun
main {
  return (3 <= 3) and (5 >= 4) and (7 != 2);
}
```
> Saída: `1` (verdadeiro)

---

### Extensão 3 — Atribuições Compostas `+=` `-=` `*=` `/=` `%=`

**Lexer:** novos tokens `TokenSomaIgual`, `TokenSubIgual`, `TokenMultIgual`, `TokenDivIgual`, `TokenModIgual`.

**Parser:** **dessugarização** — `x += 5` é convertido diretamente para `x = x + 5` na construção da AST. Nenhuma mudança necessária em semântica, codegen ou interpretador.

**Exemplo:**
```fun
var x = 19;

main {
  x += 5;   // x = 24
  x -= 2;   // x = 22
  x *= 3;   // x = 66
  x /= 4;   // x = 16
  x %= 5;   // x = 1
  return x;
}
```
> Saída: `1`

Funciona também com variáveis locais dentro de funções.

---

### Extensão 4 — Valores Booleanos `true` / `false`

**Lexer:** novos tokens `TokenTrue` e `TokenFalse`.

**Parser:** `true` → nó `Const{Value: 1}`, `false` → nó `Const{Value: 0}`.

**Representação interna:** booleanos são inteiros — `0` é falso, qualquer outro valor é verdadeiro (igual à convenção de `if`/`while`).

**Exemplo:**
```fun
fun touch() {
  return true;
}

main {
  if true or touch() {
  } else {
  }
  return false;
}
```
> Saída: `0`

---

### Extensão 5 — Operadores Lógicos `and` / `or` / `not`

**Lexer:** novos tokens `TokenAnd`, `TokenOr`, `TokenNot`.

**Precedência** (do maior para menor):
```
or  >  and  >  not  >  comparações  >  +/-  >  */%
```

#### Short-circuit evaluation

`and` e `or` implementam avaliação em curto-circuito tanto no interpretador quanto no assembly gerado.

**Codegen para `and` com short-circuit:**
```asm
  ; avalia lado esquerdo
  cmp $0, %rax
  jz  Lfalse0       ; se false, pula avaliação do direito
  ; avalia lado direito
  cmp $0, %rax
  jz  Lfalse0
  jmp Ltrue0
Ltrue0:
  mov $1, %rax
  jmp Lend0
Lfalse0:
  mov $0, %rax
Lend0:
```

**Prova do short-circuit — `or` não avalia o lado direito se o esquerdo já é verdadeiro:**
```fun
var x = 0;

fun touch() {
  x = x + 1;   // efeito colateral: incrementa x
  return true;
}

main {
  if true or touch() { } else { }
  return x;    // x deve ser 0, pois touch() nunca foi chamada
}
```
> Saída: `0`

**`not` — nega o valor booleano:**
```fun
fun nonneg(x) {
  return not x < 0;
}

main {
  return nonneg(5) and not nonneg(0 - 1);
}
```
> Saída: `1`

---

### Combinando extensões

Todas as extensões podem ser usadas juntas:

```fun
var x = 10;

main {
  x %= 4;                        // x = 2
  if (x <= 2) and (x != 0) {    // true and true
    x += 3;                      // x = 5
  } else {
    x += 100;
  }
  return x >= 5;                 // 5 >= 5 → 1
}
```
> Saída: `1`

---

## 5. Testes

O projeto possui **36 testes automatizados** cobrindo todas as fases:

| Pacote                    | Testes | O que cobre                                     |
|---------------------------|--------|-------------------------------------------------|
| `internal/lexer`          | 6      | Keywords, operadores compostos, erros léxicos   |
| `internal/parser`         | 10     | Chamadas, booleanos, precedência, dessugarização|
| `internal/semantic`       | 5      | Erros de aridade, função undef, duplicatas      |
| `internal/interpreter`    | 13     | Comportamento de todas as extensões             |
| `internal/codegen`        | —      | Estrutura do assembly gerado                    |
| `internal/integration`    | 15     | Pipeline completo com fixtures reais            |

### Fixtures de teste

**Positivos** (devem compilar e produzir saída correta):

| Arquivo             | Testa                                  | Saída |
|---------------------|----------------------------------------|-------|
| `abs.fun`           | Função com variável local e if/else    | `11`  |
| `fact.fun`          | Recursão direta (fatorial)             | `120` |
| `chain.fun`         | Função chamando outra função           | `28`  |
| `noargs.fun`        | Função sem parâmetros                  | `42`  |
| `shadow.fun`        | Shadowing de variável global           | `142` |
| `mod.fun`           | Operador módulo                        | `1`   |
| `compare.fun`       | Comparadores `<=` `>=` `!=`           | `1`   |
| `compound.fun`      | Atribuições compostas                  | `1`   |
| `bool_if.fun`       | `true`/`false` em condicionais         | —     |
| `bool_while.fun`    | Booleanos em loops                     | —     |
| `bool_func.fun`     | `not` e `and` em funções               | `1`   |
| `bool_short_or.fun` | Short-circuit de `or`                  | `0`   |
| `combined.fun`      | Todas as extensões combinadas          | `1`   |

**Negativos** (devem falhar na análise semântica):

| Arquivo              | Erro esperado                              |
|----------------------|--------------------------------------------|
| `err_fun_undef.fun`  | Chamada a função não declarada             |
| `err_fun_arity.fun`  | Número errado de argumentos na chamada     |

### Executar os testes

```bash
go test ./...
```

---

## 6. Como Usar

### Compilar o compilador

```bash
go build -o ./bin/funcc ./cmd/funcc
```

### Compilar um programa Fun

```bash
./bin/funcc testdata/fun/positive/fact.fun
# gera: testdata/fun/positive/fact.s
```

Ou especificando saída:

```bash
./bin/funcc -o /tmp/meu_programa.s meu_programa.fun
```

### Montar e executar (Linux/WSL)

```bash
as -o prog.o prog.s
ld -o prog prog.o
./prog
echo $?    # saída do programa
```

### Script completo de testes (Linux/WSL)

```bash
sh scripts/run_tests.sh
```

---

## Resumo das Extensões

| Extensão                        | Complexidade | Status |
|---------------------------------|--------------|--------|
| Operador módulo `%`             | Simples      | ✓      |
| Comparadores `<=` `>=` `!=`     | Simples      | ✓      |
| Atribuições compostas `+=` etc. | Simples      | ✓      |
| Booleanos `true` / `false`      | Simples      | ✓      |
| Lógicos `and` / `or` / `not`    | Simples      | ✓      |

**Requisito do Projeto Final:** 1 extensão média/alta **ou** 3+ extensões simples.  
**Implementado:** 5 extensões simples — requisito cumprido com folga.
