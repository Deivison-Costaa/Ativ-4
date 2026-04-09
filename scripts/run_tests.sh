#!/bin/sh
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ROOT_DIR=$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)
POSITIVE_DIR="$ROOT_DIR/testdata/fun/positive"
NEGATIVE_DIR="$ROOT_DIR/testdata/fun/negative"
TMP_DIR=$(mktemp -d)
BIN="$TMP_DIR/funcc"

trap 'rm -rf "$TMP_DIR"' EXIT HUP INT TERM

GREEN=$(printf '\033[0;32m')
RED=$(printf '\033[0;31m')
NC=$(printf '\033[0m')

check_cmd() {
    if ! command -v "$1" >/dev/null 2>&1; then
        echo "Dependencia ausente: '$1'" >&2
        exit 2
    fi
}

expected_output() {
    case "$1" in
        abs) echo "11" ;;
        bool_func) echo "1" ;;
        bool_if) echo "1" ;;
        bool_short_or) echo "0" ;;
        bool_while) echo "3" ;;
        chain) echo "28" ;;
        combined) echo "1" ;;
        compare) echo "1" ;;
        compound) echo "1" ;;
        fact) echo "120" ;;
        mod) echo "1" ;;
        noargs) echo "42" ;;
        shadow) echo "142" ;;
        *)
            echo "Nao ha saida esperada cadastrada para '$1'" >&2
            exit 2
            ;;
    esac
}

expected_error() {
    case "$1" in
        err_fun_arity) echo "Erro semantico: funcao 'inc' esperava 1 parametros, recebeu 2" ;;
        err_fun_undef) echo "Erro semantico: funcao 'foo' nao declarada" ;;
        *)
            echo "Nao ha erro esperado cadastrado para '$1'" >&2
            exit 2
            ;;
    esac
}

check_cmd go
check_cmd as
check_cmd ld

case "$(uname -s)" in
    Linux*) ;;
    *)
        echo "Este script precisa ser executado em Linux ou WSL com toolchain GNU ELF (go, as e ld)." >&2
        exit 2
        ;;
esac

echo "Compilando o compilador Fun..."
cd "$ROOT_DIR"
go build -o "$BIN" ./cmd/funcc

echo ""
echo "Executando go test..."
go test ./...

echo ""
echo "Executando testes integrados com arquivos .fun..."

PASSED=0
FAILED=0

run_program_test() {
    file="$1"
    name=$(basename "$file" .fun)
    expected=$(expected_output "$name")
    asm="$TMP_DIR/$name.s"
    obj="$TMP_DIR/$name.o"
    exe="$TMP_DIR/$name"

    "$BIN" -o "$asm" "$file"
    (
        cd "$ROOT_DIR"
        as -o "$obj" "$asm"
        ld -o "$exe" "$obj"
    )
    result=$("$exe")

    if [ "$result" = "$expected" ]; then
        printf "%s[PASS]%s %s => %s\n" "$GREEN" "$NC" "$name" "$expected"
        PASSED=$((PASSED + 1))
    else
        printf "%s[FAIL]%s %s\n" "$RED" "$NC" "$name"
        echo "  esperado: $expected"
        echo "  obtido:   $result"
        FAILED=$((FAILED + 1))
    fi
}

run_error_test() {
    file="$1"
    name=$(basename "$file" .fun)
    expected=$(expected_error "$name")
    asm="$TMP_DIR/$name.s"

    if output=$("$BIN" -o "$asm" "$file" 2>&1); then
        printf "%s[FAIL]%s %s\n" "$RED" "$NC" "$name"
        echo "  esperado erro contendo: $expected"
        echo "  compilacao concluiu sem erro"
        FAILED=$((FAILED + 1))
        return
    fi

    case "$output" in
        *"$expected"*)
            printf "%s[PASS]%s %s => erro esperado\n" "$GREEN" "$NC" "$name"
            PASSED=$((PASSED + 1))
            ;;
        *)
            printf "%s[FAIL]%s %s\n" "$RED" "$NC" "$name"
            echo "  esperado erro contendo: $expected"
            echo "  obtido: $output"
            FAILED=$((FAILED + 1))
            ;;
    esac
}

for file in "$POSITIVE_DIR"/*.fun; do
    run_program_test "$file"
done

for file in "$NEGATIVE_DIR"/*.fun; do
    run_error_test "$file"
done

echo ""
echo "================================"
echo "Resultados: $PASSED passou, $FAILED falhou"
echo "================================"

exit "$FAILED"
