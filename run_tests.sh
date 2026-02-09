#!/bin/bash
#
# Script de teste para o compilador EC1
#

EC1="./ec1"
TESTS_DIR="tests"

RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

echo "Compilando EC1..."
go build -o ec1 *.go || exit 1

run_test() {
    local input="$1"
    local expected="$2"
    local name="$3"

    echo "$input" | $EC1 -o /tmp/test_$$.s -
    as -o /tmp/test_$$.o /tmp/test_$$.s
    ld -o /tmp/test_$$ /tmp/test_$$.o
    result=$(/tmp/test_$$)
    rm -f /tmp/test_$$.s /tmp/test_$$.o /tmp/test_$$

    if [ "$result" = "$expected" ]; then
        echo -e "${GREEN}[PASS]${NC} $name: $input = $expected"
        return 0
    else
        echo -e "${RED}[FAIL]${NC} $name: $input"
        echo "  Esperado: $expected"
        echo "  Obtido:   $result"
        return 1
    fi
}

echo ""
echo "Executando testes..."
echo ""

PASSED=0
FAILED=0

# Testes
tests=(
    "42|42|Constante simples"
    "0|0|Zero"
    "333|333|Número maior"
    "(7 + 11)|18|Soma simples"
    "(11 - 7)|4|Subtração simples"
    "(6 * 7)|42|Multiplicação simples"
    "(20 / 4)|5|Divisão simples"
    "(3 + (4 + (11 + 7)))|25|Somas aninhadas"
    "(33 + (912 * 11))|10065|Soma e multiplicação"
    "((427 / 7) + (11 * (231 + 5)))|2657|Expressão complexa"
)

for test in "${tests[@]}"; do
    IFS='|' read -r input expected name <<< "$test"
    if run_test "$input" "$expected" "$name"; then
        ((PASSED++))
    else
        ((FAILED++))
    fi
done

# Testes com arquivos
if [ -f "$TESTS_DIR/input2.ec1" ]; then
    $EC1 "$TESTS_DIR/input2.ec1"
    as -o "$TESTS_DIR/input2.o" "$TESTS_DIR/input2.s"
    ld -o "$TESTS_DIR/input2" "$TESTS_DIR/input2.o"
    result=$("$TESTS_DIR/input2")
    if [ "$result" = "1065" ]; then
        echo -e "${GREEN}[PASS]${NC} input2.ec1 = 1065"
        ((PASSED++))
    else
        echo -e "${RED}[FAIL]${NC} input2.ec1: esperado 1065, obtido $result"
        ((FAILED++))
    fi
fi

if [ -f "$TESTS_DIR/input3.ec1" ]; then
    $EC1 "$TESTS_DIR/input3.ec1"
    as -o "$TESTS_DIR/input3.o" "$TESTS_DIR/input3.s"
    ld -o "$TESTS_DIR/input3" "$TESTS_DIR/input3.o"
    result=$("$TESTS_DIR/input3")
    if [ "$result" = "4120" ]; then
        echo -e "${GREEN}[PASS]${NC} input3.ec1 = 4120"
        ((PASSED++))
    else
        echo -e "${RED}[FAIL]${NC} input3.ec1: esperado 4120, obtido $result"
        ((FAILED++))
    fi
fi

echo ""
echo "================================"
echo "Resultados: $PASSED passou, $FAILED falhou"
echo "================================"

exit $FAILED
