# Atividade 04 — EC1 (Expressões Constantes 1) — Análise Léxica

Este projeto implementa um **analisador léxico** para a linguagem **EC1**.

A linguagem EC1 consiste em expressões aritméticas com operandos inteiros constantes e operadores `+`, `-`, `*`, `/`, com operações sempre entre parênteses.

## O que o programa faz

Dado um arquivo com um programa EC1, o analisador léxico:

- Ignora espaços em branco (` `, `\t`, `\n`, `\r`).
- Reconhece e gera tokens para:
  - `Numero` (sequência de um ou mais dígitos)
  - `ParenEsq` `(`
  - `ParenDir` `)`
  - `Soma` `+`
  - `Sub` `-`
  - `Mult` `*`
  - `Div` `/`
- Imprime cada token no formato:

```text
<Tipo, "lexema", posicao>
```

Onde `posicao` é o índice do caractere inicial do token na entrada.

## Erros léxicos

Se existir algum caractere fora do conjunto de dígitos, parênteses e operadores, o programa reporta erro no **stderr** e termina com código `1`:

```text
Erro lexico na posicao X
```

Obs.: a mensagem foi mantida **sem acentos** para evitar problemas de encoding no Windows/PowerShell.

## Como executar (Windows / PowerShell)

Na pasta `Ativ 4`:

### Executar lendo de um arquivo

```powershell
go run . .\tests\input.ec1
```

### Executar lendo do stdin

Use `-` como nome do arquivo:

```powershell
'(33 + (912 * 11))' | go run . -
```

## Exemplos prontos (script)

Existe um script que cria alguns arquivos de entrada e executa o programa, imprimindo a saída no terminal:

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File .\run_examples.ps1
```

## Testes

### Testes automatizados em Go

```powershell
go test ./...
```

Os testes verificam:

- Tokenização do exemplo do PDF
- Tokenização com diferentes espaços em branco
- Detecção de erro léxico e posição correta
