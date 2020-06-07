package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/threeaccents/digolang/object"

	"github.com/threeaccents/digolang/eval"

	"github.com/threeaccents/digolang/parser"

	"github.com/threeaccents/digolang/lexer"
)

const prompt = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if line == "exit" {
			break
		}

		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluated := eval.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}

	fmt.Println("Goodbye =]")
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops! error executing program.\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
