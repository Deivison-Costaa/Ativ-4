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
	outputFile := flag.String("o", "", "Arquivo de saída (padrão: mesmo nome com extensão .s)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Uso: cmd [-o arquivo.s] <arquivo.ev>\n")
		fmt.Fprintf(os.Stderr, "Compilador Cmd - Gera código assembly x86-64\n\n")
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

	prog, err := parser.ParseProgram()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := CheckProgram(prog); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cg := NewCodeGenerator()
	code := cg.Generate(prog)

	if *outputFile != "" {
		err = os.WriteFile(*outputFile, []byte(code), 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}
	} else if path == "-" {
		fmt.Print(code)
	} else {
		outPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".s"
		err = os.WriteFile(outPath, []byte(code), 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
		}
	}
}
