package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"ativ10/internal/codegen"
	"ativ10/internal/lexer"
	"ativ10/internal/parser"
	"ativ10/internal/semantic"
)

func main() {
	outputFile := flag.String("o", "", "Arquivo de saída (padrão: mesmo nome com extensão .s)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: funcc [-o arquivo.s] <arquivo.fun|->\n")
		fmt.Fprintf(os.Stderr, "Compilador Fun - Gera codigo assembly x86-64\n\n")
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

	p, err := parser.NewParser(lexer.NewLexer(string(data)))
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	prog, err := p.ParseProgram()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := semantic.CheckProgram(prog); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cg := codegen.NewCodeGenerator()
	code := cg.Generate(prog)

	if *outputFile != "" {
		err = os.WriteFile(*outputFile, []byte(code), 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}
		return
	}

	if path == "-" {
		fmt.Print(code)
		return
	}

	outPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".s"
	err = os.WriteFile(outPath, []byte(code), 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}
