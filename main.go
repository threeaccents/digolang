package main

import (
	"bytes"
	"fmt"
	"github.com/threeaccents/digolang/eval"
	"github.com/threeaccents/digolang/lexer"
	"github.com/threeaccents/digolang/parser"
	"io"
	"os"
	"os/user"
	"strings"

	"github.com/threeaccents/digolang/repl"
)

func main() {
	if len(os.Args) == 3 {
		shouldRun := os.Args[1] == "run"
		if shouldRun {
			fileName := os.Args[2]
			if !isDigoFile(fileName) {
				fmt.Println("invalid file. File must be of type .digo")
				os.Exit(1)
			}

			f, err := os.Open(fileName)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			b := new(bytes.Buffer)

			if _, err := io.Copy(b, f); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Println("tokenizing file...")
			l := lexer.New(b.String())
			fmt.Println("parsing tokens...")
			p := parser.New(l)
			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(os.Stderr, p.Errors())
				os.Exit(1)
			}

			fmt.Println("evaluating program...")

			evaluated := eval.Eval(program)
			if evaluated != nil {
				io.WriteString(os.Stdout, evaluated.Inspect())
				io.WriteString(os.Stdout, "\n")
			}

			return
		}
	}

	u, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the digo programming language.\n ", u.Username)
	repl.Start(os.Stdin, os.Stdout)
}

func isDigoFile(name string) bool {
	nameSlice := strings.Split(name, ".")
	if len(nameSlice) == 1 {
		return false
	}

	return nameSlice[1] == "digo"
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! error executing program.\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
