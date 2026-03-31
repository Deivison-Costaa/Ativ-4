# Atividade 10 — Compilador Cmd (Comandos)

Compilador completo para a linguagem **Cmd**, gerando código **assembly x86-64** para Linux.

A linguagem Cmd estende a linguagem EV (atividades anteriores) adicionando:
- Bloco principal delimitado por `{` e `}` com `return`
- Comando condicional `if/else`
- Comando de repetição `while`
- Comando de atribuição (modifica variáveis já declaradas)
- Operadores de comparação: `<`, `>`, `==`

## Linguagem Cmd

### Gramática

```
<programa> ::= <decl>* '{' <cmd>* 'return' <exp> ';' '}'
<decl>     ::= <var> '=' <exp> ';'
<cmd>      ::= <if> | <while> | <atrib>
<if>       ::= 'if' <exp> '{' <cmd>* '}' 'else' '{' <cmd>* '}'
<while>    ::= 'while' <exp> '{' <cmd>* '}'
<atrib>    ::= <var> '=' <exp> ';'
<exp>      ::= <exp_a> (('<' | '>' | '==') <exp_a>)*
<exp_a>    ::= <exp_m> (('+' | '-') <exp_m>)*
<exp_m>    ::= <prim> (('*' | '/') <prim>)*
<prim>     ::= <num> | <var> | '(' <exp> ')'
```

### Exemplo — valor absoluto do discriminante

```
a = 1;
b = 2;
c = 3;
delta = b * b - 4 * a * c;
{
  if delta < 0 {
    delta = 0 - delta;
  } else {
    delta = delta;
  }
  return delta;
}
```

Saída: `8`

### Exemplo — soma com while

```
n = 1;
m = 10;
soma = 0;
{
  while n < m {
    soma = soma + n;
    n = n + 1;
  }
  return soma;
}
```

Saída: `45`

## Componentes

| Arquivo | Descrição |
|---------|-----------|
| `token.go` | Tipos de tokens (incluindo `{`, `}`, `<`, `>`, `==`, palavras-chave) |
| `lexer.go` | Analisador léxico |
| `ast.go` | AST: expressões, comandos (`IfCmd`, `WhileCmd`, `AtribCmd`), intérprete |
| `parser.go` | Analisador sintático (descendente recursivo) |
| `semantic.go` | Verificação semântica (variáveis declaradas) |
| `codegen.go` | Gerador de código assembly x86-64 |
| `main.go` | Programa principal |
| `runtime/runtime.s` | Funções auxiliares (`imprime_num`, `sair`) |
| `tests/` | Programas de teste |

## Como usar

### Compilar o compilador

```bash
go build -o cmd .
```

### Compilar um programa Cmd

```bash
# Gera arquivo .s no mesmo diretório do .ev
./cmd tests/delta.ev

# Especificar arquivo de saída
./cmd -o saida.s tests/delta.ev

# Ler da entrada padrão
echo '{ return 42; }' | ./cmd -o /tmp/test.s -
```

### Montar e executar (a partir do diretório raiz do projeto)

```bash
as -o tests/delta.o tests/delta.s
ld -o tests/delta tests/delta.o
./tests/delta
```

### Executar testes unitários

```bash
go test ./...
```

## Geração de Código

### Operadores de comparação

```asm
# A < B  (resultado: 1 se verdadeiro, 0 se falso)
<codigo_B>
push %rax
<codigo_A>
pop %rbx
xor %rcx, %rcx
cmp %rbx, %rax    # flags baseados em A - B
setl %cl          # setz para ==, setg para >
mov %rcx, %rax
```

### Comando if/else

```asm
<codigo_cond>
cmp $0, %rax
jz LfalsoN
<codigo_then>
jmp LfimN
LfalsoN:
<codigo_else>
LfimN:
```

### Comando while

```asm
LinicioN:
<codigo_cond>
cmp $0, %rax
jz LfimN
<codigo_corpo>
jmp LinicioN
LfimN:
```

## Programas de teste

| Arquivo | Descrição | Saída esperada |
|---------|-----------|----------------|
| `tests/input.ev` | Expressão com variáveis | `60467` |
| `tests/input2.ev` | Soma de variáveis | `90` |
| `tests/input3.ev` | Expressão simples | `22` |
| `tests/delta.ev` | Valor absoluto do discriminante | `8` |
| `tests/soma.ev` | Soma 1..9 com while | `45` |
| `tests/resto.ev` | Resto da divisão por subtração | `2` |
| `tests/mdc.ev` | Máximo divisor comum | `6` |
| `tests/err.ev` | Erro léxico | — |
| `tests/err1.ev` | Erro semântico: variável não declarada | — |
| `tests/err2.ev` | Erro semântico: variável não declarada | — |

## Requisitos

- Go 1.22+
- GNU Assembler (`as`)
- GNU Linker (`ld`)
- Linux x86-64
