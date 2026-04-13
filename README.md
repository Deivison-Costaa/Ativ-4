# Compilador Fun

Compilador para a linguagem **Fun**, gerando código **assembly x86-64** para Linux.

Extensões atuais:
- booleanos com `true`, `false`, `and`, `or` e `not`, representados em tempo de execução como `0/1`
- resto de divisão com `%`
- comparadores `<=`, `>=` e `!=`
- atribuição composta com `+=`, `-=`, `*=`, `/=` e `%=`

## Estrutura

| Caminho | Descrição |
|---------|-----------|
| `cmd/funcc` | CLI do compilador |
| `internal/fun` | Tokens, AST e tipos centrais da linguagem |
| `internal/lexer` | Analisador léxico |
| `internal/parser` | Parser descendente recursivo |
| `internal/semantic` | Verificação semântica |
| `internal/interpreter` | Intérprete usado pelos testes |
| `internal/codegen` | Geração de assembly x86-64 |
| `internal/integration` | Testes de pipeline usando fixtures |
| `runtime/runtime.s` | Funções auxiliares `imprime_num` e `sair` |
| `testdata/fun` | Fixtures positivos e negativos da linguagem Fun |
| `testdata/legacy/cmd` | Casos `.ev` antigos preservados como referência |
| `scripts/` | Scripts de apoio para exemplos e testes integrados |
| `docs/relatorio.txt` | Relatório do projeto |

## Como usar

### Compilar o compilador

```bash
go build -o ./bin/funcc ./cmd/funcc
```

### Gerar assembly para um programa Fun

```bash
./bin/funcc testdata/fun/positive/abs.fun
./bin/funcc -o /tmp/abs.s testdata/fun/positive/abs.fun
echo 'main { return 42; }' | ./bin/funcc -o /tmp/test.s -
```

### Executar testes unitários

```bash
go test ./...
```

### Executar testes integrados no Linux/WSL

```bash
sh scripts/run_tests.sh
```

O script compila a CLI em um diretório temporário, valida `go test ./...`, gera assembly para os fixtures `.fun`, monta, linka e executa os binários sem sujar arquivos versionados.

## Operadores adicionais

- `true` gera valor `1`
- `false` gera valor `0`
- `not x` retorna `1` quando `x` é falso
- `a and b` e `a or b` fazem short-circuit
- `if` e `while` continuam usando a convenção `0 = falso`, `!= 0 = verdadeiro`
- `a % b` calcula o resto da divisão inteira
- `a <= b`, `a >= b` e `a != b` retornam `0` ou `1`
- `x += y`, `x -= y`, `x *= y`, `x /= y` e `x %= y` são dessugarizados para `x = x op y`

## Fixtures Fun

### Positivos

| Arquivo | Descrição | Saída esperada |
|---------|-----------|----------------|
| `testdata/fun/positive/abs.fun` | Função com variável local | `11` |
| `testdata/fun/positive/bool_func.fun` | Função retornando condição booleana | `1` |
| `testdata/fun/positive/bool_if.fun` | Uso de `if` com `and/not` | `1` |
| `testdata/fun/positive/bool_short_or.fun` | Short-circuit de `or` com efeito colateral | `0` |
| `testdata/fun/positive/bool_while.fun` | Uso de `while` com `or/not` | `3` |
| `testdata/fun/positive/noargs.fun` | Função sem parâmetros | `42` |
| `testdata/fun/positive/chain.fun` | Função chamando outra função | `28` |
| `testdata/fun/positive/mod.fun` | Uso do operador de resto `%` | `1` |
| `testdata/fun/positive/compare.fun` | Comparadores `<=`, `>=`, `!=` | `1` |
| `testdata/fun/positive/compound.fun` | Atribuições compostas em globais e locais | `1` |
| `testdata/fun/positive/combined.fun` | Combinação de booleanos, `%`, comparações e atribuição composta | `1` |
| `testdata/fun/positive/shadow.fun` | Sombra de global por parâmetro/local | `142` |
| `testdata/fun/positive/fact.fun` | Recursão direta (fatorial) | `120` |

### Negativos

| Arquivo | Descrição |
|---------|-----------|
| `testdata/fun/negative/err_fun_undef.fun` | Erro semântico: função não declarada |
| `testdata/fun/negative/err_fun_arity.fun` | Erro semântico: aridade incorreta |

## Requisitos

- Go 1.22+
- GNU Assembler (`as`)
- GNU Linker (`ld`)
- Linux x86-64 ou WSL para executar os binários gerados
