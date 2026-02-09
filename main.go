package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	outputFile := flag.String("o", "", "Arquivo de saída (padrão: stdout)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: ec1 [-o arquivo.s] <arquivo.ec1>\n")
		fmt.Fprintf(os.Stderr, "Compilador EC1 - Gera código assembly x86-64\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	path := flag.Arg(0)
	var data []byte
	var err error
	if path == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	parser, err := NewParser(NewLexer(string(data)))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	ast, err := parser.ParseProgram()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Gera código assembly
	cg := NewCodeGenerator()
	code := cg.Generate(ast)

	// Escreve a saída
	if *outputFile != "" {
		err = os.WriteFile(*outputFile, []byte(code), 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}
	} else if path != "-" && !strings.HasSuffix(path, ".ec1") {
		// Se não especificou -o e o arquivo não termina em .ec1, imprime em stdout
		fmt.Print(code)
	} else if path != "-" {
		// Gera nome do arquivo .s baseado no arquivo de entrada
		outPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".s"
		err = os.WriteFile(outPath, []byte(code), 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}
	} else {
		fmt.Print(code)
	}
}
