#!/bin/sh
set -eu

BIN="./cmd"
TESTS_DIR="tests"

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo "Compilando o compilador Fun..."
go build -o "$BIN" .

echo ""
echo "Executando go test..."
go test ./...

echo ""
echo "Executando testes integrados de programas Fun..."

PASSED=0
FAILED=0

run_program_test() {
    file="$1"
    expected="$2"

    base="${file%.fun}"
    asm="${base}.s"
    obj="${base}.o"
    exe="${base}"

    "$BIN" "$file"
    as -o "$obj" "$asm"
    ld -o "$exe" "$obj"
    result=$("./$exe")

    if [ "$result" = "$expected" ]; then
        printf "${GREEN}[PASS]${NC} %s => %s\n" "$file" "$expected"
        PASSED=$((PASSED + 1))
    else
        printf "${RED}[FAIL]${NC} %s\n" "$file"
        echo "  esperado: $expected"
        echo "  obtido:   $result"
        FAILED=$((FAILED + 1))
    fi
}

run_error_test() {
    file="$1"
    expected="$2"

    if output=$("$BIN" "$file" 2>&1); then
        printf "${RED}[FAIL]${NC} %s\n" "$file"
        echo "  esperado erro contendo: $expected"
        echo "  compilacao concluiu sem erro"
        FAILED=$((FAILED + 1))
        return
    fi

    case "$output" in
        *"$expected"*)
        printf "${GREEN}[PASS]${NC} %s => erro esperado\n" "$file"
        PASSED=$((PASSED + 1))
        ;;
        *)
        printf "${RED}[FAIL]${NC} %s\n" "$file"
        echo "  esperado erro contendo: $expected"
        echo "  obtido: $output"
        FAILED=$((FAILED + 1))
        ;;
    esac
}

run_program_test "$TESTS_DIR/abs.fun" "11"
run_program_test "$TESTS_DIR/noargs.fun" "42"
run_program_test "$TESTS_DIR/chain.fun" "28"
run_program_test "$TESTS_DIR/shadow.fun" "142"
run_program_test "$TESTS_DIR/fact.fun" "120"

run_error_test "$TESTS_DIR/err_fun_undef.fun" "Erro semantico: funcao 'foo' nao declarada"
run_error_test "$TESTS_DIR/err_fun_arity.fun" "Erro semantico: funcao 'inc' esperava 1 parametros, recebeu 2"

echo ""
echo "================================"
echo "Resultados: $PASSED passou, $FAILED falhou"
echo "================================"

exit "$FAILED"
