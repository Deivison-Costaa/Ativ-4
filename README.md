# Atividade 11 - Compilador Fun (Funções)

Compilador para a linguagem **Fun**, gerando código **assembly x86-64** para Linux.

A linguagem Fun estende a linguagem Cmd com:
- declarações globais com `var`
- bloco principal marcado por `main`
- declaração de funções com `fun`
- parâmetros formais e reais
- variáveis locais dentro de funções
- chamadas de função como expressões
- escopo local com precedência sobre variáveis globais

## Gramática

```txt
<programa> ::= <decl>* 'main' '{' <cmd>* 'return' <exp> ';' '}'
<decl>     ::= <vardecl> | <fundecl>
<fundecl>  ::= 'fun' <ident> '(' <arglist>? ')'
               '{' <vardecl>* <cmd>* 'return' <exp> ';' '}'
<arglist>  ::= <ident> | <ident> ',' <arglist>
<vardecl>  ::= 'var' <ident> '=' <exp> ';'
<cmd>      ::= <if> | <while> | <atrib>
<if>       ::= 'if' <exp> '{' <cmd>* '}' 'else' '{' <cmd>* '}'
<while>    ::= 'while' <exp> '{' <cmd>* '}'
<atrib>    ::= <ident> '=' <exp> ';'
<exp>      ::= <exp_a> (('<' | '>' | '==') <exp_a>)*
<exp_a>    ::= <exp_m> (('+' | '-') <exp_m>)*
<exp_m>    ::= <prim> (('*' | '/') <prim>)*
<prim>     ::= <num> | <ident> | '(' <exp> ')' | <fun>
<fun>      ::= <ident> '(' <params>? ')'
<params>   ::= <exp> | <exp> ',' <params>
```

## Componentes

| Arquivo | Descrição |
|---------|-----------|
| `token.go` | Tipos de tokens da linguagem Fun |
| `lexer.go` | Analisador léxico |
| `parser.go` | Parser descendente recursivo com lookahead para diferenciar variável/chamada |
| `ast.go` | AST, intérprete e suporte a escopos locais |
| `semantic.go` | Verificação semântica de variáveis, funções e aridade |
| `codegen.go` | Geração de assembly x86-64 com prólogo/epílogo de funções |
| `runtime/runtime.s` | Funções auxiliares `imprime_num` e `sair` |
| `tests/` | Programas de teste para Cmd e Fun |

## Como usar

### Compilar o compilador

```bash
go build -o cmd .
```

### Compilar um programa Fun

```bash
./cmd tests/abs.fun
./cmd -o tests/abs.s tests/abs.fun
echo 'main { return 42; }' | ./cmd -o /tmp/test.s -
```

### Montar e executar no Linux/WSL

```bash
as -o tests/abs.o tests/abs.s
ld -o tests/abs tests/abs.o
./tests/abs
```

### Executar testes unitários

```bash
go test ./...
```

### Executar testes integrados no WSL

```bash
sh run_tests.sh
```

O script `run_tests.sh` compila o compilador, gera assembly para programas `.fun`, monta, linka e executa os binários no ambiente Linux/WSL.

## Programas de teste Fun

| Arquivo | Descrição | Saída esperada |
|---------|-----------|----------------|
| `tests/abs.fun` | Função com variável local | `11` |
| `tests/noargs.fun` | Função sem parâmetros | `42` |
| `tests/chain.fun` | Função chamando outra função | `28` |
| `tests/shadow.fun` | Sombra de global por parâmetro/local | `142` |
| `tests/fact.fun` | Recursão direta (fatorial) | `120` |
| `tests/err_fun_undef.fun` | Erro semântico: função não declarada | — |
| `tests/err_fun_arity.fun` | Erro semântico: aridade incorreta | — |

## Convenção de chamada

- Argumentos são empilhados da direita para a esquerda.
- A chamada usa `call nome`.
- O chamador remove os argumentos da pilha após o retorno.
- O retorno da função é sempre colocado em `%rax`.
- Cada função usa `%rbp` como frame pointer.
- Variáveis locais e parâmetros são acessados com deslocamento relativo a `%rbp`.

## Requisitos

- Go 1.22+
- GNU Assembler (`as`)
- GNU Linker (`ld`)
- Linux x86-64 ou WSL para execução dos binários gerados
