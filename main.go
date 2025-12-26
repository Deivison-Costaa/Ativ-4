package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Uso: ec1lex <arquivo>")
		os.Exit(2)
	}

	path := os.Args[1]
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

	lex := NewLexer(string(data))
	for {
		tok, err := lex.NextToken()
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		if tok.Type == "" {
			break
		}
		fmt.Printf("<%s, \"%s\", %d>\n", tok.Type, tok.Lexeme, tok.Pos)
	}
}
